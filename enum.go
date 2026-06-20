package tinyflags

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/containeroo/tinyflags/internal/dynamic"
	"github.com/containeroo/tinyflags/internal/engine"
	"github.com/containeroo/tinyflags/internal/scalar"
)

type enumValue interface {
	~string |
		~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
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
