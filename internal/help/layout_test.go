package help

import "testing"

// TestDescriptionWidthRespectsConfiguredMaximum verifies labels leave only the remaining width for text.
func TestDescriptionWidthRespectsConfiguredMaximum(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		layout layout
		want   int
	}{
		{
			name:   "uses remaining width",
			layout: layout{indent: 2, startCol: 20, maxWidth: 40},
			want:   17,
		},
		{
			name:   "uses default width when unset",
			layout: layout{indent: 2, startCol: 20},
			want:   100,
		},
		{
			name:   "keeps a one-column minimum",
			layout: layout{indent: 2, startCol: 20, maxWidth: 10},
			want:   1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			if got := test.layout.descriptionWidth(); got != test.want {
				t.Fatalf("descriptionWidth() = %d, want %d", got, test.want)
			}
		})
	}
}
