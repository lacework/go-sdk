package array

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSort2DSlice(t *testing.T) {
	stringsTableTest := []struct {
		Input    [][]string
		Expected [][]string
	}{
		{
			Input:    [][]string{{"RoleArn", "arn:aws:iam::1234"}, {"ExternalID", "=abcd123"}, {"UpdatedBy", "TestUser"}, {"UpdatedAt", "2021-10-25T:15:18"}, {"StateDetails", "example details"}},
			Expected: [][]string{{"ExternalID", "=abcd123"}, {"RoleArn", "arn:aws:iam::1234"}, {"StateDetails", "example details"}, {"UpdatedAt", "2021-10-25T:15:18"}, {"UpdatedBy", "TestUser"}},
		},
		{
			Input:    [][]string{{"A", "example"}, {"D", "example"}, {"F", "example"}, {"E", "example"}, {"C", "example"}, {"B", "example"}},
			Expected: [][]string{{"A", "example"}, {"B", "example"}, {"C", "example"}, {"D", "example"}, {"E", "example"}, {"F", "example"}},
		},
	}

	for _, test := range stringsTableTest {
		t.Run("sort 2d array", func(t *testing.T) {
			Sort2D(test.Input)
			assert.Equal(t, test.Expected, test.Input)
		})
	}
}
