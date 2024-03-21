package offering

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"github.com/TBD54566975/tbdex-go/protocol/resource"
	"github.com/gowebpki/jcs"
	"github.com/tbd54566975/web5-go/dids/did"
	"github.com/tbd54566975/web5-go/jws"
)

const Kind = "offering"

// Offering is a resource created by a PFI to define requirements for a given currency pair offered for exchange.
type Offering struct {
	resource.Metadata `json:"metadata"`
	Data              `json:"data"`
	Signature         string `json:"signature"`
}

// Data represents the data of an Offering.
type Data struct {
	Description string        `json:"description"`
	Rate        string        `json:"payoutUnitsPerPayinUnit"`
	Payin       PayinDetails  `json:"payin,omitempty"`
	Payout      PayoutDetails `json:"payout,omitempty"`
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

type PaymentMethod struct {
	Kind                   string `json:"kind"`
	Name                   string `json:"name,omitempty"`
	Description            string `json:"description,omitempty"`
	Group                  string `json:"group,omitempty"`
	RequiredPaymentDetails string `json:"requiredPaymentDetails,omitempty"` // TODO: change to JSON Schema type
	Fee                    string `json:"fee,omitempty"`
	Min                    string `json:"min,omitempty"`
	Max                    string `json:"max,omitempty"`
}

type PayinMethod = PaymentMethod

type PayoutMethod struct {
	PaymentMethod
	EstimatedSettlementTime uint64 `json:"estimatedSettlementTime"`
}

func (o Offering) Digest() ([]byte, error) {
	payload := map[string]any{"metadata": o.Metadata, "data": o.Data}
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

func (o Offering) Sign(bearerDID did.BearerDID) (Offering, error) {
	o.From = bearerDID.URI

	digest, err := o.Digest()
	if err != nil {
		return Offering{}, fmt.Errorf("failed to sign offering: %w", err)
	}

	signature, err := jws.Sign(digest, bearerDID, jws.DetachedPayload(true))
	if err != nil {
		return Offering{}, fmt.Errorf("failed to sign offering: %w", err)
	}

	o.Signature = signature

	return o, nil
}
