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

type EntityType int

const (
	NoneEntityType EntityType = iota
	MachineDetailsEntityType
	UsersEntityType
	ImagesEntityType
	ContainersEntityType
	MachineEntityType
)

// EntityTypes is the list of available entity types
var EntityTypes = map[EntityType]string{
	NoneEntityType:           "None",
	MachineDetailsEntityType: "MachineDetails",
	MachineEntityType: "Machines",
	UsersEntityType:          "Users",
	ImagesEntityType:         "Images",
	ContainersEntityType:     "Containers",
}

// Search expects the response and the search filters
//
// e.g.
//
//   var (
//       response = &api.MachineDetailsEntityResponse{}
//       now      = time.Now().UTC()
//       before   = now.AddDate(0, 0, -7) // 7 days from ago
//       filters  = api.SearchFilter{
//           TimeFilter: &api.TimeFilter{
//               StartTime: &before,
//               EndTime:   &now,
//           },
//       }
//   )
//   lacework.V2.Entities.Search(response, filters)
//
func (svc *EntitiesService) Search(response interface{}, filters SearchFilter) error {
	var apiPath string

	switch response.(type) {
	case *MachineDetailsEntityResponse:
		apiPath = fmt.Sprintf(apiV2EntitiesSearch, EntityTypes[MachineDetailsEntityType])

	case *UsersEntityResponse:
		apiPath = fmt.Sprintf(apiV2EntitiesSearch, EntityTypes[UsersEntityType])

	case *ImagesEntityResponse:
		apiPath = fmt.Sprintf(apiV2EntitiesSearch, EntityTypes[ImagesEntityType])

	case *ContainersEntityResponse:
		apiPath = fmt.Sprintf(apiV2EntitiesSearch, EntityTypes[ContainersEntityType])

	case *MachinesEntityResponse:
		apiPath = fmt.Sprintf(apiV2EntitiesSearch, EntityTypes[MachineEntityType])

	default:
		return errors.New("missing implementation for the provided entity response")
	}

	return svc.client.RequestEncoderDecoder("POST", apiPath, filters, response)
}
