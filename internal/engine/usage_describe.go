package engine

import (
	"strings"

	"github.com/containeroo/tinyflags/internal/core"
)

// buildFlagDescription creates the full help text for a flag, including metadata
// such as allowed values, default, environment variable, deprecation, and group info.
func buildFlagDescription(flag *core.BaseFlag, globalHideEnvs bool, name string) string {
	desc := flag.Usage

	allowed := flag.AllowedValues()
	if !flag.HideAllowed && len(allowed) > 0 {
		desc += " (allowed: " + strings.Join(allowed, ", ") + ")"
	}

	if flag.Deprecated != "" {
		desc += " (deprecated: " + flag.Deprecated + ")"
	}

	if flag.ShouldShowDefaultInHelp() {
		if def := flag.Value.Default(); def != "" {
			desc += " (default: " + def + ")"
		}
	}

	if !flag.HideRequires && len(flag.Requires) > 0 {
		desc += " (requires: " + strings.Join(flag.Requires, ", ") + ")"
	}

	flag.ResolveUsageEnvKey(name, globalHideEnvs)

	if flag.ShouldShowUsageEnv(globalHideEnvs) {
		desc += " (env: " + flag.EnvKey + ")"
	}

	if !flag.HideRequired && flag.Required {
		desc += " (required)"
	}

	if flag.OneOfGroup != nil && !flag.OneOfGroup.IsHidden() {
		desc += buildGroupInfo(flag.OneOfGroup)
	}

	if flag.AllOrNone != nil && !flag.AllOrNone.IsHidden() {
		desc += buildRequireGroupInfo(flag.AllOrNone)
	}

	return desc
}

// buildGroupInfo returns group info suffix if flag belongs to a one of group.
func buildGroupInfo(group *core.OneOfGroupGroup) string {
	var b strings.Builder
	b.WriteString(" [group: ")
	if group.TitleText() != "" {
		b.WriteString(group.TitleText())
	} else {
		b.WriteString(group.Name)
	}
	b.WriteString(" (one of")
	if group.IsRequired() {
		b.WriteString(", required")
	}
	b.WriteString(")]")
	return b.String()
}

// buildRequireGroupInfo returns group info suffix if flag belongs to a require-together group.
func buildRequireGroupInfo(group *core.AllOrNoneGroup) string {
	var b strings.Builder
	b.WriteString(" [group: ")
	if group.TitleText() != "" {
		b.WriteString(group.TitleText())
	} else {
		b.WriteString(group.Name)
	}
	b.WriteString(" (all or none")
	if group.IsRequired() {
		b.WriteString(", required")
	}
	b.WriteString(")]")
	return b.String()
}
