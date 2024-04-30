package order

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/TBD54566975/tbdex-go/tbdex"
	"github.com/tbd54566975/web5-go/dids/did"
	"go.jetpack.io/typeid"
)

const Kind = "order"

type Order struct {
	Metadata  tbdex.MessageMetadata `json:"metadata"`
	Data      Data                  `json:"data"`
	Signature string                `json:"signature"`
}

type Data struct{}

// Digest computes a hash of the order
func (o Order) Digest() ([]byte, error) {
	payload := map[string]any{"metadata": o.Metadata, "data": o.Data}

	hashed, err := tbdex.DigestJSON(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to compute order digest: %w", err)
	}

	return hashed, nil
}

func (o *Order) Sign(bearerDID did.BearerDID) error {
	o.Metadata.From = bearerDID.URI

	signature, err := tbdex.Sign(o, bearerDID)
	if err != nil {
		return fmt.Errorf("failed to sign order: %w", err)
	}

	o.Signature = signature

	return nil
}

// Verify verifies the order's signature.
func (o *Order) Verify() error {
	decoded, err := tbdex.VerifySignature(o, o.Signature)
	if err != nil {
		return fmt.Errorf("failed to verify order signature: %w", err)
	}

	if decoded.SignerDID.URI != o.Metadata.From {
		return fmt.Errorf("signer: %s does not match message metadata from: %s", decoded.SignerDID.URI, o.Metadata.From)
	}

	return nil
}

// UnmarshalJSON validates and unmarshals the input data into an Order.
func (o *Order) UnmarshalJSON(data []byte) error {
	err := tbdex.Validate(tbdex.TypeMessage, data, tbdex.WithKind(Kind))
	if err != nil {
		return fmt.Errorf("invalid order: %w", err)
	}

	type shadow Order
	ret := shadow{}

	err = json.Unmarshal(data, &ret)
	if err != nil {
		return fmt.Errorf("failed to unmarshal order: %w", err)
	}

	*o = Order(ret)

	return nil
}

// Parse unmarshals the provided input into an Order and then verifies the signature.
func Parse(data []byte) (Order, error) {
	var o Order
	err := json.Unmarshal(data, &o)
	if err != nil {
		return Order{}, fmt.Errorf("failed to unmarshal order: %w", err)
	}

	err = o.Verify()
	if err != nil {
		return o, fmt.Errorf("integrity mismatch: %w", err)
	}

	return o, nil
}

// CreateOption defines a type for functions that can modify the createOptions struct.
type CreateOption func(*createOptions)

type createOptions struct {
	createdAt  time.Time
	id         string
	externalID string
	protocol   string
}

// ID can be passed to [Create] to provide a custom id.
func ID(id string) CreateOption {
	return func(r *createOptions) {
		r.id = id
	}
}

// CreatedAt can be passed to [Create] to provide a custom created at time.
func CreatedAt(t time.Time) CreateOption {
	return func(q *createOptions) {
		q.createdAt = t
	}
}

// ExternalID can be passed to [Create] to provide a external id.
func ExternalID(externalID string) CreateOption {
	return func(q *createOptions) {
		q.externalID = externalID
	}
}

func Create(from, to, exchangeID string, opts ...CreateOption) Order {
	options := createOptions{
		id:        typeid.Must(typeid.WithPrefix(Kind)).String(),
		createdAt: time.Now(),
		protocol:  "1.0",
	}

	for _, o := range opts {
		o(&options)
	}

	return Order{
		Metadata: tbdex.MessageMetadata{
			From:       from,
			To:         to,
			Kind:       Kind,
			ID:         options.id,
			ExchangeID: exchangeID,
			CreatedAt:  options.createdAt.UTC().Format(time.RFC3339),
			ExternalID: options.externalID,
			Protocol:   options.protocol,
		},
		Data: Data{},
	}
}