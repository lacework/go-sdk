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

package lacework_test

import (
	"io/ioutil"
	"os"
	"testing"

	lwint "github.com/lacework/go-sdk/internal/lacework"
	"github.com/stretchr/testify/assert"
)

func TestFileExistsWhenFileActuallyExists(t *testing.T) {
	file, err := ioutil.TempFile("", "bar")
	if assert.Nil(t, err) {
		assert.True(t, lwint.FileExists(file.Name()))
		os.Remove(file.Name())
	}
}

func TestFileExistsWhenFileIsADirectory(t *testing.T) {
	dir, err := ioutil.TempDir("", "bar")
	if assert.Nil(t, err) {
		assert.False(t, lwint.FileExists(dir))
		os.RemoveAll(dir)
	}
}

func TestFileExistsWhenFileDoesNotExists(t *testing.T) {
	assert.False(t, lwint.FileExists("file.name"))
}
