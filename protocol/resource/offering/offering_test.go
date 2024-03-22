package offering_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/TBD54566975/tbdex-go/protocol/resource/offering"
	"github.com/alecthomas/assert/v2"
	"github.com/tbd54566975/web5-go/dids/didjwk"
)

func TestCreate(t *testing.T) {
	bearerDID, err := didjwk.Create()
	assert.NoError(t, err)

	_, err = offering.Create(
		offering.WithPayin(
			"USD",
			offering.WithPayinMethod("SQUAREPAY"),
		),
		offering.WithPayout(
			"USDC",
			offering.WithPayoutMethod(
				"STORED_BALANCE",
				20*time.Minute,
			),
		),
		"1.0",
		bearerDID.URI,
	)

	assert.NoError(t, err)
}

func TestSign(t *testing.T) {
	bearerDID, err := didjwk.Create()
	assert.NoError(t, err)

	offering, err := offering.Create(
		offering.WithPayin(
			"USD",
			offering.WithPayinMethod("SQUAREPAY"),
		),
		offering.WithPayout(
			"USDC",
			offering.WithPayoutMethod(
				"STORED_BALANCE",
				20*time.Minute,
			),
		),
		"1.0",
		bearerDID.URI,
	)
	assert.NoError(t, err)

	err = offering.Sign(bearerDID)
	assert.NoError(t, err)
}

func TestValidate(t *testing.T) {
	bearerDID, _ := didjwk.Create()

	offeringMessage, _ := offering.Create(
		offering.WithPayin(
			"USD",
			offering.WithPayinMethod("SQUAREPAY"),
		),
		offering.WithPayout(
			"USDC",
			offering.WithPayoutMethod(
				"STORED_BALANCE",
				20*time.Minute,
			),
		),
		"1.0",
		bearerDID.URI,
	)

	offeringJSON, err := json.Marshal(offeringMessage)
	assert.NoError(t, err)

	err = offering.Validate(offeringJSON)
	assert.NoError(t, err)
}
