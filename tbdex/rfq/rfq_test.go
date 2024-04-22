package rfq_test

import (
	"encoding/json"
	"testing"

	"github.com/TBD54566975/tbdex-go/tbdex/offering"
	"github.com/TBD54566975/tbdex-go/tbdex/rfq"
	"github.com/alecthomas/assert/v2"
	"github.com/tbd54566975/web5-go/dids/didjwk"
	"go.jetpack.io/typeid"
)

func TestCreateRFQ(t *testing.T) {
	pfiDID, err := didjwk.Create()
	assert.NoError(t, err)

	walletDID, err := didjwk.Create()
	assert.NoError(t, err)

	offeringID, err := typeid.WithPrefix(offering.Kind)
	assert.NoError(t, err)

	rfq, err := rfq.Create(
		walletDID.URI,
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

func TestCreateRFQ_WithPrivate(t *testing.T) {
	pfiDID, _ := didjwk.Create()
	walletDID, _ := didjwk.Create()
	offeringID, _ := typeid.WithPrefix(offering.Kind)

	rfq, err := rfq.Create(
		walletDID.URI,
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

func TestRFQ_Sign(t *testing.T) {
	pfiDID, _ := didjwk.Create()

	walletDID, _ := didjwk.Create()
	offeringID, _ := typeid.WithPrefix(offering.Kind)

	r, _ := rfq.Create(
		walletDID.URI,
		pfiDID.URI,
		offeringID.String(),
		rfq.Payin("100", "STORED_BALANCE"),
		rfq.Payout("BANK_ACCOUNT"),
	)

	err := r.Sign(walletDID)
	assert.NoError(t, err)
}

func TestRFQ_UnmarshalJSON(t *testing.T) {
	pfiDID, _ := didjwk.Create()
	walletDID, _ := didjwk.Create()
	offeringID, _ := typeid.WithPrefix(offering.Kind)

	r, _ := rfq.Create(
		walletDID.URI,
		pfiDID.URI,
		offeringID.String(),
		rfq.Payin("100", "STORED_BALANCE"),
		rfq.Payout("BANK_ACCOUNT"),
	)

	_ = r.Sign(walletDID)

	bytes, err := json.Marshal(r)
	assert.NoError(t, err)

	var rfq rfq.RFQ
	err = rfq.UnmarshalJSON(bytes)
	assert.NoError(t, err)
}

func TestRFQ_Unmarshal_Invalid(t *testing.T) {
	input := []byte(`{"doo": "doo"}`)

	var rfq rfq.RFQ
	err := rfq.UnmarshalJSON(input)
	assert.Error(t, err)
}

func TestRFQ_Verify_NoPrivateDataStrict(t *testing.T) {
	pfiDID, _ := didjwk.Create()
	walletDID, _ := didjwk.Create()
	offeringID, _ := typeid.WithPrefix(offering.Kind)

	r, _ := rfq.Create(
		walletDID.URI,
		pfiDID.URI,
		offeringID.String(),
		rfq.Payin("100", "STORED_BALANCE"),
		rfq.Payout("BANK_ACCOUNT"),
	)

	_ = r.Sign(walletDID)

	bytes, err := json.Marshal(r)
	assert.NoError(t, err)

	var rfq rfq.RFQ
	err = rfq.UnmarshalJSON(bytes)
	assert.NoError(t, err)

	err = rfq.Verify(true)
	assert.NoError(t, err)
}

func TestRFQ_Verify_NoPrivateDataNotStrict(t *testing.T) {
	pfiDID, _ := didjwk.Create()
	walletDID, _ := didjwk.Create()
	offeringID, _ := typeid.WithPrefix(offering.Kind)

	r, _ := rfq.Create(
		walletDID.URI,
		pfiDID.URI,
		offeringID.String(),
		rfq.Payin("100", "STORED_BALANCE"),
		rfq.Payout("BANK_ACCOUNT"),
	)

	_ = r.Sign(walletDID)

	bytes, err := json.Marshal(r)
	assert.NoError(t, err)

	var rfq rfq.RFQ
	err = rfq.UnmarshalJSON(bytes)
	assert.NoError(t, err)

	err = rfq.Verify(false)
	assert.NoError(t, err)
}

func TestRFQ_Verify_FailsClaimsHashMismatch(t *testing.T) {
	pfiDID, _ := didjwk.Create()
	walletDID, _ := didjwk.Create()
	offeringID, _ := typeid.WithPrefix(offering.Kind)

	r, _ := rfq.Create(
		walletDID.URI,
		pfiDID.URI,
		offeringID.String(),
		rfq.Payin("100", "STORED_BALANCE"),
		rfq.Payout("BANK_ACCOUNT"),
		rfq.Claims([]string{"my_jwt"}),
	)

	_ = r.Sign(walletDID)
	r.PrivateData.Claims = []string{"different_jwt"}

	bytes, err := json.Marshal(r)
	assert.NoError(t, err)

	var rfq rfq.RFQ
	err = rfq.UnmarshalJSON(bytes)
	assert.NoError(t, err)

	err = rfq.Verify(false)
	assert.Error(t, err)
}

func TestRFQ_Verify_FailsPayoutHashMismatch(t *testing.T) {
	pfiDID, _ := didjwk.Create()
	walletDID, _ := didjwk.Create()
	offeringID, _ := typeid.WithPrefix(offering.Kind)

	r, _ := rfq.Create(
		walletDID.URI,
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

	_ = r.Sign(walletDID)
	r.PrivateData.Payout.PaymentDetails = map[string]any{
		"accountNumber": "1234567890123456",
		"routingNumber": "new_routing_number",
	}

	bytes, err := json.Marshal(r)
	assert.NoError(t, err)

	var rfq rfq.RFQ
	err = rfq.UnmarshalJSON(bytes)
	assert.NoError(t, err)

	err = rfq.Verify(false)
	assert.Error(t, err)
}

func TestRFQ_Verify_FailsPayinHashMismatch(t *testing.T) {
	pfiDID, _ := didjwk.Create()
	walletDID, _ := didjwk.Create()
	offeringID, _ := typeid.WithPrefix(offering.Kind)

	r, _ := rfq.Create(
		walletDID.URI,
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

	_ = r.Sign(walletDID)
	r.PrivateData.Payin.PaymentDetails = map[string]any{
		"accountNumber": "1234567890123456",
		"routingNumber": "new_routing_number",
	}

	bytes, err := json.Marshal(r)
	assert.NoError(t, err)

	var rfq rfq.RFQ
	err = rfq.UnmarshalJSON(bytes)
	assert.NoError(t, err)

	err = rfq.Verify(false)
	assert.Error(t, err)
}

func TestRFQ_Verify_ClaimsPrivateDataStrict(t *testing.T) {
	pfiDID, _ := didjwk.Create()
	walletDID, _ := didjwk.Create()
	offeringID, _ := typeid.WithPrefix(offering.Kind)

	r, _ := rfq.Create(
		walletDID.URI,
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

	_ = r.Sign(walletDID)

	bytes, err := json.Marshal(r)
	assert.NoError(t, err)

	var rfq rfq.RFQ
	err = rfq.UnmarshalJSON(bytes)
	assert.NoError(t, err)

	err = rfq.Verify(true)
	assert.NoError(t, err)
}

func TestRFQ_Verify_FailsMissingDataForClaimsHashStrict(t *testing.T) {
	pfiDID, _ := didjwk.Create()
	walletDID, _ := didjwk.Create()
	offeringID, _ := typeid.WithPrefix(offering.Kind)

	r, _ := rfq.Create(
		walletDID.URI,
		pfiDID.URI,
		offeringID.String(),
		rfq.Payin("100", "STORED_BALANCE"),
		rfq.Payout("BANK_ACCOUNT"),
		rfq.Claims([]string{"my_jwt"}),
	)

	_ = r.Sign(walletDID)
	r.PrivateData.Claims = nil

	bytes, err := json.Marshal(r)
	assert.NoError(t, err)

	var rfq rfq.RFQ
	err = rfq.UnmarshalJSON(bytes)
	assert.NoError(t, err)

	err = rfq.Verify(true)
	assert.Error(t, err)
}

func TestRFQ_Verify_PassesMissingDataForClaimsHashNotStrict(t *testing.T) {
	pfiDID, _ := didjwk.Create()
	walletDID, _ := didjwk.Create()
	offeringID, _ := typeid.WithPrefix(offering.Kind)

	r, _ := rfq.Create(
		walletDID.URI,
		pfiDID.URI,
		offeringID.String(),
		rfq.Payin("100", "STORED_BALANCE"),
		rfq.Payout("BANK_ACCOUNT"),
		rfq.Claims([]string{"my_jwt"}),
	)

	_ = r.Sign(walletDID)
	r.PrivateData.Claims = nil

	bytes, err := json.Marshal(r)
	assert.NoError(t, err)

	var rfq rfq.RFQ
	err = rfq.UnmarshalJSON(bytes)
	assert.NoError(t, err)

	err = rfq.Verify(false)
	assert.NoError(t, err)
}

func TestRFQ_Verify_FailsMissingDataForPayoutHashStrict(t *testing.T) {
	pfiDID, _ := didjwk.Create()
	walletDID, _ := didjwk.Create()
	offeringID, _ := typeid.WithPrefix(offering.Kind)

	r, _ := rfq.Create(
		walletDID.URI,
		pfiDID.URI,
		offeringID.String(),
		rfq.Payin("100", "STORED_BALANCE"),
		rfq.Payout("BANK_ACCOUNT", rfq.PaymentDetails(
			map[string]any{
				"accountNumber": "1234567890123456",
				"routingNumber": "123456789",
			})),
	)

	_ = r.Sign(walletDID)
	r.PrivateData.Payout.PaymentDetails = nil

	bytes, err := json.Marshal(r)
	assert.NoError(t, err)

	var rfq rfq.RFQ
	err = rfq.UnmarshalJSON(bytes)
	assert.NoError(t, err)

	err = rfq.Verify(true)
	assert.Error(t, err)
}

func TestRFQ_Verify_PassesMissingDataForPayoutHashNotStrict(t *testing.T) {
	pfiDID, _ := didjwk.Create()
	walletDID, _ := didjwk.Create()
	offeringID, _ := typeid.WithPrefix(offering.Kind)

	r, _ := rfq.Create(
		walletDID.URI,
		pfiDID.URI,
		offeringID.String(),
		rfq.Payin("100", "STORED_BALANCE"),
		rfq.Payout("BANK_ACCOUNT", rfq.PaymentDetails(
			map[string]any{
				"accountNumber": "1234567890123456",
				"routingNumber": "123456789",
			})),
	)

	_ = r.Sign(walletDID)
	r.PrivateData.Payout.PaymentDetails = nil

	bytes, err := json.Marshal(r)
	assert.NoError(t, err)

	var rfq rfq.RFQ
	err = rfq.UnmarshalJSON(bytes)
	assert.NoError(t, err)

	err = rfq.Verify(false)
	assert.NoError(t, err)
}

func TestRFQ_Verify_FailsBadSignature(t *testing.T) {
	pfiDID, _ := didjwk.Create()
	walletDID, _ := didjwk.Create()
	offeringID, _ := typeid.WithPrefix(offering.Kind)

	r, _ := rfq.Create(
		walletDID.URI,
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
	err = rfq.UnmarshalJSON(bytes)
	assert.NoError(t, err)

	err = rfq.Verify(false)
	assert.Error(t, err)
}


func TestVerify_InvalidSignature(t *testing.T) {
	pfiDID, _ := didjwk.Create()
	walletDID, _ := didjwk.Create()
	offeringID, _ := typeid.WithPrefix(offering.Kind)

	r, _ := rfq.Create(
		walletDID.URI,
		pfiDID.URI,
		offeringID.String(),
		rfq.Payin("100", "STORED_BALANCE"),
		rfq.Payout("BANK_ACCOUNT"),
	)

	_ = r.Sign(walletDID)

	r.Signature = "Invalid"

	bytes, err := json.Marshal(r)
	assert.NoError(t, err)

	var rfq rfq.RFQ
	err = rfq.UnmarshalJSON(bytes)
	assert.NoError(t, err)

	err = rfq.Verify(true)
	assert.Error(t, err)
}