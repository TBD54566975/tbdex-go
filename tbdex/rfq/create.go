package rfq

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/TBD54566975/tbdex-go/tbdex"
	"github.com/tbd54566975/web5-go/crypto"
	"go.jetpack.io/typeid"
)

// Create creates an [RFQ]
//
// # An RFQ is a resource created by a customer of the PFI to request a quote
//
// [RFQ]: https://github.com/TBD54566975/tbdex/tree/main/specs/protocol#rfq-request-for-quote
func Create(from, to, offeringID string, payin PayinMethodWithDetails, payout PayoutMethodWithDetails, opts ...CreateOption) (RFQ, error) {
	r := createOptions{
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

type createOptions struct {
	id         string
	createdAt  time.Time
	protocol   string
	externalID string
	claims     []string
}

// CreateOption is a function type used to apply options to RFQ creation.
type CreateOption func(*createOptions)

// ID can be passed to [Create] to provide a custom id.
func ID(id string) CreateOption {
	return func(r *createOptions) {
		r.id = id
	}
}

// CreatedAt can be passed to [Create] to provide a custom created at time.
func CreatedAt(t time.Time) CreateOption {
	return func(r *createOptions) {
		r.createdAt = t
	}
}

// ExternalID can be passed to [Create] to provide a custom external id.
func ExternalID(externalID string) CreateOption {
	return func(r *createOptions) {
		r.externalID = externalID
	}
}

// Claims can be passed to [Create] to provide a custom external id.
func Claims(claims []string) CreateOption {
	return func(r *createOptions) {
		r.claims = claims
	}
}

type paymentMethodOptions struct {
	details map[string]any
}

type PaymentMethodOption func(*paymentMethodOptions)

// PaymentDetails can be passed to [Payin] to provide arbitrary payment details.
func PaymentDetails(details map[string]any) PaymentMethodOption {
	return func(pm *paymentMethodOptions) {
		pm.details = details
	}
}

// Payin can be passed to [Create] to provide a payin method.
func Payin(amount, kind string, opts ...PaymentMethodOption) PayinMethodWithDetails {
	s := PayinMethodWithDetails{
		Amount: amount,
		Kind:   kind,
	}

	o := paymentMethodOptions{}
	for _, opt := range opts {
		opt(&o)
	}

	s.PaymentDetails = o.details

	return s
}

// Payout can be passed to [Create] to provide a payout method.
func Payout(kind string, opts ...PaymentMethodOption) PayoutMethodWithDetails {
	s := PayoutMethodWithDetails{
		Kind: kind,
	}

	o := paymentMethodOptions{}
	for _, opt := range opts {
		opt(&o)
	}

	s.PaymentDetails = o.details

	return s
}
