package quote_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/TBD54566975/tbdex-go/tbdex/quote"
	"github.com/TBD54566975/tbdex-go/tbdex/rfq"
	"github.com/alecthomas/assert/v2"
	"github.com/tbd54566975/web5-go/dids/didjwk"
	"github.com/tbd54566975/web5-go/jws"

	"go.jetpack.io/typeid"
)

func TestCreate(t *testing.T) {
	pfiDID, err := didjwk.Create()
	assert.NoError(t, err)

	walletDID, err := didjwk.Create()
	assert.NoError(t, err)

	rfqID, _ := typeid.WithPrefix(rfq.Kind)

	quote := quote.Create(
		pfiDID.URI,
		walletDID.URI,
		rfqID.String(),
		time.Now().UTC().Format(time.RFC3339),
		quote.NewQuoteDetails("USD", "10", quote.DetailsFee("0.1")),
		quote.NewQuoteDetails("MXN", "500"),
	)

	assert.NoError(t, err)
	assert.NotZero(t, quote.Data.ExpiresAt)
	assert.NotZero(t, quote.Data.Payin)
	assert.NotZero(t, quote.Data.Payout)
	assert.Zero(t, quote.Signature)
}

func TestSign(t *testing.T) {
	pfiDID, _ := didjwk.Create()
	walletDID, _ := didjwk.Create()
	rfqID, _ := typeid.WithPrefix(rfq.Kind)

	quote := quote.Create(
		pfiDID.URI,
		walletDID.URI,
		rfqID.String(),
		time.Now().UTC().Format(time.RFC3339),
		quote.NewQuoteDetails("USD", "10"),
		quote.NewQuoteDetails("MXN", "500"),
	)

	err := quote.Sign(pfiDID)

	assert.NoError(t, err)
	assert.NotZero(t, quote.Signature)
}

func TestUnmarshalJSON(t *testing.T) {
	pfiDID, _ := didjwk.Create()
	walletDID, _ := didjwk.Create()
	rfqID, _ := typeid.WithPrefix(rfq.Kind)

	q := quote.Create(
		pfiDID.URI,
		walletDID.URI,
		rfqID.String(),
		time.Now().UTC().Format(time.RFC3339),
		quote.NewQuoteDetails("USD", "10"),
		quote.NewQuoteDetails(
			"MXN",
			"500",
			quote.DetailsFee("0.1"),
			quote.DetailsInstruction(quote.NewPaymentInstruction(quote.Instruction("use link"))),
		),
	)

	_ = q.Sign(pfiDID)

	bytes, err := json.Marshal(q)
	assert.NoError(t, err)

	quote := quote.Quote{}
	err = quote.UnmarshalJSON(bytes)
	assert.NoError(t, err)
}

func TestUnmarshal_Invalid(t *testing.T) {
	input := []byte(`{"doo": "doo"}`)

	quote := quote.Quote{}
	err := quote.UnmarshalJSON(input)
	assert.Error(t, err)
}

func TestVerify(t *testing.T) {
	pfiDID, _ := didjwk.Create()
	walletDID, _ := didjwk.Create()
	rfqID, _ := typeid.WithPrefix(rfq.Kind)

	quote := quote.Create(
		pfiDID.URI,
		walletDID.URI,
		rfqID.String(),
		time.Now().UTC().Format(time.RFC3339),
		quote.NewQuoteDetails("USD", "10"),
		quote.NewQuoteDetails("MXN", "500"),
	)

	_ = quote.Sign(pfiDID)

	err := quote.Verify()
	assert.NoError(t, err)
}

func TestVerify_FailsChangedPayload(t *testing.T) {
	pfiDID, _ := didjwk.Create()
	walletDID, _ := didjwk.Create()
	rfqID, _ := typeid.WithPrefix(rfq.Kind)

	quote := quote.Create(
		pfiDID.URI,
		walletDID.URI,
		rfqID.String(),
		time.Now().UTC().Format(time.RFC3339),
		quote.NewQuoteDetails("USD", "10"),
		quote.NewQuoteDetails("MXN", "500"),
	)

	_ = quote.Sign(pfiDID)
	quote.Data.ExpiresAt = "badtimestamp"

	err := quote.Verify()
	assert.Error(t, err)
}

func TestVerify_InvalidSignature(t *testing.T) {
	pfiDID, _ := didjwk.Create()
	walletDID, _ := didjwk.Create()
	rfqID, _ := typeid.WithPrefix(rfq.Kind)

	quote := quote.Create(
		pfiDID.URI,
		walletDID.URI,
		rfqID.String(),
		time.Now().UTC().Format(time.RFC3339),
		quote.NewQuoteDetails("USD", "10"),
		quote.NewQuoteDetails("MXN", "500"),
	)

	_ = quote.Sign(walletDID)

	quote.Signature = "Invalid"

	err := quote.Verify()
	assert.Error(t, err)
}

func TestVerify_SignedWithWrongDID(t *testing.T) {
	pfiDID, _ := didjwk.Create()
	walletDID, _ := didjwk.Create()
	wrongDID, _ := didjwk.Create()
	rfqID, _ := typeid.WithPrefix(rfq.Kind)

	quote := quote.Create(
		pfiDID.URI,
		walletDID.URI,
		rfqID.String(),
		time.Now().UTC().Format(time.RFC3339),
		quote.NewQuoteDetails("USD", "10"),
		quote.NewQuoteDetails("MXN", "500"),
	)

	_ = quote.Sign(walletDID)

	toSign, err := quote.Digest()
	assert.NoError(t, err)

	wrongSignature, err := jws.Sign(toSign, wrongDID, jws.DetachedPayload(true))
	assert.NoError(t, err)

	quote.Signature = wrongSignature

	err = quote.Verify()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "does not match message metadata from")
}
