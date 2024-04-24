package tbdex_test

import (
	"embed"
	"encoding/json"
	"testing"

	"github.com/TBD54566975/tbdex-go/tbdex/offering"
	"github.com/TBD54566975/tbdex-go/tbdex/quote"
	"github.com/TBD54566975/tbdex-go/tbdex/rfq"
	"github.com/alecthomas/assert/v2"
)

//go:embed vectors
var embeddedVectors embed.FS

const (
	vectorsDir = "vectors/"
)

type vector struct {
	Input  string `json:"input"`
	Output any    `json:"output"`
}

func readVector(filename string) vector {
	file, err := embeddedVectors.ReadFile(vectorsDir + filename)
	if err != nil {
		panic(err)
	}

	// Unmarshal JSON data into the struct
	var v vector
	err = json.Unmarshal(file, &v)
	if err != nil {
		panic(err)
	}

	return v
}

func TestOfferingVectors(t *testing.T) {
	vector := readVector("parse-offering.json")
	res := offering.Offering{}
	err := res.Parse([]byte(vector.Input))

	assert.NoError(t, err)
}

func TestRFQVectors(t *testing.T) {
	vector := readVector("parse-rfq.json")
	res := rfq.RFQ{}
	err := res.Parse([]byte(vector.Input), true)

	assert.NoError(t, err)

	vector = readVector("parse-rfq-omit-private-data.json")
	res = rfq.RFQ{}
	err = res.Parse([]byte(vector.Input), false)

	assert.NoError(t, err)
}

func TestQuoteVectors(t *testing.T) {
	vector := readVector("parse-quote.json")
	res := quote.Quote{}
	err := res.Parse([]byte(vector.Input), true)

	assert.NoError(t, err)
}
