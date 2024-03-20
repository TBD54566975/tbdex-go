package offering_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/TBD54566975/tbdex-go/protocol/resource/offering"
	"github.com/alecthomas/assert/v2"
	"github.com/tbd54566975/web5-go/dids/didjwk"
)

func TestOffering(t *testing.T) {
	o, err := offering.Create(
		offering.PayinDetails{
			CurrencyCode: "USD",
			Methods: []offering.PayinMethod{
				{Kind: "SQUAREPAY"},
			},
		},
		offering.PayoutDetails{
			CurrencyCode: "BTC",
			Methods: []offering.PayoutMethod{
				{
					PaymentMethod: offering.PaymentMethod{Kind: "BTC_ADDRESS"},
				},
			},
		},
		"0.000032",
	)

	assert.NoError(t, err)

	bd, err := didjwk.Create()
	assert.NoError(t, err)

	so, err := o.Sign(bd)
	assert.NoError(t, err)

	b, err := json.MarshalIndent(so, "", "  ")
	assert.NoError(t, err)

	fmt.Println(string(b))

}
