//
// Author:: Salim Afiune Maya (<afiune@lacework.net>)
// Copyright:: Copyright 2023, Lacework Inc.
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

func TestGenerateContainerVulnListCacheKey(t *testing.T) {
	cases := []struct {
		filterFlagsToHash cacheFiltersToBuildVulnContainerHash
		expectedCacheKey  string
	}{
		{cacheFiltersToBuildVulnContainerHash{
			"", "", "", []string{}, []string{}},
			"vulnerability/container/v2_3285545029616131935"},
		{cacheFiltersToBuildVulnContainerHash{
			"@d", "now", "", []string{}, []string{}},
			"vulnerability/container/v2_8666301743654077811"},
		{cacheFiltersToBuildVulnContainerHash{
			"@d", "now", "", []string{"repo1", "repo2"}, []string{"reg1"}},
			"vulnerability/container/v2_2929007791209551587"},
		{cacheFiltersToBuildVulnContainerHash{
			"", "now", "", []string{}, []string{"reg1"}},
			"vulnerability/container/v2_5320155942991519168"},
		// note, this is just like the first case
		{cacheFiltersToBuildVulnContainerHash{
			"", "", "", []string{}, []string{}},
			"vulnerability/container/v2_3285545029616131935"},
	}

	// first time we test all the the test cases
	for i, kase := range cases {
		t.Run(fmt.Sprintf("first case %d", i), func(t *testing.T) {
			vulCmdState.Start = kase.filterFlagsToHash.Start
			vulCmdState.End = kase.filterFlagsToHash.End
			vulCmdState.Range = kase.filterFlagsToHash.Range
			vulCmdState.Repositories = kase.filterFlagsToHash.Repositories
			vulCmdState.Registries = kase.filterFlagsToHash.Registries

			assert.Equal(t, kase.expectedCacheKey, generateContainerVulnListCacheKey())
		})
	}

	// second time should generate the same hashes
	for i, kase := range cases {
		t.Run(fmt.Sprintf("second case %d", i), func(t *testing.T) {
			vulCmdState.Start = kase.filterFlagsToHash.Start
			vulCmdState.End = kase.filterFlagsToHash.End
			vulCmdState.Range = kase.filterFlagsToHash.Range
			vulCmdState.Repositories = kase.filterFlagsToHash.Repositories
			vulCmdState.Registries = kase.filterFlagsToHash.Registries

			assert.Equal(t, kase.expectedCacheKey, generateContainerVulnListCacheKey())
		})
	}
}
