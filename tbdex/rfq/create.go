package rfq

import (
	"encoding/base64"
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

	privateData := RFQPrivateData{}
	var hashedPayin string
	var hashedPayout string
	var hashedClaims string
	if (len(payin.PaymentDetails) != 0) || (len(payout.PaymentDetails) != 0) || (len(r.claims) != 0) {
		randomBytes, err := crypto.GenerateEntropy(crypto.Entropy128)
		if err != nil {
			return RFQ{}, err
		}

		salt := base64.RawURLEncoding.EncodeToString(randomBytes)
		privateData.Salt = salt

		if len(payin.PaymentDetails) != 0 {
			hashedPayin, err = hashData(salt, payin)
			if err != nil {
				return RFQ{}, err
			}
			privateData.Payin = &PrivatePaymentDetails{
				PaymentDetails: payin.PaymentDetails,
			}
		}

		if len(payout.PaymentDetails) != 0 {
			hashedPayout, err = hashData(salt, payout)
			if err != nil {
				return RFQ{}, err
			}

			privateData.Payout = &PrivatePaymentDetails{
				PaymentDetails: payout.PaymentDetails,
			}
		}

		if len(r.claims) != 0 {
			hashedClaims, err = hashData(salt, r.claims)
			if err != nil {
				return RFQ{}, err
			}

			privateData.Claims = r.claims
		}

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
				PaymentDetailsHash: hashedPayin,
			},
			Payout: SelectedPayoutMethod{
				Kind:               payin.Kind,
				PaymentDetailsHash: hashedPayout,
			},
			ClaimsHash: hashedClaims,
		},
		PrivateData: &privateData,
	}, nil
}

func hashData(salt string, data any) (string, error) {
	digest, err := digestData(salt, data)

	if err != nil {
		return "", err
	}

	return digest, nil
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

// Claims can be passed to [Create] to provide claims if required by the offering.
func Claims(claims []string) CreateOption {
	return func(r *createOptions) {
		r.claims = claims
	}
}

type paymentMethodOptions struct {
	details map[string]any
}

// PaymentMethodOption is a function type used to apply options to payment methods.
type PaymentMethodOption func(*paymentMethodOptions)

// PaymentDetails can be passed to [Payin] or [Payout] to provide arbitrary payment details.
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
