package unique

import (
	"sort"
	"testing"
)

func TestStringSlice(t *testing.T) {
	tests := []struct {
		description string
		input       []string
		expected    []string
	}{
		{
			"Duplicates",
			[]string{"abc", "abc", "def", "abc", "def"},
			[]string{"abc", "def"},
		},
		{
			"Empty slice",
			[]string{},
			[]string{},
		},
		{
			"No duplicates",
			[]string{"abc", "def", "rst", "xyz"},
			[]string{"abc", "def", "rst", "xyz"},
		},
	}

	for _, testCase := range tests {
		result := StringSlice(testCase.input)

		sort.Strings(result)

		for i, res := range result {
			if res != testCase.expected[i] {
				t.Errorf("Test case `%s` failed", testCase.description)
			}
		}
	}
}
