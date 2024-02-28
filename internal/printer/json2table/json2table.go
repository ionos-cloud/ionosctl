package json2table

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/Jeffail/gabs/v2"
	"github.com/ionos-cloud/ionosctl/v6/pkg/functional"
)

// ConvertJSONToTable extracts data from a JSON or struct object using specified paths.
//
// columnPathMappingPrefix: Points to a specific location in the JSON or struct object from where
// extraction should begin. If left empty, extraction will start from the root of the object.
//
// columnPathMapping: A map where each key represents the desired column name in the output
// table, and each value represents a gabs path to extract the data within the given JSON or struct.
//
// sourceData: JSON or struct from which data extraction should take place.
//
// Returns a slice of maps. Each map represents a row in the table, with each key-value
// pair in the map being equivalent to a column name and its value.
func ConvertJSONToTable(columnPathMappingPrefix string, columnPathMapping map[string]string, sourceData interface{}) ([]map[string]interface{}, error) {
	if sourceData == nil {
		return nil, fmt.Errorf("provided object cannot be nil")
	}

	var res = make([]map[string]interface{}, 0)

	objs, err := traverseJSONRoot(columnPathMappingPrefix, sourceData)
	if err != nil {
		return nil, fmt.Errorf("failed traversing the root path %s: %w", columnPathMappingPrefix, err)
	}

	if columnPathMapping == nil || len(columnPathMapping) == 0 {
		return nil, fmt.Errorf("json paths must not be empty/nil")
	}

	for _, obj := range objs {
		mappedObj := make(map[string]interface{})

		for k, v := range columnPathMapping {
			objData := obj.Path(v)
			mappedObj[k] = objData.Data()
		}

		res = append(res, mappedObj)
	}

	return res, nil
}

// traverseJSONRoot reaches a specific location in the JSON or struct object. The prefix should point to an array of
// objects. If the prefix points to an array within the elements of another array, it will flatten the resulting 2D
// array.
//
// columnPathMappingPrefix: Points to a specific location in the JSON or struct object from where
// extraction should begin. If left empty, extraction will start from the root of the object.
//
// sourceData: JSON or struct from which data extraction should take place.
//
// Returns a slice of gabs.Container. Each gabs.Container is a single object.
func traverseJSONRoot(columnPathMappingPrefix string, sourceData interface{}) ([]*gabs.Container, error) {
	jsonObj, err := json.Marshal(sourceData)
	if err != nil {
		return nil, err
	}
	if reflect.DeepEqual(jsonObj, []byte{'[', ']'}) {
		return []*gabs.Container{}, nil
	}

	parsedObj, err := gabs.ParseJSON(jsonObj)
	if err != nil {
		return nil, err
	}

	if columnPathMappingPrefix == "" {
		if reflect.TypeOf(sourceData).Kind() == reflect.Slice {
			return parsedObj.Children(), nil
		} else {
			return []*gabs.Container{parsedObj}, nil
		}
	}

	if !parsedObj.ExistsP(columnPathMappingPrefix) {
		return nil, fmt.Errorf("'%s' does not exist in [%s]",
			columnPathMappingPrefix,
			strings.Join(functional.KeysOfMap(parsedObj.ChildrenMap()), ", "))
	}

	parsedObj = parsedObj.Path(columnPathMappingPrefix)

	tempChildren := parsedObj.Children()

	if tempChildren == nil {
		return nil, fmt.Errorf("root path does not lead to an array in object: %s", columnPathMappingPrefix)
	}

	var children []*gabs.Container
	for _, child := range tempChildren {
		if _, ok := child.Data().([]interface{}); ok {
			children = append(children, child.Children()...)
		} else {
			return tempChildren, nil
		}
	}

	return children, nil
}
