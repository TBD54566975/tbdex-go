package rfq

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"github.com/gowebpki/jcs"
	"github.com/tbd54566975/web5-go/dids/did"
)

// RFQKind identifies this message kind
const RFQKind = "rfq"

// RFQ represents a request for quote message within the exchange.
type RFQ struct {
	MessageMetadata MessageMetadata `json:"metadata"`
	Data            RFQData         `json:"data"`
	PrivateData     RFQPrivateData  `json:"privateData"`
	Signature       string          `json:"signature"`
}

// RFQData encapsulates the data content of a request for quote.
type RFQData struct {
	OfferingID string               `json:"offeringId"`
	Payin      SelectedPayinMethod  `json:"payin"`
	Payout     SelectedPayoutMethod `json:"payout"`
	ClaimsHash string               `json:"claimsHash,omitempty"`
}

// RFQPrivateData contains data which can be detached from the payload without disrupting integrity.
type RFQPrivateData struct {
	Salt   string                 `json:"salt"`
	Claims []string               `json:"claims"`
	Payin  *PrivatePaymentDetails `json:"payin"`
	Payout *PrivatePaymentDetails `json:"payout"`
}

// PrivatePaymentDetails is a container for the cleartest [PaymentDetails]
type PrivatePaymentDetails struct {
	PaymentDetails map[string]any `json:"paymentDetails"`
}

// SelectedPayinMethod represents the chosen method for the pay-in
type SelectedPayinMethod struct {
	Amount             string `json:"amount"`
	Kind               string `json:"kind"`
	PaymentDetailsHash string `json:"paymentDetailsHash,omitempty"`
}

// SelectedPayoutMethod represents the chosen method for the pay-out
type SelectedPayoutMethod struct {
	Kind               string `json:"kind"`
	PaymentDetailsHash string `json:"paymentDetailsHash,omitempty"`
}

// Digest computes a hash of the resource
// A digest is the output of the hash function. It's a fixed-size string of bytes
//   - that uniquely represents the data input into the hash function. The digest is often used for
//   - data integrity checks, as any alteration in the input data results in a significantly
//   - different digest.
//     *
//   - It takes the algorithm identifier of the hash function and data to digest as input and returns
//   - the digest of the data.
func (r RFQ) Digest() ([]byte, error) {
	payload := map[string]any{"metadata": r.MessageMetadata, "data": r.Data}
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

// UnmarshalJSON validates and unmarshals the input data into an RFQ.
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

	if err := VerifySignature(r, r.Signature); err != nil { return err }

	// TODO verify private data
	// if requirePrivateData {
	// 	err = r.verifyAllPrivateData()
	// 	if err != nil {
	// 	return fmt.Errorf("failed to verify private data: %w", err)

	// 	}
	// }

	return nil
}

func (r *RFQ) verifyAllPrivateData() error {

	return nil
}

type rfq RFQ
