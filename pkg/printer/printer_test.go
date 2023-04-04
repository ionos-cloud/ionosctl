package printer

import (
	"bufio"
	"bytes"
	"net/http"
	"regexp"
	"testing"

	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	commandExecutedTestMsg = "command executed"
	statusOkTestMsg        = "Status Ok"
)

func TestNewPrinterJSON(t *testing.T) {
	var b bytes.Buffer
	viper.Set(constants.ArgOutput, "json")
	viper.Set(constants.ArgQuiet, false)
	w := bufio.NewWriter(&b)
	_, err := NewPrinterRegistry(w, w, false)
	assert.NoError(t, err)
}

func TestNewPrinterText(t *testing.T) {
	var b bytes.Buffer
	viper.Set(constants.ArgOutput, "text")
	viper.Set(constants.ArgQuiet, false)
	w := bufio.NewWriter(&b)
	_, err := NewPrinterRegistry(w, w, false)
	assert.NoError(t, err)
}

func TestPrinterPrintResultJson(t *testing.T) {
	var (
		b   bytes.Buffer
		str = `{
  "Status": "command executed"
}`
	)
	viper.Set(constants.ArgOutput, "json")
	viper.Set(constants.ArgQuiet, false)
	viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
	w := bufio.NewWriter(&b)
	reg, err := NewPrinterRegistry(w, w, false)
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
	viper.Set(constants.ArgOutput, "json")
	viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
	viper.Set(constants.ArgQuiet, false)
	w := bufio.NewWriter(&b)
	reg, err := NewPrinterRegistry(w, w, false)
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
	viper.Set(constants.ArgOutput, "json")
	viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
	viper.Set(constants.ArgQuiet, false)
	w := bufio.NewWriter(&b)
	reg, err := NewPrinterRegistry(w, w, false)
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
	viper.Set(constants.ArgOutput, "json")
	viper.Set(constants.ArgQuiet, true)
	viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
	w := bufio.NewWriter(&b)
	reg, err := NewPrinterRegistry(w, w, false)
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
	viper.Set(constants.ArgOutput, "json")
	viper.Set(constants.ArgQuiet, false)
	viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
	w := bufio.NewWriter(&b)
	reg, err := NewPrinterRegistry(w, w, false)
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
	res.ApiResponse.Header.Add("location", "https://api.test.ionos.com/cloudapi/v6/requests/123456/status")
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
	viper.Set(constants.ArgOutput, "json")
	viper.Set(constants.ArgQuiet, false)
	viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
	w := bufio.NewWriter(&b)
	reg, err := NewPrinterRegistry(w, w, false)
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
	viper.Set(constants.ArgOutput, "json")
	viper.Set(constants.ArgQuiet, false)
	w := bufio.NewWriter(&b)
	reg, err := NewPrinterRegistry(w, w, false)
	assert.NoError(t, err)
	out := reg["json"].GetStdout()
	err = w.Flush()
	assert.NoError(t, err)
	assert.True(t, out == w)
}

func TestPrinterGetStderrJSON(t *testing.T) {
	var b bytes.Buffer
	viper.Set(constants.ArgOutput, "json")
	viper.Set(constants.ArgQuiet, false)
	w := bufio.NewWriter(&b)
	reg, err := NewPrinterRegistry(w, w, false)
	assert.NoError(t, err)
	out := reg["json"].GetStderr()

	err = w.Flush()
	assert.NoError(t, err)
	assert.True(t, out == w)
}

func TestPrinterPrintResultText(t *testing.T) {
	var b bytes.Buffer
	viper.Set(constants.ArgOutput, "text")
	viper.Set(constants.ArgQuiet, false)
	viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
	w := bufio.NewWriter(&b)
	reg, err := NewPrinterRegistry(w, w, false)
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
	viper.Set(constants.ArgOutput, "text")
	viper.Set(constants.ArgQuiet, false)
	viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
	w := bufio.NewWriter(&b)
	reg, err := NewPrinterRegistry(w, w, false)
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
	res.ApiResponse.Header.Add("location", "https://api.ionos.com/cloudapi/v6/requests/123456/status")
	p.Print(res)
	err = w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`123456`)
	assert.True(t, re.Match(b.Bytes()))
}

func TestPrinterPrintResultTextResource(t *testing.T) {
	var b bytes.Buffer
	viper.Set(constants.ArgOutput, "text")
	viper.Set(constants.ArgQuiet, false)
	viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
	w := bufio.NewWriter(&b)
	reg, err := NewPrinterRegistry(w, w, false)
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
	viper.Set(constants.ArgOutput, "text")
	viper.Set(constants.ArgQuiet, false)
	viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
	w := bufio.NewWriter(&b)
	reg, err := NewPrinterRegistry(w, w, false)
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
	viper.Set(constants.ArgOutput, "text")
	viper.Set(constants.ArgQuiet, false)
	viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
	w := bufio.NewWriter(&b)
	reg, err := NewPrinterRegistry(w, w, false)
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
	viper.Set(constants.ArgOutput, "text")
	viper.Set(constants.ArgQuiet, false)
	viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
	w := bufio.NewWriter(&b)
	reg, err := NewPrinterRegistry(w, w, false)
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

func TestPrinterPrintResultTextKeyValueNoHeaders(t *testing.T) {
	var (
		b      bytes.Buffer
		tabwrt = "123   dummy   true   1.230000\n"
	)
	viper.Set(constants.ArgOutput, "text")
	viper.Set(constants.ArgQuiet, false)
	viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
	w := bufio.NewWriter(&b)
	reg, err := NewPrinterRegistry(w, w, true)
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
	assert.Truef(t, re.Match(b.Bytes()), `"%v" != "%v"`, b.Bytes(), []byte(tabwrt))
}

func TestPrinterPrintDefaultText(t *testing.T) {
	var b bytes.Buffer
	viper.Set(constants.ArgOutput, "text")
	viper.Set(constants.ArgQuiet, false)
	viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
	w := bufio.NewWriter(&b)
	reg, err := NewPrinterRegistry(w, w, false)
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
	viper.Set(constants.ArgOutput, "dummy")
	viper.Set(constants.ArgQuiet, false)
	w := bufio.NewWriter(&b)
	_, err := NewPrinterRegistry(w, w, false)
	assert.Error(t, err)
	assert.True(t, err.Error() == `unknown type format dummy. Hint: use --output json|text`)
}

func TestPrinterResultQuiet(t *testing.T) {
	var b bytes.Buffer
	viper.Set(constants.ArgOutput, "text")
	viper.Set(constants.ArgQuiet, true)
	viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
	w := bufio.NewWriter(&b)
	reg, err := NewPrinterRegistry(w, w, false)
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
	viper.Set(constants.ArgOutput, "text")
	viper.Set(constants.ArgQuiet, false)
	w := bufio.NewWriter(&b)
	reg, err := NewPrinterRegistry(w, w, false)
	assert.NoError(t, err)
	out := reg["text"].GetStdout()
	err = w.Flush()
	assert.NoError(t, err)
	assert.True(t, out == w)
}

func TestPrinterGetStderrText(t *testing.T) {
	var b bytes.Buffer
	viper.Set(constants.ArgOutput, "text")
	viper.Set(constants.ArgQuiet, false)
	w := bufio.NewWriter(&b)
	reg, err := NewPrinterRegistry(w, w, false)
	assert.NoError(t, err)
	out := reg["text"].GetStderr()
	err = w.Flush()
	assert.NoError(t, err)
	assert.True(t, out == w)
}

func TestGetRequestId(t *testing.T) {
	id, err := GetRequestId(constants.DefaultApiURL)
	assert.Error(t, err)
	assert.Empty(t, id)
	id, err = GetRequestId("https://api.test.ionos.com/cloudapi/v6/requests/test/status")
	assert.NoError(t, err)
	assert.Equal(t, "test", id)
}
