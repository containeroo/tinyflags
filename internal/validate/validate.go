package validate

import (
	"fmt"
	"strings"

	"github.com/containeroo/tinyflags/internal/core"
	"github.com/containeroo/tinyflags/internal/dynamic"
)

// CheckRequirements ensures all flag dependencies are satisfied.
func CheckRequirements(flags map[string]*core.BaseFlag) error {
	for _, fl := range flags {
		req, missing := fl.FirstMissingRequirement(flags)
		if missing {
			return fmt.Errorf("--%s requires --%s", fl.Name, req)
		}
	}
	return nil
}

// CheckRequired ensures all required static flags were set.
func CheckRequired(flags map[string]*core.BaseFlag) error {
	for _, fl := range flags {
		if fl.MissingRequired() {
			return fmt.Errorf("flag --%s is required", fl.Name)
		}
	}
	return nil
}

// CheckRequiredDynamic ensures all required dynamic flags are set for seen IDs.
func CheckRequiredDynamic(groups []*dynamic.Group) error {
	for _, g := range groups {
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

// CheckOneOfGroups ensures at most one choice per one-of group is set.
func CheckOneOfGroups(groups []*core.OneOfGroupGroup, verbose bool) error {
	for _, g := range groups {
		selections := 0
		var conflicting []string
		for _, fl := range g.Flags {
			if fl.Value.Changed() {
				selections++
				if verbose {
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
				if verbose {
					conflicting = append(conflicting, fmt.Sprintf("[%s]", joinFlagNames(grp.Flags)))
				}
			}
		}

		if selections > 1 {
			if verbose && len(conflicting) > 1 {
				return fmt.Errorf(
					"only one of the flags in group %q may be used: %s",
					g.Name, strings.Join(conflicting, " vs "),
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

// CheckAllOrNone ensures all-or-none groups are fully satisfied.
func CheckAllOrNone(groups []*core.AllOrNoneGroup) error {
	for _, g := range groups {
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

// ValidatePositionals ensures all positional arguments are valid.
func ValidatePositionals(positional []string, validate func(string) error) error {
	if len(positional) == 0 || validate == nil {
		return nil
	}
	for _, arg := range positional {
		if err := validate(arg); err != nil {
			return err
		}
	}
	return nil
}

// FinalizePositionals mutates positional arguments in place.
func FinalizePositionals(positional []string, finalize func(string) string) error {
	if len(positional) == 0 || finalize == nil {
		return nil
	}
	for i, arg := range positional {
		positional[i] = finalize(arg)
	}
	return nil
}

func joinFlagNames(flags []*core.BaseFlag) string {
	names := make([]string, 0, len(flags))
	for _, fl := range flags {
		names = append(names, "--"+fl.Name)
	}
	return strings.Join(names, ", ")
}
