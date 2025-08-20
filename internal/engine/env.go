package engine

import "strings"

// EnvKeyFunc derives the environment variable key for a flag.
// If it returns "", the flag won't be loaded from ENV unless BaseFlag.EnvKey is set.
type EnvKeyFunc func(prefix, flagName string) string

// NewReplacerEnvKeyFunc builds an EnvKeyFunc that:
// - returns "" when prefix is empty
// - applies the given replacer to the flag name
// - joins prefix + "_" + transformed name
// - upper-cases the result (if upper is true)
func NewReplacerEnvKeyFunc(replacer *strings.Replacer, upper bool) EnvKeyFunc {
	return func(prefix, name string) string {
		if prefix == "" {
			return ""
		}
		n := replacer.Replace(name)
		if upper {
			return strings.ToUpper(prefix + "_" + n)
		}
		return prefix + "_" + n
	}
}
