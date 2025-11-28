package builder

import "github.com/containeroo/tinyflags/internal/core"

// flagMeta centralizes common flag metadata mutations used by both static and dynamic builders.
type flagMeta struct {
	registry core.Registry
	bf       *core.BaseFlag
}

func (m *flagMeta) required()                { m.bf.Required = true }
func (m *flagMeta) hideRequired()            { m.bf.HideRequired = true }
func (m *flagMeta) hidden()                  { m.bf.Hidden = true }
func (m *flagMeta) deprecated(reason string) { m.bf.Deprecated = reason }

func (m *flagMeta) oneOfGroup(name string) {
	if name == "" {
		return
	}
	g := m.registry.GetOneOfGroup(name)
	g.Flags = append(g.Flags, m.bf)
	m.bf.OneOfGroup = g
}

func (m *flagMeta) allOrNoneGroup(name string) {
	if name == "" {
		return
	}
	g := m.registry.GetAllOrNoneGroup(name)
	g.Flags = append(g.Flags, m.bf)
	m.bf.AllOrNone = g
}

func (m *flagMeta) env(key string) {
	if m.bf.DisableEnv {
		panic("cannot call Env after DisableEnv")
	}
	m.bf.EnvKey = key
}

func (m *flagMeta) disableEnv() {
	if m.bf.EnvKey != "" {
		panic("cannot call DisableEnv after Env")
	}
	m.bf.DisableEnv = true
}

func (m *flagMeta) placeholder(s string) { m.bf.Placeholder = s }

func (m *flagMeta) allowed(vals ...string) {
	m.bf.Allowed = append([]string(nil), vals...)
}

func (m *flagMeta) hideAllowed() { m.bf.HideAllowed = true }

func (m *flagMeta) requires(names ...string) {
	m.bf.Requires = append([]string(nil), names...)
}

func (m *flagMeta) hideRequires() { m.bf.HideRequires = true }

func (m *flagMeta) hideDefault() { m.bf.HideDefault = true }

func (m *flagMeta) section(name string) { m.bf.Section = name }
