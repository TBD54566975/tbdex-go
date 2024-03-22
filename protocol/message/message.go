package message

import (
	"fmt"

	"github.com/tbd54566975/web5-go/dids/did"
	"github.com/tbd54566975/web5-go/jws"
)

// Metadata is the message kind agnostic data
type Metadata struct {
	From       string `json:"from"`
	To         string `json:"to"`
	Kind       string `json:"kind"`
	ID         string `json:"id"`
	ExchangeID string `json:"exchangeId"`
	ExternalID string `json:"externalId,omitempty"`
	CreatedAt  string `json:"createdAt"`
	Protocol   string `json:"protocol"`
}

// Digester is an interface for messages that can be digested
type Digester interface {
	Digest() ([]byte, error)
}

// Sign signs a message with a given bearerDID
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
