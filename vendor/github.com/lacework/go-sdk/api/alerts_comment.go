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
	"fmt"

	"github.com/pkg/errors"
)

type AlertsCommentRequest struct {
	Comment string `json:"comment"`
}

type AlertsCommentResponse struct {
	Data AlertTimeline `json:"data"`
}

func (svc *AlertsService) Comment(id int, comment string) (
	response AlertsCommentResponse,
	err error,
) {
	if comment == "" {
		err = errors.New("alert comment must be provided")
		return
	}
	err = svc.client.RequestEncoderDecoder(
		"POST",
		fmt.Sprintf(apiV2AlertsComment, id),
		AlertsCommentRequest{Comment: comment},
		&response,
	)
	return
}
