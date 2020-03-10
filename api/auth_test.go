package api_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lacework/go-sdk/api"
)

func TestWithApiV2(t *testing.T) {
	c, err := api.NewClient("test", api.WithApiV2())
	if assert.Nil(t, err) {
		assert.Equal(t, "v2", c.ApiVersion(), "API version should be v2")
	}
}

func TestWithToken(t *testing.T) {
	c, err := api.NewClient("test", api.WithToken("TOKEN"))
	if assert.Nil(t, err) {
		assert.Equal(t, "v1", c.ApiVersion(), "API version should be v2")
	}
}

func TestApiVersion(t *testing.T) {
	c, err := api.NewClient("foo")
	if assert.Nil(t, err) {
		assert.Equal(t, "v1", c.ApiVersion(), "wrong default API version")
	}
}

func TestWithApiKeys(t *testing.T) {
	c, err := api.NewClient("foo", api.WithApiKeys("KEY", "SECRET"))
	if assert.Nil(t, err) {
		assert.Equal(t, "v1", c.ApiVersion(), "wrong default API version")
	}
}

func TestWithTokenFromKeys(t *testing.T) {
	fakeAPI := NewLaceworkServer()
	fakeAPI.MockToken("TOKEN")
	defer fakeAPI.Close()

	c, err := api.NewClient("foo",
		api.WithURL(fakeAPI.URL()),
		api.WithTokenFromKeys("KEY", "SECRET"), // this option has to be the last one
	)
	if assert.Nil(t, err) {
		assert.Equal(t, "v1", c.ApiVersion(), "wrong default API version")
	}
}

func TestGenerateToken(t *testing.T) {
	fakeAPI := NewLaceworkServer()
	fakeAPI.MockToken("TOKEN")
	defer fakeAPI.Close()

	c, err := api.NewClient("foo",
		api.WithURL(fakeAPI.URL()),
		api.WithApiKeys("KEY", "SECRET"),
	)
	if assert.Nil(t, err) {
		response, err := c.GenerateToken()
		assert.Nil(t, err)
		assert.Equal(t, "TOKEN", response.Token(), "token mismatch")
	}
}

func TestGenerateTokenWithKeys(t *testing.T) {
	fakeAPI := NewLaceworkServer()
	fakeAPI.MockToken("TOKEN")
	defer fakeAPI.Close()

	c, err := api.NewClient("foo", api.WithURL(fakeAPI.URL()))
	if assert.Nil(t, err) {
		response, err := c.GenerateTokenWithKeys("KEY", "SECRET")
		assert.Nil(t, err)
		assert.Equal(t, "TOKEN", response.Token(), "token mismatch")
	}
}

func TestGenerateTokenErrorKeysMissing(t *testing.T) {
	c, err := api.NewClient("where-are-my-keys")
	if assert.Nil(t, err) {
		response, err := c.GenerateToken()
		if assert.NotNil(t, err) {
			assert.Empty(t, response, "token must be empty")
			assert.Equal(t,
				"unable to generate access token: auth keys missing",
				err.Error(),
				"error message mismatch",
			)
		}
	}
}
