package offering

import (
	"errors"
	"fmt"
	"time"

	"github.com/TBD54566975/tbdex-go/protocol/resource"
	"go.jetpack.io/typeid"
)

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
	}

	for _, opt := range opts {
		opt(&o)
	}

	if len(payin.Methods) == 0 {
		return Offering{}, errors.New("1 payin method is required.")
	}

	if len(payout.Methods) == 0 {
		return Offering{}, errors.New("1 payout method is required.")
	}

	return Offering{
		Metadata: resource.Metadata{
			Kind:      Kind,
			ID:        o.id,
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
}

type CreateOption func(*createOptions)

func WithID(id string) CreateOption {
	return func(o *createOptions) {
		o.id = id
	}
}

func WithCreatedAt(t time.Time) CreateOption {
	return func(o *createOptions) {
		o.createdAt = t
	}
}

func WithUpdatedAt(t time.Time) CreateOption {
	return func(o *createOptions) {
		o.updatedAt = t
	}
}

func WithDescription(d string) CreateOption {
	return func(o *createOptions) {
		o.description = d
	}
}

type PayinOption func(*PayinDetails)

func WithPayinMin(min string) PayinOption {
	return func(p *PayinDetails) {
		p.Min = min
	}
}

func WithPayinMax(max string) PayinOption {
	return func(p *PayinDetails) {
		p.Max = max
	}
}

type PayinMethodOption func(*PayinMethod)

func WithPayinMethodFee(fee string) PayinMethodOption {
	return func(pm *PayinMethod) {
		pm.Fee = fee
	}
}

func WithPayinMethodMin(min string) PayinMethodOption {
	return func(pm *PayinMethod) {
		pm.Min = min
	}
}

func WithPayinMethodMax(max string) PayinMethodOption {
	return func(pm *PayinMethod) {
		pm.Max = max
	}
}

func WithPayinMethodGroup(group string) PayinMethodOption {
	return func(pm *PayinMethod) {
		pm.Group = group
	}
}

func WithPayinMethodName(name string) PayinMethodOption {
	return func(pm *PayinMethod) {
		pm.Name = name
	}
}

func WithPayinMethodDescription(description string) PayinMethodOption {
	return func(pm *PayinMethod) {
		pm.Description = description
	}
}

func WithPayinMethodRequiredPaymentDetails(details string) PayinMethodOption {
	return func(pm *PayinMethod) {
		pm.RequiredPaymentDetails = details
	}
}

func WithPayinMethod(kind string, opts ...PayinMethodOption) PayinMethod {
	pm := PayinMethod{Kind: kind}

	for _, opt := range opts {
		opt(&pm)
	}

	return pm
}

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

type PayoutOption func(*PayoutDetails)

func WithPayoutMin(min string) PayoutOption {
	return func(p *PayoutDetails) {
		p.Min = min
	}
}

func WithPayoutMax(max string) PayoutOption {
	return func(p *PayoutDetails) {
		p.Max = max
	}
}

type PayoutMethodOption func(*PayoutMethod)

func WithPayoutMethodFee(fee string) PayoutMethodOption {
	return func(pm *PayoutMethod) {
		pm.Fee = fee
	}
}

func WithPayoutMethodMin(min string) PayoutMethodOption {
	return func(pm *PayoutMethod) {
		pm.Min = min
	}
}

func WithPayoutMethodMax(max string) PayoutMethodOption {
	return func(pm *PayoutMethod) {
		pm.Max = max
	}
}

func WithPayoutMethodGroup(group string) PayoutMethodOption {
	return func(pm *PayoutMethod) {
		pm.Group = group
	}
}

func WithPayoutMethodName(name string) PayoutMethodOption {
	return func(pm *PayoutMethod) {
		pm.Name = name
	}
}

func WithPayoutMethodDescription(description string) PayoutMethodOption {
	return func(pm *PayoutMethod) {
		pm.Description = description
	}
}

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
