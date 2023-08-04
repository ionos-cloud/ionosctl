package printer

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"strings"
	"text/tabwriter"

	"github.com/ionos-cloud/ionosctl/v6/internal/functional"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"golang.org/x/exp/slices"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
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

// Filter case insensitively limits the Result's columns to those specified in the parameter
// and ignores column names from `cols` that are non-existent in r.Columns
// an empty `cols` will result in no filter being applied.
func (r *Result) Filter(customCols []string) {
	if customCols == nil || len(customCols) == 0 {
		return
	}

	allColumnsLowercase := functional.Map(r.Columns, func(col string) string {
		return strings.ToLower(col)
	})
	r.Columns = functional.Filter(customCols, func(customColumn string) bool {
		return slices.Contains(allColumnsLowercase, strings.ToLower(customColumn))
	})

	// Keep old behaviour of nil slices if empty
	if len(r.Columns) == 0 {
		r.Columns = nil
	}
}

// GetHeadersAllDefault is like GetHeaders, but defaultColumns is same as allColumns.
// Useful for resources with small print table
func GetHeadersAllDefault(allColumns []string, customColumns []string) []string {
	return GetHeaders(allColumns, allColumns, customColumns)
}

// GetHeaders takes all columns of a resource and the value of the columns flag,
// returns the headers of the table. (Some legacy code might refer to these headers as "Columns")
// allColumns can be found by using structs.Names on a Print struct (i.e. structs.Names(DatacenterPrint{}))
func GetHeaders(allColumns []string, defaultColumns []string, customColumns []string) []string {
	if len(customColumns) > 0 && customColumns[0] == "all" {
		return GetHeaders(allColumns, defaultColumns, allColumns)
	}

	if customColumns == nil {
		return defaultColumns
	}

	allColumnsLowercase := functional.Map(allColumns, func(x string) string {
		return strings.ToLower(x)
	})

	var validCustomColumns []string
	for _, c := range customColumns {
		if idx := slices.Index(allColumnsLowercase, strings.ToLower(c)); idx != -1 {
			validCustomColumns = append(validCustomColumns, allColumns[idx])
		}
	}

	if len(validCustomColumns) == 0 {
		return defaultColumns
	}

	return validCustomColumns
}

func GetHeadersListAll(allColumns []string, defaultColumns []string, parentCol string, customColumns []string, argAll bool) []string {
	if argAll {
		defaultColumns = append(defaultColumns[:constants.DefaultParentIndex+1], defaultColumns[constants.DefaultParentIndex:]...)
		defaultColumns[constants.DefaultParentIndex] = parentCol
	}
	return GetHeaders(allColumns, defaultColumns, customColumns)
}

func (r *Result) PrintText(out io.Writer, noHeaders bool) error {
	var resultPrint ResultPrint

	if r.ApiResponse != nil {
		requestId, err := GetRequestId(GetRequestPath(r.ApiResponse))
		if err != nil {
			return err
		}
		resultPrint.RequestId = requestId
	}

	if r.KeyValue != nil && r.Columns != nil {
		err := printText(out, r.Columns, r.KeyValue, noHeaders)
		if err != nil {
			return err
		}
	}

	if r.Resource != "" && r.Verb != "" {
		resultPrint.Message = fmt.Sprintf("Command %s %s has been successfully executed", r.Resource, r.Verb)
	} else if r.Message != "" {
		resultPrint.Message = r.Message
	}

	if resultPrint.RequestId != nil {
		_, err := fmt.Fprintf(out, "RequestId: %v\n", resultPrint.RequestId)
		if err != nil {
			return err
		}
	}

	if resultPrint.Message != nil {
		_, err := fmt.Fprintf(out, "Status: %v\n", resultPrint.Message)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Result) PrintJSON(out io.Writer) error {
	var resultPrint ResultPrint
	if r.Resource != "" && r.Verb != "" {
		oldStringKeptForBackwardsCompatibility := ""
		if r.WaitForRequest || r.WaitForState {
			// May the lord have mercy on our souls
			oldStringKeptForBackwardsCompatibility = "& wait"
		}
		r.Message = fmt.Sprintf("Command %s %s%s has been successfully executed",
			r.Resource, r.Verb, oldStringKeptForBackwardsCompatibility)
	} else if r.Message != "" {
		resultPrint.Message = r.Message
	}
	if r.ApiResponse != nil {
		requestId, err := GetRequestId(GetRequestPath(r.ApiResponse))
		if err != nil {
			return err
		}
		resultPrint.RequestId = requestId
	}
	// wtf
	resultPrint.Output = r.OutputJSON
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
