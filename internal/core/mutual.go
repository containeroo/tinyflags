package core

// MutualExlusiveGroup enforces exclusivity among a set of flags.
// Only one flag in the group may be set.
type MutualExlusiveGroup struct {
	Name           string                   // Identifier for this group.
	Flags          []*BaseFlag              // Member flags.
	RequiredGroups []*RequiredTogetherGroup // Optional grouped sets
	titleText      string                   // Optional title to display in help.
	hidden         bool                     // Hide this group in help.
	required       bool                     // Require exactly one of the flags.
}

// Title sets a custom help heading.
func (g *MutualExlusiveGroup) Title(t string) *MutualExlusiveGroup {
	g.titleText = t
	return g
}

// TitleText returns the custom heading, if any.
func (g *MutualExlusiveGroup) TitleText() string {
	return g.titleText
}

// Hidden marks the group as omitted from help.
func (g *MutualExlusiveGroup) Hidden() *MutualExlusiveGroup {
	g.hidden = true
	return g
}

// IsHidden reports whether the group is hidden.
func (g *MutualExlusiveGroup) IsHidden() bool {
	return g.hidden
}

// Required enforces that one member must be set.
func (g *MutualExlusiveGroup) Required() *MutualExlusiveGroup {
	g.required = true
	return g
}

// IsRequired reports whether the group is required.
func (g *MutualExlusiveGroup) IsRequired() bool {
	return g.required
}

// AddGroup includes a require-together group as one exclusive member.
func (g *MutualExlusiveGroup) AddGroup(grp *RequiredTogetherGroup) *MutualExlusiveGroup {
	if grp != nil {
		g.RequiredGroups = append(g.RequiredGroups, grp)
		for _, fl := range grp.Flags {
			fl.MutualGroup = g // populate the parent mutual group reference
		}
	}
	return g
}
