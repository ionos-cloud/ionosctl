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
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/lastresponse"
	"github.com/itchyny/gojq"
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

var (
	LastSrcData = map[interface{}]interface{}(nil)
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
	err := lastresponse.SetE(sourceData)
	if err != nil {
		return "", err
	}

	if viper.IsSet(constants.ArgQuiet) {
		return "", nil
	}

	outputFormat := viper.GetString(constants.ArgOutput)
	switch outputFormat {
	case APIFormat:
		return generateJSONOutputAPI(sourceData)
	case TextFormat:
		return generateTextOutputFromJSON(columnPathMappingPrefix, sourceData, columnPathMapping, cols)
	case JSONFormat:
		return generateLegacyJSONOutput(sourceData)
	default:
		return "", fmt.Errorf(outputFormatErr, outputFormat)
	}
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
	switch outputFormat {
	case APIFormat:
		return generateJSONOutputAPI(rawSourceData)
	case TextFormat:
		return writeTableToText(convertedSourceData, cols), nil
	case JSONFormat:
		return generateLegacyJSONOutput(rawSourceData)
	default:
		return "", fmt.Errorf(outputFormatErr, outputFormat)
	}
}

func GenerateVerboseOutput(format string, a ...interface{}) string {
	if viper.IsSet(constants.ArgQuiet) {
		return ""
	}

	if !viper.IsSet(constants.ArgVerbose) {
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

// generateLegacyJSONOutput modifies the source data so that the output generated still respects the legacy JSON format.
func generateLegacyJSONOutput(sourceData interface{}) (string, error) {
	apiOut, err := json.MarshalIndent(sourceData, "", "  ")
	if err != nil {
		return "", err
	}

	var temp interface{}
	err = json.Unmarshal(apiOut, &temp)
	if err != nil {
		return "", err
	}

	query1, err := gojq.Parse("{ items: [.[] | .items] | add }")
	if err != nil {
		return "", err
	}

	query2, err := gojq.Parse(`map(select(has("properties") | not)) | { "items": . }`)
	if err != nil {
		return "", err
	}

	// I expect only one result from the query, so there is no need to loop through the results
	queryResult, _ := query1.Run(temp).Next()
	if err, ok := queryResult.(error); ok && err != nil {
		return string(apiOut) + "\n", nil
	}

	// fixes null output for embedded objects
	mappedQueryResult := queryResult.(map[string]any)
	if mappedQueryResult["items"] == nil {
		queryResult, _ = query2.Run(temp).Next()
		if err, ok := queryResult.(error); ok && err != nil {
			fmt.Println(err.Error())

			return string(apiOut) + "\n", nil
		}
	}

	return generateJSONOutputAPI(queryResult)
}
