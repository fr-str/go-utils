package conv

// Ptr returns pointer to value
func Ptr[T any](v T) *T {
	return &v
}
