package tbdex_test

import (
	"encoding/json"
	"testing"

	"github.com/TBD54566975/tbdex-go/tbdex"
	"github.com/alecthomas/assert/v2"
	"github.com/tbd54566975/web5-go/dids/didjwk"
	"go.jetpack.io/typeid"
)

func TestCreateRFQ(t *testing.T) {
	pfiDID, err := didjwk.Create()
	assert.NoError(t, err)

	walletDID, err := didjwk.Create()
	assert.NoError(t, err)

	offeringID, err := typeid.WithPrefix(tbdex.OfferingKind)
	assert.NoError(t, err)

	_, err = tbdex.CreateRFQ(
		walletDID.URI,
		pfiDID.URI,
		offeringID.String(),
		tbdex.WithRFQSelectedPayinMethod("100", "STORED_BALANCE"),
		tbdex.WithRFQSelectedPayoutMethod("BANK_ACCOUNT"),
		tbdex.WithRFQExternalID("test_1234"),
	)

	assert.NoError(t, err)
}

func TestRFQ_Sign(t *testing.T) {
	pfiDID, _ := didjwk.Create()

	walletDID, _ := didjwk.Create()
	offeringID, _ := typeid.WithPrefix(tbdex.OfferingKind)

	r, _ := tbdex.CreateRFQ(
		walletDID.URI,
		pfiDID.URI,
		offeringID.String(),
		tbdex.WithRFQSelectedPayinMethod("100", "STORED_BALANCE"),
		tbdex.WithRFQSelectedPayoutMethod("BANK_ACCOUNT"),
	)
	
	err := r.Sign(walletDID)
	assert.NoError(t, err)
}

func TestRFQ_UnmarshalJSON(t *testing.T) {
	pfiDID, _ := didjwk.Create()

	walletDID, _ := didjwk.Create()
	offeringID, _ := typeid.WithPrefix(tbdex.OfferingKind)

	r, _ := tbdex.CreateRFQ(
		walletDID.URI,
		pfiDID.URI,
		offeringID.String(),
		tbdex.WithRFQSelectedPayinMethod("100", "STORED_BALANCE"),
		tbdex.WithRFQSelectedPayoutMethod("BANK_ACCOUNT"),
	)
	
	_ = r.Sign(walletDID)

	bytes, err := json.Marshal(r)
	assert.NoError(t, err)

	var rfq tbdex.RFQ
	err = rfq.UnmarshalJSON(bytes)
	assert.NoError(t, err)
}

func TestRFQ_Unmarshal_Invalid(t *testing.T) {
	input := []byte(`{"doo": "doo"}`)

	var rfq tbdex.RFQ
	err := rfq.UnmarshalJSON(input)
	assert.Error(t, err)
}
