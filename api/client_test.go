package api_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lacework/go-sdk/api"
)

func TestNewClient(t *testing.T) {
	c, err := api.NewClient("test")
	if assert.Nil(t, err) {
		assert.Equal(t, "v1", c.ApiVersion(), "default API version should be v1")
	}
}

func TestNewClientWithOptions(t *testing.T) {
	fakeAPI := NewLaceworkServer()
	fakeAPI.ApiVersion = "v2"
	fakeAPI.MockToken("TOKEN")
	defer fakeAPI.Close()

	c, err := api.NewClient("test",
		api.WithURL(fakeAPI.URL()),
		api.WithExpirationTime(1800),
		api.WithApiV2(),
		api.WithTokenFromKeys("KEY", "SECRET"), // this option has to be the last one
	)
	if assert.Nil(t, err) {
		assert.Equal(t, "v2", c.ApiVersion(), "modified API version should be v2")
	}
}
