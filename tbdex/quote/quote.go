package quote

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/TBD54566975/tbdex-go/tbdex/cancel"
	"github.com/TBD54566975/tbdex-go/tbdex/closemsg"
	"github.com/TBD54566975/tbdex-go/tbdex/crypto"
	"github.com/TBD54566975/tbdex-go/tbdex/message"
	"github.com/TBD54566975/tbdex-go/tbdex/order"
	"github.com/TBD54566975/tbdex-go/tbdex/validator"
	"github.com/shopspring/decimal"
	"github.com/tbd54566975/web5-go/dids/did"
	"go.jetpack.io/typeid"
)

// Kind identifies this message kind
const Kind = "quote"

// ValidNext returns the valid message kinds that can follow a Quote.
func ValidNext() []string {
	return []string{order.Kind, closemsg.Kind, cancel.Kind}
}

// Quote represents a quote message within the exchange.
type Quote struct {
	Metadata  message.Metadata `json:"metadata,omitempty"`
	Data      Data             `json:"data,omitempty"`
	Signature string           `json:"signature,omitempty"`
}

// GetMetadata returns the metadata of the message
func (q Quote) GetMetadata() message.Metadata {
	return q.Metadata
}

// GetKind returns the kind of message
func (q Quote) GetKind() string {
	return Kind
}

// GetValidNext returns the valid message kinds that can follow a Quote.
func (q Quote) GetValidNext() []string {
	return ValidNext()
}

// IsValidNext checks if the kind is a valid next message kind for a Quote.
func (q Quote) IsValidNext(kind string) bool {
	for _, k := range ValidNext() {
		if k == kind {
			return true
		}
	}
	return false
}

// Data encapsulates the data content of a  quote.
type Data struct {
	ExpiresAt string       `json:"expiresAt,omitempty"`
	Rate      string       `json:"payoutUnitsPerPayinUnit,omitempty"`
	Payin     QuoteDetails `json:"payin,omitempty"`
	Payout    QuoteDetails `json:"payout,omitempty"`
}

// QuoteDetails describes the relevant information of a currency that is being sent or received
type QuoteDetails struct {
	CurrencyCode       string              `json:"currencyCode,omitempty"`
	Subtotal           string              `json:"subtotal,omitempty"`
	Fee                string              `json:"fee,omitempty"`
	Total              string              `json:"total,omitempty"`
}


// Digest computes a hash of the quote
func (q Quote) Digest() ([]byte, error) {
	payload := map[string]any{"metadata": q.Metadata, "data": q.Data}

	hashed, err := crypto.DigestJSON(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to digest quote: %w", err)
	}

	return hashed, nil
}

// Verify verifies the signature of the quote.
func (q *Quote) Verify() error {
	decoded, err := crypto.VerifySignature(q, q.Signature)
	if err != nil {
		return fmt.Errorf("failed to verify quote signature: %w", err)
	}

	if decoded.SignerDID.URI != q.Metadata.From {
		return fmt.Errorf("signer: %s does not match message metadata from: %s", decoded.SignerDID.URI, q.Metadata.From)
	}

	return nil
}

// UnmarshalJSON validates and unmarshals the input data into an Quote.
func (q *Quote) UnmarshalJSON(data []byte) error {
	err := validator.Validate(validator.TypeMessage, data, validator.WithKind(Kind))
	if err != nil {
		return fmt.Errorf("invalid quote: %w", err)
	}

	ret := quote{}
	err = json.Unmarshal(data, &ret)
	if err != nil {
		return fmt.Errorf("failed to JSON unmarshal quote: %w", err)
	}

	*q = Quote(ret)

	return nil
}

// Parse validates, parses input data into an Quote, and verifies the signature.
func Parse(data []byte) (Quote, error) {
	q := Quote{}
	if err := json.Unmarshal(data, &q); err != nil {
		return Quote{}, fmt.Errorf("failed to unmarshal Quote: %w", err)
	}

	if err := q.Verify(); err != nil {
		return Quote{}, fmt.Errorf("failed to verify Quote: %w", err)
	}

	return q, nil
}

// Create generates a new Quote with the specified parameters and options.
func Create(fromDID did.BearerDID, to, exchangeID, expiresAt string, rate string, payin, payout QuoteDetails, opts ...CreateOption) (Quote, error) {
	q := createOptions{
		id:        typeid.Must(typeid.WithPrefix(Kind)).String(),
		createdAt: time.Now(),
		protocol:  "1.0",
	}

	for _, opt := range opts {
		opt(&q)
	}

	quote := Quote{
		Metadata: message.Metadata{
			From:       fromDID.URI,
			To:         to,
			Kind:       Kind,
			ID:         q.id,
			ExchangeID: exchangeID,
			CreatedAt:  q.createdAt.UTC().Format(time.RFC3339),
			ExternalID: q.externalID,
			Protocol:   q.protocol,
		},
		Data: Data{
			ExpiresAt: expiresAt,
			Rate:      rate,
			Payin:     payin,
			Payout:    payout,
		},
	}

	signature, err := crypto.Sign(quote, fromDID)
	if err != nil {
		return Quote{}, fmt.Errorf("failed to sign quote: %w", err)
	}

	quote.Signature = signature

	return quote, nil
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
	return func(o *createOptions) {
		o.id = id
	}
}

// CreatedAt can be passed to [Create] to provide a custom created at time.
func CreatedAt(t time.Time) CreateOption {
	return func(o *createOptions) {
		o.createdAt = t
	}
}

// ExternalID can be passed to [Create] to provide a custom external id.
func ExternalID(externalID string) CreateOption {
	return func(o *createOptions) {
		o.externalID = externalID
	}
}

type quoteDetailsOptions struct {
	Fee                decimal.Decimal
}

// QuoteDetailsOption defines a type for functions that can modify the quoteDetailsOptions struct.
type QuoteDetailsOption func(*quoteDetailsOptions)

// DetailsFee is an option for [NewQuoteDetails] that allows setting a custom fee for a [QuoteDetails].
func DetailsFee(fee decimal.Decimal) QuoteDetailsOption {
	return func(q *quoteDetailsOptions) {
		q.Fee = fee
	}
}

// NewQuoteDetails creates a [QuoteDetails] object with the specified currency code, subtotal,
// and optional modifications provided through [QuoteDetailsOption] functions, such as [DetailsFee] or [DetailsInstruction].
func NewQuoteDetails(currencyCode string, subtotal decimal.Decimal, opts ...QuoteDetailsOption) QuoteDetails {
	q := quoteDetailsOptions{}
	for _, opt := range opts {
		opt(&q)
	}

	total := subtotal.Add(q.Fee)

	return QuoteDetails{
		CurrencyCode:       currencyCode,
		Subtotal:           subtotal.String(),
		Fee:                q.Fee.String(),
		Total:              total.String(),
	}
}


type quote Quote
