package tbdex

import (
	"encoding/json"
	"fmt"

	libclose "github.com/TBD54566975/tbdex-go/tbdex/closemsg"
	"github.com/TBD54566975/tbdex-go/tbdex/message"
	liborder "github.com/TBD54566975/tbdex-go/tbdex/order"
	liborderstatus "github.com/TBD54566975/tbdex-go/tbdex/orderstatus"
	libquote "github.com/TBD54566975/tbdex-go/tbdex/quote"
	librfq "github.com/TBD54566975/tbdex-go/tbdex/rfq"
)

// Message is the interface that all tbdex messages implement. Especially useful for decoding and parsing messages
// when the kind of message is not known upfront.
type Message interface {
	Digest() ([]byte, error)
	GetValidNext() []string
	GetKind() string
	GetMetadata() message.Metadata
	// TODO: uncomment these once rfq has been refactored to separate privateStrict bool
	// Verify() error
	// Parse([]byte) (Message, error)
}

type msg struct {
	Metadata message.Metadata
}

// UnmarshalMessage unmarshals a message. It uses the metadata kind to determine the type of message.
//
// # Note
//
// unmarshaling includes validation
func UnmarshalMessage(input []byte) (Message, error) {
	var m msg
	if err := json.Unmarshal(input, &m); err != nil {
		return nil, fmt.Errorf("failed to unmarshal partial message to determine kind: %w", err)
	}

	switch m.Metadata.Kind {
	case librfq.Kind:

		var rfq librfq.RFQ
		if err := json.Unmarshal(input, &rfq); err != nil {
			return nil, fmt.Errorf("failed to unmarshal rfq: %w", err)
		}

		return rfq, nil
	case libquote.Kind:
		var quote libquote.Quote
		if err := json.Unmarshal(input, &quote); err != nil {
			return nil, fmt.Errorf("failed to unmarshal quote: %w", err)
		}

		return quote, nil
	case liborder.Kind:
		var order liborder.Order
		if err := json.Unmarshal(input, &order); err != nil {
			return nil, fmt.Errorf("failed to unmarshal order: %w", err)
		}

		return order, nil

	case liborderstatus.Kind:
		var orderStatus liborderstatus.OrderStatus
		if err := json.Unmarshal(input, &orderStatus); err != nil {
			return nil, fmt.Errorf("failed to unmarshal order: %w", err)
		}

		return orderStatus, nil
	case libclose.Kind:
		var closemsg libclose.Close
		if err := json.Unmarshal(input, &closemsg); err != nil {
			return nil, fmt.Errorf("failed to unmarshal close: %w", err)
		}

		return closemsg, nil
	default:
		return nil, fmt.Errorf("unknown message kind: %v", m.Metadata.Kind)
	}
}

// ParseMessage parses a message. It uses the metadata kind to determine the type of message.
//
// # Note
//
// parsing validates the message and verifies the integrity of the message which can lead to
// a network request in order to resolve the signer's DID
func ParseMessage(input []byte) (Message, error) {
	var m msg
	err := json.Unmarshal(input, &m)

	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal partial message to determine kind: %w", err)
	}

	switch m.Metadata.Kind {
	case librfq.Kind:

		rfq, err := librfq.Parse(input, false)
		if err != nil {
			return nil, fmt.Errorf("failed to parse rfq: %w", err)
		}

		return rfq, nil
	case libquote.Kind:
		quote, err := libquote.Parse(input)
		if err != nil {
			return nil, fmt.Errorf("failed to parse quote: %w", err)
		}

		return quote, nil
	case liborder.Kind:
		order, err := liborder.Parse(input)
		if err != nil {
			return nil, fmt.Errorf("failed to parse order: %w", err)
		}

		return order, nil

	case liborderstatus.Kind:
		var orderStatus liborderstatus.OrderStatus
		if err := json.Unmarshal(input, &orderStatus); err != nil {
			return nil, fmt.Errorf("failed to unmarshal order: %w", err)
		}

		return orderStatus, nil

	case libclose.Kind:
		closemsg, err := libclose.Parse(input)
		if err != nil {
			return nil, fmt.Errorf("failed to parse close: %w", err)
		}

		return closemsg, nil
	default:
		return nil, fmt.Errorf("unknown message kind: %v", m.Metadata.Kind)
	}
}
