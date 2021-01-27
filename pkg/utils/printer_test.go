package utils

import (
	"bufio"
	"bytes"
	"io"
	"net/http"
	"regexp"
	"testing"

	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var r io.Reader

func TestNewPrinter_JSON(t *testing.T) {
	defer func(a func()) { ErrAction = a }(ErrAction)
	var b bytes.Buffer
	ErrAction = func() {}

	viper.Set(config.ArgOutput, "json")
	viper.Set(config.ArgQuiet, false)
	w := bufio.NewWriter(&b)
	NewPrinter()
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
	NewPrinter()
	err := w.Flush()
	assert.NoError(t, err)
}

func TestNewPrinter_TextQuiet(t *testing.T) {
	defer func(a func()) { ErrAction = a }(ErrAction)
	var b bytes.Buffer
	ErrAction = func() {}

	viper.Set(config.ArgOutput, "text")
	viper.Set(config.ArgQuiet, true)
	w := bufio.NewWriter(&b)
	NewPrinter()
	err := w.Flush()
	assert.NoError(t, err)
}

func TestPrinter_ResultJson(t *testing.T) {
	defer func(a func()) { ErrAction = a }(ErrAction)

	var (
		b   bytes.Buffer
		str = `{
  "Status": "command executed"
}`
	)
	ErrAction = func() {}

	viper.Set(config.ArgOutput, "json")
	viper.Set(config.ArgQuiet, false)
	w := bufio.NewWriter(&b)
	p := &Printer{
		OutputFlag: "json",
		Stdin:      r,
		Stdout:     w,
		Stderr:     w,
	}
	res := &SuccessResult{
		Message: "command executed",
	}

	p.Result(res)
	err := w.Flush()
	assert.NoError(t, err)

	re := regexp.MustCompile(str)
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
	p := &Printer{
		OutputFlag: "json",
		Stdin:      r,
		Stdout:     w,
		Stderr:     w,
	}
	res := &SuccessResult{
		Message: "",
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

	p.Result(res)
	err := w.Flush()
	assert.NoError(t, err)

	re := regexp.MustCompile(str)
	assert.True(t, re.Match(b.Bytes()))
}

func TestPrinter_ResultJsonRequestIdErr(t *testing.T) {
	defer func(a func()) { ErrAction = a }(ErrAction)

	var b bytes.Buffer
	ErrAction = func() {}

	viper.Set(config.ArgOutput, "json")
	viper.Set(config.ArgQuiet, false)
	w := bufio.NewWriter(&b)
	p := &Printer{
		OutputFlag: "json",
		Stdin:      r,
		Stdout:     w,
		Stderr:     w,
	}
	res := &SuccessResult{
		Message: "",
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

	p.Result(res)
	err := w.Flush()
	assert.NoError(t, err)

	re := regexp.MustCompile(`{"error":{},"detail":"path does not contain https://api.ionos.com/cloudapi/v5"}`)
	assert.True(t, re.Match(b.Bytes()))
}

func TestPrinter_ResultText(t *testing.T) {
	defer func(a func()) { ErrAction = a }(ErrAction)

	var b bytes.Buffer
	ErrAction = func() {}

	viper.Set(config.ArgOutput, "text")
	viper.Set(config.ArgQuiet, false)
	w := bufio.NewWriter(&b)
	p := &Printer{
		OutputFlag: "text",
		Stdin:      r,
		Stdout:     w,
		Stderr:     w,
	}
	res := &SuccessResult{
		Message: "command executed",
	}

	p.Result(res)
	err := w.Flush()
	assert.NoError(t, err)

	re := regexp.MustCompile(`command executed`)
	assert.True(t, re.Match(b.Bytes()))
}

func TestPrinter_ResultTextRequestId(t *testing.T) {
	defer func(a func()) { ErrAction = a }(ErrAction)

	var (
		b bytes.Buffer
	)

	ErrAction = func() {}

	viper.Set(config.ArgOutput, "text")
	viper.Set(config.ArgQuiet, false)
	w := bufio.NewWriter(&b)
	p := &Printer{
		OutputFlag: "text",
		Stdin:      r,
		Stdout:     w,
		Stderr:     w,
	}
	res := &SuccessResult{
		Message: "",
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

	p.Result(res)
	err := w.Flush()
	assert.NoError(t, err)

	re := regexp.MustCompile(`123456`)
	assert.True(t, re.Match(b.Bytes()))
}

func TestPrinter_ResultTextResource(t *testing.T) {
	defer func(a func()) { ErrAction = a }(ErrAction)

	var (
		b bytes.Buffer
	)

	ErrAction = func() {}

	viper.Set(config.ArgOutput, "text")
	viper.Set(config.ArgQuiet, false)
	w := bufio.NewWriter(&b)
	p := &Printer{
		OutputFlag: "text",
		Stdin:      r,
		Stdout:     w,
		Stderr:     w,
	}
	res := &SuccessResult{
		Resource: "datacenter",
		Verb:     "create",
	}
	p.Result(res)
	err := w.Flush()
	assert.NoError(t, err)

	re := regexp.MustCompile(`datacenter create command has been successfully executed`)
	assert.True(t, re.Match(b.Bytes()))
}

func TestPrinter_ResultTextKeyValue(t *testing.T) {
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
	p := &Printer{
		OutputFlag: "text",
		Stdin:      r,
		Stdout:     w,
		Stderr:     w,
	}
	keyValueMap := map[string]interface{}{
		"ID":          123,
		"Name":        "dummy",
		"Authorized":  true,
		"Age[min]":    1.23,
		"Description": "dummy",
	}
	res := &SuccessResult{
		Columns:  []string{"ID", "Name", "Authorized", "Age[min]"},
		KeyValue: []map[string]interface{}{keyValueMap},
	}

	p.Result(res)
	err := w.Flush()
	assert.NoError(t, err)

	re := regexp.MustCompile(tabwrt)
	assert.True(t, re.Match(b.Bytes()))
}

func TestPrinter_ResultDefault(t *testing.T) {
	defer func(a func()) { ErrAction = a }(ErrAction)

	var b bytes.Buffer
	ErrAction = func() {}

	viper.Set(config.ArgOutput, "dummy")
	viper.Set(config.ArgQuiet, false)
	w := bufio.NewWriter(&b)
	p := &Printer{
		OutputFlag: "dummy",
		Stdin:      r,
		Stdout:     w,
		Stderr:     w,
	}
	res := &SuccessResult{
		Message: "command executed",
	}

	p.Result(res)
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
	p := &Printer{
		OutputFlag: "text",
		Stdin:      r,
		Stdout:     w,
		Stderr:     w,
	}
	res := &SuccessResult{
		Message: "command executed",
	}

	p.Result(res)
	err := w.Flush()
	assert.NoError(t, err)

	re := regexp.MustCompile(``)
	assert.True(t, re.Match(b.Bytes()))
}

func TestPrinter_LogText(t *testing.T) {
	defer func(a func()) { ErrAction = a }(ErrAction)

	var b bytes.Buffer
	ErrAction = func() {}

	viper.Set(config.ArgOutput, "text")
	viper.Set(config.ArgQuiet, false)
	w := bufio.NewWriter(&b)
	p := &Printer{
		OutputFlag: "text",
		Stdin:      r,
		Stdout:     w,
		Stderr:     w,
	}

	p.Log()
	err := w.Flush()
	assert.NoError(t, err)

	re := regexp.MustCompile(``)
	assert.True(t, re.Match(b.Bytes()))
}

func TestPrinter_LogJSON(t *testing.T) {
	defer func(a func()) { ErrAction = a }(ErrAction)

	var b bytes.Buffer
	ErrAction = func() {}

	viper.Set(config.ArgOutput, "json")
	viper.Set(config.ArgQuiet, false)
	w := bufio.NewWriter(&b)
	p := &Printer{
		OutputFlag: "json",
		Stdin:      r,
		Stdout:     w,
		Stderr:     w,
	}

	p.Log()
	err := w.Flush()
	assert.NoError(t, err)

	re := regexp.MustCompile(``)
	assert.True(t, re.Match(b.Bytes()))
}

func TestPrinter_LogDefault(t *testing.T) {
	defer func(a func()) { ErrAction = a }(ErrAction)

	var b bytes.Buffer
	ErrAction = func() {}

	viper.Set(config.ArgOutput, "dummy")
	viper.Set(config.ArgQuiet, false)
	w := bufio.NewWriter(&b)
	p := &Printer{
		OutputFlag: "dummy",
		Stdin:      r,
		Stdout:     w,
		Stderr:     w,
	}

	p.Log()
	err := w.Flush()
	assert.NoError(t, err)

	re := regexp.MustCompile(`unknown type format dummy. Hint: use --output json|text`)
	assert.True(t, re.Match(b.Bytes()))
}
