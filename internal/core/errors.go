package core

import "fmt"

// Optional knob: value types that can enforce once-per-id expose this.
type OncePerIDToggler interface {
	EnableOncePerID()
}

// Typed error so the parser can pretty-print a friendly message.
type DuplicatePerIDError struct {
	Field string // dynamic field name, e.g. "timeout"
	ID    string // instance id, e.g. "alpha"
}

func (e *DuplicatePerIDError) Error() string {
	return fmt.Sprintf("duplicate value for field %q in ID %q", e.Field, e.ID)
}
