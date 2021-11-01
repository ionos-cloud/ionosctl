/*
	This is used for supporting completion in the CLI list requests --filters options
*/
package completer

import (
	"bytes"
	"reflect"
	"strings"

	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
)

func DataCentersFilters() []string {
	return getPropertiesName(ionoscloud.DatacenterProperties{}, ionoscloud.DatacenterElementMetadata{})
}

func getPropertiesName(params ...interface{}) []string {
	var argKind reflect.Kind
	properties := make([]string, 0)
	for _, param := range params {
		arg := reflect.TypeOf(param)
		if arg.Kind() == reflect.Ptr {
			argKind = arg.Elem().Kind()
		} else {
			argKind = arg.Kind()
		}
		if argKind == reflect.Struct {
			for i := 0; i < arg.NumField(); i++ {
				if arg.Field(i).Type.Kind() == reflect.Ptr {
					argKind = arg.Field(i).Type.Elem().Kind()
					if argKind != reflect.Slice {
						properties = append(properties, makeFirstLowerCase(arg.Field(i).Name))
					}
				}
			}
		}
	}
	return properties
}

func makeFirstLowerCase(s string) string {
	if len(s) < 2 {
		return strings.ToLower(s)
	}
	bts := []byte(s)
	return string(bytes.Join([][]byte{
		bytes.ToLower([]byte{bts[0]}), bts[1:],
	}, nil))
}
