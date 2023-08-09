package json2table

import (
	"encoding/json"
	"fmt"

	"github.com/Jeffail/gabs/v2"
)

func ConvertJSONToText(rootPath string, jsonPaths map[string]string, rootObj interface{}) ([]map[string]interface{}, error) {
	if rootObj == nil {
		return nil, fmt.Errorf("object provided cannot be nil")
	}

	var res = make([]map[string]interface{}, 0)

	objs, err := traverseJSONRoot(rootPath, rootObj)
	if err != nil {
		return nil, err
	}

	if jsonPaths == nil || len(jsonPaths) == 0 {
		return nil, fmt.Errorf("json paths must not be empty/nil")
	}

	for _, obj := range objs {
		mappedObj := make(map[string]interface{}, 0)

		for k, v := range jsonPaths {
			if !obj.ExistsP(v) {
				return nil, fmt.Errorf("wrong path provided: %s", v)
			}

			objData := obj.Path(v)
			mappedObj[k] = objData.Data()
		}

		res = append(res, mappedObj)
	}

	return res, nil
}

func traverseJSONRoot(rootPath string, obj interface{}) ([]*gabs.Container, error) {
	jsonObj, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	parsedObj, err := gabs.ParseJSON(jsonObj)
	if err != nil {
		return nil, err
	}

	if rootPath == "" {
		return []*gabs.Container{parsedObj}, nil
	}

	if !parsedObj.ExistsP(rootPath) {
		return nil, fmt.Errorf("root path does not exist in object: %s", rootPath)
	}

	parsedObj = parsedObj.Path(rootPath)
	children := parsedObj.Children()

	if children == nil {
		return nil, fmt.Errorf("root path does not lead to an array in object: %s", rootPath)
	}

	return children, nil
}
