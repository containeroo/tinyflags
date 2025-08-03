package core

// AllOrNoneGroup enforces that all member flags must be set if any is set.
type AllOrNoneGroup struct {
	Name      string      // Identifier for this group.
	Flags     []*BaseFlag // Member flags.
	titleText string      // Optional title to display in help.
	hidden    bool        // Hide this group in help.
	required  bool        // If true, at least one must be set.
}

// Title sets a custom help heading.
func (g *AllOrNoneGroup) Title(t string) *AllOrNoneGroup {
	g.titleText = t
	return g
}

// TitleText returns the custom heading, if any.
func (g *AllOrNoneGroup) TitleText() string {
	return g.titleText
}

// Hidden marks the group as omitted from help.
func (g *AllOrNoneGroup) Hidden() *AllOrNoneGroup {
	g.hidden = true
	return g
}

// IsHidden reports whether the group is hidden.
func (g *AllOrNoneGroup) IsHidden() bool {
	return g.hidden
}

// Required enforces that at least one of the group flags must be set.
func (g *AllOrNoneGroup) Required() *AllOrNoneGroup {
	g.required = true
	return g
}

// IsRequired reports whether the group is required.
func (g *AllOrNoneGroup) IsRequired() bool {
	return g.required
}
