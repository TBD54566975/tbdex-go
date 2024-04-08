package tbdex

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"go.jetpack.io/typeid"
)

type PayinMethodWithPrivate struct {
	Amount         string         `json:"amount"`
	Kind           string         `json:"kind"`
	PaymentDetails map[string]any `json:"paymentDetails"`
}

type PayoutMethodWithPrivate struct {
	Kind           string         `json:"kind"`
	PaymentDetails map[string]any `json:"paymentDetails"`
}

type rfqHashes struct {
	PayinHash  string
	PayoutHash string
	ClaimsHash string
}

// CreateRFQ creates an [RFQ]
//
// # An RFQ is a resource created by a customer of the PFI to request a quote
//
// [Offering]: https://github.com/TBD54566975/tbdex/tree/main/specs/protocol#rfq-request-for-quote
func CreateRFQ(from, to, offeringID string, payin PayinMethodWithPrivate, payout PayoutMethodWithPrivate, opts ...CreateRFQOption) (RFQ, error) {
	defaultID, err := typeid.WithPrefix(RFQKind)
	if err != nil {
		return RFQ{}, fmt.Errorf("failed to generate default id: %w", err)
	}

	r := createRFQOptions{
		id:         defaultID.String(),
		createdAt:  time.Now(),
		protocol:   "1.0",
		externalID: "",
	}

	for _, opt := range opts {
		opt(&r)
	}

	hashedData, privateData, err := hashData(&payin, &payout, r.claims)
    if err != nil {
		return RFQ{}, fmt.Errorf("failed to hash private data: %w", err)
	}

	return RFQ{
		MessageMetadata: MessageMetadata{
			From:       from,
			To:         to,
			Kind:       RFQKind,
			ID:         r.id,
			ExchangeID: r.id,
			CreatedAt:  r.createdAt.UTC().Format(time.RFC3339),
			ExternalID: r.externalID,
			Protocol:   r.protocol,
		},
		Data: RFQData{
			OfferingID: offeringID,
			Payin:      SelectedPayinMethod{
                Amount: payin.Amount,
                Kind: payin.Kind,
                PaymentDetailsHash: hashedData.PayinHash,
            },
			Payout:     SelectedPayoutMethod{
                Kind: payin.Kind,
                PaymentDetailsHash: hashedData.PayoutHash,
            },
			ClaimsHash: hashedData.ClaimsHash,
		},
		PrivateData: privateData,
	}, nil
}

func hashData(payin *PayinMethodWithPrivate, payout *PayoutMethodWithPrivate, claims []string) (rfqHashes, RFQPrivateData, error) {
    randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return rfqHashes{}, RFQPrivateData{}, err
	}

	salt := base64.URLEncoding.EncodeToString(randomBytes)

    payinHash, err := digestData(salt, payin)
    if err != nil {
		return rfqHashes{}, RFQPrivateData{}, err
	}
    payoutHash, err := digestData(salt, payout)
    if err != nil {
		return rfqHashes{}, RFQPrivateData{}, err
	}
    claimsHash, err := digestData(salt, claims)
    if err != nil {
		return rfqHashes{}, RFQPrivateData{}, err
	}


    h := rfqHashes{
        PayinHash: payinHash,
        PayoutHash: payoutHash,
        ClaimsHash: claimsHash,
    }

    p := RFQPrivateData {
        Salt: salt,
        Payin: &PrivatePaymentDetails{
            PaymentDetails: payin.PaymentDetails,
        },
        Payout: &PrivatePaymentDetails{
            PaymentDetails: payout.PaymentDetails,
        },
        Claims: claims,
    }

    return h, p, nil
}

func digestData(salt string, data any) (string, error) {
    digestible := []interface{}{salt, data}

	byteArray, err := Digest(digestible)
	if err != nil {
		return "", err
	}

	encodedString := base64.URLEncoding.EncodeToString(byteArray)

	return encodedString, nil
}

// CreateRFQOption is a function type used to apply options to RFQ creation.
type createRFQOptions struct {
	id         string
	createdAt  time.Time
	protocol   string
	externalID string
	claims     []string
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

// WithRFQExternalID can be passed to [CreateRFQ] to provide a custom external id.
func WithRFQExternalID(externalID string) CreateRFQOption {
	return func(r *createRFQOptions) {
		r.externalID = externalID
	}
}

// WithRFQClaims can be passed to [CreateRFQ] to provide a custom external id.
func WithRFQClaims(claims []string) CreateRFQOption {
	return func(r *createRFQOptions) {
		r.claims = claims
	}
}

// PayinMethodWithPrivateOption is a function type used to apply options to [PayinMethodWithPrivate] creation.
type PayinMethodWithPrivateOption func(*PayinMethodWithPrivate)

// PayoutMethodWithPrivateOption is a function type used to apply options to [PayoutMethodWithPrivate] creation.
type PayoutMethodWithPrivateOption func(*PayoutMethodWithPrivate)

// WithPayinMethodWithPrivate can be passed to [WithRFQSelectedPayinMethod] to provide a hash of the payment details.
func WithPayinMethodWithPrivate(detailsPrivate map[string]any) PayinMethodWithPrivateOption {
	return func(pm *PayinMethodWithPrivate) {
		pm.PaymentDetails = detailsPrivate
	}
}

// WithPayoutMethodWithPrivate can be passed to [WithRFQSelectedPayoutMethod] to provide a hash of the payment details.
func WithPayoutMethodWithPrivate(detailsPrivate map[string]any) PayoutMethodWithPrivateOption {
	return func(pm *PayoutMethodWithPrivate) {
		pm.PaymentDetails = detailsPrivate
	}
}

// WithRFQSelectedPayinMethod can be passed to [Create] to provide a payin method.
func WithRFQSelectedPayinMethod(amount, kind string, opts ...PayinMethodWithPrivateOption) PayinMethodWithPrivate {
	s := PayinMethodWithPrivate{
		Amount: amount,
		Kind:   kind,
	}

	for _, opt := range opts {
		opt(&s)
	}

	return s
}

// WithRFQSelectedPayoutMethod can be passed to [Create] to provide a payout method.
func WithRFQSelectedPayoutMethod(kind string, opts ...PayoutMethodWithPrivateOption) PayoutMethodWithPrivate {
	s := PayoutMethodWithPrivate{
		Kind: kind,
	}

	for _, opt := range opts {
		opt(&s)
	}

	return s
}
