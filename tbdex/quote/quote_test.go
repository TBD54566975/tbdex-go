package quote_test

import (
	"testing"
	"time"

	"github.com/TBD54566975/tbdex-go/tbdex/quote"
	"github.com/TBD54566975/tbdex-go/tbdex/rfq"
	"github.com/alecthomas/assert/v2"
	"github.com/tbd54566975/web5-go/dids/didjwk"
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
		quote.NewQuoteDetails("USD", "10"),
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
