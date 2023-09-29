package tabheaders

import (
	"fmt"
	"strings"

	"github.com/ionos-cloud/ionosctl/v6/internal/functional"
	"golang.org/x/exp/slices"
)

// GetHeadersAllDefault is like GetHeaders, but defaultColumns is same as allColumns.
// Useful for resources with small print table
func GetHeadersAllDefault(allColumns []string, customColumns []string) []string {
	return GetHeaders(allColumns, allColumns, customColumns)
}

// GetHeaders takes all columns of a resource and the value of the columns flag,
// returns the headers of the table. (Some legacy code might refer to these headers as "Columns")
func GetHeaders(allColumns []string, defaultColumns []string, customColumns []string) []string {
	if len(customColumns) > 0 && customColumns[0] == "all" {
		return GetHeaders(allColumns, defaultColumns, allColumns)
	}

	if customColumns == nil {
		return defaultColumns
	}

	allColumnsLowercase := functional.Map(allColumns, func(x string) string {
		return strings.ToLower(x)
	})

	var validCustomColumns []string
	for _, c := range customColumns {
		if idx := slices.Index(allColumnsLowercase, strings.ToLower(c)); idx != -1 {
			validCustomColumns = append(validCustomColumns, allColumns[idx])
		}
	}

	if len(validCustomColumns) == 0 {
		return defaultColumns
	}

	return validCustomColumns
}

func ColsMessage(cols []string) string {
	return fmt.Sprintf("Set of columns to be printed on output \nAvailable columns: %v", cols)
}
