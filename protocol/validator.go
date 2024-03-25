package protocol

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/santhosh-tekuri/jsonschema/v5"
)

// DataType represents the type of data being validated e.g. resource or message
type DataType = string

const (
	TypeResource      DataType = "resource" // TypeResource represents tbdex resource
	TypeMessage       DataType = "message"  // TypeMessage represents tbdex message
	definitionsSchema          = "definitions.json"
	schemaDir                  = "json-schemas/"
	schemaHost                 = "https://tbdex.dev/"
	schemaExtension            = ".schema.json"
)

//go:embed json-schemas
var embeddedSchemas embed.FS
var schemaCompiler *jsonschema.Compiler
var schemaMap map[string]*jsonschema.Schema = make(map[string]*jsonschema.Schema)

// init does the following:
//   - initializes the schema compiler
//   - loads the shared definitions schema
//   - loads the resource and message schemas
func init() {
	schemaCompiler = jsonschema.NewCompiler()
	schemaCompiler.Draft = jsonschema.Draft7

	definitions, err := embeddedSchemas.Open(schemaDir + definitionsSchema)
	if err != nil {
		panic(err)
	}

	err = schemaCompiler.AddResource(schemaHost+definitionsSchema, definitions)
	if err != nil {
		panic(err)
	}

	for _, schemaName := range []string{TypeResource, TypeMessage} {
		_, err = loadSchema(schemaName)
		if err != nil {
			panic(err)
		}
	}
}

type validateOptions struct {
	kind string
}

// ValidateOption is a function that sets an option for the Validate function
type ValidateOption func(*validateOptions)

// WithKind sets the kind option for the Validate function
func WithKind(kind string) ValidateOption {
	return func(o *validateOptions) {
		o.kind = kind
	}
}

// Validate validates the input provided in two phases:
//  1. Validate the general structure of the resource or message based on the Type.
//  2. Validate the specific structure of the resource or message based on the Kind.
//
// A Kind can be optionally specified in order to fail early if the input's Kind does match
// what was provided. This is useful when the Kind is known ahead of time. If the Kind is not
// specified, validation will proceed to phase 2 using metadata.kind.
//
// Note: Kind-specific schemas are lazily loaded the first time they are needed and then
// cached for future use.
func Validate(dataType DataType, input []byte, opts ...ValidateOption) error {
	var options validateOptions
	for _, o := range opts {
		o(&options)
	}

	var v any
	err := json.Unmarshal(input, &v)
	if err != nil {
		return fmt.Errorf("failed to unmarshal resource: %w", err)
	}

	typeSchema := schemaMap[TypeResource]
	err = typeSchema.Validate(v)
	if err != nil {
		return fmt.Errorf("failed to validate resource: %w", err)
	}

	resource, ok := v.(map[string]any)
	if !ok {
		return fmt.Errorf("expected resource to be an object: %w", err)
	}

	metadata, _ := resource["metadata"].(map[string]any)
	kind, _ := metadata["kind"].(string)

	if options.kind != "" && kind != options.kind {
		return errors.New("kind mismatch")
	}

	kindSchema, err := loadSchema(kind)
	if err != nil {
		return fmt.Errorf("failed to validate resource: %w", err)
	}

	err = kindSchema.Validate(resource["data"])
	if err != nil {
		return fmt.Errorf("failed to validate resource: %w", err)
	}

	return nil
}

// loadSchema loads the schema with the given name. If the schema has already been loaded,
// it is returned from the cache.
func loadSchema(schemaName string) (*jsonschema.Schema, error) {
	schema, ok := schemaMap[schemaName]
	if ok {
		return schema, nil
	}

	schemaPath := schemaDir + schemaName + schemaExtension
	schemaFile, err := embeddedSchemas.Open(schemaPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load schema file: %w", err)
	}

	schemaURL := schemaHost + schemaPath
	err = schemaCompiler.AddResource(schemaURL, schemaFile)
	if err != nil {
		return nil, fmt.Errorf("failed to add schema as resource: %w", err)
	}

	schema, err = schemaCompiler.Compile(schemaURL)
	if err != nil {
		return nil, fmt.Errorf("failed to compile schema: %w", err)
	}

	schemaMap[schemaName] = schema

	return schema, nil
}
