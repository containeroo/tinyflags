package dynamic

type dynamicHelpValue[T any] struct {
	def     string
	changed bool
}

func (v *dynamicHelpValue[T]) Set(string) error { return nil } // not used
func (v *dynamicHelpValue[T]) Get() any         { return nil } // not used
func (v *dynamicHelpValue[T]) Changed() bool    { return v.changed }
func (v *dynamicHelpValue[T]) Default() string  { return v.def }
