package functional

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

// Fold accumulates the result of f into acc and returns acc by applying f over each element in the slice
func Fold[T any, Acc any](xs []T, f func(Acc, T) Acc, acc Acc) Acc {
	for _, x := range xs {
		acc = f(acc, x)
	}
	return acc
}

// Map applies a function to each element of a slice and returns the modified slice without considering the index of each element.
func Map[T comparable, K any](s []T, f func(T) K) []K {
	return MapIdx(s, func(_ int, t T) K {
		return f(t)
	})
}

// MapIdx applies a function to each element and index of a slice, returning the modified slice with consideration of the index.
func MapIdx[V comparable, R any](s []V, f func(int, V) R) []R {
	sm := make([]R, len(s))
	for i, v := range s {
		sm[i] = f(i, v)
	}
	return sm
}
