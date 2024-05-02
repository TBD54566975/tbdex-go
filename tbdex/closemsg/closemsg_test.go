package closemsg_test

import (
	"encoding/json"
	"testing"

	"github.com/TBD54566975/tbdex-go/tbdex/closemsg"
	"github.com/TBD54566975/tbdex-go/tbdex/rfq"
	"github.com/alecthomas/assert/v2"
	"github.com/tbd54566975/web5-go/dids/didjwk"
	"github.com/tbd54566975/web5-go/jws"
	"go.jetpack.io/typeid"
)

func TestCreate(t *testing.T) {
	pfiDID, err := didjwk.Create()
	assert.NoError(t, err)

	walletDID, err := didjwk.Create()
	assert.NoError(t, err)

	rfqID, _ := typeid.WithPrefix(rfq.Kind)

	c, err := closemsg.Create(
		pfiDID,
		walletDID.URI,
		rfqID.String(),
		closemsg.Reason("card declined"),
		closemsg.Success(false),
	)

	assert.NoError(t, err)
	assert.NotZero(t, c.Data.Reason)
	assert.Equal(t, "card declined", c.Data.Reason)
	assert.False(t, c.Data.Success)
	assert.NotZero(t, c.Signature)
}

func TestUnmarshalJSON(t *testing.T) {
	pfiDID, _ := didjwk.Create()
	walletDID, _ := didjwk.Create()
	rfqID, _ := typeid.WithPrefix(rfq.Kind)

	message, _ := closemsg.Create(
		pfiDID,
		walletDID.URI,
		rfqID.String(),
		closemsg.Reason("test"),
	)

	bytes, err := json.Marshal(message)
	assert.NoError(t, err)

	c := closemsg.Close{}
	err = json.Unmarshal(bytes, &c)
	assert.NoError(t, err)
}

func TestUnmarshal_Invalid(t *testing.T) {
	input := []byte(`{"doo": "doo"}`)

	c := closemsg.Close{}
	err := json.Unmarshal(input, &c)
	assert.Error(t, err)
}

func TestUnmarshal_Empty(t *testing.T) {
	input := []byte(`{"metadata":{},"data":{},"signature":"eyJhbGciOiJFZERTQSIsImtpZCI6ImRpZDpkaHQ6M3doZnRncGJkamloeDl6ZTl0ZG41NzV6cXptNHF3Y2NldG5mMXliaWlidXphZDdycm15eSMwIn0..ZvoVDuSrqqdXsSXgqB-U26tAU1WqUqqU_KpD1KvdYocIcmTsshjUASEwM_lUz1UnGglqkWeCIrHqrm9NNGDqBw"}`)

	c := closemsg.Close{}
	_ = json.Unmarshal(input, &c)

	assert.Zero(t, c.Metadata)
	assert.Zero(t, c.Data)
}

func TestVerify(t *testing.T) {
	pfiDID, _ := didjwk.Create()
	walletDID, _ := didjwk.Create()
	rfqID, _ := typeid.WithPrefix(rfq.Kind)

	c, _ := closemsg.Create(
		pfiDID,
		walletDID.URI,
		rfqID.String(),
	)

	err := c.Verify()
	assert.NoError(t, err)
}

func TestVerify_FailsChangedPayload(t *testing.T) {
	pfiDID, _ := didjwk.Create()
	walletDID, _ := didjwk.Create()
	rfqID, _ := typeid.WithPrefix(rfq.Kind)

	c, _ := closemsg.Create(
		pfiDID,
		walletDID.URI,
		rfqID.String(),
		closemsg.Reason("test"),
	)

	c.Data.Reason = "new reason"

	err := c.Verify()
	assert.Error(t, err)
}

func TestVerify_InvalidSignature(t *testing.T) {
	pfiDID, _ := didjwk.Create()
	walletDID, _ := didjwk.Create()
	rfqID, _ := typeid.WithPrefix(rfq.Kind)

	c, _ := closemsg.Create(
		pfiDID,
		walletDID.URI,
		rfqID.String(),
		closemsg.Reason("test"),
	)

	c.Signature = "Invalid"

	err := c.Verify()
	assert.Error(t, err)
}

func TestVerify_SignedWithWrongDID(t *testing.T) {
	pfiDID, _ := didjwk.Create()
	walletDID, _ := didjwk.Create()
	wrongDID, _ := didjwk.Create()
	rfqID, _ := typeid.WithPrefix(rfq.Kind)

	c, _ := closemsg.Create(
		pfiDID,
		walletDID.URI,
		rfqID.String(),
		closemsg.Reason("test"),
	)

	toSign, err := c.Digest()
	assert.NoError(t, err)

	wrongSignature, err := jws.Sign(toSign, wrongDID, jws.DetachedPayload(true))
	assert.NoError(t, err)

	c.Signature = wrongSignature

	err = c.Verify()
	assert.Error(t, err)
}
