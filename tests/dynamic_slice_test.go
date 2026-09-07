package tinyflags_test

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/containeroo/tinyflags"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func parseStatusCodes(raw string) ([]int, error) {
	if raw == "none" {
		return nil, nil
	}
	first, last, isRange := strings.Cut(raw, "-")
	start, err := strconv.Atoi(first)
	if err != nil {
		return nil, err
	}
	end := start
	if isRange {
		end, err = strconv.Atoi(last)
		if err != nil {
			return nil, err
		}
	}
	if start < 100 || end > 599 || end < start {
		return nil, fmt.Errorf("invalid status range")
	}
	codes := make([]int, 0, end-start+1)
	for code := start; code <= end; code++ {
		codes = append(codes, code)
	}
	return codes, nil
}

func TestCustomDynamicSlice(t *testing.T) {
	t.Parallel()
	fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
	group := fs.DynamicGroup("http")
	flag := tinyflags.DynamicSlice(group, "expected-status-codes", []int{200}, "Expected HTTP status codes", parseStatusCodes, strconv.Itoa)
	require.NoError(t, fs.Parse([]string{
		"--http.web.expected-status-codes=200-299, 304",
		"--http.web.expected-status-codes=404",
		"--http.api.expected-status-codes=201",
		"--http.empty.expected-status-codes=none",
	}))
	expected := make([]int, 100)
	for i := range expected {
		expected[i] = 200 + i
	}
	expected = append(expected, 304, 404)
	assert.Equal(t, expected, tinyflags.GetOrDefaultDynamic[[]int](group, "web", "expected-status-codes"))
	assert.Equal(t, []int{201}, flag.MustGet("api"))
	assert.Equal(t, []int{200}, tinyflags.GetOrDefaultDynamic[[]int](group, "missing", "expected-status-codes"))
	assert.Empty(t, tinyflags.GetOrDefaultDynamic[[]int](group, "empty", "expected-status-codes"))
	assert.True(t, flag.Has("empty"))
	assert.True(t, flag.Changed())
	require.NoError(t, fs.Parse(nil))
	assert.False(t, flag.Has("web"))
	assert.False(t, flag.Changed())
	assert.Equal(t, []int{200}, tinyflags.GetOrDefaultDynamic[[]int](group, "web", "expected-status-codes"))
}

func TestCustomDynamicSliceHooks(t *testing.T) {
	t.Parallel()
	fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
	group := fs.DynamicGroup("http")
	var validated []int
	flag := tinyflags.DynamicSlice(group, "codes", []int{200}, "Status codes", parseStatusCodes, strconv.Itoa).
		Delimiter(";").
		Validate(func(v int) error { validated = append(validated, v); return nil }).
		Finalize(func(v int) int { return v + 1 }).
		FinalizeWithID(func(id string, v int) int { require.Equal(t, "web", id); return v + 10 })
	require.NoError(t, fs.Parse([]string{"--http.web.codes=200-202; 304"}))
	assert.Equal(t, []int{200, 201, 202, 304}, validated)
	assert.Equal(t, []int{211, 212, 213, 315}, flag.MustGet("web"))
}

func TestCustomDynamicSliceErrors(t *testing.T) {
	t.Parallel()
	for _, raw := range []string{"bad", "299-200", "200-202", "200,,201"} {
		t.Run(raw, func(t *testing.T) {
			fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
			flag := tinyflags.DynamicSlice(fs.DynamicGroup("http"), "codes", []int{200}, "Status codes", parseStatusCodes, strconv.Itoa)
			flag.Validate(func(v int) error {
				if v == 201 {
					return fmt.Errorf("rejected status")
				}
				return nil
			})
			require.ErrorContains(t, fs.Parse([]string{"--http.web.codes=" + raw}), "invalid value")
		})
	}
}
