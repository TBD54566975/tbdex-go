package order_test

import (
	"encoding/json"
	"testing"

	"github.com/TBD54566975/tbdex-go/tbdex/order"
	"github.com/alecthomas/assert/v2"
)

func TestUnmarshal(t *testing.T) {
	vector := `{"metadata":{"kind":"order","to":"did:jwk:eyJrdHkiOiJPS1AiLCJhbGciOiJFZERTQSIsImtpZCI6InY4M3dmLW9ETi1idUxzam5uWFBOQ21rRlJyMDFpV3ZaTHdCNmRnNE0wbWciLCJjcnYiOiJFZDI1NTE5IiwieCI6IkU2NjJfRVM2ZW9ReE9EcFludTdmUVA4OVBrX3p2Z3NRZjdUaTZlMVhSRG8ifQ","from":"did:jwk:eyJrdHkiOiJPS1AiLCJhbGciOiJFZERTQSIsImtpZCI6IkZqMm80LUpmOFhCeFJmSTdZQlRuZGVGQ3Q0V3lROEdYU05lMjVqRjZOUUkiLCJjcnYiOiJFZDI1NTE5IiwieCI6IkhHdkFHTHljVjYzSV9ONEpQX2JqazRmNVRrU19qeGJHQ1A2RUtHSGlqMGsifQ","id":"order_01hwpc95zkfhd8fsfdreatfqj7","exchangeId":"rfq_01hwpc95zke09a9zdq6b78a2qv","createdAt":"2024-04-29T21:10:32.563160","protocol":"1.0"},"data":{},"signature":"eyJhbGciOiJFZERTQSIsImtpZCI6ImRpZDpqd2s6ZXlKcmRIa2lPaUpQUzFBaUxDSmhiR2NpT2lKRlpFUlRRU0lzSW10cFpDSTZJa1pxTW04MExVcG1PRmhDZUZKbVNUZFpRbFJ1WkdWR1EzUTBWM2xST0VkWVUwNWxNalZxUmpaT1VVa2lMQ0pqY25ZaU9pSkZaREkxTlRFNUlpd2llQ0k2SWtoSGRrRkhUSGxqVmpZelNWOU9ORXBRWDJKcWF6Um1OVlJyVTE5cWVHSkhRMUEyUlV0SFNHbHFNR3NpZlEjMCJ9..Uiu_nsMRcD5F2WA7gcahX61M20lEEttUpMFSCQZNuXR42RK2z_qqjYjk85EZ1M_ILywe2DtubfZZwwFuQbcAAg"}`
	var o order.Order
	err := json.Unmarshal([]byte(vector), &o)

	assert.NoError(t, err)
	assert.NotZero(t, o.Metadata)
	assert.NotZero(t, o.Metadata.To)
	assert.NotZero(t, o.Metadata.From)
	assert.NotZero(t, o.Signature)
	assert.Zero(t, o.Data)
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
		var o order.Order
		err := json.Unmarshal([]byte(v), &o)
		assert.Error(t, err)
	}
}
