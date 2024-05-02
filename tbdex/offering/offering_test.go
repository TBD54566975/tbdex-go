package offering_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/TBD54566975/tbdex-go/tbdex/offering"
	"github.com/alecthomas/assert/v2"
	"github.com/tbd54566975/web5-go/dids/didjwk"
	"github.com/tbd54566975/web5-go/jws"
)

func TestCreate(t *testing.T) {
	pfiDID, err := didjwk.Create()
	assert.NoError(t, err)

	_, err = offering.Create(
		pfiDID,
		offering.NewPayin(
			"USD",
			[]offering.PayinMethod{offering.NewPayinMethod("SQUAREPAY")},
		),
		offering.NewPayout(
			"USDC",
			[]offering.PayoutMethod{offering.NewPayoutMethod("STORED_BALANCE", 20*time.Minute)},
		),
		"1.0",
	)

	assert.NoError(t, err)
}

func TestSign(t *testing.T) {
	bearerDID, err := didjwk.Create()
	assert.NoError(t, err)

	offering, err := offering.Create(
		bearerDID,
		offering.NewPayin(
			"USD",
			[]offering.PayinMethod{offering.NewPayinMethod("SQUAREPAY")},
		),
		offering.NewPayout(
			"USDC",
			[]offering.PayoutMethod{offering.NewPayoutMethod("STORED_BALANCE", 20*time.Minute)},
		),
		"1.0",
	)
	assert.NoError(t, err)
	assert.NotZero(t, offering.Signature)
}

func TestUnmarshal(t *testing.T) {
	bearerDID, err := didjwk.Create()
	assert.NoError(t, err)

	o, err := offering.Create(
		bearerDID,
		offering.NewPayin(
			"USD",
			[]offering.PayinMethod{offering.NewPayinMethod("SQUAREPAY")},
		),
		offering.NewPayout(
			"MXN",
			[]offering.PayoutMethod{offering.NewPayoutMethod(
				"STORED_BALANCE",
				20*time.Minute,
				offering.RequiredDetails(`{
					"$schema": "http://json-schema.org/draft-07/schema#",
					"additionalProperties": false,
					"properties": {
						"clabe": {
							"type": "string"
						},
						"fullName": {
							"type": "string"
						}
					},
					"required": ["clabe", "fullName"]
				}`),
			)},
		),
		"1.0",
	)

	assert.NoError(t, err)

	bytes, err := json.Marshal(o)
	assert.NoError(t, err)

	var o2 offering.Offering
	err = o2.UnmarshalJSON(bytes)
	assert.NoError(t, err)
}

func TestUnmarshal_Invalid(t *testing.T) {
	input := []byte(`{"doo": "doo"}`)

	var o offering.Offering
	err := json.Unmarshal(input, &o)
	assert.Error(t, err)
}

func TestVerify(t *testing.T) {
	bearerDID, err := didjwk.Create()
	assert.NoError(t, err)

	o, err := offering.Create(
		bearerDID,
		offering.NewPayin(
			"BTC",
			[]offering.PayinMethod{offering.NewPayinMethod("BTC_ADDRESS")},
		),
		offering.NewPayout(
			"USDC",
			[]offering.PayoutMethod{offering.NewPayoutMethod("STORED_BALANCE", 20*time.Minute)},
		),
		"60000.00",
	)

	assert.NoError(t, err)

	err = o.Verify()
	assert.NoError(t, err)
}

func TestVerify_InvalidSignature(t *testing.T) {
	bearerDID, err := didjwk.Create()
	assert.NoError(t, err)

	o, err := offering.Create(
		bearerDID,
		offering.NewPayin(
			"BTC",
			[]offering.PayinMethod{offering.NewPayinMethod("BTC_ADDRESS")},
		),
		offering.NewPayout(
			"USDC",
			[]offering.PayoutMethod{offering.NewPayoutMethod("STORED_BALANCE", 20*time.Minute)},
		),
		"60000.00",
	)

	assert.NoError(t, err)

	o.Signature = "invalid"

	err = o.Verify()
	assert.Error(t, err)
}

func TestVerify_SignedWithWrongDID(t *testing.T) {
	bearerDID, _ := didjwk.Create()
	wrongDID, _ := didjwk.Create()

	o, err := offering.Create(
		bearerDID,
		offering.NewPayin(
			"BTC",
			[]offering.PayinMethod{offering.NewPayinMethod("BTC_ADDRESS")},
		),
		offering.NewPayout(
			"USDC",
			[]offering.PayoutMethod{offering.NewPayoutMethod("STORED_BALANCE", 20*time.Minute)},
		),
		"60000.00",
	)

	assert.NoError(t, err)

	toSign, err := o.Digest()
	assert.NoError(t, err)

	wrongSignature, err := jws.Sign(toSign, wrongDID, jws.DetachedPayload(true))
	assert.NoError(t, err)

	o.Signature = wrongSignature

	err = o.Verify()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "does not match resource metadata from")
}
