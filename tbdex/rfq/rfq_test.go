package rfq_test

import (
	"encoding/json"
	"testing"

	"github.com/TBD54566975/tbdex-go/tbdex/offering"
	"github.com/TBD54566975/tbdex-go/tbdex/rfq"
	"github.com/alecthomas/assert/v2"
	"github.com/tbd54566975/web5-go/dids/didjwk"
	"github.com/tbd54566975/web5-go/jws"
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
	_ = json.Unmarshal([]byte(input), &rfq)

	assert.Zero(t, rfq.Metadata)
	assert.Zero(t, rfq.Data)
}

func TestUnmarshal_Invalid(t *testing.T) {
	input := []byte(`{"doo": "doo"}`)

	var rfq rfq.RFQ
	err := json.Unmarshal(input, &rfq)
	assert.Error(t, err)
}

func TestVerify_NoPrivateDataStrict(t *testing.T) {
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

	err = rfq.Verify(true)
	assert.NoError(t, err)
}

func TestVerify_NoPrivateDataNotStrict(t *testing.T) {
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

	err = rfq.Verify(false)
	assert.NoError(t, err)
}

func TestVerify_FailsClaimsHashMismatch(t *testing.T) {
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

	err = rfq.Verify(false)
	assert.Error(t, err)
}

func TestVerify_FailsPayoutHashMismatch(t *testing.T) {
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

	err = rfq.Verify(false)
	assert.Error(t, err)
}

func TestVerify_FailsPayinHashMismatch(t *testing.T) {
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

	err = rfq.Verify(false)
	assert.Error(t, err)
}

func TestVerify_ClaimsPrivateDataStrict(t *testing.T) {
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

	err = rfq.Verify(true)
	assert.NoError(t, err)
}

func TestVerify_FailsMissingDataForClaimsHashStrict(t *testing.T) {
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

	err = rfq.Verify(true)
	assert.Error(t, err)
}

func TestVerify_PassesMissingDataForClaimsHashNotStrict(t *testing.T) {
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

	err = rfq.Verify(false)
	assert.NoError(t, err)
}

func TestVerify_FailsMissingDataForPayoutHashStrict(t *testing.T) {
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

	err = rfq.Verify(true)
	assert.Error(t, err)
}

func TestVerify_PassesMissingDataForPayoutHashNotStrict(t *testing.T) {
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

	err = rfq.Verify(false)
	assert.NoError(t, err)
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

	err = rfq.Verify(false)
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

	err = rfq.Verify(true)
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

	var RFQ rfq.RFQ
	err = json.Unmarshal(bytes, &RFQ)
	assert.NoError(t, err)

	err = RFQ.Verify(true)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "does not match message metadata from")
}
