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

import (
	"net/url"

	"github.com/pkg/errors"
)

const (
	LQLDeleteBadInputMsg string = "Please specify a valid query ID."
)

type LQLDeleteResponse struct {
	Ok      bool             `json:"ok"`
	Message LQLDeleteMessage `json:"message"`
}

type LQLDeleteMessage struct {
	ID string `json:"lqlDeleted"`
}

func (svc *LQLService) DeleteQuery(queryID string) (
	response LQLDeleteResponse,
	err error,
) {
	if queryID == "" {
		err = errors.New(LQLDeleteBadInputMsg)
		return
	}
	var uri string = ApiLQL + "?LQL_ID=" + url.QueryEscape(queryID)

	err = svc.client.RequestDecoder("DELETE", uri, nil, &response)
	return
}
