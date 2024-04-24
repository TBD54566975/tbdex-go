package quote

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/TBD54566975/tbdex-go/tbdex"
	"github.com/tbd54566975/web5-go/dids/did"
	"go.jetpack.io/typeid"
)

// Kind identifies this message kind
const Kind = "quote"

// Quote represents a quote message within the exchange.
type Quote struct {
	MessageMetadata tbdex.MessageMetadata `json:"metadata"`
	Data            Data                  `json:"data"`
	Signature       string                `json:"signature"`
}

// Data encapsulates the data content of a  quote.
type Data struct {
	ExpiresAt string       `json:"expiresAt"`
	Payin     QuoteDetails `json:"payin"`
	Payout    QuoteDetails `json:"payout"`
}

// QuoteDetails describes the relevant information of a currency that is being sent or received
type QuoteDetails struct {
	CurrencyCode       string              `json:"currencyCode"`
	Amount             string              `json:"amount"`
	Fee                string              `json:"fee,omitempty"`
	PaymentInstruction *PaymentInstruction `json:"paymentInstruction,omitempty"`
}

// PaymentInstruction contains instructions with plain text and/or a link
type PaymentInstruction struct {
	Link        string `json:"link,omitempty"`
	Instruction string `json:"instruction,omitempty"`
}

// Digest computes a hash of the quote
func (q Quote) Digest() ([]byte, error) {
	payload := map[string]any{"metadata": q.MessageMetadata, "data": q.Data}

	hashed, err := tbdex.DigestJSON(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to digest quote: %w", err)
	}

	return hashed, nil
}

// Sign cryptographically signs the quote using DID's private key
func (q *Quote) Sign(bearerDID did.BearerDID) error {
	q.MessageMetadata.From = bearerDID.URI

	signature, err := tbdex.Sign(q, bearerDID)
	if err != nil {
		return fmt.Errorf("failed to sign rfq: %w", err)
	}

	q.Signature = signature

	return nil
}

// Verify verifies the signature of the quote.
func (q *Quote) Verify() error {
	decoded, err := tbdex.VerifySignature(q, q.Signature)
	if err != nil {
		return fmt.Errorf("failed to verify quote signature: %w", err)
	}

	if decoded.SignerDID.URI != q.MessageMetadata.From {
		return fmt.Errorf("signer: %s does not match message metadata from: %s", decoded.SignerDID.URI, q.MessageMetadata.From)
	}

	return nil
}

// UnmarshalJSON validates and unmarshals the input data into an Quote.
func (q *Quote) UnmarshalJSON(data []byte) error {
	err := tbdex.Validate(tbdex.TypeMessage, data, tbdex.WithKind(Kind))
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

// Parse validates, parses input data into an RFQ, and verifies the signature and private data.
func (q *Quote) Parse(data []byte, privateDataStrict bool) error {
	err := q.UnmarshalJSON(data)
	if err != nil {
		return fmt.Errorf("failed to unmarshal Quote: %w", err)
	}

	err = q.Verify()
	if err != nil {
		return fmt.Errorf("failed to verify Quote: %w", err)
	}

	return nil
}

// Create generates a new Quote with the specified parameters and options.
func Create(from, to, exchangeID, expiresAt string, payin, payout QuoteDetails, opts ...CreateOption) Quote {
	q := createOptions{
		id:        typeid.Must(typeid.WithPrefix(Kind)).String(),
		createdAt: time.Now(),
		protocol:  "1.0",
	}

	return Quote{
		MessageMetadata: tbdex.MessageMetadata{
			From:       from,
			To:         to,
			Kind:       Kind,
			ID:         q.id,
			ExchangeID: exchangeID,
			CreatedAt:  q.createdAt.UTC().Format(time.RFC3339),
			ExternalID: q.externalID,
			Protocol:   q.protocol,
		},
		Data: Data{ExpiresAt: expiresAt, Payin: payin, Payout: payout},
	}
}

type createOptions struct {
	id         string
	createdAt  time.Time
	protocol   string
	externalID string
	exchangeID string
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

// ExchangeID can be passed to [Create] to provide a custom exchange id.
func ExchangeID(exchangeID string) CreateOption {
	return func(q *createOptions) {
		q.exchangeID = exchangeID
	}
}

type quoteDetailsOptions struct {
	Fee                string
	PaymentInstruction *PaymentInstruction
}

// QuoteDetailsOption defines a type for functions that can modify the quoteDetailsOptions struct.
type QuoteDetailsOption func(*quoteDetailsOptions)

// DetailsFee is an option for [NewQuoteDetails] that allows setting a custom fee for a [QuoteDetails].
func DetailsFee(fee string) QuoteDetailsOption {
	return func(q *quoteDetailsOptions) {
		q.Fee = fee
	}
}

// DetailsInstruction is an option for NewQuoteDetails that allows setting a custom [PaymentInstruction]
// for a [QuoteDetails].
func DetailsInstruction(p *PaymentInstruction) QuoteDetailsOption {
	return func(q *quoteDetailsOptions) {
		q.PaymentInstruction = p
	}
}

// NewQuoteDetails creates a [QuoteDetails] object with the specified currency code, amount,
// and optional modifications provided through [QuoteDetailsOption] functions.
func NewQuoteDetails(currencyCode string, amount string, opts ...QuoteDetailsOption) QuoteDetails {
	q := quoteDetailsOptions{}
	for _, opt := range opts {
		opt(&q)
	}
	return QuoteDetails{
		CurrencyCode:       currencyCode,
		Amount:             amount,
		Fee:                q.Fee,
		PaymentInstruction: q.PaymentInstruction,
	}
}

type paymentInstructionOptions struct {
	Link        string
	Instruction string
}

// PaymentInstructionOptions defines a type for functions that can modify the paymentInstructionOptions struct.
type PaymentInstructionOptions func(*paymentInstructionOptions)

// InstructionLink is an option for [NewPaymentInstruction] that allows setting a custom link.
func InstructionLink(link string) PaymentInstructionOptions {
	return func(p *paymentInstructionOptions) {
		p.Link = link
	}
}

// Instruction is an option for [NewPaymentInstruction] that allows setting custom text.
func Instruction(instruction string) PaymentInstructionOptions {
	return func(p *paymentInstructionOptions) {
		p.Instruction = instruction
	}
}

// NewPaymentInstruction creates a new [PaymentInstruction] using the provided options.
func NewPaymentInstruction(opts ...PaymentInstructionOptions) *PaymentInstruction {
	p := paymentInstructionOptions{}
	for _, opt := range opts {
		opt(&p)
	}
	return &PaymentInstruction{
		Link:        p.Link,
		Instruction: p.Instruction,
	}
}

type quote Quote
