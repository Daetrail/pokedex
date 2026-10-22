package util

import (
	"slices"
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  []string
	}{
		{
			name:  "extra space between words",
			input: "  hello world  ",
			want:  []string{"hello", "world"},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := CleanInput(c.input)
			if !slices.Equal(got, c.want) {
				t.Errorf("got: %v, want: %v", got, c.want)
			}
		})
	}
}
