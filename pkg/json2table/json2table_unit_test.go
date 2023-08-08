//go:build unit
// +build unit

package json2table_test

import (
	"encoding/json"
	"testing"

	"github.com/ionos-cloud/ionosctl/v6/pkg/json2table"
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
			"int32":  testInt32,
			"string": testStr,
		},
	}

	expectedResultComplexStruct = []map[string]interface{}{
		{
			"innerInts32":  []int32{testInt32, testInt32},
			"innerStrings": []string{testStr, testStr},
			"strings":      []string{testStr, testStr, testStr},
			"int64":        testInt64,
			"string":       testStr,
			"float64":      testFloat64,
		},
	}

	expectedResultInnerBasicStructs = []map[string]interface{}{
		{
			"int32":  testInt32,
			"string": testStr,
		},
		{
			"int32":  testInt32,
			"string": testStr,
		},
	}
)

func TestConvertJSONToText(t *testing.T) {
	t.Run("Convert JSON to TEXT with basic struct", testConvertJSONToTextWithBasicStruct)
	t.Run("Convert JSON to TEXT with complex struct", testConvertJSONToTextWithComplexStruct)
	t.Run("Convert JSON to TEXT with inner basic structs", testConvertJSONToTextWithInnerBasicStructs)
}

func testConvertJSONToTextWithBasicStruct(t *testing.T) {
	res, err := json2table.ConvertJSONToText("", innerStructJsonPaths, testInnerStruct)
	assert.NoError(t, err)
	assert.Equal(t, true, compareMapSlices(res, expectedResultBasicStruct))
}

func testConvertJSONToTextWithComplexStruct(t *testing.T) {
	res, err := json2table.ConvertJSONToText("", outerStructJsonPaths, testOuterStruct)
	assert.NoError(t, err)
	assert.Equal(t, true, compareMapSlices(res, expectedResultComplexStruct))
}

func testConvertJSONToTextWithInnerBasicStructs(t *testing.T) {
	res, err := json2table.ConvertJSONToText("innerStructs", innerStructJsonPaths, testOuterStruct)
	assert.NoError(t, err)
	assert.Equal(t, true, compareMapSlices(res, expectedResultInnerBasicStructs))
}

func compareMapSlices(mapSlice1 []map[string]interface{}, mapSlice2 []map[string]interface{}) bool {
	m1, err := json.Marshal(mapSlice1)
	if err != nil {
		return false
	}

	m2, err := json.Marshal(mapSlice2)
	if err != nil {
		return false
	}

	return string(m1) == string(m2)
}
