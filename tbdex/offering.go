package tbdex

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"github.com/gowebpki/jcs"
	"github.com/tbd54566975/web5-go/dids/did"
)

// OfferingKind distinguishes between different resource kinds
const OfferingKind = "offering"

// Offering is a resource created by a PFI to define requirements for a given currency pair offered for exchange.
type Offering struct {
	ResourceMetadata `json:"metadata"`
	OfferingData     `json:"data"`
	Signature        string `json:"signature"`
}

// OfferingData represents the data of an Offering.
type OfferingData struct {
	Description string                `json:"description"`
	Rate        string                `json:"payoutUnitsPerPayinUnit"`
	Payin       OfferingPayinDetails  `json:"payin,omitempty"`
	Payout      OfferingPayoutDetails `json:"payout,omitempty"`
}

// OfferingPayinDetails represents the details of the payin part of an Offering.
type OfferingPayinDetails struct {
	CurrencyCode string                `json:"currencyCode"`
	Min          string                `json:"min,omitempty"`
	Max          string                `json:"max,omitempty"`
	Methods      []OfferingPayinMethod `json:"methods,omitempty"`
}

// OfferingPayoutDetails represents the details of the payout part of an Offering.
type OfferingPayoutDetails struct {
	CurrencyCode string                 `json:"currencyCode"`
	Min          string                 `json:"min,omitempty"`
	Max          string                 `json:"max,omitempty"`
	Methods      []OfferingPayoutMethod `json:"methods,omitempty"`
}

// OfferingPaymentMethod represents a single payment option on an Offering.
type OfferingPaymentMethod struct {
	Kind                   string `json:"kind"`
	Name                   string `json:"name,omitempty"`
	Description            string `json:"description,omitempty"`
	Group                  string `json:"group,omitempty"`
	RequiredPaymentDetails string `json:"requiredPaymentDetails,omitempty"` // TODO: change to JSON Schema type
	Fee                    string `json:"fee,omitempty"`
	Min                    string `json:"min,omitempty"`
	Max                    string `json:"max,omitempty"`
}

// OfferingPayinMethod is an alias for PaymentMethod.
type OfferingPayinMethod = OfferingPaymentMethod

// OfferingPayoutMethod contains all the fields from PaymentMethod, in addition to estimated settlement time.
type OfferingPayoutMethod struct {
	OfferingPaymentMethod
	EstimatedSettlementTime uint64 `json:"estimatedSettlementTime"`
}

// Digest computes a hash of the resource
// A digest is the output of the hash function. It's a fixed-size string of bytes
//   - that uniquely represents the data input into the hash function. The digest is often used for
//   - data integrity checks, as any alteration in the input data results in a significantly
//   - different digest.
//     *
//   - It takes the algorithm identifier of the hash function and data to digest as input and returns
//   - the digest of the data.
func (o Offering) Digest() ([]byte, error) {
	payload := map[string]any{"metadata": o.ResourceMetadata, "data": o.OfferingData}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal offering: %w", err)
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
func (o Offering) Sign(bearerDID did.BearerDID) error {
	o.From = bearerDID.URI

	signature, err := Sign(o, bearerDID)
	if err != nil {
		return fmt.Errorf("failed to sign offering: %w", err)
	}

	o.Signature = signature

	return nil
}

// UnmarshalJSON validates and unmarshals the input data into an Offering.
func (o *Offering) UnmarshalJSON(data []byte) error {
	err := Validate(TypeResource, data, WithKind(OfferingKind))
	if err != nil {
		return fmt.Errorf("invalid offering: %w", err)
	}

	off := offering{}
	err = json.Unmarshal(data, &off)
	if err != nil {
		return fmt.Errorf("failed to unmarshal offering: %w", err)
	}

	*o = Offering(off)

	return nil
}

type offering Offering
