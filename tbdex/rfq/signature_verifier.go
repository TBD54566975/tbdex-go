package rfq

import (
	"github.com/TBD54566975/tbdex-go/tbdex"
	"github.com/tbd54566975/web5-go/jws"
)

var VerifySignatureFunc = func(r *RFQ, signature string) (*jws.Decoded, error) {
	return tbdex.VerifySignature(r, signature)
}
