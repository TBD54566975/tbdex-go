package tbdex

import (
	"encoding/json"
	"fmt"

	libclose "github.com/TBD54566975/tbdex-go/tbdex/closemsg"
	"github.com/TBD54566975/tbdex-go/tbdex/message"
	liborder "github.com/TBD54566975/tbdex-go/tbdex/order"
	libquote "github.com/TBD54566975/tbdex-go/tbdex/quote"
	librfq "github.com/TBD54566975/tbdex-go/tbdex/rfq"
)

// Message is the interface that all tbdex messages implement. Especially useful for decoding and parsing messages
// when the kind of message is not known upfront.
type Message interface {
	ValidNext() []string
	Kind() string
	Digest() ([]byte, error)
	// TODO: uncomment these once rfq has been refactored to separate privateStrict bool
	// Verify() error
	// Parse([]byte) (Message, error)
}

type msg struct {
	Metadata message.Metadata
}

// DecodeMessage unmarshals a message. It uses the metadata kind to determine the type of message.
// Note: unmarshaling includes validation
func DecodeMessage(data []byte) (Message, error) {
	var msg msg
	err := json.Unmarshal(data, &msg)

	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal partial message to determine kind: %w", err)
	}

	switch msg.Metadata.Kind {
	case librfq.Kind:

		var rfq librfq.RFQ
		if err := json.Unmarshal(data, &rfq); err != nil {
			return nil, fmt.Errorf("failed to unmarshal rfq: %w", err)
		}

		return rfq, nil
	case libquote.Kind:
		var quote libquote.Quote
		if err := json.Unmarshal(data, &quote); err != nil {
			return nil, fmt.Errorf("failed to unmarshal quote: %w", err)
		}

		return quote, nil
	case liborder.Kind:
		var order liborder.Order
		if err := json.Unmarshal(data, &order); err != nil {
			return nil, fmt.Errorf("failed to unmarshal order: %w", err)
		}

		return order, nil
	case libclose.Kind:
		var closemsg libclose.Close
		if err := json.Unmarshal(data, &closemsg); err != nil {
			return nil, fmt.Errorf("failed to unmarshal close: %w", err)
		}

		return closemsg, nil
	default:
		return nil, fmt.Errorf("unknown message kind: %v", msg.Metadata.Kind)
	}
}

func ParseMessage(data []byte) (Message, error) {
	var m msg
	err := json.Unmarshal(data, &m)

	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal partial message to determine kind: %w", err)
	}

	switch m.Metadata.Kind {
	case librfq.Kind:

		rfq, err := librfq.Parse(data, false)
		if err != nil {
			return nil, fmt.Errorf("failed to parse rfq: %w", err)
		}

		return rfq, nil
	case libquote.Kind:
		quote, err := libquote.Parse(data)
		if err != nil {
			return nil, fmt.Errorf("failed to parse quote: %w", err)
		}

		return quote, nil
	case liborder.Kind:
		order, err := liborder.Parse(data)
		if err != nil {
			return nil, fmt.Errorf("failed to parse order: %w", err)
		}

		return order, nil

	case libclose.Kind:
		closemsg, err := libclose.Parse(data)
		if err != nil {
			return nil, fmt.Errorf("failed to parse close: %w", err)
		}

		return closemsg, nil
	default:
		return nil, fmt.Errorf("unknown message kind: %v", m.Metadata.Kind)
	}
}
