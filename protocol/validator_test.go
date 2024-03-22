package protocol

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/TBD54566975/tbdex-go/protocol/resource/offering"
	"github.com/alecthomas/assert"
	"github.com/tbd54566975/web5-go/dids/didjwk"
)

func TestCompiler(t *testing.T) {
	bearerDID, err := didjwk.Create()
	assert.NoError(t, err)

	offering, _ := offering.Create(
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

	offeringJSON, err := json.Marshal(offering)
	assert.NoError(t, err)

	var v interface{}
	err = json.Unmarshal(offeringJSON, &v)
	assert.NoError(t, err)

    schema := ValidatorMap["resource"]
	err = schema.Validate(v)
	assert.NoError(t, err)
}
