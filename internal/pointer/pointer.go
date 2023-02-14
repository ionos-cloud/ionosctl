package pointer

// From is a helper routine that allocates a new any value
// to store v and returns a pointer to it.
func From[T any](v T) *T {
	return &v
}
