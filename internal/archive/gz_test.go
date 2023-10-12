package archive

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGunzip(t *testing.T) {
	testingDir, err := os.MkdirTemp("", "test-cdk-internal-archive-")
	if err != nil {
		return
	}
	defer os.RemoveAll(testingDir)
	testGZFile := filepath.Join("test_resources", "test_archives.tar.gz")
	testTargetFile := filepath.Join(testingDir, "test_archives.tar")

	err = Gunzip(testGZFile, testTargetFile)
	if err != nil {
		return
	}
	_, err = os.Stat(testTargetFile)
	assert.True(t, !os.IsNotExist(err))
}
