package orderinstructions

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/TBD54566975/tbdex-go/tbdex/cancel"
	"github.com/TBD54566975/tbdex-go/tbdex/closemsg"
	"github.com/TBD54566975/tbdex-go/tbdex/crypto"
	"github.com/TBD54566975/tbdex-go/tbdex/message"
	"github.com/TBD54566975/tbdex-go/tbdex/validator"
	"github.com/tbd54566975/web5-go/dids/did"
	"go.jetpack.io/typeid"
)

// Kind is the value used within a message's metadata.kind
const Kind = "orderinstructions"

// ValidNext returns the valid message kinds that can follow an orderinstructions.
func ValidNext() []string {
	return []string{closemsg.Kind, cancel.Kind}
}

// OrderInstructions represents a tbdex [orderinstructions] message.
//
// [orderinstructions]: https://github.com/TBD54566975/tbdex/tree/main/specs/protocol#orderinstructions
type OrderInstructions struct {
	Metadata  message.Metadata `json:"metadata,omitempty"`
	Data      Data             `json:"data,omitempty"`
	Signature string           `json:"signature,omitempty"`
}

// GetMetadata returns the metadata of the message
func (o OrderInstructions) GetMetadata() message.Metadata {
	return o.Metadata
}

// GetKind returns the kind of message
func (o OrderInstructions) GetKind() string {
	return o.Metadata.Kind
}

// GetValidNext returns the valid message kinds that can follow an orderinstructions.
func (o OrderInstructions) GetValidNext() []string {
	return ValidNext()
}

// IsValidNext checks if the kind is a valid next message kind for an orderinstructions.
func (o OrderInstructions) IsValidNext(kind string) bool {
	for _, k := range ValidNext() {
		if k == kind {
			return true
		}
	}
	return false
}

// Data represents the data field of an orderinstructions message.
type Data struct {
	Payin  *PaymentInstruction `json:"payin,omitempty"`
	Payout *PaymentInstruction `json:"payout,omitempty"`
}

// PaymentInstruction contains instructions with plain text and/or a link
type PaymentInstruction struct {
	Link        string `json:"link,omitempty"`
	Instruction string `json:"instruction,omitempty"`
}

// Digest computes a hash of the orderinstructions
func (o OrderInstructions) Digest() ([]byte, error) {
	payload := map[string]any{"metadata": o.Metadata, "data": o.Data}

	hashed, err := crypto.DigestJSON(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to compute orderinstructions digest: %w", err)
	}

	return hashed, nil
}

// Verify verifies the orderinstructions's signature.
func (o *OrderInstructions) Verify() error {
	decoded, err := crypto.VerifySignature(o, o.Signature)
	if err != nil {
		return fmt.Errorf("failed to verify orderinstructions signature: %w", err)
	}

	if decoded.SignerDID.URI != o.Metadata.From {
		return fmt.Errorf("signer: %s does not match message metadata from: %s", decoded.SignerDID.URI, o.Metadata.From)
	}

	return nil
}

// UnmarshalJSON validates and unmarshals the input data into an OrderInstructions.
func (o *OrderInstructions) UnmarshalJSON(data []byte) error {
	err := validator.Validate(validator.TypeMessage, data, validator.WithKind(Kind))
	if err != nil {
		return fmt.Errorf("invalid orderinstructions: %w", err)
	}

	ret := orderinstructions{}
	err = json.Unmarshal(data, &ret)
	if err != nil {
		return fmt.Errorf("failed to unmarshal orderinstructions: %w", err)
	}

	*o = OrderInstructions(ret)

	return nil
}

// Parse unmarshals the provided input into an OrderInstructions and then verifies the signature.
func Parse(data []byte) (OrderInstructions, error) {
	var o OrderInstructions
	err := json.Unmarshal(data, &o)
	if err != nil {
		return OrderInstructions{}, fmt.Errorf("failed to unmarshal orderinstructions: %w", err)
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
	createdAt         time.Time
	id                string
	externalID        string
	protocol          string
	payinInstruction  *PaymentInstruction
	payoutInstruction *PaymentInstruction
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

// Create creates a new orderinstructions message. The following are generated by default unless custom values are provided:
//   - created at time is set to the current time
//   - protocol is set to "1.0"
//   - id is autogenerated
func Create(fromDID did.BearerDID, to, exchangeID string, opts ...CreateOption) (OrderInstructions, error) {
	options := createOptions{
		id:        typeid.Must(typeid.WithPrefix(Kind)).String(),
		createdAt: time.Now(),
		protocol:  "1.0",
	}

	for _, o := range opts {
		o(&options)
	}

	o := OrderInstructions{
		Metadata: message.Metadata{
			From:       fromDID.URI,
			To:         to,
			Kind:       Kind,
			ID:         options.id,
			ExchangeID: exchangeID,
			CreatedAt:  options.createdAt.UTC().Format(time.RFC3339),
			ExternalID: options.externalID,
			Protocol:   options.protocol,
		},
		Data: Data{
			Payin:  options.payinInstruction,
			Payout: options.payoutInstruction,
		},
	}

	signature, err := crypto.Sign(o, fromDID)
	if err != nil {
		return OrderInstructions{}, fmt.Errorf("failed to sign orderinstructions: %w", err)
	}

	o.Signature = signature

	return o, nil
}

// PayinInstruction is an option that allows setting a custom payin [PaymentInstruction]
func PayinInstruction(opts ...PaymentInstructionOptions) CreateOption {
	p := paymentInstructionOptions{}
	for _, opt := range opts {
		opt(&p)
	}

	return func(q *createOptions) {
		q.payinInstruction = &PaymentInstruction{
			Instruction: p.Instruction,
			Link:        p.Link,
		}
	}
}

// PayoutInstruction is an option that allows setting a custom payout [PaymentInstruction]
func PayoutInstruction(opts ...PaymentInstructionOptions) CreateOption {
	p := paymentInstructionOptions{}
	for _, opt := range opts {
		opt(&p)
	}

	return func(q *createOptions) {
		q.payoutInstruction = &PaymentInstruction{
			Instruction: p.Instruction,
			Link:        p.Link,
		}
	}
}

type paymentInstructionOptions struct {
	Link        string
	Instruction string
}

// PaymentInstructionOptions defines a type for functions that can modify the paymentInstructionOptions struct.
type PaymentInstructionOptions func(*paymentInstructionOptions)

// Link is an option for [NewInstruction] that allows setting a custom link.
func Link(link string) PaymentInstructionOptions {
	return func(p *paymentInstructionOptions) {
		p.Link = link
	}
}

// Instruction is an option for [NewInstruction] that allows setting custom text.
func Instruction(instruction string) PaymentInstructionOptions {
	return func(p *paymentInstructionOptions) {
		p.Instruction = instruction
	}
}

type orderinstructions OrderInstructions
