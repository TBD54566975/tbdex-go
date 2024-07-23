package tbdex_test

import (
	"testing"

	"github.com/TBD54566975/tbdex-go/tbdex"
	"github.com/TBD54566975/tbdex-go/tbdex/cancel"
	"github.com/TBD54566975/tbdex-go/tbdex/closemsg"
	"github.com/TBD54566975/tbdex-go/tbdex/order"
	"github.com/TBD54566975/tbdex-go/tbdex/orderinstructions"
	"github.com/TBD54566975/tbdex-go/tbdex/orderstatus"
	"github.com/TBD54566975/tbdex-go/tbdex/quote"
	"github.com/TBD54566975/tbdex-go/tbdex/rfq"
	"github.com/alecthomas/assert/v2"
)

func TestParseMessage(t *testing.T) {
	t.Run("rfq", func(t *testing.T) {
		vector := `{"metadata":{"kind":"rfq","to":"did:jwk:eyJrdHkiOiJPS1AiLCJhbGciOiJFZERTQSIsImtpZCI6Im1ENEYzNlVGNlUxT2FiT19TVEZJZ2tWX0R3b3pWeXVwbDFLeS1Xd25zUUkiLCJjcnYiOiJFZDI1NTE5IiwieCI6Ikh2X2JVcUE5bkR6dmJ1bkUxem5DREhybXdrdGo2Q1llTWl4TVBDUlg4Z00ifQ","from":"did:jwk:eyJrdHkiOiJPS1AiLCJhbGciOiJFZERTQSIsImtpZCI6InFsYnFMMFplZUFOcWV0UDRUS1d3RHl5d2o5cDg2b3k3cmZQQTlGNTdnRlEiLCJjcnYiOiJFZDI1NTE5IiwieCI6IjMxQWhJY1FLMjVXS2pYbzVDWWx0bVQ1SGpDaWZvemx6SzJUQ3lqdjVaWjQifQ","id":"rfq_01hwztehxhe139magy0a18mzms","exchangeId":"rfq_01hwztehxhe139magy0a18mzms","createdAt":"2024-05-03T18:11:18.577263Z","protocol":"1.0"},"data":{"offeringId":"offering_01hwztehxdezgajyyc95te7vbw","payin":{"amount":"100","kind":"DEBIT_CARD","paymentDetailsHash":"pO-bFytOXtqFsYi1fZicSb9HWGKGz5-SwDM5pYEq6QU"},"payout":{"kind":"DEBIT_CARD","paymentDetailsHash":"pO-bFytOXtqFsYi1fZicSb9HWGKGz5-SwDM5pYEq6QU"},"claimsHash":"1_FSPTu5xlVU08wgi1P-hfi77ec6sGmTUT6aRE_jOjE"},"privateData":{"salt":"tNnXU3KS5I8WLqn83Ikd3g","payin":{"paymentDetails":{"cardNumber":"0123456789012345","expiryDate":"01/21","cardHolderName":"John Meme","cvv":"123"}},"payout":{"paymentDetails":{"cardNumber":"0123456789012345","expiryDate":"01/21","cardHolderName":"John Meme","cvv":"123"}},"claims":[]},"signature":"eyJhbGciOiJFZERTQSIsImtpZCI6ImRpZDpqd2s6ZXlKcmRIa2lPaUpQUzFBaUxDSmhiR2NpT2lKRlpFUlRRU0lzSW10cFpDSTZJbkZzWW5GTU1GcGxaVUZPY1dWMFVEUlVTMWQzUkhsNWQybzVjRGcyYjNrM2NtWlFRVGxHTlRkblJsRWlMQ0pqY25ZaU9pSkZaREkxTlRFNUlpd2llQ0k2SWpNeFFXaEpZMUZMTWpWWFMycFlielZEV1d4MGJWUTFTR3BEYVdadmVteDZTekpVUTNscWRqVmFXalFpZlEjMCJ9..EmitT-FhIRkHG2761i5pujiLtGbDkEekFw2j6shE_Ni72sOVz4dipgktqQYAg4hJdB6D-F7BMv1lrvtO_IalCg"}`

		msg, err := tbdex.ParseMessage([]byte(vector))
		assert.NoError(t, err)

		rfq, ok := msg.(rfq.RFQ)
		assert.True(t, ok)
		assert.NotZero(t, rfq)
	})

	t.Run("quote", func(t *testing.T) {
		vector := `{"metadata":{"from":"did:jwk:eyJrdHkiOiJPS1AiLCJjcnYiOiJFZDI1NTE5IiwieCI6Ind2b1dtcUx6cC1OSlhwemxTNWYzUEpsaHFsMThaOXZRR1FpQTRRRE9qckkifQ","to":"did:jwk:eyJrdHkiOiJPS1AiLCJjcnYiOiJFZDI1NTE5IiwieCI6IjBXS1p0aEpla0U2UFdXcVZTQXBaVzB5dEoxOHEzVzU1cm13RmJNWktRamcifQ","kind":"quote","id":"quote_01j3erbrnyf4z9g77t225ekavn","exchangeId":"rfq_01j3erbrnyf4yr5p7xwm0b7dk3","createdAt":"2024-07-23T02:57:37Z","protocol":"1.0"},"data":{"expiresAt":"2024-07-23T02:57:37Z","payoutUnitsPerPayinUnit":"16.665","payin":{"currencyCode":"USD","subtotal":"10","fee":"0.1","total":"10.1"},"payout":{"currencyCode":"MXN","subtotal":"500","fee":"0","total":"500"}},"signature":"eyJhbGciOiJFZERTQSIsImtpZCI6ImRpZDpqd2s6ZXlKcmRIa2lPaUpQUzFBaUxDSmpjbllpT2lKRlpESTFOVEU1SWl3aWVDSTZJbmQyYjFkdGNVeDZjQzFPU2xod2VteFROV1l6VUVwc2FIRnNNVGhhT1haUlIxRnBRVFJSUkU5cWNra2lmUSMwIn0..sK85SbI4QgrHXOHMtCWnECmbMjHxzv3ID6zhU-84PkeHOrdyWPF5nFRMCssjN2YQNc65pK2vJBOwFEwwA_dsCQ"}`
		msg, err := tbdex.ParseMessage([]byte(vector))
		assert.NoError(t, err)

		quote, ok := msg.(quote.Quote)
		assert.True(t, ok)
		assert.NotZero(t, quote)
	})

	t.Run("order", func(t *testing.T) {
		vector := `{"metadata":{"kind":"order","to":"did:jwk:eyJrdHkiOiJPS1AiLCJhbGciOiJFZERTQSIsImtpZCI6IkJmZ1hkZzAydTZETWZ2WVdJRUJDOFBYSGhUR2xMSnJFS05SdjM3N252YWciLCJjcnYiOiJFZDI1NTE5IiwieCI6Ik01NzZTVkFTNVkzY3g2ZlNPMk9EMGtZR3BHN3BLOElqM29pM0NGSzhPMzgifQ","from":"did:jwk:eyJrdHkiOiJPS1AiLCJhbGciOiJFZERTQSIsImtpZCI6InNQT1YzSnZzQldFVzdzTUtXWTBkZmZDTEt3VG1UcENreVFZcmdTOWN0S3MiLCJjcnYiOiJFZDI1NTE5IiwieCI6Inp1MF84NVN0bHRfbnhJQjMxZms4bnBNT2JoVnNkbHowbWhDel9yQzFTMFkifQ","id":"order_01hwr7hh2afa79xfr5ny7f1sg4","exchangeId":"rfq_01hwr7hh26e5s8fmebhrd29n2c","createdAt":"2024-04-30T19:26:12.041866Z","protocol":"1.0"},"data":{},"signature":"eyJhbGciOiJFZERTQSIsImtpZCI6ImRpZDpqd2s6ZXlKcmRIa2lPaUpQUzFBaUxDSmhiR2NpT2lKRlpFUlRRU0lzSW10cFpDSTZJbk5RVDFZelNuWnpRbGRGVnpkelRVdFhXVEJrWm1aRFRFdDNWRzFVY0VOcmVWRlpjbWRUT1dOMFMzTWlMQ0pqY25ZaU9pSkZaREkxTlRFNUlpd2llQ0k2SW5wMU1GODROVk4wYkhSZmJuaEpRak14Wm1zNGJuQk5UMkpvVm5Oa2JIb3diV2hEZWw5eVF6RlRNRmtpZlEjMCJ9..8MK4kzG2IqqMGqyp79j-TT6Jl341YFtvCjl5V-kM46N9sbbfCMgZNSGLaQh5bpN0qr4zsfJxmCOA1GkOIhCgDA"}`
		msg, err := tbdex.ParseMessage([]byte(vector))
		assert.NoError(t, err)

		order, ok := msg.(order.Order)
		assert.True(t, ok)
		assert.NotZero(t, order)
	})

	t.Run("orderinstructions", func(t *testing.T) {
		vector := `{"metadata":{"from":"did:jwk:eyJrdHkiOiJPS1AiLCJjcnYiOiJFZDI1NTE5IiwieCI6IjZ6eU5iUVlKNm9CY1VUS0J1enRSdTE5bmQySVVmU0tqMFVHUy1nU2I5eWMifQ","to":"did:jwk:eyJrdHkiOiJPS1AiLCJjcnYiOiJFZDI1NTE5IiwieCI6InVxSENzMjJsMFgtbUFoQzR5ZmhFZm9mYVRrMDBFQktkaW5pbEJPR3p4c2MifQ","kind":"orderinstructions","id":"orderinstructions_01j3erdh2ceq5rqdzcdc3xysam","exchangeId":"rfq_01j3erdh2ceq59brmcgd1vqd2y","createdAt":"2024-07-23T02:58:35Z","protocol":"1.0"},"data":{"payin":{"link":"http://example.com/payin/123","instruction":"deducted from stored balance"},"payout":{"link":"http://example.com/payout/123","instruction":"sent to your bank account"}},"signature":"eyJhbGciOiJFZERTQSIsImtpZCI6ImRpZDpqd2s6ZXlKcmRIa2lPaUpQUzFBaUxDSmpjbllpT2lKRlpESTFOVEU1SWl3aWVDSTZJalo2ZVU1aVVWbEtObTlDWTFWVVMwSjFlblJTZFRFNWJtUXlTVlZtVTB0cU1GVkhVeTFuVTJJNWVXTWlmUSMwIn0..d-hh6ydEq0-w7LY4BB0Sseq3-GFEvo5BBErDauXQKto_cO72ztdqeDTiZr1I-uKrqobUV_fV4d-3S9GqvTEdBw"}`
		msg, err := tbdex.ParseMessage([]byte(vector))
		assert.NoError(t, err)

		oi, ok := msg.(orderinstructions.OrderInstructions)
		assert.True(t, ok)
		assert.NotZero(t, oi)
	})

	t.Run("orderstatus", func(t *testing.T) {
		vector := `{"metadata":{"from":"did:jwk:eyJrdHkiOiJPS1AiLCJjcnYiOiJFZDI1NTE5IiwieCI6IlhOOVB3T3R2MUEwV1VxcFp6YTFHeHFQM1VoRUxDZWhsdXBPcnp3dEl6NGMifQ","to":"did:jwk:eyJrdHkiOiJPS1AiLCJjcnYiOiJFZDI1NTE5IiwieCI6ImFHdVlsbmxiRHVZVmZMSVRJeDJkdWxPajhMVkFGaF83QzlEN2VQaU4xRnMifQ","kind":"orderstatus","id":"orderstatus_01j2hk4kt5f8jrsnmt24fnty86","exchangeId":"rfq_01j2hk4kt5f8jbb105b3kaj0jp","createdAt":"2024-07-11T19:08:21Z","protocol":"1.0"},"data":{"status":"PAYIN_INITIATED","details":"CC Payment Initiated"},"signature":"eyJhbGciOiJFZERTQSIsImtpZCI6ImRpZDpqd2s6ZXlKcmRIa2lPaUpQUzFBaUxDSmpjbllpT2lKRlpESTFOVEU1SWl3aWVDSTZJbGhPT1ZCM1QzUjJNVUV3VjFWeGNGcDZZVEZIZUhGUU0xVm9SVXhEWldoc2RYQlBjbnAzZEVsNk5HTWlmUSMwIn0..9ieXFwyxUq13OXGFQuhNjZK7rznTghkhH4McnwJ7GqBtWiQBxwvbpMEE8ZWezMhKsbLuBbLi3UBfcrdU-3VyCA"}`
		msg, err := tbdex.ParseMessage([]byte(vector))
		assert.NoError(t, err)

		orderStatus, ok := msg.(orderstatus.OrderStatus)
		assert.True(t, ok)
		assert.NotZero(t, orderStatus)
	})

	t.Run("close", func(t *testing.T) {
		vector := `{"metadata":{"kind":"close","to":"did:jwk:eyJrdHkiOiJPS1AiLCJhbGciOiJFZERTQSIsImtpZCI6IkhCVUltRVI1cm4wUzBnN1R3WEoxTHRGTU5MdnZndUNPb2RjSE1VQ3l3eGciLCJjcnYiOiJFZDI1NTE5IiwieCI6ImRxSVhMX0cxTHFkNVIwTV90blo3NC1pUDc5OXNjelhGS0pobUxDVEE4aTgifQ","from":"did:jwk:eyJrdHkiOiJPS1AiLCJhbGciOiJFZERTQSIsImtpZCI6ImYwc3BWVFo3ZUNMRVJ3RThvdjh3ZDk3WDBsNWJDcUlpV29YLUk2WVp6aW8iLCJjcnYiOiJFZDI1NTE5IiwieCI6IkFqVDJ6REdXRGdNMFZlaG93QVFSYkE1d3o1QVRDV1pFSXg0ODc1VEdsbFEifQ","id":"close_01hx0ahsy7fwns6pdznxy5yn2v","exchangeId":"rfq_01hx0ahsy3e0hv78q6r7zvkgtz","createdAt":"2024-05-03T22:52:42.311081Z","protocol":"1.0"},"data":{"reason":"reason"},"signature":"eyJhbGciOiJFZERTQSIsImtpZCI6ImRpZDpqd2s6ZXlKcmRIa2lPaUpQUzFBaUxDSmhiR2NpT2lKRlpFUlRRU0lzSW10cFpDSTZJbVl3YzNCV1ZGbzNaVU5NUlZKM1JUaHZkamgzWkRrM1dEQnNOV0pEY1VscFYyOVlMVWsyV1ZwNmFXOGlMQ0pqY25ZaU9pSkZaREkxTlRFNUlpd2llQ0k2SWtGcVZESjZSRWRYUkdkTk1GWmxhRzkzUVZGU1lrRTFkM28xUVZSRFYxcEZTWGcwT0RjMVZFZHNiRkVpZlEjMCJ9..u497k8waDvyCVNF55SqGoaRNEdftfVCjWkbLSDePfVNSsNU_LAdntfiZLI6BLL9fXMVfR1-ZCUOqad4s-ZwgDQ"}`
		msg, err := tbdex.ParseMessage([]byte(vector))
		assert.NoError(t, err)

		closemsg, ok := msg.(closemsg.Close)
		assert.True(t, ok)
		assert.NotZero(t, closemsg)
	})

	t.Run("cancel", func(t *testing.T) {
		vector := `{"metadata":{"from":"did:jwk:eyJrdHkiOiJPS1AiLCJjcnYiOiJFZDI1NTE5IiwieCI6InFSUG1TSnRnUmhocklldHphSG1mUnJyaXVMaXhqS29EeDhNeFduREZRaU0ifQ","to":"did:jwk:eyJrdHkiOiJPS1AiLCJjcnYiOiJFZDI1NTE5IiwieCI6ImxqaDdqbUs2WFY2aVktUnZBRVQ1cEhva21Zem9jZnFhVmc0ODc0MHlwOHcifQ","kind":"cancel","id":"cancel_01j2fejf5eenyrt6d6xjdkh7ed","exchangeId":"rfq_01j2fejf5eeny8gycyf1ft8x3j","createdAt":"2024-07-10T23:10:03Z","protocol":"1.0"},"data":{"reason":"I don't want to do this anymore"},"signature":"eyJhbGciOiJFZERTQSIsImtpZCI6ImRpZDpqd2s6ZXlKcmRIa2lPaUpQUzFBaUxDSmpjbllpT2lKRlpESTFOVEU1SWl3aWVDSTZJbkZTVUcxVFNuUm5VbWhvY2tsbGRIcGhTRzFtVW5KeWFYVk1hWGhxUzI5RWVEaE5lRmR1UkVaUmFVMGlmUSMwIn0..sn8IoOx3bmAeaaZPPY3i2BqKS0h9Eydrp_Zkx8czQCmOvhNquXBKxMqH2nd2nsK5XuO_Poqv70aHDCAJblCeCg"}`
		msg, err := tbdex.ParseMessage([]byte(vector))
		assert.NoError(t, err)

		cancel, ok := msg.(cancel.Cancel)
		assert.True(t, ok)
		assert.NotZero(t, cancel)
	})
}

func TestUnmarshalMessage(t *testing.T) {
	t.Run("rfq", func(t *testing.T) {
		rfqvector := `{"metadata":{"kind":"rfq","to":"did:jwk:eyJrdHkiOiJPS1AiLCJhbGciOiJFZERTQSIsImtpZCI6Im1ENEYzNlVGNlUxT2FiT19TVEZJZ2tWX0R3b3pWeXVwbDFLeS1Xd25zUUkiLCJjcnYiOiJFZDI1NTE5IiwieCI6Ikh2X2JVcUE5bkR6dmJ1bkUxem5DREhybXdrdGo2Q1llTWl4TVBDUlg4Z00ifQ","from":"did:jwk:eyJrdHkiOiJPS1AiLCJhbGciOiJFZERTQSIsImtpZCI6InFsYnFMMFplZUFOcWV0UDRUS1d3RHl5d2o5cDg2b3k3cmZQQTlGNTdnRlEiLCJjcnYiOiJFZDI1NTE5IiwieCI6IjMxQWhJY1FLMjVXS2pYbzVDWWx0bVQ1SGpDaWZvemx6SzJUQ3lqdjVaWjQifQ","id":"rfq_01hwztehxhe139magy0a18mzms","exchangeId":"rfq_01hwztehxhe139magy0a18mzms","createdAt":"2024-05-03T18:11:18.577263Z","protocol":"1.0"},"data":{"offeringId":"offering_01hwztehxdezgajyyc95te7vbw","payin":{"amount":"100","kind":"DEBIT_CARD","paymentDetailsHash":"pO-bFytOXtqFsYi1fZicSb9HWGKGz5-SwDM5pYEq6QU"},"payout":{"kind":"DEBIT_CARD","paymentDetailsHash":"pO-bFytOXtqFsYi1fZicSb9HWGKGz5-SwDM5pYEq6QU"},"claimsHash":"1_FSPTu5xlVU08wgi1P-hfi77ec6sGmTUT6aRE_jOjE"},"privateData":{"salt":"tNnXU3KS5I8WLqn83Ikd3g","payin":{"paymentDetails":{"cardNumber":"0123456789012345","expiryDate":"01/21","cardHolderName":"John Meme","cvv":"123"}},"payout":{"paymentDetails":{"cardNumber":"0123456789012345","expiryDate":"01/21","cardHolderName":"John Meme","cvv":"123"}},"claims":[]},"signature":"eyJhbGciOiJFZERTQSIsImtpZCI6ImRpZDpqd2s6ZXlKcmRIa2lPaUpQUzFBaUxDSmhiR2NpT2lKRlpFUlRRU0lzSW10cFpDSTZJbkZzWW5GTU1GcGxaVUZPY1dWMFVEUlVTMWQzUkhsNWQybzVjRGcyYjNrM2NtWlFRVGxHTlRkblJsRWlMQ0pqY25ZaU9pSkZaREkxTlRFNUlpd2llQ0k2SWpNeFFXaEpZMUZMTWpWWFMycFlielZEV1d4MGJWUTFTR3BEYVdadmVteDZTekpVUTNscWRqVmFXalFpZlEjMCJ9..EmitT-FhIRkHG2761i5pujiLtGbDkEekFw2j6shE_Ni72sOVz4dipgktqQYAg4hJdB6D-F7BMv1lrvtO_IalCg"}`

		msg, err := tbdex.UnmarshalMessage([]byte(rfqvector))
		assert.NoError(t, err)

		rfq, ok := msg.(rfq.RFQ)
		assert.True(t, ok)
		assert.NotZero(t, rfq)
	})

	t.Run("quote", func(t *testing.T) {
		vector := `{"metadata":{"from":"did:jwk:eyJrdHkiOiJPS1AiLCJjcnYiOiJFZDI1NTE5IiwieCI6Ind2b1dtcUx6cC1OSlhwemxTNWYzUEpsaHFsMThaOXZRR1FpQTRRRE9qckkifQ","to":"did:jwk:eyJrdHkiOiJPS1AiLCJjcnYiOiJFZDI1NTE5IiwieCI6IjBXS1p0aEpla0U2UFdXcVZTQXBaVzB5dEoxOHEzVzU1cm13RmJNWktRamcifQ","kind":"quote","id":"quote_01j3erbrnyf4z9g77t225ekavn","exchangeId":"rfq_01j3erbrnyf4yr5p7xwm0b7dk3","createdAt":"2024-07-23T02:57:37Z","protocol":"1.0"},"data":{"expiresAt":"2024-07-23T02:57:37Z","payoutUnitsPerPayinUnit":"16.665","payin":{"currencyCode":"USD","subtotal":"10","fee":"0.1","total":"10.1"},"payout":{"currencyCode":"MXN","subtotal":"500","fee":"0","total":"500"}},"signature":"eyJhbGciOiJFZERTQSIsImtpZCI6ImRpZDpqd2s6ZXlKcmRIa2lPaUpQUzFBaUxDSmpjbllpT2lKRlpESTFOVEU1SWl3aWVDSTZJbmQyYjFkdGNVeDZjQzFPU2xod2VteFROV1l6VUVwc2FIRnNNVGhhT1haUlIxRnBRVFJSUkU5cWNra2lmUSMwIn0..sK85SbI4QgrHXOHMtCWnECmbMjHxzv3ID6zhU-84PkeHOrdyWPF5nFRMCssjN2YQNc65pK2vJBOwFEwwA_dsCQ"}`
		msg, err := tbdex.UnmarshalMessage([]byte(vector))
		assert.NoError(t, err)

		quote, ok := msg.(quote.Quote)
		assert.True(t, ok)
		assert.NotZero(t, quote)
	})

	t.Run("order", func(t *testing.T) {
		vector := `{"metadata":{"kind":"order","to":"did:jwk:eyJrdHkiOiJPS1AiLCJhbGciOiJFZERTQSIsImtpZCI6IkJmZ1hkZzAydTZETWZ2WVdJRUJDOFBYSGhUR2xMSnJFS05SdjM3N252YWciLCJjcnYiOiJFZDI1NTE5IiwieCI6Ik01NzZTVkFTNVkzY3g2ZlNPMk9EMGtZR3BHN3BLOElqM29pM0NGSzhPMzgifQ","from":"did:jwk:eyJrdHkiOiJPS1AiLCJhbGciOiJFZERTQSIsImtpZCI6InNQT1YzSnZzQldFVzdzTUtXWTBkZmZDTEt3VG1UcENreVFZcmdTOWN0S3MiLCJjcnYiOiJFZDI1NTE5IiwieCI6Inp1MF84NVN0bHRfbnhJQjMxZms4bnBNT2JoVnNkbHowbWhDel9yQzFTMFkifQ","id":"order_01hwr7hh2afa79xfr5ny7f1sg4","exchangeId":"rfq_01hwr7hh26e5s8fmebhrd29n2c","createdAt":"2024-04-30T19:26:12.041866Z","protocol":"1.0"},"data":{},"signature":"eyJhbGciOiJFZERTQSIsImtpZCI6ImRpZDpqd2s6ZXlKcmRIa2lPaUpQUzFBaUxDSmhiR2NpT2lKRlpFUlRRU0lzSW10cFpDSTZJbk5RVDFZelNuWnpRbGRGVnpkelRVdFhXVEJrWm1aRFRFdDNWRzFVY0VOcmVWRlpjbWRUT1dOMFMzTWlMQ0pqY25ZaU9pSkZaREkxTlRFNUlpd2llQ0k2SW5wMU1GODROVk4wYkhSZmJuaEpRak14Wm1zNGJuQk5UMkpvVm5Oa2JIb3diV2hEZWw5eVF6RlRNRmtpZlEjMCJ9..8MK4kzG2IqqMGqyp79j-TT6Jl341YFtvCjl5V-kM46N9sbbfCMgZNSGLaQh5bpN0qr4zsfJxmCOA1GkOIhCgDA"}`
		msg, err := tbdex.UnmarshalMessage([]byte(vector))
		assert.NoError(t, err)

		order, ok := msg.(order.Order)
		assert.True(t, ok)
		assert.NotZero(t, order)
	})

	t.Run("orderinstructions", func(t *testing.T) {
		vector := `{"metadata":{"from":"did:jwk:eyJrdHkiOiJPS1AiLCJjcnYiOiJFZDI1NTE5IiwieCI6IjZ6eU5iUVlKNm9CY1VUS0J1enRSdTE5bmQySVVmU0tqMFVHUy1nU2I5eWMifQ","to":"did:jwk:eyJrdHkiOiJPS1AiLCJjcnYiOiJFZDI1NTE5IiwieCI6InVxSENzMjJsMFgtbUFoQzR5ZmhFZm9mYVRrMDBFQktkaW5pbEJPR3p4c2MifQ","kind":"orderinstructions","id":"orderinstructions_01j3erdh2ceq5rqdzcdc3xysam","exchangeId":"rfq_01j3erdh2ceq59brmcgd1vqd2y","createdAt":"2024-07-23T02:58:35Z","protocol":"1.0"},"data":{"payin":{"link":"http://example.com/payin/123","instruction":"deducted from stored balance"},"payout":{"link":"http://example.com/payout/123","instruction":"sent to your bank account"}},"signature":"eyJhbGciOiJFZERTQSIsImtpZCI6ImRpZDpqd2s6ZXlKcmRIa2lPaUpQUzFBaUxDSmpjbllpT2lKRlpESTFOVEU1SWl3aWVDSTZJalo2ZVU1aVVWbEtObTlDWTFWVVMwSjFlblJTZFRFNWJtUXlTVlZtVTB0cU1GVkhVeTFuVTJJNWVXTWlmUSMwIn0..d-hh6ydEq0-w7LY4BB0Sseq3-GFEvo5BBErDauXQKto_cO72ztdqeDTiZr1I-uKrqobUV_fV4d-3S9GqvTEdBw"}`
		msg, err := tbdex.UnmarshalMessage([]byte(vector))
		assert.NoError(t, err)

		orderInstructions, ok := msg.(orderinstructions.OrderInstructions)
		assert.True(t, ok)
		assert.NotZero(t, orderInstructions)
	})

	t.Run("orderstatus", func(t *testing.T) {
		vector := `{"metadata":{"from":"did:jwk:eyJrdHkiOiJPS1AiLCJjcnYiOiJFZDI1NTE5IiwieCI6IlhOOVB3T3R2MUEwV1VxcFp6YTFHeHFQM1VoRUxDZWhsdXBPcnp3dEl6NGMifQ","to":"did:jwk:eyJrdHkiOiJPS1AiLCJjcnYiOiJFZDI1NTE5IiwieCI6ImFHdVlsbmxiRHVZVmZMSVRJeDJkdWxPajhMVkFGaF83QzlEN2VQaU4xRnMifQ","kind":"orderstatus","id":"orderstatus_01j2hk4kt5f8jrsnmt24fnty86","exchangeId":"rfq_01j2hk4kt5f8jbb105b3kaj0jp","createdAt":"2024-07-11T19:08:21Z","protocol":"1.0"},"data":{"status":"PAYIN_INITIATED","details":"CC Payment Initiated"},"signature":"eyJhbGciOiJFZERTQSIsImtpZCI6ImRpZDpqd2s6ZXlKcmRIa2lPaUpQUzFBaUxDSmpjbllpT2lKRlpESTFOVEU1SWl3aWVDSTZJbGhPT1ZCM1QzUjJNVUV3VjFWeGNGcDZZVEZIZUhGUU0xVm9SVXhEWldoc2RYQlBjbnAzZEVsNk5HTWlmUSMwIn0..9ieXFwyxUq13OXGFQuhNjZK7rznTghkhH4McnwJ7GqBtWiQBxwvbpMEE8ZWezMhKsbLuBbLi3UBfcrdU-3VyCA"}`
		msg, err := tbdex.UnmarshalMessage([]byte(vector))
		assert.NoError(t, err)

		orderStatus, ok := msg.(orderstatus.OrderStatus)
		assert.True(t, ok)
		assert.NotZero(t, orderStatus)
	})

	t.Run("close", func(t *testing.T) {
		vector := `{"metadata":{"kind":"close","to":"did:jwk:eyJrdHkiOiJPS1AiLCJhbGciOiJFZERTQSIsImtpZCI6IkhCVUltRVI1cm4wUzBnN1R3WEoxTHRGTU5MdnZndUNPb2RjSE1VQ3l3eGciLCJjcnYiOiJFZDI1NTE5IiwieCI6ImRxSVhMX0cxTHFkNVIwTV90blo3NC1pUDc5OXNjelhGS0pobUxDVEE4aTgifQ","from":"did:jwk:eyJrdHkiOiJPS1AiLCJhbGciOiJFZERTQSIsImtpZCI6ImYwc3BWVFo3ZUNMRVJ3RThvdjh3ZDk3WDBsNWJDcUlpV29YLUk2WVp6aW8iLCJjcnYiOiJFZDI1NTE5IiwieCI6IkFqVDJ6REdXRGdNMFZlaG93QVFSYkE1d3o1QVRDV1pFSXg0ODc1VEdsbFEifQ","id":"close_01hx0ahsy7fwns6pdznxy5yn2v","exchangeId":"rfq_01hx0ahsy3e0hv78q6r7zvkgtz","createdAt":"2024-05-03T22:52:42.311081Z","protocol":"1.0"},"data":{"reason":"reason"},"signature":"eyJhbGciOiJFZERTQSIsImtpZCI6ImRpZDpqd2s6ZXlKcmRIa2lPaUpQUzFBaUxDSmhiR2NpT2lKRlpFUlRRU0lzSW10cFpDSTZJbVl3YzNCV1ZGbzNaVU5NUlZKM1JUaHZkamgzWkRrM1dEQnNOV0pEY1VscFYyOVlMVWsyV1ZwNmFXOGlMQ0pqY25ZaU9pSkZaREkxTlRFNUlpd2llQ0k2SWtGcVZESjZSRWRYUkdkTk1GWmxhRzkzUVZGU1lrRTFkM28xUVZSRFYxcEZTWGcwT0RjMVZFZHNiRkVpZlEjMCJ9..u497k8waDvyCVNF55SqGoaRNEdftfVCjWkbLSDePfVNSsNU_LAdntfiZLI6BLL9fXMVfR1-ZCUOqad4s-ZwgDQ"}`
		msg, err := tbdex.UnmarshalMessage([]byte(vector))
		assert.NoError(t, err)

		closemsg, ok := msg.(closemsg.Close)
		assert.True(t, ok)
		assert.NotZero(t, closemsg)
	})

	t.Run("cancel", func(t *testing.T) {
		vector := `{"metadata":{"from":"did:jwk:eyJrdHkiOiJPS1AiLCJjcnYiOiJFZDI1NTE5IiwieCI6InFSUG1TSnRnUmhocklldHphSG1mUnJyaXVMaXhqS29EeDhNeFduREZRaU0ifQ","to":"did:jwk:eyJrdHkiOiJPS1AiLCJjcnYiOiJFZDI1NTE5IiwieCI6ImxqaDdqbUs2WFY2aVktUnZBRVQ1cEhva21Zem9jZnFhVmc0ODc0MHlwOHcifQ","kind":"cancel","id":"cancel_01j2fejf5eenyrt6d6xjdkh7ed","exchangeId":"rfq_01j2fejf5eeny8gycyf1ft8x3j","createdAt":"2024-07-10T23:10:03Z","protocol":"1.0"},"data":{"reason":"I don't want to do this anymore"},"signature":"eyJhbGciOiJFZERTQSIsImtpZCI6ImRpZDpqd2s6ZXlKcmRIa2lPaUpQUzFBaUxDSmpjbllpT2lKRlpESTFOVEU1SWl3aWVDSTZJbkZTVUcxVFNuUm5VbWhvY2tsbGRIcGhTRzFtVW5KeWFYVk1hWGhxUzI5RWVEaE5lRmR1UkVaUmFVMGlmUSMwIn0..sn8IoOx3bmAeaaZPPY3i2BqKS0h9Eydrp_Zkx8czQCmOvhNquXBKxMqH2nd2nsK5XuO_Poqv70aHDCAJblCeCg"}`
		msg, err := tbdex.UnmarshalMessage([]byte(vector))
		assert.NoError(t, err)

		cancel, ok := msg.(cancel.Cancel)
		assert.True(t, ok)
		assert.NotZero(t, cancel)
	})
}
