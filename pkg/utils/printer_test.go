package utils

import (
	"bufio"
	"bytes"
	"net/http"
	"regexp"
	"testing"

	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestNewPrinter_JSON(t *testing.T) {
	defer func(a func()) { ErrAction = a }(ErrAction)
	var b bytes.Buffer
	ErrAction = func() {}

	viper.Set(config.ArgOutput, "json")
	viper.Set(config.ArgQuiet, false)
	w := bufio.NewWriter(&b)
	NewPrinterRegistry(w, w)
	err := w.Flush()
	assert.NoError(t, err)
}

func TestNewPrinter_Text(t *testing.T) {
	defer func(a func()) { ErrAction = a }(ErrAction)
	var b bytes.Buffer
	ErrAction = func() {}

	viper.Set(config.ArgOutput, "text")
	viper.Set(config.ArgQuiet, false)
	w := bufio.NewWriter(&b)
	NewPrinterRegistry(w, w)
	err := w.Flush()
	assert.NoError(t, err)
}

func TestPrinter_PrintResultJson(t *testing.T) {
	defer func(a func()) { ErrAction = a }(ErrAction)

	var (
		b   bytes.Buffer
		str = `{
  "Message": "command executed"
}`
	)
	ErrAction = func() {}

	viper.Set(config.ArgOutput, "json")
	viper.Set(config.ArgQuiet, false)
	w := bufio.NewWriter(&b)
	reg := NewPrinterRegistry(w, w)
	p := reg["json"]
	p.SetStderr(w)
	p.SetStdout(w)
	res := Result{
		Message: "command executed",
	}

	p.Print(res)
	err := w.Flush()
	assert.NoError(t, err)

	re := regexp.MustCompile(str)
	assert.True(t, re.Match(b.Bytes()))
}

func TestPrinter_PrintStandardResultJson(t *testing.T) {
	defer func(a func()) { ErrAction = a }(ErrAction)

	var (
		b   bytes.Buffer
		str = `{
  "Message": "datacenter create command has been successfully executed"
}`
	)
	ErrAction = func() {}

	viper.Set(config.ArgOutput, "json")
	viper.Set(config.ArgQuiet, false)
	w := bufio.NewWriter(&b)
	reg := NewPrinterRegistry(w, w)
	p := reg["json"]
	p.SetStderr(w)
	p.SetStdout(w)
	res := Result{
		Resource: "datacenter",
		Verb:     "create",
	}

	p.Print(res)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(str)
	assert.True(t, re.Match(b.Bytes()))
}

func TestPrinter_PrintDefaultJson(t *testing.T) {
	defer func(a func()) { ErrAction = a }(ErrAction)

	var (
		b   bytes.Buffer
		str = `{
  "Message": "command executed"
}`
	)
	ErrAction = func() {}

	viper.Set(config.ArgOutput, "json")
	viper.Set(config.ArgQuiet, false)
	w := bufio.NewWriter(&b)
	reg := NewPrinterRegistry(w, w)
	p := reg["json"]
	p.SetStderr(w)
	p.SetStdout(w)
	input := "command executed"

	p.Print(input)
	err := w.Flush()
	assert.NoError(t, err)

	re := regexp.MustCompile(str)
	assert.True(t, re.Match(b.Bytes()))
}

func TestPrinter_PrintDefaultQuietJson(t *testing.T) {
	defer func(a func()) { ErrAction = a }(ErrAction)

	var (
		b bytes.Buffer
	)
	ErrAction = func() {}

	viper.Set(config.ArgOutput, "json")
	viper.Set(config.ArgQuiet, true)
	w := bufio.NewWriter(&b)
	reg := NewPrinterRegistry(w, w)
	p := reg["json"]
	p.SetStderr(w)
	p.SetStdout(w)
	input := "command executed"

	p.Print(input)
	err := w.Flush()
	assert.NoError(t, err)

	re := regexp.MustCompile(``)
	assert.True(t, re.Match(b.Bytes()))
}

func TestPrinter_ResultJsonRequestId(t *testing.T) {
	defer func(a func()) { ErrAction = a }(ErrAction)

	var (
		b   bytes.Buffer
		str = `{
  "RequestId": "123456"
}`
	)

	ErrAction = func() {}

	viper.Set(config.ArgOutput, "json")
	viper.Set(config.ArgQuiet, false)
	w := bufio.NewWriter(&b)
	reg := NewPrinterRegistry(w, w)
	p := reg["json"]
	p.SetStderr(w)
	p.SetStdout(w)
	res := Result{
		ApiResponse: &resources.Response{
			APIResponse: ionoscloud.APIResponse{
				Message: "Status OK",
				Response: &http.Response{
					Header: map[string][]string{},
				},
			},
		},
	}
	res.ApiResponse.Header.Add("location", "https://api.ionos.com/cloudapi/v5/requests/123456/status")

	p.Print(res)
	err := w.Flush()
	assert.NoError(t, err)

	re := regexp.MustCompile(str)
	assert.True(t, re.Match(b.Bytes()))
}

func TestPrinter_ResultJsonRequestIdErr(t *testing.T) {
	defer func(a func()) { ErrAction = a }(ErrAction)

	var (
		b      bytes.Buffer
		strErr = `{
  "Error": {},
  "Detail": "path does not contain https://api.ionos.com/cloudapi/v5"
}
`
	)
	ErrAction = func() {}

	viper.Set(config.ArgOutput, "json")
	viper.Set(config.ArgQuiet, false)
	w := bufio.NewWriter(&b)
	reg := NewPrinterRegistry(w, w)
	p := reg["json"]
	p.SetStderr(w)
	p.SetStdout(w)
	res := Result{
		ApiResponse: &resources.Response{
			APIResponse: ionoscloud.APIResponse{
				Message: "Status OK",
				Response: &http.Response{
					Header: map[string][]string{},
				},
			},
		},
	}
	res.ApiResponse.Header.Add("location", "https://api.ionos.com/requests/123456/status")

	p.Print(res)
	err := w.Flush()
	assert.NoError(t, err)

	re := regexp.MustCompile(strErr)
	assert.True(t, re.Match(b.Bytes()))
}

func TestPrinter_GetStdoutJSON(t *testing.T) {
	defer func(a func()) { ErrAction = a }(ErrAction)

	var (
		b bytes.Buffer
	)
	ErrAction = func() {}

	viper.Set(config.ArgOutput, "json")
	viper.Set(config.ArgQuiet, false)
	w := bufio.NewWriter(&b)
	reg := NewPrinterRegistry(w, w)
	out := reg["json"].GetStdout()

	err := w.Flush()
	assert.NoError(t, err)
	assert.True(t, out == w)
}

func TestPrinter_GetStderrJSON(t *testing.T) {
	defer func(a func()) { ErrAction = a }(ErrAction)

	var (
		b bytes.Buffer
	)
	ErrAction = func() {}

	viper.Set(config.ArgOutput, "json")
	viper.Set(config.ArgQuiet, false)
	w := bufio.NewWriter(&b)
	reg := NewPrinterRegistry(w, w)
	out := reg["json"].GetStderr()

	err := w.Flush()
	assert.NoError(t, err)
	assert.True(t, out == w)
}

func TestPrinter_PrintResultText(t *testing.T) {
	defer func(a func()) { ErrAction = a }(ErrAction)

	var b bytes.Buffer
	ErrAction = func() {}

	viper.Set(config.ArgOutput, "text")
	viper.Set(config.ArgQuiet, false)
	w := bufio.NewWriter(&b)
	reg := NewPrinterRegistry(w, w)
	p := reg["text"]
	p.SetStderr(w)
	p.SetStdout(w)
	res := Result{
		Message: "command executed",
	}

	p.Print(res)
	err := w.Flush()
	assert.NoError(t, err)

	re := regexp.MustCompile(`command executed`)
	assert.True(t, re.Match(b.Bytes()))
}

func TestPrinter_PrintResultTextRequestId(t *testing.T) {
	defer func(a func()) { ErrAction = a }(ErrAction)

	var (
		b bytes.Buffer
	)

	ErrAction = func() {}

	viper.Set(config.ArgOutput, "text")
	viper.Set(config.ArgQuiet, false)
	w := bufio.NewWriter(&b)
	reg := NewPrinterRegistry(w, w)
	p := reg["text"]
	p.SetStderr(w)
	p.SetStdout(w)
	res := Result{
		ApiResponse: &resources.Response{
			APIResponse: ionoscloud.APIResponse{
				Message: "Status OK",
				Response: &http.Response{
					Header: map[string][]string{},
				},
			},
		},
	}
	res.ApiResponse.Header.Add("location", "https://api.ionos.com/cloudapi/v5/requests/123456/status")

	p.Print(res)
	err := w.Flush()
	assert.NoError(t, err)

	re := regexp.MustCompile(`123456`)
	assert.True(t, re.Match(b.Bytes()))
}

func TestPrinter_PrintResultTextResource(t *testing.T) {
	defer func(a func()) { ErrAction = a }(ErrAction)

	var (
		b bytes.Buffer
	)

	ErrAction = func() {}

	viper.Set(config.ArgOutput, "text")
	viper.Set(config.ArgQuiet, false)
	w := bufio.NewWriter(&b)
	reg := NewPrinterRegistry(w, w)
	p := reg["text"]
	p.SetStderr(w)
	p.SetStdout(w)
	res := Result{
		Resource: "datacenter",
		Verb:     "create",
	}
	p.Print(res)
	err := w.Flush()
	assert.NoError(t, err)

	re := regexp.MustCompile(`datacenter create command has been successfully executed`)
	assert.True(t, re.Match(b.Bytes()))
}

func TestPrinter_PrintResultTextKeyValue(t *testing.T) {
	defer func(a func()) { ErrAction = a }(ErrAction)

	var (
		b      bytes.Buffer
		tabwrt = `ID    Name    Authorized   Age\[min]
123   dummy   true         1.230000`
	)

	ErrAction = func() {}

	viper.Set(config.ArgOutput, "text")
	viper.Set(config.ArgQuiet, false)
	w := bufio.NewWriter(&b)
	reg := NewPrinterRegistry(w, w)
	p := reg["text"]
	p.SetStderr(w)
	p.SetStdout(w)
	keyValueMap := map[string]interface{}{
		"ID":          123,
		"Name":        "dummy",
		"Authorized":  true,
		"Age[min]":    1.23,
		"Description": "dummy",
	}
	res := Result{
		Columns:  []string{"ID", "Name", "Authorized", "Age[min]"},
		KeyValue: []map[string]interface{}{keyValueMap},
	}

	p.Print(res)
	err := w.Flush()
	assert.NoError(t, err)

	re := regexp.MustCompile(tabwrt)
	assert.True(t, re.Match(b.Bytes()))
}

func TestPrinter_PrintDefaultText(t *testing.T) {
	defer func(a func()) { ErrAction = a }(ErrAction)

	var (
		b bytes.Buffer
	)

	ErrAction = func() {}

	viper.Set(config.ArgOutput, "text")
	viper.Set(config.ArgQuiet, false)
	w := bufio.NewWriter(&b)
	reg := NewPrinterRegistry(w, w)
	p := reg["text"]
	p.SetStderr(w)
	p.SetStdout(w)
	input := "dummy"

	p.Print(input)
	err := w.Flush()
	assert.NoError(t, err)

	re := regexp.MustCompile(input)
	assert.True(t, re.Match(b.Bytes()))
}

func TestPrinter_UnknownFormatType(t *testing.T) {
	defer func(a func()) { ErrAction = a }(ErrAction)

	var b bytes.Buffer
	ErrAction = func() {}

	viper.Set(config.ArgOutput, "dummy")
	viper.Set(config.ArgQuiet, false)
	w := bufio.NewWriter(&b)
	NewPrinterRegistry(w, w)
	err := w.Flush()
	assert.NoError(t, err)

	re := regexp.MustCompile(`unknown type format dummy. Hint: use --output json|text`)
	assert.True(t, re.Match(b.Bytes()))
}

func TestPrinter_ResultQuiet(t *testing.T) {
	defer func(a func()) { ErrAction = a }(ErrAction)

	var b bytes.Buffer
	ErrAction = func() {}

	viper.Set(config.ArgOutput, "text")
	viper.Set(config.ArgQuiet, true)
	w := bufio.NewWriter(&b)
	reg := NewPrinterRegistry(w, w)
	p := reg["text"]
	p.SetStderr(w)
	p.SetStdout(w)
	res := Result{
		Message: "command executed",
	}

	p.Print(res)
	err := w.Flush()
	assert.NoError(t, err)

	re := regexp.MustCompile(``)
	assert.True(t, re.Match(b.Bytes()))
}

func TestPrinter_GetStdoutText(t *testing.T) {
	defer func(a func()) { ErrAction = a }(ErrAction)

	var (
		b bytes.Buffer
	)
	ErrAction = func() {}

	viper.Set(config.ArgOutput, "text")
	viper.Set(config.ArgQuiet, false)
	w := bufio.NewWriter(&b)
	reg := NewPrinterRegistry(w, w)
	out := reg["text"].GetStdout()

	err := w.Flush()
	assert.NoError(t, err)
	assert.True(t, out == w)
}

func TestPrinter_GetStderrText(t *testing.T) {
	defer func(a func()) { ErrAction = a }(ErrAction)

	var (
		b bytes.Buffer
	)
	ErrAction = func() {}

	viper.Set(config.ArgOutput, "text")
	viper.Set(config.ArgQuiet, false)
	w := bufio.NewWriter(&b)
	reg := NewPrinterRegistry(w, w)
	out := reg["text"].GetStderr()

	err := w.Flush()
	assert.NoError(t, err)
	assert.True(t, out == w)
}
