package tinyflags

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"strings"
)

type numeric interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

type ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64 | ~string
}

// Positive returns a validator that requires a numeric value greater than zero.
func Positive[T numeric]() func(T) error {
	return func(value T) error {
		var zero T
		if isNaN(value) || value <= zero {
			return errors.New("must be greater than zero")
		}
		return nil
	}
}

// NonNegative returns a validator that requires a numeric value greater than or equal to zero.
func NonNegative[T numeric]() func(T) error {
	return func(value T) error {
		var zero T
		if isNaN(value) || value < zero {
			return errors.New("must not be negative")
		}
		return nil
	}
}

// AtLeast returns a validator that requires a value greater than or equal to min.
func AtLeast[T ordered](min T) func(T) error {
	if isNaN(min) {
		panic("tinyflags: AtLeast minimum must be an ordered value")
	}

	return func(value T) error {
		if isNaN(value) || value < min {
			return fmt.Errorf("must be at least %v", min)
		}
		return nil
	}
}

// AtMost returns a validator that requires a value less than or equal to max.
func AtMost[T ordered](max T) func(T) error {
	if isNaN(max) {
		panic("tinyflags: AtMost maximum must be an ordered value")
	}

	return func(value T) error {
		if isNaN(value) || value > max {
			return fmt.Errorf("must be at most %v", max)
		}
		return nil
	}
}

// Between returns a validator that requires a value to be within the inclusive min/max range.
// It panics if min is greater than max or either bound is not ordered.
func Between[T ordered](min, max T) func(T) error {
	if isNaN(min) || isNaN(max) || min > max {
		panic("tinyflags: Between minimum must be less than or equal to maximum")
	}

	return func(value T) error {
		if isNaN(value) || value < min || value > max {
			return fmt.Errorf("must be between %v and %v (inclusive)", min, max)
		}
		return nil
	}
}

func isNaN[T ordered](value T) bool {
	rv := reflect.ValueOf(value)
	if rv.Kind() != reflect.Float32 && rv.Kind() != reflect.Float64 {
		return false
	}

	return math.IsNaN(rv.Float())
}

// NotBlank returns a validator that rejects empty or whitespace-only strings.
func NotBlank() func(string) error {
	return func(value string) error {
		if strings.TrimSpace(value) == "" {
			return errors.New("must not be blank")
		}
		return nil
	}
}
