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
	"bytes"

	"go.uber.org/zap"
)

// LQLService is a service that interacts with the LQL
// endpoints from the Lacework Server
type LQLService struct {
	client *Client
}

func (svc *LQLService) Query(query string) (
	response map[string]interface{},
	err error,
) {
	request, err := svc.client.NewRequest(
		"POST",
		apiLQLQuery,
		bytes.NewReader([]byte(query)),
	)
	if err != nil {
		return
	}

	// LQL requires the Content-Type to be 'text/plain'
	svc.client.log.Debug("setting custom header for LQL", zap.String("Content-Type", "text/plain"))
	request.Header.Set("Content-Type", "text/plain")

	res, err := svc.client.DoDecoder(request, &response)
	if err != nil {
		return
	}
	defer res.Body.Close()

	return
}
