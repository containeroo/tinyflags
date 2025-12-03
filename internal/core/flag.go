package core

// Flag is the minimal interface exposed by built-in flag handles.
// It reports whether a value was provided and returns the resulting pointer.
type Flag[T any] interface {
	Changed() bool
	Value() *T
}
