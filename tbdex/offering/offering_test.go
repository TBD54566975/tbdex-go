package offering_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/TBD54566975/tbdex-go/tbdex/offering"
	"github.com/alecthomas/assert/v2"
	"github.com/tbd54566975/web5-go/dids/did"
	"github.com/tbd54566975/web5-go/dids/didjwk"
	"github.com/tbd54566975/web5-go/jws"
)

func TestCreate(t *testing.T) {
	pfiDID, err := didjwk.Create()
	assert.NoError(t, err)

	_, err = offering.Create(
		offering.NewPayin(
			"USD",
			[]offering.PayinMethod{offering.NewPayinMethod("SQUAREPAY")},
		),
		offering.NewPayout(
			"USDC",
			[]offering.PayoutMethod{offering.NewPayoutMethod("STORED_BALANCE", 20*time.Minute)},
		),
		"1.0",
		offering.From(pfiDID),
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

func TestVerify_SimpleOffering(t *testing.T) {
	portableDIDStr := `{"uri":"did:jwk:eyJrdHkiOiJPS1AiLCJjcnYiOiJFZDI1NTE5IiwieCI6ImpWM3BCeUtmYXkwTHdFQ2lwZ3VUU1MyWExNTDNWQ0UzNnBUVFltZk5ONTQifQ","privateKeys":[{"kty":"OKP","crv":"Ed25519","d":"aIdFbVAIgnqnrH-TDLyZVAEP9QD6vt5C9fhUkPystB-NXekHIp9rLQvAQKKmC5NJLZcswvdUITfqlNNiZ803ng","x":"jV3pByKfay0LwECipguTSS2XLML3VCE36pTTYmfNN54"}],"document":{"@context":["https://www.w3.org/ns/did/v1"],"id":"did:jwk:eyJrdHkiOiJPS1AiLCJjcnYiOiJFZDI1NTE5IiwieCI6ImpWM3BCeUtmYXkwTHdFQ2lwZ3VUU1MyWExNTDNWQ0UzNnBUVFltZk5ONTQifQ","verificationMethod":[{"id":"did:jwk:eyJrdHkiOiJPS1AiLCJjcnYiOiJFZDI1NTE5IiwieCI6ImpWM3BCeUtmYXkwTHdFQ2lwZ3VUU1MyWExNTDNWQ0UzNnBUVFltZk5ONTQifQ#0","type":"JsonWebKey","controller":"did:jwk:eyJrdHkiOiJPS1AiLCJjcnYiOiJFZDI1NTE5IiwieCI6ImpWM3BCeUtmYXkwTHdFQ2lwZ3VUU1MyWExNTDNWQ0UzNnBUVFltZk5ONTQifQ","publicKeyJwk":{"kty":"OKP","crv":"Ed25519","x":"jV3pByKfay0LwECipguTSS2XLML3VCE36pTTYmfNN54"}}],"assertionMethod":["did:jwk:eyJrdHkiOiJPS1AiLCJjcnYiOiJFZDI1NTE5IiwieCI6ImpWM3BCeUtmYXkwTHdFQ2lwZ3VUU1MyWExNTDNWQ0UzNnBUVFltZk5ONTQifQ#0"],"authentication":["did:jwk:eyJrdHkiOiJPS1AiLCJjcnYiOiJFZDI1NTE5IiwieCI6ImpWM3BCeUtmYXkwTHdFQ2lwZ3VUU1MyWExNTDNWQ0UzNnBUVFltZk5ONTQifQ#0"],"capabilityDelegation":["did:jwk:eyJrdHkiOiJPS1AiLCJjcnYiOiJFZDI1NTE5IiwieCI6ImpWM3BCeUtmYXkwTHdFQ2lwZ3VUU1MyWExNTDNWQ0UzNnBUVFltZk5ONTQifQ#0"],"capabilityInvocation":["did:jwk:eyJrdHkiOiJPS1AiLCJjcnYiOiJFZDI1NTE5IiwieCI6ImpWM3BCeUtmYXkwTHdFQ2lwZ3VUU1MyWExNTDNWQ0UzNnBUVFltZk5ONTQifQ#0"]},"metadata":null}`
	var portableDID did.PortableDID
	err := json.Unmarshal([]byte(portableDIDStr), &portableDID)
	assert.NoError(t, err)
	pfiDID, err := did.FromPortableDID(portableDID)
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
		offering.From(pfiDID),
	)
	assert.NoError(t, err)

	// ---
	offeringJSON, err := json.MarshalIndent(offering, "", "    ")
	if err != nil {
		fmt.Println("Error serializing JSON:", err)
		return
	}
	fmt.Println(string(offeringJSON))
}

func Test_TestVector(t *testing.T) {
	offeringStr := `{"metadata":{"from":"did:jwk:eyJrdHkiOiJPS1AiLCJjcnYiOiJFZDI1NTE5IiwieCI6InhfeU1ZZ2RNODhPZ0pPZS1zMFN6aHN0UUwwQ0h4SGdFelRYblA4U3RNZnMifQ","kind":"offering","id":"offering_01j2fkvz7efqjt6k248trnsx6s","createdAt":"2024-07-11T00:42:38Z","updatedAt":"2024-07-11T00:42:38Z","protocol":"1.0"},"data":{"description":"USDC for USD","payoutUnitsPerPayinUnit":"1.0","payin":{"currencyCode":"USD","min":"0.1","max":"1000","methods":[{"kind":"DEBIT_CARD","requiredPaymentDetails":{"$schema":"http://json-schema.org/draft-07/schema#","type":"object","properties":{"cardNumber":{"type":"string","description":"The 16-digit debit card number","minLength":16,"maxLength":16},"expiryDate":{"type":"string","description":"The expiry date of the card in MM/YY format","pattern":"^(0[1-9]|1[0-2])\\/([0-9]{2})$"},"cardHolderName":{"type":"string","description":"Name of the cardholder as it appears on the card"},"cvv":{"type":"string","description":"The 3-digit CVV code","minLength":3,"maxLength":3}},"required":["cardNumber","expiryDate","cardHolderName","cvv"],"additionalProperties":false}}]},"payout":{"currencyCode":"USDC","max":"5000","methods":[{"kind":"STORED_BALANCE","estimatedSettlementTime":1200}]},"requiredClaims":{"id":"foo","name":"kyccredential","purpose":"To verify the identity of the user","input_descriptors":[{"id":"1","name":"KYC Information","purpose":"To verify the identity of the user","constraints":{"fields":[{"path":["$.type[0]"],"filter":{"type":"string","pattern":"KYC"}}]}}]},"cancellation":{"enabled":false}},"signature":"eyJhbGciOiJFZERTQSIsImtpZCI6ImRpZDpqd2s6ZXlKcmRIa2lPaUpQUzFBaUxDSmpjbllpT2lKRlpESTFOVEU1SWl3aWVDSTZJbmhmZVUxWloyUk5PRGhQWjBwUFpTMXpNRk42YUhOMFVVd3dRMGg0U0dkRmVsUllibEE0VTNSTlpuTWlmUSMwIn0..4TtQVGurrHzk4_IJgH7zZmlDzn354M67YVVu-n21IAW52-AyPdz9W13efslj9k5y49zIFjkg76yoHFUfL-yeAg"}`
	var offering offering.Offering
	err := json.Unmarshal([]byte(offeringStr), &offering)
	assert.NoError(t, err)

	err = offering.Verify()
	assert.NoError(t, err)
}
