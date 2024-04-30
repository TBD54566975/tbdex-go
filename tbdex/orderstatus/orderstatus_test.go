package orderstatus_test

import (
	"encoding/json"
	"testing"

	"github.com/TBD54566975/tbdex-go/tbdex/orderstatus"
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
		"processing",
	)

	assert.NoError(t, err)
	assert.NotZero(t, orderstatus.Data.OrderStatus)
	assert.NotZero(t, orderstatus.Signature)
}

func TestUnmarshalJSON(t *testing.T) {
	pfiDID, _ := didjwk.Create()
	walletDID, _ := didjwk.Create()
	rfqID, _ := typeid.WithPrefix(rfq.Kind)

	message, _ := orderstatus.Create(
		pfiDID,
		walletDID.URI,
		rfqID.String(),
		"processing",
	)

	bytes, err := json.Marshal(message)
	assert.NoError(t, err)

	os := orderstatus.OrderStatus{}
	err = os.UnmarshalJSON(bytes)
	assert.NoError(t, err)
}

func TestUnmarshal_Invalid(t *testing.T) {
	input := []byte(`{"doo": "doo"}`)

	os := orderstatus.OrderStatus{}
	err := os.UnmarshalJSON(input)
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
		"processing",
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
		"processing",
	)

	os.Data.OrderStatus = "failed"

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
		"processing",
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
		"processing",
	)

	toSign, err := os.Digest()
	assert.NoError(t, err)

	wrongSignature, err := jws.Sign(toSign, wrongDID, jws.DetachedPayload(true))
	assert.NoError(t, err)

	os.Signature = wrongSignature

	err = os.Verify()
	assert.Error(t, err)
}
