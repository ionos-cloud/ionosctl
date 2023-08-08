//go:build unit
// +build unit

package jsontabwriter_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/jsontabwriter"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

type innerStruct struct {
	Str     *string `json:"string,omitempty"`
	Integer *int32  `json:"int32,omitempty"`
}

type outerStruct struct {
	Str               *string        `json:"string,omitempty"`
	Integer           *int64         `json:"int64,omitempty"`
	FloatingPoint     *float64       `json:"float64,omitempty"`
	Strs              *[]string      `json:"strings,omitempty"`
	InnerStructs      *[]innerStruct `json:"innerStructs,omitempty"`
	InnerStructSingle *innerStruct   `json:"innerStruct,omitempty"`
}

var (
	testStr     = "test-string"
	testInt32   = int32(123)
	testInt64   = int64(123543)
	testFloat64 = float64(12.5)

	innerStructJsonPaths = map[string]string{
		"int32":  "int32",
		"string": "string",
	}

	outerStructJsonPaths = map[string]string{
		"string":       "string",
		"int64":        "int64",
		"float64":      "float64",
		"strings":      "strings",
		"innerStrings": "innerStructs.*.string",
		"innerInts32":  "innerStructs.*.int32",
	}

	testInnerStruct = innerStruct{
		Str:     &testStr,
		Integer: &testInt32,
	}

	testOuterStruct = outerStruct{
		Str:               &testStr,
		Integer:           &testInt64,
		FloatingPoint:     &testFloat64,
		Strs:              &[]string{testStr, testStr, testStr},
		InnerStructs:      &[]innerStruct{testInnerStruct, testInnerStruct},
		InnerStructSingle: &testInnerStruct,
	}

	testFormatString = "Testing with %v, %v"
)

func TestGenerateOutput(t *testing.T) {
	t.Run("Generate TEXT output with basic struct", testGenerateTextOutputWithBasicStruct)
	t.Run("Generate JSON output with basic struct", testGenerateJSONOutputWithBasicStruct)
	t.Run("Generate TEXT output with complex struct", testGenerateTextOutputWithComplexStruct)
	t.Run("Generate JSON output with complex struct", testGenerateJSONOutputWithComplexStruct)
	t.Run("Generate TEXT output with inner basic structs", testGenerateTextOutputWithInnerBasicStructs)
	t.Run("Generate JSON output with inner basic structs", testGenerateJSONOutputWithInnerBasicStructs)
}

func TestGenerateVerboseOutput(t *testing.T) {
	t.Run("Generate TEXT verbose output with verbosity SET", testGenerateTextVerboseOutputSet)
	t.Run("Generate JSON verbose output with verbosity SET", testGenerateJSONVerboseOutputSet)
	t.Run("Generate TEXT verbose output with verbosity NOT SET", testGenerateTextVerboseOutputNotSet)
	t.Run("Generate JSON verbose output with verbosity NOT SET", testGenerateJSONVerboseOutputNotSet)
}

func TestGenerateLogOutput(t *testing.T) {
	t.Run("Generate TEXT log output", testGenerateTextLogOutput)
	t.Run("Generate JSON log output", testGenerateJSONLogOutput)
}

func testGenerateTextOutputWithBasicStruct(t *testing.T) {
	viper.Reset()

	viper.Set(constants.ArgOutput, "text")

	out, err := jsontabwriter.GenerateOutput("", innerStructJsonPaths, testInnerStruct, []string{"int32", "string"})
	assert.NoError(t, err)

	fmt.Println(out)
}

func testGenerateJSONOutputWithBasicStruct(t *testing.T) {
	viper.Reset()

	viper.Set(constants.ArgOutput, "json")

	out, err := jsontabwriter.GenerateOutput("", innerStructJsonPaths, testInnerStruct, []string{"int32", "string"})
	assert.NoError(t, err)

	fmt.Println(out)
}

func testGenerateTextOutputWithComplexStruct(t *testing.T) {
	viper.Reset()

	viper.Set(constants.ArgOutput, "text")

	out, err := jsontabwriter.GenerateOutput("", outerStructJsonPaths, testOuterStruct,
		[]string{"int64", "string", "float64", "strings", "innerStrings", "innerInts32"})
	assert.NoError(t, err)

	fmt.Println(out)
}

func testGenerateJSONOutputWithComplexStruct(t *testing.T) {
	viper.Reset()

	viper.Set(constants.ArgOutput, "json")

	out, err := jsontabwriter.GenerateOutput("", outerStructJsonPaths, testOuterStruct,
		[]string{"int32", "string", "float64", "strings", "innerStrings", "innerInts32"})
	assert.NoError(t, err)

	fmt.Println(out)
}

func testGenerateTextOutputWithInnerBasicStructs(t *testing.T) {
	viper.Reset()

	viper.Set(constants.ArgOutput, "text")
	out, err := jsontabwriter.GenerateOutput("innerStructs", innerStructJsonPaths, testOuterStruct, []string{"int32", "string"})
	assert.NoError(t, err)

	fmt.Println(out)
}

func testGenerateJSONOutputWithInnerBasicStructs(t *testing.T) {
	viper.Reset()

	viper.Set(constants.ArgOutput, "json")

	out, err := jsontabwriter.GenerateOutput("innerStructs", innerStructJsonPaths, testOuterStruct, []string{"int32", "string"})
	assert.NoError(t, err)

	fmt.Println(out)
}

func testGenerateTextVerboseOutputSet(t *testing.T) {
	viper.Reset()

	viper.Set(constants.ArgVerbose, true)
	viper.Set(constants.ArgOutput, "text")
	n, err := fmt.Printf(jsontabwriter.GenerateVerboseOutput(testFormatString, testInt32, testStr))
	assert.NoError(t, err)
	assert.Equal(t, len(fmt.Sprintf("[INFO] "+testFormatString+"\n", testInt32, testStr)), n)
}

func testGenerateJSONVerboseOutputSet(t *testing.T) {
	j, _ := json.MarshalIndent(fmt.Sprintf("[INFO] "+testFormatString, testInt32, testStr), "", "\t")

	viper.Reset()

	viper.Set(constants.ArgVerbose, true)
	viper.Set(constants.ArgOutput, "json")
	n, err := fmt.Printf(jsontabwriter.GenerateVerboseOutput(testFormatString, testInt32, testStr))
	assert.NoError(t, err)
	assert.Equal(t, len(string(j)+"\n"), n)
}

func testGenerateTextVerboseOutputNotSet(t *testing.T) {
	viper.Reset()

	viper.Set(constants.ArgOutput, "text")
	n, err := fmt.Printf(jsontabwriter.GenerateVerboseOutput(testFormatString, testInt32, testStr))
	assert.NoError(t, err)
	assert.Equal(t, 0, n)
}

func testGenerateJSONVerboseOutputNotSet(t *testing.T) {
	viper.Reset()

	viper.Set(constants.ArgOutput, "json")
	n, err := fmt.Printf(jsontabwriter.GenerateVerboseOutput(testFormatString, testInt32, testStr))
	assert.NoError(t, err)
	assert.Equal(t, 0, n)
}

func testGenerateTextLogOutput(t *testing.T) {
	viper.Reset()

	viper.Set(constants.ArgOutput, "text")
	n, err := fmt.Printf(jsontabwriter.GenerateLogOutput(testStr))
	assert.NoError(t, err)
	assert.Equal(t, len(testStr+"\n"), n)
}

func testGenerateJSONLogOutput(t *testing.T) {
	j, _ := json.MarshalIndent(testStr, "", "\t")

	viper.Reset()

	viper.Set(constants.ArgOutput, "json")
	n, err := fmt.Printf(jsontabwriter.GenerateLogOutput(testStr))
	assert.NoError(t, err)
	assert.Equal(t, len(string(j)+"\n"), n)
}
