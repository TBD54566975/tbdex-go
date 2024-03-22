package protocol

import (
	"fmt"
	"os"

	"github.com/santhosh-tekuri/jsonschema/v5"
)

var ValidatorMap map[string]*jsonschema.Schema = make(map[string]*jsonschema.Schema)

type schemaObject struct {
	URL  string
	Path string
}

var schemas = map[string]schemaObject{
	"definitions": {"https://tbdex.dev/definitions.json", "../tbdex/hosted/json-schemas/definitions.json"},
	"resource":    {"https://tbdex.dev/resource.schema.json", "../tbdex/hosted/json-schemas/resource.schema.json"},
	"offering":    {"https://tbdex.dev/offering.schema.json", "../tbdex/hosted/json-schemas/offering.schema.json"},
	"balance":     {"https://tbdex.dev/balance.schema.json", "../tbdex/hosted/json-schemas/balance.schema.json"},
	"message":     {"https://tbdex.dev/message.schema.json", "../tbdex/hosted/json-schemas/message.schema.json"},
	"order":       {"https://tbdex.dev/order.schema.json", "../tbdex/hosted/json-schemas/order.schema.json"},
	"orderstatus": {"https://tbdex.dev/orderstatus.schema.json", "../tbdex/hosted/json-schemas/orderstatus.schema.json"},
	"quote":       {"https://tbdex.dev/Quote.schema.json", "../tbdex/hosted/json-schemas/Quote.schema.json"},
	"rfq":         {"https://tbdex.dev/rfq.schema.json", "../tbdex/hosted/json-schemas/rfq.schema.json"},
	"close":       {"https://tbdex.dev/close.schema.json", "../tbdex/hosted/json-schemas/close.schema.json"},
}

func init() {
    compiler := jsonschema.NewCompiler()

	for name, schemaObject := range schemas {
		file, err := os.Open(schemaObject.Path)
		if err != nil {
            fmt.Println("Failed to open schema file", err)
		}

		if err := compiler.AddResource(schemaObject.URL, file); err != nil {
            fmt.Println("Failed to add schema file", err)
		}

		file.Close()

        schema, err := compiler.Compile(schemas["resource"].URL)
        if err != nil {
            fmt.Println("Failed to compile schema", err)
        }
        ValidatorMap[name] = schema
	}
}
