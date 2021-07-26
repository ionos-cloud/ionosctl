package printer

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ionos-cloud/ionosctl/pkg/resources/v5"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/spf13/viper"
)

// Type defines an formatter format.
type Type string

func (p Type) String() string {
	return string(p)
}

const (
	// TypeJSON defines a JSON formatter.
	TypeJSON = Type("json")

	// TypeText defines a human readable formatted formatter.
	TypeText = Type("text")
)

type Registry map[string]PrintService

func NewPrinterRegistry(out, outErr io.Writer) (Registry, error) {
	if viper.GetString(config.ArgOutput) != TypeJSON.String() &&
		viper.GetString(config.ArgOutput) != TypeText.String() {
		return nil, errors.New(fmt.Sprintf(unknownTypeFormatErr, viper.GetString(config.ArgOutput)))
	}

	return Registry{
		TypeJSON.String(): &JSONPrinter{
			Stderr: outErr,
			Stdout: out,
		},
		TypeText.String(): &TextPrinter{
			Stderr: outErr,
			Stdout: out,
		},
	}, nil
}

type PrintService interface {
	Print(interface{}) error

	GetStdout() io.Writer
	SetStdout(io.Writer)
	GetStderr() io.Writer
	SetStderr(io.Writer)
}

type JSONPrinter struct {
	Stdout io.Writer
	Stderr io.Writer
}

func (p *JSONPrinter) Print(v interface{}) error {
	var resultPrint ResultPrint

	switch v.(type) {
	case Result:
		if v.(Result).Resource != "" && v.(Result).Verb != "" {
			resultPrint.Message = standardSuccessMsg(v.(Result).Resource, v.(Result).Verb, v.(Result).WaitForRequest, v.(Result).WaitForState)
		} else if v.(Result).Message != "" {
			resultPrint.Message = v.(Result).Message
		}
		if v.(Result).ApiResponse != nil {
			requestId, err := GetRequestId(GetRequestPath(v.(Result).ApiResponse))
			if err != nil {
				return err
			}
			if requestId != nil {
				resultPrint.RequestId = *requestId
			}
		}
		resultPrint.Output = v.(Result).OutputJSON
		if !structs.IsZero(resultPrint) {
			err := WriteJSON(&resultPrint, p.Stdout)
			if err != nil {
				return err
			}
		}
	default:
		var msg DefaultMsgPrint
		msg.Message = v
		err := WriteJSON(&msg, p.Stdout)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *JSONPrinter) GetStdout() io.Writer {
	return p.Stdout
}

func (p *JSONPrinter) SetStdout(writer io.Writer) {
	p.Stdout = writer
}

func (p *JSONPrinter) GetStderr() io.Writer {
	return p.Stderr
}

func (p *JSONPrinter) SetStderr(writer io.Writer) {
	p.Stderr = writer
}

type TextPrinter struct {
	Stdout io.Writer
	Stderr io.Writer
}

func (p *TextPrinter) Print(v interface{}) error {
	var resultPrint ResultPrint

	switch v.(type) {
	case Result:
		if v.(Result).Resource != "" && v.(Result).Verb != "" {
			resultPrint.Message = standardSuccessMsg(v.(Result).Resource, v.(Result).Verb, v.(Result).WaitForRequest, v.(Result).WaitForState)
		} else if v.(Result).Message != "" {
			resultPrint.Message = v.(Result).Message
		}
		if v.(Result).ApiResponse != nil {
			requestId, err := GetRequestId(GetRequestPath(v.(Result).ApiResponse))
			if err != nil {
				return err
			}
			if requestId != nil {
				resultPrint.RequestId = *requestId
			}
		}
		if v.(Result).KeyValue != nil && v.(Result).Columns != nil {
			err := printText(p.Stdout, v.(Result).Columns, v.(Result).KeyValue)
			if err != nil {
				return err
			}
		}
		if resultPrint.RequestId != nil {
			requestIdMsg(p.Stdout, "%v", resultPrint.RequestId)
		}
		if resultPrint.Message != nil {
			statusMsg(p.Stdout, "%v", resultPrint.Message)
		}
	case string:
		if strings.HasSuffix(v.(string), "\n") {
			if _, err := fmt.Fprintf(p.Stdout, "%v", v); err != nil {
				return err
			}
		} else {
			if _, err := fmt.Fprintf(p.Stdout, "%v\n", v); err != nil {
				return err
			}
		}
	default:
		_, err := fmt.Fprintf(p.Stdout, "%v\n", v)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *TextPrinter) GetStdout() io.Writer {
	return p.Stdout
}

func (p *TextPrinter) SetStdout(writer io.Writer) {
	p.Stdout = writer
}

func (p *TextPrinter) GetStderr() io.Writer {
	return p.Stderr
}

func (p *TextPrinter) SetStderr(writer io.Writer) {
	p.Stderr = writer
}

type Result struct {
	Message        string
	Resource       string
	Verb           string
	WaitForRequest bool
	WaitForState   bool

	Columns    []string
	KeyValue   []map[string]interface{}
	OutputJSON interface{}

	ApiResponse *v5.Response
}

type ResultPrint struct {
	Message   interface{} `json:"Status,omitempty"`
	RequestId interface{} `json:"RequestId,omitempty"`
	Output    interface{} `json:"Resources,omitempty"`
}

type DefaultMsgPrint struct {
	Message interface{} `json:"Message,omitempty"`
}

var (
	standardSuccessMessages     = "Command %s %s has been successfully executed"
	waitStandardSuccessMessages = "Command %s %s & wait have been successfully executed"
	unknownTypeFormatErr        = "unknown type format %s. Hint: use --output json|text"
)

func standardSuccessMsg(resource, verb string, waitRequest, waitState bool) string {
	if waitRequest || waitState {
		return fmt.Sprintf(waitStandardSuccessMessages, resource, verb)
	}
	return fmt.Sprintf(standardSuccessMessages, resource, verb)
}

func requestIdMsg(writer io.Writer, msg string, args ...interface{}) {
	fmt.Fprintf(writer, "RequestId: %s\n", fmt.Sprintf(msg, args...))
}

func statusMsg(writer io.Writer, msg string, args ...interface{}) {
	fmt.Fprintf(writer, "Status: %s\n", fmt.Sprintf(msg, args...))
}

func printText(out io.Writer, cols []string, keyValueMap []map[string]interface{}) error {
	w := new(tabwriter.Writer)
	w.Init(out, 5, 0, 3, ' ', tabwriter.StripEscape)

	headers := []string{}
	for _, col := range cols {
		headers = append(headers, col)
	}
	fmt.Fprintln(w, strings.Join(headers, "\t"))

	for _, r := range keyValueMap {
		values := []interface{}{}
		formats := []string{}

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
		fmt.Fprintf(w, format+"\n", values...)
	}

	return w.Flush()
}

func WriteJSON(item interface{}, writer io.Writer) error {
	j, err := json.MarshalIndent(item, "", "  ")
	if err != nil {
		return err
	}
	fmt.Fprintf(writer, "%s\n", string(j))
	return nil
}

func GetRequestId(path string) (*string, error) {
	if !strings.Contains(path, config.DefaultApiURL) {
		return nil, errors.New("path does not contain " + config.DefaultApiURL)
	}
	str := strings.Split(path, "/")
	return &str[len(str)-2], nil
}

func GetRequestPath(r *v5.Response) string {
	if r != nil {
		return r.Header.Get("location")
	}
	return ""
}
