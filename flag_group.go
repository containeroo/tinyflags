package tinyflags

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
