package tinyflags

type Flag[T any] struct {
	builderBase[T]
}
