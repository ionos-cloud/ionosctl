package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/fatih/color"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/resources"
	"github.com/sirupsen/logrus"
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

type Printer struct {
	OutputFlag string
	Stdin      io.Reader
	Stdout     io.Writer
	Stderr     io.Writer
}

func NewPrinter() *Printer {
	var out, outErr io.Writer
	if viper.GetBool(config.ArgQuiet) {
		var execOut bytes.Buffer
		out = &execOut
		outErr = os.Stderr
	} else {
		out = os.Stdout
		outErr = os.Stderr
	}

	if viper.GetString(config.ArgOutput) != PrinterTypeJSON.String() &&
		viper.GetString(config.ArgOutput) != PrinterTypeText.String() {
		err := errors.New(fmt.Sprintf(unknownTypeFormatErr, viper.GetString(config.ArgOutput)))
		CheckError(err, outErr)
	}

	return &Printer{
		OutputFlag: viper.GetString(config.ArgOutput),
		Stdin:      os.Stdin,
		Stdout:     out,
		Stderr:     outErr,
	}
}

func (p *Printer) Result(res *SuccessResult) {
	if viper.GetBool(config.ArgQuiet) {
		return
	}

	var msg PrintResult

	if res.Resource != "" && res.Verb != "" {
		msg.Message = fmt.Sprintf(standardSuccessMessages, res.Resource, res.Verb)
	} else {
		msg.Message = res.Message
	}
	if res.ApiResponse != nil {
		path, err := getRequestId(res.ApiResponse.Header.Get("location"))
		CheckError(err, p.Stderr)
		msg.Response = path
	}

	switch p.OutputFlag {
	case PrinterTypeJSON.String():
		msg.Output = res.OutputJSON
		err := writeJSON(&msg, p.Stdout)
		CheckError(err, p.Stderr)
	case PrinterTypeText.String():
		if res.KeyValue != nil {
			err := printText(p.Stdout, res.Columns, res.KeyValue)
			CheckError(err, p.Stderr)
		}
		if msg.Response != nil {
			requestIdMsg(p.Stdout, "%s", msg.Response)
		}
		if msg.Message != "" {
			succesMsg(p.Stdout, msg.Message)
		}
	default:
		err := errors.New(fmt.Sprintf(unknownTypeFormatErr, p.OutputFlag))
		CheckError(err, p.Stderr)
	}
}

func (p *Printer) Log() *Logger {
	l := logrus.New()
	l.SetOutput(p.Stdout)

	switch p.OutputFlag {
	case PrinterTypeJSON.String():
		jsonFormat := new(logrus.JSONFormatter)
		jsonFormat.PrettyPrint = true
		l.SetFormatter(jsonFormat)
	case PrinterTypeText.String():
		txtFormat := new(logrus.TextFormatter)
		txtFormat.ForceColors = true
		l.SetFormatter(txtFormat)
	default:
		err := errors.New(fmt.Sprintf(unknownTypeFormatErr, p.OutputFlag))
		CheckError(err, p.Stderr)
	}

	return &Logger{
		Logger: l,
	}
}

/*
	Printer Result
*/

type SuccessResult struct {
	Message  string
	Resource string
	Verb     string

	Columns    []string
	KeyValue   []map[string]interface{}
	OutputJSON interface{}

	ApiResponse *resources.Response
}

type PrintResult struct {
	Message  string      `json:"Status,omitempty"`
	Response interface{} `json:"RequestId,omitempty"`
	Output   interface{} `json:"Resources,omitempty"`
}

var (
	standardSuccessMessages = "%s %s command has been successfully executed"
	unknownTypeFormatErr    = "unknown type format %s. Hint: use --output json|text"
)

func requestIdMsg(writer io.Writer, msg string, args ...interface{}) {
	colorWarn := color.BlueString("Request Id")
	fmt.Fprintf(writer, "\u2714 %s: %s\n", colorWarn, fmt.Sprintf(msg, args...))
}

func succesMsg(writer io.Writer, msg string, args ...interface{}) {
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

func getRequestId(path string) (string, error) {
	if !strings.Contains(path, config.DefaultApiURL) {
		return "", errors.New("path does not contain " + config.DefaultApiURL)
	}
	str := strings.Split(path, "/")
	return str[len(str)-2], nil
}

/*
	Printer Log
*/

type Logger struct {
	*logrus.Logger
}
