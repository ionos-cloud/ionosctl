package jsontabwriter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
	"text/tabwriter"

	client2 "github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/json2table"
	"github.com/spf13/viper"
)

var outputFormatErr = fmt.Errorf("invalid format")

var datacenter = map[string]string{
	"DatacenterId": "id",
	"Name":         "properties.name",
	"Location":     "properties.location",
	"CpuFamily":    "properties.cpuArchitecture.*.cpuFamily",
	"State":        "metadata.state",
	"Version":      "properties.version",
}

func TestGen(t *testing.T) {
	client, err := client2.NewTestClient(os.Getenv("IONOS_USERNAME"), os.Getenv("IONOS_PASSWORD"), "", "")
	if err != nil {
		fmt.Println(err)
		return
	}

	//dcs, _, err := client.CloudClient.DataCentersApi.DatacentersFindById(context.Background(), "374a0987-d594-48e7-b599-e08e6cf95012").Depth(0).Execute()
	dcs, _, err := client.CloudClient.DataCentersApi.DatacentersGet(context.Background()).Execute()
	if err != nil {
		fmt.Println(err)
		return
	}

	viper.Set(constants.ArgOutput, "text")
	out, _ := GenerateOutput("items", dcs)
	fmt.Println(out)
}

func GenerateOutput(rootPath string, obj interface{}) (string, error) {
	if viper.GetString(constants.ArgOutput) == "json" {
		return generateJSONOutput(obj)
	}

	if viper.GetString(constants.ArgOutput) == "text" {
		return generateTextOutput(rootPath, obj, datacenter)
	}

	return "", outputFormatErr
}

func generateJSONOutput(obj interface{}) (string, error) {
	out, err := json.MarshalIndent(obj, "", "\t")
	if err != nil {
		return "", err
	}

	return string(out), nil
}

func generateTextOutput(rootPath string, obj interface{}, objJsonPaths map[string]string) (string, error) {
	cols := viper.GetStringSlice(constants.ArgCols)
	text, err := json2table.ConvertJSONToText(rootPath, objJsonPaths, obj)
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

	if viper.GetString(constants.ArgOutput) == "json" {
		return generateJSONVerboseOutput(format, a...)
	}

	if viper.GetString(constants.ArgOutput) == "text" {
		return generateTextVerboseOutput(format, a...)
	}

	return ""
}

func generateJSONVerboseOutput(format string, a ...interface{}) string {
	msg := fmt.Sprintf("[INFO] "+format, a...)

	out, err := json.MarshalIndent(map[string]string{
		"Message": msg,
	}, "", "\t")
	if err != nil {
		return ""
	}

	return fmt.Sprintln(string(out))
}

func generateTextVerboseOutput(format string, a ...interface{}) string {
	return fmt.Sprintf("[INFO] "+format+"\n", a...)
}
