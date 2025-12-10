// go
package jsontabwriter_test

import (
	"encoding/json"
	"testing"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func unmarshal(t *testing.T, s string) interface{} {
	var v interface{}
	err := json.Unmarshal([]byte(s), &v)
	assert.NoError(t, err)
	return v
}

func TestQuery_SelectFirstID_APIJSON(t *testing.T) {
	viper.Reset()
	viper.Set(constants.ArgOutput, jsontabwriter.APIFormat)
	viper.Set(constants.FlagQuery, "items[0].id")

	src := unmarshal(t, `{"items":[{"id":"a1"},{"id":"a2"}]}`)
	out, err := jsontabwriter.GenerateOutput("", nil, src, nil)
	assert.NoError(t, err)
	assert.Equal(t, "\"a1\"\n", out)
}

func TestQuery_FilterByVersion_LegacyJSON(t *testing.T) {
	viper.Reset()
	viper.Set(constants.ArgOutput, jsontabwriter.JSONFormat)
	viper.Set(constants.FlagQuery, "items[?properties.version>=`2`].id")

	// Two paged responses -> legacy merger will produce a single items array
	pageData := []interface{}{
		unmarshal(t, `{"items":[{"id":"a1","properties":{"version":1}},{"id":"a2","properties":{"version":3}}]}`),
		unmarshal(t, `{"items":[{"id":"a3","properties":{"version":2}}]}`),
	}

	out, err := jsontabwriter.GenerateOutput("", nil, pageData, nil)
	assert.NoError(t, err)
	// Expect ids with version >= 2: a2, a3
	expected := "[\n  \"a2\",\n  \"a3\"\n]\n"
	assert.Equal(t, expected, out)
}

func TestQuery_ProjectionMatrix_APIJSON(t *testing.T) {
	viper.Reset()
	viper.Set(constants.ArgOutput, jsontabwriter.APIFormat)
	viper.Set(constants.FlagQuery, "items[?properties.version>=`2`].[properties.name, properties.location]")

	src := unmarshal(t, `{
    "items":[
      {"properties":{"version":1,"name":"N1","location":"loc-a"}},
      {"properties":{"version":2,"name":"N2","location":"loc-b"}},
      {"properties":{"version":5,"name":"N5","location":"loc-c"}}
    ]}`)

	out, err := jsontabwriter.GenerateOutput("", nil, src, nil)
	assert.NoError(t, err)
	expected := "[\n  [\n    \"N2\",\n    \"loc-b\"\n  ],\n  [\n    \"N5\",\n    \"loc-c\"\n  ]\n]\n"
	assert.Equal(t, expected, out)
}

func TestQuery_InvalidExpression_Error(t *testing.T) {
	viper.Reset()
	viper.Set(constants.ArgOutput, jsontabwriter.APIFormat)
	viper.Set(constants.FlagQuery, "items[?properties.version>=") // malformed

	src := unmarshal(t, `{"items":[{"properties":{"version":1}}]}`)

	_, err := jsontabwriter.GenerateOutput("", nil, src, nil)
	assert.Error(t, err)
}

func TestQuery_ScalarResult_APIJSON(t *testing.T) {
	viper.Reset()
	viper.Set(constants.ArgOutput, jsontabwriter.APIFormat)
	viper.Set(constants.FlagQuery, "items|length(@)")

	src := unmarshal(t, `{"items":[{"id":"x"},{"id":"y"},{"id":"z"}]}`)

	out, err := jsontabwriter.GenerateOutput("", nil, src, nil)
	assert.NoError(t, err)
	// Length is 3
	assert.Equal(t, "3\n", out)
}

func TestQuery_NestedField_APIJSON(t *testing.T) {
	viper.Reset()
	viper.Set(constants.ArgOutput, jsontabwriter.APIFormat)
	viper.Set(constants.FlagQuery, "items[].metadata.lastModifiedDate")

	src := unmarshal(t, `{
    "items":[
      {"metadata":{"lastModifiedDate":"2025-09-26T10:27:55Z"}},
      {"metadata":{"lastModifiedDate":"2024-01-01T00:00:00Z"}}
    ]}`)

	out, err := jsontabwriter.GenerateOutput("", nil, src, nil)
	assert.NoError(t, err)
	expected := "[\n  \"2025-09-26T10:27:55Z\",\n  \"2024-01-01T00:00:00Z\"\n]\n"
	assert.Equal(t, expected, out)
}
