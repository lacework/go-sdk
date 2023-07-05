package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrepareUpload(t *testing.T) {
	dir, err := os.MkdirTemp("", "example")
	assert.Nil(t, err)
	defer os.RemoveAll(dir)
	file := filepath.Join(dir, "upload.json")
	err = os.WriteFile(file, []byte("{\"somekey\": \"somedata\"}"), 0666)
	assert.Nil(t, err)
	request, err := prepareUpload("somename", file)
	assert.Nil(t, err)
	assert.Equal(t, "somename", request.Feature)
	assert.Equal(t, "somedata", request.FeatureData["somekey"])
}
