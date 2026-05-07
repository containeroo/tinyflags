package core

import "strings"

// SliceInputConfig centralizes delimiter and empty-item behavior for slice values.
type SliceInputConfig struct {
	Delimiter  string
	AllowEmpty bool
}

// Split breaks a raw slice input into chunks using the configured delimiter.
func (c *SliceInputConfig) Split(raw string) ([]string, error) {
	return strings.Split(raw, c.Delimiter), nil
}
