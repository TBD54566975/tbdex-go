package orderstatus_test

import (
	"encoding/json"
	"testing"

	"github.com/TBD54566975/tbdex-go/tbdex/closemsg"
	"github.com/TBD54566975/tbdex-go/tbdex/order"
	"github.com/TBD54566975/tbdex-go/tbdex/orderstatus"
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

	orderstatus, err := orderstatus.Create(
		pfiDID,
		walletDID.URI,
		rfqID.String(),
		"PAYIN_INITIATED",
		orderstatus.Details("CC Payment Initiated"),
	)

	assert.NoError(t, err)

	assert.NotZero(t, orderstatus.Data.Status)
	assert.Equal(t, "PAYIN_INITIATED", orderstatus.Data.Status)
	assert.NotZero(t, orderstatus.Signature)
}

func TestUnmarshal(t *testing.T) {
	pfiDID, _ := didjwk.Create()
	walletDID, _ := didjwk.Create()
	rfqID, _ := typeid.WithPrefix(rfq.Kind)

	message, _ := orderstatus.Create(
		pfiDID,
		walletDID.URI,
		rfqID.String(),
		"PAYIN_INITIATED",
	)

	bytes, err := json.Marshal(message)
	assert.NoError(t, err)

	os := orderstatus.OrderStatus{}
	err = json.Unmarshal(bytes, &os)
	assert.NoError(t, err)
}

func TestUnmarshal_Empty(t *testing.T) {
	input := []byte(`{"metadata":{},"data":{},"signature":"eyJhbGciOiJFZERTQSIsImtpZCI6ImRpZDpkaHQ6M3doZnRncGJkamloeDl6ZTl0ZG41NzV6cXptNHF3Y2NldG5mMXliaWlidXphZDdycm15eSMwIn0..ZvoVDuSrqqdXsSXgqB-U26tAU1WqUqqU_KpD1KvdYocIcmTsshjUASEwM_lUz1UnGglqkWeCIrHqrm9NNGDqBw"}`)

	os := orderstatus.OrderStatus{}

	_ = json.Unmarshal(input, &os)

	assert.Zero(t, os.Metadata)
	assert.Zero(t, os.Data)
}

func TestUnmarshal_Invalid(t *testing.T) {
	input := []byte(`{"doo": "doo"}`)

	os := orderstatus.OrderStatus{}
	err := json.Unmarshal(input, &os)
	assert.Error(t, err)
}

func TestVerify(t *testing.T) {
	pfiDID, _ := didjwk.Create()
	walletDID, _ := didjwk.Create()
	rfqID, _ := typeid.WithPrefix(rfq.Kind)

	os, _ := orderstatus.Create(
		pfiDID,
		walletDID.URI,
		rfqID.String(),
		"PAYIN_INITIATED",
	)

	err := os.Verify()
	assert.NoError(t, err)
}

func TestVerify_FailsChangedPayload(t *testing.T) {
	pfiDID, _ := didjwk.Create()
	walletDID, _ := didjwk.Create()
	rfqID, _ := typeid.WithPrefix(rfq.Kind)

	os, _ := orderstatus.Create(
		pfiDID,
		walletDID.URI,
		rfqID.String(),
		"PAYIN_INITIATED",
	)

	os.Data.Status = "PAYIN_FAILED"

	err := os.Verify()
	assert.Error(t, err)
}

func TestVerify_InvalidSignature(t *testing.T) {
	pfiDID, _ := didjwk.Create()
	walletDID, _ := didjwk.Create()
	rfqID, _ := typeid.WithPrefix(rfq.Kind)

	os, _ := orderstatus.Create(
		pfiDID,
		walletDID.URI,
		rfqID.String(),
		"PAYIN_INITIATED",
	)

	os.Signature = "Invalid"

	err := os.Verify()
	assert.Error(t, err)
}

func TestVerify_SignedWithWrongDID(t *testing.T) {
	pfiDID, _ := didjwk.Create()
	walletDID, _ := didjwk.Create()
	wrongDID, _ := didjwk.Create()
	rfqID, _ := typeid.WithPrefix(rfq.Kind)

	os, _ := orderstatus.Create(
		pfiDID,
		walletDID.URI,
		rfqID.String(),
		"PAYIN_INITIATED",
	)

	toSign, err := os.Digest()
	assert.NoError(t, err)

	wrongSignature, err := jws.Sign(toSign, wrongDID, jws.DetachedPayload(true))
	assert.NoError(t, err)

	os.Signature = wrongSignature

	err = os.Verify()
	assert.Error(t, err)
}

func TestIsValidNext(t *testing.T) {

	pfiDID, _ := didjwk.Create()
	walletDID, _ := didjwk.Create()
	rfqID, _ := typeid.WithPrefix(rfq.Kind)

	os, err := orderstatus.Create(
		pfiDID,
		walletDID.URI,
		rfqID.String(),
		"PAYIN_INITIATED",
	)

	assert.NoError(t, err)

	assert.False(t, os.IsValidNext(rfq.Kind))
	assert.False(t, os.IsValidNext(quote.Kind))
	assert.False(t, os.IsValidNext(order.Kind))

	// orderstatus can only be followed by another orderstatus or close
	assert.True(t, os.IsValidNext(orderstatus.Kind))
	assert.True(t, os.IsValidNext(closemsg.Kind))
}
