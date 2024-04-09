package rfq

import (
	"encoding/json"
	"fmt"

	"github.com/TBD54566975/tbdex-go/tbdex"
	"github.com/tbd54566975/web5-go/dids/did"
)

// RFQKind identifies this message kind
const RFQKind = "rfq"

// RFQ represents a request for quote message within the exchange.
type RFQ struct {
	MessageMetadata tbdex.MessageMetadata `json:"metadata"`
	Data            RFQData               `json:"data"`
	PrivateData     *RFQPrivateData       `json:"privateData,omitempty"`
	Signature       string                `json:"signature"`
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

// Digest computes a hash of the rfq
func (r RFQ) Digest() ([]byte, error) {
	payload := map[string]any{"metadata": r.MessageMetadata, "data": r.Data}

	hashed, err := tbdex.DigestJSON(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to digest rfq: %w", err)
	}

	return hashed, nil
}

// Sign cryptographically signs the RFQ using DID's private key
func (r *RFQ) Sign(bearerDID did.BearerDID) error {
	r.MessageMetadata.From = bearerDID.URI

	signature, err := tbdex.Sign(r, bearerDID)
	if err != nil {
		return fmt.Errorf("failed to sign rfq: %w", err)
	}

	r.Signature = signature

	return nil
}

// UnmarshalJSON validates and unmarshals the input data into an RFQ.
func (r *RFQ) UnmarshalJSON(data []byte) error {
	err := tbdex.Validate(tbdex.TypeMessage, data, tbdex.WithKind(RFQKind))
	if err != nil {
		return fmt.Errorf("invalid rfq: %w", err)
	}

	ret := rfq{}
	err = json.Unmarshal(data, &ret)
	if err != nil {
		return fmt.Errorf("failed to JSON unmarshal rfq: %w", err)
	}

	*r = RFQ(ret)

	_, err = tbdex.VerifySignature(r, r.Signature)
	if err != nil {
		return fmt.Errorf("failed to verify RFQ signature: %w", err)
	}

	// TODO add check when decoded.SignerDID is implemented
	// if decoded.SignerDID != r.MessageMetadata.From {
	// 	return errors.New("signer: %w does not match message metadata from: %w", decoded.Header.SignerDID, r.MessageMetadata.From)
	// }

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
