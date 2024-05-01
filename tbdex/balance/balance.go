package balance

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/TBD54566975/tbdex-go/tbdex"
	"github.com/tbd54566975/web5-go/dids/did"
	"go.jetpack.io/typeid"
)

// Kind distinguishes between different resource kinds
const Kind = "balance"

// Balance is a resource to communicate the amounts of each currency held by the PFI on behalf of its customer.
type Balance struct {
	*tbdex.ResourceMetadata `json:"metadata,omitempty"`
	*Data                   `json:"data,omitempty"`
	Signature               string `json:"signature,omitempty"`
}

// Data represents the data of a Balance.
type Data struct {
	CurrencyCode string `json:"currencyCode,omitempty"`
	Available    string `json:"available,omitempty"`
}

// ID is a unique identifier for a Balance.
type ID struct {
	typeid.TypeID[ID]
}

// Prefix returns the prefix for the Balance ID.
func (id ID) Prefix() string { return Kind }

// Digest computes a hash of the resource
func (b Balance) Digest() ([]byte, error) {
	payload := map[string]any{"metadata": b.ResourceMetadata, "data": b.Data}

	hashed, err := tbdex.DigestJSON(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to digest balance: %w", err)
	}

	return hashed, nil
}

// Create a Balance object
func Create(fromDID did.BearerDID, currencyCode, availableAmount string, opts ...CreateOption) (Balance, error) {
	o := createOptions{
		id:        typeid.Must(typeid.New[ID]()),
		createdAt: time.Now(),
		updatedAt: time.Now(),
		protocol:  "1.0",
	}

	for _, opt := range opts {
		opt(&o)
	}

	b := Balance{
		ResourceMetadata: &tbdex.ResourceMetadata{
			From:      fromDID.URI,
			Kind:      Kind,
			ID:        o.id.String(),
			CreatedAt: o.createdAt.UTC().Format(time.RFC3339),
			UpdatedAt: o.updatedAt.UTC().Format(time.RFC3339),
			Protocol:  o.protocol,
		},
		Data: &Data{
			CurrencyCode: currencyCode,
			Available:    availableAmount,
		},
	}

	signature, err := tbdex.Sign(b, fromDID)
	if err != nil {
		return Balance{}, fmt.Errorf("failed to sign balance: %w", err)
	}

	b.Signature = signature

	return b, nil
}

type createOptions struct {
	id        ID
	createdAt time.Time
	updatedAt time.Time
	protocol  string
}

// CreateOption implements functional options pattern for [Create].
type CreateOption func(*createOptions)

// CreatedAt can be passed to [Create] to provide a custom created at time.
func CreatedAt(t time.Time) CreateOption {
	return func(o *createOptions) {
		o.createdAt = t
	}
}

// UpdatedAt can be passed to [Create] to provide a custom updated at time.
func UpdatedAt(t time.Time) CreateOption {
	return func(o *createOptions) {
		o.updatedAt = t
	}
}

// UnmarshalJSON validates and unmarshals the input data into a Balance.
func (b *Balance) UnmarshalJSON(data []byte) error {
	err := tbdex.Validate(tbdex.TypeResource, data, tbdex.WithKind(Kind))
	if err != nil {
		return fmt.Errorf("invalid balance: %w", err)
	}

	balance := balance{}
	err = json.Unmarshal(data, &balance)
	if err != nil {
		return fmt.Errorf("failed to JSON unmarshal balance: %w", err)
	}

	*b = Balance(balance)

	return nil
}

// Verify verifies the signature of the Balance.
func (b *Balance) Verify() error {
	decoded, err := tbdex.VerifySignature(b, b.Signature)
	if err != nil {
		return fmt.Errorf("failed to verify Balance signature: %w", err)
	}

	if decoded.SignerDID.URI != b.ResourceMetadata.From {
		return fmt.Errorf("signer: %s does not match resource metadata from: %s", decoded.SignerDID.URI, b.ResourceMetadata.From)
	}

	return nil
}

// Parse validates, parses input data into a Balance, and verifies the signature.
func Parse(data []byte) (Balance, error) {
	balance := Balance{}
	if err := balance.UnmarshalJSON(data); err != nil {
		return Balance{}, fmt.Errorf("failed to unmarshal Balance: %w", err)
	}

	if err := balance.Verify(); err != nil {
		return Balance{}, fmt.Errorf("failed to verify Balance: %w", err)
	}
	return balance, nil
}

type balance Balance
