package lwcomponent

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	capturer "github.com/lacework/go-sdk/internal/capturer"
	"github.com/lacework/go-sdk/lwlogger"
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
		err = DownloadFile(file.Name(), fmt.Sprintf("%s%s", server.URL, urlPath))
		assert.Nil(t, err)

		buf, err := os.ReadFile(file.Name())
		assert.Nil(t, err)
		assert.Equal(t, content, string(buf))
	})

	t.Run("retry on error", func(t *testing.T) {
		var (
			count int = 0
		)

		mux.HandleFunc("/err", func(w http.ResponseWriter, r *http.Request) {
			if assert.Equal(t, "GET", r.Method, "Get() should be a GET method") {
				count += 1
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		})

		logsCaptured := capturer.CaptureOutput(func() {
			log = lwlogger.New("INFO").Sugar()
			err = DownloadFile(file.Name(), fmt.Sprintf("%s%s", server.URL, "/err"))
		})
		assert.NotNil(t, err)
		assert.Equal(t, DefaultMaxRetry+1, count)

		assert.Contains(t, logsCaptured, "WARN RESTY Get")
		assert.Contains(t, logsCaptured, "/err\": EOF")
		assert.Contains(t, logsCaptured, "Attempt 4") // the fifth attempt will error
		assert.Contains(t, logsCaptured, "ERROR RESTY Get")
		assert.Contains(t, logsCaptured, "Failed to download component")
		assert.Contains(t, logsCaptured, "trace_info")
	})

	t.Run("url error", func(t *testing.T) {
		logsCaptured := capturer.CaptureOutput(func() {
			log = lwlogger.New("INFO").Sugar()
			err = DownloadFile(file.Name(), "")
		})
		assert.NotNil(t, err)
		assert.False(t, os.IsTimeout(err))

		assert.Contains(t, logsCaptured, "WARN RESTY Get")
		assert.Contains(t, logsCaptured, "Attempt 4") // the fifth attempt will error
		assert.Contains(t, logsCaptured, "ERROR RESTY Get")
		assert.Contains(t, logsCaptured, "Failed to download component")
		assert.Contains(t, logsCaptured, "trace_info")
	})
}
