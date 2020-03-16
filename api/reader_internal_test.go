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

package api

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	Foo string `json:"foo"`
	Bar int    `json:"bar"`
}

func TestJsonReader(t *testing.T) {
	var subject = testStruct{"foo", 1}

	reader, err := jsonReader(subject)
	if assert.Nil(t, err) {
		buf := new(bytes.Buffer)
		buf.ReadFrom(reader)
		assert.Equal(t,
			"{\"foo\":\"foo\",\"bar\":1}\n",
			buf.String(),
			"unexpected streaming encoder")
	}
}
