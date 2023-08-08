package format

import (
	"bytes"
	"fmt"
	"unicode"
)

// SpaceUpperCase add a space each occurrence of an upper case letter.
//
// res := format.SpaceUpperCase("myExampleString")
// -> my Example String
func SpaceUpperCase(s string) string {
	buf := &bytes.Buffer{}
	prev := false
	spaces := []*unicode.RangeTable{unicode.White_Space, unicode.Space}
	for i, rune := range s {
		if unicode.IsOneOf(spaces, rune) {
			prev = true
		}

		// if previous value is a space skip
		if prev {
			prev = false
			continue
		}
		if unicode.IsUpper(rune) && i > 0 {
			buf.WriteRune(' ')
		}
		buf.WriteRune(rune)
	}
	return buf.String()
}

// Truncate reduce a string to specified char limit. Append ellipsis.
//
// res := format.Truncate("myExampleString", 10)
//
//	-> myExampleS...
func Truncate(s string, max int) string {
	if max > len(s) {
		return s
	}

	return fmt.Sprintf("%s...", s[:max])
}
