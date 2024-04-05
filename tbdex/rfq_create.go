package tbdex

import (
	"fmt"
	"time"

	"go.jetpack.io/typeid"
)

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

func WithRFQID(id string) CreateRFQOption {
	return func(r *createRFQOptions) {
		r.id = id
	}
}

func WithRFQCreatedAt(t time.Time) CreateRFQOption {
	return func(r *createRFQOptions) {
		r.createdAt = t
	}
}

func WithRFQClaimsHash(claimsHash string) CreateRFQOption {
	return func(r *createRFQOptions) {
		r.claimsHash = claimsHash
	}
}

type SelectedPayinMethodOption func(*SelectedPayinMethod)
type SelectedPayoutMethodOption func(*SelectedPayoutMethod)

func SelectedPayinMethodDetailsHash(detailsHash string) SelectedPayinMethodOption {
	return func(pm *SelectedPayinMethod) {
		pm.PaymentDetailsHash = detailsHash
	}
}

func SelectedPayoutMethodDetailsHash(detailsHash string) SelectedPayoutMethodOption {
	return func(pm *SelectedPayoutMethod) {
		pm.PaymentDetailsHash = detailsHash
	}
}

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

func WithRFQSelectedPayoutMethod(kind string, opts ...SelectedPayoutMethodOption) SelectedPayoutMethod {
	s := SelectedPayoutMethod{
		Kind: kind,
	}

	for _, opt := range opts {
		opt(&s)
	}

	return s
}
