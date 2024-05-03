package crypto

import (
	"crypto/sha256"
	"encoding/base64"
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

// DigestJSON generates a SHA-256 hash of the canonicalized input payload.
func DigestJSON(payload any) ([]byte, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to JSON marshal payload: %w", err)
	}

	canonicalized, err := jcs.Transform(payloadBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to canonicalize payload: %w", err)
	}

	hash := sha256.Sum256(canonicalized)

	return hash[:], nil
}

// VerifyDigest verifies that the digest of a given payload matches the expected digest.
func VerifyDigest(expectedDigest string, payload any) error {
	digestByteArray, err := DigestJSON(payload)
	if err != nil {
		return fmt.Errorf("failed to digest while verifying: %w", err)
	}
	digestEncodedString := base64.RawURLEncoding.EncodeToString(digestByteArray)

	if digestEncodedString != expectedDigest {
		return fmt.Errorf("digested payload: %v does not equal expected expectedDigest: %v", digestEncodedString, expectedDigest)
	}
	return nil
}

// VerifySignature verifies the given signature and the signed payload
func VerifySignature(digester Digester, signature string) (*jws.Decoded, error) {
	if signature == "" {
		return nil, errors.New("could not verify signature because signature is empty")
	}

	payload, err := digester.Digest()
	if err != nil {
		return nil, err
	}

	decoded, err := jws.Verify(signature, jws.Payload(payload))
	if err != nil {
		return nil, err
	}

	return &decoded, nil
}
