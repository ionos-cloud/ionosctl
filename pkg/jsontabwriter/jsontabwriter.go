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

func GenerateOutput(rootPath string, jsonPaths map[string]string, obj interface{}) (string, error) {
	outputFormat := viper.GetString(constants.ArgOutput)

	if outputFormat == "json" {
		return generateJSONOutput(obj)
	}

	if outputFormat == "text" {
		return generateTextOutput(rootPath, obj, jsonPaths)
	}

	return "", fmt.Errorf(outputFormatErr, outputFormat)
}

func generateJSONOutput(obj interface{}) (string, error) {
	out, err := json.MarshalIndent(obj, "", "\t")
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s\n", string(out)), nil
}

func generateTextOutput(rootPath string, obj interface{}, jsonPaths map[string]string) (string, error) {
	cols := viper.GetStringSlice(constants.ArgCols)
	text, err := json2table.ConvertJSONToText(rootPath, jsonPaths, obj)
	if err != nil {
		return "", err
	}

	var buff = new(bytes.Buffer)
	var w = new(tabwriter.Writer)
	w.Init(buff, 5, 0, 3, ' ', tabwriter.StripEscape)

	if !viper.IsSet(constants.ArgNoHeaders) {
		_, err := fmt.Fprintln(w, strings.Join(cols, "\t"))
		if err != nil {
			return "", nil
		}
	}

	for _, obj := range text {
		var formats []string
		var values []interface{}

		for _, col := range cols {
			field := obj[col]

			switch field.(type) {
			case []interface{}:
				formats = append(formats, "%s")

				temp := make([]string, 0)
				for _, val := range field.([]interface{}) {
					temp = append(temp, fmt.Sprintf("%v", val))
				}

				field = strings.Join(temp, ", ")
			default:
				formats = append(formats, "%v")
			}

			values = append(values, field)
		}

		format := strings.Join(formats, "\t")
		_, err := fmt.Fprintf(w, format+"\n", values...)
		if err != nil {
			return "", err
		}
	}

	if err := w.Flush(); err != nil {
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

func GenerateLogOutput(a interface{}) string {
	outputFormat := viper.GetString(constants.ArgOutput)

	if outputFormat == "json" {
		return generateJSONLogOutput(a)
	}

	if outputFormat == "text" {
		return generateTextLogOutput(a)
	}

	return ""
}

func generateJSONLogOutput(a interface{}) string {
	out, err := json.MarshalIndent(a, "", "\t")
	if err != nil {
		return ""
	}

	return fmt.Sprintf("%s\n", string(out))
}

func generateTextLogOutput(a interface{}) string {
	return fmt.Sprintf("%v\n", a)
}
