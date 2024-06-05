package offering

import (
	"encoding/json"
	"fmt"

	"github.com/TBD54566975/tbdex-go/tbdex/crypto"
	"github.com/TBD54566975/tbdex-go/tbdex/resource"
	"github.com/TBD54566975/tbdex-go/tbdex/validator"
	"github.com/tbd54566975/web5-go/dids/did"
	"github.com/tbd54566975/web5-go/pexv2"
	"go.jetpack.io/typeid"
)

// Kind distinguishes between different resource kinds
const Kind = "offering"

// Offering is a resource created by a PFI to define requirements for a given currency pair offered for exchange.
type Offering struct {
	Metadata  resource.Metadata `json:"metadata,omitempty"`
	Data      Data              `json:"data,omitempty"`
	Signature string            `json:"signature,omitempty"`
}

// Data represents the data of an Offering.
type Data struct {
	Description    string                        `json:"description,omitempty"`
	Rate           string                        `json:"payoutUnitsPerPayinUnit,omitempty"`
	Payin          *PayinDetails                 `json:"payin,omitempty"`
	Payout         *PayoutDetails                `json:"payout,omitempty"`
	RequiredClaims *pexv2.PresentationDefinition `json:"requiredClaims,omitempty"`
}

// PayinDetails represents the details of the payin part of an Offering.
type PayinDetails struct {
	CurrencyCode string        `json:"currencyCode,omitempty"`
	Min          string        `json:"min,omitempty"`
	Max          string        `json:"max,omitempty"`
	Methods      []PayinMethod `json:"methods,omitempty"`
}

// PayoutDetails represents the details of the payout part of an Offering.
type PayoutDetails struct {
	CurrencyCode string         `json:"currencyCode,omitempty"`
	Min          string         `json:"min,omitempty"`
	Max          string         `json:"max,omitempty"`
	Methods      []PayoutMethod `json:"methods,omitempty"`
}

// PaymentMethod is an interface for PayinMethod and PayoutMethod.
type PaymentMethod interface {
	GetKind() string
	GetRequiredPaymentDetails() json.RawMessage
}

// GetKind returns the kind of the PayinMethod.
func (pm PayinMethod) GetKind() string {
    return pm.Kind
}

// GetRequiredPaymentDetails returns the required payment details of the PayinMethod.
func (pm PayinMethod) GetRequiredPaymentDetails() json.RawMessage {
    return pm.RequiredPaymentDetails
}

// GetKind returns the kind of the PayoutMethod.
func (pm PayoutMethod) GetKind() string {
    return pm.Kind
}

// GetRequiredPaymentDetails returns the required payment details of the PayoutMethod.
func (pm PayoutMethod) GetRequiredPaymentDetails() json.RawMessage {
    return pm.RequiredPaymentDetails
}

// PayinMethod represents a single payment option on an Offering.
type PayinMethod struct {
	Kind                   string          `json:"kind,omitempty"`
	Name                   string          `json:"name,omitempty"`
	Description            string          `json:"description,omitempty"`
	Group                  string          `json:"group,omitempty"`
	RequiredPaymentDetails json.RawMessage `json:"requiredPaymentDetails,omitempty"` // TODO: change to JSON Schema type
	Fee                    string          `json:"fee,omitempty"`
	Min                    string          `json:"min,omitempty"`
	Max                    string          `json:"max,omitempty"`
}

// PayoutMethod contains all the fields from PaymentMethod, in addition to estimated settlement time.
type PayoutMethod struct {
	Kind                    string          `json:"kind,omitempty"`
	Name                    string          `json:"name,omitempty"`
	Description             string          `json:"description,omitempty"`
	Group                   string          `json:"group,omitempty"`
	RequiredPaymentDetails  json.RawMessage `json:"requiredPaymentDetails,omitempty"` // TODO: change to JSON Schema type
	Fee                     string          `json:"fee,omitempty"`
	Min                     string          `json:"min,omitempty"`
	Max                     string          `json:"max,omitempty"`
	EstimatedSettlementTime uint64          `json:"estimatedSettlementTime,omitempty"`
}

// ID is a unique identifier for an Offering.
type ID struct {
	typeid.TypeID[ID]
}

// Prefix returns the prefix for the Offering ID.
func (id ID) Prefix() string { return Kind }

// Digest computes a hash of the resource
// A digest is the output of the hash function. It's a fixed-size string of bytes
//   - that uniquely represents the data input into the hash function. The digest is often used for
//   - data integrity checks, as any alteration in the input data results in a significantly
//   - different digest.
//     *
//   - It takes the algorithm identifier of the hash function and data to digest as input and returns
//   - the digest of the data.
func (o Offering) Digest() ([]byte, error) {
	payload := map[string]any{"metadata": o.Metadata, "data": o.Data}

	hashed, err := crypto.DigestJSON(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to digest offering: %w", err)
	}

	return hashed, nil
}

// Sign cryptographically signs the Resource using DID's private key
func (o *Offering) Sign(bearerDID did.BearerDID) error {
	o.Metadata.From = bearerDID.URI

	signature, err := crypto.Sign(o, bearerDID)
	if err != nil {
		return fmt.Errorf("failed to sign offering: %w", err)
	}

	o.Signature = signature

	return nil
}

// UnmarshalJSON validates and unmarshals the input data into an Offering.
func (o *Offering) UnmarshalJSON(data []byte) error {
	err := validator.Validate(validator.TypeResource, data, validator.WithKind(Kind))
	if err != nil {
		return fmt.Errorf("invalid offering: %w", err)
	}

	off := offering{}
	err = json.Unmarshal(data, &off)
	if err != nil {
		return fmt.Errorf("failed to JSON unmarshal offering: %w", err)
	}

	*o = Offering(off)

	return nil
}

// Verify verifies the signature of the Offering.
func (o *Offering) Verify() error {
	decoded, err := crypto.VerifySignature(o, o.Signature)
	if err != nil {
		return fmt.Errorf("failed to verify Offering signature: %w", err)
	}

	if decoded.SignerDID.URI != o.Metadata.From {
		return fmt.Errorf("signer: %s does not match resource metadata from: %s", decoded.SignerDID.URI, o.Metadata.From)
	}

	return nil
}

// Parse validates, parses input data into an Offering, and verifies the signature.
func (o *Offering) Parse(data []byte) error {
	if err := json.Unmarshal(data, &o); err != nil {
		return fmt.Errorf("failed to unmarshal Offering: %w", err)
	}

	if err := o.Verify(); err != nil {
		return fmt.Errorf("failed to verify Offering: %w", err)
	}
	return nil
}

type offering Offering
