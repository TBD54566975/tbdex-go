package rfq

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/TBD54566975/tbdex-go/tbdex"
	"github.com/tbd54566975/web5-go/crypto"
	"go.jetpack.io/typeid"
)

// CreateRFQ creates an [RFQ]
//
// # An RFQ is a resource created by a customer of the PFI to request a quote
//
// [RFQ]: https://github.com/TBD54566975/tbdex/tree/main/specs/protocol#rfq-request-for-quote
func CreateRFQ(from, to, offeringID string, payin PayinMethodWithDetails, payout PayoutMethodWithDetails, opts ...CreateRFQOption) (RFQ, error) {
	r := createRFQOptions{
		id:         typeid.Must(typeid.WithPrefix(RFQKind)).String(),
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
		MessageMetadata: tbdex.MessageMetadata{
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
			Payin: SelectedPayinMethod{
				Amount:             payin.Amount,
				Kind:               payin.Kind,
				PaymentDetailsHash: hashedData.PayinHash,
			},
			Payout: SelectedPayoutMethod{
				Kind:               payin.Kind,
				PaymentDetailsHash: hashedData.PayoutHash,
			},
			ClaimsHash: hashedData.ClaimsHash,
		},
		PrivateData: &privateData,
	}, nil
}

func hashData(payin *PayinMethodWithDetails, payout *PayoutMethodWithDetails, claims []string) (rfqHashes, RFQPrivateData, error) {
	randomBytes, err := crypto.GenerateEntropy(crypto.Entropy128)
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
		PayinHash:  payinHash,
		PayoutHash: payoutHash,
		ClaimsHash: claimsHash,
	}

	p := RFQPrivateData{
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
	digestible := []any{salt, data}

	byteArray, err := tbdex.DigestJSON(digestible)
	if err != nil {
		return "", err
	}

	encodedString := base64.URLEncoding.EncodeToString(byteArray)

	return encodedString, nil
}

// PayinMethodWithDetails is used to create the payin method for an RFQ
type PayinMethodWithDetails struct {
	Amount         string         `json:"amount"`
	Kind           string         `json:"kind"`
	PaymentDetails map[string]any `json:"paymentDetails"`
}

// PayoutMethodWithDetails is used to create the payout method for an RFQ
type PayoutMethodWithDetails struct {
	Kind           string         `json:"kind"`
	PaymentDetails map[string]any `json:"paymentDetails"`
}

type rfqHashes struct {
	PayinHash  string
	PayoutHash string
	ClaimsHash string
}

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

// PayinMethodWithDetailsOption is a function type used to apply options to [PayinMethodWithDetails] creation.
type PayinMethodWithDetailsOption func(*PayinMethodWithDetails)

// PayoutMethodWithDetailsOption is a function type used to apply options to [PayoutMethodWithDetails] creation.
type PayoutMethodWithDetailsOption func(*PayoutMethodWithDetails)

// WithPayinMethodWithDetails can be passed to [WithRFQSelectedPayinMethod] to provide arbitrary payment details.
func WithPayinMethodWithDetails(details map[string]any) PayinMethodWithDetailsOption {
	return func(pm *PayinMethodWithDetails) {
		pm.PaymentDetails = details
	}
}

// WithPayoutMethodWithDetails can be passed to [WithRFQSelectedPayoutMethod] to provide arbitrary payment details.
func WithPayoutMethodWithDetails(details map[string]any) PayoutMethodWithDetailsOption {
	return func(pm *PayoutMethodWithDetails) {
		pm.PaymentDetails = details
	}
}

// WithRFQSelectedPayinMethod can be passed to [Create] to provide a payin method.
func WithRFQSelectedPayinMethod(amount, kind string, opts ...PayinMethodWithDetailsOption) PayinMethodWithDetails {
	s := PayinMethodWithDetails{
		Amount: amount,
		Kind:   kind,
	}

	for _, opt := range opts {
		opt(&s)
	}

	return s
}

// WithRFQSelectedPayoutMethod can be passed to [Create] to provide a payout method.
func WithRFQSelectedPayoutMethod(kind string, opts ...PayoutMethodWithDetailsOption) PayoutMethodWithDetails {
	s := PayoutMethodWithDetails{
		Kind: kind,
	}

	for _, opt := range opts {
		opt(&s)
	}

	return s
}
