package archive

import (
	"os"
	"path/filepath"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestUnTar(t *testing.T) {
	testingDir, err := os.MkdirTemp("", "test-cdk-internal-archive-")
	if err != nil {
		return
	}
	defer os.RemoveAll(testingDir)
	testFilePath := filepath.Join("test_resources", "test_archives.tar")
	err = UnTar(testFilePath, testingDir)
	if err != nil {
		return
	}
	regFile := filepath.Join(testingDir, "test_archive", "test_file.txt")
	dirFile := filepath.Join(testingDir, "test_archive", "test_folder")
	symlinkFile := filepath.Join(testingDir, "test_archive", "test_soft_line")

	f, err := os.Stat(regFile)
	isRegFile := !os.IsNotExist(err) && !f.IsDir()

	f, err = os.Stat(dirFile)
	isDirFile := !os.IsNotExist(err) && f.IsDir()

	f, err = os.Lstat(symlinkFile)
	isSymlinkFile := !os.IsNotExist(err) && (f.Mode()&os.ModeSymlink == os.ModeSymlink)

	assert.True(t, isRegFile && isDirFile && isSymlinkFile)
}