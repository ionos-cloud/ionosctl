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

	"github.com/fatih/structs"
	client2 "github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/viper"
)

var outputFormatErr = fmt.Errorf("invalid format")

type Formatter struct {
	Json interface{}
	Text []interface{}
}

func TestGen(t *testing.T) {
	//var test interface{}
	var cols = []string{"DatacenterId", "Name", "Location", "CpuFamily", "State"}

	client, err := client2.NewTestClient(os.Getenv("IONOS_USERNAME"), os.Getenv("IONOS_PASSWORD"), "", "")
	if err != nil {
		fmt.Println(err)
		return
	}

	//dcs, _, err := client.CloudClient.DataCentersApi.DatacentersFindById(context.Background(), "374a0987-d594-48e7-b599-e08e6cf95012").Depth(0).Execute()
	dcs, _, err := client.CloudClient.DataCentersApi.DatacentersGet(context.Background()).Execute()
	//dcs := []interface{}{"create", "delete", "whatever"}
	if err != nil {
		fmt.Println(err)
		return
	}

	obj := Formatter{
		Json: dcs,
		//Text: mapToTextOutputStruct(dcs),
	}

	viper.Set(constants.ArgOutput, "json")
	out, _ := obj.GenerateOutput(false, cols)
	fmt.Println(out)
}

func (f Formatter) GenerateOutput(noHeaders bool, cols []string) (string, error) {
	if viper.GetString(constants.ArgOutput) == "json" {
		return f.generateJSONOutput()
	}

	if viper.GetString(constants.ArgOutput) == "text" {
		return f.generateTextOutput(noHeaders, cols)
	}

	return "", outputFormatErr
}

type simpleDC struct {
	DatacenterId string `json:"DatacenterId,omitempty"`
	Name         string `json:"Name,omitempty"`
	Location     string `json:"Location,omitempty"`
	CpuFamily    string `json:"CpuFamily,omitempty"`
	State        string `json:"State,omitempty"`
}

func mapToTextOutputStruct(dc ionoscloud.Datacenter) []interface{} {
	texts := make([]interface{}, 0)
	text := simpleDC{
		DatacenterId: *dc.Id,
		Name:         *dc.Properties.Name,
		Location:     *dc.Properties.Location,
		CpuFamily:    *(*dc.Properties.GetCpuArchitecture())[0].CpuFamily,
		State:        *dc.Metadata.State,
	}

	texts = append(texts, text)

	return texts
}

// NOTE: does it require a wrapper to not create breaking changes in users' scripts? To be used with raw, api returned structures
func (f Formatter) generateJSONOutput() (string, error) {
	out, err := json.MarshalIndent(f.Json, "", "\t")
	if err != nil {
		return "", err
	}

	return string(out), nil
}

// NOTE: to be used with user-defined structures (json output should return more info?)
func (f Formatter) generateTextOutput(noHeaders bool, cols []string) (string, error) {
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

		fields := structs.Map(obj)
		// Properties -> Name

		for _, col := range cols {
			field := fields[col]

			switch field.(type) {
			case []string:
				formats = append(formats, "%s")
				field = strings.Join(field.([]string), "\t")
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
