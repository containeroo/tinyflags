package core

// OneOfGroupGroup enforces exclusivity among a set of flags.
// Only one flag in the group may be set.
type OneOfGroupGroup struct {
	Name           string            // Identifier for this group.
	Flags          []*BaseFlag       // Member flags.
	RequiredGroups []*AllOrNoneGroup // Optional grouped sets
	titleText      string            // Optional title to display in help.
	hidden         bool              // Hide this group in help.
	required       bool              // Require exactly one of the flags.
}

// Title sets a custom help heading.
func (g *OneOfGroupGroup) Title(t string) *OneOfGroupGroup {
	g.titleText = t
	return g
}

// TitleText returns the custom heading, if any.
func (g *OneOfGroupGroup) TitleText() string {
	return g.titleText
}

// Hidden marks the group as omitted from help.
func (g *OneOfGroupGroup) Hidden() *OneOfGroupGroup {
	g.hidden = true
	return g
}

// IsHidden reports whether the group is hidden.
func (g *OneOfGroupGroup) IsHidden() bool {
	return g.hidden
}

// Required enforces that one member must be set.
func (g *OneOfGroupGroup) Required() *OneOfGroupGroup {
	g.required = true
	return g
}

// IsRequired reports whether the group is required.
func (g *OneOfGroupGroup) IsRequired() bool {
	return g.required
}

// AddGroup includes a OneOfGroup group as one exclusive member.
func (g *OneOfGroupGroup) AddGroup(grp *AllOrNoneGroup) *OneOfGroupGroup {
	if grp != nil {
		g.RequiredGroups = append(g.RequiredGroups, grp)
	}
	return g
}
