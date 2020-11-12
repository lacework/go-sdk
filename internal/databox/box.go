//go:generate go run generator/main.go
//
// Author:: Salim Afiune Maya (<afiune@lacework.net>)
// Copyright:: Copyright 2020, Lacework Inc.
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

package databox

import (
	"fmt"
	"strings"
)

// Embed a global data box
var box = &data{box: make(map[string][]byte)}

// Add a file to the global box
func Add(file string, content []byte) {
	box.Add(file, content)
}

// Get a file from the global box
func Get(file string) ([]byte, bool) {
	return box.Get(file)
}

// Data box definition
type data struct {
	box map[string][]byte
}

// Add a file to the box
func (d *data) Add(file string, content []byte) {
	d.box[file] = content
}

// Get a file from the box
func (d *data) Get(file string) ([]byte, bool) {
	if !strings.HasPrefix(file, "/") {
		file = fmt.Sprintf("/%s", file)
	}

	f, ok := d.box[file]
	return f, ok
}
