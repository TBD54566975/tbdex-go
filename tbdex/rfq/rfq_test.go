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
	err = rfq.ValidateAndUnmarshalJSON(bytes, false)
	assert.NoError(t, err)
}

func TestRFQ_UnmarshalJSON_FailsIncompletePrivateData(t *testing.T) {
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
	err = rfq.ValidateAndUnmarshalJSON(bytes, true)
	assert.Error(t, err)
}

func TestRFQ_UnmarshalJSON_VerifiesPresentPrivateData(t *testing.T) {
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
	err = rfq.ValidateAndUnmarshalJSON(bytes, false)
	assert.NoError(t, err)
}

func TestRFQ_Unmarshal_Invalid(t *testing.T) {
	input := []byte(`{"doo": "doo"}`)

	var rfq rfq.RFQ
	err := rfq.ValidateAndUnmarshalJSON(input, false)
	assert.Error(t, err)
}
