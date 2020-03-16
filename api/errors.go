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
	"io/ioutil"
	"net/http"
)

// erorResponse handles errors caused by a Lacework API request
type errorResponse struct {
	Response *http.Response
	Message  string
}

type apiErrorResponse struct {
	Ok   bool
	Data struct {
		Message       string
		StatusMessage string
	}
}

// Message extracts the message from an api error response
func (r *apiErrorResponse) Message() string {
	if r != nil {
		return r.Data.Message
	}
	return ""
}

// Error fulfills the built-in error interface function
func (r *errorResponse) Error() string {
	return fmt.Sprintf("[%v] %v: %d %s",
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
		data, err = ioutil.ReadAll(r.Body)
	)
	if err == nil && len(data) > 0 {
		// try to unmarshal the api error response
		apiErrRes := &apiErrorResponse{}
		if err := json.Unmarshal(data, apiErrRes); err == nil {
			errRes.Message = apiErrRes.Message()
		} else {
			errRes.Message = string(data)
		}
	}

	return errRes
}
