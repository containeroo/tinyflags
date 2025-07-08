package tinyflags

type mockValue struct {
	val         string
	defaultVal  string
	changed     bool
	errToReturn error

	setCalled bool
	getCalled bool
}

func (v *mockValue) Set(s string) error {
	v.setCalled = true
	v.changed = true
	v.val = s
	return v.errToReturn
}

func (v *mockValue) Get() any {
	v.getCalled = true
	return v.val
}

func (v *mockValue) Default() string {
	return v.defaultVal
}

func (v *mockValue) IsChanged() bool {
	return v.changed
}
