package jsontabwriter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"text/tabwriter"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table"
	"github.com/jmespath/go-jmespath"
	"github.com/spf13/viper"
)

const (
	outputFormatErr = "invalid format: %s"

	// JSONFormat defines the legacy JSON format. It will eventually be removed and fully replaced by APIFormat
	// (in terms of behavior, the name will be kept)
	JSONFormat = "json"
	// TextFormat defines a human-readable format.
	TextFormat = "text"
	// APIFormat defines the API matching JSON format. This will be removed once its behavior will be moved to JSONFormat
	APIFormat = "api-json"
)

// GenerateOutput converts and formats source data into printable output.
//
// columnPathMappingPrefix: Points to a specific location in the JSON or struct object from where
// extraction should begin. If left empty, it will start from the root of the object.
//
// columnPathMapping: A map where each key represents the desired column name in the output
// table, and each value represents a JSON path to extract the data within the given JSON or struct.
//
// sourceData: JSON or struct from which data should be extracted, converted and formatted.
//
// cols: The columns that need to be printed
//
// Returns a ready-to-print string, which has the source data in either human-readable/table or JSON format
//
// TODO: remove cols as function parameter once --cols flag fix is ready
func GenerateOutput(
	columnPathMappingPrefix string, columnPathMapping map[string]string, sourceData interface{}, cols []string,
) (string, error) {
	if viper.IsSet(constants.ArgQuiet) {
		return "", nil
	}

	outputFormat := viper.GetString(constants.ArgOutput)

	// Generate the output in the requested format first.
	var formatted interface{}
	var err error
	switch outputFormat {
	case APIFormat:
		// For api-json, keep sourceData as-is (will be marshaled after filter).
		formatted = sourceData
	case TextFormat:
		if viper.IsSet(constants.FlagQuery) {
			return "", fmt.Errorf("JMESPath filtering (--query) is not supported with text output. Use -o api-json or json format instead")
		}
		// Text format produces a string directly; no filtering afterward.
		return generateTextOutputFromJSON(columnPathMappingPrefix, sourceData, columnPathMapping, cols)
	case JSONFormat:
		// Generate the legacy JSON structure (wraps into { "items": [...] }).
		formatted, err = generateLegacyJSONOutputAsData(sourceData)
		if err != nil {
			return "", err
		}
	default:
		return "", fmt.Errorf(outputFormatErr, outputFormat)
	}

	// Apply JMESPath filter if --query is set.
	if viper.IsSet(constants.FlagQuery) {
		expr := viper.GetString(constants.FlagQuery)
		if expr != "" {
			filtered, err := applyJMESPathFilter(formatted, expr)
			if err != nil {
				return "", fmt.Errorf("failed applying filter %q: %w", expr, err)
			}
			formatted = filtered
		}
	}

	// Marshal the (possibly filtered) data to JSON string.
	return generateJSONOutputAPI(formatted)
}

// GenerateOutputPreconverted is just like GenerateOutput, but it assumes that the source data has already been converted
// from JSON to table format. It is recommended when certain table columns cannot be automatically extracted from source
// data and require to be manually populated.
//
// rawSourceData: JSON or struct which will be used for JSON formatted output
//
// convertedSourceData: Table which will be used for human-readable output
//
// cols: The columns that need to be printed
//
// Returns a ready-to-print string, which has the source data in either human-readable/table or JSON format
func GenerateOutputPreconverted(
	rawSourceData interface{}, convertedSourceData []map[string]interface{}, cols []string,
) (string, error) {
	if viper.IsSet(constants.ArgQuiet) {
		return "", nil
	}

	outputFormat := viper.GetString(constants.ArgOutput)

	var formatted interface{}
	var err error
	switch outputFormat {
	case APIFormat:
		formatted = rawSourceData
	case TextFormat:
		return writeTableToText(convertedSourceData, cols), nil
	case JSONFormat:
		formatted, err = generateLegacyJSONOutputAsData(rawSourceData)
		if err != nil {
			return "", err
		}
	default:
		return "", fmt.Errorf(outputFormatErr, outputFormat)
	}

	if viper.IsSet(constants.FlagQuery) {
		expr := viper.GetString(constants.FlagQuery)
		if expr != "" {
			filtered, err := applyJMESPathFilter(formatted, expr)
			if err != nil {
				return "", fmt.Errorf("failed applying filter %q: %w", expr, err)
			}
			formatted = filtered
		}
	}

	return generateJSONOutputAPI(formatted)
}

// generateLegacyJSONOutputAsData converts source data into legacy JSON output structure.
// (generally used for list --all and other merged outputs, in combination with -o json)
//
//   - If source data is a slice of objects that each contain an "items" field which is a slice,
//     then all those "items" slices are merged into a single "items" slice in the output object.
//   - Otherwise, if source data is a slice of objects, any object that contains a "properties" field
//     is excluded, and the remaining objects are included in an "items" slice in the output object.
//   - If source data is not a slice, it is returned as-is.
func generateLegacyJSONOutputAsData(sourceData interface{}) (interface{}, error) {
	raw, err := json.Marshal(sourceData)
	if err != nil {
		return nil, fmt.Errorf("failed converting source data to JSON: %w", err)
	}

	var temp interface{}
	if err := json.Unmarshal(raw, &temp); err != nil {
		return nil, fmt.Errorf("unmarshal source data for legacy JSON output: %w", err)
	}

	slice, ok := temp.([]interface{})
	if !ok {
		return temp, nil
	}

	merged := make([]interface{}, 0)
	foundItemsSlice := false
	for _, elem := range slice {
		m, isMap := elem.(map[string]interface{})
		if !isMap {
			continue
		}
		itemsVal, hasItems := m["items"]
		if !hasItems {
			continue
		}
		itemsSlice, isSlice := itemsVal.([]interface{})
		if !isSlice {
			continue
		}
		foundItemsSlice = true
		merged = append(merged, itemsSlice...)
	}

	if foundItemsSlice {
		return map[string]interface{}{"items": merged}, nil
	}

	fallback := make([]interface{}, 0, len(slice))
	for _, elem := range slice {
		m, isMap := elem.(map[string]interface{})
		if isMap {
			if _, hasProps := m["properties"]; hasProps {
				continue
			}
		}
		fallback = append(fallback, elem)
	}

	return map[string]interface{}{"items": fallback}, nil
}

// applyJMESPathFilter compiles and applies a JMESPath expression to source.
// source may be an SDK struct or slice; it is marshaled to JSON and unmarshaled
// into interface{} so jmespath can operate on maps/slices/primitives.
func applyJMESPathFilter(source interface{}, expr string) (interface{}, error) {
	// convert structs -> JSON-compatible map/slice representation
	b, err := json.Marshal(source)
	if err != nil {
		return nil, fmt.Errorf("marshal source for filter: %w", err)
	}

	var data interface{}
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, fmt.Errorf("unmarshal source for filter: %w", err)
	}

	res, err := jmespath.Search(expr, data)
	if err != nil {
		return nil, fmt.Errorf("search jmespath expression: %w", err)
	}

	// return the filtered result (may be nil, scalar, object, or array)
	return res, nil
}

func GenerateVerboseOutput(format string, a ...interface{}) string {
	if viper.IsSet(constants.ArgQuiet) {
		return ""
	}

	if viper.GetInt(constants.ArgVerbose) == 0 {
		return ""
	}

	msg := fmt.Sprintf("[INFO] "+format, a...)

	if viper.GetString(constants.ArgOutput) == JSONFormat {
		out, _ := json.MarshalIndent(map[string]string{"Message": msg}, "", "  ")

		return string(out)
	}

	return GenerateRawOutput(msg)
}

// GenerateLogOutput is similar to fmt.Sprintf, but it will return the string in either JSON or text format.
func GenerateLogOutput(format string, a ...interface{}) string {
	msg := fmt.Sprintf(format, a...)

	return GenerateRawOutput(msg)
}

// GenerateRawOutput takes a generic variable and converts it to plain JSON or human-readable output without additional
// formatting fluff.
func GenerateRawOutput(a interface{}) string {
	if viper.IsSet(constants.ArgQuiet) {
		return ""
	}

	switch viper.GetString(constants.ArgOutput) {
	case APIFormat, JSONFormat:
		// Since generateJSONOutputAPI will only error out if an unsupported JSON type (e.g. chan, function values,
		// complex numbers or cyclic structs) or value (e.g. math.Inf()), which are not typically used in the API/SDKs,
		// I believe this error can be completely ignored in this use case.
		out, _ := generateJSONOutputAPI(a)

		return out
	case TextFormat:
		return fmt.Sprintf("%v\n", a)
	default:
		return ""
	}
}

// generateJSONOutputAPI marshals source data into JSON format, with indent.
func generateJSONOutputAPI(sourceData interface{}) (string, error) {
	out, err := json.MarshalIndent(sourceData, "", "  ")
	if err != nil {
		return "", err
	}

	return string(out) + "\n", nil
}

// generateTextOutputFromJSON converts JSON/struct object into human-readable format
//
// columnPathMappingPrefix: Points to a specific location in the JSON or struct object from where
// extraction should begin. If left empty, it will start from the root of the object.
//
// columnPathMapping: A map where each key represents the desired column name in the output
// table, and each value represents a JSON path to extract the data within the given JSON or struct.
//
// sourceData: JSON or struct from which data should be extracted, converted and formatted.
//
// cols: The columns that need to be in the table.
func generateTextOutputFromJSON(
	columnPathMappingPrefix string, sourceData interface{}, columnPathMapping map[string]string, cols []string,
) (string, error) {
	table, err := json2table.ConvertJSONToTable(columnPathMappingPrefix, columnPathMapping, sourceData)
	if err != nil {
		return "", fmt.Errorf("failed converting source data to table using %+v: %w", columnPathMapping, err)
	}

	return writeTableToText(table, cols), nil
}

// writeTableToText converts the tabled data (column-value associations) into an actual text table.
//
// table: Each map represents a row in the table, with each key-value pair in the map being equivalent to a
// column name and its value.
//
// cols: The columns that need to be in the table.
func writeTableToText(table []map[string]interface{}, cols []string) string {
	var buff = new(bytes.Buffer)
	var w = new(tabwriter.Writer)

	w.Init(buff, 5, 0, 3, ' ', tabwriter.StripEscape)

	updatedCols := eliminateEmptyCols(cols, table)
	if updatedCols == nil {
		return ""
	}

	if !viper.GetBool(constants.ArgNoHeaders) {
		fmt.Fprintln(w, strings.Join(updatedCols, "\t"))
	}

	for _, t := range table {
		format, values := convertTableToText(updatedCols, t)
		fmt.Fprintf(w, format+"\n", values...)
	}

	w.Flush()

	return buff.String()
}

// convertTableToText generates the basic formats string and corresponding values slice, based on the tabled data.
//
// cols: The columns that need to be in the table.
//
// table: Each map represents a row in the table, with each key-value pair in the map being equivalent to a
// column name and its value.
func convertTableToText(cols []string, table map[string]interface{}) (formats string, values []interface{}) {
	formatsSlice := make([]string, 0)
	values = make([]interface{}, 0)

	for _, col := range cols {
		field := table[col]
		formatsSlice = append(formatsSlice, "%v")

		if field == nil {
			values = append(values, "")

			continue
		}

		switch fieldType := reflect.TypeOf(field).Kind(); fieldType {
		case reflect.Slice:
			temp := make([]string, 0)
			for i := 0; i < reflect.ValueOf(field).Len(); i++ {
				temp = append(temp, fmt.Sprintf("%v", reflect.ValueOf(field).Index(i)))
			}

			field = strings.Join(temp, ", ")
		case reflect.Float64:
			if field == float64(int64(field.(float64))) {
				field = int64(field.(float64))
			}
		default:
		}

		values = append(values, field)
	}

	return strings.Join(formatsSlice, "\t"), values
}

// eliminateEmptyCols filters the columns so that there will be no empty columns in the final table.
//
// cols: The columns that need to be in the table.
//
// table: Each map represents a row in the table, with each key-value pair in the map being equivalent to a
// column name and its value.
func eliminateEmptyCols(cols []string, table []map[string]interface{}) []string {
	var newCols []string

	for _, c := range cols {
		for _, elem := range table {
			e, ok := elem[c]
			if !ok || e == nil {
				continue
			}

			if reflect.TypeOf(e).Kind() == reflect.Slice && reflect.ValueOf(e).Len() == 0 {
				continue
			}

			if s, ok := e.(string); ok && s == "" {
				continue
			}

			newCols = append(newCols, c)
			break
		}
	}

	return newCols
}
