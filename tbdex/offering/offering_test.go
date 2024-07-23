package offering_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/TBD54566975/tbdex-go/tbdex/offering"
	"github.com/alecthomas/assert/v2"
	"github.com/tbd54566975/web5-go/dids/didjwk"
	"github.com/tbd54566975/web5-go/jws"
	"github.com/tbd54566975/web5-go/pexv2"
)

func TestCreate(t *testing.T) {
	pfiDID, err := didjwk.Create()
	assert.NoError(t, err)

	pd := pexv2.PresentationDefinition{
		ID:      "foo",
		Name:    "kyccredential",
		Purpose: "To verify the identity of the user",
		InputDescriptors: []pexv2.InputDescriptor{
			{
				ID:      "1",
				Name:    "KYC Information",
				Purpose: "To verify the identity of the user",
				Constraints: pexv2.Constraints{
					Fields: []pexv2.Field{
						{
							Path: []string{"$.type[0]"},
							Filter: &pexv2.Filter{
								Type:    "string",
								Pattern: "KYC",
							},
						},
					},
				},
			},
		},
	}

	_, err = offering.Create(
		offering.NewPayin(
			"USD",
			[]offering.PayinMethod{
				offering.NewPayinMethod(
					"DEBIT_CARD",
					offering.RequiredDetails(`{
					"$schema": "http://json-schema.org/draft-07/schema#",
					"type": "object",
					"properties": {
						"cardNumber": {
						"type": "string",
						"description": "The 16-digit debit card number",
						"minLength": 16,
						"maxLength": 16
						},
						"expiryDate": {
						"type": "string",
						"description": "The expiry date of the card in MM/YY format",
						"pattern": "^(0[1-9]|1[0-2])\\/([0-9]{2})$"
						},
						"cardHolderName": {
						"type": "string",
						"description": "Name of the cardholder as it appears on the card"
						},
						"cvv": {
						"type": "string",
						"description": "The 3-digit CVV code",
						"minLength": 3,
						"maxLength": 3
						}
					},
					"required": [
						"cardNumber",
						"expiryDate",
						"cardHolderName",
						"cvv"
					],
					"additionalProperties": false
					}`))},
			offering.Min("0.1"),
			offering.Max("1000"),
		),
		offering.NewPayout(
			"USDC",
			[]offering.PayoutMethod{offering.NewPayoutMethod("STORED_BALANCE", 20*time.Minute)},
			offering.Max("5000"),
		),
		"1.0",
		offering.NewCancellationDetails(false),
		offering.From(pfiDID),
		offering.RequiredClaims(pd),
	)

	assert.NoError(t, err)
}

func TestSign(t *testing.T) {
	bearerDID, err := didjwk.Create()
	assert.NoError(t, err)

	offering, err := offering.Create(
		offering.NewPayin(
			"USD",
			[]offering.PayinMethod{offering.NewPayinMethod("SQUAREPAY")},
		),
		offering.NewPayout(
			"USDC",
			[]offering.PayoutMethod{offering.NewPayoutMethod("STORED_BALANCE", 20*time.Minute)},
		),
		"1.0",
		offering.NewCancellationDetails(false),
		offering.From(bearerDID),
	)
	assert.NoError(t, err)
	assert.NotZero(t, offering.Signature)
}

func TestUnmarshal(t *testing.T) {
	bearerDID, err := didjwk.Create()
	assert.NoError(t, err)

	o, err := offering.Create(
		offering.NewPayin(
			"USD",
			[]offering.PayinMethod{offering.NewPayinMethod("SQUAREPAY")},
		),
		offering.NewPayout(
			"MXN",
			[]offering.PayoutMethod{offering.NewPayoutMethod(
				"STORED_BALANCE",
				20*time.Minute,
				offering.RequiredDetails(`{
					"$schema": "http://json-schema.org/draft-07/schema#",
					"additionalProperties": false,
					"properties": {
						"clabe": {
							"type": "string"
						},
						"fullName": {
							"type": "string"
						}
					},
					"required": ["clabe", "fullName"]
				}`),
			)},
		),
		"1.0",
		offering.NewCancellationDetails(false),
		offering.From(bearerDID),
	)

	assert.NoError(t, err)

	bytes, err := json.Marshal(o)
	assert.NoError(t, err)

	var o2 offering.Offering
	err = json.Unmarshal(bytes, &o2)
	assert.NoError(t, err)
}

func TestUnmarshal_Invalid(t *testing.T) {
	input := []byte(`{"doo": "doo"}`)

	var o offering.Offering
	err := json.Unmarshal(input, &o)
	assert.Error(t, err)
}

func TestUnmarshal_Empty(t *testing.T) {
	input := []byte(`{"metadata":{},"data":{},"signature":"eyJhbGciOiJFZERTQSIsImtpZCI6ImRpZDpkaHQ6M3doZnRncGJkamloeDl6ZTl0ZG41NzV6cXptNHF3Y2NldG5mMXliaWlidXphZDdycm15eSMwIn0..ZvoVDuSrqqdXsSXgqB-U26tAU1WqUqqU_KpD1KvdYocIcmTsshjUASEwM_lUz1UnGglqkWeCIrHqrm9NNGDqBw"}`)

	var o offering.Offering
	_ = json.Unmarshal(input, &o)

	assert.Zero(t, o.Metadata)
	assert.Zero(t, o.Data)
}

func TestVerify(t *testing.T) {
	bearerDID, err := didjwk.Create()
	assert.NoError(t, err)

	o, err := offering.Create(
		offering.NewPayin(
			"BTC",
			[]offering.PayinMethod{offering.NewPayinMethod("BTC_ADDRESS")},
		),
		offering.NewPayout(
			"USDC",
			[]offering.PayoutMethod{offering.NewPayoutMethod("STORED_BALANCE", 20*time.Minute)},
		),
		"60000.00",
		offering.NewCancellationDetails(false),
		offering.From(bearerDID),
	)

	assert.NoError(t, err)

	err = o.Verify()
	assert.NoError(t, err)
}

func TestVerify_InvalidSignature(t *testing.T) {
	bearerDID, err := didjwk.Create()
	assert.NoError(t, err)

	o, err := offering.Create(
		offering.NewPayin(
			"BTC",
			[]offering.PayinMethod{offering.NewPayinMethod("BTC_ADDRESS")},
		),
		offering.NewPayout(
			"USDC",
			[]offering.PayoutMethod{offering.NewPayoutMethod("STORED_BALANCE", 20*time.Minute)},
		),
		"60000.00",
		offering.NewCancellationDetails(false),
		offering.From(bearerDID),
	)

	assert.NoError(t, err)

	o.Signature = "invalid"

	err = o.Verify()
	assert.Error(t, err)
}

func TestVerify_SignedWithWrongDID(t *testing.T) {
	bearerDID, _ := didjwk.Create()
	wrongDID, _ := didjwk.Create()

	o, err := offering.Create(
		offering.NewPayin(
			"BTC",
			[]offering.PayinMethod{offering.NewPayinMethod("BTC_ADDRESS")},
		),
		offering.NewPayout(
			"USDC",
			[]offering.PayoutMethod{offering.NewPayoutMethod("STORED_BALANCE", 20*time.Minute)},
		),
		"60000.00",
		offering.NewCancellationDetails(false),
		offering.From(bearerDID),
	)

	assert.NoError(t, err)

	toSign, err := o.Digest()
	assert.NoError(t, err)

	wrongSignature, err := jws.Sign(toSign, wrongDID, jws.DetachedPayload(true))
	assert.NoError(t, err)

	o.Signature = wrongSignature

	err = o.Verify()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "does not match resource metadata from")
}
