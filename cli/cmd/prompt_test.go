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

package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatSecret(t *testing.T) {
	assert.Equal(t,
		formatSecret(4, "_ab4c34d2df97babcd"),
		"**************abcd",
		"secrets are not being formatted correctly")

	assert.Equal(t, formatSecret(0, "_ab4c34d2df97babcd"), "******************")
	assert.Equal(t, formatSecret(1, "_ab4c34d2df97babcd"), "*****************d")
	assert.Equal(t, formatSecret(2, "_ab4c34d2df97babcd"), "****************cd")
	assert.Equal(t, formatSecret(3, "_ab4c34d2df97babcd"), "***************bcd")
	assert.Equal(t, formatSecret(4, "_ab4c34d2df97babcd"), "**************abcd")
	assert.Equal(t, formatSecret(5, "_ab4c34d2df97babcd"), "*************babcd")
	assert.Equal(t, formatSecret(6, "_ab4c34d2df97babcd"), "************7babcd")
	assert.Equal(t, formatSecret(7, "_ab4c34d2df97babcd"), "***********97babcd")
	assert.Equal(t, formatSecret(8, "_ab4c34d2df97babcd"), "**********f97babcd")
	assert.Equal(t, formatSecret(9, "_ab4c34d2df97babcd"), "*********df97babcd")
	assert.Equal(t, formatSecret(10, "_ab4c34d2df97babcd"), "********2df97babcd")
	assert.Equal(t, formatSecret(11, "_ab4c34d2df97babcd"), "*******d2df97babcd")
	assert.Equal(t, formatSecret(12, "_ab4c34d2df97babcd"), "******4d2df97babcd")
	assert.Equal(t, formatSecret(13, "_ab4c34d2df97babcd"), "*****34d2df97babcd")
	assert.Equal(t, formatSecret(14, "_ab4c34d2df97babcd"), "****c34d2df97babcd")
	assert.Equal(t, formatSecret(15, "_ab4c34d2df97babcd"), "***4c34d2df97babcd")
	assert.Equal(t, formatSecret(16, "_ab4c34d2df97babcd"), "**b4c34d2df97babcd")
	assert.Equal(t, formatSecret(17, "_ab4c34d2df97babcd"), "*ab4c34d2df97babcd")
	assert.Equal(t, formatSecret(18, "_ab4c34d2df97babcd"), "_ab4c34d2df97babcd")
	assert.Equal(t, formatSecret(20, "_ab4c34d2df97babcd"), "_ab4c34d2df97babcd")

	// empty string
	assert.Equal(t, formatSecret(0, ""), "")
	assert.Equal(t, formatSecret(10, ""), "")
}
