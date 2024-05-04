package validator_test

import (
	"testing"

	"github.com/TBD54566975/tbdex-go/tbdex/validator"
	"github.com/alecthomas/assert"
)

func TestValidate_Invalid(t *testing.T) {
	err := validator.Validate(validator.TypeResource, []byte(`{"foo": "bar"}`))
	assert.Error(t, err)
}
