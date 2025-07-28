package dynamic

type BoolValue struct {
	*DynamicScalarValue[bool]
	Strict *bool
}

func (b *BoolValue) Base() *DynamicScalarValue[bool] { return b.DynamicScalarValue }
func (b *BoolValue) IsStrictBool() bool              { return b.Strict != nil && *b.Strict }
