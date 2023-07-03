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

package api_test

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	subject "github.com/lacework/go-sdk/api"
)

func TestVersionMatchVERSIONfile(t *testing.T) {
	expectedVersion, err := os.ReadFile("../VERSION")
	assert.Nil(t, err)
	assert.Equalf(t, strings.TrimSuffix(string(expectedVersion), "\n"), subject.Version,
		"api/version.go doesn't match with VERSION file; run scripts/version_updater.sh")
}
