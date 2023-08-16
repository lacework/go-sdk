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
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// errorResponse handles errors caused by a Lacework API request
type errorResponse struct {
	Response *http.Response
	Message  string
}

type apiErrorResponse struct {
	Ok        bool   `json:"ok"`
	V2Message string `json:"message"`
	Data      struct {
		Message       string `json:"message"`
		StatusMessage string `json:"statusMessage"`
		ErrorMsg      string `json:"ErrorMsg"`
	} `json:"data"`
}

// Message extracts the message from an api error response
func (r *apiErrorResponse) Message() string {
	if r != nil {
		if r.Data.ErrorMsg != "" {
			return r.Data.ErrorMsg
		}
		if r.Data.Message != "" {
			return r.Data.Message
		}
		if r.V2Message != "" && r.V2Message != "SUCCESS" {
			return r.V2Message
		}
		if r.Data.StatusMessage != "" {
			return r.Data.StatusMessage
		}
	}
	return ""
}

// Error fulfills the built-in error interface function
func (r *errorResponse) Error() string {
	return fmt.Sprintf("\n  [%v] %v\n  [%d] %s",
		r.Response.Request.Method,
		r.Response.Request.URL,
		r.Response.StatusCode,
		r.Message,
	)
}

// checkResponse checks the provided response and generates an Error
func checkErrorInResponse(r *http.Response) error {
	if c := r.StatusCode; c >= 200 && c <= 299 {
		return nil
	}

	var (
		errRes    = &errorResponse{Response: r}
		data, err = io.ReadAll(r.Body)
	)
	if err == nil && len(data) > 0 {
		// try to unmarshal the api error response
		apiErrRes := &apiErrorResponse{}
		if err := json.Unmarshal(data, apiErrRes); err != nil {
			errRes.Message = string(data)
			return errRes
		}

		var (
			apiErrResMessage = apiErrRes.Message()
			statusText       = http.StatusText(r.StatusCode)
		)

		// try our best to parse the error message
		if apiErrResMessage != "" {
			errRes.Message = apiErrResMessage
			return errRes
		}

		// if it is empty, try to display the Status Code in pretty Text
		if statusText != "" {
			errRes.Message = statusText
			return errRes
		}

		// if we couldn't even decode the StatusCode... well, lets just be transparent
		errRes.Message = "Unknown"
	}

	return errRes
}
