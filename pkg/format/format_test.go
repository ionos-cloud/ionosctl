package format

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
	"text/tabwriter"

	"github.com/Jeffail/gabs/v2"
	client2 "github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/spf13/viper"
)

var outputFormatErr = fmt.Errorf("invalid format")

type formatter struct {
	Json interface{}
	Text []map[string]interface{}
}

type HeaderJsonPath struct {
	Header string
	Path   string
}

var datacenter = []HeaderJsonPath{
	{"DatacenterId", "id"},
	{"Name", "properties.name"},
	{"Location", "properties.location"},
	{"CpuFamily", "properties.cpuArchitecture.*.cpuFamily"},
	{"State", "metadata.state"},
	{"Version", "properties.version"},
}

func TestGen(t *testing.T) {
	var cols = []string{"DatacenterId", "Name", "Location", "CpuFamily", "Version"}

	client, err := client2.NewTestClient(os.Getenv("IONOS_USERNAME"), os.Getenv("IONOS_PASSWORD"), "", "")
	if err != nil {
		fmt.Println(err)
		return
	}

	dcs, _, err := client.CloudClient.DataCentersApi.DatacentersFindById(context.Background(), "374a0987-d594-48e7-b599-e08e6cf95012").Depth(0).Execute()
	//dcs, _, err := client.CloudClient.DataCentersApi.DatacentersGet(context.Background()).Execute()
	if err != nil {
		fmt.Println(err)
		return
	}

	viper.Set(constants.ArgOutput, "text")
	out, _ := GenerateOutput(false, cols, "", dcs)
	fmt.Println(out)
}

func GenerateOutput(noHeaders bool, cols []string, rootPath string, obj interface{}) (string, error) {
	objToText, err := mapObjectToText(rootPath, datacenter, obj)
	if err != nil {
		return "", err
	}

	format := formatter{
		Json: obj,
		Text: objToText,
	}

	if viper.GetString(constants.ArgOutput) == "json" {
		return format.generateJSONOutput()
	}

	if viper.GetString(constants.ArgOutput) == "text" {
		return format.generateTextOutput(noHeaders, cols)
	}

	return "", outputFormatErr
}

func mapObjectToText(rootPath string, headerPaths []HeaderJsonPath, rootObj interface{}) ([]map[string]interface{}, error) {
	var res = make([]map[string]interface{}, 0)

	objs, err := traverseRoot(rootPath, rootObj)
	if err != nil {
		return nil, err
	}

	for _, obj := range objs {
		mappedObj := make(map[string]interface{}, 0)

		for _, iter := range headerPaths {
			if !obj.ExistsP(iter.Path) {
				return nil, fmt.Errorf("wrong path provided: %s", iter.Path)
			}

			objData := obj.Path(iter.Path)
			mappedObj[iter.Header] = objData.Data()
		}

		res = append(res, mappedObj)
	}

	return res, nil
}

func traverseRoot(rootPath string, obj interface{}) ([]*gabs.Container, error) {
	jsonObj, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	parsedObj, err := gabs.ParseJSON(jsonObj)
	if err != nil {
		return nil, err
	}

	if rootPath == "" {
		return []*gabs.Container{parsedObj}, nil
	}

	if !parsedObj.ExistsP(rootPath) {
		return nil, fmt.Errorf("root path does not exist in object: %s", rootPath)
	}

	parsedObj = parsedObj.Path(rootPath)
	children := parsedObj.Children()

	if children == nil {
		return nil, fmt.Errorf("root path does not lead to an array in object: %s", rootPath)
	}

	return children, nil
}

func (f formatter) generateJSONOutput() (string, error) {
	out, err := json.MarshalIndent(f.Json, "", "\t")
	if err != nil {
		return "", err
	}

	return string(out), nil
}

func (f formatter) generateTextOutput(noHeaders bool, cols []string) (string, error) {
	var buff = new(bytes.Buffer)
	var w = new(tabwriter.Writer)
	w.Init(buff, 5, 0, 3, ' ', tabwriter.StripEscape)

	if !noHeaders {
		_, err := fmt.Fprintln(w, strings.Join(cols, "\t"))
		if err != nil {
			return "", nil
		}
	}

	for _, obj := range f.Text {
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
