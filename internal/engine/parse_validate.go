package engine

import "fmt"

// checkRequirements ensures all flag dependencies are satisfied.
func (f *FlagSet) checkRequirements() error {
	for _, fl := range f.staticFlagsMap {
		req, missing := fl.FirstMissingRequirement(f.staticFlagsMap)
		if missing {
			return fmt.Errorf("--%s requires --%s", fl.Name, req)
		}
	}
	return nil
}

// checkPositionals ensures all positional arguments are valid.
func (f *FlagSet) checkPositionals() error {
	if err := f.validatePositionals(); err != nil {
		return f.handleError(err)
	}
	if err := f.finalizePositionals(); err != nil {
		return f.handleError(err)
	}
	return nil
}

// checkOneOfGroups ensures only one flag per group is set.
func (f *FlagSet) checkOneOfGroups() error {
	if f.showHelp != nil && *f.showHelp {
		return nil
	}
	if f.showVersion != nil && *f.showVersion {
		return nil
	}

	for _, g := range f.oneOfGroup {
		selections := 0
		var conflicting []string
		for _, fl := range g.Flags {
			if fl.Value.Changed() {
				selections++
				if f.oneOfVerbose {
					conflicting = append(conflicting, "--"+fl.Name)
				}
			}
		}
		for _, grp := range g.RequiredGroups {
			changed := 0
			for _, fl := range grp.Flags {
				if fl.Value.Changed() {
					changed++
				}
			}
			if changed == len(grp.Flags) && len(grp.Flags) > 0 {
				selections++
				if f.oneOfVerbose {
					conflicting = append(conflicting, fmt.Sprintf("[%s]", joinFlagNames(grp.Flags)))
				}
			}
		}

		if selections > 1 {
			if f.oneOfVerbose && len(conflicting) > 1 {
				return fmt.Errorf(
					"only one of the flags in group %q may be used: %s",
					g.Name, joinConflictNames(conflicting),
				)
			}
			return fmt.Errorf("only one of the flags in group %q may be used", g.Name)
		}
		if g.IsRequired() && selections == 0 {
			return fmt.Errorf("one of the flags in group %q must be set", g.Name)
		}
	}
	return nil
}

// checkAllOrNone ensures all required flags were set.
func (f *FlagSet) checkAllOrNone() error {
	if f.showHelp != nil && *f.showHelp {
		return nil
	}
	if f.showVersion != nil && *f.showVersion {
		return nil
	}

	for _, g := range f.allOrNoneGroup {
		changed := 0
		for _, fl := range g.Flags {
			if fl.Value.Changed() {
				changed++
			}
		}
		if changed > 0 && changed != len(g.Flags) {
			return fmt.Errorf("flags %s must be set together", joinFlagNames(g.Flags))
		}
		if g.IsRequired() && changed == 0 {
			return fmt.Errorf("flags %s must be set together", joinFlagNames(g.Flags))
		}
	}
	return nil
}

// validatePositionals ensures all positional arguments are valid.
func (f *FlagSet) validatePositionals() error {
	if len(f.positional) == 0 || f.validatePositional == nil {
		return nil
	}
	for _, arg := range f.positional {
		if err := f.validatePositional(arg); err != nil {
			return err
		}
	}
	return nil
}

// finalizePositionals mutates positional arguments.
func (f *FlagSet) finalizePositionals() error {
	if len(f.positional) == 0 || f.finalizePositional == nil {
		return nil
	}
	for i, arg := range f.positional {
		f.positional[i] = f.finalizePositional(arg)
	}
	return nil
}

// checkRequired ensures all required flags were set.
func (f *FlagSet) checkRequired() error {
	for _, fl := range f.staticFlagsMap {
		if fl.MissingRequired() {
			return fmt.Errorf("flag --%s is required", fl.Name)
		}
	}
	return nil
}

// checkRequiredDynamic ensures all required dynamic flags are set for each
// existing instance of every dynamic group. Errors use --group.id.flag format.
func (f *FlagSet) checkRequiredDynamic() error {
	if len(f.dynamicGroups()) == 0 {
		return nil
	}

	for _, g := range f.dynamicGroups() {
		ids := g.Instances()
		if len(ids) == 0 {
			continue
		}

		items := g.Items()
		if len(items) == 0 {
			continue
		}

		for _, id := range ids {
			for field, item := range items {
				bf := item.Flag
				if bf == nil || !bf.Required {
					continue
				}
				if _, ok := item.Value.GetAny(id); !ok {
					return fmt.Errorf("flag --%s.%s.%s is required", g.Name(), id, field)
				}
			}
		}
	}
	return nil
}
