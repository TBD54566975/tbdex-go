package cancel

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
const Kind = "cancel"

// ValidNext returns the valid message kinds that can follow a cancel.
func ValidNext() []string {
	// todo hardcoded orderstatus kind here because otherwise i am introducing a circular dependency :(
	return []string{"orderstatus", closemsg.Kind}
}

// Cancel represents a cancel message within the exchange.
type Cancel struct {
	Metadata  message.Metadata `json:"metadata,omitempty"`
	Data      Data             `json:"data,omitempty"`
	Signature string           `json:"signature,omitempty"`
}

// GetMetadata returns the metadata of the message.
func (c Cancel) GetMetadata() message.Metadata {
	return c.Metadata
}

// GetKind returns the kind of message.
func (c Cancel) GetKind() string {
	return c.Metadata.Kind
}

// GetValidNext returns the kinds of messages that can follow a cancel.
func (c Cancel) GetValidNext() []string {
	return ValidNext()
}

// IsValidNext checks if the kind is a valid next message kind for a cancel.
func (c Cancel) IsValidNext(kind string) bool {
	for _, k := range ValidNext() {
		if k == kind {
			return true
		}
	}
	return false
}

// Data encapsulates the data content of a cancel.
type Data struct {
	Reason string `json:"reason,omitempty"`
}

// Digest computes a hash of the message
func (c Cancel) Digest() ([]byte, error) {
	payload := map[string]any{"metadata": c.Metadata, "data": c.Data}

	hashed, err := crypto.DigestJSON(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to digest cancel: %w", err)
	}

	return hashed, nil
}

// Verify verifies the signature of the Cancel.
func (c *Cancel) Verify() error {
	decoded, err := crypto.VerifySignature(c, c.Signature)
	if err != nil {
		return fmt.Errorf("failed to verify cancel signature: %w", err)
	}

	if decoded.SignerDID.URI != c.Metadata.From {
		return fmt.Errorf("signer: %s does not match message metadata from: %s", decoded.SignerDID.URI, c.Metadata.From)
	}

	return nil
}

// UnmarshalJSON validates and unmarshals the input data into a Cancel.
func (c *Cancel) UnmarshalJSON(data []byte) error {
	err := validator.Validate(validator.TypeMessage, data, validator.WithKind(Kind))
	if err != nil {
		return fmt.Errorf("invalid cancel: %w", err)
	}

	ret := cancelType{}
	err = json.Unmarshal(data, &ret)
	if err != nil {
		return fmt.Errorf("failed to JSON unmarshal cancel: %w", err)
	}

	*c = Cancel(ret)

	return nil
}

// Parse validates and unmarshals the input data into a Cancel.
func Parse(data []byte) (Cancel, error) {
	c := Cancel{}
	if err := json.Unmarshal(data, &c); err != nil {
		return Cancel{}, fmt.Errorf("failed to unmarshal Cancel: %w", err)
	}

	if err := c.Verify(); err != nil {
		return Cancel{}, fmt.Errorf("failed to verify Cancel: %w", err)
	}

	return c, nil
}

// Create creates a new Cancel message.
func Create(fromDID did.BearerDID, to, exchangeID string, opts ...CreateOption) (Cancel, error) {
	o := createOptions{
		id:        typeid.Must(typeid.WithPrefix(Kind)).String(),
		createdAt: time.Now(),
		protocol:  "1.0",
	}

	for _, opt := range opts {
		opt(&o)
	}

	c := Cancel{
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
		Data: Data{Reason: o.reason},
	}

	signature, err := crypto.Sign(c, fromDID)
	if err != nil {
		return Cancel{}, fmt.Errorf("failed to sign cancel: %w", err)
	}

	c.Signature = signature

	return c, nil
}

type createOptions struct {
	id         string
	createdAt  time.Time
	protocol   string
	externalID string
	reason     string
}

// CreateOption defines a type for functions that can modify the createOptions struct.
type CreateOption func(*createOptions)

// ID can be passed to [Create] to provide a custom id.
func ID(id string) CreateOption {
	return func(c *createOptions) {
		c.id = id
	}
}

// CreatedAt can be passed to [Create] to provide a custom created at time.
func CreatedAt(t time.Time) CreateOption {
	return func(c *createOptions) {
		c.createdAt = t
	}
}

// ExternalID can be passed to [Create] to provide a custom external id.
func ExternalID(externalID string) CreateOption {
	return func(c *createOptions) {
		c.externalID = externalID
	}
}

// Reason can be passed to [Create] to provide a custom reason.
func Reason(reason string) CreateOption {
	return func(c *createOptions) {
		c.reason = reason
	}
}

type cancelType Cancel
