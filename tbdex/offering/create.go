package offering

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/TBD54566975/tbdex-go/tbdex"
	"github.com/tbd54566975/web5-go/dids/did"
	"github.com/tbd54566975/web5-go/pexv2"
	"go.jetpack.io/typeid"
)

// Create creates an [Offering]
//
// An Offering is a resource created by a PFI to define requirements for a given currency pair offered for exchange.
//
// [Offering]: https://github.com/TBD54566975/tbdex/tree/main/specs/protocol#offering
func Create(fromDID did.BearerDID, payin *PayinDetails, payout *PayoutDetails, rate string, opts ...CreateOption) (Offering, error) {
	o := createOptions{
		id:          typeid.Must(typeid.New[ID]()),
		createdAt:   time.Now(),
		updatedAt:   time.Now(),
		description: fmt.Sprintf("%s for %s", payout.CurrencyCode, payin.CurrencyCode),
		protocol:    "1.0",
	}

	for _, opt := range opts {
		opt(&o)
	}

	if len(payin.Methods) == 0 {
		return Offering{}, errors.New("at least 1 payin method is required")
	}

	if len(payout.Methods) == 0 {
		return Offering{}, errors.New("at least 1 payout method is required")
	}

	offering := Offering{
		ResourceMetadata: &tbdex.ResourceMetadata{
			From:      fromDID.URI,
			Kind:      Kind,
			ID:        o.id.String(),
			CreatedAt: o.createdAt.UTC().Format(time.RFC3339),
			UpdatedAt: o.updatedAt.UTC().Format(time.RFC3339),
			Protocol:  o.protocol,
		},
		Data: &Data{
			Payin:          payin,
			Payout:         payout,
			Rate:           rate,
			Description:    o.description,
			RequiredClaims: o.requiredClaims,
		},
	}

	signature, err := tbdex.Sign(offering, fromDID)
	if err != nil {
		return Offering{}, fmt.Errorf("failed to sign offering: %w", err)
	}

	offering.Signature = signature

	return offering, nil
}

// NewPayin creates PayinDetails
func NewPayin(currencyCode string, methods []PayinMethod, opts ...PaymentOption) *PayinDetails {
	return &PayinDetails{
		CurrencyCode: currencyCode,
		Methods:      methods,
	}
}

// NewPayinMethod creates PayinMethod
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

// NewPayout creates PayoutDetails
func NewPayout(currencyCode string, methods []PayoutMethod, opts ...PaymentOption) *PayoutDetails {
	o := paymentOptions{}
	for _, opt := range opts {
		opt(&o)
	}

	return &PayoutDetails{
		CurrencyCode: currencyCode,
		Min:          o.Min,
		Max:          o.Max,
		Methods:      methods,
	}
}

// NewPayoutMethod creates PayoutMethod
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
	id             ID
	createdAt      time.Time
	updatedAt      time.Time
	description    string
	protocol       string
	requiredClaims *pexv2.PresentationDefinition
}

// CreateOption implements functional options pattern for [Create].
type CreateOption func(*createOptions)

// OfferingID can be passed to [CreateOffering] to provide a custom id.
func OfferingID(id ID) CreateOption {
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

// UpdatedAt can be passed to [Create] to provide a custom updated at time.
func UpdatedAt(t time.Time) CreateOption {
	return func(o *createOptions) {
		o.updatedAt = t
	}
}

// Description can be passed to [Create] to provide a custom description.
func Description(d string) CreateOption {
	return func(o *createOptions) {
		o.description = d
	}
}

type paymentOptions struct {
	Min string
	Max string
}

// PaymentOption implements functional options pattern for Payin and Payout
type PaymentOption func(*paymentOptions)

// Min can be passed to [Create] to provide a custom min payin amount.
func Min(min string) PaymentOption {
	return func(p *paymentOptions) {
		p.Min = min
	}
}

// Max can be passed to [Create] to provide a custom max payin amount.
func Max(max string) PaymentOption {
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
	RequiredPaymentDetails json.RawMessage
}

// PaymentMethodOption implements functional options pattern for PayinMethod and PayoutMethod.
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

// RequiredDetails can be passed to [Create] to provide a custom payin method required payment details.
func RequiredDetails(details string) PaymentMethodOption {
	return func(pm *paymentMethodOptions) {
		pm.RequiredPaymentDetails = json.RawMessage(details)
	}
}

// RequiredClaims can be used to set claims on the offering being created
func RequiredClaims(claims pexv2.PresentationDefinition) CreateOption {
	return func(o *createOptions) {
		o.requiredClaims = &claims
	}
}
