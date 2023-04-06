package api_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/lacework"
	"github.com/stretchr/testify/assert"
)

func TestGetFeatureFlagsMatchingPrefix(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.UseApiV2()
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()
	fakeServer.MockAPI("FeatureFlags/PUBLIC.sca", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		_, err := fmt.Fprintf(w, generateResponse())
		assert.Nil(t, err)
	})
	c, err := api.NewClient("test",
		api.WithApiV2(),
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)
	response, err := c.V2.FeatureFlags.GetFeatureFlagsMatchingPrefix("PUBLIC.sca")
	assert.Nil(t, err)
	assert.Equal(t, response.Data, api.FeatureFlags{Flags: []api.FeatureFlag{"flag1", "flag2", "flag3"}})
}

func generateResponse() string {
	return `
{
	"data": {
		"flags": [
			"flag1",
			"flag2",
			"flag3"
		]
	}
}
`
}
