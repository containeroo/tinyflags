package dynamic

// BoolValue wraps a dynamic boolean value with optional strict mode.
type BoolValue struct {
	*DynamicScalarValue[bool]       // Underlying parsed values and metadata
	strictMode                *bool // Pointer to shared strict mode flag
	hideStrict                *bool
}

// Base returns the underlying DynamicScalarValue.
func (b *BoolValue) Base() *DynamicScalarValue[bool] {
	return b.DynamicScalarValue
}

// IsStrictBool reports whether strict mode is enabled.
func (b *BoolValue) IsStrictBool() bool {
	return b.strictMode != nil && *b.strictMode
}

// IsStrictHidden reports whether the flag is hidden from usage output.
func (b *BoolValue) IsStrictHidden() bool {
	return b.hideStrict != nil && *b.hideStrict
}
