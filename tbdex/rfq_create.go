package tbdex

import (
	"fmt"
	"time"

	"go.jetpack.io/typeid"
)

// CreateRFQ creates an [RFQ]
//
// An RFQ is a resource created by a customer of the PFI to request a quote
//
// [Offering]: https://github.com/TBD54566975/tbdex/tree/main/specs/protocol#rfq-request-for-quote
func CreateRFQ(from, to, offeringID string, payin SelectedPayinMethod, payout SelectedPayoutMethod, opts ...CreateRFQOption) (RFQ, error) {
	defaultID, err := typeid.WithPrefix(RFQKind)
	if err != nil {
		return RFQ{}, fmt.Errorf("failed to generate default id: %w", err)
	}

	r := createRFQOptions{
		id:          defaultID.String(),
		createdAt:   time.Now(),
		protocol:    "1.0",
		externalID: "",
		claimsHash: "",
	}

	for _, opt := range opts {
		opt(&r)
	}

	return RFQ{
		MessageMetadata: MessageMetadata{
			From: from,
			To: to,
			Kind:      RFQKind,
			ID:        r.id,
			ExchangeID: r.id,
			CreatedAt: r.createdAt.UTC().Format(time.RFC3339),
			ExternalID: r.externalID,
			Protocol:  r.protocol,
		},
		RFQData: RFQData{
				OfferingID: offeringID,
				Payin: payin,
				Payout: payout,
				ClaimsHash: r.claimsHash,
		},
	}, nil
}

// CreateRFQOption is a function type used to apply options to RFQ creation.
type createRFQOptions struct {
	id         string
	createdAt  time.Time
	claimsHash string
	protocol   string
	externalID string
}

// CreateRFQOption is a function type used to apply options to RFQ creation.
type CreateRFQOption func(*createRFQOptions)

// WithRFQID can be passed to [CreateRFQ] to provide a custom id.
func WithRFQID(id string) CreateRFQOption {
	return func(r *createRFQOptions) {
		r.id = id
	}
}

// WithRFQCreatedAt can be passed to [CreateRFQ] to provide a custom created at time.
func WithRFQCreatedAt(t time.Time) CreateRFQOption {
	return func(r *createRFQOptions) {
		r.createdAt = t
	}
}

// WithRFQClaimsHash can be passed to [CreateRFQ] to provide a hash of the claims data.
func WithRFQClaimsHash(claimsHash string) CreateRFQOption {
	return func(r *createRFQOptions) {
		r.claimsHash = claimsHash
	}
}

// WithRFQExternalID can be passed to [CreateRFQ] to provide a custom external id.
func WithRFQExternalID(externalID string) CreateRFQOption {
	return func(r *createRFQOptions) {
		r.externalID = externalID
	}
}

// SelectedPayinMethodOption is a function type used to apply options to [SelectedPayinMethod] creation.
type SelectedPayinMethodOption func(*SelectedPayinMethod)

// SelectedPayoutMethodOption is a function type used to apply options to [SelectedPayoutMethod] creation.
type SelectedPayoutMethodOption func(*SelectedPayoutMethod)

// WithSelectedPayinMethodDetailsHash can be passed to [WithRFQSelectedPayinMethod] to provide a hash of the payment details.
func WithSelectedPayinMethodDetailsHash(detailsHash string) SelectedPayinMethodOption {
	return func(pm *SelectedPayinMethod) {
		pm.PaymentDetailsHash = detailsHash
	}
}

// WithSelectedPayoutMethodDetailsHash can be passed to [WithRFQSelectedPayoutMethod] to provide a hash of the payment details.
func WithSelectedPayoutMethodDetailsHash(detailsHash string) SelectedPayoutMethodOption {
	return func(pm *SelectedPayoutMethod) {
		pm.PaymentDetailsHash = detailsHash
	}
}

// WithRFQSelectedPayinMethod can be passed to [Create] to provide a payin method.
func WithRFQSelectedPayinMethod(amount, kind string, opts ...SelectedPayinMethodOption) SelectedPayinMethod {
	s := SelectedPayinMethod{
		Amount: amount,
		Kind: kind,
	}

	for _, opt := range opts {
		opt(&s)
	}

	return s
}

// WithRFQSelectedPayoutMethod can be passed to [Create] to provide a payout method.
func WithRFQSelectedPayoutMethod(kind string, opts ...SelectedPayoutMethodOption) SelectedPayoutMethod {
	s := SelectedPayoutMethod{
		Kind: kind,
	}

	for _, opt := range opts {
		opt(&s)
	}

	return s
}
