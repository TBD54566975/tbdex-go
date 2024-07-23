package orderinstructions_test

import (
	"encoding/json"
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/tbd54566975/web5-go/dids/didjwk"
	"github.com/tbd54566975/web5-go/jws"

	"go.jetpack.io/typeid"

	"github.com/TBD54566975/tbdex-go/tbdex/cancel"
	"github.com/TBD54566975/tbdex-go/tbdex/closemsg"
	"github.com/TBD54566975/tbdex-go/tbdex/order"
	"github.com/TBD54566975/tbdex-go/tbdex/orderinstructions"
	"github.com/TBD54566975/tbdex-go/tbdex/orderstatus"
	"github.com/TBD54566975/tbdex-go/tbdex/quote"
	"github.com/TBD54566975/tbdex-go/tbdex/rfq"
)

func TestCreate(t *testing.T) {
	pfiDID, err := didjwk.Create()
	assert.NoError(t, err)

	walletDID, err := didjwk.Create()
	assert.NoError(t, err)

	rfqID, _ := typeid.WithPrefix(rfq.Kind)

	oi, err := orderinstructions.Create(
		pfiDID,
		walletDID.URI,
		rfqID.String(),
		orderinstructions.PayinInstruction(
			orderinstructions.Instruction("deducted from stored balance"),
			orderinstructions.Link("http://example.com/payin/123"),
		),
		orderinstructions.PayoutInstruction(
			orderinstructions.Instruction("sent to your bank account"),
			orderinstructions.Link("http://example.com/payout/123"),
		),
	)
	assert.NoError(t, err)

	assert.NotZero(t, oi.Data.Payin)
	assert.NotZero(t, oi.Data.Payout)
	assert.NotZero(t, oi.Signature)
}

func TestUnmarshal(t *testing.T) {
	pfiDID, _ := didjwk.Create()
	walletDID, _ := didjwk.Create()
	rfqID, _ := typeid.WithPrefix(rfq.Kind)

	oi, err := orderinstructions.Create(
		pfiDID,
		walletDID.URI,
		rfqID.String(),
		orderinstructions.PayinInstruction(
			orderinstructions.Instruction("deducted from stored balance"),
			orderinstructions.Link("http://example.com/payin/123"),
		),
		orderinstructions.PayoutInstruction(
			orderinstructions.Instruction("sent to your bank account"),
			orderinstructions.Link("http://example.com/payout/123"),
		),
	)
	assert.NoError(t, err)

	bytes, err := json.Marshal(oi)
	assert.NoError(t, err)

	unmarshaledOi := orderinstructions.OrderInstructions{}
	err = json.Unmarshal(bytes, &unmarshaledOi)
	assert.NoError(t, err)
}

func TestUnmarshal_Empty(t *testing.T) {
	input := []byte(`{"metadata":{},"data":{},"signature":"eyJhbGciOiJFZERTQSIsImtpZCI6ImRpZDpkaHQ6M3doZnRncGJkamloeDl6ZTl0ZG41NzV6cXptNHF3Y2NldG5mMXliaWlidXphZDdycm15eSMwIn0..ZvoVDuSrqqdXsSXgqB-U26tAU1WqUqqU_KpD1KvdYocIcmTsshjUASEwM_lUz1UnGglqkWeCIrHqrm9NNGDqBw"}`)

	oi := orderinstructions.OrderInstructions{}
	_ = json.Unmarshal(input, &oi)

	assert.Zero(t, oi.Metadata)
	assert.Zero(t, oi.Data)
}

func TestUnmarshal_Invalid(t *testing.T) {
	input := []byte(`{"doo": "doo"}`)

	oi := orderinstructions.OrderInstructions{}
	err := json.Unmarshal(input, &oi)
	assert.Error(t, err)
}

func TestVerify(t *testing.T) {
	pfiDID, _ := didjwk.Create()
	walletDID, _ := didjwk.Create()
	rfqID, _ := typeid.WithPrefix(rfq.Kind)

	oi, err := orderinstructions.Create(
		pfiDID,
		walletDID.URI,
		rfqID.String(),
		orderinstructions.PayinInstruction(
			orderinstructions.Instruction("deducted from stored balance"),
			orderinstructions.Link("http://example.com/payin/123"),
		),
		orderinstructions.PayoutInstruction(
			orderinstructions.Instruction("sent to your bank account"),
			orderinstructions.Link("http://example.com/payout/123"),
		),
	)
	assert.NoError(t, err)

	err = oi.Verify()
	assert.NoError(t, err)
}

func TestVerify_FailsChangedPayload(t *testing.T) {
	pfiDID, _ := didjwk.Create()
	walletDID, _ := didjwk.Create()
	rfqID, _ := typeid.WithPrefix(rfq.Kind)

	oi, err := orderinstructions.Create(
		pfiDID,
		walletDID.URI,
		rfqID.String(),
		orderinstructions.PayinInstruction(
			orderinstructions.Instruction("deducted from stored balance"),
			orderinstructions.Link("http://example.com/payin/123"),
		),
		orderinstructions.PayoutInstruction(
			orderinstructions.Instruction("sent to your bank account"),
			orderinstructions.Link("http://example.com/payout/123"),
		),
	)
	assert.NoError(t, err)

	oi.Data.Payin.Link = "http://hackerlink.com/payin/123"

	err = oi.Verify()
	assert.Error(t, err)
}

func TestVerify_InvalidSignature(t *testing.T) {
	pfiDID, _ := didjwk.Create()
	walletDID, _ := didjwk.Create()
	rfqID, _ := typeid.WithPrefix(rfq.Kind)

	oi, err := orderinstructions.Create(
		pfiDID,
		walletDID.URI,
		rfqID.String(),
		orderinstructions.PayinInstruction(
			orderinstructions.Instruction("deducted from stored balance"),
			orderinstructions.Link("http://example.com/payin/123"),
		),
		orderinstructions.PayoutInstruction(
			orderinstructions.Instruction("sent to your bank account"),
			orderinstructions.Link("http://example.com/payout/123"),
		),
	)
	assert.NoError(t, err)

	oi.Signature = "Invalid"

	err = oi.Verify()
	assert.Error(t, err)
}

func TestVerify_SignedWithWrongDID(t *testing.T) {
	pfiDID, _ := didjwk.Create()
	walletDID, _ := didjwk.Create()
	wrongDID, _ := didjwk.Create()
	rfqID, _ := typeid.WithPrefix(rfq.Kind)

	oi, err := orderinstructions.Create(
		pfiDID,
		walletDID.URI,
		rfqID.String(),
		orderinstructions.PayinInstruction(
			orderinstructions.Instruction("deducted from stored balance"),
			orderinstructions.Link("http://example.com/payin/123"),
		),
		orderinstructions.PayoutInstruction(
			orderinstructions.Instruction("sent to your bank account"),
			orderinstructions.Link("http://example.com/payout/123"),
		),
	)
	assert.NoError(t, err)

	toSign, err := oi.Digest()
	assert.NoError(t, err)

	wrongSignature, err := jws.Sign(toSign, wrongDID, jws.DetachedPayload(true))
	assert.NoError(t, err)

	oi.Signature = wrongSignature

	err = oi.Verify()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "does not match message metadata from")
}

func TestIsValidNext(t *testing.T) {

	pfiDID, _ := didjwk.Create()
	walletDID, _ := didjwk.Create()
	rfqID, _ := typeid.WithPrefix(rfq.Kind)

	oi, err := orderinstructions.Create(
		pfiDID,
		walletDID.URI,
		rfqID.String(),
		orderinstructions.PayinInstruction(
			orderinstructions.Instruction("deducted from stored balance"),
			orderinstructions.Link("http://example.com/payin/123"),
		),
		orderinstructions.PayoutInstruction(
			orderinstructions.Instruction("sent to your bank account"),
			orderinstructions.Link("http://example.com/payout/123"),
		),
	)
	assert.NoError(t, err)

	assert.False(t, oi.IsValidNext(rfq.Kind))
	assert.False(t, oi.IsValidNext(quote.Kind))
	assert.False(t, oi.IsValidNext(order.Kind))

	// orderinstructions can only be followed by orderstatus or cancel or close
	assert.True(t, oi.IsValidNext(orderstatus.Kind))
	assert.True(t, oi.IsValidNext(cancel.Kind))
	assert.True(t, oi.IsValidNext(closemsg.Kind))

}
