package rfq

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/TBD54566975/tbdex-go/tbdex"
)

// Kind identifies this message kind
const Kind = "rfq"

// RFQ represents a request for quote message within the exchange.
type RFQ struct {
	Metadata    tbdex.MessageMetadata `json:"metadata,omitempty"`
	Data        Data                  `json:"data,omitempty"`
	PrivateData *PrivateData          `json:"privateData,omitempty"`
	Signature   string                `json:"signature,omitempty"`
}

// UnmarshalJSON validates and unmarshals the input data into an RFQ.
func (r *RFQ) UnmarshalJSON(data []byte) error {
	err := tbdex.Validate(tbdex.TypeMessage, data, tbdex.WithKind(Kind))
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

// Verify verifies the signature and private data hashes of the RFQ.
func (r *RFQ) Verify(privateDataStrict bool) error {
	decoded, err := tbdex.VerifySignature(r, r.Signature)
	if err != nil {
		return fmt.Errorf("failed to verify RFQ signature: %w", err)
	}

	if decoded.SignerDID.URI != r.Metadata.From {
		return fmt.Errorf("signer: %s does not match message metadata from: %s", decoded.SignerDID.URI, r.Metadata.From)
	}

	err = r.verifyPrivateData(privateDataStrict)
	if err != nil {
		return fmt.Errorf("failed to verify private data: %w", err)
	}

	return nil
}

// Parse validates, parses input data into an RFQ, and verifies the signature and private data.
func Parse(data []byte, privateDataStrict bool) (RFQ, error) {
	r := RFQ{}
	if err := json.Unmarshal(data, &r); err != nil {
		return RFQ{}, fmt.Errorf("failed to unmarshal RFQ: %w", err)
	}

	if err := r.Verify(privateDataStrict); err != nil {
		return RFQ{}, fmt.Errorf("failed to verify RFQ: %w", err)
	}

	return r, nil
}

// Digest computes a hash of the rfq
func (r RFQ) Digest() ([]byte, error) {
	payload := map[string]any{"metadata": r.Metadata, "data": r.Data}

	hashed, err := tbdex.DigestJSON(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to digest rfq: %w", err)
	}

	return hashed, nil
}

func (r *RFQ) verifyPrivateData(strict bool) error {
	if !strict && r.PrivateData == nil {
		return nil
	}

	if r.Data.ClaimsHash != "" {
		if strict && len(r.PrivateData.Claims) == 0 {
			return errors.New("strict verification: claims hash is set but claims are missing")

		}
		if len(r.PrivateData.Claims) != 0 {
			payload := []any{r.PrivateData.Salt, r.PrivateData.Claims}
			if err := tbdex.VerifyDigest(r.Data.ClaimsHash, payload); err != nil {
				return fmt.Errorf("failed to verify claims: %w", err)
			}
		}
	}

	if r.Data.Payin.PaymentDetailsHash != "" {
		if strict && (r.PrivateData.Payin.PaymentDetails == nil) {
			return errors.New("strict verification: payin details hash is set but payin details are missing")

		}
		if r.PrivateData.Payin.PaymentDetails != nil {
			payload := []any{r.PrivateData.Salt, r.PrivateData.Payin.PaymentDetails}
			if err := tbdex.VerifyDigest(r.Data.Payin.PaymentDetailsHash, payload); err != nil {
				return fmt.Errorf("failed to verify payin: %w", err)
			}
		}
	}

	if r.Data.Payout.PaymentDetailsHash != "" {
		if strict && (r.PrivateData.Payout.PaymentDetails == nil) {
			return errors.New("strict verification: payout details hash is set but payout details are missing")

		}
		if r.PrivateData.Payout.PaymentDetails != nil {
			payload := []any{r.PrivateData.Salt, r.PrivateData.Payout.PaymentDetails}
			if err := tbdex.VerifyDigest(r.Data.Payout.PaymentDetailsHash, payload); err != nil {
				return fmt.Errorf("failed to verify payout: %w", err)
			}
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
