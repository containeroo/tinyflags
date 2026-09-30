package utils

import (
	"math"
	"testing"
)

func TestParseBytes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  uint64
	}{
		{name: "plainBytes", input: "33554432", want: 33_554_432},
		{name: "bytesSuffix", input: "512B", want: 512},
		{name: "decimalKilobytes", input: "32KB", want: 32_000},
		{name: "decimalMegabytes", input: "32MB", want: 32_000_000},
		{name: "binaryKibibytes", input: "32KiB", want: 32 << 10},
		{name: "binaryMebibytes", input: "32MiB", want: 32 << 20},
		{name: "binaryGibibytes", input: "2GiB", want: 2 << 30},
		{name: "shortUnit", input: "512M", want: 512_000_000},
		{name: "caseInsensitive", input: "64mib", want: 64 << 20},
		{name: "surroundingSpace", input: " 32 MiB ", want: 32 << 20},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := ParseBytes(tt.input)
			if err != nil {
				t.Fatalf("ParseBytes(%q) returned error: %v", tt.input, err)
			}
			if got != tt.want {
				t.Fatalf("ParseBytes(%q) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}

func TestParseBytesRejectsInvalidValues(t *testing.T) {
	t.Parallel()

	tests := []string{
		"",
		"MiB",
		"-1MiB",
		"1.5GiB",
		"32XB",
		"18446744073709551615EiB",
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			t.Parallel()
			if _, err := ParseBytes(input); err == nil {
				t.Fatalf("ParseBytes(%q) unexpectedly succeeded", input)
			}
		})
	}
}

func TestParseBytesMaximumValue(t *testing.T) {
	t.Parallel()

	got, err := ParseBytes("18446744073709551615")
	if err != nil {
		t.Fatalf("ParseBytes(max uint64) returned error: %v", err)
	}
	if got != uint64(math.MaxUint64) {
		t.Fatalf("ParseBytes(max uint64) = %d, want %d", got, uint64(math.MaxUint64))
	}
}
