package functional

import (
	"errors"
)

/*
 * TODO: Remove this pkg once sdk-go-bundle is imported, and import github.com/ionos-cloud/sdk-go-bundle/shared/functional.go
 */

// ApplyAndAggregateErrors applies the provided function for each element of the slice
// If the function returns an error, it accumulates the error and continues execution
// After all elements are processed, it returns the aggregated errors if any
func ApplyAndAggregateErrors[T any](xs []T, f func(T) error) error {
	return Fold(
		xs,
		func(errs error, x T) error {
			err := f(x)
			if err != nil {
				errs = errors.Join(errs, err)
			}
			return errs
		},
		nil,
	)
}

// ApplyOrFail tries applying the provided function for each element of the slice
// If the function returns an error, we break execution and return the error
func ApplyOrFail[T any](xs []T, f func(T) error) error {
	return Fold(
		xs,
		func(err error, x T) error {
			// accumulate the error. If it's not nil break out of the fold
			if err != nil {
				return err
			}
			return f(x)
		},
		nil,
	)
}

// Filter applies a function to each element of a slice, returning a new slice with only the elements for which the function returns true.
func Filter[T any](xs []T, f func(T) bool) []T {
	result := make([]T, 0, len(xs))
	for _, x := range xs {
		if f(x) {
			result = append(result, x)
		}
	}
	return result
}

// Fold accumulates the result of f into acc and returns acc by applying f over each element in the slice
func Fold[T any, Acc any](xs []T, f func(Acc, T) Acc, acc Acc) Acc {
	for _, x := range xs {
		acc = f(acc, x)
	}
	return acc
}

// Map applies a function to each element of a slice and returns the modified slice without considering the index of each element.
func Map[T any, K any](s []T, f func(T) K) []K {
	return MapIdx(s, func(_ int, t T) K {
		return f(t)
	})
}

// MapIdx applies a function to each element and index of a slice, returning the modified slice with consideration of the index.
func MapIdx[V any, R any](s []V, f func(int, V) R) []R {
	sm := make([]R, len(s))
	for i, v := range s {
		sm[i] = f(i, v)
	}
	return sm
}

func KeysOfMap[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
