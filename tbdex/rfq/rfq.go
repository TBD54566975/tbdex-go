package rfq

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/TBD54566975/tbdex-go/tbdex/closemsg"
	"github.com/TBD54566975/tbdex-go/tbdex/crypto"
	"github.com/TBD54566975/tbdex-go/tbdex/message"
	"github.com/TBD54566975/tbdex-go/tbdex/offering"
	"github.com/TBD54566975/tbdex-go/tbdex/quote"
	"github.com/TBD54566975/tbdex-go/tbdex/validator"
	"github.com/tbd54566975/web5-go/pexv2"
	"github.com/tbd54566975/web5-go/vc"

	// jsonschema "github.com/xeipuuv/gojsonschema" // todo remove this one
	jsonschema "github.com/santhosh-tekuri/jsonschema/v5"
)

// Kind identifies this message kind
const Kind = "rfq"

// ValidNext returns the valid message kinds that can follow a RFQ.
func ValidNext() []string {
	return []string{quote.Kind, closemsg.Kind}
}

// RFQ represents a request for quote message within the exchange.
type RFQ struct {
	Metadata    message.Metadata `json:"metadata,omitempty"`
	Data        Data             `json:"data,omitempty"`
	PrivateData *PrivateData     `json:"privateData,omitempty"`
	Signature   string           `json:"signature,omitempty"`
}

// GetMetadata returns the metadata of the message
func (r RFQ) GetMetadata() message.Metadata {
	return r.Metadata
}

// GetKind returns the kind of message
func (r RFQ) GetKind() string {
	return Kind
}

// GetValidNext returns the valid message kinds that can follow a RFQ.
func (r RFQ) GetValidNext() []string {
	return ValidNext()
}

// UnmarshalJSON validates and unmarshals the input data into an RFQ.
func (r *RFQ) UnmarshalJSON(data []byte) error {
	err := validator.Validate(validator.TypeMessage, data, validator.WithKind(Kind))
	if err != nil {
		return fmt.Errorf("invalid rfq: %w", err)
	}

	ret := rfq{}
	err = json.Unmarshal(data, &ret)
	if err != nil {
		return fmt.Errorf("failed to JSON unmarshal rfq: %w", err)
	}

	*r = RFQ(ret)

	return nil
}

// Verify verifies the signature of the RFQ.
func (r *RFQ) Verify() error {
	decoded, err := crypto.VerifySignature(r, r.Signature)
	if err != nil {
		return fmt.Errorf("failed to verify RFQ signature: %w", err)
	}

	if decoded.SignerDID.URI != r.Metadata.From {
		return fmt.Errorf("signer: %s does not match message metadata from: %s", decoded.SignerDID.URI, r.Metadata.From)
	}

	return nil
}

// Scrub verifies the private data and returns an RFQ without private data for storage, as well as private data for separate processing
// todo: allow passing in custom type for PrivateData.Payin and PrivateData.Payout when PrivateData is genericized. https://github.com/TBD54566975/tbdex-go/issues/50
func (r *RFQ) Scrub() (RFQ, PrivateData, error) {

	err := r.verifyPrivateData()
	if err != nil {
		return RFQ{}, PrivateData{}, fmt.Errorf("failed to verify private data: %w", err)
	}

	privateData := *r.PrivateData

	// copy rfq
	r.PrivateData = nil

	scrubbed := RFQ{
		Metadata:  r.Metadata,
		Data:      r.Data,
		Signature: r.Signature,
	}

	return scrubbed, privateData, nil

}

// Parse validates, parses input data into an RFQ, and verifies the signature and private data.
func Parse(data []byte) (RFQ, error) {
	r := RFQ{}
	if err := json.Unmarshal(data, &r); err != nil {
		return RFQ{}, fmt.Errorf("failed to unmarshal RFQ: %w", err)
	}

	if err := r.Verify(); err != nil {
		return RFQ{}, fmt.Errorf("failed to verify RFQ: %w", err)
	}

	return r, nil
}

// Digest computes a hash of the rfq
func (r RFQ) Digest() ([]byte, error) {
	payload := map[string]any{"metadata": r.Metadata, "data": r.Data}

	hashed, err := crypto.DigestJSON(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to digest rfq: %w", err)
	}

	return hashed, nil
}

func (r *RFQ) verifyPrivateData() error {
	if r.PrivateData == nil {
		return errors.New("private data is missing")
	}

	if r.Data.ClaimsHash != "" {
		if len(r.PrivateData.Claims) == 0 {
			return errors.New("verification: claims hash is set but claims are missing")
		}
		payload := []any{r.PrivateData.Salt, r.PrivateData.Claims}
		if err := crypto.VerifyDigest(r.Data.ClaimsHash, payload); err != nil {
			return fmt.Errorf("failed to verify claims: %w", err)
		}
	}

	if r.Data.Payin.PaymentDetailsHash != "" {
		if r.PrivateData.Payin.PaymentDetails == nil {
			return errors.New("verification: payin details hash is set but payin details are missing")

		}
		payload := []any{r.PrivateData.Salt, r.PrivateData.Payin.PaymentDetails}
		if err := crypto.VerifyDigest(r.Data.Payin.PaymentDetailsHash, payload); err != nil {
			return fmt.Errorf("failed to verify payin: %w", err)
		}
	}

	if r.Data.Payout.PaymentDetailsHash != "" {
		if r.PrivateData.Payout.PaymentDetails == nil {
			return errors.New("verification: payout details hash is set but payout details are missing")
		}
		payload := []any{r.PrivateData.Salt, r.PrivateData.Payout.PaymentDetails}
		if err := crypto.VerifyDigest(r.Data.Payout.PaymentDetailsHash, payload); err != nil {
			return fmt.Errorf("failed to verify payout: %w", err)
		}
	}

	return nil
}

func (r *RFQ) VerifyOfferingRequirements(offering offering.Offering) error {
	if offering.Metadata.Protocol != r.Metadata.Protocol {
		return fmt.Errorf("protocol version mismatch. (rfq) %s != %s (offering)", r.Metadata.Protocol, offering.Metadata.Protocol)
	}

	if offering.Metadata.ID != r.Data.OfferingID {
		return fmt.Errorf("offering ID mismatch. (rfq) %s != %s (offering)", r.Data.OfferingID, offering.Metadata.ID)
	}

	// Verifying payin amount is less than maximum
	payinAmount, err := strconv.ParseFloat(r.Data.Payin.Amount, 64)
	if err != nil {
		return fmt.Errorf("failed to parse payin amount: %w", err)
	}

	if offering.Data.Payin.Max != "" {
		maxAmount, err := strconv.ParseFloat(offering.Data.Payin.Max, 64)
		if err != nil {
			return fmt.Errorf("failed to parse offering payin max amount: %w", err)
		}

		if payinAmount > maxAmount {
			return fmt.Errorf("payin amount exceeds maximum. (rfq) %f > %f (offering)", payinAmount, maxAmount)
		}
	}

	// Verify payin amount is greater than minimum
	if offering.Data.Payin.Min != "" {
		minAmount, err := strconv.ParseFloat(offering.Data.Payin.Min, 64)
		if err != nil {
			return fmt.Errorf("failed to parse offering payin min amount: %w", err)
		}

		if payinAmount < minAmount {
			return fmt.Errorf("rfq payin amount is below minimum. (rfq) %s < %s (offering)", r.Data.Payin.Amount, offering.Data.Payin.Min)
		}
	}

	// Verify payin method
	err = r.verifyPaymentMethod(
		r.Data.Payin.Kind,
		r.Data.Payin.PaymentDetailsHash,
		r.PrivateData,
		convertPayinMethods(offering.Data.Payin.Methods),
		"payin",
	)
	if err != nil {
		return fmt.Errorf("failed to verify payin method: %w", err)
	}

	// Verify payout method
	err = r.verifyPaymentMethod(
		r.Data.Payout.Kind,
		r.Data.Payout.PaymentDetailsHash,
		r.PrivateData,
		convertPayoutMethods(offering.Data.Payout.Methods),
		"payout",
	)

	if err != nil {
		return fmt.Errorf("failed to verify payout method: %w", err)
	}

	err = r.verifyClaims(offering.Data.RequiredClaims)
	if err != nil {
		return fmt.Errorf("failed to verify claims: %w", err)
	}

	return nil

}

func convertPayinMethods(payinMethods []offering.PayinMethod) []offering.PaymentMethod {
	paymentMethods := make([]offering.PaymentMethod, len(payinMethods))
	for i, method := range payinMethods {
		paymentMethods[i] = method
	}
	return paymentMethods
}

func convertPayoutMethods(payoutMethods []offering.PayoutMethod) []offering.PaymentMethod {
	paymentMethods := make([]offering.PaymentMethod, len(payoutMethods))
	for i, method := range payoutMethods {
		paymentMethods[i] = method
	}
	return paymentMethods
}

func (r *RFQ) verifyPaymentMethod(
	selectedPaymentKind string,
	selectedPaymentDetailsHash string,
	privateData *PrivateData,
	offeringPaymentMethods []offering.PaymentMethod,
	payDirection PayDirection,
) error {
	// Filter allowed payment methods by selectedPaymentKind
	var allowedPaymentMethods []offering.PaymentMethod
	for _, paymentMethod := range offeringPaymentMethods {
		if paymentMethod.GetKind() == selectedPaymentKind {
			allowedPaymentMethods = append(allowedPaymentMethods, paymentMethod)
		}
	}

	// If no matching payment methods are found, throw an error
	if len(allowedPaymentMethods) == 0 {
		var allowedPaymentMethodKinds []string
		for _, paymentMethod := range offeringPaymentMethods {
			allowedPaymentMethodKinds = append(allowedPaymentMethodKinds, paymentMethod.GetKind())
		}
		return fmt.Errorf(
			"offering does not support rfq's %sMethod kind. (rfq) %s was not found in: [%s] (offering)",
			payDirection, selectedPaymentKind, strings.Join(allowedPaymentMethodKinds, ", "),
		)
	}

	for _, paymentMethod := range allowedPaymentMethods {
		if paymentMethod.GetRequiredPaymentDetails() == nil {
			// If requiredPaymentDetails is omitted, and selectedPaymentDetailsHash is also omitted, we have a match
			if selectedPaymentDetailsHash == "" {
				return nil
			}
			// selectedPaymentDetailsHash is present even though requiredPaymentDetails is omitted.
			return fmt.Errorf("rfq paymentDetails must be omitted when offering requiredPaymentDetails is omitted")
		}

		// todo do we actually want to return error in this case where requiredPaymentDetails is present but privateData is nil?
		if privateData == nil {
			return fmt.Errorf("offering requiredPaymentDetails is present but rfq private data is omitted")
		}

		// requiredPaymentDetails is present and privateData is present, so Rfq's payment details must match
		if payDirection == PayIn {
			err := validatePaymentDetailSchema(paymentMethod, privateData.Payin.PaymentDetails, payDirection)
			if err != nil {
				return err
			}
		}

		if payDirection == PayOut {
			err := validatePaymentDetailSchema(paymentMethod, privateData.Payout.PaymentDetails, payDirection)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func validatePaymentDetailSchema(paymentMethod offering.PaymentMethod, paymentDetails map[string]any, payDirection PayDirection) error {
	if paymentDetails == nil {
		return fmt.Errorf("rfq %s paymentDetails are missing", payDirection)
	}

	sch, err := json.Marshal(paymentMethod.GetRequiredPaymentDetails())
	if err != nil {
		return fmt.Errorf("failed to marshal requiredPaymentDetails: %w", err)
	}

	// todo idk what the id should be?
	schema, err := jsonschema.CompileString("schema.json", string(sch))
	if err != nil {
		return fmt.Errorf("failed to compile requiredPaymentDetails schema: %w", err)
	}

	err = schema.Validate(paymentDetails)
	if err != nil {
		return fmt.Errorf("failed to validate %sMethod's paymentDetails: %w", payDirection, err)
	}
	return nil
}

func (r *RFQ) verifyClaims(requiredClaims *pexv2.PresentationDefinition) error {
	if requiredClaims == nil {
		return nil
	}

	credentials, err := pexv2.SelectCredentials(r.PrivateData.Claims, *requiredClaims)

	if err != nil {
		return fmt.Errorf("failed to select credentials: %w", err)
	}

	if len(credentials) == 0 {
		return fmt.Errorf("claims do not fulfill the offering's requirements")
	}

	for _, cred := range credentials {
		_, err = vc.Verify[vc.Claims](cred)

		if err != nil {
			return fmt.Errorf("failed to verify credential: %w", err)
		}
	}

	return nil
}

type PayDirection string

const (
	PayIn  PayDirection = "payin"
	PayOut PayDirection = "payout"
)

// Data encapsulates the data content of a request for quote.
type Data struct {
	OfferingID string               `json:"offeringId"`
	Payin      ScrubbedPayinMethod  `json:"payin"`
	Payout     ScrubbedPayoutMethod `json:"payout"`
	ClaimsHash string               `json:"claimsHash,omitempty"`
}

// PrivateData contains data which can be detached from the payload without disrupting integrity.
type PrivateData struct {
	Salt   string                `json:"salt,omitempty"`
	Claims []string              `json:"claims,omitempty"`
	Payin  PrivatePaymentDetails `json:"payin,omitempty"`
	Payout PrivatePaymentDetails `json:"payout,omitempty"`
}

// IsZero checks if struct is empty
func (p PrivateData) IsZero() bool {
	v := reflect.ValueOf(p)
	return v.IsZero()
}

// PrivatePaymentDetails contains the private payment details. used in [PrivateData]
type PrivatePaymentDetails struct {
	PaymentDetails PaymentMethodDetails `json:"paymentDetails,omitempty"`
}

// PayinMethod is used to create the payin method for an RFQ
type PayinMethod struct {
	Amount         string               `json:"amount"`
	Kind           string               `json:"kind"`
	PaymentDetails PaymentMethodDetails `json:"paymentDetails"`
}

// Scrub extracts the private data from the payin method and replaces it with a hash
func (p *PayinMethod) Scrub(salt string, privateData *PrivateData) (ScrubbedPayinMethod, error) {
	scrubbedPayin := ScrubbedPayinMethod{Amount: p.Amount, Kind: p.Kind}

	if p.PaymentDetails == nil {
		return scrubbedPayin, nil
	}

	hash, err := computeHash(salt, p.PaymentDetails)
	if err != nil {
		return ScrubbedPayinMethod{}, fmt.Errorf("failed to compute hash: %w", err)
	}

	scrubbedPayin.PaymentDetailsHash = hash
	privateData.Payin = PrivatePaymentDetails{PaymentDetails: p.PaymentDetails}

	return scrubbedPayin, nil
}

// PayoutMethod is used to create the payout method for an RFQ
type PayoutMethod struct {
	Kind           string               `json:"kind"`
	PaymentDetails PaymentMethodDetails `json:"paymentDetails"`
}

// Scrub extracts the private data from the payout method and replaces it with a hash
func (p *PayoutMethod) Scrub(salt string, privateData *PrivateData) (ScrubbedPayoutMethod, error) {
	scrubbedPayout := ScrubbedPayoutMethod{Kind: p.Kind}
	if p.PaymentDetails == nil {
		return scrubbedPayout, nil
	}

	hash, err := computeHash(salt, p.PaymentDetails)
	if err != nil {
		return ScrubbedPayoutMethod{}, fmt.Errorf("failed to compute hash: %w", err)
	}

	scrubbedPayout.PaymentDetailsHash = hash
	privateData.Payout = PrivatePaymentDetails{PaymentDetails: p.PaymentDetails}

	return scrubbedPayout, nil
}

// ClaimsSet is a set of claims
type ClaimsSet []string

// Scrub extracts claims from the payin method and replaces it with a hash
func (c ClaimsSet) Scrub(salt string, privateData *PrivateData) (string, error) {
	if len(c) == 0 {
		return "", nil
	}

	scrubbedClaims, err := computeHash(salt, c)
	if err != nil {
		return "", err
	}

	privateData.Claims = c

	return scrubbedClaims, nil
}

// ScrubbedPayinMethod represents the chosen method for the pay-in
type ScrubbedPayinMethod struct {
	Amount             string `json:"amount"`
	Kind               string `json:"kind"`
	PaymentDetailsHash string `json:"paymentDetailsHash,omitempty"`
}

// ScrubbedPayoutMethod represents the chosen method for the pay-out
type ScrubbedPayoutMethod struct {
	Kind               string `json:"kind"`
	PaymentDetailsHash string `json:"paymentDetailsHash,omitempty"`
}

// PaymentMethodDetails is a map populated with the required payment details specified by an offering
type PaymentMethodDetails map[string]any

type rfq RFQ
