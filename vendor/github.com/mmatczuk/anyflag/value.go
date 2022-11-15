package anyflag

import (
	"fmt"
	"reflect"
)

// Value is a generic pflag.Value for a T.
type Value[T any] struct {
	parse func(val string) (T, error)
	value *T
}

// NewValue returns a new Value[T] with the given value, pointer to a T, and a parse function.
func NewValue[T any](val T, p *T, parse func(val string) (T, error)) *Value[T] {
	v := new(Value[T])
	v.parse = parse
	v.value = p
	*v.value = val
	return v
}

func (v *Value[T]) Set(val string) error {
	var err error
	*v.value, err = v.parse(val)
	return err
}

func (v *Value[T]) Type() string {
	return reflect.TypeOf(v).Name()
}

func (v *Value[T]) String() string {
	if v.value == nil || reflect.ValueOf(*v.value).IsZero() {
		return ""
	}
	return fmt.Sprint(*v.value)
}
