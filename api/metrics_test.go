package api_test

import (
	"os"
	"testing"

	"github.com/lacework/go-sdk/v2/api"
	"github.com/lacework/go-sdk/v2/internal/lacework"
	"github.com/stretchr/testify/assert"
)

func TestDisableTelemetry(t *testing.T) {
	os.Setenv("LW_TELEMETRY_DISABLE", "true")
	defer os.Setenv("LW_TELEMETRY_DISABLE", "")

	fakeServer := lacework.MockServer()
	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)

	assert.Nil(t, err)

	err = c.V2.Metrics.Send(api.MetricEvent{})

	assert.Nil(t, err)
}
