/*

The allocate library provides helper functions for allocation/initializing structs.

@author Colorado Reed (colorado at dotdashpay ... com)

This library was seeded by this discussion in the golang-nuts mailing list:
https://groups.google.com/forum/#!topic/golang-nuts/Wd9jiZswwMU
*/

package allocate

import (
	"fmt"
	"reflect"
)

// Zero allocates an input structure such that all pointer fields
// are fully allocated, i.e. rather than having a nil value,
// the pointer contains a pointer to an initialized value,
// e.g. an *int field will be a pointer to 0 instead of a nil pointer.
//
// Zero does not allocate private fields.
func Zero(inputIntf interface{}) error {
	indirectVal := reflect.Indirect(reflect.ValueOf(inputIntf))

	if !indirectVal.CanSet() {
		return fmt.Errorf("Input interface is not addressable (can't Set the memory address): %#v",
			inputIntf)
	}
	if indirectVal.Kind() != reflect.Struct {
		return fmt.Errorf("allocate.Zero currently only works with [pointers to] structs, not type %v",
			indirectVal.Kind())
	}

	// allocate each of the structs fields
	var err error
	for i := 0; i < indirectVal.NumField(); i++ {
		field := indirectVal.Field(i)

		// pre-allocate pointer fields
		if field.Kind() == reflect.Ptr && field.IsNil() {
			if field.CanSet() {
				field.Set(reflect.New(field.Type().Elem()))
			}
		}

		indirectField := reflect.Indirect(field)
		switch indirectField.Kind() {
		case reflect.Map:
			indirectField.Set(reflect.MakeMap(indirectField.Type()))
		case reflect.Struct:
			// recursively allocate each of the structs embedded fields
			if field.Kind() == reflect.Ptr {
				err = Zero(field.Interface())
			} else {
				// field of Struct can always use field.Addr()
				fieldAddr := field.Addr()
				if fieldAddr.CanInterface() {
					err = Zero(fieldAddr.Interface())
				} else {
					err = fmt.Errorf("struct field can't interface, %#v", fieldAddr)
				}
			}
		}
		if err != nil {
			return err
		}
	}
	return err
}

// MustZero will panic instead of return error
func MustZero(inputIntf interface{}) {
	err := Zero(inputIntf)
	if err != nil {
		panic(err)
	}
}

// TODO(cjrd)
// Add an allocate.Random() function that assigns random values rather than nil values
