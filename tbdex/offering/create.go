package offering

import (
	"errors"
	"fmt"
	"time"

	"github.com/TBD54566975/tbdex-go/tbdex"
	"github.com/tbd54566975/web5-go/pexv2"
	"go.jetpack.io/typeid"
)

// Create creates an [Offering]
//
// An Offering is a resource created by a PFI to define requirements for a given currency pair offered for exchange.
//
// [Offering]: https://github.com/TBD54566975/tbdex/tree/main/specs/protocol#offering
func Create(payin PayinDetails, payout PayoutDetails, rate string, opts ...CreateOption) (Offering, error) {
	defaultID, err := typeid.WithPrefix(Kind)
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
		ResourceMetadata: tbdex.ResourceMetadata{
			Kind:      Kind,
			ID:        o.id,
			CreatedAt: o.createdAt.UTC().Format(time.RFC3339),
			UpdatedAt: o.updatedAt.UTC().Format(time.RFC3339),
			Protocol:  o.protocol,
		},
		Data: Data{
			Payin:       payin,
			Payout:      payout,
			Rate:        rate,
			Description: o.description,
		},
	}, nil
}

func NewPayin(currencyCode string, methods []PayinMethod, opts ...PaymentOption) PayinDetails {
	return PayinDetails{
		CurrencyCode: currencyCode,
		Methods:      methods,
	}
}

func NewPayinMethod(kind string, opts ...PaymentMethodOption) PayinMethod {
	o := paymentMethodOptions{}
	for _, opt := range opts {
		opt(&o)
	}

	return PayinMethod{
		Kind:        kind,
		Min:         o.Min,
		Max:         o.Max,
		Group:       o.Group,
		Fee:         o.Fee,
		Name:        o.Name,
		Description: o.Description,
	}
}

func NewPayout(currencyCode string, methods []PayoutMethod, opts ...PaymentOption) PayoutDetails {
	o := paymentOptions{}
	for _, opt := range opts {
		opt(&o)
	}

	return PayoutDetails{
		CurrencyCode: currencyCode,
		Min:          o.Min,
		Max:          o.Max,
		Methods:      methods,
	}
}

func NewPayoutMethod(kind string, estimatedSettlementTime time.Duration, opts ...PaymentMethodOption) PayoutMethod {
	o := paymentMethodOptions{}
	for _, opt := range opts {
		opt(&o)
	}

	pm := PayoutMethod{
		Kind:                    kind,
		Min:                     o.Min,
		Max:                     o.Max,
		Group:                   o.Group,
		Fee:                     o.Fee,
		Name:                    o.Name,
		Description:             o.Description,
		RequiredPaymentDetails:  o.RequiredPaymentDetails,
		EstimatedSettlementTime: uint64(estimatedSettlementTime.Abs().Seconds()),
	}

	return pm
}

type createOptions struct {
	id             string
	createdAt      time.Time
	updatedAt      time.Time
	description    string
	protocol       string
	requiredClaims pexv2.PresentationDefinition
}

// CreateOption implements functional options pattern for [Create].
type CreateOption func(*createOptions)

// WithID can be passed to [CreateOffering] to provide a custom id.
func ID(id string) CreateOption {
	return func(o *createOptions) {
		o.id = id
	}
}

// WithCreatedAt can be passed to [Create] to provide a custom created at time.
func CreatedAt(t time.Time) CreateOption {
	return func(o *createOptions) {
		o.createdAt = t
	}
}

// WithUpdatedAt can be passed to [Create] to provide a custom updated at time.
func UpdatedAt(t time.Time) CreateOption {
	return func(o *createOptions) {
		o.updatedAt = t
	}
}

// WithDescription can be passed to [Create] to provide a custom description.
func Description(d string) CreateOption {
	return func(o *createOptions) {
		o.description = d
	}
}

type paymentOptions struct {
	Min string
	Max string
}

// OfferingPayinOption implements functional options pattern for [PayinDetails].
type PaymentOption func(*paymentOptions)

// WithPayinMin can be passed to [Create] to provide a custom min payin amount.
func WithMin(min string) PaymentOption {
	return func(p *paymentOptions) {
		p.Min = min
	}
}

// WithPayinMax can be passed to [Create] to provide a custom max payin amount.
func WithMax(max string) PaymentOption {
	return func(p *paymentOptions) {
		p.Max = max
	}
}

type paymentMethodOptions struct {
	Min                    string
	Max                    string
	Group                  string
	Fee                    string
	Name                   string
	Description            string
	RequiredPaymentDetails string
}

// PaymentMethodOption implements functional options pattern for [PayinMethod].
type PaymentMethodOption func(*paymentMethodOptions)

// MethodFee can be passed to [Create] to provide a custom payin method fee.
func MethodFee(fee string) PaymentMethodOption {
	return func(pm *paymentMethodOptions) {
		pm.Fee = fee
	}
}

// MethodMin can be passed to [Create] to provide a custom min payin method amount.
func MethodMin(min string) PaymentMethodOption {
	return func(pm *paymentMethodOptions) {
		pm.Min = min
	}
}

// MethodMax can be passed to [Create] to provide a custom max payin method amount.
func MethodMax(max string) PaymentMethodOption {
	return func(pm *paymentMethodOptions) {
		pm.Max = max
	}
}

// MethodGroup can be passed to [Create] to provide a custom payin method group.
func MethodGroup(group string) PaymentMethodOption {
	return func(pm *paymentMethodOptions) {
		pm.Group = group
	}
}

// MethodName can be passed to [Create] to provide a custom payin method name.
func MethodName(name string) PaymentMethodOption {
	return func(pm *paymentMethodOptions) {
		pm.Name = name
	}
}

// MethodDescription can be passed to [Create] to provide a custom payin method description.
func MethodDescription(description string) PaymentMethodOption {
	return func(pm *paymentMethodOptions) {
		pm.Description = description
	}
}

// MethodRequiredPaymentDetails can be passed to [Create] to provide a custom payin method required payment details.
func RequiredDetails(details string) PaymentMethodOption {
	return func(pm *paymentMethodOptions) {
		pm.RequiredPaymentDetails = details
	}
}

func RequiredClaims(claims pexv2.PresentationDefinition) CreateOption {
	return func(o *createOptions) {
		o.requiredClaims = claims
	}
}
