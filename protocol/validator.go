package protocol

import (
	"os"

	"github.com/santhosh-tekuri/jsonschema/v5"
)

var schemaPaths = map[string]string{
    "definitions": "../tbdex/hosted/json-schemas/definitions.json",
    "resource":    "../tbdex/hosted/json-schemas/resource.schema.json",
    "offering":    "../tbdex/hosted/json-schemas/offering.schema.json",
    "balance":     "../tbdex/hosted/json-schemas/balance.schema.json",
    "message":     "../tbdex/hosted/json-schemas/message.schema.json",
    "order":       "../tbdex/hosted/json-schemas/order.schema.json",
    "orderstatus": "../tbdex/hosted/json-schemas/orderstatus.schema.json",
    "Quote":       "../tbdex/hosted/json-schemas/quote.schema.json",
    "rfq":         "../tbdex/hosted/json-schemas/rfq.schema.json",
    "close":       "../tbdex/hosted/json-schemas/close.schema.json",
}

func Compiler() (*jsonschema.Compiler, error) {
    compiler := jsonschema.NewCompiler()
    
    for _, path := range schemaPaths {
        file, err := os.Open(path)
		if err != nil {
            return nil, err
        }

        if err := compiler.AddResource(path, file); err != nil {
            return nil, err
        }

        file.Close()
	}
	
    return compiler, nil
}

