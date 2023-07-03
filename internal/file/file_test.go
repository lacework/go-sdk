//
// Author:: Salim Afiune Maya (<afiune@lacework.net>)
// Copyright:: Copyright 2021, Lacework Inc.
// License:: Apache License, Version 2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package file_test

import (
	"os"
	"path"
	"testing"

	"github.com/lacework/go-sdk/internal/file"
	"github.com/stretchr/testify/assert"
)

func TestFileExistsWhenFileActuallyExists(t *testing.T) {
	f, err := os.CreateTemp("", "bar")
	defer os.Remove(f.Name())
	if assert.Nil(t, err) {
		assert.True(t, file.FileExists(f.Name()))
	}
}

func TestFileExistsWhenFileIsADirectory(t *testing.T) {
	dir, err := os.MkdirTemp("", "t")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)
	assert.False(t, file.FileExists(dir))
}

func TestFileExistsButWeDontHavePermissions(t *testing.T) {
	// @afiune don't run it in Codefresh since the pipeline runs as root
	// and root has access to all files and directories
	if os.Getenv("CI") != "" {
		t.Skip("skipping since CI runs as root, and root has permissions")
		return
	}

	dir, err := os.MkdirTemp("", "t")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)

	// create a directory that we can't read
	os.Mkdir(path.Join(dir, "protected"), 0700)
	os.WriteFile(path.Join(dir, "protected", "bubulubu"), []byte("data"), 0644)
	os.Chmod(path.Join(dir, "protected"), 0000)

	assert.False(t, file.FileExists(path.Join(dir, "protected", "bubulubu")))
}

func TestFileExistsWhenFileDoesNotExists(t *testing.T) {
	assert.False(t, file.FileExists("file.name"))
}
