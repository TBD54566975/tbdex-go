package tbdex_test

import (
	"testing"

	"github.com/TBD54566975/tbdex-go/tbdex"
	"github.com/alecthomas/assert/v2"
)

func TestValidate_Invalid(t *testing.T) {
	err := tbdex.Validate(tbdex.TypeResource, []byte(`{"foo": "bar"}`))
	assert.Error(t, err)
}
