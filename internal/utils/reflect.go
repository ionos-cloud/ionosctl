package utils

import (
	"reflect"
)

// GetPropertiesKVSet converts a struct to a map[string]interface{}.
// It will only include fields that are set and are not nil.
// It will also recursively convert any nested structs.
//
//		type MyStruct struct {
//		    Field1 int
//		    Field2 string
//		}
//
//	  instance := MyStruct{Field1: 42, Field2: "Hello"}
//	  result := GetPropertiesKVSet(instance) // map[string]interface{}{"Field1": 42, "Field2": "Hello"}
func GetPropertiesKVSet(i interface{}) map[string]interface{} {
	properties := make(map[string]interface{})
	if i == nil {
		return properties
	}

	v := reflect.Indirect(reflect.ValueOf(i))
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)

		if fieldValue.Kind() == reflect.Ptr && !fieldValue.IsNil() {
			fieldValue = fieldValue.Elem()
		}

		if fieldValue.IsZero() {
			continue
		}

		// If the field is a struct, recursively call the function
		if fieldValue.Kind() == reflect.Struct {
			properties[field.Name] = GetPropertiesKVSet(fieldValue.Interface())
		} else {
			properties[field.Name] = fieldValue.Interface()
		}
	}
	return properties
}
