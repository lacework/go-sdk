package lwcomponent_test

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/lacework/go-sdk/lwcomponent"
	"github.com/stretchr/testify/assert"
)

func TestStageTarGzCommit(t *testing.T) {
	var (
		expectedSig string = "expected sig"
		name               = "test-commit"
	)

	t.Run("ok", func(t *testing.T) {
		stage, err := lwcomponent.NewStageTarGz(name, "", 0)
		assert.Nil(t, err)
		defer stage.Close()

		target, err := os.MkdirTemp("", "apiInfo-Ok")
		if err != nil {
			panic(err)
		}
		defer os.RemoveAll(target)

		fileName := filepath.Join(stage.Directory(), lwcomponent.SignatureFile)

		os.WriteFile(fileName, []byte(expectedSig), os.ModePerm)

		os.Mkdir(filepath.Join(stage.Directory(), "nested"), os.ModePerm)
		os.WriteFile(filepath.Join(stage.Directory(), "nested", lwcomponent.SignatureFile), []byte(expectedSig), os.ModePerm)

		err = stage.Commit(target)
		assert.Nil(t, err)

		data, err := os.ReadFile(filepath.Join(target, lwcomponent.SignatureFile))
		assert.NotNil(t, data)
		assert.Nil(t, err)
		assert.Equal(t, expectedSig, string(data))

		data, err = os.ReadFile(filepath.Join(target, "nested", lwcomponent.SignatureFile))
		assert.NotNil(t, data)
		assert.Nil(t, err)
		assert.Equal(t, expectedSig, string(data))
	})

	t.Run("target doesn't exist", func(t *testing.T) {
		stage, err := lwcomponent.NewStageTarGz(name, "", 0)
		assert.Nil(t, err)
		defer stage.Close()

		err = stage.Commit("")
		assert.NotNil(t, err)
	})
}

func TestStagingTarGzSignature(t *testing.T) {
	var (
		expectedSig string = "test signature"
		name               = "test-sig"
	)

	t.Run("ok", func(t *testing.T) {
		stage, err := lwcomponent.NewStageTarGz(name, "", 0)
		assert.Nil(t, err)
		defer stage.Close()

		fileName := filepath.Join(stage.Directory(), lwcomponent.SignatureFile)

		os.WriteFile(fileName, []byte(expectedSig), os.ModePerm)

		result, err := stage.Signature()
		assert.Nil(t, err)
		assert.Equal(t, expectedSig, string(result))
	})

	t.Run("no signature file", func(t *testing.T) {
		stage, err := lwcomponent.NewStageTarGz(name, "", 0)
		assert.Nil(t, err)
		defer stage.Close()

		_, err = stage.Signature()
		assert.NotNil(t, err)
	})
}

func TestStagingTarGzUnpack(t *testing.T) {
	var (
		name          = "stageTarGzUnpack"
		componentData = "component data"
		sigData       = "signature data"
	)

	stage, err := lwcomponent.NewStageTarGz(name, fmt.Sprintf("https://127.0.0.1/%s.tar.gz", name), 0)
	assert.Nil(t, err)
	assert.NotNil(t, stage)
	defer stage.Close()

	makeGzip(name, makeTar(name, "1.1.1", stage.Directory(), componentData, sigData))

	stage.Unpack()
}

func makeTar(name string, version string, dir string, data string, sig string) string {
	tarname := fmt.Sprintf("%s.tar", name)
	path := filepath.Join(dir, tarname)

	tarfile, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer tarfile.Close()

	tarball := tar.NewWriter(tarfile)
	defer tarball.Close()

	files := []struct {
		Name, Body string
	}{
		{name, data},
		{lwcomponent.VersionFile, version},
		{lwcomponent.SignatureFile, sig},
	}
	for _, file := range files {
		hdr := &tar.Header{
			Name: file.Name,
			Mode: 0600,
			Size: int64(len(file.Body)),
		}

		if err := tarball.WriteHeader(hdr); err != nil {
			panic(err)
		}

		if _, err := tarball.Write([]byte(file.Body)); err != nil {
			panic(err)
		}
	}

	return path
}

func makeGzip(name, path string) string {
	reader, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	target := fmt.Sprintf("%s.gz", path)

	writer, err := os.Create(target)
	if err != nil {
		panic(err)
	}
	defer writer.Close()

	archiver := gzip.NewWriter(writer)
	archiver.Name = name
	defer archiver.Close()

	_, err = io.Copy(archiver, reader)
	if err != nil {
		panic(err)
	}

	return target
}
