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

var (
	commandExecutedTestMsg = "command executed"
	statusOkTestMsg        = "Status Ok"
)

func TestNewPrinterJSON(t *testing.T) {
	var b bytes.Buffer

	viper.Set(config.ArgOutput, "json")
	viper.Set(config.ArgQuiet, false)
	w := bufio.NewWriter(&b)
	_, err := NewPrinterRegistry(w, w)
	assert.NoError(t, err)
}

func TestNewPrinterText(t *testing.T) {
	var b bytes.Buffer

	viper.Set(config.ArgOutput, "text")
	viper.Set(config.ArgQuiet, false)
	w := bufio.NewWriter(&b)
	_, err := NewPrinterRegistry(w, w)
	assert.NoError(t, err)
}

func TestPrinterPrintResultJson(t *testing.T) {
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
		Message: commandExecutedTestMsg,
	}

	p.Print(res)
	err = w.Flush()
	assert.NoError(t, err)

	re := regexp.MustCompile(str)
	assert.True(t, re.Match(b.Bytes()))
}

func TestPrinterPrintStandardResultJson(t *testing.T) {
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

func TestPrinterPrintDefaultJson(t *testing.T) {
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

	p.Print(commandExecutedTestMsg)
	err = w.Flush()
	assert.NoError(t, err)

	re := regexp.MustCompile(str)
	assert.True(t, re.Match(b.Bytes()))
}

func TestPrinterPrintDefaultQuietJson(t *testing.T) {
	var b bytes.Buffer

	viper.Set(config.ArgOutput, "json")
	viper.Set(config.ArgQuiet, true)
	w := bufio.NewWriter(&b)
	reg, err := NewPrinterRegistry(w, w)
	assert.NoError(t, err)
	p := reg["json"]
	p.SetStderr(w)
	p.SetStdout(w)

	p.Print(commandExecutedTestMsg)
	err = w.Flush()
	assert.NoError(t, err)

	re := regexp.MustCompile(``)
	assert.True(t, re.Match(b.Bytes()))
}

func TestPrinterResultJsonRequestId(t *testing.T) {
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
				Message: statusOkTestMsg,
				Response: &http.Response{
					Header: map[string][]string{},
				},
			},
		},
	}
	res.ApiResponse.Header.Add("location", "https://api.test.ionos.com/cloudapi/v5/requests/123456/status")

	p.Print(res)
	err = w.Flush()
	assert.NoError(t, err)

	re := regexp.MustCompile(str)
	assert.True(t, re.Match(b.Bytes()))
}

func TestPrinterResultJsonRequestIdErr(t *testing.T) {
	var (
		b      bytes.Buffer
		strErr = `https://api.ionos.com/servers/123456 does not contain requestId`
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
				Message: statusOkTestMsg,
				Response: &http.Response{
					Header: map[string][]string{},
				},
			},
		},
	}
	res.ApiResponse.Header.Add("location", "https://api.ionos.com/servers/123456")

	err = p.Print(res)
	assert.Error(t, err)
	assert.EqualError(t, err, strErr)
}

func TestPrinterGetStdoutJSON(t *testing.T) {
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

func TestPrinterGetStderrJSON(t *testing.T) {
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

func TestPrinterPrintResultText(t *testing.T) {
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
		Message: commandExecutedTestMsg,
	}

	p.Print(res)
	err = w.Flush()
	assert.NoError(t, err)

	re := regexp.MustCompile(commandExecutedTestMsg)
	assert.True(t, re.Match(b.Bytes()))
}

func TestPrinterPrintResultTextRequestId(t *testing.T) {
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
				Message: statusOkTestMsg,
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

func TestPrinterPrintResultTextResource(t *testing.T) {
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
		Resource:       "datacenter",
		Verb:           "create",
		WaitForRequest: false,
	}
	p.Print(res)
	err = w.Flush()
	assert.NoError(t, err)

	re := regexp.MustCompile(`Command datacenter create has been successfully executed`)
	assert.True(t, re.Match(b.Bytes()))
}

func TestPrinterPrintResultTextWaitResource(t *testing.T) {
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
		Resource:       "datacenter",
		Verb:           "create",
		WaitForRequest: true,
	}
	p.Print(res)
	err = w.Flush()
	assert.NoError(t, err)

	re := regexp.MustCompile(`Command datacenter create & wait have been successfully executed`)
	assert.True(t, re.Match(b.Bytes()))
}

func TestPrinterPrintResultTextWaitStateResource(t *testing.T) {
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
		Resource:       "datacenter",
		Verb:           "create",
		WaitForRequest: false,
		WaitForState:   true,
	}
	p.Print(res)
	err = w.Flush()
	assert.NoError(t, err)

	re := regexp.MustCompile(`Command datacenter create & wait have been successfully executed`)
	assert.True(t, re.Match(b.Bytes()))
}

func TestPrinterPrintResultTextKeyValue(t *testing.T) {
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

func TestPrinterPrintDefaultText(t *testing.T) {
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

func TestPrinterUnknownFormatType(t *testing.T) {
	var b bytes.Buffer

	viper.Set(config.ArgOutput, "dummy")
	viper.Set(config.ArgQuiet, false)
	w := bufio.NewWriter(&b)
	_, err := NewPrinterRegistry(w, w)
	assert.Error(t, err)

	assert.True(t, err.Error() == `unknown type format dummy. Hint: use --output json|text`)
}

func TestPrinterResultQuiet(t *testing.T) {
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
		Message: commandExecutedTestMsg,
	}

	p.Print(res)
	err = w.Flush()
	assert.NoError(t, err)

	re := regexp.MustCompile(``)
	assert.True(t, re.Match(b.Bytes()))
}

func TestPrinterGetStdoutText(t *testing.T) {
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

func TestPrinterGetStderrText(t *testing.T) {
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

func TestGetRequestId(t *testing.T) {
	id, err := GetRequestId(config.DefaultApiURL)
	assert.Error(t, err)
	assert.Empty(t, id)
	id, err = GetRequestId("https://api.test.ionos.com/cloudapi/v5/requests/test/status")
	assert.NoError(t, err)
	assert.Equal(t, "test", id)
}
