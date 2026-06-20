package tinyflags

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/containeroo/tinyflags/internal/dynamic"
	"github.com/containeroo/tinyflags/internal/engine"
	"github.com/containeroo/tinyflags/internal/scalar"
)

type enumValue interface {
	~string |
		~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

// EnumChoice maps one user-facing enum name to its typed value.
type EnumChoice[T enumValue] struct {
	Name  string
	Value T
}

// Choice returns an enum choice for named enum helpers.
func Choice[T enumValue](name string, value T) EnumChoice[T] {
	return EnumChoice[T]{Name: name, Value: value}
}

// EnumVar defines a typed enum flag and binds it to the given pointer.
// Only values listed in allowed are accepted.
func EnumVar[T enumValue](f *FlagSet, ptr *T, name string, def T, usage string, allowed ...T) *scalar.ScalarFlag[T] {
	flag := engine.RegisterStaticScalar(
		f.impl,
		ptr,
		name,
		usage,
		def,
		parseEnumValue[T],
		formatEnumValue[T],
	)
	return flag.Choices(allowed...)
}

// Enum defines a typed enum flag and returns its handle.
// Only values listed in allowed are accepted.
func Enum[T enumValue](f *FlagSet, name string, def T, usage string, allowed ...T) *scalar.ScalarFlag[T] {
	return EnumVar(f, new(T), name, def, usage, allowed...)
}

// DynamicEnum defines a typed enum flag on a dynamic group.
// Only values listed in allowed are accepted.
func DynamicEnum[T enumValue](g *DynamicGroup, field string, def T, usage string, allowed ...T) *dynamic.ScalarFlag[T] {
	return dynamic.Enum(g, field, def, usage, allowed...)
}

// EnumMapVar defines a typed enum flag with user-facing names and binds it to the given pointer.
func EnumMapVar[T enumValue](f *FlagSet, ptr *T, name string, def T, usage string, choices ...EnumChoice[T]) *scalar.ScalarFlag[T] {
	parse, format, allowed := enumChoiceHooks(choices)
	flag := engine.RegisterStaticScalar(f.impl, ptr, name, usage, def, parse, format)
	flag.Allowed(allowed...)
	return flag
}

// EnumMap defines a typed enum flag with user-facing names and returns its handle.
func EnumMap[T enumValue](f *FlagSet, name string, def T, usage string, choices ...EnumChoice[T]) *scalar.ScalarFlag[T] {
	return EnumMapVar(f, new(T), name, def, usage, choices...)
}

// DynamicEnumMap defines a typed enum flag with user-facing names on a dynamic group.
func DynamicEnumMap[T enumValue](g *DynamicGroup, field string, def T, usage string, choices ...EnumChoice[T]) *dynamic.ScalarFlag[T] {
	return dynamic.EnumMap(g, field, def, usage, dynamicEnumChoices(choices)...)
}

func parseEnumValue[T enumValue](raw string) (T, error) {
	var zero T
	typ := reflect.TypeOf(zero)
	out := reflect.New(typ).Elem()

	switch typ.Kind() {
	case reflect.String:
		out.SetString(raw)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v, err := strconv.ParseInt(raw, 10, typ.Bits())
		if err != nil {
			return zero, err
		}
		out.SetInt(v)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v, err := strconv.ParseUint(raw, 10, typ.Bits())
		if err != nil {
			return zero, err
		}
		out.SetUint(v)
	default:
		return zero, fmt.Errorf("unsupported enum kind %s", typ.Kind())
	}

	return out.Interface().(T), nil
}

func formatEnumValue[T enumValue](v T) string {
	value := reflect.ValueOf(v)

	switch value.Kind() {
	case reflect.String:
		return value.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(value.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(value.Uint(), 10)
	default:
		return fmt.Sprintf("%v", v)
	}
}

func enumChoiceHooks[T enumValue](choices []EnumChoice[T]) (func(string) (T, error), func(T) string, []string) {
	names := make([]string, 0, len(choices))
	byName := make(map[string]T, len(choices))
	byValue := make(map[T]string, len(choices))

	for _, choice := range choices {
		names = append(names, choice.Name)
		byName[choice.Name] = choice.Value
		if _, exists := byValue[choice.Value]; !exists {
			byValue[choice.Value] = choice.Name
		}
	}

	parse := func(raw string) (T, error) {
		val, ok := byName[raw]
		if ok {
			return val, nil
		}
		var zero T
		return zero, fmt.Errorf("must be one of: %s", strings.Join(names, ", "))
	}

	format := func(v T) string {
		name, ok := byValue[v]
		if ok {
			return name
		}
		return formatEnumValue(v)
	}

	return parse, format, names
}

func dynamicEnumChoices[T enumValue](choices []EnumChoice[T]) []dynamic.EnumChoice[T] {
	out := make([]dynamic.EnumChoice[T], len(choices))
	for i, choice := range choices {
		out[i] = dynamic.EnumChoice[T]{Name: choice.Name, Value: choice.Value}
	}
	return out
}
