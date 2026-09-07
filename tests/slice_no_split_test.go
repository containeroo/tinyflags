package tinyflags_test

import (
	"strconv"
	"strings"
	"testing"

	"github.com/containeroo/tinyflags"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSliceNoSplit(t *testing.T) {
	t.Parallel()
	for _, splitAgain := range []bool{false, true} {
		name := "disabled"
		if splitAgain {
			name = "reenabled"
		}
		t.Run(name, func(t *testing.T) {
			fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
			static := fs.StringSlice("values", []string{"default"}, "Values").Delimiter("|").NoSplit().TrimSpace()
			dynamic := fs.DynamicGroup("http").StringSlice("values", []string{"default"}, "Values").Delimiter("|").NoSplit().TrimSpace()
			if splitAgain {
				static.Delimiter("|")
				dynamic.Delimiter("|")
			}
			args := []string{"--values= a|b ", "--values=c,d", "--http.web.values= a|b ", "--http.web.values=c,d"}
			expected := []string{"a|b", "c,d"}
			if splitAgain {
				expected = []string{"a", "b", "c,d"}
			}
			for range 2 {
				require.NoError(t, fs.Parse(args))
				assert.Equal(t, expected, *static.Value())
				assert.Equal(t, expected, dynamic.MustGet("web"))
			}
		})
	}
}

func TestSliceNoSplitEmptyAndWhitespace(t *testing.T) {
	t.Parallel()
	for _, allowEmpty := range []bool{false, true} {
		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		static := fs.StringSlice("values", nil, "Values").NoSplit().PreserveSpace()
		dynamic := fs.DynamicGroup("http").StringSlice("values", nil, "Values").NoSplit().PreserveSpace()
		if allowEmpty {
			static.AllowEmpty()
			dynamic.AllowEmpty()
		}
		require.NoError(t, fs.Parse([]string{"--values= a,b ", "--http.web.values= a,b "}))
		assert.Equal(t, []string{" a,b "}, *static.Value())
		assert.Equal(t, []string{" a,b "}, dynamic.MustGet("web"))
		for _, arg := range []string{"--values=", "--http.web.values="} {
			err := fs.Parse([]string{arg})
			if allowEmpty {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, "empty values are not allowed")
			}
		}
	}
}

func TestCustomDynamicSliceNoSplit(t *testing.T) {
	t.Parallel()
	fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
	var inputs []string
	flag := tinyflags.DynamicSlice(fs.DynamicGroup("http"), "codes", []int{200}, "Codes",
		func(raw string) ([]int, error) {
			inputs = append(inputs, raw)
			var out []int
			for _, part := range strings.Split(raw, ",") {
				codes, err := parseStatusCodes(part)
				if err != nil {
					return nil, err
				}
				out = append(out, codes...)
			}
			return out, nil
		}, strconv.Itoa).NoSplit()
	require.NoError(t, fs.Parse([]string{"--http.web.codes=200-202,304", "--http.web.codes=404"}))
	assert.Equal(t, []string{"200-202,304", "404"}, inputs)
	assert.Equal(t, []int{200, 201, 202, 304, 404}, flag.MustGet("web"))
}
