package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAlertComment(t *testing.T) {
	tests := []struct {
		name             string
		args             []string
		commentFlag      string
		formatFlag       string
	}{
		{
			name:             "Valid Comment with Default Format",
			args:             []string{"123"},
			commentFlag:      "This is a test comment.",
		},
		{
			name:             "Valid Comment with Markdown Format",
			args:             []string{"456"},
			commentFlag:      "This is a **markdown** comment.",
			formatFlag:       "markdown",
		},
		{
			name:             "Valid Comment with Plain Text Format",
			args:             []string{"789"},
			commentFlag:      "This is a plain text comment.",
			formatFlag:       "plaintext",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			alertCmdState.Comment = test.commentFlag
			alertCmdState.Format = test.formatFlag
			err := commentAlert(nil, test.args)
			assert.Nil(t, err)
		})
	}
}
