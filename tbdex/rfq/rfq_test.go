package rfq_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/TBD54566975/tbdex-go/tbdex/offering"
	"github.com/TBD54566975/tbdex-go/tbdex/rfq"
	"github.com/alecthomas/assert/v2"
	"github.com/tbd54566975/web5-go/dids/did"
	"github.com/tbd54566975/web5-go/dids/didjwk"
	"github.com/tbd54566975/web5-go/jws"
	"github.com/tbd54566975/web5-go/pexv2"
	"go.jetpack.io/typeid"
)

func TestCreate(t *testing.T) {
	pfiDID, err := didjwk.Create()
	walletDID, _ := didjwk.Create()
	assert.NoError(t, err)

	offeringID, err := typeid.WithPrefix(offering.Kind)
	assert.NoError(t, err)

	rfq, err := rfq.Create(
		walletDID,
		pfiDID.URI,
		offeringID.String(),
		rfq.Payin("100", "STORED_BALANCE"),
		rfq.Payout("BANK_ACCOUNT"),
		rfq.ExternalID("test_1234"),
	)

	assert.NoError(t, err)
	assert.Zero(t, rfq.PrivateData)
	assert.Zero(t, rfq.Data.Payin.PaymentDetailsHash)
	assert.Zero(t, rfq.Data.Payin.PaymentDetailsHash)
	assert.Zero(t, rfq.Data.ClaimsHash)
}

func TestCreate_WithPrivate(t *testing.T) {
	pfiDID, _ := didjwk.Create()
	walletDID, _ := didjwk.Create()
	offeringID, _ := typeid.WithPrefix(offering.Kind)

	rfq, err := rfq.Create(
		walletDID,
		pfiDID.URI,
		offeringID.String(),
		rfq.Payin("100", "STORED_BALANCE"),
		rfq.Payout("BANK_ACCOUNT", rfq.PaymentDetails(
			map[string]interface{}{
				"accountNumber": "1234567890123456",
				"routingNumber": "123456789",
			})),
		rfq.Claims([]string{"my_jwt"}),
	)

	assert.NoError(t, err)
	assert.Zero(t, rfq.Data.Payin.PaymentDetailsHash)
	assert.NotZero(t, rfq.Data.Payout.PaymentDetailsHash)
	assert.NotZero(t, rfq.Data.ClaimsHash)
}

func TestSign(t *testing.T) {
	pfiDID, _ := didjwk.Create()
	walletDID, _ := didjwk.Create()
	offeringID, _ := typeid.WithPrefix(offering.Kind)

	r, err := rfq.Create(
		walletDID,
		pfiDID.URI,
		offeringID.String(),
		rfq.Payin("100", "STORED_BALANCE"),
		rfq.Payout("BANK_ACCOUNT"),
	)

	assert.NoError(t, err)
	assert.NotZero(t, r.Signature)
}

func TestUnmarshal(t *testing.T) {
	pfiDID, _ := didjwk.Create()
	walletDID, _ := didjwk.Create()
	offeringID, _ := typeid.WithPrefix(offering.Kind)

	r, _ := rfq.Create(
		walletDID,
		pfiDID.URI,
		offeringID.String(),
		rfq.Payin("100", "STORED_BALANCE"),
		rfq.Payout("BANK_ACCOUNT"),
	)

	bytes, err := json.Marshal(r)
	assert.NoError(t, err)

	var rfq rfq.RFQ
	err = json.Unmarshal(bytes, &rfq)
	assert.NoError(t, err)
}

func TestUnmarshal_Empty(t *testing.T) {
	input := []byte(`{"metadata":{},"data":{},"signature":"eyJhbGciOiJFZERTQSIsImtpZCI6ImRpZDpkaHQ6M3doZnRncGJkamloeDl6ZTl0ZG41NzV6cXptNHF3Y2NldG5mMXliaWlidXphZDdycm15eSMwIn0..ZvoVDuSrqqdXsSXgqB-U26tAU1WqUqqU_KpD1KvdYocIcmTsshjUASEwM_lUz1UnGglqkWeCIrHqrm9NNGDqBw"}`)

	var rfq rfq.RFQ
	_ = json.Unmarshal(input, &rfq)

	assert.Zero(t, rfq.Metadata)
	assert.Zero(t, rfq.Data)
}

func TestUnmarshal_Invalid(t *testing.T) {
	input := []byte(`{"doo": "doo"}`)

	var rfq rfq.RFQ
	err := json.Unmarshal(input, &rfq)
	assert.Error(t, err)
}

func TestScrub_FailsNoPrivateData(t *testing.T) {
	pfiDID, _ := didjwk.Create()
	walletDID, _ := didjwk.Create()
	offeringID, _ := typeid.WithPrefix(offering.Kind)

	r, _ := rfq.Create(
		walletDID,
		pfiDID.URI,
		offeringID.String(),
		rfq.Payin("100", "STORED_BALANCE"),
		rfq.Payout("BANK_ACCOUNT"),
	)

	bytes, err := json.Marshal(r)
	assert.NoError(t, err)

	var rfq rfq.RFQ
	err = json.Unmarshal(bytes, &rfq)
	assert.NoError(t, err)

	err = rfq.Verify()
	assert.NoError(t, err)
	_, _, err = rfq.Scrub()
	assert.Error(t, err)
}

func TestScrub_FailsClaimsHashMismatch(t *testing.T) {
	pfiDID, _ := didjwk.Create()
	walletDID, _ := didjwk.Create()
	offeringID, _ := typeid.WithPrefix(offering.Kind)

	r, _ := rfq.Create(
		walletDID,
		pfiDID.URI,
		offeringID.String(),
		rfq.Payin("100", "STORED_BALANCE"),
		rfq.Payout("BANK_ACCOUNT"),
		rfq.Claims([]string{"my_jwt"}),
	)

	r.PrivateData.Claims = []string{"different_jwt"}

	bytes, err := json.Marshal(r)
	assert.NoError(t, err)

	var rfq rfq.RFQ
	err = json.Unmarshal(bytes, &rfq)
	assert.NoError(t, err)

	err = rfq.Verify()
	assert.NoError(t, err)
	_, _, err = rfq.Scrub()
	assert.Error(t, err)
}

func TestScrub_FailsPayoutHashMismatch(t *testing.T) {
	pfiDID, _ := didjwk.Create()
	walletDID, _ := didjwk.Create()
	offeringID, _ := typeid.WithPrefix(offering.Kind)

	r, _ := rfq.Create(
		walletDID,
		pfiDID.URI,
		offeringID.String(),
		rfq.Payin("100", "STORED_BALANCE"),
		rfq.Payout("BANK_ACCOUNT", rfq.PaymentDetails(
			map[string]any{
				"accountNumber": "1234567890123456",
				"routingNumber": "123456789",
			})),
		rfq.Claims([]string{"my_jwt"}),
	)

	r.PrivateData.Payout.PaymentDetails = map[string]any{
		"accountNumber": "1234567890123456",
		"routingNumber": "new_routing_number",
	}

	bytes, err := json.Marshal(r)
	assert.NoError(t, err)

	var rfq rfq.RFQ
	err = json.Unmarshal(bytes, &rfq)
	assert.NoError(t, err)

	err = rfq.Verify()
	assert.NoError(t, err)
	_, _, err = rfq.Scrub()
	assert.Error(t, err)
}

func TestScrub_FailsPayinHashMismatch(t *testing.T) {
	pfiDID, _ := didjwk.Create()
	walletDID, _ := didjwk.Create()
	offeringID, _ := typeid.WithPrefix(offering.Kind)

	r, _ := rfq.Create(
		walletDID,
		pfiDID.URI,
		offeringID.String(),
		rfq.Payin("100", "BANK_ACCOUNT", rfq.PaymentDetails(
			map[string]any{
				"accountNumber": "1234567890123456",
				"routingNumber": "123456789",
			})),
		rfq.Payout("STORED_BALANCE"),
		rfq.Claims([]string{"my_jwt"}),
	)

	r.PrivateData.Payin.PaymentDetails = map[string]any{
		"accountNumber": "1234567890123456",
		"routingNumber": "new_routing_number",
	}

	bytes, err := json.Marshal(r)
	assert.NoError(t, err)

	var rfq rfq.RFQ
	err = json.Unmarshal(bytes, &rfq)
	assert.NoError(t, err)

	err = rfq.Verify()
	assert.NoError(t, err)
	_, _, err = rfq.Scrub()
	assert.Error(t, err)
}

func TestScrub_ClaimsPrivateData(t *testing.T) {
	pfiDID, _ := didjwk.Create()
	walletDID, _ := didjwk.Create()
	offeringID, _ := typeid.WithPrefix(offering.Kind)

	r, _ := rfq.Create(
		walletDID,
		pfiDID.URI,
		offeringID.String(),
		rfq.Payin("100", "STORED_BALANCE"),
		rfq.Payout("BANK_ACCOUNT", rfq.PaymentDetails(
			map[string]any{
				"accountNumber": "1234567890123456",
				"routingNumber": "123456789",
			})),
		rfq.Claims([]string{"my_jwt"}),
	)

	bytes, err := json.Marshal(r)
	assert.NoError(t, err)

	var rfq rfq.RFQ
	err = json.Unmarshal(bytes, &rfq)
	assert.NoError(t, err)

	err = rfq.Verify()
	assert.NoError(t, err)
	_, _, err = rfq.Scrub()
	assert.NoError(t, err)
}

func TestScrub_FailsMissingDataForClaimsHash(t *testing.T) {
	pfiDID, _ := didjwk.Create()
	walletDID, _ := didjwk.Create()
	offeringID, _ := typeid.WithPrefix(offering.Kind)

	r, _ := rfq.Create(
		walletDID,
		pfiDID.URI,
		offeringID.String(),
		rfq.Payin("100", "STORED_BALANCE"),
		rfq.Payout("BANK_ACCOUNT"),
		rfq.Claims([]string{"my_jwt"}),
	)

	r.PrivateData.Claims = nil

	bytes, err := json.Marshal(r)
	assert.NoError(t, err)

	var rfq rfq.RFQ
	err = json.Unmarshal(bytes, &rfq)
	assert.NoError(t, err)

	err = rfq.Verify()
	assert.NoError(t, err)
	_, _, err = rfq.Scrub()
	assert.Error(t, err)
}

func TestScrub_FailsMissingDataForPayoutHash(t *testing.T) {
	pfiDID, _ := didjwk.Create()
	walletDID, _ := didjwk.Create()
	offeringID, _ := typeid.WithPrefix(offering.Kind)

	r, _ := rfq.Create(
		walletDID,
		pfiDID.URI,
		offeringID.String(),
		rfq.Payin("100", "STORED_BALANCE"),
		rfq.Payout("BANK_ACCOUNT", rfq.PaymentDetails(
			map[string]any{
				"accountNumber": "1234567890123456",
				"routingNumber": "123456789",
			})),
	)

	r.PrivateData.Payout.PaymentDetails = nil

	bytes, err := json.Marshal(r)
	assert.NoError(t, err)

	var rfq rfq.RFQ
	err = json.Unmarshal(bytes, &rfq)
	assert.NoError(t, err)

	err = rfq.Verify()
	assert.NoError(t, err)
	_, _, err = rfq.Scrub()
	assert.Error(t, err)
}

func TestVerify_FailsBadSignature(t *testing.T) {
	pfiDID, _ := didjwk.Create()
	walletDID, _ := didjwk.Create()
	offeringID, _ := typeid.WithPrefix(offering.Kind)

	r, _ := rfq.Create(
		walletDID,
		pfiDID.URI,
		offeringID.String(),
		rfq.Payin("100", "STORED_BALANCE"),
		rfq.Payout("BANK_ACCOUNT", rfq.PaymentDetails(
			map[string]any{
				"accountNumber": "1234567890123456",
				"routingNumber": "123456789",
			})),
		rfq.Claims([]string{"my_jwt"}),
	)

	r.Signature = "bad signature"

	bytes, err := json.Marshal(r)
	assert.NoError(t, err)

	var rfq rfq.RFQ
	err = json.Unmarshal(bytes, &rfq)
	assert.NoError(t, err)

	err = rfq.Verify()
	assert.Error(t, err)
}

func TestVerify_InvalidSignature(t *testing.T) {
	pfiDID, _ := didjwk.Create()
	walletDID, _ := didjwk.Create()
	offeringID, _ := typeid.WithPrefix(offering.Kind)

	r, _ := rfq.Create(
		walletDID,
		pfiDID.URI,
		offeringID.String(),
		rfq.Payin("100", "STORED_BALANCE"),
		rfq.Payout("BANK_ACCOUNT"),
	)

	r.Signature = "Invalid"

	bytes, err := json.Marshal(r)
	assert.NoError(t, err)

	var rfq rfq.RFQ
	err = json.Unmarshal(bytes, &rfq)
	assert.NoError(t, err)

	err = rfq.Verify()
	assert.Error(t, err)
}

func TestVerify_SignedWithWrongDID(t *testing.T) {
	pfiDID, _ := didjwk.Create()
	walletDID, _ := didjwk.Create()
	wrongDID, _ := didjwk.Create()
	offeringID, _ := typeid.WithPrefix(offering.Kind)

	r, _ := rfq.Create(
		walletDID,
		pfiDID.URI,
		offeringID.String(),
		rfq.Payin("100", "STORED_BALANCE"),
		rfq.Payout("BANK_ACCOUNT"),
	)

	toSign, err := r.Digest()
	assert.NoError(t, err)

	wrongSignature, err := jws.Sign(toSign, wrongDID, jws.DetachedPayload(true))
	assert.NoError(t, err)

	r.Signature = wrongSignature

	bytes, err := json.Marshal(r)
	assert.NoError(t, err)

	var rfq rfq.RFQ
	err = json.Unmarshal(bytes, &rfq)
	assert.NoError(t, err)

	err = rfq.Verify()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "does not match message metadata from")
}

func TestVerifyOfferingRequirements_Pass(t *testing.T) {
	pfiDID, _ := didjwk.Create()
	walletDID, _ := didjwk.Create()

	offering, err := getSimpleOffering(pfiDID)
	assert.NoError(t, err)

	r, _ := rfq.Create(
		walletDID,
		pfiDID.URI,
		offering.Metadata.ID,
		rfq.Payin("100", "SQUAREPAY"),
		rfq.Payout("STORED_BALANCE"),
	)

	err = r.VerifyOfferingRequirements(offering)
	assert.NoError(t, err)
}

func TestVerifyOfferingRequirements_OfferingIdMismatch(t *testing.T) {
	pfiDID, _ := didjwk.Create()
	walletDID, _ := didjwk.Create()

	offering, err := getSimpleOffering(pfiDID)
	assert.NoError(t, err)

	r, _ := rfq.Create(
		walletDID,
		pfiDID.URI,
		"wrong_offering_id",
		rfq.Payin("100", "SQUAREPAY"),
		rfq.Payout("STORED_BALANCE"),
	)

	err = r.VerifyOfferingRequirements(offering)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "offering ID mismatch")
}

func TestVerifyOfferingRequirements_PayinMoreThanMax(t *testing.T) {
	pfiDID, _ := didjwk.Create()
	walletDID, _ := didjwk.Create()

	offering, err := getSimpleOffering(pfiDID)
	assert.NoError(t, err)

	r, _ := rfq.Create(
		walletDID,
		pfiDID.URI,
		offering.Metadata.ID,
		rfq.Payin("99999", "SQUAREPAY"),
		rfq.Payout("STORED_BALANCE"),
	)

	err = r.VerifyOfferingRequirements(offering)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "payin amount exceeds maximum")
}

func TestVerifyOfferingRequirements_PayinLessThanMin(t *testing.T) {
	pfiDID, _ := didjwk.Create()
	walletDID, _ := didjwk.Create()

	offering, err := getSimpleOffering(pfiDID)
	assert.NoError(t, err)

	r, _ := rfq.Create(
		walletDID,
		pfiDID.URI,
		offering.Metadata.ID,
		rfq.Payin("1", "SQUAREPAY"),
		rfq.Payout("STORED_BALANCE"),
	)

	err = r.VerifyOfferingRequirements(offering)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "payin amount is below minimum")
}

func TestVerifyOfferingRequirements_VerifyPayinMethodFail_PayinMethodKindNotSupported(t *testing.T) {
	pfiDID, _ := didjwk.Create()
	walletDID, _ := didjwk.Create()

	offering, err := getSimpleOffering(pfiDID)
	assert.NoError(t, err)

	r, _ := rfq.Create(
		walletDID,
		pfiDID.URI,
		offering.Metadata.ID,
		rfq.Payin("100","AFTERPAY"),
		rfq.Payout("STORED_BALANCE"),
	)

	err = r.VerifyOfferingRequirements(offering)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "offering does not support rfq's payinMethod kind")
}

func TestVerifyOfferingRequirements_VerifyPayinMethodFail_OfferingRequiredPaymentDetailsNil(t *testing.T) {
	pfiDID, _ := didjwk.Create()
	walletDID, _ := didjwk.Create()

	offering, err := getSimpleOffering(pfiDID)
	assert.NoError(t, err)

	r, _ := rfq.Create(
		walletDID,
		pfiDID.URI,
		offering.Metadata.ID,
		rfq.Payin(
			"100",
			"SQUAREPAY",
			rfq.PaymentDetails(map[string]any{"accountNumber": "1234567890123456"}),
		),
		rfq.Payout("STORED_BALANCE"),
	)

	err = r.VerifyOfferingRequirements(offering)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "paymentDetails must be omitted when offering requiredPaymentDetails is omitted")
}

func TestVerifyOfferingRequirements_VerifyPayinMethodFail_PayinPrivateDataNil(t *testing.T) {
	pfiDID, _ := didjwk.Create()
	walletDID, _ := didjwk.Create()

	offering, err := offering.Create(
		offering.NewPayin(
			"MXN",
			[]offering.PayinMethod{
				offering.NewPayinMethod(
					"SPEI",
					offering.RequiredDetails(`{
						"$schema": "http://json-schema.org/draft-07/schema#",
						"additionalProperties": false,
						"properties": {
							"clabe": {
								"type": "string"
							}
						},
						"required": ["clabe"]
					}`),
				),
			},
		),
		offering.NewPayout(
			"USDC",
			[]offering.PayoutMethod{offering.NewPayoutMethod("STORED_BALANCE", 20*time.Minute)},
		),
		"16.0",
		offering.From(pfiDID),
	)
	assert.NoError(t, err)

	r, _ := rfq.Create(
		walletDID,
		pfiDID.URI,
		offering.Metadata.ID,
		rfq.Payin("100", "SPEI"),
		rfq.Payout("STORED_BALANCE"),
	)

	err = r.VerifyOfferingRequirements(offering)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "offering requiredPaymentDetails is present but rfq private data is omitted")
}

// todo not sure if we need this test. see rfq.go#301
func TestVerifyOfferingRequirements_VerifyPayinMethodFail_RequiredPrivateDataNil(t *testing.T) {
	pfiDID, _ := didjwk.Create()
	walletDID, _ := didjwk.Create()

	offering, err := offering.Create(
		offering.NewPayin(
			"USD",
			[]offering.PayinMethod{
				offering.NewPayinMethod(
					"SPEI",
					offering.RequiredDetails(`{
						"$schema": "http://json-schema.org/draft-07/schema#",
						"additionalProperties": false,
						"properties": {
							"clabe": {
								"type": "string"
							}
						},
						"required": ["clabe"]
					}`),
				),
			},
		),
		offering.NewPayout(
			"USDC",
			[]offering.PayoutMethod{offering.NewPayoutMethod("STORED_BALANCE", 20*time.Minute)},
		),
		"1.0",
		offering.From(pfiDID),
	)
	assert.NoError(t, err)

	r, _ := rfq.Create(
		walletDID,
		pfiDID.URI,
		offering.Metadata.ID,
		rfq.Payin(
			"100",
			"SPEI",
		),
		rfq.Payout("STORED_BALANCE", rfq.PaymentDetails(map[string]any{"accountNumber": "1234567890123456"})),
	)

	err = r.VerifyOfferingRequirements(offering)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "rfq payin paymentDetails are missing")
}

func TestVerifyOfferingRequirements_VerifyPayinMethodFail_SchemaValidationFails(t *testing.T) {
	pfiDID, _ := didjwk.Create()
	walletDID, _ := didjwk.Create()

	offering, err := offering.Create(
		offering.NewPayin(
			"USD",
			[]offering.PayinMethod{
				offering.NewPayinMethod(
					"SPEI",
					offering.RequiredDetails(`{
						"$schema": "http://json-schema.org/draft-07/schema#",
						"additionalProperties": false,
						"properties": {
							"clabe": {
								"type": "string"
							}
						},
						"required": ["clabe"]
					}`),
				),
			},
		),
		offering.NewPayout(
			"USDC",
			[]offering.PayoutMethod{offering.NewPayoutMethod("STORED_BALANCE", 20*time.Minute)},
		),
		"1.0",
		offering.From(pfiDID),
	)
	assert.NoError(t, err)

	r, _ := rfq.Create(
		walletDID,
		pfiDID.URI,
		offering.Metadata.ID,
		rfq.Payin(
			"100",
			"SPEI",
			rfq.PaymentDetails(map[string]any{"NotClabe": 1234567890123456}),
		),
		rfq.Payout("STORED_BALANCE"),
	)

	err = r.VerifyOfferingRequirements(offering)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to validate payinMethod's paymentDetails")
}

func TestVerifyOfferingRequirements_VerifyPayoutMethodFail_SchemaValidationFails(t *testing.T) {
	pfiDID, _ := didjwk.Create()
	walletDID, _ := didjwk.Create()

	offering, err := offering.Create(
		offering.NewPayin(
			"USD",
			[]offering.PayinMethod{
				offering.NewPayinMethod(
					"SPEI",
				),
			},
		),
		offering.NewPayout(
			"USDC",
			[]offering.PayoutMethod{
				offering.NewPayoutMethod(
					"STORED_BALANCE",
					20*time.Minute,
					offering.RequiredDetails(`{
						"$schema": "http://json-schema.org/draft-07/schema#",
						"additionalProperties": false,
						"properties": {
							"name": {
								"type": "string"
							}
						},
						"required": ["name"]
					}`),
				)},
		),
		"1.0",
		offering.From(pfiDID),
	)
	assert.NoError(t, err)

	r, _ := rfq.Create(
		walletDID,
		pfiDID.URI,
		offering.Metadata.ID,
		rfq.Payin(
			"100",
			"SPEI",
		),
		rfq.Payout(
			"STORED_BALANCE",
			rfq.PaymentDetails(map[string]any{"favoriteColor": "purple"}),
		),
	)

	err = r.VerifyOfferingRequirements(offering)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to validate payoutMethod's paymentDetails")
}

func TestVerifyOfferingRequirements_VerifyClaimsPass(t *testing.T) {
	pfiDID, _ := didjwk.Create()
	walletDID, _ := didjwk.Create()

	pd := getPresentationDefinition()

	vcJwt := "eyJhbGciOiJFZERTQSIsImtpZCI6ImRpZDpqd2s6ZXlKcmRIa2lPaUpQUzFBaUxDSmpjbllpT2lKRlpESTFOVEU1SWl3aWVDSTZJa013VlhrNGN6SnVla1JhTjFsdmRXeFFUM010YjBabFEwWkxaVzAzY1hST1FWVTNiVzB0TjNaaFkxRWlmUSMwIiwidHlwIjoiSldUIn0.eyJpc3MiOiJkaWQ6andrOmV5SnJkSGtpT2lKUFMxQWlMQ0pqY25ZaU9pSkZaREkxTlRFNUlpd2llQ0k2SWtNd1ZYazRjekp1ZWtSYU4xbHZkV3hRVDNNdGIwWmxRMFpMWlcwM2NYUk9RVlUzYlcwdE4zWmhZMUVpZlEiLCJqdGkiOiJ1cm46dmM6dXVpZDo3MjllMTc2ZS1mYjNlLTQyOTktOWI3Yi02MGQzNmVkMzQxNmUiLCJuYmYiOjE3MTQxNjE1NjEsInN1YiI6ImRpZDpqd2s6ZXlKcmRIa2lPaUpQUzFBaUxDSmpjbllpT2lKRlpESTFOVEU1SWl3aWVDSTZJa013VlhrNGN6SnVla1JhTjFsdmRXeFFUM010YjBabFEwWkxaVzAzY1hST1FWVTNiVzB0TjNaaFkxRWlmUSIsInZjIjp7IkBjb250ZXh0IjpbImh0dHBzOi8vd3d3LnczLm9yZy8yMDE4L2NyZWRlbnRpYWxzL3YxIl0sInR5cGUiOlsiVmVyaWZpYWJsZUNyZWRlbnRpYWwiLCJTdHJlZXRDcmVkZW50aWFsIl0sImlzc3VlciI6ImRpZDpqd2s6ZXlKcmRIa2lPaUpQUzFBaUxDSmpjbllpT2lKRlpESTFOVEU1SWl3aWVDSTZJa013VlhrNGN6SnVla1JhTjFsdmRXeFFUM010YjBabFEwWkxaVzAzY1hST1FWVTNiVzB0TjNaaFkxRWlmUSIsImNyZWRlbnRpYWxTdWJqZWN0Ijp7ImlkIjoiZGlkOmp3azpleUpyZEhraU9pSlBTMUFpTENKamNuWWlPaUpGWkRJMU5URTVJaXdpZUNJNklrTXdWWGs0Y3pKdWVrUmFOMWx2ZFd4UVQzTXRiMFpsUTBaTFpXMDNjWFJPUVZVM2JXMHROM1poWTFFaWZRIiwibmFtZSI6IlNhdG9zaGkgVGFjb21vdG8ifSwiaWQiOiJ1cm46dmM6dXVpZDo3MjllMTc2ZS1mYjNlLTQyOTktOWI3Yi02MGQzNmVkMzQxNmUiLCJpc3N1YW5jZURhdGUiOiIyMDI0LTA0LTI2VDE5OjU5OjIxWiJ9fQ.MoPYXkSASXEgySIc59HnSN8576cu5q8QC5tCG3PKr3j-glvZNa12j_P563FUVzx7PeFD3QkJne1RYBDOj3OcBw"

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
		offering.RequiredClaims(pd),
	)
	assert.NoError(t, err)

	r, _ := rfq.Create(
		walletDID,
		pfiDID.URI,
		offering.Metadata.ID,
		rfq.Payin("100", "SQUAREPAY"),
		rfq.Payout("STORED_BALANCE"),
		rfq.Claims([]string{vcJwt}),
	)

	err = r.VerifyOfferingRequirements(offering)
	assert.NoError(t, err)
}

func TestVerifyOfferingRequirements_VerifyClaimsFail(t *testing.T) {
	pfiDID, _ := didjwk.Create()
	walletDID, _ := didjwk.Create()

	pd := getPresentationDefinition()

	vcJwt := "eyJ0eXAiOiJKV1QiLCJhbGciOiJFZERTQSIsImtpZCI6ImRpZDprZXk6ejZNa2ZIQk5tYzRjQ2hOdER2YXIxOFRaaWdmWHQ3UDY1RkJHd3FRVHR4UW9RUG5HI3o2TWtmSEJObWM0Y0NoTnREdmFyMThUWmlnZlh0N1A2NUZCR3dxUVR0eFFvUVBuRyJ9.eyJpc3MiOiJkaWQ6a2V5Ono2TWtmSEJObWM0Y0NoTnREdmFyMThUWmlnZlh0N1A2NUZCR3dxUVR0eFFvUVBuRyIsInN1YiI6ImRpZDprZXk6ejZNa2ZIQk5tYzRjQ2hOdER2YXIxOFRaaWdmWHQ3UDY1RkJHd3FRVHR4UW9RUG5HIiwidmMiOnsiQGNvbnRleHQiOlsiaHR0cHM6Ly93d3cudzMub3JnLzIwMTgvY3JlZGVudGlhbHMvdjEiXSwidHlwZSI6WyJWZXJpZmlhYmxlQ3JlZGVudGlhbCIsIlN0cmVldENyZWQiXSwiaWQiOiJ1cm46dXVpZDoxM2Q1YTg3YS1kY2Y1LTRmYjktOWUyOS0wZTYyZTI0YzQ0ODYiLCJpc3N1ZXIiOiJkaWQ6a2V5Ono2TWtmSEJObWM0Y0NoTnREdmFyMThUWmlnZlh0N1A2NUZCR3dxUVR0eFFvUVBuRyIsImlzc3VhbmNlRGF0ZSI6IjIwMjMtMTItMDdUMTc6MTk6MTNaIiwiY3JlZGVudGlhbFN1YmplY3QiOnsiaWQiOiJkaWQ6a2V5Ono2TWtmSEJObWM0Y0NoTnREdmFyMThUWmlnZlh0N1A2NUZCR3dxUVR0eFFvUVBuRyIsIm90aGVydGhpbmciOiJvdGhlcnN0dWZmIn19fQ.FVvL3z8LHJXm7lGX2bGFvH_U-bTyoheRbLzE7zIk_P1BKwRYeW4sbYNzsovFX59twXrnpF-hHkqVVsejSljxDw"
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
		offering.RequiredClaims(pd),
	)
	assert.NoError(t, err)

	r, _ := rfq.Create(
		walletDID,
		pfiDID.URI,
		offering.Metadata.ID,
		rfq.Payin("100", "SQUAREPAY"),
		rfq.Payout("STORED_BALANCE"),
		rfq.Claims([]string{vcJwt}),
	)

	err = r.VerifyOfferingRequirements(offering)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "claims do not fulfill the offering's requirements")
}

func getPresentationDefinition() pexv2.PresentationDefinition {
	return pexv2.PresentationDefinition{
		ID: "test_pd",
		InputDescriptors: []pexv2.InputDescriptor{
			{
				ID: "test_input_descriptor",
				Constraints: pexv2.Constraints{
					Fields: []pexv2.Field{
						{
							Path: []string{"$.vc.credentialSubject.name"},
							Filter: &pexv2.Filter{
								Const: "Satoshi Tacomoto",
							},
						},
					},
				},
			},
		},
	}
}

func getSimpleOffering(pfiDID did.BearerDID) (offering.Offering, error) {
	offering, err := offering.Create(
		offering.NewPayin(
			"USD",
			[]offering.PayinMethod{offering.NewPayinMethod("SQUAREPAY")},
			offering.Min("5"),
			offering.Max("100"),
		),
		offering.NewPayout(
			"USDC",
			[]offering.PayoutMethod{offering.NewPayoutMethod("STORED_BALANCE", 20*time.Minute)},
		),
		"1.0",
		offering.From(pfiDID),
	)
	return offering, err
}
