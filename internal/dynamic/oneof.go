package dynamic

import (
	"fmt"
	"slices"
	"strings"

	"github.com/containeroo/tinyflags/internal/core"
)

// Field identifies a registered dynamic scalar, boolean, or slice without exposing its value type.
type Field interface {
	// dynamicValue returns the registered per-instance storage, sealing field handles to this package.
	dynamicValue() core.DynamicValue
}

// OneOfGroup constrains a set of field handles independently for each active instance.
type OneOfGroup struct {
	// owner identifies the dynamic group whose instances are validated.
	owner *Group
	// fields preserves the registered choices in declaration order.
	fields []core.DynamicValue
	// help supplies the existing one-of annotations without registering a global constraint.
	help *core.OneOfGroupGroup
}

// OneOf permits at most one supplied field per instance and panics on invalid field handles.
func (g *Group) OneOf(fields ...Field) *OneOfGroup {
	if len(fields) == 0 {
		panic("tinyflags: dynamic OneOf requires at least one field")
	}
	choice := &OneOfGroup{owner: g, help: &core.OneOfGroupGroup{}}
	var names []string
	for _, field := range fields {
		value := g.oneOfField(field)
		if slices.Contains(choice.fields, value) {
			panic("tinyflags: duplicate field in dynamic OneOf: " + value.FieldName())
		}
		choice.fields = append(choice.fields, value)
		names = append(names, value.FieldName())
	}
	choice.help.Name = strings.Join(names, " | ") + " per instance"
	for _, value := range choice.fields {
		flag := g.items[value.FieldName()].Flag
		choice.help.Flags = append(choice.help.Flags, flag)
		flag.OneOfGroups = append(flag.OneOfGroups, choice.help)
	}
	g.oneOf = append(g.oneOf, choice)
	return choice
}

// oneOfField rejects nil, zero-value, and foreign-group handles before registration changes state.
func (g *Group) oneOfField(field Field) core.DynamicValue {
	if field == nil {
		panic("tinyflags: nil field in dynamic OneOf")
	}
	value := field.dynamicValue()
	if value == nil {
		panic("tinyflags: nil field in dynamic OneOf")
	}
	item, exists := g.items[value.FieldName()]
	if !exists {
		panic("tinyflags: dynamic OneOf field does not belong to group " + g.name)
	}
	if !sameFieldStorage(item.Value, value) {
		panic("tinyflags: dynamic OneOf field does not belong to group " + g.name)
	}
	return item.Value
}

// Required requires exactly one supplied member for every active instance of the owning group.
func (g *OneOfGroup) Required() *OneOfGroup {
	g.help.Required()
	return g
}

// ValidateOneOf checks per-instance choices after command-line and environment parsing finish.
func (g *Group) ValidateOneOf() error {
	if len(g.oneOf) == 0 {
		return nil
	}
	for _, id := range g.Instances() {
		for _, choice := range g.oneOf {
			if err := choice.validate(id); err != nil {
				return err
			}
		}
	}
	return nil
}

// validate counts explicit presence rather than truthiness or nonzero values for one instance.
func (g *OneOfGroup) validate(id string) error {
	count := 0
	for _, field := range g.fields {
		if _, supplied := field.GetAny(id); supplied {
			count++
		}
	}
	if count > 1 {
		return fmt.Errorf("only one of %s may be set", g.fieldNames(id))
	}
	if !g.help.IsRequired() {
		return nil
	}
	if count == 0 {
		return fmt.Errorf("one of %s must be set", g.fieldNames(id))
	}
	return nil
}

// fieldNames formats qualified member names only when a constraint needs a diagnostic.
func (g *OneOfGroup) fieldNames(id string) string {
	names := make([]string, 0, len(g.fields))
	for _, field := range g.fields {
		names = append(names, "--"+g.owner.name+"."+id+"."+field.FieldName())
	}
	return strings.Join(names, ", ")
}

// dynamicValue exposes a scalar's storage to handle-based constraints without generic type assertions.
func (f *ScalarFlag[T]) dynamicValue() core.DynamicValue {
	if f == nil {
		return nil
	}
	if f.item == nil {
		return nil
	}
	return f.item
}

// dynamicValue exposes a slice's storage to handle-based constraints without inspecting its elements.
func (f *SliceFlag[T]) dynamicValue() core.DynamicValue {
	if f == nil {
		return nil
	}
	if f.item == nil {
		return nil
	}
	return f.item
}

// dynamicValue exposes a boolean's storage while keeping explicitly false values distinguishable.
func (f *BoolFlag) dynamicValue() core.DynamicValue {
	if f == nil {
		return nil
	}
	if f.item == nil {
		return nil
	}
	return f.item
}

// sameFieldStorage recognizes the boolean wrapper around a handle's scalar storage.
func sameFieldStorage(registered, candidate core.DynamicValue) bool {
	if registered == candidate {
		return true
	}
	boolean, wrapped := registered.(*BoolValue)
	if !wrapped {
		return false
	}
	return boolean.DynamicScalarValue == candidate
}
