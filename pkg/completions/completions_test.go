package completions_test

import (
	"fmt"
	"testing"

	"github.com/ionos-cloud/ionosctl/v6/internal/pointer"
	"github.com/ionos-cloud/ionosctl/v6/pkg/completions"
	"github.com/ionos-cloud/ionosctl/v6/pkg/json2table"
	"github.com/stretchr/testify/assert"
)

type innerStruct struct {
	Str     *string  `json:"string,omitempty"`
	Integer *int32   `json:"int32,omitempty"`
	Float   *float64 `json:"float64,omitempty"`
}

type outerStruct struct {
	InnerStructs *[]innerStruct `json:"innerStructs,omitempty"`
}

var (
	innerStructJsonPaths = map[string]string{
		"int32":   "int32",
		"string":  "string",
		"float64": "float64",
	}

	testInnerStruct = innerStruct{
		Str:     pointer.From("test-string"),
		Integer: pointer.From(int32(123)),
		Float:   pointer.From(float64(12.5)),
	}

	testOuterStruct = outerStruct{
		InnerStructs: &[]innerStruct{testInnerStruct, testInnerStruct},
	}

	testConvertedStruct []map[string]interface{}

	expectedOutputBasic             = []string{"123", "123"}
	expectedOutputWithInfo          = []string{"123\t test-string", "123\t test-string"}
	expectedOutputWithMoreInfo      = []string{"123\t test-string 12.5", "123\t test-string 12.5"}
	expectedOutputWithFormattedInfo = []string{"123\t (test-string)", "123\t (test-string)"}
)

func TestCompleter(t *testing.T) {
	var err error

	testConvertedStruct, err = json2table.ConvertJSONToTable("innerStructs", innerStructJsonPaths, testOuterStruct)
	assert.NoError(t, err)

	t.Run("initialize new completer", testNewCompleter)
	t.Run("initialize new completer FAIL", testNewCompleterFail)
	t.Run("add info to completer", testCompleterAddInfo)
	t.Run("add info to completer WITH FORMATTING", testCompleterAddInfoWithFormat)
	t.Run("add info to completer FAIL", testCompleterAddInfoFail)
	t.Run("add MORE info to completer", testCompleterAddMoreInfo)
	t.Run("add MORE info to completer FAIL", testCompleterAddMoreInfoFail)
}

func testNewCompleter(t *testing.T) {
	completer := completions.NewCompleter(testConvertedStruct, "int32")
	assert.NotEqual(t, completer, completions.Completer{})
	assert.Equal(t, expectedOutputBasic, completer.ToString())

	printCompleter(completer)
}

func testNewCompleterFail(t *testing.T) {
	completer := completions.NewCompleter(testConvertedStruct, "intz")
	assert.Equal(t, completer, completions.Completer{})
}

func testCompleterAddInfo(t *testing.T) {
	completer := completions.NewCompleter(testConvertedStruct, "int32").AddInfo("string")
	assert.NotEqual(t, completer, completions.Completer{})
	assert.Equal(t, expectedOutputWithInfo, completer.ToString())

	printCompleter(completer)
}

func testCompleterAddInfoFail(t *testing.T) {
	completer := completions.NewCompleter(testConvertedStruct, "int32").AddInfo("stringz")
	assert.Equal(t, completer, completions.Completer{})
}

func testCompleterAddInfoWithFormat(t *testing.T) {
	completer := completions.NewCompleter(testConvertedStruct, "int32").AddInfo("string", "(%v)")
	assert.NotEqual(t, completer, completions.Completer{})
	assert.Equal(t, expectedOutputWithFormattedInfo, completer.ToString())
	printCompleter(completer)
}

func testCompleterAddMoreInfo(t *testing.T) {
	completer := completions.NewCompleter(testConvertedStruct, "int32").AddInfo("string").AddInfo("float64")
	assert.NotEqual(t, completer, completions.Completer{})
	assert.Equal(t, expectedOutputWithMoreInfo, completer.ToString())

	printCompleter(completer)
}

func testCompleterAddMoreInfoFail(t *testing.T) {
	completer := completions.NewCompleter(testConvertedStruct, "int32").AddInfo("string").AddInfo("floatz")
	assert.Equal(t, completer, completions.Completer{})
}

func printCompleter(completer completions.Completer) {
	for _, out := range completer.ToString() {
		fmt.Println(out)
	}
}
