package core

// Registry defines the interface for registering and managing flags and dynamic groups.
type Registry interface {
	RegisterFlag(name string, bf *BaseFlag)                       // RegisterFlag registers a static flag by name.
	RegisterDynamic(group, field string, item DynamicValue) error // RegisterDynamic registers a dynamic flag under a group and field.
	GetGroup(name string) *MutualGroup                            // GetGroup retrieves a mutual exclusion group by name.
	Groups() []*MutualGroup                                       // Groups returns all defined mutual exclusion groups.
	DefaultDelimiter() string                                     // DefaultDelimiter returns the default delimiter used for slice values.
}

// Value is the interface all static flag values must implement (scalar or slice).
type Value interface {
	Set(string) error // Set parses and sets the value from the provided string.
	Get() any         // Get returns the current value as an `any`.
	Changed() bool    // Changed returns true if the flag was set explicitly.
	Default() string  // Default returns the default value as a string.
}

// DynamicValue is the interface implemented by dynamic flag containers.
// It allows setting values under a given ID dynamically (e.g., --http.alpha.port).
type DynamicValue interface {
	Set(id string, val string) error
}

// DynamicItemValues is implemented by dynamic containers to expose all parsed values.
type DynamicItemValues interface {
	ValuesAny() map[string]any // ValuesAny returns a map of parsed values keyed by ID.
}

// StrictBool is an optional interface implemented by bool flags
// that support the `--flag` / `--no-flag` style syntax.
type StrictBool interface {
	IsStrictBool() bool
}

// SliceMarker is a marker interface used to distinguish slice-type flags.
// It has no methods and is used only for internal type assertions.
type SliceMarker interface {
	isSlice()
}
