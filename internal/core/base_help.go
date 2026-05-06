package core

import "strings"

// AllowedValues returns the explicit or implied allowed values for help output.
func (f *BaseFlag) AllowedValues() []string {
	if len(f.Allowed) > 0 {
		out := make([]string, len(f.Allowed))
		copy(out, f.Allowed)
		return out
	}
	if bv, ok := f.Value.(StrictBool); ok && bv.IsStrictBool() {
		return []string{"true", "false"}
	}
	return nil
}

// ShouldShowDefaultInHelp reports whether the default belongs in help output.
func (f *BaseFlag) ShouldShowDefaultInHelp() bool {
	if f == nil || f.Value == nil || f.HideDefault {
		return false
	}
	if bv, ok := f.Value.(StrictBool); ok && !bv.IsStrictBool() {
		return false
	}
	if _, ok := f.Value.(Incrementable); ok {
		return false
	}
	return true
}

// ResolveUsageEnvKey sets a derived env key for help output when one should be shown.
func (f *BaseFlag) ResolveUsageEnvKey(prefix string, globalHideEnvs bool) {
	if f == nil || globalHideEnvs || f.DisableEnv || f.HideEnv || f.EnvKey != "" || prefix == "" {
		return
	}
	f.EnvKey = strings.ToUpper(prefix + "_" + strings.ReplaceAll(f.Name, "-", "_"))
}

// ShouldShowUsageEnv reports whether the env key belongs in help output.
func (f *BaseFlag) ShouldShowUsageEnv(globalHideEnvs bool) bool {
	return f != nil && !globalHideEnvs && !f.DisableEnv && !f.HideEnv && f.EnvKey != ""
}

// UsagePlaceholder returns the placeholder that should appear in help output.
func (f *BaseFlag) UsagePlaceholder() string {
	isBool := false
	isStrict := false

	if bv, ok := f.Value.(StrictBool); ok {
		isBool = true
		isStrict = bv.IsStrictBool()
	}

	if isBool && !isStrict {
		return ""
	}
	if f.Placeholder != "" {
		return f.Placeholder
	}
	if allowed := f.AllowedValues(); len(allowed) > 0 {
		return "<" + strings.Join(allowed, "|") + ">"
	}
	if isStrict {
		return "<true|false>"
	}
	if _, ok := f.Value.(Incrementable); ok {
		return ""
	}
	placeholder := strings.ToUpper(f.Name)
	if _, ok := f.Value.(SliceMarker); ok {
		placeholder += "..."
	}
	return placeholder
}
