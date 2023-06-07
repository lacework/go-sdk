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

type alertCommentFormat int

const (
	AlertCommentFormatPlaintext alertCommentFormat = iota
	AlertCommentFormatMarkdown
)

// String returns the string representation of an Alert comment format
func (i alertCommentFormat) String() string {
	return AlertCommentFormats[i]
}

type alertCommentFormats map[alertCommentFormat]string

// AlertCommentFormats is the list of available Alert comment formats
var AlertCommentFormats = alertCommentFormats{
	AlertCommentFormatPlaintext: "plaintext",
	AlertCommentFormatMarkdown:  "markdown",
}

func (acf alertCommentFormats) GetOrderedFormatStrings() []string {
	formats := []string{}
	for i := 0; i < len(acf); i++ {
		formats = append(formats, acf[alertCommentFormat(i)])
	}
	return formats
}

type AlertsCommentRequest struct {
	Comment string `json:"comment"`
	Format  string `json:"format"`
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

func (svc *AlertsService) CommentFromRequest(id int, request AlertsCommentRequest) (
	response AlertsCommentResponse,
	err error,
) {
	if request.Comment == "" {
		err = errors.New("alert comment must be provided")
		return
	}
	err = svc.client.RequestEncoderDecoder(
		"POST",
		fmt.Sprintf(apiV2AlertsComment, id),
		request,
		&response,
	)
	return
}
