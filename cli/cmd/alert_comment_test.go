package cmd

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/lacework"
	"github.com/stretchr/testify/assert"
)

func TestAlertComment(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		commentFlag string
		formatFlag  string
	}{
		{
			name:        "Valid Comment with Markdown Format",
			args:        []string{"12345"},
			commentFlag: "This is a **markdown** comment.",
			formatFlag:  "markdown",
		},
		{
			name:        "Valid Comment with Plain Text Format",
			args:        []string{"12344"},
			commentFlag: "This is a plain text comment.",
			formatFlag:  "plaintext",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			id, _ := strconv.Atoi(test.args[0])
			fakeServer := lacework.MockServer()
			fakeServer.MockAPI(
				fmt.Sprintf("Alerts/%d/comment", id),
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "POST", r.Method, "should be a POST method")
					fmt.Fprint(w, "{}")
				},
			)
			defer fakeServer.Close()

			c, err := api.NewClient("test",
				api.WithToken("TOKEN"),
				api.WithURL(fakeServer.URL()),
			)
			assert.Nil(t, err)

			_, err = c.V2.Alerts.CommentFromRequest(id,
				api.AlertsCommentRequest{Comment: test.commentFlag, Format: test.formatFlag})
			assert.Nil(t, err)
		})
	}
}

func TestInvalidCommentNoFormat(t *testing.T) {
	alertCmdState.Comment = "This is a test comment."
	alertCmdState.Format = ""
	err := commentAlert(nil, []string{"12545"})
	assert.Equal(t, err.Error(), "unable to process alert comment format: EOF")
}
