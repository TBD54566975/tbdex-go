package tbdex_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/TBD54566975/tbdex-go/tbdex"
	"github.com/alecthomas/assert/v2"
	"github.com/tbd54566975/web5-go/dids/didjwk"
)

func TestCreate(t *testing.T) {
	_, err := tbdex.CreateOffering(
		tbdex.WithOfferingPayin(
			"USD",
			tbdex.WithOfferingPayinMethod("SQUAREPAY"),
		),
		tbdex.WithOfferingPayout(
			"USDC",
			tbdex.WithOfferingPayoutMethod(
				"STORED_BALANCE",
				20*time.Minute,
			),
		),
		"1.0",
	)

	assert.NoError(t, err)
}

func TestSign(t *testing.T) {
	bearerDID, err := didjwk.Create()
	assert.NoError(t, err)

	offering, err := tbdex.CreateOffering(
		tbdex.WithOfferingPayin(
			"USD",
			tbdex.WithOfferingPayinMethod("SQUAREPAY"),
		),
		tbdex.WithOfferingPayout(
			"USDC",
			tbdex.WithOfferingPayoutMethod(
				"STORED_BALANCE",
				20*time.Minute,
			),
		),
		"1.0",
	)
	assert.NoError(t, err)

	err = offering.Sign(bearerDID)
	assert.NoError(t, err)
}

func TestUnmarshal(t *testing.T) {
	bearerDID, _ := didjwk.Create()

	o, _ := tbdex.CreateOffering(
		tbdex.WithOfferingPayin(
			"USD",
			tbdex.WithOfferingPayinMethod("SQUAREPAY"),
		),
		tbdex.WithOfferingPayout(
			"USDC",
			tbdex.WithOfferingPayoutMethod(
				"STORED_BALANCE",
				20*time.Minute,
			),
		),
		"1.0",
	)

	err := o.Sign(bearerDID)
	assert.NoError(t, err)

	bytes, err := json.Marshal(o)
	assert.NoError(t, err)

	var o2 tbdex.Offering
	err = o2.UnmarshalJSON(bytes)
	assert.NoError(t, err)
}

func TestUnmarshal_Invalid(t *testing.T) {
	input := []byte(`{"doo": "doo"}`)

	var o tbdex.Offering
	err := json.Unmarshal(input, &o)

	assert.Error(t, err)
}
