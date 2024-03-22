package validator

import (
	"fmt"
	"strings"

	"github.com/TBD54566975/tbdex-go"
	"github.com/santhosh-tekuri/jsonschema/v5"
)

// ValidatorMap contains a map of json schema name to schema validator
var ValidatorMap map[string]*jsonschema.Schema = make(map[string]*jsonschema.Schema)

var fileNames = []string{
	"definitions.json",
	"resource.schema.json",
	"offering.schema.json",
	"balance.schema.json",
	"message.schema.json",
	"order.schema.json",
	"orderstatus.schema.json",
	"quote.schema.json",
	"rfq.schema.json",
	"close.schema.json",
}

func init() {
	compiler := jsonschema.NewCompiler()

    for _, name := range fileNames {
        path := "tbdex/hosted/json-schemas/" + name
		URL := "https://tbdex.dev/" + name

        file, err := tbdex.EmbeddedFiles.Open(path)
		if err != nil {
            file.Close()
			panic(fmt.Sprintf("Failed to open schema file: %s", err))
		}

        if err := compiler.AddResource(URL, file); err != nil {
            file.Close()
			panic(fmt.Sprintf("Failed to add schema file: %s", err))
		}

        file.Close()

        if(name == "definitions.json") {
            continue
        }

        schema, err := compiler.Compile(URL)
		if err != nil {
			panic(fmt.Sprintf("Failed to compile schema: %s", err))
		}

		ValidatorMap[strings.Split(name, ".")[0]] = schema
    }
}
