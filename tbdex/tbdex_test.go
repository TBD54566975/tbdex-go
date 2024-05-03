package tbdex_test

import (
	"testing"

	"github.com/TBD54566975/tbdex-go/tbdex"
	"github.com/TBD54566975/tbdex-go/tbdex/rfq"
	"github.com/alecthomas/assert/v2"
)

func TestParseMessage(t *testing.T) {
	t.Run("rfq", func(t *testing.T) {
		rfqvector := `{"metadata":{"kind":"rfq","to":"did:jwk:eyJrdHkiOiJPS1AiLCJhbGciOiJFZERTQSIsImtpZCI6Im1ENEYzNlVGNlUxT2FiT19TVEZJZ2tWX0R3b3pWeXVwbDFLeS1Xd25zUUkiLCJjcnYiOiJFZDI1NTE5IiwieCI6Ikh2X2JVcUE5bkR6dmJ1bkUxem5DREhybXdrdGo2Q1llTWl4TVBDUlg4Z00ifQ","from":"did:jwk:eyJrdHkiOiJPS1AiLCJhbGciOiJFZERTQSIsImtpZCI6InFsYnFMMFplZUFOcWV0UDRUS1d3RHl5d2o5cDg2b3k3cmZQQTlGNTdnRlEiLCJjcnYiOiJFZDI1NTE5IiwieCI6IjMxQWhJY1FLMjVXS2pYbzVDWWx0bVQ1SGpDaWZvemx6SzJUQ3lqdjVaWjQifQ","id":"rfq_01hwztehxhe139magy0a18mzms","exchangeId":"rfq_01hwztehxhe139magy0a18mzms","createdAt":"2024-05-03T18:11:18.577263Z","protocol":"1.0"},"data":{"offeringId":"offering_01hwztehxdezgajyyc95te7vbw","payin":{"amount":"100","kind":"DEBIT_CARD","paymentDetailsHash":"pO-bFytOXtqFsYi1fZicSb9HWGKGz5-SwDM5pYEq6QU"},"payout":{"kind":"DEBIT_CARD","paymentDetailsHash":"pO-bFytOXtqFsYi1fZicSb9HWGKGz5-SwDM5pYEq6QU"},"claimsHash":"1_FSPTu5xlVU08wgi1P-hfi77ec6sGmTUT6aRE_jOjE"},"privateData":{"salt":"tNnXU3KS5I8WLqn83Ikd3g","payin":{"paymentDetails":{"cardNumber":"0123456789012345","expiryDate":"01/21","cardHolderName":"John Meme","cvv":"123"}},"payout":{"paymentDetails":{"cardNumber":"0123456789012345","expiryDate":"01/21","cardHolderName":"John Meme","cvv":"123"}},"claims":[]},"signature":"eyJhbGciOiJFZERTQSIsImtpZCI6ImRpZDpqd2s6ZXlKcmRIa2lPaUpQUzFBaUxDSmhiR2NpT2lKRlpFUlRRU0lzSW10cFpDSTZJbkZzWW5GTU1GcGxaVUZPY1dWMFVEUlVTMWQzUkhsNWQybzVjRGcyYjNrM2NtWlFRVGxHTlRkblJsRWlMQ0pqY25ZaU9pSkZaREkxTlRFNUlpd2llQ0k2SWpNeFFXaEpZMUZMTWpWWFMycFlielZEV1d4MGJWUTFTR3BEYVdadmVteDZTekpVUTNscWRqVmFXalFpZlEjMCJ9..EmitT-FhIRkHG2761i5pujiLtGbDkEekFw2j6shE_Ni72sOVz4dipgktqQYAg4hJdB6D-F7BMv1lrvtO_IalCg"}`

		msg, err := tbdex.ParseMessage([]byte(rfqvector))
		assert.NoError(t, err)

		rfq, ok := msg.(rfq.RFQ)
		assert.True(t, ok)
		assert.NotZero(t, rfq)
	})

	t.Run("quote", func(t *testing.T) {})

	t.Run("order", func(t *testing.T) {})
	t.Run("close", func(t *testing.T) {})
}

func TestDecodeMessage(t *testing.T) {
	t.Run("rfq", func(t *testing.T) {
		rfqvector := `{"metadata":{"kind":"rfq","to":"did:jwk:eyJrdHkiOiJPS1AiLCJhbGciOiJFZERTQSIsImtpZCI6Im1ENEYzNlVGNlUxT2FiT19TVEZJZ2tWX0R3b3pWeXVwbDFLeS1Xd25zUUkiLCJjcnYiOiJFZDI1NTE5IiwieCI6Ikh2X2JVcUE5bkR6dmJ1bkUxem5DREhybXdrdGo2Q1llTWl4TVBDUlg4Z00ifQ","from":"did:jwk:eyJrdHkiOiJPS1AiLCJhbGciOiJFZERTQSIsImtpZCI6InFsYnFMMFplZUFOcWV0UDRUS1d3RHl5d2o5cDg2b3k3cmZQQTlGNTdnRlEiLCJjcnYiOiJFZDI1NTE5IiwieCI6IjMxQWhJY1FLMjVXS2pYbzVDWWx0bVQ1SGpDaWZvemx6SzJUQ3lqdjVaWjQifQ","id":"rfq_01hwztehxhe139magy0a18mzms","exchangeId":"rfq_01hwztehxhe139magy0a18mzms","createdAt":"2024-05-03T18:11:18.577263Z","protocol":"1.0"},"data":{"offeringId":"offering_01hwztehxdezgajyyc95te7vbw","payin":{"amount":"100","kind":"DEBIT_CARD","paymentDetailsHash":"pO-bFytOXtqFsYi1fZicSb9HWGKGz5-SwDM5pYEq6QU"},"payout":{"kind":"DEBIT_CARD","paymentDetailsHash":"pO-bFytOXtqFsYi1fZicSb9HWGKGz5-SwDM5pYEq6QU"},"claimsHash":"1_FSPTu5xlVU08wgi1P-hfi77ec6sGmTUT6aRE_jOjE"},"privateData":{"salt":"tNnXU3KS5I8WLqn83Ikd3g","payin":{"paymentDetails":{"cardNumber":"0123456789012345","expiryDate":"01/21","cardHolderName":"John Meme","cvv":"123"}},"payout":{"paymentDetails":{"cardNumber":"0123456789012345","expiryDate":"01/21","cardHolderName":"John Meme","cvv":"123"}},"claims":[]},"signature":"eyJhbGciOiJFZERTQSIsImtpZCI6ImRpZDpqd2s6ZXlKcmRIa2lPaUpQUzFBaUxDSmhiR2NpT2lKRlpFUlRRU0lzSW10cFpDSTZJbkZzWW5GTU1GcGxaVUZPY1dWMFVEUlVTMWQzUkhsNWQybzVjRGcyYjNrM2NtWlFRVGxHTlRkblJsRWlMQ0pqY25ZaU9pSkZaREkxTlRFNUlpd2llQ0k2SWpNeFFXaEpZMUZMTWpWWFMycFlielZEV1d4MGJWUTFTR3BEYVdadmVteDZTekpVUTNscWRqVmFXalFpZlEjMCJ9..EmitT-FhIRkHG2761i5pujiLtGbDkEekFw2j6shE_Ni72sOVz4dipgktqQYAg4hJdB6D-F7BMv1lrvtO_IalCg"}`

		msg, err := tbdex.DecodeMessage([]byte(rfqvector))
		assert.NoError(t, err)

		rfq, ok := msg.(rfq.RFQ)
		assert.True(t, ok)
		assert.NotZero(t, rfq)
	})

	t.Run("quote", func(t *testing.T) {})

	t.Run("order", func(t *testing.T) {})
	t.Run("close", func(t *testing.T) {})
}
