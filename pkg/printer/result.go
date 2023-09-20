package printer

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/v6/pkg/utils"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
)

type Result struct {
	Message        string
	Resource       string
	Verb           string
	WaitForRequest bool
	WaitForState   bool
	Columns        []string
	KeyValue       []map[string]interface{}
	OutputJSON     interface{}

	ApiResponse *resources.Response
}

// TODO: identical name to printText. Hard to decipher behaviour
func (prt *Result) PrintText(out io.Writer, noHeaders bool) error {
	var resultPrint ResultPrint
	if prt.Resource != "" && prt.Verb != "" {
		resultPrint.Message = standardSuccessMsg(prt.Resource, prt.Verb, prt.WaitForRequest, prt.WaitForState)
	} else if prt.Message != "" {
		resultPrint.Message = prt.Message
	}
	if prt.ApiResponse != nil {
		requestId, err := utils.GetRequestId(utils.GetRequestPath(prt.ApiResponse))
		if err != nil {
			return err
		}
		resultPrint.RequestId = requestId
	}
	if prt.KeyValue != nil && prt.Columns != nil {
		err := printText(out, prt.Columns, prt.KeyValue, noHeaders)
		if err != nil {
			return err
		}
	}
	if resultPrint.RequestId != nil {
		requestIdMsg(out, "%v", resultPrint.RequestId)
	}
	if resultPrint.Message != nil {
		statusMsg(out, "%v", resultPrint.Message)
	}
	return nil
}

func (prt *Result) PrintJSON(out io.Writer) error {
	var resultPrint ResultPrint
	if prt.Resource != "" && prt.Verb != "" {
		resultPrint.Message = standardSuccessMsg(prt.Resource, prt.Verb, prt.WaitForRequest, prt.WaitForState)
	} else if prt.Message != "" {
		resultPrint.Message = prt.Message
	}
	if prt.ApiResponse != nil {
		requestId, err := utils.GetRequestId(utils.GetRequestPath(prt.ApiResponse))
		if err != nil {
			return err
		}
		resultPrint.RequestId = requestId
	}
	resultPrint.Output = prt.OutputJSON
	if !structs.IsZero(resultPrint) {
		j, err := json.MarshalIndent(&resultPrint, "", "  ")
		if err != nil {
			return err
		}
		_, err = fmt.Fprintf(out, "%s\n", string(j))
		if err != nil {
			return err
		}
	}
	return nil
}

type ResultPrint struct {
	Message   interface{} `json:"Status,omitempty"`
	RequestId interface{} `json:"RequestId,omitempty"`
	Output    interface{} `json:"items,omitempty"`
}

var (
	standardSuccessMessages     = "Command %s %s has been successfully executed"
	waitStandardSuccessMessages = "Command %s %s & wait have been successfully executed"
)

func standardSuccessMsg(resource, verb string, waitRequest, waitState bool) string {
	if waitRequest || waitState {
		return fmt.Sprintf(waitStandardSuccessMessages, resource, verb)
	}
	return fmt.Sprintf(standardSuccessMessages, resource, verb)
}

func requestIdMsg(writer io.Writer, msg string, args ...interface{}) {
	_, err := fmt.Fprintf(writer, "RequestId: %s\n", fmt.Sprintf(msg, args...))
	if err != nil {
		return
	}
}

func statusMsg(writer io.Writer, msg string, args ...interface{}) {
	_, err := fmt.Fprintf(writer, "Status: %s\n", fmt.Sprintf(msg, args...))
	if err != nil {
		return
	}
}

func printText(out io.Writer, cols []string, keyValueMap []map[string]interface{}, noHeaders bool) error {
	w := new(tabwriter.Writer)
	w.Init(out, 5, 0, 3, ' ', tabwriter.StripEscape)

	if !noHeaders {
		var headers []string
		for _, col := range cols {
			headers = append(headers, col)
		}
		_, err := fmt.Fprintln(w, strings.Join(headers, "\t"))
		if err != nil {
			return err
		}
	}

	for _, r := range keyValueMap {
		var values []interface{}
		var formats []string

		for _, col := range cols {
			v := r[col]

			switch v.(type) {
			case string:
				formats = append(formats, "%s")
			case int:
				formats = append(formats, "%d")
			case float64:
				formats = append(formats, "%f")
			case bool:
				formats = append(formats, "%v")
			case []string:
				formats = append(formats, "%s")
				v = strings.Join(v.([]string), ",")
			default:
				formats = append(formats, "%v")
			}

			values = append(values, v)

		}
		format := strings.Join(formats, "\t")
		_, err := fmt.Fprintf(w, format+"\n", values...)
		if err != nil {
			return err
		}
	}

	return w.Flush()
}
