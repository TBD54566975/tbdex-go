package tbdex

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gowebpki/jcs"
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

// Digest generates a SHA-256 hash of the canonicalized input payload.
func Digest(payload interface{}) ([]byte, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	canonicalized, err := jcs.Transform(payloadBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to canonicalize payload: %w", err)
	}

	hash := sha256.Sum256(canonicalized)

	return hash[:], nil
}

// VerifySignature verifies the given signature and the signed payload
func VerifySignature(digester Digester, signature string) error {
	if signature == "" {
		return errors.New("could not verify message signature because signature is empty")
	}

	payload, err := digester.Digest()
	if err != nil {
		return err
	}

	_, err = jws.Verify(signature, jws.Payload(payload))
	if err != nil {
		return err
	}

	return nil
}
