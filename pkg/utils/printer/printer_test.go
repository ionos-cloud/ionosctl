package printer

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
	var b bytes.Buffer

	viper.Set(config.ArgOutput, "json")
	viper.Set(config.ArgQuiet, false)
	w := bufio.NewWriter(&b)
	_, err := NewPrinterRegistry(w, w)
	assert.NoError(t, err)
}

func TestNewPrinter_Text(t *testing.T) {
	var b bytes.Buffer

	viper.Set(config.ArgOutput, "text")
	viper.Set(config.ArgQuiet, false)
	w := bufio.NewWriter(&b)
	_, err := NewPrinterRegistry(w, w)
	assert.NoError(t, err)
}

func TestPrinter_PrintResultJson(t *testing.T) {
	var (
		b   bytes.Buffer
		str = `{
  "Status": "command executed"
}`
	)

	viper.Set(config.ArgOutput, "json")
	viper.Set(config.ArgQuiet, false)
	w := bufio.NewWriter(&b)
	reg, err := NewPrinterRegistry(w, w)
	assert.NoError(t, err)
	p := reg["json"]
	p.SetStderr(w)
	p.SetStdout(w)
	res := Result{
		Message: "command executed",
	}

	p.Print(res)
	err = w.Flush()
	assert.NoError(t, err)

	re := regexp.MustCompile(str)
	assert.True(t, re.Match(b.Bytes()))
}

func TestPrinter_PrintStandardResultJson(t *testing.T) {
	var (
		b   bytes.Buffer
		str = `{
  "Status": "Command datacenter create has been successfully executed"
}`
	)

	viper.Set(config.ArgOutput, "json")
	viper.Set(config.ArgQuiet, false)
	w := bufio.NewWriter(&b)
	reg, err := NewPrinterRegistry(w, w)
	assert.NoError(t, err)
	p := reg["json"]
	p.SetStderr(w)
	p.SetStdout(w)
	res := Result{
		Resource: "datacenter",
		Verb:     "create",
	}

	p.Print(res)
	err = w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(str)
	assert.True(t, re.Match(b.Bytes()))
}

func TestPrinter_PrintDefaultJson(t *testing.T) {
	var (
		b   bytes.Buffer
		str = `{
  "Message": "command executed"
}`
	)

	viper.Set(config.ArgOutput, "json")
	viper.Set(config.ArgQuiet, false)
	w := bufio.NewWriter(&b)
	reg, err := NewPrinterRegistry(w, w)
	assert.NoError(t, err)
	p := reg["json"]
	p.SetStderr(w)
	p.SetStdout(w)
	input := "command executed"

	p.Print(input)
	err = w.Flush()
	assert.NoError(t, err)

	re := regexp.MustCompile(str)
	assert.True(t, re.Match(b.Bytes()))
}

func TestPrinter_PrintDefaultQuietJson(t *testing.T) {
	var b bytes.Buffer

	viper.Set(config.ArgOutput, "json")
	viper.Set(config.ArgQuiet, true)
	w := bufio.NewWriter(&b)
	reg, err := NewPrinterRegistry(w, w)
	assert.NoError(t, err)
	p := reg["json"]
	p.SetStderr(w)
	p.SetStdout(w)
	input := "command executed"

	p.Print(input)
	err = w.Flush()
	assert.NoError(t, err)

	re := regexp.MustCompile(``)
	assert.True(t, re.Match(b.Bytes()))
}

func TestPrinter_ResultJsonRequestId(t *testing.T) {
	var (
		b   bytes.Buffer
		str = `{
  "RequestId": "123456"
}`
	)

	viper.Set(config.ArgOutput, "json")
	viper.Set(config.ArgQuiet, false)
	w := bufio.NewWriter(&b)
	reg, err := NewPrinterRegistry(w, w)
	assert.NoError(t, err)
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
	err = w.Flush()
	assert.NoError(t, err)

	re := regexp.MustCompile(str)
	assert.True(t, re.Match(b.Bytes()))
}

func TestPrinter_ResultJsonRequestIdErr(t *testing.T) {
	var (
		b      bytes.Buffer
		strErr = `path does not contain https://api.ionos.com/cloudapi/v5`
	)

	viper.Set(config.ArgOutput, "json")
	viper.Set(config.ArgQuiet, false)
	w := bufio.NewWriter(&b)
	reg, err := NewPrinterRegistry(w, w)
	assert.NoError(t, err)
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

	err = p.Print(res)
	assert.Error(t, err)
	assert.True(t, err.Error() == strErr)
}

func TestPrinter_GetStdoutJSON(t *testing.T) {
	var b bytes.Buffer

	viper.Set(config.ArgOutput, "json")
	viper.Set(config.ArgQuiet, false)
	w := bufio.NewWriter(&b)
	reg, err := NewPrinterRegistry(w, w)
	assert.NoError(t, err)
	out := reg["json"].GetStdout()

	err = w.Flush()
	assert.NoError(t, err)
	assert.True(t, out == w)
}

func TestPrinter_GetStderrJSON(t *testing.T) {
	var b bytes.Buffer

	viper.Set(config.ArgOutput, "json")
	viper.Set(config.ArgQuiet, false)
	w := bufio.NewWriter(&b)
	reg, err := NewPrinterRegistry(w, w)
	assert.NoError(t, err)
	out := reg["json"].GetStderr()

	err = w.Flush()
	assert.NoError(t, err)
	assert.True(t, out == w)
}

func TestPrinter_PrintResultText(t *testing.T) {
	var b bytes.Buffer

	viper.Set(config.ArgOutput, "text")
	viper.Set(config.ArgQuiet, false)
	w := bufio.NewWriter(&b)
	reg, err := NewPrinterRegistry(w, w)
	assert.NoError(t, err)
	p := reg["text"]
	p.SetStderr(w)
	p.SetStdout(w)
	res := Result{
		Message: "command executed",
	}

	p.Print(res)
	err = w.Flush()
	assert.NoError(t, err)

	re := regexp.MustCompile(`command executed`)
	assert.True(t, re.Match(b.Bytes()))
}

func TestPrinter_PrintResultTextRequestId(t *testing.T) {
	var b bytes.Buffer

	viper.Set(config.ArgOutput, "text")
	viper.Set(config.ArgQuiet, false)
	w := bufio.NewWriter(&b)
	reg, err := NewPrinterRegistry(w, w)
	assert.NoError(t, err)
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
	err = w.Flush()
	assert.NoError(t, err)

	re := regexp.MustCompile(`123456`)
	assert.True(t, re.Match(b.Bytes()))
}

func TestPrinter_PrintResultTextResource(t *testing.T) {
	var b bytes.Buffer

	viper.Set(config.ArgOutput, "text")
	viper.Set(config.ArgQuiet, false)
	w := bufio.NewWriter(&b)
	reg, err := NewPrinterRegistry(w, w)
	assert.NoError(t, err)
	p := reg["text"]
	p.SetStderr(w)
	p.SetStdout(w)
	res := Result{
		Resource: "datacenter",
		Verb:     "create",
		WaitFlag: false,
	}
	p.Print(res)
	err = w.Flush()
	assert.NoError(t, err)

	re := regexp.MustCompile(`Command datacenter create has been successfully executed`)
	assert.True(t, re.Match(b.Bytes()))
}

func TestPrinter_PrintResultTextWaitResource(t *testing.T) {
	var b bytes.Buffer

	viper.Set(config.ArgOutput, "text")
	viper.Set(config.ArgQuiet, false)
	w := bufio.NewWriter(&b)
	reg, err := NewPrinterRegistry(w, w)
	assert.NoError(t, err)
	p := reg["text"]
	p.SetStderr(w)
	p.SetStdout(w)
	res := Result{
		Resource: "datacenter",
		Verb:     "create",
		WaitFlag: true,
	}
	p.Print(res)
	err = w.Flush()
	assert.NoError(t, err)

	re := regexp.MustCompile(`Command datacenter create and request have been successfully executed`)
	assert.True(t, re.Match(b.Bytes()))
}

func TestPrinter_PrintResultTextKeyValue(t *testing.T) {
	var (
		b      bytes.Buffer
		tabwrt = `ID    Name    Authorized   Age\[min]
123   dummy   true         1.230000`
	)

	viper.Set(config.ArgOutput, "text")
	viper.Set(config.ArgQuiet, false)
	w := bufio.NewWriter(&b)
	reg, err := NewPrinterRegistry(w, w)
	assert.NoError(t, err)
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
	err = w.Flush()
	assert.NoError(t, err)

	re := regexp.MustCompile(tabwrt)
	assert.True(t, re.Match(b.Bytes()))
}

func TestPrinter_PrintDefaultText(t *testing.T) {
	var b bytes.Buffer

	viper.Set(config.ArgOutput, "text")
	viper.Set(config.ArgQuiet, false)
	w := bufio.NewWriter(&b)
	reg, err := NewPrinterRegistry(w, w)
	assert.NoError(t, err)
	p := reg["text"]
	p.SetStderr(w)
	p.SetStdout(w)
	input := "dummy"

	p.Print(input)
	err = w.Flush()
	assert.NoError(t, err)

	re := regexp.MustCompile(input)
	assert.True(t, re.Match(b.Bytes()))
}

func TestPrinter_UnknownFormatType(t *testing.T) {
	var b bytes.Buffer

	viper.Set(config.ArgOutput, "dummy")
	viper.Set(config.ArgQuiet, false)
	w := bufio.NewWriter(&b)
	_, err := NewPrinterRegistry(w, w)
	assert.Error(t, err)

	assert.True(t, err.Error() == `unknown type format dummy. Hint: use --output json|text`)
}

func TestPrinter_ResultQuiet(t *testing.T) {
	var b bytes.Buffer

	viper.Set(config.ArgOutput, "text")
	viper.Set(config.ArgQuiet, true)
	w := bufio.NewWriter(&b)
	reg, err := NewPrinterRegistry(w, w)
	assert.NoError(t, err)
	p := reg["text"]
	p.SetStderr(w)
	p.SetStdout(w)
	res := Result{
		Message: "command executed",
	}

	p.Print(res)
	err = w.Flush()
	assert.NoError(t, err)

	re := regexp.MustCompile(``)
	assert.True(t, re.Match(b.Bytes()))
}

func TestPrinter_GetStdoutText(t *testing.T) {
	var b bytes.Buffer

	viper.Set(config.ArgOutput, "text")
	viper.Set(config.ArgQuiet, false)
	w := bufio.NewWriter(&b)
	reg, err := NewPrinterRegistry(w, w)
	assert.NoError(t, err)
	out := reg["text"].GetStdout()

	err = w.Flush()
	assert.NoError(t, err)
	assert.True(t, out == w)
}

func TestPrinter_GetStderrText(t *testing.T) {
	var b bytes.Buffer

	viper.Set(config.ArgOutput, "text")
	viper.Set(config.ArgQuiet, false)
	w := bufio.NewWriter(&b)
	reg, err := NewPrinterRegistry(w, w)
	assert.NoError(t, err)
	out := reg["text"].GetStderr()

	err = w.Flush()
	assert.NoError(t, err)
	assert.True(t, out == w)
}
