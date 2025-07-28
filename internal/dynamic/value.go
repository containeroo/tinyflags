package dynamic

type placeholderValue struct {
	def string
}

func (p *placeholderValue) Set(string) error { return nil }
func (p *placeholderValue) Get() any         { return nil }
func (p *placeholderValue) Changed() bool    { return false }
func (p *placeholderValue) Default() string  { return p.def }

type slicePlaceholder struct {
	def string
}

func (v *slicePlaceholder) Set(string) error { return nil }
func (v *slicePlaceholder) Get() any         { return nil }
func (v *slicePlaceholder) Changed() bool    { return false }
func (v *slicePlaceholder) Default() string  { return v.def }
func (v *slicePlaceholder) IsSlice()         {}

type boolPlaceholder struct {
	def    string
	strict *bool
}

func (v *boolPlaceholder) Set(string) error   { return nil }
func (v *boolPlaceholder) Get() any           { return nil }
func (v *boolPlaceholder) Changed() bool      { return false }
func (v *boolPlaceholder) Default() string    { return v.def }
func (v *boolPlaceholder) IsStrictBool() bool { return v.strict != nil && *v.strict }
