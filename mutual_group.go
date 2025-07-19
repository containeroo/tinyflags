package tinyflags

type mutualGroup struct {
	name     string      // internal group name
	flags    []*baseFlag // registered flags
	title    string      // shown in help output
	hidden   bool        // hides group from help output
	required bool        // at least one flag must be set
}

func (g *mutualGroup) Title(t string) *mutualGroup {
	g.title = t
	return g
}

func (g *mutualGroup) Hidden() *mutualGroup {
	g.hidden = true
	return g
}

func (g *mutualGroup) Required() *mutualGroup {
	g.required = true
	return g
}
