package client_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lacework/go-sdk/client"
)

func TestWithApiV2(t *testing.T) {
	c, err := client.New("test", client.WithApiV2())
	if assert.Nil(t, err) {
		assert.Equal(t, "v2", c.ApiVersion(), "API version should be v2")
	}
}

func TestWithToken(t *testing.T) {
	c, err := client.New("test", client.WithToken("TOKEN"))
	if assert.Nil(t, err) {
		assert.Equal(t, "v1", c.ApiVersion(), "API version should be v2")
	}
}

func TestApiVersion(t *testing.T) {
	c, err := client.New("foo")
	if assert.Nil(t, err) {
		assert.Equal(t, "v1", c.ApiVersion(), "wrong default API version")
	}
}
