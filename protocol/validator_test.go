package protocol

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/TBD54566975/tbdex-go/protocol/resource/offering"
	"github.com/alecthomas/assert"
	"github.com/santhosh-tekuri/jsonschema/v5"
	_ "github.com/santhosh-tekuri/jsonschema/v5/httploader"
	"github.com/tbd54566975/web5-go/dids/didjwk"
)

func TestValidator(t *testing.T) {
    definitionsPath := "../tbdex/hosted/json-schemas/definitions.json"
	resourcePath := "../tbdex/hosted/json-schemas/resource.schema.json"

    definitionsFile, err := os.Open(definitionsPath)
    resourceFile, err := os.Open(resourcePath)
    if err != nil {
		panic(err)
	}


	// 1. Load the JSON schema
	compiler := jsonschema.NewCompiler()
	err = compiler.AddResource(definitionsPath, definitionsFile)

	assert.NoError(t, err)

    if err := compiler.AddResource(resourcePath, resourceFile); err != nil {
		panic(err)
	}
    schema, err := compiler.Compile(resourcePath)
	assert.NoError(t, err)

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

    offeringJSON, err := json.Marshal(offering)
	if err != nil {
		panic(err)
	}

    // 4. Validate the JSON against the schema
	var v interface{}
	err = json.Unmarshal([]byte(offeringJSON), &v)
    assert.NoError(t, err)

	err = schema.Validate(v)
    assert.NoError(t, err)
}

func TestCompiler(t *testing.T) {
    compiler, err := Compiler()
    assert.NoError(t, err)

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

    offeringJSON, err := json.Marshal(offering)
	if err != nil {
		panic(err)
	}

    // 4. Validate the JSON against the schema
    schema, err := compiler.Compile(schemaPaths["resource"])
    assert.NoError(t, err)

    var v interface{}
	err = json.Unmarshal([]byte(offeringJSON), &v)
    assert.NoError(t, err)

	err = schema.Validate(v)
    assert.NoError(t, err)
}
