package core

import "strings"

var envKeyPartReplacer = strings.NewReplacer("-", "_", ".", "_", "/", "_")

// NormalizeEnvKeyPart converts one flag-path segment to its canonical ENV form.
func NormalizeEnvKeyPart(s string) string {
	return strings.ToUpper(envKeyPartReplacer.Replace(s))
}

// DynamicEnvKey builds the canonical environment key for a dynamic flag.
func DynamicEnvKey(prefix, group, id, field string) string {
	if prefix == "" {
		return ""
	}
	parts := []string{
		NormalizeEnvKeyPart(prefix),
		NormalizeEnvKeyPart(group),
		id,
		NormalizeEnvKeyPart(field),
	}
	return strings.Join(parts, "_")
}

// DynamicEnvID converts the ENV instance segment back to a dynamic flag ID.
func DynamicEnvID(segment string) string {
	return strings.ToLower(segment)
}
