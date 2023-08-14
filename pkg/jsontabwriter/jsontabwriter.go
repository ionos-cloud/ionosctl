package jsontabwriter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/json2table"
	"github.com/spf13/viper"
)

var outputFormatErr = "invalid format: %s"

// TODO: remove cols as function parameter once --cols flag fix is ready
func GenerateOutput(rootPath string, jsonPaths map[string]string, obj interface{}, cols []string) (string, error) {
	outputFormat := viper.GetString(constants.ArgOutput)

	if outputFormat == "json" {
		return generateJSONOutput(obj)
	}

	if outputFormat == "text" {
		return generateTextOutput(rootPath, obj, jsonPaths, cols, false)
	}

	return "", fmt.Errorf(outputFormatErr, outputFormat)
}

func GenerateOutputPreconverted(obj interface{}, convertedObj []map[string]interface{}, cols []string) (string, error) {
	outputFormat := viper.GetString(constants.ArgOutput)

	if outputFormat == "json" {
		return generateJSONOutput(obj)
	}

	if outputFormat == "text" {
		return generateTextOutput("", convertedObj, nil, cols, true)
	}

	return "", fmt.Errorf(outputFormatErr, outputFormat)
}

func generateJSONOutput(obj interface{}) (string, error) {
	out, err := json.MarshalIndent(obj, "", "\t")
	if err != nil {
		return "", err
	}

	return string(out) + "\n", nil
}

func generateTextOutput(rootPath string, obj interface{}, jsonPaths map[string]string, cols []string, skipConversion bool) (string, error) {
	var table []map[string]interface{}
	var err error

	if skipConversion {
		table = obj.([]map[string]interface{})
	} else {
		table, err = json2table.ConvertJSONToTable(rootPath, jsonPaths, obj)
		if err != nil {
			return "", err
		}
	}

	var buff = new(bytes.Buffer)
	var w = new(tabwriter.Writer)
	w.Init(buff, 5, 0, 3, ' ', tabwriter.StripEscape)

	if !viper.IsSet(constants.ArgNoHeaders) {
		if _, err = fmt.Fprintln(w, strings.Join(cols, "\t")); err != nil {
			return "", nil
		}
	}

	for _, t := range table {
		format, values := convertTableToText(cols, t)
		if _, err = fmt.Fprintf(w, format+"\n", values...); err != nil {
			return "", err
		}
	}

	if err = w.Flush(); err != nil {
		return "", err
	}

	return buff.String(), nil
}

func GenerateVerboseOutput(format string, a ...interface{}) string {
	if !viper.IsSet(constants.ArgVerbose) {
		return ""
	}

	msg := fmt.Sprintf("[INFO] "+format, a...)

	return GenerateLogOutput(msg)
}

func GenerateLogOutput(format string, a ...interface{}) string {
	outputFormat := viper.GetString(constants.ArgOutput)
	msg := fmt.Sprintf(format, a...)

	if outputFormat == "json" {
		out, err := generateJSONOutput(msg)
		if err != nil {
			return ""
		}

		return out
	}

	if outputFormat == "text" {
		return msg + "\n"
	}

	return ""
}

func convertTableToText(cols []string, table map[string]interface{}) (format string, values []interface{}) {
	formats := make([]string, 0)
	values = make([]interface{}, 0)

	for _, col := range cols {
		field := table[col]
		formats = append(formats, "%v")

		if field == nil {
			values = append(values, "")

			continue
		}

		switch field.(type) {
		case []interface{}:
			temp := make([]string, 0)
			for _, val := range field.([]interface{}) {
				temp = append(temp, fmt.Sprintf("%v", val))
			}

			field = strings.Join(temp, ", ")
		case float64:
			temp := field.(float64)

			if temp == float64(int64(temp)) {
				field = int64(temp)
			}
		}

		values = append(values, field)
	}

	return strings.Join(formats, "\t"), values
}