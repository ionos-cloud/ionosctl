package utils

import (
	"reflect"
)

func GetPropertiesKVSet(params ...interface{}) map[string]interface{} {
	var argKind reflect.Kind
	propertiesKVset := map[string]interface{}{}
	for _, param := range params {
		arg := reflect.TypeOf(param)
		value := reflect.ValueOf(param)
		if arg.Kind() == reflect.Ptr {
			argKind = arg.Elem().Kind()
		} else {
			argKind = arg.Kind()
		}
		if argKind == reflect.Struct {
			for i := 0; i < arg.NumField(); i++ {
				if value.Field(i).Type().Kind() == reflect.Ptr {
					if !value.Field(i).IsZero() && !value.Field(i).Elem().IsZero() {
						if value.Field(i).Elem().Type().Kind() == reflect.Struct {
							propertiesKVset[arg.Field(i).Name] = GetPropertiesKVSet(value.Field(i).Elem().Interface())
						} else {
							// Add Name of the field from struct and the value, if set
							propertiesKVset[arg.Field(i).Name] = value.Field(i).Elem().Interface()
						}
					}
				}
			}
		}
	}
	return propertiesKVset
}
