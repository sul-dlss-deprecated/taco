package validators

import (
	"strings"

	"github.com/santhosh-tekuri/jsonschema"
)

// BuildSchema builds the schema for deposits from the encoded strings
func BuildSchema(uri string, files []string) *jsonschema.Schema {
	compiler := jsonschema.NewCompiler()
	for _, file := range files {
		if err := compiler.AddResource(file, strings.NewReader(Data[file])); err != nil {
			panic(err)
		}
	}

	schema, err := compiler.Compile(uri)
	if err != nil {
		panic(err)
	}
	return schema
}
