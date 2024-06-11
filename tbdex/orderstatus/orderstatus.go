package orderstatus

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/TBD54566975/tbdex-go/tbdex/closemsg"
	"github.com/TBD54566975/tbdex-go/tbdex/crypto"
	"github.com/TBD54566975/tbdex-go/tbdex/message"
	"github.com/TBD54566975/tbdex-go/tbdex/validator"
	"github.com/tbd54566975/web5-go/dids/did"
	"go.jetpack.io/typeid"
)

// Kind identifies this message kind
const Kind = "orderstatus"

// ValidNext returns the valid next message kinds that can follow an orderstatus
func ValidNext() []string {
	return []string{Kind, closemsg.Kind}
}

// OrderStatus represents an order status message within the exchange.
type OrderStatus struct {
	Metadata  message.Metadata `json:"metadata,omitempty"`
	Data      Data             `json:"data,omitempty"`
	Signature string           `json:"signature,omitempty"`
}

// GetMetadata returns the metadata of the message
func (os OrderStatus) GetMetadata() message.Metadata {
	return os.Metadata
}

// GetKind returns the kind of message
func (os OrderStatus) GetKind() string {
	return os.Metadata.Kind
}

// GetValidNext returns the valid next message kinds that can follow an orderstatus
func (os OrderStatus) GetValidNext() []string {
	return []string{Kind, closemsg.Kind}
}

// IsValidNext checks if the kind is a valid next message kind for an orderstatus
func (os OrderStatus) IsValidNext(kind string) bool {
	for _, k := range ValidNext() {
		if k == kind {
			return true
		}
	}
	return false
}

// Data encapsulates the data content of an order status.
type Data struct {
	OrderStatus string `json:"orderStatus,omitempty"`
}

// Digest computes a hash of the message
func (os OrderStatus) Digest() ([]byte, error) {
	payload := map[string]any{"metadata": os.Metadata, "data": os.Data}

	hashed, err := crypto.DigestJSON(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to digest order status: %w", err)
	}

	return hashed, nil
}

// Verify verifies the signature of the OrderStatus.
func (os *OrderStatus) Verify() error {
	decoded, err := crypto.VerifySignature(os, os.Signature)
	if err != nil {
		return fmt.Errorf("failed to verify OrderStatus signature: %w", err)
	}

	if decoded.SignerDID.URI != os.Metadata.From {
		return fmt.Errorf("signer: %s does not match message metadata from: %s", decoded.SignerDID.URI, os.Metadata.From)
	}

	return nil
}

// UnmarshalJSON validates and unmarshals the input data into an OrderStatus.
func (os *OrderStatus) UnmarshalJSON(data []byte) error {
	err := validator.Validate(validator.TypeMessage, data, validator.WithKind(Kind))
	if err != nil {
		return fmt.Errorf("invalid order status: %w", err)
	}

	ret := orderStatus{}
	err = json.Unmarshal(data, &ret)
	if err != nil {
		return fmt.Errorf("failed to JSON unmarshal order status: %w", err)
	}

	*os = OrderStatus(ret)

	return nil
}

// Parse validates and unmarshals the input data into an OrderStatus.
func Parse(data []byte) (OrderStatus, error) {
	os := OrderStatus{}
	if err := json.Unmarshal(data, &os); err != nil {
		return OrderStatus{}, fmt.Errorf("failed to unmarshal order status: %w", err)

	}

	if err := os.Verify(); err != nil {
		return OrderStatus{}, fmt.Errorf("failed to verify order status: %w", err)
	}

	return os, nil
}

// Create creates a new OrderStatus message.
func Create(fromDID did.BearerDID, to, exchangeID, orderStatus string, opts ...CreateOption) (OrderStatus, error) {
	o := createOptions{
		id:        typeid.Must(typeid.WithPrefix(Kind)).String(),
		createdAt: time.Now(),
		protocol:  "1.0",
	}

	for _, opt := range opts {
		opt(&o)
	}

	os := OrderStatus{
		Metadata: message.Metadata{
			From:       fromDID.URI,
			To:         to,
			Kind:       Kind,
			ID:         o.id,
			ExchangeID: exchangeID,
			CreatedAt:  o.createdAt.UTC().Format(time.RFC3339),
			ExternalID: o.externalID,
			Protocol:   o.protocol,
		},
		Data: Data{OrderStatus: orderStatus},
	}

	signature, err := crypto.Sign(os, fromDID)
	if err != nil {
		return OrderStatus{}, fmt.Errorf("failed to sign order status: %w", err)
	}

	os.Signature = signature

	return os, nil
}

type createOptions struct {
	id         string
	createdAt  time.Time
	protocol   string
	externalID string
}

// CreateOption defines a type for functions that can modify the createOptions struct.
type CreateOption func(*createOptions)

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

// ExternalID can be passed to [Create] to provide a custom external id.
func ExternalID(externalID string) CreateOption {
	return func(q *createOptions) {
		q.externalID = externalID
	}
}

type orderStatus OrderStatus
