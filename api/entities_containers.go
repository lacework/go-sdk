//
// Author:: Salim Afiune Maya (<afiune@lacework.net>)
// Copyright:: Copyright 2023, Lacework Inc.
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
	"time"

	"github.com/lacework/go-sdk/internal/array"
)

// ListContainers returns a list of Active Containers from the last 7 days
func (svc *EntitiesService) ListContainers() (response ContainersEntityResponse, err error) {
	now := time.Now().UTC()
	before := now.AddDate(0, 0, -7) // 7 days from ago
	err = svc.Search(&response,
		SearchFilter{
			TimeFilter: &TimeFilter{
				StartTime: &before,
				EndTime:   &now,
			},
		},
	)
	return
}

// ListContainersWithFilters returns a list of Active Containers based on a user defined filter
func (svc *EntitiesService) ListContainersWithFilters(filters SearchFilter) (response ContainersEntityResponse, err error) {
	err = svc.Search(&response, filters)
	return
}

// ListAllContainers iterates over all pages to return all active container information at once
func (svc *EntitiesService) ListAllContainers() (response ContainersEntityResponse, err error) {
	response, err = svc.ListContainers()
	if err != nil {
		return
	}

	var (
		all    []ContainerEntity
		pageOk bool
	)
	for {
		all = append(all, response.Data...)

		pageOk, err = svc.client.NextPage(&response)
		if err == nil && pageOk {
			continue
		}
		break
	}

	response.Data = all
	response.ResetPaging()
	return
}

// ListAllContainersWithFilters iterates over all pages to return all active container information at once based on a user defined filter
func (svc *EntitiesService) ListAllContainersWithFilters(filters SearchFilter) (response ContainersEntityResponse, err error) {
	response, err = svc.ListContainersWithFilters(filters)
	if err != nil {
		return
	}

	var (
		all    []ContainerEntity
		pageOk bool
	)

	for {
		all = append(all, response.Data...)

		pageOk, err = svc.client.NextPage(&response)
		if err == nil && pageOk {
			continue
		}
		break
	}

	response.Data = all
	response.ResetPaging()
	return
}

type ContainersEntityResponse struct {
	Data   []ContainerEntity `json:"data"`
	Paging V2Pagination      `json:"paging"`
}

// Fulfill Pageable interface (look at api/v2.go)
func (r ContainersEntityResponse) PageInfo() *V2Pagination {
	return &r.Paging
}
func (r *ContainersEntityResponse) ResetPaging() {
	r.Paging = V2Pagination{}
}

// Total returns the total number of active containers
func (r *ContainersEntityResponse) Total() int {
	uniqMIDs := []int{}
	for _, container := range r.Data {
		if !array.ContainsInt(uniqMIDs, container.Mid) {
			uniqMIDs = append(uniqMIDs, container.Mid)
		}
	}
	return len(uniqMIDs)
}

// Count returns the number of active containers with the provided image ID
func (r *ContainersEntityResponse) Count(imageID string) int {
	uniqMIDs := []int{}
	for _, container := range r.Data {
		if container.ImageID == imageID &&
			!array.ContainsInt(uniqMIDs, container.Mid) {
			uniqMIDs = append(uniqMIDs, container.Mid)
		}
	}
	return len(uniqMIDs)
}

type ContainerEntity struct {
	ContainerName  string                 `json:"containerName"`
	ImageID        string                 `json:"imageId"`
	Mid            int                    `json:"mid"`
	StartTime      time.Time              `json:"startTime"`
	EndTime        time.Time              `json:"endTime"`
	PodName        string                 `json:"podName"`
	PropsContainer map[string]interface{} `json:"propsContainer"`
	Tags           map[string]interface{} `json:"tags"`
}
