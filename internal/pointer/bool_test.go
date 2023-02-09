package pointer_test

import (
	"testing"

	"github.com/aws/smithy-go/ptr"
	"github.com/lacework/go-sdk/internal/pointer"
	"github.com/stretchr/testify/assert"
)

func TestCompareBoolPtr(t *testing.T) {
	boolPtrTableTests := []struct {
		Name       string
		Input      *bool
		Comparison bool
		Expected   bool
	}{
		{
			Name:       "nil bool ptr returns false",
			Input:      nil,
			Comparison: true,
			Expected:   false,
		},
		{
			Name:       "true bool ptr returns true when comparison is true",
			Input:      ptr.Bool(true),
			Comparison: true,
			Expected:   true,
		},
		{
			Name:       "false bool ptr returns false when comparison is true",
			Input:      ptr.Bool(false),
			Comparison: true,
			Expected:   false,
		},
		{
			Name:       "false bool ptr returns true when comparison is false",
			Input:      ptr.Bool(false),
			Comparison: false,
			Expected:   true,
		},
		{
			Name:       "true bool ptr returns false when comparison is false",
			Input:      ptr.Bool(true),
			Comparison: false,
			Expected:   false,
		},
	}

	for _, ptt := range boolPtrTableTests {
		t.Run(ptt.Name, func(t *testing.T) {
			res := pointer.CompareBoolPtr(ptt.Input, ptt.Comparison)
			assert.Equal(t, res, ptt.Expected)
		})
	}
}
