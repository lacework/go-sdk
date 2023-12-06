package cmd

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteHCLOutputLocation(t *testing.T) {
	t.Run("should write output with existing directory", func(*testing.T) {
		d, err := os.MkdirTemp("", "locationTest")
		if err != nil {
			panic(err)
		}
		defer os.RemoveAll(d)

		_, err = writeHclOutput("", d, "")
		assert.NoError(t, err)
	})
	t.Run("should create missing location directory", func(*testing.T) {
		d, err := os.MkdirTemp("", "locationTest")
		if err != nil {
			panic(err)
		}
		defer os.RemoveAll(d)

		location := fmt.Sprintf("%s/newplace", d)
		_, err = writeHclOutput("", location, "")
		assert.NoError(t, err)

		statOut, err := os.Stat(location)
		assert.NoError(t, err)
		assert.True(t, statOut.IsDir())
	})
	t.Run("should fail on existing location of type file", func(*testing.T) {
		d, err := os.MkdirTemp("", "locationTest")
		if err != nil {
			panic(err)
		}
		defer os.RemoveAll(d)

		fileName := fmt.Sprintf("%s/testfile", d)
		if err := os.WriteFile(fileName, []byte("test"), os.FileMode(0744)); err != nil {
			panic(err)
		}
		_, err = writeHclOutput("", fileName, "")
		assert.Error(t, err)
	})
	t.Run("should write to homedir location when not supplied", func(*testing.T) {
		d, err := os.MkdirTemp("", "locationTest")
		if err != nil {
			panic(err)
		}
		h := os.Getenv("HOME")
		os.Setenv("HOME", d)
		defer func() {
			os.Setenv("HOME", h)
			os.RemoveAll(d)
		}()

		_, err = writeHclOutput("", "", "")
		assert.NoError(t, err)
		statOut, err := os.Stat(fmt.Sprintf("%s/lacework/main.tf", d))
		assert.NoError(t, err)
		assert.True(t, !statOut.IsDir())
	})
}

func TestValidateOutputLocation(t *testing.T) {
	t.Run("should validate existing dir location", func(t *testing.T) {
		d, err := os.MkdirTemp("", "locationTest")
		if err != nil {
			panic(err)
		}
		defer os.RemoveAll(d)
		assert.NoError(t, validateOutputLocation(d))
	})
	t.Run("should validate non-existent location", func(t *testing.T) {
		assert.NoError(t, validateOutputLocation("/i/dont/exist"))
	})
	t.Run("should not validate existing file location", func(t *testing.T) {
		d, err := os.MkdirTemp("", "locationTest")
		if err != nil {
			panic(err)
		}
		defer os.RemoveAll(d)

		fileName := fmt.Sprintf("%s/testfile", d)
		if err := os.WriteFile(fileName, []byte("test"), os.FileMode(0744)); err != nil {
			panic(err)
		}
		assert.Error(t, validateOutputLocation(fileName))
	})
}
