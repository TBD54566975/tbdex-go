package offering

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/TBD54566975/tbdex-go/protocol/resource"
	"github.com/gowebpki/jcs"
	"github.com/tbd54566975/web5-go/dids/did"
	"github.com/tbd54566975/web5-go/jws"
	"go.jetpack.io/typeid"
)

const Kind = "offering"

// Offering is a resource created by a PFI to define requirements for a given currency pair offered for exchange.
type Offering struct {
	resource.Metadata `json:"metadata"`
	Data              `json:"data"`
	Signature         string `json:"signature"`
}

type Data struct {
	Description string        `json:"description"`
	Rate        string        `json:"payoutUnitsPerPayinUnit"`
	Payin       PayinDetails  `json:"payin,omitempty"`
	Payout      PayoutDetails `json:"payout,omitempty"`
}

type PayinDetails struct {
	CurrencyCode string        `json:"currencyCode"`
	Min          string        `json:"min,omitempty"`
	Max          string        `json:"max,omitempty"`
	Methods      []PayinMethod `json:"methods,omitempty"`
}

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

type createOptions struct {
	id          string
	createdAt   time.Time
	updatedAt   time.Time
	description string
}

type CreateOption func(*createOptions)

func ID(id string) CreateOption {
	return func(o *createOptions) {
		o.id = id
	}
}

func CreatedAt(t time.Time) CreateOption {
	return func(o *createOptions) {
		o.createdAt = t
	}
}

func UpdatedAt(t time.Time) CreateOption {
	return func(o *createOptions) {
		o.updatedAt = t
	}
}

func Description(d string) CreateOption {
	return func(o *createOptions) {
		o.description = d
	}
}

func Create(payin PayinDetails, payout PayoutDetails, rate string, opts ...CreateOption) (Offering, error) {
	defaultID, err := typeid.WithPrefix(Kind)
	if err != nil {
		return Offering{}, fmt.Errorf("failed to generate default id: %w", err)
	}

	o := createOptions{
		id:          defaultID.String(),
		createdAt:   time.Now(),
		updatedAt:   time.Now(),
		description: fmt.Sprintf("%s for %s", payout.CurrencyCode, payin.CurrencyCode),
	}

	for _, opt := range opts {
		opt(&o)
	}

	if len(payin.Methods) == 0 {
		return Offering{}, errors.New("1 payin method is required.")
	}

	if len(payout.Methods) == 0 {
		return Offering{}, errors.New("1 payout method is required.")
	}

	return Offering{
		Metadata: resource.Metadata{
			Kind:      Kind,
			ID:        o.id,
			CreatedAt: o.createdAt.UTC().Format(time.RFC3339),
			UpdatedAt: o.updatedAt.UTC().Format(time.RFC3339),
		},
		Data: Data{
			Payin:       payin,
			Payout:      payout,
			Rate:        rate,
			Description: o.description,
		},
	}, nil
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
