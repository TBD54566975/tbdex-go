package tbdex

import (
	"encoding/json"
	"errors"

	libclose "github.com/TBD54566975/tbdex-go/tbdex/closemsg"
	"github.com/TBD54566975/tbdex-go/tbdex/message"
	liborder "github.com/TBD54566975/tbdex-go/tbdex/order"
	libquote "github.com/TBD54566975/tbdex-go/tbdex/quote"
	librfq "github.com/TBD54566975/tbdex-go/tbdex/rfq"
)

type Message interface {
	ValidNext() []string
	Kind() string
	Digest() ([]byte, error)
	// Verify() error
	// Parse([]byte) (Message, error)
}

type msg struct {
	Metadata message.Metadata
}

func DecodeMessage(data []byte) (Message, error) {
	var msg msg
	err := json.Unmarshal(data, &msg)

	if err != nil {
		return nil, err
	}

	switch msg.Metadata.Kind {
	case librfq.Kind:

		var rfq librfq.RFQ
		if err := json.Unmarshal(data, &rfq); err != nil {
			return nil, err
		}

		return rfq, nil
	case libquote.Kind:
		var quote libquote.Quote
		if err := json.Unmarshal(data, &quote); err != nil {
			return nil, err
		}

		return quote, nil
	case liborder.Kind:
		var order liborder.Order
		if err := json.Unmarshal(data, &order); err != nil {
			return nil, err
		}

		return order, nil
	case libclose.Kind:
		var closemsg libclose.Close
		if err := json.Unmarshal(data, &closemsg); err != nil {
			return nil, err
		}

		return closemsg, nil
	default:
		return nil, errors.New("unknown message kind")
	}
}

func ParseMessage(data []byte) (Message, error) {
	var m msg
	err := json.Unmarshal(data, &m)

	if err != nil {
		return nil, err
	}

	switch m.Metadata.Kind {
	case librfq.Kind:

		rfq, err := librfq.Parse(data, false)
		if err != nil {
			return nil, err
		}

		return rfq, nil
	case libquote.Kind:
		quote, err := libquote.Parse(data)
		if err != nil {
			return nil, err
		}

		return quote, nil
	default:
		return nil, errors.New("unknown message kind")
	}
}
