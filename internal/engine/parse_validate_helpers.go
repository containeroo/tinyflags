package engine

import (
	"strings"

	"github.com/containeroo/tinyflags/internal/core"
)

func joinFlagNames(flags []*core.BaseFlag) string {
	names := make([]string, 0, len(flags))
	for _, fl := range flags {
		names = append(names, "--"+fl.Name)
	}
	return strings.Join(names, ", ")
}

func joinConflictNames(names []string) string {
	return strings.Join(names, " vs ")
}
