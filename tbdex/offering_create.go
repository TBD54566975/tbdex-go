package tbdex

import (
	"errors"
	"fmt"
	"time"

	"go.jetpack.io/typeid"
)

// CreateOffering creates an [Offering]
//
// An Offering is a resource created by a PFI to define requirements for a given currency pair offered for exchange.
//
// [Offering]: https://github.com/TBD54566975/tbdex/tree/main/specs/protocol#offering
func CreateOffering(payin OfferingPayinDetails, payout OfferingPayoutDetails, rate string, opts ...CreateOfferingOption) (Offering, error) {
	defaultID, err := typeid.WithPrefix(OfferingKind)
	if err != nil {
		return Offering{}, fmt.Errorf("failed to generate default id: %w", err)
	}

	o := createOfferingOptions{
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
			CreatedAt: o.createdAt.UTC().Format(time.RFC3339),
			UpdatedAt: o.updatedAt.UTC().Format(time.RFC3339),
			Protocol:  o.protocol,
		},
		OfferingData: OfferingData{
			Payin:       payin,
			Payout:      payout,
			Rate:        rate,
			Description: o.description,
		},
	}, nil
}

type createOfferingOptions struct {
	id          string
	createdAt   time.Time
	updatedAt   time.Time
	description string
	protocol    string
}

// CreateOfferingOption implements functional options pattern for [Create].
type CreateOfferingOption func(*createOfferingOptions)

// WithOfferingID can be passed to [CreateOffering] to provide a custom id.
func WithOfferingID(id string) CreateOfferingOption {
	return func(o *createOfferingOptions) {
		o.id = id
	}
}

// WithOfferingCreatedAt can be passed to [Create] to provide a custom created at time.
func WithOfferingCreatedAt(t time.Time) CreateOfferingOption {
	return func(o *createOfferingOptions) {
		o.createdAt = t
	}
}

// WithOfferingUpdatedAt can be passed to [Create] to provide a custom updated at time.
func WithOfferingUpdatedAt(t time.Time) CreateOfferingOption {
	return func(o *createOfferingOptions) {
		o.updatedAt = t
	}
}

// WithOfferingDescription can be passed to [Create] to provide a custom description.
func WithOfferingDescription(d string) CreateOfferingOption {
	return func(o *createOfferingOptions) {
		o.description = d
	}
}

// OfferingPayinOption implements functional options pattern for [OfferingPayinDetails].
type OfferingPayinOption func(*OfferingPayinDetails)

// WithOfferingPayinMin can be passed to [Create] to provide a custom min payin amount.
func WithOfferingPayinMin(min string) OfferingPayinOption {
	return func(p *OfferingPayinDetails) {
		p.Min = min
	}
}

// WithOfferingPayinMax can be passed to [Create] to provide a custom max payin amount.
func WithOfferingPayinMax(max string) OfferingPayinOption {
	return func(p *OfferingPayinDetails) {
		p.Max = max
	}
}

// OfferingPayinMethodOption implements functional options pattern for [OfferingPayinMethod].
type OfferingPayinMethodOption func(*OfferingPayinMethod)

// WithOfferingPayinMethodFee can be passed to [Create] to provide a custom payin method fee.
func WithOfferingPayinMethodFee(fee string) OfferingPayinMethodOption {
	return func(pm *OfferingPayinMethod) {
		pm.Fee = fee
	}
}

// WithOfferingPayinMethodMin can be passed to [Create] to provide a custom min payin method amount.
func WithOfferingPayinMethodMin(min string) OfferingPayinMethodOption {
	return func(pm *OfferingPayinMethod) {
		pm.Min = min
	}
}

// WithOfferingPayinMethodMax can be passed to [Create] to provide a custom max payin method amount.
func WithOfferingPayinMethodMax(max string) OfferingPayinMethodOption {
	return func(pm *OfferingPayinMethod) {
		pm.Max = max
	}
}

// WithOfferingPayinMethodGroup can be passed to [Create] to provide a custom payin method group.
func WithOfferingPayinMethodGroup(group string) OfferingPayinMethodOption {
	return func(pm *OfferingPayinMethod) {
		pm.Group = group
	}
}

// WithOfferingPayinMethodName can be passed to [Create] to provide a custom payin method name.
func WithOfferingPayinMethodName(name string) OfferingPayinMethodOption {
	return func(pm *OfferingPayinMethod) {
		pm.Name = name
	}
}

// WithOfferingPayinMethodDescription can be passed to [Create] to provide a custom payin method description.
func WithOfferingPayinMethodDescription(description string) OfferingPayinMethodOption {
	return func(pm *OfferingPayinMethod) {
		pm.Description = description
	}
}

// WithOfferingPayinMethodRequiredPaymentDetails can be passed to [Create] to provide a custom payin method required payment details.
func WithOfferingPayinMethodRequiredPaymentDetails(details string) OfferingPayinMethodOption {
	return func(pm *OfferingPayinMethod) {
		pm.RequiredPaymentDetails = details
	}
}

// WithOfferingPayinMethod can be passed to [Create] to provide a custom payin method.
func WithOfferingPayinMethod(kind string, opts ...OfferingPayinMethodOption) OfferingPayinMethod {
	pm := OfferingPayinMethod{Kind: kind}

	for _, opt := range opts {
		opt(&pm)
	}

	return pm
}

// WithOfferingPayin can be passed to [Create] to provide a custom payin.
func WithOfferingPayin(currencyCode string, payinMethod OfferingPayinMethod, opts ...OfferingPayinOption) OfferingPayinDetails {
	p := OfferingPayinDetails{
		CurrencyCode: currencyCode,
		Methods:      []OfferingPayinMethod{payinMethod},
	}

	for _, opt := range opts {
		opt(&p)
	}

	return p
}

// OfferingPayoutOption implements functional options pattern for [OfferingPayoutDetails].
type OfferingPayoutOption func(*OfferingPayoutDetails)

// WithOfferingPayoutMin can be passed to [Create] to provide a custom min payout amount.
func WithOfferingPayoutMin(min string) OfferingPayoutOption {
	return func(p *OfferingPayoutDetails) {
		p.Min = min
	}
}

// WithOfferingPayoutMax can be passed to [Create] to provide a custom max payout amount.
func WithOfferingPayoutMax(max string) OfferingPayoutOption {
	return func(p *OfferingPayoutDetails) {
		p.Max = max
	}
}

// OfferingPayoutMethodOption implements functional options pattern for [OfferingPayoutMethod].
type OfferingPayoutMethodOption func(*OfferingPayoutMethod)

// WithOfferingPayoutMethodFee can be passed to [Create] to provide a custom payout method fee.
func WithOfferingPayoutMethodFee(fee string) OfferingPayoutMethodOption {
	return func(pm *OfferingPayoutMethod) {
		pm.Fee = fee
	}
}

// WithOfferingPayoutMethodMin can be passed to [Create] to provide a custom min payout method amount.
func WithOfferingPayoutMethodMin(min string) OfferingPayoutMethodOption {
	return func(pm *OfferingPayoutMethod) {
		pm.Min = min
	}
}

// WithOfferingPayoutMethodMax can be passed to [Create] to provide a custom max payout method amount.
func WithOfferingPayoutMethodMax(max string) OfferingPayoutMethodOption {
	return func(pm *OfferingPayoutMethod) {
		pm.Max = max
	}
}

// WithOfferingPayoutMethodGroup can be passed to [Create] to provide a custom payout method group.
func WithOfferingPayoutMethodGroup(group string) OfferingPayoutMethodOption {
	return func(pm *OfferingPayoutMethod) {
		pm.Group = group
	}
}

// WithOfferingPayoutMethodName can be passed to [Create] to provide a custom payout method name.
func WithOfferingPayoutMethodName(name string) OfferingPayoutMethodOption {
	return func(pm *OfferingPayoutMethod) {
		pm.Name = name
	}
}

// WithOfferingPayoutMethodDescription can be passed to [Create] to provide a custom payout method description.
func WithOfferingPayoutMethodDescription(description string) OfferingPayoutMethodOption {
	return func(pm *OfferingPayoutMethod) {
		pm.Description = description
	}
}

// WithOfferingPayoutMethod can be passed to [Create] to provide a custom payout method.
func WithOfferingPayoutMethod(kind string, estimatedSettlementTime time.Duration, opts ...OfferingPayoutMethodOption) OfferingPayoutMethod {
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
func WithOfferingPayout(currencyCode string, payoutMethod OfferingPayoutMethod, opts ...OfferingPayoutOption) OfferingPayoutDetails {
	p := OfferingPayoutDetails{
		CurrencyCode: currencyCode,
		Methods:      []OfferingPayoutMethod{payoutMethod},
	}

	for _, opt := range opts {
		opt(&p)
	}

	return p
}
