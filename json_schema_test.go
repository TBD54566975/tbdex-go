package jsonvalidator


import (
	"encoding/json"
	"testing"

	"github.com/alecthomas/assert"
)

func TestValidator(t *testing.T) {
	instance := `{"foo": "bar"}`

	var v interface{}
	err := json.Unmarshal([]byte(instance), &v)
	assert.NoError(t, err)

	schema := ValidatorMap["resource"]
	err = schema.Validate(v)

	assert.Error(t, err)
}

