package rfq

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/TBD54566975/tbdex-go/tbdex/closemsg"
	"github.com/TBD54566975/tbdex-go/tbdex/crypto"
	"github.com/TBD54566975/tbdex-go/tbdex/message"
	_offering "github.com/TBD54566975/tbdex-go/tbdex/offering"
	"github.com/TBD54566975/tbdex-go/tbdex/quote"
	"github.com/TBD54566975/tbdex-go/tbdex/validator"
	"github.com/shopspring/decimal"
	"github.com/tbd54566975/web5-go/pexv2"
	"github.com/tbd54566975/web5-go/vc"

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

// VerifyOfferingRequirements verifies that the RFQ meets the requirements of the provided offering.
// Specifically this includes the following checks:
//   - payin method is present in the offering
//   - payin amount is within the offering or payin method's min and max
//   - payin details satisfy the offering's required payment details
//   - payout method is present in the offering
//   - payout details satisfy the offering's required payment details
//   - claims satisfy the offering's required claims
func (rfq *RFQ) VerifyOfferingRequirements(offering _offering.Offering) error {
	if rfq.Data.OfferingID != offering.Metadata.ID {
		return fmt.Errorf("rfq's offering id does not match offering used to evaluate rfq")
	}

	payinAmount, err := decimal.NewFromString(rfq.Data.Payin.Amount)
	if err != nil {
		return fmt.Errorf("failed to parse rfq payin amount: %w", err)
	}

	var selectedPayinMethod *_offering.PayinMethod
	for _, method := range offering.Data.Payin.Methods {

		if method.Kind == rfq.Data.Payin.Kind {
			selectedPayinMethod = &method
			break
		}
	}

	if selectedPayinMethod == nil {
		return errors.New("rfq payin method not found in offering")
	}

	var min *decimal.Decimal
	maybeMins := []string{selectedPayinMethod.Min, offering.Data.Payin.Min}
	for _, maybeMin := range maybeMins {
		if maybeMin != "" {
			minAmount, err := decimal.NewFromString(maybeMin)
			if err != nil {
				return fmt.Errorf("failed to parse min amount: %w", err)
			}
			min = &minAmount
			break
		}
	}

	if min != nil {
		if payinAmount.LessThan(*min) {
			return fmt.Errorf("rfq payin amount is less than offering's minimum amount")
		}
	}

	var max *decimal.Decimal
	maybeMaxes := []string{selectedPayinMethod.Max, offering.Data.Payin.Max}

	for _, maybeMax := range maybeMaxes {
		if maybeMax != "" {
			maxAmount, err := decimal.NewFromString(maybeMax)
			if err != nil {
				return fmt.Errorf("failed to parse max amount: %w", err)
			}
			max = &maxAmount
			break
		}
	}

	if max != nil {
		if payinAmount.GreaterThan(*max) {
			return fmt.Errorf("rfq payin amount is greater than offering's max amount")
		}
	}

	if selectedPayinMethod.RequiredPaymentDetails == nil && rfq.PrivateData != nil && rfq.PrivateData.Payin.PaymentDetails != nil {
		return errors.New("rfq contains unexpected payin details")
	}

	if selectedPayinMethod.RequiredPaymentDetails != nil {
		if rfq.PrivateData == nil || rfq.PrivateData.Payin.PaymentDetails == nil {
			return errors.New("rfq does not contain expected payin details")
		}

		payinDetails, err := json.Marshal(rfq.PrivateData.Payin.PaymentDetails)
		if err != nil {
			return fmt.Errorf("failed to json marshal rfq payin details: %w", err)
		}

		schema, err := jsonschema.CompileString("payin", string(selectedPayinMethod.RequiredPaymentDetails))
		if err != nil {
			return fmt.Errorf("failed to compile offering's required payment details")
		}

		if err := schema.Validate(payinDetails); err != nil {
			return fmt.Errorf("rfq payin details do not satisfy offering's requirements: %w", err)
		}
	}

	var selectedPayoutMethod *_offering.PayoutMethod
	for _, method := range offering.Data.Payout.Methods {

		if method.Kind == rfq.Data.Payout.Kind {
			selectedPayoutMethod = &method
			break
		}
	}

	if selectedPayoutMethod == nil {
		return errors.New("rfq's selected payout method is not present in offering")
	}

	if selectedPayoutMethod.RequiredPaymentDetails == nil && rfq.PrivateData != nil && rfq.PrivateData.Payout.PaymentDetails != nil {
		return errors.New("rfq contains unexpected payout details")
	}

	if selectedPayoutMethod.RequiredPaymentDetails != nil {
		if rfq.PrivateData == nil || rfq.PrivateData.Payout.PaymentDetails == nil {
			return errors.New("rfq does not contain expected payout details")
		}

		payinDetails, err := json.Marshal(rfq.PrivateData.Payout.PaymentDetails)
		if err != nil {
			return fmt.Errorf("failed to json marshal rfq payout details: %w", err)
		}

		schema, err := jsonschema.CompileString("payout", string(selectedPayoutMethod.RequiredPaymentDetails))
		if err != nil {
			return fmt.Errorf("failed to compile offering's required payment details")
		}

		if err := schema.Validate(payinDetails); err != nil {
			return fmt.Errorf("rfq payout details do not satisfy offering's requirements: %w", err)
		}
	}

	if offering.Data.RequiredClaims != nil {
		if err := rfq.verifyClaims(offering.Data.RequiredClaims); err != nil {
			return fmt.Errorf("rfq claims do not satisfy offering's requirements")
		}
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

func (r *RFQ) verifyClaims(requiredClaims *pexv2.PresentationDefinition) error {
	if requiredClaims == nil {
		return errors.New("required claims cannot be nil")
	}

	if r.PrivateData == nil || r.PrivateData.Claims == nil {
		return errors.New("rfq claims is nil")
	}

	credentials, err := pexv2.SelectCredentials(r.PrivateData.Claims, *requiredClaims)

	if err != nil {
		return fmt.Errorf("failed to select credentials: %w", err)
	}

	if len(credentials) == 0 {
		return errors.New("claims do not fulfill the offering's requirements")
	}

	for _, cred := range credentials {
		_, err = vc.Verify[vc.Claims](cred)

		if err != nil {
			return fmt.Errorf("failed to verify credential: %w", err)
		}
	}

	return nil
}

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
