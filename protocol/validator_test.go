package protocol_test

import (
	"testing"

	"github.com/TBD54566975/tbdex-go/protocol"
	"github.com/alecthomas/assert/v2"
)

func TestValidate(t *testing.T) {
	err := protocol.Validate(protocol.TypeResource, []byte(`{"foo": "bar"}`))
	assert.Error(t, err)
}
