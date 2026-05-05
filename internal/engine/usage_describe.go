package engine

import (
	"strings"

	"github.com/containeroo/tinyflags/internal/core"
)

// buildFlagDescription creates the full help text for a flag, including metadata
// such as allowed values, default, environment variable, deprecation, and group info.
func buildFlagDescription(flag *core.BaseFlag, globalHideEnvs bool, name string) string {
	desc := flag.Usage

	var allowed []string
	if len(flag.Allowed) > 0 {
		allowed = append(allowed, flag.Allowed...)
	} else if bv, ok := flag.Value.(core.StrictBool); ok && bv.IsStrictBool() {
		allowed = append(allowed, "true", "false")
	}
	if !flag.HideAllowed && len(allowed) > 0 {
		desc += " (Allowed: " + strings.Join(allowed, ", ") + ")"
	}

	showDefault := true
	if bv, ok := flag.Value.(core.StrictBool); ok && !bv.IsStrictBool() {
		showDefault = false
	}
	if _, ok := flag.Value.(core.Incrementable); ok {
		showDefault = false
	}

	if flag.Deprecated != "" {
		desc += " [DEPRECATED: " + flag.Deprecated + "]"
	}

	if flag.Value != nil && showDefault && !flag.HideDefault {
		if def := flag.Value.Default(); def != "" {
			desc += " (Default: " + def + ")"
		}
	}

	if !flag.HideRequires && len(flag.Requires) > 0 {
		desc += " (Requires: " + strings.Join(flag.Requires, ", ") + ")"
	}

	if shouldInjectEnvKey(flag, globalHideEnvs, name) {
		flag.EnvKey = strings.ToUpper(name + "_" + strings.ReplaceAll(flag.Name, "-", "_"))
	}

	if shouldShowEnv(flag, globalHideEnvs) {
		desc += " (Env: " + flag.EnvKey + ")"
	}

	if !flag.HideRequired && flag.Required {
		desc += " (Required)"
	}

	if flag.OneOfGroup != nil && !flag.OneOfGroup.IsHidden() {
		desc += buildGroupInfo(flag.OneOfGroup)
	}

	if flag.AllOrNone != nil && !flag.AllOrNone.IsHidden() {
		desc += buildRequireGroupInfo(flag.AllOrNone)
	}

	return desc
}

// shouldInjectEnvKey decides whether to compute EnvKey from prefix.
func shouldInjectEnvKey(flag *core.BaseFlag, globalHideEnvs bool, prefix string) bool {
	return !globalHideEnvs && !flag.DisableEnv && !flag.HideEnv && flag.EnvKey == "" && prefix != ""
}

// shouldShowEnv decides whether to include EnvKey in help.
func shouldShowEnv(flag *core.BaseFlag, globalHideEnvs bool) bool {
	return !globalHideEnvs && !flag.DisableEnv && !flag.HideEnv && flag.EnvKey != ""
}

// buildGroupInfo returns group info suffix if flag belongs to a one of group.
func buildGroupInfo(group *core.OneOfGroupGroup) string {
	var b strings.Builder
	b.WriteString(" [Group: ")
	if group.TitleText() != "" {
		b.WriteString(group.TitleText())
	} else {
		b.WriteString(group.Name)
	}
	b.WriteString(" (One Of)")
	if group.IsRequired() {
		b.WriteString(" - required")
	}
	b.WriteString("]")
	return b.String()
}

// buildRequireGroupInfo returns group info suffix if flag belongs to a require-together group.
func buildRequireGroupInfo(group *core.AllOrNoneGroup) string {
	var b strings.Builder
	b.WriteString(" [Group: ")
	if group.TitleText() != "" {
		b.WriteString(group.TitleText())
	} else {
		b.WriteString(group.Name)
	}
	b.WriteString(" (All Or None)")
	if group.IsRequired() {
		b.WriteString(" - required")
	}
	b.WriteString("]")
	return b.String()
}
