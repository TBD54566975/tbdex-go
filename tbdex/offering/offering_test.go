package offering_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/TBD54566975/tbdex-go/tbdex/offering"
	"github.com/alecthomas/assert/v2"
	"github.com/tbd54566975/web5-go/dids/didjwk"
)

func TestCreate(t *testing.T) {
	_, err := offering.Create(
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
	assert.NoError(t, err)

	err = offering.Sign(bearerDID)
	assert.NoError(t, err)
}

func TestUnmarshal(t *testing.T) {
	bearerDID, err := didjwk.Create()
	assert.NoError(t, err)

	requiredPayoutDetails := offering.RequiredDetails(`{
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
	}`)

	o, err := offering.Create(
		offering.NewPayin(
			"USD",
			[]offering.PayinMethod{offering.NewPayinMethod("SQUAREPAY")},
		),
		offering.NewPayout(
			"USDC",
			[]offering.PayoutMethod{offering.NewPayoutMethod("STORED_BALANCE", 20*time.Minute, requiredPayoutDetails)},
		),
		"1.0",
	)

	assert.NoError(t, err)

	err = o.Sign(bearerDID)
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
