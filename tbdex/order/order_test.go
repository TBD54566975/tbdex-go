package order_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/TBD54566975/tbdex-go/tbdex/order"
	"github.com/TBD54566975/tbdex-go/tbdex/rfq"
	"github.com/alecthomas/assert/v2"
	"github.com/tbd54566975/web5-go/dids/didjwk"
	"go.jetpack.io/typeid"
)

func TestUnmarshal(t *testing.T) {
	vector := `{"metadata":{"kind":"order","to":"did:jwk:eyJrdHkiOiJPS1AiLCJhbGciOiJFZERTQSIsImtpZCI6InY4M3dmLW9ETi1idUxzam5uWFBOQ21rRlJyMDFpV3ZaTHdCNmRnNE0wbWciLCJjcnYiOiJFZDI1NTE5IiwieCI6IkU2NjJfRVM2ZW9ReE9EcFludTdmUVA4OVBrX3p2Z3NRZjdUaTZlMVhSRG8ifQ","from":"did:jwk:eyJrdHkiOiJPS1AiLCJhbGciOiJFZERTQSIsImtpZCI6IkZqMm80LUpmOFhCeFJmSTdZQlRuZGVGQ3Q0V3lROEdYU05lMjVqRjZOUUkiLCJjcnYiOiJFZDI1NTE5IiwieCI6IkhHdkFHTHljVjYzSV9ONEpQX2JqazRmNVRrU19qeGJHQ1A2RUtHSGlqMGsifQ","id":"order_01hwpc95zkfhd8fsfdreatfqj7","exchangeId":"rfq_01hwpc95zke09a9zdq6b78a2qv","createdAt":"2024-04-29T21:10:32.563160","protocol":"1.0"},"data":{},"signature":"eyJhbGciOiJFZERTQSIsImtpZCI6ImRpZDpqd2s6ZXlKcmRIa2lPaUpQUzFBaUxDSmhiR2NpT2lKRlpFUlRRU0lzSW10cFpDSTZJa1pxTW04MExVcG1PRmhDZUZKbVNUZFpRbFJ1WkdWR1EzUTBWM2xST0VkWVUwNWxNalZxUmpaT1VVa2lMQ0pqY25ZaU9pSkZaREkxTlRFNUlpd2llQ0k2SWtoSGRrRkhUSGxqVmpZelNWOU9ORXBRWDJKcWF6Um1OVlJyVTE5cWVHSkhRMUEyUlV0SFNHbHFNR3NpZlEjMCJ9..Uiu_nsMRcD5F2WA7gcahX61M20lEEttUpMFSCQZNuXR42RK2z_qqjYjk85EZ1M_ILywe2DtubfZZwwFuQbcAAg"}`

	var o order.Order
	err := json.Unmarshal([]byte(vector), &o)

	fmt.Println(o.Metadata.CreatedAt)

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

func TestParse(t *testing.T) {
	// generated using tbdex-dart
	vector := `{"metadata":{"kind":"order","to":"did:jwk:eyJrdHkiOiJPS1AiLCJhbGciOiJFZERTQSIsImtpZCI6IkJmZ1hkZzAydTZETWZ2WVdJRUJDOFBYSGhUR2xMSnJFS05SdjM3N252YWciLCJjcnYiOiJFZDI1NTE5IiwieCI6Ik01NzZTVkFTNVkzY3g2ZlNPMk9EMGtZR3BHN3BLOElqM29pM0NGSzhPMzgifQ","from":"did:jwk:eyJrdHkiOiJPS1AiLCJhbGciOiJFZERTQSIsImtpZCI6InNQT1YzSnZzQldFVzdzTUtXWTBkZmZDTEt3VG1UcENreVFZcmdTOWN0S3MiLCJjcnYiOiJFZDI1NTE5IiwieCI6Inp1MF84NVN0bHRfbnhJQjMxZms4bnBNT2JoVnNkbHowbWhDel9yQzFTMFkifQ","id":"order_01hwr7hh2afa79xfr5ny7f1sg4","exchangeId":"rfq_01hwr7hh26e5s8fmebhrd29n2c","createdAt":"2024-04-30T19:26:12.041866Z","protocol":"1.0"},"data":{},"signature":"eyJhbGciOiJFZERTQSIsImtpZCI6ImRpZDpqd2s6ZXlKcmRIa2lPaUpQUzFBaUxDSmhiR2NpT2lKRlpFUlRRU0lzSW10cFpDSTZJbk5RVDFZelNuWnpRbGRGVnpkelRVdFhXVEJrWm1aRFRFdDNWRzFVY0VOcmVWRlpjbWRUT1dOMFMzTWlMQ0pqY25ZaU9pSkZaREkxTlRFNUlpd2llQ0k2SW5wMU1GODROVk4wYkhSZmJuaEpRak14Wm1zNGJuQk5UMkpvVm5Oa2JIb3diV2hEZWw5eVF6RlRNRmtpZlEjMCJ9..8MK4kzG2IqqMGqyp79j-TT6Jl341YFtvCjl5V-kM46N9sbbfCMgZNSGLaQh5bpN0qr4zsfJxmCOA1GkOIhCgDA"}`

	o, err := order.Parse([]byte(vector))

	assert.NoError(t, err)
	assert.NotZero(t, o.Metadata)
	assert.NotZero(t, o.Metadata.To)
	assert.NotZero(t, o.Metadata.From)
	assert.NotZero(t, o.Signature)
	assert.Zero(t, o.Data)
}

func TestParse_Invalid(t *testing.T) {
	// signature is kaka
	vector := `{"metadata":{"kind":"order","to":"did:jwk:eyJrdHkiOiJPS1AiLCJhbGciOiJFZERTQSIsImtpZCI6InY4M3dmLW9ETi1idUxzam5uWFBOQ21rRlJyMDFpV3ZaTHdCNmRnNE0wbWciLCJjcnYiOiJFZDI1NTE5IiwieCI6IkU2NjJfRVM2ZW9ReE9EcFludTdmUVA4OVBrX3p2Z3NRZjdUaTZlMVhSRG8ifQ","from":"did:jwk:eyJrdHkiOiJPS1AiLCJhbGciOiJFZERTQSIsImtpZCI6IkZqMm80LUpmOFhCeFJmSTdZQlRuZGVGQ3Q0V3lROEdYU05lMjVqRjZOUUkiLCJjcnYiOiJFZDI1NTE5IiwieCI6IkhHdkFHTHljVjYzSV9ONEpQX2JqazRmNVRrU19qeGJHQ1A2RUtHSGlqMGsifQ","id":"order_01hwpc95zkfhd8fsfdreatfqj7","exchangeId":"rfq_01hwpc95zke09a9zdq6b78a2qv","createdAt":"2024-04-29T21:10:32.563160","protocol":"1.0"},"data":{},"signature":"eyJhbGciOiJFZERTQSIsImtpZCI6ImRpZDpqd2s6ZXlKcmRIa2lPaUpQUzFBaUxDSmhiR2NpT2lKRlpFUlRRU0lzSW10cFpDSTZJa1pxTW04MExVcG1PRmhDZUZKbVNUZFpRbFJ1WkdWR1EzUTBWM2xST0VkWVUwNWxNalZxUmpaT1VVa2lMQ0pqY25ZaU9pSkZaREkxTlRFNUlpd2llQ0k2SWtoSGRrRkhUSGxqVmpZelNWOU9ORXBRWDJKcWF6Um1OVlJyVTE5cWVHSkhRMUEyUlV0SFNHbHFNR3NpZlEjMCJ9..Uiu_nsMRcD5F2WA7gcahX61M20lEEttUpMFSCQZNuXR42RK2z_qqjYjk85EZ1M_ILywe2DtubfZZwwFuQbcAAg"}`

	_, err := order.Parse([]byte(vector))
	assert.Error(t, err)
}

func TestVerify_Invalid(t *testing.T) {
	// signature is kaka
	vector := `{"metadata":{"kind":"order","to":"did:jwk:eyJrdHkiOiJPS1AiLCJhbGciOiJFZERTQSIsImtpZCI6InY4M3dmLW9ETi1idUxzam5uWFBOQ21rRlJyMDFpV3ZaTHdCNmRnNE0wbWciLCJjcnYiOiJFZDI1NTE5IiwieCI6IkU2NjJfRVM2ZW9ReE9EcFludTdmUVA4OVBrX3p2Z3NRZjdUaTZlMVhSRG8ifQ","from":"did:jwk:eyJrdHkiOiJPS1AiLCJhbGciOiJFZERTQSIsImtpZCI6IkZqMm80LUpmOFhCeFJmSTdZQlRuZGVGQ3Q0V3lROEdYU05lMjVqRjZOUUkiLCJjcnYiOiJFZDI1NTE5IiwieCI6IkhHdkFHTHljVjYzSV9ONEpQX2JqazRmNVRrU19qeGJHQ1A2RUtHSGlqMGsifQ","id":"order_01hwpc95zkfhd8fsfdreatfqj7","exchangeId":"rfq_01hwpc95zke09a9zdq6b78a2qv","createdAt":"2024-04-29T21:10:32.563160","protocol":"1.0"},"data":{},"signature":"eyJhbGciOiJFZERTQSIsImtpZCI6ImRpZDpqd2s6ZXlKcmRIa2lPaUpQUzFBaUxDSmhiR2NpT2lKRlpFUlRRU0lzSW10cFpDSTZJa1pxTW04MExVcG1PRmhDZUZKbVNUZFpRbFJ1WkdWR1EzUTBWM2xST0VkWVUwNWxNalZxUmpaT1VVa2lMQ0pqY25ZaU9pSkZaREkxTlRFNUlpd2llQ0k2SWtoSGRrRkhUSGxqVmpZelNWOU9ORXBRWDJKcWF6Um1OVlJyVTE5cWVHSkhRMUEyUlV0SFNHbHFNR3NpZlEjMCJ9..Uiu_nsMRcD5F2WA7gcahX61M20lEEttUpMFSCQZNuXR42RK2z_qqjYjk85EZ1M_ILywe2DtubfZZwwFuQbcAAg"}`

	var o order.Order
	err := o.UnmarshalJSON([]byte(vector))
	assert.NoError(t, err)

	err = o.Verify()
	assert.Error(t, err)
}

func TestVerify(t *testing.T) {
	// generated using tbdex-dart
	vector := `{"metadata":{"kind":"order","to":"did:jwk:eyJrdHkiOiJPS1AiLCJhbGciOiJFZERTQSIsImtpZCI6IkJmZ1hkZzAydTZETWZ2WVdJRUJDOFBYSGhUR2xMSnJFS05SdjM3N252YWciLCJjcnYiOiJFZDI1NTE5IiwieCI6Ik01NzZTVkFTNVkzY3g2ZlNPMk9EMGtZR3BHN3BLOElqM29pM0NGSzhPMzgifQ","from":"did:jwk:eyJrdHkiOiJPS1AiLCJhbGciOiJFZERTQSIsImtpZCI6InNQT1YzSnZzQldFVzdzTUtXWTBkZmZDTEt3VG1UcENreVFZcmdTOWN0S3MiLCJjcnYiOiJFZDI1NTE5IiwieCI6Inp1MF84NVN0bHRfbnhJQjMxZms4bnBNT2JoVnNkbHowbWhDel9yQzFTMFkifQ","id":"order_01hwr7hh2afa79xfr5ny7f1sg4","exchangeId":"rfq_01hwr7hh26e5s8fmebhrd29n2c","createdAt":"2024-04-30T19:26:12.041866Z","protocol":"1.0"},"data":{},"signature":"eyJhbGciOiJFZERTQSIsImtpZCI6ImRpZDpqd2s6ZXlKcmRIa2lPaUpQUzFBaUxDSmhiR2NpT2lKRlpFUlRRU0lzSW10cFpDSTZJbk5RVDFZelNuWnpRbGRGVnpkelRVdFhXVEJrWm1aRFRFdDNWRzFVY0VOcmVWRlpjbWRUT1dOMFMzTWlMQ0pqY25ZaU9pSkZaREkxTlRFNUlpd2llQ0k2SW5wMU1GODROVk4wYkhSZmJuaEpRak14Wm1zNGJuQk5UMkpvVm5Oa2JIb3diV2hEZWw5eVF6RlRNRmtpZlEjMCJ9..8MK4kzG2IqqMGqyp79j-TT6Jl341YFtvCjl5V-kM46N9sbbfCMgZNSGLaQh5bpN0qr4zsfJxmCOA1GkOIhCgDA"}`

	var o order.Order
	err := o.UnmarshalJSON([]byte(vector))
	assert.NoError(t, err)

	err = o.Verify()
	assert.NoError(t, err)
}

func TestCreate(t *testing.T) {
	alice, err := didjwk.Create()
	assert.NoError(t, err)

	pfi, err := didjwk.Create()
	assert.NoError(t, err)

	exchangeID := typeid.Must(typeid.WithPrefix(rfq.Kind))
	o := order.Create(alice.URI, pfi.URI, exchangeID.String())

	assert.NotZero(t, o)
	assert.NotZero(t, o.Metadata.ID)
	assert.Equal(t, alice.URI, o.Metadata.From)
	assert.Equal(t, pfi.URI, o.Metadata.To)
	assert.Equal(t, exchangeID.String(), o.Metadata.ExchangeID)
}

func TestSign(t *testing.T) {
	alice, err := didjwk.Create()
	assert.NoError(t, err)

	pfi, err := didjwk.Create()
	assert.NoError(t, err)

	exchangeID := typeid.Must(typeid.WithPrefix(rfq.Kind))
	o := order.Create(alice.URI, pfi.URI, exchangeID.String())

	err = o.Sign(alice)
	assert.NoError(t, err)

	err = o.Verify()
	assert.NoError(t, err)
}
