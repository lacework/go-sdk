//
// Author:: Salim Afiune Maya (<afiune@lacework.net>)
// Copyright:: Copyright 2022, Lacework Inc.
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

type EntitiesService struct {
	client *Client
}

// Search expects the response and the search filters
//
// e.g.
//
//   response := api.MachineDetailsResponse{}
//   lacework.V2.Entities.Search(response, api.SearchFilter{})
//
func (svc *EntitiesService) Search(response interface{}, filters SearchFilter) error {
	var apiPath string

	switch response.(type) {
	case *MachineDetailsResponse:
		apiPath = fmt.Sprintf(apiV2EntitiesSearch, "MachineDetails")
	default:
		return errors.New("missing implementation for the provided entity response")
	}

	return svc.client.RequestEncoderDecoder("POST", apiPath, filters, response)
}
