package archive

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDetectTGZAndUnpackWhenFileIsTGZ(t *testing.T) {
	testingDir, err := os.MkdirTemp("", "test-cdk-internal-archive-")
	if err != nil {
		return
	}
	defer os.RemoveAll(testingDir)
	testFilePath := filepath.Join("test_resources", "test_archive.tar.gz")
	err = DetectTGZAndUnpack(testFilePath, testingDir)
	if err != nil {
		return
	}

	extractPath := filepath.Join(testingDir, "test_archive")
	file, err := os.Open(extractPath)
	if err != nil {
		return
	}
	fileInfo, err := file.Stat()
	if err != nil {
		return
	}
	assert.True(t, fileInfo.IsDir())
}

func TestDetectTGZAndUnpackWhenFileIsNotTGZ(t *testing.T) {
	testingDir, err := os.MkdirTemp("", "test-cdk-internal-archive-")
	if err != nil {
		return
	}
	defer os.RemoveAll(testingDir)
	testFilePath := filepath.Join(testingDir, "test_file")
	contents := []byte("testing\n")
	err = os.WriteFile(testFilePath, contents, 0644)
	if err != nil {
		return
	}

	err = DetectTGZAndUnpack(testFilePath, testingDir)
	if err != nil {
		return
	}

	file, err := os.Open(testFilePath)
	if err != nil {
		return
	}
	fileInfo, err := file.Stat()
	if err != nil {
		return
	}
	assert.False(t, fileInfo.IsDir())
}

func TestFileIsGZWhenFileIsGZ(t *testing.T) {
	testFilePath := filepath.Join("test_resources", "test_archive.tar.gz")
	isGZ, err := FileIsGZ(testFilePath)
	if err != nil {
		return
	}
	assert.True(t, isGZ)
}

func TestFileIsGZWhenFileIsNotGZ(t *testing.T) {
	testingDir, err := os.MkdirTemp("", "test-cdk-internal-archive-")
	if err != nil {
		return
	}
	defer os.RemoveAll(testingDir)
	testFilePath := filepath.Join(testingDir, "test_file")
	contents := []byte("testing\n")
	err = os.WriteFile(testFilePath, contents, 0644)
	if err != nil {
		return
	}
	isGZ, err := FileIsGZ(testFilePath)
	if err != nil {
		return
	}
	assert.False(t, isGZ)
}

func TestFileIsTarWhenFileIsTar(t *testing.T) {
	testFilePath := filepath.Join("test_resources", "test_archive.tar")
	isTar, err := FileIsTar(testFilePath)
	if err != nil {
		return
	}
	assert.True(t, isTar)
}

func TestFileIsTarWhenFileIsNotTar(t *testing.T) {
	testingDir, err := os.MkdirTemp("", "test-cdk-internal-archive-")
	if err != nil {
		return
	}
	defer os.RemoveAll(testingDir)
	testFilePath := filepath.Join(testingDir, "test_file")
	contents := []byte("testing\n")
	err = os.WriteFile(testFilePath, contents, 0644)
	if err != nil {
		return
	}
	isTar, err := FileIsTar(testFilePath)
	if err != nil {
		return
	}
	assert.False(t, isTar)
}
