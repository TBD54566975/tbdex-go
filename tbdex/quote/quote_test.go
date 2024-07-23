package quote_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/alecthomas/assert/v2"
	"github.com/shopspring/decimal"
	"github.com/tbd54566975/web5-go/dids/didjwk"
	"github.com/tbd54566975/web5-go/jws"

	"go.jetpack.io/typeid"

	"github.com/TBD54566975/tbdex-go/tbdex/closemsg"
	"github.com/TBD54566975/tbdex-go/tbdex/order"
	"github.com/TBD54566975/tbdex-go/tbdex/orderstatus"
	"github.com/TBD54566975/tbdex-go/tbdex/quote"
	"github.com/TBD54566975/tbdex-go/tbdex/rfq"
)

func TestCreate(t *testing.T) {
	pfiDID, err := didjwk.Create()
	assert.NoError(t, err)

	walletDID, err := didjwk.Create()
	assert.NoError(t, err)

	rfqID, _ := typeid.WithPrefix(rfq.Kind)

	quote, err := quote.Create(
		pfiDID,
		walletDID.URI,
		rfqID.String(),
		time.Now().UTC().Format(time.RFC3339),
		"16.665",
		quote.NewQuoteDetails(
			"USD",
			decimal.RequireFromString("10"),
			quote.DetailsFee(decimal.RequireFromString("0.1")),
		),
		quote.NewQuoteDetails("MXN", decimal.RequireFromString("500")),
	)
	assert.NoError(t, err)

	assert.NoError(t, err)
	assert.NotZero(t, quote.Data.Rate)
	assert.NotZero(t, quote.Data.ExpiresAt)
	assert.NotZero(t, quote.Data.Payin)
	assert.NotZero(t, quote.Data.Payout)
	assert.NotZero(t, quote.Signature)
	assert.NotZero(t, quote.Data.Payin.Total)
}

func TestUnmarshal(t *testing.T) {
	pfiDID, _ := didjwk.Create()
	walletDID, _ := didjwk.Create()
	rfqID, _ := typeid.WithPrefix(rfq.Kind)

	q, err := quote.Create(
		pfiDID,
		walletDID.URI,
		rfqID.String(),
		time.Now().UTC().Format(time.RFC3339),
		"16.665",
		quote.NewQuoteDetails(
			"USD",
			decimal.RequireFromString("10"),
			quote.DetailsFee(decimal.RequireFromString("0.1")),
		),
		quote.NewQuoteDetails(
			"MXN",
			decimal.RequireFromString("500"),
		))
	assert.NoError(t, err)

	bytes, err := json.Marshal(q)
	assert.NoError(t, err)

	quote := quote.Quote{}
	err = json.Unmarshal(bytes, &quote)
	assert.NoError(t, err)
}

func TestUnmarshal_Empty(t *testing.T) {
	input := []byte(`{"metadata":{},"data":{},"signature":"eyJhbGciOiJFZERTQSIsImtpZCI6ImRpZDpkaHQ6M3doZnRncGJkamloeDl6ZTl0ZG41NzV6cXptNHF3Y2NldG5mMXliaWlidXphZDdycm15eSMwIn0..ZvoVDuSrqqdXsSXgqB-U26tAU1WqUqqU_KpD1KvdYocIcmTsshjUASEwM_lUz1UnGglqkWeCIrHqrm9NNGDqBw"}`)

	quote := quote.Quote{}
	_ = json.Unmarshal(input, &quote)

	assert.Zero(t, quote.Metadata)
	assert.Zero(t, quote.Data)
}

func TestUnmarshal_Invalid(t *testing.T) {
	input := []byte(`{"doo": "doo"}`)

	quote := quote.Quote{}
	err := json.Unmarshal(input, &quote)
	assert.Error(t, err)
}

func TestVerify(t *testing.T) {
	pfiDID, _ := didjwk.Create()
	walletDID, _ := didjwk.Create()
	rfqID, _ := typeid.WithPrefix(rfq.Kind)

	quote, err := quote.Create(
		pfiDID,
		walletDID.URI,
		rfqID.String(),
		time.Now().UTC().Format(time.RFC3339),
		"16.665",
		quote.NewQuoteDetails("USD", decimal.RequireFromString("10")),
		quote.NewQuoteDetails("MXN", decimal.RequireFromString("500")),
	)
	assert.NoError(t, err)

	err = quote.Verify()
	assert.NoError(t, err)
}

func TestVerify_FailsChangedPayload(t *testing.T) {
	pfiDID, _ := didjwk.Create()
	walletDID, _ := didjwk.Create()
	rfqID, _ := typeid.WithPrefix(rfq.Kind)

	quote, err := quote.Create(
		pfiDID,
		walletDID.URI,
		rfqID.String(),
		time.Now().UTC().Format(time.RFC3339),
		"16.665",
		quote.NewQuoteDetails("USD", decimal.RequireFromString("10")),
		quote.NewQuoteDetails("MXN", decimal.RequireFromString("500")),
	)
	assert.NoError(t, err)

	quote.Data.ExpiresAt = "badtimestamp"

	err = quote.Verify()
	assert.Error(t, err)
}

func TestVerify_InvalidSignature(t *testing.T) {
	pfiDID, _ := didjwk.Create()
	walletDID, _ := didjwk.Create()
	rfqID, _ := typeid.WithPrefix(rfq.Kind)

	quote, err := quote.Create(
		pfiDID,
		walletDID.URI,
		rfqID.String(),
		time.Now().UTC().Format(time.RFC3339),
		"16.665",
		quote.NewQuoteDetails("USD", decimal.RequireFromString("10")),
		quote.NewQuoteDetails("MXN", decimal.RequireFromString("500")),
	)
	assert.NoError(t, err)

	quote.Signature = "Invalid"

	err = quote.Verify()
	assert.Error(t, err)
}

func TestVerify_SignedWithWrongDID(t *testing.T) {
	pfiDID, _ := didjwk.Create()
	walletDID, _ := didjwk.Create()
	wrongDID, _ := didjwk.Create()
	rfqID, _ := typeid.WithPrefix(rfq.Kind)

	quote, err := quote.Create(
		pfiDID,
		walletDID.URI,
		rfqID.String(),
		time.Now().UTC().Format(time.RFC3339),
		"16.665",
		quote.NewQuoteDetails("USD", decimal.RequireFromString("10")),
		quote.NewQuoteDetails("MXN", decimal.RequireFromString("500")),
	)
	assert.NoError(t, err)

	toSign, err := quote.Digest()
	assert.NoError(t, err)

	wrongSignature, err := jws.Sign(toSign, wrongDID, jws.DetachedPayload(true))
	assert.NoError(t, err)

	quote.Signature = wrongSignature

	err = quote.Verify()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "does not match message metadata from")
}

func TestIsValidNext(t *testing.T) {

	pfiDID, _ := didjwk.Create()
	walletDID, _ := didjwk.Create()
	rfqID, _ := typeid.WithPrefix(rfq.Kind)

	q, err := quote.Create(
		pfiDID,
		walletDID.URI,
		rfqID.String(),
		time.Now().UTC().Format(time.RFC3339),
		"16.665",
		quote.NewQuoteDetails("USD", decimal.RequireFromString("10")),
		quote.NewQuoteDetails("MXN", decimal.RequireFromString("500")),
	)
	assert.NoError(t, err)

	assert.False(t, q.IsValidNext(rfq.Kind))
	assert.False(t, q.IsValidNext(quote.Kind))
	assert.False(t, q.IsValidNext(orderstatus.Kind))

	// quote can only be followed by order or close
	assert.True(t, q.IsValidNext(order.Kind))
	assert.True(t, q.IsValidNext(closemsg.Kind))
}
