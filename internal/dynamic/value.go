package dynamic

type dynamicSliceValue[T any] struct {
	def     string
	changed bool
}

func (v *dynamicSliceValue[T]) Set(string) error { return nil } // not used
func (v *dynamicSliceValue[T]) Get() any         { return nil } // not used
func (v *dynamicSliceValue[T]) Changed() bool    { return v.changed }
func (v *dynamicSliceValue[T]) Default() string  { return v.def }

func (v *dynamicSliceValue[T]) IsSlice() {}

type dynamicBoolValue[T any] struct {
	def        string
	changed    bool
	strictMode *bool
}

func (v *dynamicBoolValue[T]) Changed() bool      { return v.changed }
func (v *dynamicBoolValue[T]) Default() string    { return v.def }
func (v *dynamicBoolValue[T]) Set(string) error   { return nil }
func (v *dynamicBoolValue[T]) Get() any           { return nil }
func (v *dynamicBoolValue[T]) IsStrictBool() bool { return v.strictMode != nil && *v.strictMode }
