package cancel_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/TBD54566975/tbdex-go/tbdex/cancel"
	"github.com/TBD54566975/tbdex-go/tbdex/closemsg"
	"github.com/TBD54566975/tbdex-go/tbdex/order"
	"github.com/TBD54566975/tbdex-go/tbdex/orderstatus"
	"github.com/TBD54566975/tbdex-go/tbdex/quote"
	"github.com/TBD54566975/tbdex-go/tbdex/rfq"

	"github.com/alecthomas/assert/v2"
	"github.com/tbd54566975/web5-go/dids/didjwk"
	"go.jetpack.io/typeid"
)

func TestUnmarshal(t *testing.T) {
	vector := `{"metadata":{"from":"did:jwk:eyJrdHkiOiJPS1AiLCJjcnYiOiJFZDI1NTE5IiwieCI6InFSUG1TSnRnUmhocklldHphSG1mUnJyaXVMaXhqS29EeDhNeFduREZRaU0ifQ","to":"did:jwk:eyJrdHkiOiJPS1AiLCJjcnYiOiJFZDI1NTE5IiwieCI6ImxqaDdqbUs2WFY2aVktUnZBRVQ1cEhva21Zem9jZnFhVmc0ODc0MHlwOHcifQ","kind":"cancel","id":"cancel_01j2fejf5eenyrt6d6xjdkh7ed","exchangeId":"rfq_01j2fejf5eeny8gycyf1ft8x3j","createdAt":"2024-07-10T23:10:03Z","protocol":"1.0"},"data":{"reason":"I don't want to do this anymore"},"signature":"eyJhbGciOiJFZERTQSIsImtpZCI6ImRpZDpqd2s6ZXlKcmRIa2lPaUpQUzFBaUxDSmpjbllpT2lKRlpESTFOVEU1SWl3aWVDSTZJbkZTVUcxVFNuUm5VbWhvY2tsbGRIcGhTRzFtVW5KeWFYVk1hWGhxUzI5RWVEaE5lRmR1UkVaUmFVMGlmUSMwIn0..sn8IoOx3bmAeaaZPPY3i2BqKS0h9Eydrp_Zkx8czQCmOvhNquXBKxMqH2nd2nsK5XuO_Poqv70aHDCAJblCeCg"}`

	var c cancel.Cancel
	err := json.Unmarshal([]byte(vector), &c)

	fmt.Println(c.Metadata.CreatedAt)

	assert.NoError(t, err)
	assert.NotZero(t, c.Metadata)
	assert.NotZero(t, c.Metadata.To)
	assert.NotZero(t, c.Metadata.From)
	assert.NotZero(t, c.Signature)
	assert.NotZero(t, c.Data)
}

func TestUnmarshal_Empty(t *testing.T) {
	vector := `{"metadata":{},"data":{},"signature":"eyJhbGciOiJFZERTQSIsImtpZCI6ImRpZDpqd2s6ZXlKcmRIa2lPaUpQUzFBaUxDSmpjbllpT2lKRlpESTFOVEU1SWl3aWVDSTZJbkZTVUcxVFNuUm5VbWhvY2tsbGRIcGhTRzFtVW5KeWFYVk1hWGhxUzI5RWVEaE5lRmR1UkVaUmFVMGlmUSMwIn0..sn8IoOx3bmAeaaZPPY3i2BqKS0h9Eydrp_Zkx8czQCmOvhNquXBKxMqH2nd2nsK5XuO_Poqv70aHDCAJblCeCg"}`

	var c cancel.Cancel
	_ = json.Unmarshal([]byte(vector), &c)

	assert.Zero(t, c.Metadata)
	assert.Zero(t, c.Data)
}

func TestUnmarshal_Invalid(t *testing.T) {
	vectors := []string{
		"hi",
		"{}",
		"[]",
		`{"metadata": {"kind": "hoarder"}}`,
		`{"metadata": {"kind": "hoarder", "from": "hehe"}}`,
	}

	for _, v := range vectors {
		var c cancel.Cancel
		err := json.Unmarshal([]byte(v), &c)
		assert.Error(t, err)
	}
}

func TestParse(t *testing.T) {
	// generated using tbdex-dart
	vector := `{"metadata":{"from":"did:jwk:eyJrdHkiOiJPS1AiLCJjcnYiOiJFZDI1NTE5IiwieCI6InFSUG1TSnRnUmhocklldHphSG1mUnJyaXVMaXhqS29EeDhNeFduREZRaU0ifQ","to":"did:jwk:eyJrdHkiOiJPS1AiLCJjcnYiOiJFZDI1NTE5IiwieCI6ImxqaDdqbUs2WFY2aVktUnZBRVQ1cEhva21Zem9jZnFhVmc0ODc0MHlwOHcifQ","kind":"cancel","id":"cancel_01j2fejf5eenyrt6d6xjdkh7ed","exchangeId":"rfq_01j2fejf5eeny8gycyf1ft8x3j","createdAt":"2024-07-10T23:10:03Z","protocol":"1.0"},"data":{"reason":"I don't want to do this anymore"},"signature":"eyJhbGciOiJFZERTQSIsImtpZCI6ImRpZDpqd2s6ZXlKcmRIa2lPaUpQUzFBaUxDSmpjbllpT2lKRlpESTFOVEU1SWl3aWVDSTZJbkZTVUcxVFNuUm5VbWhvY2tsbGRIcGhTRzFtVW5KeWFYVk1hWGhxUzI5RWVEaE5lRmR1UkVaUmFVMGlmUSMwIn0..sn8IoOx3bmAeaaZPPY3i2BqKS0h9Eydrp_Zkx8czQCmOvhNquXBKxMqH2nd2nsK5XuO_Poqv70aHDCAJblCeCg"}`

	c, err := cancel.Parse([]byte(vector))

	assert.NoError(t, err)
	assert.NotZero(t, c.Metadata)
	assert.NotZero(t, c.Metadata.To)
	assert.NotZero(t, c.Metadata.From)
	assert.NotZero(t, c.Signature)
	assert.NotZero(t, c.Data)
}

func TestParse_Invalid(t *testing.T) {
	// signature is kaka
	vector := `{"metadata":{"from":"did:jwk:eyJrdHkiOiJPS1AiLCJjcnYiOiJFZDI1NTE5IiwieCI6InFSUG1TSnRnUmhocklldHphSG1mUnJyaXVMaXhqS29EeDhNeFduREZRaU0ifQ","to":"did:jwk:eyJrdHkiOiJPS1AiLCJjcnYiOiJFZDI1NTE5IiwieCI6ImxqaDdqbUs2WFY2aVktUnZBRVQ1cEhva21Zem9jZnFhVmc0ODc0MHlwOHcifQ","kind":"cancel","id":"cancel_01j2fejf5eenyrt6d6xjdkh7ed","exchangeId":"rfq_01j2fejf5eeny8gycyf1ft8x3j","createdAt":"2024-07-10T23:10:03Z","protocol":"1.0"},"data":{"reason":"I don't want to do this anymore"},"signature":"lol"}`

	_, err := cancel.Parse([]byte(vector))
	assert.Error(t, err)
}

func TestVerify_Invalid(t *testing.T) {
	// signature is kaka
	vector := `{"metadata":{"from":"did:jwk:eyJrdHkiOiJPS1AiLCJjcnYiOiJFZDI1NTE5IiwieCI6InFSUG1TSnRnUmhocklldHphSG1mUnJyaXVMaXhqS29EeDhNeFduREZRaU0ifQ","to":"did:jwk:eyJrdHkiOiJPS1AiLCJjcnYiOiJFZDI1NTE5IiwieCI6ImxqaDdqbUs2WFY2aVktUnZBRVQ1cEhva21Zem9jZnFhVmc0ODc0MHlwOHcifQ","kind":"cancel","id":"cancel_01j2fejf5eenyrt6d6xjdkh7ed","exchangeId":"rfq_01j2fejf5eeny8gycyf1ft8x3j","createdAt":"2024-07-10T23:10:03Z","protocol":"1.0"},"data":{"reason":"I don't want to do this anymore"},"signature":"lol"}`

	var c cancel.Cancel
	err := json.Unmarshal([]byte(vector), &c)
	assert.NoError(t, err)

	err = c.Verify()
	assert.Error(t, err)
}

func TestVerify(t *testing.T) {
	vector := `{"metadata":{"from":"did:jwk:eyJrdHkiOiJPS1AiLCJjcnYiOiJFZDI1NTE5IiwieCI6InFSUG1TSnRnUmhocklldHphSG1mUnJyaXVMaXhqS29EeDhNeFduREZRaU0ifQ","to":"did:jwk:eyJrdHkiOiJPS1AiLCJjcnYiOiJFZDI1NTE5IiwieCI6ImxqaDdqbUs2WFY2aVktUnZBRVQ1cEhva21Zem9jZnFhVmc0ODc0MHlwOHcifQ","kind":"cancel","id":"cancel_01j2fejf5eenyrt6d6xjdkh7ed","exchangeId":"rfq_01j2fejf5eeny8gycyf1ft8x3j","createdAt":"2024-07-10T23:10:03Z","protocol":"1.0"},"data":{"reason":"I don't want to do this anymore"},"signature":"eyJhbGciOiJFZERTQSIsImtpZCI6ImRpZDpqd2s6ZXlKcmRIa2lPaUpQUzFBaUxDSmpjbllpT2lKRlpESTFOVEU1SWl3aWVDSTZJbkZTVUcxVFNuUm5VbWhvY2tsbGRIcGhTRzFtVW5KeWFYVk1hWGhxUzI5RWVEaE5lRmR1UkVaUmFVMGlmUSMwIn0..sn8IoOx3bmAeaaZPPY3i2BqKS0h9Eydrp_Zkx8czQCmOvhNquXBKxMqH2nd2nsK5XuO_Poqv70aHDCAJblCeCg"}`

	var c cancel.Cancel
	err := json.Unmarshal([]byte(vector), &c)
	assert.NoError(t, err)

	err = c.Verify()
	assert.NoError(t, err)
}

func TestCreate(t *testing.T) {
	alice, err := didjwk.Create()
	assert.NoError(t, err)

	pfi, err := didjwk.Create()
	assert.NoError(t, err)

	exchangeID := typeid.Must(typeid.WithPrefix(rfq.Kind))
	c, err := cancel.Create(
		alice,
		pfi.URI,
		exchangeID.String(),
		cancel.Reason("I don't want to do this anymore"),
	)
	assert.NoError(t, err)

	j, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(j))

	assert.NotZero(t, c)
	assert.NotZero(t, c.Metadata.ID)
	assert.Equal(t, alice.URI, c.Metadata.From)
	assert.Equal(t, pfi.URI, c.Metadata.To)
	assert.Equal(t, exchangeID.String(), c.Metadata.ExchangeID)
}

func TestSign(t *testing.T) {
	alice, err := didjwk.Create()
	assert.NoError(t, err)
	pfi, err := didjwk.Create()
	assert.NoError(t, err)

	exchangeID := typeid.Must(typeid.WithPrefix(rfq.Kind))
	c, err := cancel.Create(
		alice,
		pfi.URI,
		exchangeID.String(),
		cancel.Reason("I don't want to do this anymore"),
	)
	assert.NoError(t, err)

	err = c.Verify()
	assert.NoError(t, err)
}

func TestIsValidNext(t *testing.T) {

	alice, err := didjwk.Create()
	assert.NoError(t, err)
	pfi, err := didjwk.Create()
	assert.NoError(t, err)

	exchangeID := typeid.Must(typeid.WithPrefix(rfq.Kind))
	c, err := cancel.Create(alice, pfi.URI, exchangeID.String())
	assert.NoError(t, err)

	assert.False(t, c.IsValidNext(rfq.Kind))
	assert.False(t, c.IsValidNext(quote.Kind))
	assert.False(t, c.IsValidNext(order.Kind))

	// cancel can only be followed by orderstatus or close
	assert.True(t, c.IsValidNext(orderstatus.Kind))
	assert.True(t, c.IsValidNext(closemsg.Kind))
}
