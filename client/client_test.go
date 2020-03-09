package client_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lacework/go-sdk/client"
)

func TestNew(t *testing.T) {
	c, err := client.New("test")
	if assert.Nil(t, err) {
		assert.Equal(t, "v1", c.ApiVersion(), "default API version should be v1")
	}
}

func TestNewWithOptions(t *testing.T) {
	fakeAPI := NewLaceworkServer()
	fakeAPI.ApiVersion = "v2"
	fakeAPI.MockToken("TOKEN")
	defer fakeAPI.Close()

	c, err := client.New("test",
		client.WithURL(fakeAPI.URL()),
		client.WithExpirationTime(1800),
		client.WithApiV2(),
		client.WithTokenFromKeys("KEY", "SECRET"), // this option has to be the last one
	)
	if assert.Nil(t, err) {
		assert.Equal(t, "v2", c.ApiVersion(), "modified API version should be v2")
	}
}
