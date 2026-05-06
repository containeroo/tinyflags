package engine

import (
	"strings"

	"github.com/containeroo/tinyflags/internal/core"
)

// joinFlagNames formats flag names for validation errors.
func joinFlagNames(flags []*core.BaseFlag) string {
	names := make([]string, 0, len(flags))
	for _, fl := range flags {
		names = append(names, "--"+fl.Name)
	}
	return strings.Join(names, ", ")
}

// joinConflictNames formats conflicting names for one-of errors.
func joinConflictNames(names []string) string {
	return strings.Join(names, " vs ")
}
