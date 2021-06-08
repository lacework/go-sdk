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

package domain_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	subject "github.com/lacework/go-sdk/internal/domain"
)

func TestDomains(t *testing.T) {
	cases := []struct {
		URL              string
		expectedAccount  string
		expectedCluster  string
		expectedInternal bool
		expectedString   string
		expectedError    string
	}{
		{URL: "account.lacework.net",
			expectedAccount:  "account",
			expectedCluster:  "",
			expectedString:   "account",
			expectedInternal: false},
		{URL: "account.fra.lacework.net",
			expectedAccount:  "account",
			expectedCluster:  "fra",
			expectedString:   "account.fra",
			expectedInternal: false},
		{URL: "abc.abc.corp.lacework.net",
			expectedAccount:  "abc",
			expectedCluster:  "abc",
			expectedString:   "abc.abc.corp",
			expectedInternal: true},
		{URL: "http://account.lacework.net",
			expectedAccount:  "account",
			expectedCluster:  "",
			expectedString:   "account",
			expectedInternal: false},
		{URL: "https://account.lacework.net",
			expectedAccount:  "account",
			expectedCluster:  "",
			expectedString:   "account",
			expectedInternal: false},
		{URL: "https://account.lacework.net/foo/bar",
			expectedAccount:  "account",
			expectedCluster:  "",
			expectedString:   "account",
			expectedInternal: false},
		{URL: "https://account.fra.lacework.net",
			expectedAccount:  "account",
			expectedCluster:  "fra",
			expectedString:   "account.fra",
			expectedInternal: false},
		{URL: "https://my-super-long-account.devc.corp.lacework.net/bubulubu",
			expectedAccount:  "my-super-long-account",
			expectedCluster:  "devc",
			expectedString:   "my-super-long-account.devc.corp",
			expectedInternal: true},

		// Errors!!!!
		{URL: "",
			expectedError: "domain not supported"},
		{URL: "account.lacework.com",
			expectedError: "domain not supported"},
		{URL: "account.c.not-corp.lacework.net",
			expectedError: "unable to detect if domain is internal"},
		{URL: "too.many.sub.domains.corp.lacework.net",
			expectedError: "unable to detect domain information"},
	}
	for i, kase := range cases {
		t.Run(fmt.Sprintf("test case %d", i), func(t *testing.T) {
			d, err := subject.New(kase.URL)
			if err != nil {
				assert.Equal(t, kase.expectedError, err.Error())
			}
			assert.Equal(t, kase.expectedAccount, d.Account)
			assert.Equal(t, kase.expectedCluster, d.Cluster)
			assert.Equal(t, kase.expectedInternal, d.Internal)
			assert.Equal(t, kase.expectedString, d.String())
		})
	}
}
