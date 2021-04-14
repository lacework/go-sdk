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

	"github.com/lacework/go-sdk/api"
)

func TestSplitAzureSubscriptionsApiResponse(t *testing.T) {
	cases := []struct {
		subject  api.CompAzureSubscriptions
		expected cliComplianceAzureInfo
	}{
		// empty subscriptions will return empty cli info
		{
			api.CompAzureSubscriptions{},
			cliComplianceAzureInfo{Subscriptions: make([]cliComplianceIDAlias, 0)},
		},
		// real test case with NO alias
		{
			api.CompAzureSubscriptions{
				Tenant:        "ABCCC123-abc-123-AB12-XYZ987",
				Subscriptions: []string{"subscription-id-1", "subscription-id-2", "subscription-id-3", "subscription-id-4"},
			},
			cliComplianceAzureInfo{
				Tenant: cliComplianceIDAlias{"ABCCC123-abc-123-AB12-XYZ987", ""},
				Subscriptions: []cliComplianceIDAlias{
					cliComplianceIDAlias{"subscription-id-1", ""},
					cliComplianceIDAlias{"subscription-id-2", ""},
					cliComplianceIDAlias{"subscription-id-3", ""},
					cliComplianceIDAlias{"subscription-id-4", ""},
				},
			},
		},
		// real test case with alias
		{
			api.CompAzureSubscriptions{
				Tenant: "ABCCC123-abc-123-AB12-XYZ987 (cool.org.alias.example.com)",
				Subscriptions: []string{
					"id-1 (a test subscription)",
					"xmen-subscription (serious alias)",
					"disney-movies (Maybe Production)",
					"foo (bar)",
				},
			},
			cliComplianceAzureInfo{
				Tenant: cliComplianceIDAlias{"ABCCC123-abc-123-AB12-XYZ987", "cool.org.alias.example.com"},
				Subscriptions: []cliComplianceIDAlias{
					cliComplianceIDAlias{"id-1", "a test subscription"},
					cliComplianceIDAlias{"xmen-subscription", "serious alias"},
					cliComplianceIDAlias{"disney-movies", "Maybe Production"},
					cliComplianceIDAlias{"foo", "bar"},
				},
			},
		},
	}
	for i, kase := range cases {
		t.Run(fmt.Sprintf("test case %d", i), func(t *testing.T) {
			assert.Equalf(t,
				kase.expected, splitAzureSubscriptionsApiResponse(kase.subject),
				"there is a problem with this test case, please check",
			)
		})
	}
}
