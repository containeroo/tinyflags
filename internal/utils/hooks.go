package utils

// ApplyValueHooks runs validation and finalization hooks on a parsed value.
func ApplyValueHooks[T any](val T, validate func(T) error, finalize func(T) T) (T, error) {
	if validate != nil {
		if err := validate(val); err != nil {
			return val, err
		}
	}
	if finalize != nil {
		val = finalize(val)
	}
	return val, nil
}
