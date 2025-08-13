package goserversdk

// ToPtr converts a value to a pointer.
// Avoid passing pointers to this function.
func ToPtr[T any](v T) *T {
	return &v
}
