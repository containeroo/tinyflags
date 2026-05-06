package core

import (
	"strings"

	"github.com/containeroo/tinyflags/internal/utils"
)

// SliceInputConfig centralizes delimiter and empty-item behavior for slice values.
type SliceInputConfig struct {
	Delimiter  string
	StrictDel  bool
	AllowEmpty bool
}

// Split breaks a raw slice input into chunks, validating delimiter usage first when configured.
func (c *SliceInputConfig) Split(raw string) ([]string, error) {
	if c.StrictDel {
		if err := utils.CheckMixedDelimiters(raw, c.Delimiter); err != nil {
			return nil, err
		}
	}
	return strings.Split(raw, c.Delimiter), nil
}
