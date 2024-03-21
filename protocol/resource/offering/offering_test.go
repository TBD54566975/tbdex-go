package offering_test

import (
	"testing"
	"time"

	"github.com/TBD54566975/tbdex-go/protocol/resource/offering"
	"github.com/alecthomas/assert/v2"
	"github.com/tbd54566975/web5-go/dids/didjwk"
)

func TestOffering(t *testing.T) {
	_, err := offering.Create(
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
	)

	assert.NoError(t, err)
}

func TestSign(t *testing.T) {
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
	)
	assert.NoError(t, err)

	bearerDID, err := didjwk.Create()
	assert.NoError(t, err)

	err = offering.Sign(bearerDID)
	assert.NoError(t, err)
}
