package tinyflags

// SetDescIndent sets the indentation for the description block.
func (f *FlagSet) SetDescIndent(n int) { f.Layout().SetDescIndent(n) }

// DescIndent returns the current indent used for the description block.
func (f *FlagSet) DescIndent() int { return f.impl.DescIndent() }

// SetDescWidth sets the wrapping width for the description block.
func (f *FlagSet) SetDescWidth(max int) { f.Layout().SetDescWidth(max) }

// DescWidth returns the wrapping width for the description block.
func (f *FlagSet) DescWidth() int { return f.impl.DescWidth() }

// SetStaticUsageIndent sets the indentation for static flag usage lines.
func (f *FlagSet) SetStaticUsageIndent(n int) { f.Layout().SetStaticUsageIndent(n) }

// StaticUsageIndent returns the static usage indentation.
func (f *FlagSet) StaticUsageIndent() int { return f.impl.StaticUsageIndent() }

// SetStaticUsageColumn sets the column at which static flag descriptions begin.
func (f *FlagSet) SetStaticUsageColumn(col int) { f.Layout().SetStaticUsageColumn(col) }

// StaticUsageColumn returns the description column for static flags.
func (f *FlagSet) StaticUsageColumn() int { return f.impl.StaticUsageColumn() }

// SetStaticUsageWidth sets the max wrapping width for static flag descriptions.
func (f *FlagSet) SetStaticUsageWidth(maxWidth int) { f.Layout().SetStaticUsageWidth(maxWidth) }

// StaticUsageWidth returns the wrapping width for static flag descriptions.
func (f *FlagSet) StaticUsageWidth() int { return f.impl.StaticUsageWidth() }

// StaticAutoUsageColumn computes a good usage column for static flags.
func (f *FlagSet) StaticAutoUsageColumn(padding int) int {
	return f.impl.StaticAutoUsageColumn(padding)
}

// SetStaticUsageNote adds a note after the static flag block.
func (f *FlagSet) SetStaticUsageNote(s string) { f.Layout().SetStaticUsageNote(s) }

// StaticUsageNote returns the static flag section note.
func (f *FlagSet) StaticUsageNote() string { return f.impl.StaticUsageNote() }

// SetDynamicUsageIndent sets the indentation for dynamic flag usage lines.
func (f *FlagSet) SetDynamicUsageIndent(n int) { f.Layout().SetDynamicUsageIndent(n) }

// DynamicUsageIndent returns the dynamic flag usage indent.
func (f *FlagSet) DynamicUsageIndent() int { return f.impl.DynamicUsageIndent() }

// SetDynamicUsageColumn sets the column at which dynamic flag descriptions begin.
func (f *FlagSet) SetDynamicUsageColumn(col int) { f.Layout().SetDynamicUsageColumn(col) }

// DynamicUsageColumn returns the description column for dynamic flags.
func (f *FlagSet) DynamicUsageColumn() int { return f.impl.DynamicUsageColumn() }

// SetDynamicUsageWidth sets the max wrapping width for dynamic flags.
func (f *FlagSet) SetDynamicUsageWidth(max int) { f.Layout().SetDynamicUsageWidth(max) }

// DynamicUsageWidth returns the wrapping width for dynamic flag descriptions.
func (f *FlagSet) DynamicUsageWidth() int { return f.impl.DynamicUsageWidth() }

// DynamicAutoUsageColumn computes a good usage column for dynamic flags.
func (f *FlagSet) DynamicAutoUsageColumn(padding int) int {
	return f.impl.DynamicAutoUsageColumn(padding)
}

// SetDynamicUsageNote adds a note after the dynamic flag block.
func (f *FlagSet) SetDynamicUsageNote(s string) { f.Layout().SetDynamicUsageNote(s) }

// DynamicUsageNote returns the dynamic flag section note.
func (f *FlagSet) DynamicUsageNote() string { return f.impl.DynamicUsageNote() }

// SetNoteIndent sets the indentation for help notes.
func (f *FlagSet) SetNoteIndent(n int) { f.Layout().SetNoteIndent(n) }

// NoteIndent returns the note section indentation.
func (f *FlagSet) NoteIndent() int { return f.impl.NoteIndent() }

// SetNoteWidth sets the wrapping width for help notes.
func (f *FlagSet) SetNoteWidth(max int) { f.Layout().SetNoteWidth(max) }

// NoteWidth returns the wrapping width for help notes.
func (f *FlagSet) NoteWidth() int { return f.impl.NoteWidth() }
