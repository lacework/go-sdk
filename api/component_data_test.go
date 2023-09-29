package api_test

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/lacework"
	"github.com/stretchr/testify/assert"
)

func TestUploadFiles(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()
	fakeServer.MockAPI("ComponentData/requestUpload", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.NotNil(t, r.Body)
		body := httpBodySniffer(r)
		assert.Contains(t, body, "sast")
		assert.Contains(t, body, "doc-set")
		_, err := fmt.Fprintf(w, generateInitialResponse())
		assert.Nil(t, err)
	})
	fakeServer.MockAPI("ComponentData/completeUpload", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.NotNil(t, r.Body)
		body := httpBodySniffer(r)
		assert.Contains(t, body, "SOME-GUID")
		_, err := fmt.Fprintf(w, generateCompleteResponse())
		assert.Nil(t, err)
	})
	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)
	guid, err := c.V2.ComponentData.UploadFiles("doc-set", []string{"sast"}, []string{})
	assert.Nil(t, err)
	assert.Equal(t, "SOME-GUID", guid)
}

func TestDefaultUrlType(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()
	fakeServer.MockAPI("ComponentData/requestUpload", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.NotNil(t, r.Body)
		body := httpBodySniffer(r)
		assert.Contains(t, body, api.URL_TYPE_DEFAULT)
		_, err := fmt.Fprintf(w, generateInitialResponse())
		assert.Nil(t, err)
	})
	fakeServer.MockAPI("ComponentData/completeUpload", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.NotNil(t, r.Body)
		body := httpBodySniffer(r)
		assert.Contains(t, body, api.URL_TYPE_DEFAULT)
		_, err := fmt.Fprintf(w, generateCompleteResponse())
		assert.Nil(t, err)
	})
	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)
	guid, err := c.V2.ComponentData.UploadFiles("doc-set", []string{"sast"}, []string{})
	assert.Nil(t, err)
	assert.Equal(t, "SOME-GUID", guid)
}

func TestSastTablesUrlType(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()
	fakeServer.MockAPI("ComponentData/requestUpload", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.NotNil(t, r.Body)
		body := httpBodySniffer(r)
		assert.Contains(t, body, api.URL_TYPE_SAST_TABLES)
		_, err := fmt.Fprintf(w, generateInitialResponse())
		assert.Nil(t, err)
	})
	fakeServer.MockAPI("ComponentData/completeUpload", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.NotNil(t, r.Body)
		body := httpBodySniffer(r)
		assert.Contains(t, body, api.URL_TYPE_SAST_TABLES)
		_, err := fmt.Fprintf(w, generateCompleteResponse())
		assert.Nil(t, err)
	})
	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)
	guid, err := c.V2.ComponentData.UploadSastTables("doc-set", []string{})
	assert.Nil(t, err)
	assert.Equal(t, "SOME-GUID", guid)
}

func TestDoWithExponentialBackoffAlwaysFailing(t *testing.T) {
	waited := 0
	err := api.DoWithExponentialBackoff(func() error {
		return errors.New("failed")
	}, func(x int) {
		waited += x
	})
	assert.NotNil(t, err)
	assert.Equal(t, 62, waited)
}

func TestDoWithExponentialBackoffSucceedsOnThirdAttempt(t *testing.T) {
	waited := 0
	attempt := 1
	err := api.DoWithExponentialBackoff(func() error {
		if attempt == 3 {
			return nil
		}
		attempt += 1
		return errors.New("failed")
	}, func(x int) {
		waited += x
	})
	assert.Nil(t, err)
	assert.Equal(t, 6, waited)
}

func generateInitialResponse() string {
	return `
{
	"data": {
		"guid": "SOME-GUID",
		"uploadMethods": [
			{
				"method": "AwsS3",
				"info": {}
			}
		]
	}
}
`
}

func generateCompleteResponse() string {
	return `
{
	"data": {
		"guid": "SOME-GUID"
	}
}
`
}
