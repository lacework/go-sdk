package api_test

import (
	"os"
	"testing"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/lacework"
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

	res, err := c.V2.Metrics.Send(api.Honeyvent{})

	assert.Nil(t, err)
	assert.Equal(t, res.Data[0].TraceID, "Telemetry Disabled")
}
