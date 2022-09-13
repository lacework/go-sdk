package format

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSpaceUpperCase(t *testing.T) {
	stringsTableTest := []struct {
		Input    string
		Expected string
	}{
		{
			Input:    "myExampleString",
			Expected: "my Example String",
		},
		{
			Input:    "MyExampleString",
			Expected: "My Example String",
		},
		{
			Input:    "my Example String",
			Expected: "my Example String",
		},

		{
			Input:    "my ExAmplE StRing",
			Expected: "my Ex Ampl E St Ring",
		},
		{
			Input:    "",
			Expected: "",
		},
	}

	for _, test := range stringsTableTest {
		t.Run("Add space before upper case", func(t *testing.T) {
			result := SpaceUpperCase(test.Input)
			assert.Equal(t, result, test.Expected)

		})
	}
}

func TestTruncate(t *testing.T) {
	stringsTableTest := []struct {
		Input    string
		Max      int
		Expected string
	}{
		{
			Input:    longString,
			Expected: "AWS Error...",
			Max:      9,
		},
		{
			Input:    longString,
			Expected: "AWS Error: Error assuming role. INTG_GUID:TECHALLY...",
			Max:      50,
		},
		{
			Input:    longString,
			Expected: "...",
			Max:      0,
		},
		{
			Input:    shortString,
			Expected: "A...",
			Max:      1,
		},
		{
			Input:    shortString,
			Expected: "AWS",
			Max:      50,
		},
		{
			Input:    "",
			Expected: "",
			Max:      5,
		},
	}

	for _, test := range stringsTableTest {
		t.Run("Add space before upper case", func(t *testing.T) {
			result := Truncate(test.Input, test.Max)
			assert.Equal(t, result, test.Expected)

		})
	}
}

var longString = "AWS Error: Error assuming role. INTG_GUID:TECHALLY_E839836BC385C452E68B3CA7EB45BA0E7BDA39CCF65673A [Error msg: 403 Forbidden]"
var shortString = "AWS"
