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

func parse_offering(t *testing.T) {
	vector := readVector("parse-offering.json")
	res := offering.Offering{}
	err := res.Parse([]byte(vector.Input))

	assert.NoError(t, err)
}

func parse_balance(t *testing.T) {
	vector := readVector("parse-balance.json")
	_, err := balance.Parse([]byte(vector.Input))

	assert.NoError(t, err)
}

func parse_rfq(t *testing.T) {
	vector := readVector("parse-rfq.json")
	_, err := rfq.Parse([]byte(vector.Input))

	assert.NoError(t, err)

	vector = readVector("parse-rfq-omit-private-data.json")
	_, err = rfq.Parse([]byte(vector.Input))

	assert.NoError(t, err)
}

func parse_quote(t *testing.T) {
	vector := readVector("parse-quote.json")
	_, err := quote.Parse([]byte(vector.Input))

	assert.NoError(t, err)
}

func parse_order(t *testing.T) {
	vector := readVector("parse-order.json")
	_, err := order.Parse([]byte(vector.Input))

	assert.NoError(t, err)
}

func parse_orderstatus(t *testing.T) {
	vector := readVector("parse-orderstatus.json")
	_, err := orderstatus.Parse([]byte(vector.Input))

	assert.NoError(t, err)
}

func parse_close(t *testing.T) {
	vector := readVector("parse-close.json")
	_, err := closemsg.Parse([]byte(vector.Input))

	assert.NoError(t, err)
}

func parse_cancel(t *testing.T) {
	vector := readVector("parse-cancel.json")
	_, err := cancel.Parse([]byte(vector.Input))

	assert.NoError(t, err)
}

func TestAllParsers(t *testing.T) {
	t.Run("parse_offering", parse_offering)
	t.Run("parse_balance", parse_balance)
	t.Run("parse_rfq", parse_rfq)
	t.Run("parse_quote", parse_quote)
	t.Run("parse_order", parse_order)
	t.Run("parse_orderstatus", parse_orderstatus)
	t.Run("parse_close", parse_close)
	t.Run("parse_cancel", parse_cancel)
}
