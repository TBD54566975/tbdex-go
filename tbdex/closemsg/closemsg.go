package closemsg

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/TBD54566975/tbdex-go/tbdex"
	"github.com/tbd54566975/web5-go/dids/did"
	"go.jetpack.io/typeid"
)

// Kind identifies this message kind
const Kind = "close"

// Close represents a close message within the exchange.
type Close struct {
	Metadata  tbdex.MessageMetadata `json:"metadata,omitempty"`
	Data      Data                  `json:"data,omitempty"`
	Signature string                `json:"signature,omitempty"`
}

// Data encapsulates the data content of a close.
type Data struct {
	Reason  string `json:"reason,omitempty"`
	Success bool   `json:"success,omitempty"`
}

// Digest computes a hash of the message
func (c Close) Digest() ([]byte, error) {
	payload := map[string]any{"metadata": c.Metadata, "data": c.Data}

	hashed, err := tbdex.DigestJSON(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to digest close: %w", err)
	}

	return hashed, nil
}

// Verify verifies the signature of the Close.
func (c *Close) Verify() error {
	decoded, err := tbdex.VerifySignature(c, c.Signature)
	if err != nil {
		return fmt.Errorf("failed to verify close signature: %w", err)
	}

	if decoded.SignerDID.URI != c.Metadata.From {
		return fmt.Errorf("signer: %s does not match message metadata from: %s", decoded.SignerDID.URI, c.Metadata.From)
	}

	return nil
}

// UnmarshalJSON validates and unmarshals the input data into a Close.
func (c *Close) UnmarshalJSON(data []byte) error {
	err := tbdex.Validate(tbdex.TypeMessage, data, tbdex.WithKind(Kind))
	if err != nil {
		return fmt.Errorf("invalid close: %w", err)
	}

	ret := closeType{}
	err = json.Unmarshal(data, &ret)
	if err != nil {
		return fmt.Errorf("failed to JSON unmarshal close: %w", err)
	}

	*c = Close(ret)

	return nil
}

// Parse validates and unmarshals the input data into a Close.
func Parse(data []byte) (Close, error) {
	os := Close{}
	err := os.UnmarshalJSON(data)
	if err != nil {
		return Close{}, err

	}

	err = os.Verify()
	if err != nil {
		return Close{}, fmt.Errorf("failed to verify close: %w", err)
	}

	return os, nil
}

// Create creates a new Close message.
func Create(fromDID did.BearerDID, to, exchangeID string, opts ...CreateOption) (Close, error) {
	o := createOptions{
		id:        typeid.Must(typeid.WithPrefix(Kind)).String(),
		createdAt: time.Now(),
		protocol:  "1.0",
	}

	for _, opt := range opts {
		opt(&o)
	}

	c := Close{
		Metadata: tbdex.MessageMetadata{
			From:       fromDID.URI,
			To:         to,
			Kind:       Kind,
			ID:         o.id,
			ExchangeID: exchangeID,
			CreatedAt:  o.createdAt.UTC().Format(time.RFC3339),
			ExternalID: o.externalID,
			Protocol:   o.protocol,
		},
		Data: Data{Reason: o.reason, Success: o.success},
	}

	signature, err := tbdex.Sign(c, fromDID)
	if err != nil {
		return Close{}, fmt.Errorf("failed to sign close: %w", err)
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
	success    bool
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

// Success can be passed to [Create] to provide a custom success.
func Success(success bool) CreateOption {
	return func(c *createOptions) {
		c.success = success
	}
}

type closeType Close
