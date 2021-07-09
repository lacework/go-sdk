//
// Author:: Salim Afiune Maya (<afiune@lacework.net>)
// Copyright:: Copyright 2020-2021, Lacework Inc.
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

package format_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lacework/go-sdk/internal/format"
)

func TestSecret(t *testing.T) {
	assert.Equal(t,
		format.Secret(4, "_ab4c34d2df97babcd"),
		"**************abcd",
		"secrets are not being formatted correctly")

	assert.Equal(t, format.Secret(0, "_ab4c34d2df97babcd"), "******************")
	assert.Equal(t, format.Secret(1, "_ab4c34d2df97babcd"), "*****************d")
	assert.Equal(t, format.Secret(2, "_ab4c34d2df97babcd"), "****************cd")
	assert.Equal(t, format.Secret(3, "_ab4c34d2df97babcd"), "***************bcd")
	assert.Equal(t, format.Secret(4, "_ab4c34d2df97babcd"), "**************abcd")
	assert.Equal(t, format.Secret(5, "_ab4c34d2df97babcd"), "*************babcd")
	assert.Equal(t, format.Secret(6, "_ab4c34d2df97babcd"), "************7babcd")
	assert.Equal(t, format.Secret(7, "_ab4c34d2df97babcd"), "***********97babcd")
	assert.Equal(t, format.Secret(8, "_ab4c34d2df97babcd"), "**********f97babcd")
	assert.Equal(t, format.Secret(9, "_ab4c34d2df97babcd"), "*********df97babcd")
	assert.Equal(t, format.Secret(10, "_ab4c34d2df97babcd"), "********2df97babcd")
	assert.Equal(t, format.Secret(11, "_ab4c34d2df97babcd"), "*******d2df97babcd")
	assert.Equal(t, format.Secret(12, "_ab4c34d2df97babcd"), "******4d2df97babcd")
	assert.Equal(t, format.Secret(13, "_ab4c34d2df97babcd"), "*****34d2df97babcd")
	assert.Equal(t, format.Secret(14, "_ab4c34d2df97babcd"), "****c34d2df97babcd")
	assert.Equal(t, format.Secret(15, "_ab4c34d2df97babcd"), "***4c34d2df97babcd")
	assert.Equal(t, format.Secret(16, "_ab4c34d2df97babcd"), "**b4c34d2df97babcd")
	assert.Equal(t, format.Secret(17, "_ab4c34d2df97babcd"), "*ab4c34d2df97babcd")
	assert.Equal(t, format.Secret(18, "_ab4c34d2df97babcd"), "_ab4c34d2df97babcd")
	assert.Equal(t, format.Secret(20, "_ab4c34d2df97babcd"), "_ab4c34d2df97babcd")

	// empty string
	assert.Equal(t, format.Secret(0, ""), "")
	assert.Equal(t, format.Secret(10, ""), "")
}
