package offering

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"github.com/TBD54566975/tbdex-go/tbdex"
	"github.com/gowebpki/jcs"
	"github.com/tbd54566975/web5-go/dids/did"
	"github.com/tbd54566975/web5-go/pexv2"
	"go.jetpack.io/typeid"
)

// Kind distinguishes between different resource kinds
const Kind = "offering"

// Offering is a resource created by a PFI to define requirements for a given currency pair offered for exchange.
type Offering struct {
	tbdex.ResourceMetadata `json:"metadata"`
	Data                   `json:"data"`
	Signature              string `json:"signature"`
}

// Data represents the data of an Offering.
type Data struct {
	Description    string                        `json:"description"`
	Rate           string                        `json:"payoutUnitsPerPayinUnit"`
	Payin          PayinDetails                  `json:"payin,omitempty"`
	Payout         PayoutDetails                 `json:"payout,omitempty"`
	RequiredClaims *pexv2.PresentationDefinition `json:"requiredClaims,omitempty"`
}

// PayinDetails represents the details of the payin part of an Offering.
type PayinDetails struct {
	CurrencyCode string        `json:"currencyCode"`
	Min          string        `json:"min,omitempty"`
	Max          string        `json:"max,omitempty"`
	Methods      []PayinMethod `json:"methods,omitempty"`
}

// PayoutDetails represents the details of the payout part of an Offering.
type PayoutDetails struct {
	CurrencyCode string         `json:"currencyCode"`
	Min          string         `json:"min,omitempty"`
	Max          string         `json:"max,omitempty"`
	Methods      []PayoutMethod `json:"methods,omitempty"`
}

// PayinMethod represents a single payment option on an Offering.
type PayinMethod struct {
	Kind                   string          `json:"kind"`
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
	Kind                    string          `json:"kind"`
	Name                    string          `json:"name,omitempty"`
	Description             string          `json:"description,omitempty"`
	Group                   string          `json:"group,omitempty"`
	RequiredPaymentDetails  json.RawMessage `json:"requiredPaymentDetails,omitempty"` // TODO: change to JSON Schema type
	Fee                     string          `json:"fee,omitempty"`
	Min                     string          `json:"min,omitempty"`
	Max                     string          `json:"max,omitempty"`
	EstimatedSettlementTime uint64          `json:"estimatedSettlementTime"`
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
	payload := map[string]any{"metadata": o.ResourceMetadata, "data": o.Data}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to JSON marshal offering: %w", err)
	}

	canonicalized, err := jcs.Transform(payloadBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to canonicalize offering: %w", err)
	}

	hasher := sha256.New()
	_, err = hasher.Write(canonicalized)
	if err != nil {
		return nil, fmt.Errorf("failed to compute digest: %w", err)

	}

	return hasher.Sum(nil), nil
}

// Sign cryptographically signs the Resource using DID's private key
func (o *Offering) Sign(bearerDID did.BearerDID) error {
	o.From = bearerDID.URI

	signature, err := tbdex.Sign(o, bearerDID)
	if err != nil {
		return fmt.Errorf("failed to sign offering: %w", err)
	}

	o.Signature = signature

	return nil
}

// UnmarshalJSON validates and unmarshals the input data into an Offering.
func (o *Offering) UnmarshalJSON(data []byte) error {
	err := tbdex.Validate(tbdex.TypeResource, data, tbdex.WithKind(Kind))
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
	decoded, err := tbdex.VerifySignature(o, o.Signature)
	if err != nil {
		return fmt.Errorf("failed to verify Offering signature: %w", err)
	}

	if decoded.SignerDID.URI != o.ResourceMetadata.From {
		return fmt.Errorf("signer: %s does not match message metadata from: %s", decoded.SignerDID.URI, o.ResourceMetadata.From)
	}

	return nil
}

// Parse validates, parses input data into an Offering, and verifies the signature.
func (o *Offering) Parse(data []byte) error {
	if err := o.UnmarshalJSON(data); err != nil {
		return fmt.Errorf("failed to unmarshal Offering: %w", err)
	}

	if err := o.Verify(); err != nil {
		return fmt.Errorf("failed to verify Offering: %w", err)
	}
	return nil
}

type offering Offering
