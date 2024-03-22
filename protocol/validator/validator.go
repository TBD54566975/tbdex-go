package validator

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/TBD54566975/tbdex-go"
	"github.com/santhosh-tekuri/jsonschema/v5"
)

// ValidatorMap contains a map of json schema name to schema validator
var ValidatorMap map[string]*jsonschema.Schema = make(map[string]*jsonschema.Schema)

type schemaFile struct {
	Name string
	Schema []byte
}

var schemaFiles = []schemaFile {
	{"definitions.json", tbdex.DefinitionsSchema},
	{"resource.schema.json", tbdex.ResourceSchema},
	{"offering.schema.json", tbdex.OfferingSchema},
	{"balance.schema.json", tbdex.BalanceSchema},
	{"message.schema.json", tbdex.MessageSchema},
	{"order.schema.json", tbdex.OrderSchema},
	{"orderstatus.schema.json", tbdex.OrderStatusSchema},
	{"quote.schema.json", tbdex.QuoteSchema},
	{"rfq.schema.json", tbdex.RFQSchema},
	{"close.schema.json", tbdex.CloseSchema},
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
