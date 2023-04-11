package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"

	"github.com/getkin/kin-openapi/openapi3"
	"golang.org/x/exp/slices"
)

func main() {
	openAPIFile := flag.String("s", "dataplatform.yaml", "Path to the OpenAPI spec file (required)")
	operationID := flag.String("i", "", "Operation ID to generate the CLI command for (required or -l)")
	listOperationIDs := flag.Bool("l", false, "List all operation IDs instead of generating a command (required or -i)")
	genHelpers := flag.Bool("p", false, "Generate the duplicated printing helper functions: Make(obj)Print, allCols, PrintObject, etc")

	out := flag.String("o", "generated_command", "Output file")
	flag.Parse()

	if *openAPIFile == "" {
		flag.Usage()
		os.Exit(1)
	}

	swagger, err := openapi3.NewLoader().LoadFromFile(*openAPIFile)
	if err != nil {
		panic(err)
	}

	if *listOperationIDs {
		printOperationIDs(swagger)
		return
	}

	if *operationID == "" {
		flag.Usage()
		os.Exit(1)
	}

	operation, _, method, err := findOperation(swagger, *operationID)
	if err != nil {
		panic(err)
	}

	flags, err := extractFlags(operation)
	if err != nil {
		panic(err)
	}

	cols, err := extractCols(operation)
	if err != nil {
		panic(err)
	}

	command := CLICommand{
		FunctionName: pascalCase(*operationID),
		Namespace:    strings.TrimSuffix(strings.TrimSuffix(filepath.Base(*openAPIFile), ".yaml"), ".json"),
		Resource:     strings.TrimSuffix(extractResource(*operationID), "s"),
		Verb:         strings.ToLower(method),
		Aliases:      createAliases(method),
		ShortDesc:    operation.Summary,
		InitClient:   "true",
		Flags:        flags,
		Columns:      cols,
	}

	if command.Verb == "post" {
		command.Verb = "create"
	}
	if command.Verb == "patch" || command.Verb == "put" {
		command.Verb = "update"
	}
	if command.Verb == "delete" {
		command.Verb = "remove"
	}
	if command.Verb == "get" {
		command.Verb = "list OR get // TODO"
	}
	command.Example = fmt.Sprintf("ionosctl %s %s %s\",// TODO: Add required flags or improve gen script", command.Namespace, command.Resource, command.Verb)

	if *genHelpers {
		tmplHelpers, err := template.New("cli-helpers").Parse(printHelpersTmpl)
		if err != nil {
			panic(err)
		}
		var bufHelpers bytes.Buffer
		err = tmplHelpers.Execute(&bufHelpers, command)
		if err != nil {
			panic(err)
		}
		command.Helpers = bufHelpers.String()
	}

	var templateFunctions = template.FuncMap{
		"requiredFlagsExample": requiredFlagsExample,
	}

	tmpl, err := template.New("cli-command").Funcs(templateFunctions).Parse(cliCommandTemplate)
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, command)
	if err != nil {
		panic(err)
	}

	// Write the output to a file or print it
	f, err := os.Create(fmt.Sprintf("%s.go", *out))
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = f.Write(buf.Bytes())
	if err != nil {
		panic(err)
	}
}

// Func called if `-l` set
func printOperationIDs(swagger *openapi3.T) {
	for _, pathItem := range swagger.Paths {
		for _, operation := range pathItem.Operations() {
			fmt.Println(operation.OperationID)
		}
	}
}

type CLICommand struct {
	FunctionName string
	Namespace    string
	Resource     string
	Verb         string
	Aliases      string
	ShortDesc    string
	Example      string
	InitClient   string
	Flags        []Flag
	Helpers      string // MakePrint, allCols and other duplicated code related to printing.
	Columns      []Col  // Columns for the response object
}

type Flag struct {
	Name        string
	ShortName   string
	Type        string
	Default     string
	Description string
	Required    bool
}

type Col struct {
	Name string
	Type string
}

func parseFlagDescription(desc string) string {
	return strings.TrimSuffix(strings.ReplaceAll(desc, "\n", ""), ".")
}

func findOperation(swagger *openapi3.T, operationID string) (*openapi3.Operation, string, string, error) {
	for path, pathItem := range swagger.Paths {
		for method, operation := range pathItem.Operations() {
			if operation.OperationID == operationID {
				return operation, path, method, nil
			}
		}
	}
	return nil, "", "", fmt.Errorf("operation not found: %s", operationID)
}

func extractCols(operation *openapi3.Operation) ([]Col, error) {
	var cols []Col

	response := operation.Responses.Get(200)
	if response == nil {
		return nil, fmt.Errorf("no 200 response found")
	}

	content := response.Value.Content["application/json"]
	if content == nil {
		return nil, fmt.Errorf("no 'application/json' content found")
	}

	schema := content.Schema
	if schema == nil {
		return nil, fmt.Errorf("no schema found")
	}

	cols, err := extractNestedCols(schema.Value.Properties, "")
	if err != nil {
		return nil, err
	}

	return cols, nil
}

func extractNestedCols(properties map[string]*openapi3.SchemaRef, parentName string) ([]Col, error) {
	var cols []Col

	for propName, prop := range properties {
		fullPropName := propName
		if parentName != "" {
			fullPropName = strings.Title(propName)
		}

		if prop.Value.Type == "object" {
			nestedCols, err := extractNestedCols(prop.Value.Properties, fullPropName)
			if err != nil {
				return nil, err
			}
			cols = append(cols, nestedCols...)
		} else {
			col := Col{
				Name: strings.Title(fullPropName),
				Type: goTypeFromSchema(prop.Value.Type),
			}
			cols = append(cols, col)
		}
	}

	return cols, nil
}

func extractFlags(operation *openapi3.Operation) ([]Flag, error) {
	var flags []Flag

	j, _ := operation.MarshalJSON()
	log.Printf("Operation: %s\n", j)
	// Iterate through operation parameters
	for _, paramRef := range operation.Parameters {
		param := paramRef.Value

		// Extract parameter properties
		flag := Flag{
			Name:        param.Name,
			ShortName:   "", // You can provide a custom mapping for short names or leave it empty
			Type:        flagTypeFromSchema(param.Schema.Value.Type),
			Default:     flagDefaultFromSchema(param.Schema.Value),
			Description: parseFlagDescription(param.Description),
			Required:    param.Required,
		}
		flags = append(flags, flag)
	}

	// Iterate through request body properties
	if operation.RequestBody != nil {
		content := operation.RequestBody.Value.Content["application/json"]
		if content != nil {
			for propName, prop := range content.Schema.Value.Properties {
				flags = appendNestedFlags(flags, propName, prop, content)
			}
		}
	}

	return flags, nil
}

func appendNestedFlags(flags []Flag, propName string, prop *openapi3.SchemaRef, content *openapi3.MediaType) []Flag {
	if prop.Value.Properties == nil {
		flag := Flag{
			Name:        "constants.Flag" + strings.Title(propName),
			ShortName:   "",
			Type:        flagTypeFromSchema(prop.Value.Type),
			Default:     flagDefaultFromSchema(prop.Value),
			Description: parseFlagDescription(prop.Value.Description),
			Required:    slices.Contains(prop.Value.Required, propName),
		}
		flags = append(flags, flag)
	} else {
		for nestedPropName, nestedProp := range prop.Value.Properties {
			flags = appendNestedFlags(flags, nestedPropName, nestedProp, content)
		}
	}

	return flags
}

func flagTypeFromSchema(schema string) string {
	switch schema {
	case "integer":
		return "Int"
	case "number":
		return "Float64"
	case "string":
		return "String"
	case "boolean":
		return "Bool"
	case "array":
		return "StringSlice"
	case "object":
		return "StringToString"
	default:
		return "String" // Default to string for unknown types
	}
}

func goTypeFromSchema(schema string) string {
	switch schema {
	case "integer":
		return "int"
	case "number":
		return "float"
	case "string":
		return "string"
	case "boolean":
		return "bool"
	case "array":
		return "[]string"
	case "object":
		return "interface{}"
	default:
		return "string" // Default to string for unknown types
	}
}
func flagDefaultFromSchema(schema *openapi3.Schema) string {
	switch schema.Type {
	case "integer":
		return "0"
	case "number":
		return "0.0"
	case "string":
		return `""`
	case "boolean":
		return "false"
	case "array":
		return "[]string{}"
	case "object":
		return "map[string]string{}"
	default:
		return `""`
	}
}

// used to get resource name from operationId, e.g. clustersKubeconfigFindByClusterId => kubeconfig
func extractResource(operationID string) string {
	parts := splitOnCapitalLetters(operationID)

	if len(parts) == 2 {
		return strings.ToLower(parts[0])
	}

	return strings.ToLower(parts[1])
}

func createAliases(verb string) string {
	verb = strings.ToLower(verb)
	switch verb {
	case "delete":
		return "{\"del\", \"d\"}"
	case "list":
		return "{\"ls\"}"
	case "update":
		return "{\"u\", \"patch\"}"
	case "get":
		return "{\"g\"}"
	case "create":
		return "{\"c\", \"post\"}"
	default:
		return "{}"
	}
}

func splitOnCapitalLetters(s string) []string {
	var words []string
	start := 0
	for i := 1; i < len(s); i++ {
		if unicode.IsUpper(rune(s[i])) {
			words = append(words, s[start:i])
			start = i
		}
	}
	words = append(words, s[start:])
	return words
}

func pascalCase(s string) string {
	words := strings.Split(s, "_")
	for i, word := range words {
		words[i] = strings.Title(word)
	}
	return strings.Join(words, "")
}

const cliCommandTemplate = `package main

import (
	"context"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
)

func {{.FunctionName}}Cmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "{{.Namespace}}",
		Resource:  "{{.Resource}}",
		Verb:      "{{.Verb}}",
		Aliases:   []string{{.Aliases}},
		ShortDesc: "{{.ShortDesc}}",
		Example:   "ionosctl {{.Namespace}} {{.Resource}} {{.Verb}} {{requiredFlagsExample .Flags}}",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			/* TODO: Delete/modify me for --all
			 * err := core.CheckRequiredFlagsSets(c.Command, c.NS, []string{constants.ArgAll}, []string{constants.Flag<Parent>Id}, []string{constants.ArgAll, constants.Flag<Parent>Id})
			 * if err != nil {
			 * 	return err
			 * }
             * */

			// TODO: If no --all, mark individual flags as required{{range .Flags}}{{if .Required}}
			err = c.Command.Command.MarkFlagRequired("{{.Name}}")
			if err != nil {
				return err
			}{{end}}{{end}}

			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			// Implement the actual command logic here
		},
		InitClient: {{.InitClient}},
	})

	{{range .Flags}}
	cmd.Add{{.Type}}Flag({{.Name}}, "{{.ShortName}}", {{.Default}}, "{{.Description}}"{{if .Required}}, core.RequiredFlagOption(){{end}}){{end}}

	return cmd
}

{{.Helpers}}
`

const printHelpersTmpl = `// Helper functions for printing {{.Resource}}

func get{{.Resource}}Print(c *core.CommandConfig, dcs *[]ionoscloud.{{.Resource}}ResponseData) printer.Result {
	r := printer.Result{}
	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

	if c != nil && dcs != nil {
		r.OutputJSON = dcs
		r.KeyValue = make{{.Resource}}PrintObj(dcs)
		r.Columns = printer.GetHeaders(allCols, defCols, cols)
	}
	return r
}

type {{.Resource}}Print struct {
{{range .Columns}}	{{.Name}} {{.Type}} ` + "`json:\"{{.Name}},omitempty\"`" + `
{{end}}
}

var allCols = structs.Names({{.Resource}}Print{})
var defCols = allCols[:len(allCols)-3]

func make{{.Resource}}PrintObj(data *[]ionoscloud.{{.Resource}}ResponseData) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(*data))

	for _, item := range *data {
		var printObj {{.Resource}}Print
		printObj.Id = *item.GetId()

		// Fill in the rest of the fields from the response object
		{{range .Columns}}
		if propertiesOk, ok := item.GetPropertiesOk(); ok && propertiesOk != nil {
			printObj.{{.Name}} = *propertiesOk.Get{{.Name}}()
		}{{end}}

		o := structs.Map(printObj)
		out = append(out, o)
	}
	return out
}
`

func requiredFlagsExample(flags []Flag) string {
	var flagExample []string
	for _, flag := range flags {
		if flag.Required {
			flagExample = append(flagExample, fmt.Sprintf("--%s <%s>", flag.Name, flag.Type))
		}
	}
	return strings.Join(flagExample, " ")
}
