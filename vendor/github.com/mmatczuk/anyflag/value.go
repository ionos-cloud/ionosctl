package anyflag

import (
	"fmt"
	"reflect"

	"github.com/spf13/pflag"
)

// Value is a generic pflag.Value for a T.
type Value[T any] struct {
	parse  func(val string) (T, error)
	value  *T
	redact func(T) string
}

// NewValue returns a new Value[T] with the given value, pointer to a T, and a parse function.
func NewValue[T any](val T, p *T, parse func(val string) (T, error)) *Value[T] {
	v := new(Value[T])
	v.parse = parse
	v.value = p
	*v.value = val
	return v
}

// NewValueWithRedact returns a new Value[T] and additionally sets custom String() function.
// Redact primary purpose is to redact passwords to prevent them from leaking in logs.
func NewValueWithRedact[T any](val T, p *T, parse func(val string) (T, error), redact func(T) string) *Value[T] {
	v := NewValue(val, p, parse)
	v.redact = redact
	return v
}

// Unredacted returns a copy of Value[T] without redact function.
func (v *Value[T]) Unredacted() pflag.Value {
	if v.redact == nil {
		return v
	}

	vv := *v
	vv.redact = nil
	return &vv
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
	if v.redact != nil {
		return v.redact(*v.value)
	}

	if v.value == nil || reflect.ValueOf(*v.value).IsZero() {
		return ""
	}
	return fmt.Sprint(*v.value)
}
