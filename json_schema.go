package jsonvalidator

import _ "embed"

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v5"
)

//go:embed jsonschema/definitions.json
var definitionsSchema []byte

//go:embed jsonschema/resource.schema.json
var resourceSchema []byte

//go:embed jsonschema/offering.schema.json
var offeringSchema []byte

//go:embed jsonschema/balance.schema.json
var balanceSchema []byte

//go:embed jsonschema/message.schema.json
var messageSchema []byte

//go:embed jsonschema/rfq.schema.json
var rfqSchema []byte

//go:embed jsonschema/quote.schema.json
var quoteSchema []byte

//go:embed jsonschema/order.schema.json
var orderSchema []byte

//go:embed tbdex/hosted/json-schemas/close.schema.json
var closeSchema []byte

//go:embed tbdex/hosted/json-schemas/orderstatus.schema.json
var orderStatusSchema []byte

// ValidatorMap contains a map of json schema name to schema validator
var ValidatorMap map[string]*jsonschema.Schema = make(map[string]*jsonschema.Schema)

type schemaFile struct {
	Name string
	Schema []byte
}

var schemaFiles = []schemaFile {
	{"definitions.json", definitionsSchema},
	{"resource.schema.json", resourceSchema},
	{"offering.schema.json", offeringSchema},
	{"balance.schema.json", balanceSchema},
	{"message.schema.json", messageSchema},
	{"order.schema.json", orderSchema},
	{"orderstatus.schema.json", orderStatusSchema},
	{"quote.schema.json", quoteSchema},
	{"rfq.schema.json", rfqSchema},
	{"close.schema.json", closeSchema},
}

func init() {
	compiler := jsonschema.NewCompiler()


    for _, schemaFile := range schemaFiles {
		URL := "https://tbdex.dev/" + schemaFile.Name

		reader := bytes.NewReader(schemaFile.Schema)
        if err := compiler.AddResource(URL, reader); err != nil {
			panic(fmt.Sprintf("Failed to add schema file: %s", err))
		}


        if(schemaFile.Name == "definitions.json") {
            continue
        }

        schema, err := compiler.Compile(URL)
		if err != nil {
			panic(fmt.Sprintf("Failed to compile schema: %s", err))
		}

		ValidatorMap[strings.Split(schemaFile.Name, ".")[0]] = schema
    }
}
