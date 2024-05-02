package rfq

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/TBD54566975/tbdex-go/tbdex"
	"github.com/tbd54566975/web5-go/crypto"
	"github.com/tbd54566975/web5-go/dids/did"
	"go.jetpack.io/typeid"
)

// Create creates an [RFQ]
//
// # An RFQ is a resource created by a customer of the PFI to request a quote
//
// [RFQ]: https://github.com/TBD54566975/tbdex/tree/main/specs/protocol#rfq-request-for-quote
func Create(fromDID did.BearerDID, to, offeringID string, payin PayinMethod, payout PayoutMethod, opts ...CreateOption) (RFQ, error) {
	r := createOptions{
		id:        typeid.Must(typeid.WithPrefix(Kind)).String(),
		createdAt: time.Now(),
		protocol:  "1.0",
	}

	for _, opt := range opts {
		opt(&r)
	}

	randomBytes, err := crypto.GenerateEntropy(crypto.Entropy128)
	if err != nil {
		return RFQ{}, err
	}

	salt := base64.RawURLEncoding.EncodeToString(randomBytes)
	privateData := PrivateData{}

	scrubbedPayin, err := payin.Scrub(salt, &privateData)
	if err != nil {
		return RFQ{}, fmt.Errorf("failed to scrub payin: %w", err)
	}

	scrubbedPayout, err := payout.Scrub(salt, &privateData)
	if err != nil {
		return RFQ{}, fmt.Errorf("failed to scrub payout: %w", err)
	}

	scrubbedClaims, err := r.claims.Scrub(salt, &privateData)
	if err != nil {
		return RFQ{}, fmt.Errorf("failed to scrub claims: %w", err)
	}

	rfq := RFQ{
		MessageMetadata: tbdex.MessageMetadata{
			From:       fromDID.URI,
			To:         to,
			Kind:       Kind,
			ID:         r.id,
			ExchangeID: r.id,
			CreatedAt:  r.createdAt.UTC().Format(time.RFC3339),
			ExternalID: r.externalID,
			Protocol:   r.protocol,
		},
		Data: Data{
			OfferingID: offeringID,
			Payin:      scrubbedPayin,
			Payout:     scrubbedPayout,
			ClaimsHash: scrubbedClaims,
		},
	}

	if !privateData.IsZero() {
		privateData.Salt = salt
		rfq.PrivateData = &privateData
	}

	signature, err := tbdex.Sign(rfq, fromDID)
	if err != nil {
		return RFQ{}, fmt.Errorf("failed to sign rfq: %w", err)
	}

	rfq.Signature = signature

	return rfq, nil
}

type createOptions struct {
	id         string
	createdAt  time.Time
	protocol   string
	externalID string
	claims     ClaimsSet
}

// CreateOption is a function type used to apply options to RFQ creation.
type CreateOption func(*createOptions)

// ID can be passed to [Create] to provide a custom id.
func ID(id string) CreateOption {
	return func(r *createOptions) {
		r.id = id
	}
}

// CreatedAt can be passed to [Create] to provide a custom created at time.
func CreatedAt(t time.Time) CreateOption {
	return func(r *createOptions) {
		r.createdAt = t
	}
}

// ExternalID can be passed to [Create] to provide a custom external id.
func ExternalID(externalID string) CreateOption {
	return func(r *createOptions) {
		r.externalID = externalID
	}
}

// Claims can be passed to [Create] to provide claims if required by the offering.
func Claims(claims []string) CreateOption {
	return func(r *createOptions) {
		r.claims = claims
	}
}

type paymentMethodOptions struct {
	details map[string]any
}

// PaymentMethodOption is a function type used to apply options to payment methods.
type PaymentMethodOption func(*paymentMethodOptions)

// PaymentDetails can be passed to [Payin] or [Payout] to provide arbitrary payment details.
func PaymentDetails(details map[string]any) PaymentMethodOption {
	return func(pm *paymentMethodOptions) {
		pm.details = details
	}
}

// Payin can be passed to [Create] to provide a payin method.
func Payin(amount, kind string, opts ...PaymentMethodOption) PayinMethod {
	s := PayinMethod{Amount: amount, Kind: kind}

	o := paymentMethodOptions{}
	for _, opt := range opts {
		opt(&o)
	}

	s.PaymentDetails = o.details

	return s
}

// Payout can be passed to [Create] to provide a payout method.
func Payout(kind string, opts ...PaymentMethodOption) PayoutMethod {
	s := PayoutMethod{Kind: kind}

	o := paymentMethodOptions{}
	for _, opt := range opts {
		opt(&o)
	}

	s.PaymentDetails = o.details

	return s
}

func computeHash(salt string, data any) (string, error) {
	byteArray, err := tbdex.DigestJSON([]any{salt, data})
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(byteArray), nil
}
