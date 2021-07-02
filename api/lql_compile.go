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

type LQLCompileResponse struct {
	Data    []map[string]interface{} `json:"data"`
	Ok      bool                     `json:"ok"`
	Message string                   `json:"message"`
}

func (svc *LQLService) Compile(q string) (
	response LQLCompileResponse,
	err error,
) {
	var v2Query Query
	if v2Query, err = ParseQuery(q); err != nil {
		return
	}
	v1Query := map[string]string{"query_text": v2Query.QueryText}
	err = svc.client.RequestEncoderDecoder("POST", apiLQLCompile, v1Query, &response)
	return
}
