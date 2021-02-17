package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/fatih/color"
	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/resources"
	"github.com/spf13/viper"
)

// Type defines an formatter format.
type PrinterType string

func (p PrinterType) String() string {
	return string(p)
}

const (
	// PrinterTypeJSON defines a JSON formatter.
	PrinterTypeJSON = PrinterType("json")

	// PrinterTypeText defines a human readable formatted formatter.
	PrinterTypeText = PrinterType("text")
)

type PrinterRegistry map[string]PrintService

func NewPrinterRegistry(out, outErr io.Writer) PrinterRegistry {
	if viper.GetString(config.ArgOutput) != PrinterTypeJSON.String() &&
		viper.GetString(config.ArgOutput) != PrinterTypeText.String() {
		err := errors.New(fmt.Sprintf(unknownTypeFormatErr, viper.GetString(config.ArgOutput)))
		CheckError(err, outErr)
	}

	return PrinterRegistry{
		PrinterTypeJSON.String(): &JSONPrinter{
			Stderr: outErr,
			Stdout: out,
		},
		PrinterTypeText.String(): &TextPrinter{
			Stderr: outErr,
			Stdout: out,
		},
	}
}

type PrintService interface {
	Print(interface{})

	GetStdout() io.Writer
	SetStdout(io.Writer)
	GetStderr() io.Writer
	SetStderr(io.Writer)
}

type JSONPrinter struct {
	Stdout io.Writer
	Stderr io.Writer
}

func (p *JSONPrinter) Print(v interface{}) {
	var resultPrint ResultPrint

	switch v.(type) {
	case Result:
		if v.(Result).Resource != "" && v.(Result).Verb != "" {
			resultPrint.Message = standardSuccessMsg(v.(Result).Resource, v.(Result).Verb, v.(Result).WaitFlag)
		} else if v.(Result).Message != "" {
			resultPrint.Message = v.(Result).Message
		}
		if v.(Result).ApiResponse != nil {
			requestId, err := GetRequestId(GetRequestPath(v.(Result).ApiResponse))
			CheckError(err, p.Stderr)
			if requestId != nil {
				resultPrint.RequestId = *requestId
			}
		}
		resultPrint.Output = v.(Result).OutputJSON
		if !structs.IsZero(resultPrint) {
			err := writeJSON(&resultPrint, p.Stdout)
			CheckError(err, p.Stderr)
		}
	default:
		var msg DefaultMsgPrint
		msg.Message = v
		err := writeJSON(&msg, p.Stdout)
		CheckError(err, p.Stderr)
	}
	return
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

func (p *TextPrinter) Print(v interface{}) {
	var resultPrint ResultPrint

	switch v.(type) {
	case Result:
		if v.(Result).Resource != "" && v.(Result).Verb != "" {
			resultPrint.Message = standardSuccessMsg(v.(Result).Resource, v.(Result).Verb, v.(Result).WaitFlag)
		} else if v.(Result).Message != "" {
			resultPrint.Message = v.(Result).Message
		}
		if v.(Result).ApiResponse != nil {
			requestId, err := GetRequestId(GetRequestPath(v.(Result).ApiResponse))
			CheckError(err, p.Stderr)
			if requestId != nil {
				resultPrint.RequestId = *requestId
			}
		}
		if v.(Result).KeyValue != nil && v.(Result).Columns != nil {
			err := printText(p.Stdout, v.(Result).Columns, v.(Result).KeyValue)
			CheckError(err, p.Stderr)
		}
		if resultPrint.RequestId != nil {
			requestIdMsg(p.Stdout, "%v", resultPrint.RequestId)
		}
		if resultPrint.Message != nil {
			statusMsg(p.Stdout, "%v", resultPrint.Message)
		}
	default:
		_, err := fmt.Fprintf(p.Stdout, "%v\n", v)
		CheckError(err, p.Stderr)
	}

	return
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
	Message  string
	Resource string
	Verb     string
	WaitFlag bool

	Columns    []string
	KeyValue   []map[string]interface{}
	OutputJSON interface{}

	ApiResponse *resources.Response
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
	waitStandardSuccessMessages = "Command %s %s and request have been successfully executed"
	unknownTypeFormatErr        = "unknown type format %s. Hint: use --output json|text"
)

func standardSuccessMsg(resource, verb string, wait bool) string {
	if wait {
		return fmt.Sprintf(waitStandardSuccessMessages, resource, verb)
	}
	return fmt.Sprintf(standardSuccessMessages, resource, verb)
}

func requestIdMsg(writer io.Writer, msg string, args ...interface{}) {
	colorWarn := color.BlueString("RequestId")
	fmt.Fprintf(writer, "\u2714 %s: %s\n", colorWarn, fmt.Sprintf(msg, args...))
}

func statusMsg(writer io.Writer, msg string, args ...interface{}) {
	colorWarn := color.GreenString("Status")
	fmt.Fprintf(writer, "\u2714 %s: %s\n", colorWarn, fmt.Sprintf(msg, args...))
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

func writeJSON(item interface{}, writer io.Writer) error {
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

func GetRequestPath(r *resources.Response) string {
	if r != nil {
		return r.Header.Get("location")
	}
	return ""
}
