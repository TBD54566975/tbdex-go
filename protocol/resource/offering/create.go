package offering

import (
	"errors"
	"fmt"
	"time"

	"github.com/TBD54566975/tbdex-go/protocol/resource"
	"go.jetpack.io/typeid"
)

// Create creates an [Offering]
//
// An Offering is a resource created by a PFI to define requirements for a given currency pair offered for exchange.
//
// [Offering]: https://github.com/TBD54566975/tbdex/tree/main/specs/protocol#offering
func Create(payin PayinDetails, payout PayoutDetails, rate string, from string, opts ...CreateOption) (Offering, error) {
	defaultID, err := typeid.WithPrefix(Kind)
	if err != nil {
		return Offering{}, fmt.Errorf("failed to generate default id: %w", err)
	}

	o := createOptions{
		id:          defaultID.String(),
		createdAt:   time.Now(),
		updatedAt:   time.Now(),
		description: fmt.Sprintf("%s for %s", payout.CurrencyCode, payin.CurrencyCode),
		protocol: "1.0",
	}

	for _, opt := range opts {
		opt(&o)
	}

	if len(payin.Methods) == 0 {
		return Offering{}, errors.New("1 payin method is required")
	}

	if len(payout.Methods) == 0 {
		return Offering{}, errors.New("1 payout method is required")
	}

	return Offering{
		Metadata: resource.Metadata{
			Kind:      Kind,
			ID:        o.id,
			From: from,
			CreatedAt: o.createdAt.UTC().Format(time.RFC3339),
			UpdatedAt: o.updatedAt.UTC().Format(time.RFC3339),
		},
		Data: Data{
			Payin:       payin,
			Payout:      payout,
			Rate:        rate,
			Description: o.description,
		},
	}, nil
}

type createOptions struct {
	id          string
	createdAt   time.Time
	updatedAt   time.Time
	description string
	protocol string
}

// CreateOption implements functional options pattern for [Create].
type CreateOption func(*createOptions)

// WithID can be passed to [Create] to provide a custom id.
func WithID(id string) CreateOption {
	return func(o *createOptions) {
		o.id = id
	}
}

// WithCreatedAt can be passed to [Create] to provide a custom created at time.
func WithCreatedAt(t time.Time) CreateOption {
	return func(o *createOptions) {
		o.createdAt = t
	}
}

// WithUpdatedAt can be passed to [Create] to provide a custom updated at time.
func WithUpdatedAt(t time.Time) CreateOption {
	return func(o *createOptions) {
		o.updatedAt = t
	}
}

// WithDescription can be passed to [Create] to provide a custom description.
func WithDescription(d string) CreateOption {
	return func(o *createOptions) {
		o.description = d
	}
}

// PayinOption implements functional options pattern for [PayinDetails].
type PayinOption func(*PayinDetails)

// WithPayinMin can be passed to [Create] to provide a custom min payin amount.
func WithPayinMin(min string) PayinOption {
	return func(p *PayinDetails) {
		p.Min = min
	}
}

// WithPayinMax can be passed to [Create] to provide a custom max payin amount.
func WithPayinMax(max string) PayinOption {
	return func(p *PayinDetails) {
		p.Max = max
	}
}

// PayinMethodOption implements functional options pattern for [PayinMethod].
type PayinMethodOption func(*PayinMethod)

// WithPayinMethodFee can be passed to [Create] to provide a custom payin method fee.
func WithPayinMethodFee(fee string) PayinMethodOption {
	return func(pm *PayinMethod) {
		pm.Fee = fee
	}
}

// WithPayinMethodMin can be passed to [Create] to provide a custom min payin method amount.
func WithPayinMethodMin(min string) PayinMethodOption {
	return func(pm *PayinMethod) {
		pm.Min = min
	}
}

// WithPayinMethodMax can be passed to [Create] to provide a custom max payin method amount.
func WithPayinMethodMax(max string) PayinMethodOption {
	return func(pm *PayinMethod) {
		pm.Max = max
	}
}

// WithPayinMethodGroup can be passed to [Create] to provide a custom payin method group.
func WithPayinMethodGroup(group string) PayinMethodOption {
	return func(pm *PayinMethod) {
		pm.Group = group
	}
}

// WithPayinMethodName can be passed to [Create] to provide a custom payin method name.
func WithPayinMethodName(name string) PayinMethodOption {
	return func(pm *PayinMethod) {
		pm.Name = name
	}
}

// WithPayinMethodDescription can be passed to [Create] to provide a custom payin method description.
func WithPayinMethodDescription(description string) PayinMethodOption {
	return func(pm *PayinMethod) {
		pm.Description = description
	}
}

// WithPayinMethodRequiredPaymentDetails can be passed to [Create] to provide a custom payin method required payment details.
func WithPayinMethodRequiredPaymentDetails(details string) PayinMethodOption {
	return func(pm *PayinMethod) {
		pm.RequiredPaymentDetails = details
	}
}

// WithPayinMethod can be passed to [Create] to provide a custom payin method.
func WithPayinMethod(kind string, opts ...PayinMethodOption) PayinMethod {
	pm := PayinMethod{Kind: kind}

	for _, opt := range opts {
		opt(&pm)
	}

	return pm
}

// WithPayin can be passed to [Create] to provide a custom payin.
func WithPayin(currencyCode string, payinMethod PayinMethod, opts ...PayinOption) PayinDetails {
	p := PayinDetails{
		CurrencyCode: currencyCode,
		Methods:      []PayinMethod{payinMethod},
	}

	for _, opt := range opts {
		opt(&p)
	}

	return p
}

// PayoutOption implements functional options pattern for [PayoutDetails].
type PayoutOption func(*PayoutDetails)

// WithPayoutMin can be passed to [Create] to provide a custom min payout amount.
func WithPayoutMin(min string) PayoutOption {
	return func(p *PayoutDetails) {
		p.Min = min
	}
}

// WithPayoutMax can be passed to [Create] to provide a custom max payout amount.
func WithPayoutMax(max string) PayoutOption {
	return func(p *PayoutDetails) {
		p.Max = max
	}
}

// PayoutMethodOption implements functional options pattern for [PayoutMethod].
type PayoutMethodOption func(*PayoutMethod)

// WithPayoutMethodFee can be passed to [Create] to provide a custom payout method fee.
func WithPayoutMethodFee(fee string) PayoutMethodOption {
	return func(pm *PayoutMethod) {
		pm.Fee = fee
	}
}

// WithPayoutMethodMin can be passed to [Create] to provide a custom min payout method amount.
func WithPayoutMethodMin(min string) PayoutMethodOption {
	return func(pm *PayoutMethod) {
		pm.Min = min
	}
}

// WithPayoutMethodMax can be passed to [Create] to provide a custom max payout method amount.
func WithPayoutMethodMax(max string) PayoutMethodOption {
	return func(pm *PayoutMethod) {
		pm.Max = max
	}
}

// WithPayoutMethodGroup can be passed to [Create] to provide a custom payout method group.
func WithPayoutMethodGroup(group string) PayoutMethodOption {
	return func(pm *PayoutMethod) {
		pm.Group = group
	}
}

// WithPayoutMethodName can be passed to [Create] to provide a custom payout method name.
func WithPayoutMethodName(name string) PayoutMethodOption {
	return func(pm *PayoutMethod) {
		pm.Name = name
	}
}

// WithPayoutMethodDescription can be passed to [Create] to provide a custom payout method description.
func WithPayoutMethodDescription(description string) PayoutMethodOption {
	return func(pm *PayoutMethod) {
		pm.Description = description
	}
}

// WithPayoutMethod can be passed to [Create] to provide a custom payout method.
func WithPayoutMethod(kind string, estimatedSettlementTime time.Duration, opts ...PayoutMethodOption) PayoutMethod {
	pm := PayoutMethod{
		PaymentMethod:           PaymentMethod{Kind: kind},
		EstimatedSettlementTime: uint64(estimatedSettlementTime.Abs().Seconds()),
	}

	for _, opt := range opts {
		opt(&pm)
	}

	return pm
}

// WithPayout can be passed to [Create] to provide a custom payout.
func WithPayout(currencyCode string, payoutMethod PayoutMethod, opts ...PayoutOption) PayoutDetails {
	p := PayoutDetails{
		CurrencyCode: currencyCode,
		Methods:      []PayoutMethod{payoutMethod},
	}

	for _, opt := range opts {
		opt(&p)
	}

	return p
}
