package tinyflags

// Flag is a generic type for defining scalar.
type Flag[T any] struct {
	fs  *FlagSet  // The parent FlagSet this flag belongs to
	bf  *baseFlag // Internal flag metadata and state
	ptr *T        // Pointer to the destination variable for storing the parsed value
}

// Env sets the environment variable name that can override this flag.
// If not set explicitly, a default name is generated based on the flag name.
func (b *Flag[T]) Env(key string) *Flag[T] {
	if b.bf.disableEnv {
		panic("tinyflags: cannot call Env(...) after DisableEnv() on the same flag")
	}
	b.bf.envKey = key
	return b
}

// DisableEnv disables reading from environment variables for this flag,
// even if a global prefix is configured via the FlagSet.
func (b *Flag[T]) DisableEnv() *Flag[T] {
	if b.bf.envKey != "" {
		panic("tinyflags: cannot call DisableEnv() after Env(...) on the same flag")
	}
	b.bf.disableEnv = true
	return b
}

// Group assigns this flag to a mutual exclusion group.
// Flags in the same group are mutually exclusive and cannot be used together.
func (b *Flag[T]) Group(name string) *Flag[T] {
	b.fs.Group(name, b.bf)
	return b
}

// Deprecated marks the flag as deprecated, and the provided reason is shown in help output.
func (b *Flag[T]) Deprecated(reason string) *Flag[T] {
	b.bf.deprecated = reason
	return b
}

// Required marks this flag as required. Parsing will fail if it is not provided.
func (b *Flag[T]) Required() *Flag[T] {
	b.bf.required = true
	return b
}

// Value returns a pointer to the underlying parsed value.
func (b *Flag[T]) Value() *T {
	return b.ptr
}

// Metavar sets the placeholder name shown in help output for the flag value.
// If not set, it defaults to the uppercase flag name.
func (b *Flag[T]) Metavar(s string) *Flag[T] {
	b.bf.metavar = s
	return b
}

// Choices restricts the allowed values for this flag to a predefined set.
func (b *Flag[T]) Choices(allowed ...T) *Flag[T] {
	if bv, ok := b.bf.value.(*FlagItem[T]); ok {
		// Build validator from list
		bv.SetValidator(func(v T) bool {
			for _, a := range allowed {
				if bv.format(a) == bv.format(v) {
					return true
				}
			}
			return false
		}, allowed)

		// Convert allowed values to string for help text
		b.bf.allowed = make([]string, len(allowed))
		for i, x := range allowed {
			b.bf.allowed[i] = bv.format(x)
		}
	}
	return b
}

// Validator adds a custom validation function for the flag value.
func (b *Flag[T]) Validator(fn func(T) bool) *Flag[T] {
	if bv, ok := b.bf.value.(*FlagItem[T]); ok {
		bv.SetValidator(fn, nil) // no predefined list
	}
	return b
}

// Hidden marks the flag as hidden, and it will not be shown in help output.
func (b *Flag[T]) Hidden() *Flag[T] {
	b.bf.hidden = true
	return b
}

// SliceFlag is a specialized builder for slice-type flags.
// It embeds FlagBuilder and provides an additional Delimiter method
// to control how multiple values are split from a single argument.
type SliceFlag[T any] struct {
	Flag[T] // Inherits all scalar flag methods
}

// Delimiter sets the separator for slice values (e.g., "," or ":").
func (b *SliceFlag[T]) Delimiter(sep string) *SliceFlag[T] {
	if d, ok := b.bf.value.(HasDelimiter); ok {
		d.SetDelimiter(sep)
	}
	return b
}
