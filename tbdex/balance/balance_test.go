package balance_test

import (
	"encoding/json"
	"testing"

	"github.com/TBD54566975/tbdex-go/tbdex/balance"
	"github.com/alecthomas/assert/v2"
	"github.com/tbd54566975/web5-go/dids/didjwk"
	"github.com/tbd54566975/web5-go/jws"
)

func TestCreate(t *testing.T) {
	bearerDID, err := didjwk.Create()
	assert.NoError(t, err)

	b, err := balance.Create(bearerDID, "USD", "100.00")
	assert.NoError(t, err)
	assert.NotZero(t, b.Available)
	assert.NotZero(t, b.Signature)
}

func TestUnmarshal(t *testing.T) {
	bearerDID, err := didjwk.Create()
	assert.NoError(t, err)

	b, _ := balance.Create(bearerDID, "USD", "100.00")

	bytes, err := json.Marshal(b)
	assert.NoError(t, err)

	var b2 balance.Balance
	err = b2.UnmarshalJSON(bytes)
	assert.NoError(t, err)
}

func TestUnmarshal_Invalid(t *testing.T) {
	cases := []string{
		"hi",
		"{}",
		"[]",
		`{"metadata": {"kind": "hoarder"}}`,
		`{"metadata": {"kind": "hoarder", "from": "hehe"}}`,
	}

	for _, v := range cases {
		var b balance.Balance
		err := b.UnmarshalJSON([]byte(v))
		assert.Error(t, err)
	}
}

func TestVerify(t *testing.T) {
	bearerDID, err := didjwk.Create()
	assert.NoError(t, err)

	b, err := balance.Create(bearerDID, "USD", "100.00")

	assert.NoError(t, err)

	err = b.Verify()
	assert.NoError(t, err)
}

func TestVerify_InvalidSignature(t *testing.T) {
	bearerDID, err := didjwk.Create()
	assert.NoError(t, err)

	b, err := balance.Create(bearerDID, "USD", "100.00")
	assert.NoError(t, err)

	b.Signature = "invalid"

	err = b.Verify()
	assert.Error(t, err)
}

func TestVerify_SignedWithWrongDID(t *testing.T) {
	bearerDID, _ := didjwk.Create()
	wrongDID, _ := didjwk.Create()

	b, err := balance.Create(bearerDID, "USD", "100.00")
	assert.NoError(t, err)

	toSign, err := b.Digest()
	assert.NoError(t, err)

	wrongSignature, err := jws.Sign(toSign, wrongDID, jws.DetachedPayload(true))
	assert.NoError(t, err)

	b.Signature = wrongSignature

	err = b.Verify()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "does not match resource metadata from")
}
