package tbdex_test

import (
	"embed"
	"encoding/json"
	"testing"

	"github.com/TBD54566975/tbdex-go/tbdex/offering"
	"github.com/TBD54566975/tbdex-go/tbdex/rfq"
	"github.com/alecthomas/assert/v2"
)

//go:embed vectors
var embeddedVectors embed.FS

const (
	vectorsDir = "vectors/"
)

type offeringVector struct {
	Input  string            `json:"input"`
	Output offering.Offering `json:"output"`
}

type rfqVector struct {
	Input  string  `json:"input"`
	Output rfq.RFQ `json:"output"`
}

// var vectorsMap map[string]vector = make(map[string]vector)

func readVector[T any](filename string) T {
	file, err := embeddedVectors.ReadFile(vectorsDir + filename)
	if err != nil {
		panic(err)
	}

	// Unmarshal JSON data into the struct
	var vector T
	err = json.Unmarshal(file, &vector)
	if err != nil {
		panic(err)
	}

	return vector
}

func TestOfferingVectors(t *testing.T) {
	vector := readVector[offeringVector]("parse-offering.json")
	input := offering.Offering{}
	err := input.UnmarshalJSON([]byte(vector.Input))

	assert.NoError(t, err)
	assert.Equal(t, input, vector.Output)
}

func TestRFQVectors(t *testing.T) {
	vector := readVector[rfqVector]("parse-rfq.json")
	input := rfq.RFQ{}
	err := input.ValidateAndUnmarshalJSON([]byte(vector.Input), false)

	assert.NoError(t, err)
	assert.Equal(t, input, vector.Output)

	vector = readVector[rfqVector]("parse-rfq-omit-private-data.json")
	input = rfq.RFQ{}
	err = input.ValidateAndUnmarshalJSON([]byte(vector.Input), false)

	assert.NoError(t, err)
	assert.Equal(t, input, vector.Output)
}
