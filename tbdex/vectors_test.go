package tbdex_test

// type vector struct {
// 	Input  string `json:"input"`
// 	Output any    `json:"output"`
// }

// func readVector(filename string) vector {
// 	file, err := os.ReadFile("../spec/hosted/test-vectors/protocol/vectors/" + filename)
// 	if err != nil {
// 		log.Fatalf("Failed to read %s: %v", filename, err)
// 	}

// 	// Unmarshal JSON data into the struct
// 	var v vector
// 	err = json.Unmarshal(file, &v)
// 	if err != nil {
// 		log.Fatalf("Failed to unmarshal %s: %v", filename, err)
// 	}

// 	return v
// }

// func TestOfferingVectors(t *testing.T) {
// 	vector := readVector("parse-offering.json")
// 	res := offering.Offering{}
// 	err := res.Parse([]byte(vector.Input))

// 	assert.NoError(t, err)
// }

// func TestBalanceVectors(t *testing.T) {
// 	vector := readVector("parse-balance.json")
// 	_, err := balance.Parse([]byte(vector.Input))

// 	assert.NoError(t, err)
// }

// func TestRFQVectors(t *testing.T) {
// 	vector := readVector("parse-rfq.json")
// 	_, err := rfq.Parse([]byte(vector.Input), true)

// 	assert.NoError(t, err)

// 	vector = readVector("parse-rfq-omit-private-data.json")
// 	_, err = rfq.Parse([]byte(vector.Input), false)

// 	assert.NoError(t, err)
// }

// func TestQuoteVectors(t *testing.T) {
// 	vector := readVector("parse-quote.json")
// 	_, err := quote.Parse([]byte(vector.Input))

// 	assert.NoError(t, err)
// }

// func TestOrderVectors(t *testing.T) {
// 	vector := readVector("parse-order.json")
// 	_, err := order.Parse([]byte(vector.Input))

// 	assert.NoError(t, err)
// }

// func TestOrderStatusVectors(t *testing.T) {
// 	vector := readVector("parse-orderstatus.json")
// 	_, err := orderstatus.Parse([]byte(vector.Input))

// 	assert.NoError(t, err)
// }

// func TestCloseVectors(t *testing.T) {
// 	vector := readVector("parse-close.json")
// 	_, err := closemsg.Parse([]byte(vector.Input))

// 	assert.NoError(t, err)
// }
