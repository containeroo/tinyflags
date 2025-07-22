package core

// MutualGroup enforces exclusivity among a set of flags.
// Only one flag in the group may be set.
type MutualGroup struct {
	Name      string      // Identifier for this group.
	Flags     []*BaseFlag // Member flags.
	titleText string      // Optional title to display in help.
	hidden    bool        // Hide this group in help.
	required  bool        // Require exactly one of the flags.
}

// Title sets a custom help heading.
func (g *MutualGroup) Title(t string) *MutualGroup {
	g.titleText = t
	return g
}

// TitleText returns the custom heading, if any.
func (g *MutualGroup) TitleText() string {
	return g.titleText
}

// Hidden marks the group as omitted from help.
func (g *MutualGroup) Hidden() *MutualGroup {
	g.hidden = true
	return g
}

// IsHidden reports whether the group is hidden.
func (g *MutualGroup) IsHidden() bool {
	return g.hidden
}

// Required enforces that one member must be set.
func (g *MutualGroup) Required() *MutualGroup {
	g.required = true
	return g
}

// IsRequired reports whether the group is required.
func (g *MutualGroup) IsRequired() bool {
	return g.required
}
