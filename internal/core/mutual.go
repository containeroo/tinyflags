package core

// MutualGroup defines a mutually exclusive group of flags.
// Only one flag in a group can be set at a time.
type MutualGroup struct {
	Name      string      // Name is the group identifier.
	Flags     []*BaseFlag // Flags is the list of flags that belong to this group.
	titleText string      // titleText is the custom title shown in help output (optional).
	hidden    bool        // isHidden marks the group as hidden from help output.
	required  bool        // isRequired marks the group as required (one flag must be set).
}

// Title sets the custom title text shown in the help output.
func (g *MutualGroup) Title(t string) *MutualGroup {
	g.titleText = t
	return g
}

// TitleText returns the group's custom title text.
func (g *MutualGroup) TitleText() string {
	return g.titleText
}

// Hidden marks the group as hidden from help output.
func (g *MutualGroup) Hidden() *MutualGroup {
	g.hidden = true
	return g
}

// IsHidden returns true if the group is marked as hidden.
func (g *MutualGroup) IsHidden() bool {
	return g.hidden
}

// Required marks the group as required (at least one flag must be set).
func (g *MutualGroup) Required() *MutualGroup {
	g.required = true
	return g
}

// IsRequired returns true if the group is marked as required.
func (g *MutualGroup) IsRequired() bool {
	return g.required
}
