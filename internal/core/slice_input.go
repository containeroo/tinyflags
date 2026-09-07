package core

import "strings"

// SliceInputConfig centralizes delimiter and item normalization for slice values.
type SliceInputConfig struct {
	Delimiter  string
	NoSplit    bool
	AllowEmpty bool
	TrimSpace  bool
}

// Split breaks a raw slice input into chunks using the configured delimiter.
func (c *SliceInputConfig) Split(raw string) ([]string, error) {
	if c.NoSplit {
		return []string{raw}, nil
	}
	return strings.Split(raw, c.Delimiter), nil
}

// Normalize prepares one split item before parsing.
func (c *SliceInputConfig) Normalize(raw string) string {
	if c.TrimSpace {
		return strings.TrimSpace(raw)
	}
	return raw
}
