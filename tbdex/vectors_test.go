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

type vectorDetails struct {
	Filename string
	Type     any
}

type tbdexType interface {
	rfq.RFQ | offering.Offering
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

func RFQVectors(t *testing.T) {
	vector := readVector[rfqVector]("parse-rfq.json")
	input := rfq.RFQ{}
	input.ValidateAndUnmarshalJSON([]byte(vector.Input), false)

	assert.Equal(t, input, vector.Output)

	// rfq := rfq.RFQ{}
	// rfq.ValidateAndUnmarshalJSON([]byte(vectorsMap["parse-rfq-omit-private-data"].Input), false)

	// assert.Equal(t, rfq, vectorsMap["parse-rfq"].Output)
}
