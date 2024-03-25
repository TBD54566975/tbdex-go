package tbdex

import (
	"errors"
	"fmt"
	"time"

	"go.jetpack.io/typeid"
)

// Create creates an [Offering]
//
// An Offering is a resource created by a PFI to define requirements for a given currency pair offered for exchange.
//
// [Offering]: https://github.com/TBD54566975/tbdex/tree/main/specs/protocol#offering
func CreateOffering(payin OfferingPayinDetails, payout OfferingPayoutDetails, rate string, from string, opts ...CreateOption) (Offering, error) {
	defaultID, err := typeid.WithPrefix(OfferingKind)
	if err != nil {
		return Offering{}, fmt.Errorf("failed to generate default id: %w", err)
	}

	o := createOptions{
		id:          defaultID.String(),
		createdAt:   time.Now(),
		updatedAt:   time.Now(),
		description: fmt.Sprintf("%s for %s", payout.CurrencyCode, payin.CurrencyCode),
		protocol:    "1.0",
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
		ResourceMetadata: ResourceMetadata{
			Kind:      OfferingKind,
			ID:        o.id,
			From:      from,
			CreatedAt: o.createdAt.UTC().Format(time.RFC3339),
			UpdatedAt: o.updatedAt.UTC().Format(time.RFC3339),
		},
		OfferingData: OfferingData{
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
	protocol    string
}

// CreateOption implements functional options pattern for [Create].
type CreateOption func(*createOptions)

// WithID can be passed to [Create] to provide a custom id.
func WithOfferingID(id string) CreateOption {
	return func(o *createOptions) {
		o.id = id
	}
}

// WithOfferingCreatedAt can be passed to [Create] to provide a custom created at time.
func WithOfferingCreatedAt(t time.Time) CreateOption {
	return func(o *createOptions) {
		o.createdAt = t
	}
}

// WithOfferingUpdatedAt can be passed to [Create] to provide a custom updated at time.
func WithOfferingUpdatedAt(t time.Time) CreateOption {
	return func(o *createOptions) {
		o.updatedAt = t
	}
}

// WithOfferingDescription can be passed to [Create] to provide a custom description.
func WithOfferingDescription(d string) CreateOption {
	return func(o *createOptions) {
		o.description = d
	}
}

// PayinOption implements functional options pattern for [OfferingPayinDetails].
type PayinOption func(*OfferingPayinDetails)

// WithOfferingPayinMin can be passed to [Create] to provide a custom min payin amount.
func WithOfferingPayinMin(min string) PayinOption {
	return func(p *OfferingPayinDetails) {
		p.Min = min
	}
}

// WithOfferingPayinMax can be passed to [Create] to provide a custom max payin amount.
func WithOfferingPayinMax(max string) PayinOption {
	return func(p *OfferingPayinDetails) {
		p.Max = max
	}
}

// PayinMethodOption implements functional options pattern for [OfferingPayinMethod].
type PayinMethodOption func(*OfferingPayinMethod)

// WithOfferingPayinMethodFee can be passed to [Create] to provide a custom payin method fee.
func WithOfferingPayinMethodFee(fee string) PayinMethodOption {
	return func(pm *OfferingPayinMethod) {
		pm.Fee = fee
	}
}

// WithOfferingPayinMethodMin can be passed to [Create] to provide a custom min payin method amount.
func WithOfferingPayinMethodMin(min string) PayinMethodOption {
	return func(pm *OfferingPayinMethod) {
		pm.Min = min
	}
}

// WithOfferingPayinMethodMax can be passed to [Create] to provide a custom max payin method amount.
func WithOfferingPayinMethodMax(max string) PayinMethodOption {
	return func(pm *OfferingPayinMethod) {
		pm.Max = max
	}
}

// WithOfferingPayinMethodGroup can be passed to [Create] to provide a custom payin method group.
func WithOfferingPayinMethodGroup(group string) PayinMethodOption {
	return func(pm *OfferingPayinMethod) {
		pm.Group = group
	}
}

// WithOfferingPayinMethodName can be passed to [Create] to provide a custom payin method name.
func WithOfferingPayinMethodName(name string) PayinMethodOption {
	return func(pm *OfferingPayinMethod) {
		pm.Name = name
	}
}

// WithOfferingPayinMethodDescription can be passed to [Create] to provide a custom payin method description.
func WithOfferingPayinMethodDescription(description string) PayinMethodOption {
	return func(pm *OfferingPayinMethod) {
		pm.Description = description
	}
}

// WithOfferingPayinMethodRequiredPaymentDetails can be passed to [Create] to provide a custom payin method required payment details.
func WithOfferingPayinMethodRequiredPaymentDetails(details string) PayinMethodOption {
	return func(pm *OfferingPayinMethod) {
		pm.RequiredPaymentDetails = details
	}
}

// WithOfferingPayinMethod can be passed to [Create] to provide a custom payin method.
func WithOfferingPayinMethod(kind string, opts ...PayinMethodOption) OfferingPayinMethod {
	pm := OfferingPayinMethod{Kind: kind}

	for _, opt := range opts {
		opt(&pm)
	}

	return pm
}

// WithOfferingPayin can be passed to [Create] to provide a custom payin.
func WithOfferingPayin(currencyCode string, payinMethod OfferingPayinMethod, opts ...PayinOption) OfferingPayinDetails {
	p := OfferingPayinDetails{
		CurrencyCode: currencyCode,
		Methods:      []OfferingPayinMethod{payinMethod},
	}

	for _, opt := range opts {
		opt(&p)
	}

	return p
}

// PayoutOption implements functional options pattern for [OfferingPayoutDetails].
type PayoutOption func(*OfferingPayoutDetails)

// WithOfferingPayoutMin can be passed to [Create] to provide a custom min payout amount.
func WithOfferingPayoutMin(min string) PayoutOption {
	return func(p *OfferingPayoutDetails) {
		p.Min = min
	}
}

// WithOfferingPayoutMax can be passed to [Create] to provide a custom max payout amount.
func WithOfferingPayoutMax(max string) PayoutOption {
	return func(p *OfferingPayoutDetails) {
		p.Max = max
	}
}

// PayoutMethodOption implements functional options pattern for [OfferingPayoutMethod].
type PayoutMethodOption func(*OfferingPayoutMethod)

// WithOfferingPayoutMethodFee can be passed to [Create] to provide a custom payout method fee.
func WithOfferingPayoutMethodFee(fee string) PayoutMethodOption {
	return func(pm *OfferingPayoutMethod) {
		pm.Fee = fee
	}
}

// WithOfferingPayoutMethodMin can be passed to [Create] to provide a custom min payout method amount.
func WithOfferingPayoutMethodMin(min string) PayoutMethodOption {
	return func(pm *OfferingPayoutMethod) {
		pm.Min = min
	}
}

// WithOfferingPayoutMethodMax can be passed to [Create] to provide a custom max payout method amount.
func WithOfferingPayoutMethodMax(max string) PayoutMethodOption {
	return func(pm *OfferingPayoutMethod) {
		pm.Max = max
	}
}

// WithOfferingPayoutMethodGroup can be passed to [Create] to provide a custom payout method group.
func WithOfferingPayoutMethodGroup(group string) PayoutMethodOption {
	return func(pm *OfferingPayoutMethod) {
		pm.Group = group
	}
}

// WithOfferingPayoutMethodName can be passed to [Create] to provide a custom payout method name.
func WithOfferingPayoutMethodName(name string) PayoutMethodOption {
	return func(pm *OfferingPayoutMethod) {
		pm.Name = name
	}
}

// WithOfferingPayoutMethodDescription can be passed to [Create] to provide a custom payout method description.
func WithOfferingPayoutMethodDescription(description string) PayoutMethodOption {
	return func(pm *OfferingPayoutMethod) {
		pm.Description = description
	}
}

// WithOfferingPayoutMethod can be passed to [Create] to provide a custom payout method.
func WithOfferingPayoutMethod(kind string, estimatedSettlementTime time.Duration, opts ...PayoutMethodOption) OfferingPayoutMethod {
	pm := OfferingPayoutMethod{
		OfferingPaymentMethod:   OfferingPaymentMethod{Kind: kind},
		EstimatedSettlementTime: uint64(estimatedSettlementTime.Abs().Seconds()),
	}

	for _, opt := range opts {
		opt(&pm)
	}

	return pm
}

// WithOfferingPayout can be passed to [Create] to provide a custom payout.
func WithOfferingPayout(currencyCode string, payoutMethod OfferingPayoutMethod, opts ...PayoutOption) OfferingPayoutDetails {
	p := OfferingPayoutDetails{
		CurrencyCode: currencyCode,
		Methods:      []OfferingPayoutMethod{payoutMethod},
	}

	for _, opt := range opts {
		opt(&p)
	}

	return p
}
