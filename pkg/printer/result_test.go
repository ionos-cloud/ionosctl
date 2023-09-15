package printer

import (
	"bytes"
	"testing"

	"github.com/ionos-cloud/ionosctl/v6/pkg/tabheaders"
	"github.com/stretchr/testify/assert"
)

func TestGetHeaders(t *testing.T) {
	type args struct {
		allColumns     []string
		defaultColumns []string
		customColumns  []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "no custom columns",
			args: struct {
				allColumns     []string
				defaultColumns []string
				customColumns  []string
			}{
				allColumns:     []string{"col1", "col2"},
				defaultColumns: []string{"col1"},
				customColumns:  nil,
			},
			want: []string{"col1"},
		},
		{
			name: "all columns as custom",
			args: struct {
				allColumns     []string
				defaultColumns []string
				customColumns  []string
			}{
				allColumns:     []string{"col1", "col2"},
				defaultColumns: []string{"col1"},
				customColumns:  []string{"col1", "col2"},
			},
			want: []string{"col1", "col2"},
		},
		{
			name: "custom order",
			args: struct {
				allColumns     []string
				defaultColumns []string
				customColumns  []string
			}{
				allColumns:     []string{"col1", "col2", "col3", "col4"},
				defaultColumns: []string{"col1"},
				customColumns:  []string{"col2", "col1", "col4"},
			},
			want: []string{"col2", "col1", "col4"},
		},
		{
			name: "some invalid custom cols",
			args: struct {
				allColumns     []string
				defaultColumns []string
				customColumns  []string
			}{
				allColumns:     []string{"col2", "col3"},
				defaultColumns: []string{"col2"},
				customColumns:  []string{"col2", "col1", "col4", "col3"},
			},
			want: []string{"col2", "col3"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, tabheaders.GetHeaders(tt.args.allColumns, tt.args.defaultColumns, tt.args.customColumns), "GetHeaders(%v, %v, %v)", tt.args.allColumns, tt.args.defaultColumns, tt.args.customColumns)
		})
	}
}

func TestPrintText(t *testing.T) {
	buf := new(bytes.Buffer)
	headers := []string{
		"String", "Int", "Float64", "Bool", "StringSlice",
	}
	kvmap := []map[string]interface{}{
		{
			"String":      "string",
			"Int":         int(1),
			"Float64":     float64(1.123),
			"Bool":        true,
			"StringSlice": []string{"a", "b"},
		},
	}

	assert.NoError(t, printText(buf, headers, kvmap, false))

	expectOut := `String   Int   Float64    Bool   StringSlice
string   1     1.123000   true   a,b
`
	assert.Equal(t, expectOut, buf.String())
}
