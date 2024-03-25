package tbdex

import (
	"fmt"

	"github.com/tbd54566975/web5-go/dids/did"
	"github.com/tbd54566975/web5-go/jws"
)

// Digester is an interface for resources that can be digested
type Digester interface {
	Digest() ([]byte, error)
}

// Sign signs a resource with a given bearerDID
func Sign(digester Digester, bearerDID did.BearerDID) (string, error) {
	digest, err := digester.Digest()
	if err != nil {
		return "", fmt.Errorf("failed to compute digest: %w", err)
	}

	signature, err := jws.Sign(digest, bearerDID, jws.DetachedPayload(true))
	if err != nil {
		return "", fmt.Errorf("failed to compute signature: %w", err)
	}

	return signature, nil
}
