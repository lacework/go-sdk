package lwcomponent_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/lacework/go-sdk/lwcomponent"
	"github.com/stretchr/testify/assert"
)

func TestDownloadFile(t *testing.T) {
	var (
		urlPath string = "/lw-cdk-store/catalog/component-example/1.0.0/component-example-linux-amd64.tar.gz"
		content string = "CDK component"
	)

	file, err := os.CreateTemp("", "lwcomponent-downloadFile")
	assert.Nil(t, err)
	defer file.Close()

	mux := http.NewServeMux()

	server := httptest.NewServer(mux)
	defer server.Close()

	mux.HandleFunc(urlPath, func(w http.ResponseWriter, r *http.Request) {
		if assert.Equal(t, "GET", r.Method, "Get() should be a GET method") {
			fmt.Fprint(w, content)
		}
	})

	t.Run("happy path", func(t *testing.T) {
		err = lwcomponent.DownloadFile(file.Name(), fmt.Sprintf("%s%s", server.URL, urlPath), 0)
		assert.Nil(t, err)

		buf, err := os.ReadFile(file.Name())
		assert.Nil(t, err)
		assert.Equal(t, content, string(buf))
	})

	t.Run("timeout error", func(t *testing.T) {
		err = lwcomponent.DownloadFile(file.Name(), fmt.Sprintf("%s%s", server.URL, urlPath), 1*time.Microsecond)
		assert.NotNil(t, err)
		assert.True(t, os.IsTimeout(err))
	})

	t.Run("non-timeout error", func(t *testing.T) {
		err = lwcomponent.DownloadFile(file.Name(), "", 0)
		assert.NotNil(t, err)
		assert.False(t, os.IsTimeout(err))
	})
}
