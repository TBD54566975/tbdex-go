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

	_, err = rfq.CreateRFQ(
		walletDID.URI,
		pfiDID.URI,
		offeringID.String(),
		rfq.WithRFQSelectedPayinMethod("100", "STORED_BALANCE"),
		rfq.WithRFQSelectedPayoutMethod("BANK_ACCOUNT"),
		rfq.WithRFQExternalID("test_1234"),
	)

	assert.NoError(t, err)
}

func TestCreateRFQ_WithPrivate(t *testing.T) {
	pfiDID, _ := didjwk.Create()
	walletDID, _ := didjwk.Create()
	offeringID, _ := typeid.WithPrefix(offering.Kind)

	rfq, err := rfq.CreateRFQ(
		walletDID.URI,
		pfiDID.URI,
		offeringID.String(),
		rfq.WithRFQSelectedPayinMethod("100", "STORED_BALANCE"),
		rfq.WithRFQSelectedPayoutMethod("BANK_ACCOUNT", rfq.WithPayoutMethodWithPrivate(
			map[string]interface{}{
			"accountNumber":     "1234567890123456",
			"routingNumber":     "123456789",
		},)),
		rfq.WithRFQClaims([]string{"my_jwt"}),
	)

	assert.NoError(t, err)
	assert.NotZero(t, rfq.Data.Payout.PaymentDetailsHash)
	assert.NotZero(t, rfq.Data.ClaimsHash)
}

func TestRFQ_Sign(t *testing.T) {
	pfiDID, _ := didjwk.Create()

	walletDID, _ := didjwk.Create()
	offeringID, _ := typeid.WithPrefix(offering.Kind)

	r, _ := rfq.CreateRFQ(
		walletDID.URI,
		pfiDID.URI,
		offeringID.String(),
		rfq.WithRFQSelectedPayinMethod("100", "STORED_BALANCE"),
		rfq.WithRFQSelectedPayoutMethod("BANK_ACCOUNT"),
	)
	
	err := r.Sign(walletDID)
	assert.NoError(t, err)
}

func TestRFQ_UnmarshalJSON(t *testing.T) {
	pfiDID, _ := didjwk.Create()

	walletDID, _ := didjwk.Create()
	offeringID, _ := typeid.WithPrefix(offering.Kind)

	r, _ := rfq.CreateRFQ(
		walletDID.URI,
		pfiDID.URI,
		offeringID.String(),
		rfq.WithRFQSelectedPayinMethod("100", "STORED_BALANCE"),
		rfq.WithRFQSelectedPayoutMethod("BANK_ACCOUNT", rfq.WithPayoutMethodWithPrivate(
			map[string]interface{}{
			"accountNumber":     "1234567890123456",
			"routingNumber":     "123456789",
		},)),
		rfq.WithRFQClaims([]string{"my_jwt"}),
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
