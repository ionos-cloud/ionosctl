package printer

import (
	"encoding/json"
	"fmt"
	"golang.org/x/exp/slices"
	"io"
	"regexp"
	"strings"
	"text/tabwriter"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/services/cloudapi-v6/resources"
)

func ColsMessage(cols []string) string {
	return fmt.Sprintf("Set of columns to be printed on output \nAvailable columns: %v", cols)
}

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

// GetHeadersAllDefault is like GetHeaders, but defaultColumns is same as allColumns.
// Useful for resources with small print table
func GetHeadersAllDefault(allColumns []string, customColumns []string) []string {
	return GetHeaders(allColumns, allColumns, customColumns)
}

// GetHeaders takes all columns of a resource and the value of the columns flag,
// returns the headers of the table. (Some legacy code might refer to these headers as "Columns")
//
// allColumns can be found by using structs.Names on a Print struct (i.e. structs.Names(DatacenterPrint{}))
func GetHeaders(allColumns []string, defaultColumns []string, customColumns []string) []string {
	if customColumns == nil {
		return defaultColumns
	}

	var validCustomColumns []string
	for _, c := range customColumns {
		if slices.Contains(allColumns, c) {
			validCustomColumns = append(validCustomColumns, c)
		}
	}

	if len(validCustomColumns) == 0 {
		return defaultColumns
	}

	return validCustomColumns
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
		requestId, err := GetRequestId(GetRequestPath(prt.ApiResponse))
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
		requestId, err := GetRequestId(GetRequestPath(prt.ApiResponse))
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
	Output    interface{} `json:"Resources,omitempty"`
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

			values = append(values, v)

			switch v.(type) {
			case string:
				formats = append(formats, "%s")
			case int:
				formats = append(formats, "%d")
			case float64:
				formats = append(formats, "%f")
			case bool:
				formats = append(formats, "%v")
			default:
				formats = append(formats, "%v")
			}
		}
		format := strings.Join(formats, "\t")
		_, err := fmt.Fprintf(w, format+"\n", values...)
		if err != nil {
			return err
		}
	}

	return w.Flush()
}

var requestPathRegex = regexp.MustCompile(`https?://[a-zA-Z0-9./-]+/requests/([a-z0-9-]+)/status`)

func GetRequestId(path string) (string, error) {
	if !requestPathRegex.MatchString(path) {
		return "", fmt.Errorf("%s does not contain requestId", path)
	}
	return requestPathRegex.FindStringSubmatch(path)[1], nil
}

func GetId(r *resources.Response) string {
	if id, err := GetRequestId(GetRequestPath(r)); err == nil {
		return id
	}
	return ""
}

func GetRequestPath(r *resources.Response) string {
	if r != nil && r.Header != nil && len(r.Header) > 0 {
		return r.Header.Get("location")
	}
	return ""
}
