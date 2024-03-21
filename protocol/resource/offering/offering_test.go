package offering_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/TBD54566975/tbdex-go/protocol/resource/offering"
	"github.com/alecthomas/assert/v2"
)

func TestOffering(t *testing.T) {
	o, err := offering.Create(
		offering.WithPayin(
			"USD",
			offering.WithPayinMethod("SQUAREPAY"),
		),
		offering.WithPayout(
			"USDC",
			offering.WithPayoutMethod(
				"STORED_BALANCE",
				time.Duration(20*time.Minute),
			),
		),
		"1.0",
	)

	assert.NoError(t, err)
	j, err := json.MarshalIndent(o, "", "  ")
	assert.NoError(t, err)

	t.Log(string(j))
}
