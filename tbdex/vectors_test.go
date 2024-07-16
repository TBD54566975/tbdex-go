package tbdex_test

import (
	"encoding/json"
	"log"
	"os"
	"testing"

	"github.com/TBD54566975/tbdex-go/tbdex/balance"
	"github.com/TBD54566975/tbdex-go/tbdex/cancel"
	"github.com/TBD54566975/tbdex-go/tbdex/closemsg"
	"github.com/TBD54566975/tbdex-go/tbdex/offering"
	"github.com/TBD54566975/tbdex-go/tbdex/order"
	"github.com/TBD54566975/tbdex-go/tbdex/orderstatus"
	"github.com/TBD54566975/tbdex-go/tbdex/quote"
	"github.com/TBD54566975/tbdex-go/tbdex/rfq"
	"github.com/alecthomas/assert/v2"
)

type vector struct {
	Input  string `json:"input"`
	Output any    `json:"output"`
}

func readVector(filename string) vector {
	file, err := os.ReadFile("../spec/hosted/test-vectors/protocol/vectors/" + filename)
	if err != nil {
		log.Fatalf("Failed to read %s: %v", filename, err)
	}

	// Unmarshal JSON data into the struct
	var v vector
	err = json.Unmarshal(file, &v)
	if err != nil {
		log.Fatalf("Failed to unmarshal %s: %v", filename, err)
	}

	return v
}

func TestOfferingTbdexTestVectors(t *testing.T) {
	vector := readVector("parse-offering.json")
	res := offering.Offering{}
	err := res.Parse([]byte(vector.Input))

	assert.NoError(t, err)
}

func TestBalanceTbdexTestVectors(t *testing.T) {
	vector := readVector("parse-balance.json")
	_, err := balance.Parse([]byte(vector.Input))

	assert.NoError(t, err)
}

func TestRFQTbdexTestVectors(t *testing.T) {
	vector := readVector("parse-rfq.json")
	_, err := rfq.Parse([]byte(vector.Input))

	assert.NoError(t, err)

	vector = readVector("parse-rfq-omit-private-data.json")
	_, err = rfq.Parse([]byte(vector.Input))

	assert.NoError(t, err)
}

func TestQuoteTbdexTestVectors(t *testing.T) {
	vector := readVector("parse-quote.json")
	_, err := quote.Parse([]byte(vector.Input))

	assert.NoError(t, err)
}

func TestOrderTbdexTestVectors(t *testing.T) {
	vector := readVector("parse-order.json")
	_, err := order.Parse([]byte(vector.Input))

	assert.NoError(t, err)
}

func TestOrderStatusTbdexTestVectors(t *testing.T) {
	vector := readVector("parse-orderstatus.json")
	_, err := orderstatus.Parse([]byte(vector.Input))

	assert.NoError(t, err)
}

func TestCloseTbdexTestVectors(t *testing.T) {
	vector := readVector("parse-close.json")
	_, err := closemsg.Parse([]byte(vector.Input))

	assert.NoError(t, err)
}

func TestCancelTbdexTestVectors(t *testing.T) {
	vector := readVector("parse-cancel.json")
	_, err := cancel.Parse([]byte(vector.Input))

	assert.NoError(t, err)
}
