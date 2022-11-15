package anyflag

import (
	"fmt"
	"reflect"
	"strings"
)

// SliceValue is a generic pflag.SliceValue for a slice of T.
type SliceValue[T any] struct {
	parse   func(val string) (T, error)
	value   *[]T
	changed bool
}

// NewSliceValue returns a new SliceValue[T] with the given value, pointer to a slice of T, and a parse function.
func NewSliceValue[T any](val []T, p *[]T, parse func(val string) (T, error)) *SliceValue[T] {
	sv := new(SliceValue[T])
	sv.parse = parse
	sv.value = p
	*sv.value = val
	return sv
}

func (s *SliceValue[T]) Set(val string) error {
	ss := strings.Split(val, ",")
	out := make([]T, len(ss))
	for i, d := range ss {
		var err error
		out[i], err = s.parse(d)
		if err != nil {
			return err
		}

	}
	if !s.changed {
		*s.value = out
	} else {
		*s.value = append(*s.value, out...)
	}
	s.changed = true
	return nil
}

func (s *SliceValue[T]) Type() string {
	return reflect.TypeOf(s).Name()
}

func (s *SliceValue[T]) String() string {
	out := make([]string, len(*s.value))
	for i, d := range *s.value {
		out[i] = fmt.Sprint(d)
	}
	return "[" + strings.Join(out, ",") + "]"
}

func (s *SliceValue[T]) fromString(val string) (T, error) {
	return s.parse(val)
}

func (s *SliceValue[T]) toString(val T) string {
	return fmt.Sprint(val)
}

func (s *SliceValue[T]) Append(val string) error {
	i, err := s.fromString(val)
	if err != nil {
		return err
	}
	*s.value = append(*s.value, i)
	return nil
}

func (s *SliceValue[T]) Replace(val []string) error {
	out := make([]T, len(val))
	for i, d := range val {
		var err error
		out[i], err = s.fromString(d)
		if err != nil {
			return err
		}
	}
	*s.value = out
	return nil
}

func (s *SliceValue[T]) GetSlice() []string {
	out := make([]string, len(*s.value))
	for i, d := range *s.value {
		out[i] = s.toString(d)
	}
	return out
}
