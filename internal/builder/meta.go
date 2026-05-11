package builder

import "github.com/containeroo/tinyflags/internal/core"

// flagMeta centralizes common flag metadata mutations used by both static and dynamic builders.
type flagMeta struct {
	registry core.Registry
	bf       *core.BaseFlag
}

// required marks the flag as required.
func (m *flagMeta) required() { m.bf.Required = true }

// hideRequired hides the required marker in help output.
func (m *flagMeta) hideRequired() { m.bf.HideRequired = true }

// hidden hides the flag from help output.
func (m *flagMeta) hidden() { m.bf.Hidden = true }

// deprecated marks the flag as deprecated.
func (m *flagMeta) deprecated(reason string) { m.bf.Deprecated = reason }

// oneOfGroup attaches the flag to a one-of group.
func (m *flagMeta) oneOfGroup(name string) {
	if name == "" {
		return
	}
	g := m.registry.GetOneOfGroup(name)
	g.Flags = appendBaseFlagUnique(g.Flags, m.bf)
	m.bf.OneOfGroups = appendOneOfGroupUnique(m.bf.OneOfGroups, g)
}

// helpOneOfGroups overrides which one-of groups should appear in help.
func (m *flagMeta) helpOneOfGroups(names ...string) {
	m.bf.HelpOneOfSet = true
	m.bf.HelpOneOf = m.bf.HelpOneOf[:0]
	for _, name := range names {
		if name == "" {
			continue
		}
		g := m.registry.GetOneOfGroup(name)
		m.bf.HelpOneOf = appendOneOfGroupUnique(m.bf.HelpOneOf, g)
	}
}

// allOrNoneGroup attaches the flag to an all-or-none group.
func (m *flagMeta) allOrNoneGroup(name string) {
	if name == "" {
		return
	}
	g := m.registry.GetAllOrNoneGroup(name)
	g.Flags = append(g.Flags, m.bf)
	m.bf.AllOrNone = g
}

// env assigns the environment variable key for the flag.
func (m *flagMeta) env(key string) {
	if m.bf.DisableEnv {
		panic("cannot call Env after DisableEnv")
	}
	m.bf.EnvKey = key
}

// disableEnv disables environment variable parsing for the flag.
func (m *flagMeta) disableEnv() {
	if m.bf.EnvKey != "" {
		panic("cannot call DisableEnv after Env")
	}
	m.bf.DisableEnv = true
}

// hideEnv hides the environment variable hint in help output.
func (m *flagMeta) hideEnv() { m.bf.HideEnv = true }

// placeholder sets the usage placeholder for the flag.
func (m *flagMeta) placeholder(s string) { m.bf.Placeholder = s }

// allowed records the allowed values for the flag.
func (m *flagMeta) allowed(vals ...string) {
	m.bf.Allowed = append([]string(nil), vals...)
}

// hideAllowed hides the allowed-values suffix in help output.
func (m *flagMeta) hideAllowed() { m.bf.HideAllowed = true }

// requires records other flags required alongside this flag.
func (m *flagMeta) requires(names ...string) {
	m.bf.Requires = append([]string(nil), names...)
}

// hideRequires hides the requires suffix in help output.
func (m *flagMeta) hideRequires() { m.bf.HideRequires = true }

// hideDefault hides the default value suffix in help output.
func (m *flagMeta) hideDefault() { m.bf.HideDefault = true }

// section assigns the help section for the flag.
func (m *flagMeta) section(name string) { m.bf.Section = name }

// maskFn sets the masking function for overridden values.
func (m *flagMeta) maskFn(fn func(any) any) { m.bf.MaskFn = fn }

func appendBaseFlagUnique(flags []*core.BaseFlag, target *core.BaseFlag) []*core.BaseFlag {
	for _, flag := range flags {
		if flag == target {
			return flags
		}
	}
	return append(flags, target)
}

func appendOneOfGroupUnique(groups []*core.OneOfGroupGroup, target *core.OneOfGroupGroup) []*core.OneOfGroupGroup {
	for _, group := range groups {
		if group == target {
			return groups
		}
	}
	return append(groups, target)
}
