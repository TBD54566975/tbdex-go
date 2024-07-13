package tbdex_test

import (
	"testing"

	"github.com/TBD54566975/tbdex-go/tbdex"
	"github.com/TBD54566975/tbdex-go/tbdex/cancel"
	"github.com/TBD54566975/tbdex-go/tbdex/closemsg"
	"github.com/TBD54566975/tbdex-go/tbdex/order"
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
		vector := `{"metadata":{"from":"did:jwk:eyJrdHkiOiJPS1AiLCJjcnYiOiJFZDI1NTE5IiwieCI6IlZGZ0JFaWhFMnMxVlJnRXU3N0JlU2w4MmNRdWtoMExYbFhPVGptN2lXdkkifQ","to":"did:jwk:eyJrdHkiOiJPS1AiLCJjcnYiOiJFZDI1NTE5IiwieCI6IlJWSWNLclhnd1FKckJ6QnhPRWFkZnFyMWJaRDRhSDlMcHpSS1cxN09BYlkifQ","kind":"quote","id":"quote_01j1y0epn3eb9bhtf52nr9hjyr","exchangeId":"rfq_01j1y0epn3eb8r7nkq6x3f3zhq","createdAt":"2024-07-04T04:36:15Z","protocol":"1.0"},"data":{"expiresAt":"2024-07-04T04:36:15Z","payoutUnitsPerPayinUnit":"16.665","payin":{"currencyCode":"USD","subtotal":"10","fee":"0","total":"10"},"payout":{"currencyCode":"MXN","subtotal":"500","fee":"0","total":"500"}},"signature":"eyJhbGciOiJFZERTQSIsImtpZCI6ImRpZDpqd2s6ZXlKcmRIa2lPaUpQUzFBaUxDSmpjbllpT2lKRlpESTFOVEU1SWl3aWVDSTZJbFpHWjBKRmFXaEZNbk14VmxKblJYVTNOMEpsVTJ3NE1tTlJkV3RvTUV4WWJGaFBWR3B0TjJsWGRra2lmUSMwIn0..40Ud5k5dviIDqmniJejKLC9N7rtF9jgztThwyux3xsJBKTX76lHy-LHPE9ORqvse_wPA2OTMsu4STyHPBA50AQ"}`
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

	t.Run("orderstatus", func(t *testing.T) {
		vector := `{"metadata":{"kind":"orderstatus","to":"did:jwk:eyJrdHkiOiJPS1AiLCJhbGciOiJFZERTQSIsImtpZCI6ImRwRF9PMkhkNFRJMlptclk1SjF6LWQteFAwOHFiYi03ZDhQU0hiZGh0dWsiLCJjcnYiOiJFZDI1NTE5IiwieCI6IkhXU3dmalpIQzUtYkN0U2hPdHM2LW1QWDZsOS1hTXljVU80NzA1QzVYdE0ifQ","from":"did:jwk:eyJrdHkiOiJPS1AiLCJhbGciOiJFZERTQSIsImtpZCI6ImRhRmc5bUtpZ0YyUDJNdGI2VnZVNW9KZTRnM3BaNWg1Wk9HeGkzdUZUdkUiLCJjcnYiOiJFZDI1NTE5IiwieCI6InFSUE5tdEsyQTRDQlhVcXpHZlo4NFNVVG1pamZ3MUoxLUg4d3ZqdDFuT3MifQ","id":"orderstatus_01hy3xq944f059fzcvx9yw9fp7","exchangeId":"rfq_01hy3xq940edcv0p1ejt26t230","createdAt":"2024-05-17T18:41:09.763818Z","protocol":"1.0"},"data":{"orderStatus":"order status"},"signature":"eyJhbGciOiJFZERTQSIsImtpZCI6ImRpZDpqd2s6ZXlKcmRIa2lPaUpQUzFBaUxDSmhiR2NpT2lKRlpFUlRRU0lzSW10cFpDSTZJbVJoUm1jNWJVdHBaMFl5VURKTmRHSTJWblpWTlc5S1pUUm5NM0JhTldnMVdrOUhlR2t6ZFVaVWRrVWlMQ0pqY25ZaU9pSkZaREkxTlRFNUlpd2llQ0k2SW5GU1VFNXRkRXN5UVRSRFFsaFZjWHBIWmxvNE5GTlZWRzFwYW1aM01Vb3hMVWc0ZDNacWRERnVUM01pZlEjMCJ9..ZFnfWhdXHj-tNjtV5H62UaLA141BBkJi3a7MCB0pSLQdviGKbHMNqpEzX11zJ5Hhq0Yo0AXHtuQIGmphSrrECQ"}`
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
		vector := `{"metadata":{"kind":"quote","to":"did:jwk:eyJrdHkiOiJPS1AiLCJhbGciOiJFZERTQSIsImtpZCI6Im9NQzVLd0hPS1kzd3NUV0dzVExsVldGcVRxV3c4SGlSVEtIZnAxOGxRWG8iLCJjcnYiOiJFZDI1NTE5IiwieCI6IlhyNTljd004US1UZV9hOU5uZ19Vc2RKYlo2SnVicHpSMk5zR2VTeU5pZVEifQ","from":"did:jwk:eyJrdHkiOiJPS1AiLCJhbGciOiJFZERTQSIsImtpZCI6ImV4UFJlNlJ0NDNhcWF6cFFJMWo0dlYwNGpaWmhhZndEcjFFSGNwNWtDZDgiLCJjcnYiOiJFZDI1NTE5IiwieCI6IjRDU3c1clAxbTNDSk1udHBiak5Zdm9GYjZQMU5wZ0ctR2pYRnBrMERDZVUifQ","id":"quote_01hx0b1nmqejkb17w0zaa40bvs","exchangeId":"rfq_01hx0b1nmkeq8sjaskn2bhgc7m","createdAt":"2024-05-03T23:01:22.198580Z","protocol":"1.0"},"data":{"expiresAt":"2022-01-01T00:00:00Z","payoutUnitsPerPayinUnit":"0.00001681351","payin":{"currencyCode":"AUD","subtotal":"100","fee":"0.01","total":"100.01","paymentInstruction":{"link":"https://block.xyz","instruction":"payin instruction"}},"payout":{"currencyCode":"BTC","subtotal":"0.12","fee":"0.02","total":"0.14","paymentInstruction":{"link":"https://block.xyz","instruction":"payout instruction"}}},"signature":"eyJhbGciOiJFZERTQSIsImtpZCI6ImRpZDpqd2s6ZXlKcmRIa2lPaUpQUzFBaUxDSmhiR2NpT2lKRlpFUlRRU0lzSW10cFpDSTZJbVY0VUZKbE5sSjBORE5oY1dGNmNGRkpNV28wZGxZd05HcGFXbWhoWm5kRWNqRkZTR053Tld0RFpEZ2lMQ0pqY25ZaU9pSkZaREkxTlRFNUlpd2llQ0k2SWpSRFUzYzFjbEF4YlRORFNrMXVkSEJpYWs1WmRtOUdZalpRTVU1d1owY3RSMnBZUm5Cck1FUkRaVlVpZlEjMCJ9..H_LuRN3wdYhH29HSjrV4SzSa4LHY-f71BouYutWXEcSS4hOASC1iCFN6wZWCnI65LFWZDdLXK8fzRGHlACjNAg"}`
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

	t.Run("orderstatus", func(t *testing.T) {
		vector := `{"metadata":{"kind":"orderstatus","to":"did:jwk:eyJrdHkiOiJPS1AiLCJhbGciOiJFZERTQSIsImtpZCI6ImRwRF9PMkhkNFRJMlptclk1SjF6LWQteFAwOHFiYi03ZDhQU0hiZGh0dWsiLCJjcnYiOiJFZDI1NTE5IiwieCI6IkhXU3dmalpIQzUtYkN0U2hPdHM2LW1QWDZsOS1hTXljVU80NzA1QzVYdE0ifQ","from":"did:jwk:eyJrdHkiOiJPS1AiLCJhbGciOiJFZERTQSIsImtpZCI6ImRhRmc5bUtpZ0YyUDJNdGI2VnZVNW9KZTRnM3BaNWg1Wk9HeGkzdUZUdkUiLCJjcnYiOiJFZDI1NTE5IiwieCI6InFSUE5tdEsyQTRDQlhVcXpHZlo4NFNVVG1pamZ3MUoxLUg4d3ZqdDFuT3MifQ","id":"orderstatus_01hy3xq944f059fzcvx9yw9fp7","exchangeId":"rfq_01hy3xq940edcv0p1ejt26t230","createdAt":"2024-05-17T18:41:09.763818Z","protocol":"1.0"},"data":{"orderStatus":"order status"},"signature":"eyJhbGciOiJFZERTQSIsImtpZCI6ImRpZDpqd2s6ZXlKcmRIa2lPaUpQUzFBaUxDSmhiR2NpT2lKRlpFUlRRU0lzSW10cFpDSTZJbVJoUm1jNWJVdHBaMFl5VURKTmRHSTJWblpWTlc5S1pUUm5NM0JhTldnMVdrOUhlR2t6ZFVaVWRrVWlMQ0pqY25ZaU9pSkZaREkxTlRFNUlpd2llQ0k2SW5GU1VFNXRkRXN5UVRSRFFsaFZjWHBIWmxvNE5GTlZWRzFwYW1aM01Vb3hMVWc0ZDNacWRERnVUM01pZlEjMCJ9..ZFnfWhdXHj-tNjtV5H62UaLA141BBkJi3a7MCB0pSLQdviGKbHMNqpEzX11zJ5Hhq0Yo0AXHtuQIGmphSrrECQ"}`
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
