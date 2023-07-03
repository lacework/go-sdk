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

	"github.com/lacework/go-sdk/api"
	"github.com/stretchr/testify/assert"
)

func TestContainerRegistriesToTable(t *testing.T) {
	cases := []struct {
		Data          []api.ContainerRegistryRaw
		expectedTable [][]string
	}{
		{nil, nil},
		{[]api.ContainerRegistryRaw{
			api.NewContainerRegistry("mock1", api.InlineScannerContainerRegistry, api.InlineScannerData{}),
		}, [][]string{{"", "mock1", "ContVulnCfg", "Enabled", "Ok"}}},
		{[]api.ContainerRegistryRaw{
			api.NewContainerRegistry("mock2", api.ProxyScannerContainerRegistry, api.ProxyScannerData{}),
		}, [][]string{{"", "mock2", "ContVulnCfg", "Enabled", "Ok"}}},
		{[]api.ContainerRegistryRaw{
			api.NewContainerRegistry("mock3", api.GcpGarContainerRegistry, api.GcpGcrData{}),
		}, [][]string{{"", "mock3", "ContVulnCfg", "Enabled", "Pending"}}},
		{[]api.ContainerRegistryRaw{
			api.NewContainerRegistry("mock1", api.InlineScannerContainerRegistry, api.InlineScannerData{}),
			api.NewContainerRegistry("mock2", api.ProxyScannerContainerRegistry, api.ProxyScannerData{}),
			api.NewContainerRegistry("mock3", api.GcpGarContainerRegistry, api.GcpGcrData{}),
		}, [][]string{
			{"", "mock1", "ContVulnCfg", "Enabled", "Ok"},
			{"", "mock2", "ContVulnCfg", "Enabled", "Ok"},
			{"", "mock3", "ContVulnCfg", "Enabled", "Pending"},
		}},
	}

	for i, kase := range cases {
		t.Run(fmt.Sprintf("test case %d", i), func(t *testing.T) {
			subjectTable := containerRegistriesToTable(kase.Data)
			assert.Equal(t, kase.expectedTable, subjectTable)
		})
	}
}
