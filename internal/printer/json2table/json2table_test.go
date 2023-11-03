package json2table_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table"
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

	expectedResultBasicStruct = []map[string]interface{}{
		{
			"int32":  float64(testInt32), // because unmarshalling into an interface makes all numbers into float64
			"string": testStr,
		},
	}

	expectedResultComplexStruct = []map[string]interface{}{
		{
			"innerInts32":  []interface{}{float64(testInt32), float64(testInt32)},
			"innerStrings": []interface{}{testStr, testStr},
			"strings":      []interface{}{testStr, testStr, testStr},
			"int64":        float64(testInt64),
			"string":       testStr,
			"float64":      testFloat64,
		},
	}

	expectedResultInnerBasicStructs = []map[string]interface{}{
		{
			"int32":  float64(testInt32),
			"string": testStr,
		},
		{
			"int32":  float64(testInt32),
			"string": testStr,
		},
	}

	wrongPaths = map[string]string{
		"random": "random.path",
		"path":   "path.random",
	}

	expectedResultWrongPaths = []map[string]interface{}{
		{
			"random": nil,
			"path":   nil,
		},
	}
)

func TestConvertJSONToText(t *testing.T) {
	t.Run("Convert JSON to TEXT with basic struct", testConvertJSONToTextWithBasicStruct)
	t.Run("Convert JSON to TEXT with complex struct", testConvertJSONToTextWithComplexStruct)
	t.Run("Convert JSON to TEXT with inner basic structs", testConvertJSONToTextWithInnerBasicStructs)
	t.Run("Convert JSON to TEXT with wrong paths", testConvertJSONToTextWithWrongPaths)
	t.Run("FAIL Convert JSON to TEXT with empty paths", testFailConvertJSONToTextWithEmptyPaths)
	t.Run("FAIL Convert JSON to TEXT with wrong root", testFailConvertJSONToTextWithWrongRoot)
	t.Run("FAIL Convert JSON to TEXT with wrong root destination", testFailConvertJSONToTextWithWrongRootDestination)
	t.Run("FAIL Convert JSON to TEXT with unsupported JSON value", testFailConvertJSONToTextWithUnsupportedJSONValue)
	t.Run("FAIL Convert JSON to TEXT with unsupported JSON type", testFailConvertJSONToTextWithUnsupportedJSONType)
	t.Run("FAIL Convert JSON to TEXT with empty JSON", testFailConvertJSONToTextWithEmptyJSON)
}

func testConvertJSONToTextWithBasicStruct(t *testing.T) {
	res, err := json2table.ConvertJSONToTable("", innerStructJsonPaths, testInnerStruct)
	assert.NoError(t, err)
	assert.ElementsMatch(t, expectedResultBasicStruct, res)
}

func testConvertJSONToTextWithComplexStruct(t *testing.T) {
	res, err := json2table.ConvertJSONToTable("", outerStructJsonPaths, testOuterStruct)
	assert.NoError(t, err)
	assert.ElementsMatch(t, expectedResultComplexStruct, res)
}

func testConvertJSONToTextWithInnerBasicStructs(t *testing.T) {
	res, err := json2table.ConvertJSONToTable("innerStructs", innerStructJsonPaths, testOuterStruct)
	assert.NoError(t, err)
	assert.ElementsMatch(t, expectedResultInnerBasicStructs, res)
}

func testFailConvertJSONToTextWithWrongRoot(t *testing.T) {
	_, err := json2table.ConvertJSONToTable("random.root.path", innerStructJsonPaths, testOuterStruct)
	if assert.Error(t, err) {
		fmt.Println(err)
	}
}

func testFailConvertJSONToTextWithWrongRootDestination(t *testing.T) {
	_, err := json2table.ConvertJSONToTable("int64", innerStructJsonPaths, testOuterStruct)
	if assert.Error(t, err) {
		fmt.Println(err)
	}
}

func testFailConvertJSONToTextWithUnsupportedJSONValue(t *testing.T) {
	_, err := json2table.ConvertJSONToTable("", innerStructJsonPaths, math.Inf(1))
	if assert.Error(t, err) {
		fmt.Println(err)
	}
}

func testFailConvertJSONToTextWithUnsupportedJSONType(t *testing.T) {
	_, err := json2table.ConvertJSONToTable("", innerStructJsonPaths, make(chan int))
	if assert.Error(t, err) {
		fmt.Println(err)
	}
}

func testFailConvertJSONToTextWithEmptyJSON(t *testing.T) {
	_, err := json2table.ConvertJSONToTable("", innerStructJsonPaths, nil)
	if assert.Error(t, err) {
		fmt.Println(err)
	}
}

func testFailConvertJSONToTextWithEmptyPaths(t *testing.T) {
	_, err := json2table.ConvertJSONToTable("", nil, testInnerStruct)
	if assert.Error(t, err) {
		fmt.Println(err)
	}
}

func testConvertJSONToTextWithWrongPaths(t *testing.T) {
	res, err := json2table.ConvertJSONToTable("", wrongPaths, testInnerStruct)
	assert.NoError(t, err)
	assert.ElementsMatch(t, expectedResultWrongPaths, res)
}
