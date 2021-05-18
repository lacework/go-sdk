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

package cmd

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAPICleanupEndpoint(t *testing.T) {

	cases := []struct {
		endpoint         string
		expectedEndpoint string
	}{
		{"", ""},
		{"/foo", "foo"},
		{"/external/bar", "external/bar"},
		{"/v1/bubulubu", "v1/bubulubu"},
		{"/v2/bubulubu", "v2/bubulubu"},

		{"/api/v2/schemas", "v2/schemas"},
		{"/api/v1/external/foo", "v1/external/foo"},
		{"api/v1/endpoint", "v1/endpoint"},
		{"api/v2/coolendpoint", "v2/coolendpoint"},
	}

	for i, kase := range cases {
		t.Run(fmt.Sprintf("test case %d", i), func(t *testing.T) {
			assert.Equal(t, kase.expectedEndpoint, cleanupEndpoint(kase.endpoint))
		})
	}
}
