package tbdex

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gowebpki/jcs"
	"github.com/tbd54566975/web5-go/dids/did"
	"go.jetpack.io/typeid"
)

// RFQKind identifies this message kind
const RFQKind = "rfq"

type RFQ struct {
	MessageMetadata MessageMetadata `json:"metadata"`
	RFQData         RFQData `json:"data"`
	Signature       string `json:"signature"`
}

type RFQData struct {
	OfferingID string               `json:"offeringId"`
	Payin      SelectedPayinMethod  `json:"payin"`
	Payout     SelectedPayoutMethod `json:"payout"`
	ClaimsHash string               `json:"claimsHash,omitempty"`
}

type SelectedPayinMethod struct {
	Amount             string `json:"amount"`
	Kind               string `json:"kind"`
	PaymentDetailsHash string `json:"paymentDetailsHash,omitempty"`
}

type SelectedPayoutMethod struct {
	Kind               string `json:"kind"`
	PaymentDetailsHash string `json:"paymentDetailsHash,omitempty"`
}

type createRFQOptions struct {
	id         string
	createdAt  time.Time
	claimsHash string
	protocol   string
	externalID string
}

type CreateRFQOption func(*createRFQOptions)

func (rfq RFQ) Digest() ([]byte, error) {
	payload := map[string]any{"metadata": rfq.MessageMetadata, "data": rfq.RFQData}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal rfq: %w", err)
	}

	canonicalized, err := jcs.Transform(payloadBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to canonicalize rfq: %w", err)
	}

	hasher := sha256.New()
	_, err = hasher.Write(canonicalized)
	if err != nil {
		return nil, fmt.Errorf("failed to compute digest: %w", err)
	}

	return hasher.Sum(nil), nil
}

// Sign cryptographically signs the RFQ using DID's private key
func (r *RFQ) Sign(bearerDID did.BearerDID) error {
	r.MessageMetadata.From = bearerDID.URI

	signature, err := Sign(r, bearerDID)
	if err != nil {
		return fmt.Errorf("failed to sign rfq: %w", err)
	}

	r.Signature = signature

	return nil
}

func (r *RFQ) UnmarshalJSON(data []byte) error {
	err := Validate(TypeMessage, data, WithKind(RFQKind))
	if err != nil {
		return fmt.Errorf("invalid rfq: %w", err)
	}

	ret := rfq{}
	err = json.Unmarshal(data, &ret)
	if err != nil {
		return fmt.Errorf("failed to unmarshal rfq: %w", err)
	}

	*r = RFQ(ret)

	return nil
}

type rfq RFQ

func CreateRFQ(from, to, offeringID string, payin SelectedPayinMethod, payout SelectedPayoutMethod, opts ...CreateRFQOption) (RFQ, error) {
	defaultID, err := typeid.WithPrefix(RFQKind)
	if err != nil {
		return RFQ{}, fmt.Errorf("failed to generate default id: %w", err)
	}

	r := createRFQOptions{
		id:          defaultID.String(),
		createdAt:   time.Now(),
		protocol:    "1.0",
		externalID: "",
		claimsHash: "",
	}

	for _, opt := range opts {
		opt(&r)
	}

	return RFQ{
		MessageMetadata: MessageMetadata{
			From: from,
			To: to,
			Kind:      RFQKind,
			ID:        r.id,
			ExchangeID: r.id,
			CreatedAt: r.createdAt.UTC().Format(time.RFC3339),
			ExternalID: r.externalID,
			Protocol:  r.protocol,
		},
		RFQData: RFQData{
				OfferingID: offeringID,
				Payin: payin,
				Payout: payout,
				ClaimsHash: r.claimsHash,
		},
	}, nil
}

func WithRFQID(id string) CreateRFQOption {
	return func(r *createRFQOptions) {
		r.id = id
	}
}

func WithRFQCreatedAt(t time.Time) CreateRFQOption {
	return func(r *createRFQOptions) {
		r.createdAt = t
	}
}

func WithRFQClaimsHash(claimsHash string) CreateRFQOption {
	return func(r *createRFQOptions) {
		r.claimsHash = claimsHash
	}
}

func WithRFQSelectedPayinMethod(amount, kind string, opts ...SelectedPayinMethodOption) SelectedPayinMethod {
	s := SelectedPayinMethod{
		Amount: amount,
		Kind: kind,
	}

	for _, opt := range opts {
		opt(&s)
	}

	return s
}

func WithRFQSelectedPayoutMethod(kind string, opts ...SelectedPayoutMethodOption) SelectedPayoutMethod {
	s := SelectedPayoutMethod{
		Kind: kind,
	}

	for _, opt := range opts {
		opt(&s)
	}

	return s
}

type SelectedPayinMethodOption func(*SelectedPayinMethod)
type SelectedPayoutMethodOption func(*SelectedPayoutMethod)

func SelectedPayinMethodDetailsHash(detailsHash string) SelectedPayinMethodOption {
	return func(pm *SelectedPayinMethod) {
		pm.PaymentDetailsHash = detailsHash
	}
}

func SelectedPayoutMethodDetailsHash(detailsHash string) SelectedPayoutMethodOption {
	return func(pm *SelectedPayoutMethod) {
		pm.PaymentDetailsHash = detailsHash
	}
}
