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
	"testing"

	"github.com/lacework/go-sdk/api"
	"github.com/stretchr/testify/assert"
)

var (
	lqlQueryTests = []TestLQLQuery{
		TestLQLQuery{
			Name: "empty-blob",
			Input: &api.LQLQuery{
				QueryBlob: ``,
			},
			Output: "unable to translate query blob",
			Expected: &api.LQLQuery{
				QueryText: ``,
				QueryBlob: ``,
			},
		},
		TestLQLQuery{
			Name: "junk-blob",
			Input: &api.LQLQuery{
				QueryBlob: `this is junk`,
			},
			Output: "unable to translate query blob",
			Expected: &api.LQLQuery{
				QueryText: ``,
				QueryBlob: `this is junk`,
			},
		},
		TestLQLQuery{
			Name: "json-blob",
			Input: &api.LQLQuery{
				QueryBlob: `{"QUERY_TEXT": "my_lql(CloudTrailRawEvents e) { SELECT INSERT_ID LIMIT 10 }"}`,
			},
			Output: nil,
			Expected: &api.LQLQuery{
				QueryText: `my_lql(CloudTrailRawEvents e) { SELECT INSERT_ID LIMIT 10 }`,
				QueryBlob: `{"QUERY_TEXT": "my_lql(CloudTrailRawEvents e) { SELECT INSERT_ID LIMIT 10 }"}`,
			},
		},
		TestLQLQuery{
			Name: "lql-blob",
			Input: &api.LQLQuery{
				QueryBlob: `my_lql(CloudTrailRawEvents e) { SELECT INSERT_ID LIMIT 10 }`,
			},
			Output: nil,
			Expected: &api.LQLQuery{
				QueryText: `my_lql(CloudTrailRawEvents e) { SELECT INSERT_ID LIMIT 10 }`,
				QueryBlob: `my_lql(CloudTrailRawEvents e) { SELECT INSERT_ID LIMIT 10 }`,
			},
		},
		TestLQLQuery{
			Name: "overwrite-blob",
			Input: &api.LQLQuery{
				QueryBlob: `my_lql(CloudTrailRawEvents e) { SELECT INSERT_ID LIMIT 10 }`,
				QueryText: `should not overwrite`,
			},
			Output: nil,
			Expected: &api.LQLQuery{
				QueryText: `should not overwrite`,
				QueryBlob: `my_lql(CloudTrailRawEvents e) { SELECT INSERT_ID LIMIT 10 }`,
			},
		},
	}
)

type TestLQLQuery struct {
	Name     string
	Input    *api.LQLQuery
	Output   interface{}
	Expected *api.LQLQuery
}

func TestLQLQueryTranslate(t *testing.T) {
	for _, lqlQueryTest := range lqlQueryTests {
		t.Run(lqlQueryTest.Name, func(t *testing.T) {
			if err := lqlQueryTest.Input.Translate(); err == nil {
				assert.Equal(t, lqlQueryTest.Output, err)
			} else {
				assert.Equal(t, lqlQueryTest.Output, err.Error())
			}
			assert.Equal(t, lqlQueryTest.Expected, lqlQueryTest.Input)
		})
	}
}
