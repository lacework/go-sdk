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

import "time"

// ListImages returns a list of UserEntity from the last 7 days
func (svc *EntitiesService) ListImages() (response ImagesEntityResponse, err error) {
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

// ListImagesWithFilters returns a list of UserEntity based on a user defined filter
func (svc *EntitiesService) ListImagesWithFilters(filters SearchFilter) (response ImagesEntityResponse, err error) {
	err = svc.Search(&response, filters)
	return
}

// ListAllImages iterates over all pages to return all images information at once
func (svc *EntitiesService) ListAllImages() (response ImagesEntityResponse, err error) {
	response, err = svc.ListImages()
	if err != nil {
		return
	}

	var (
		all    []ImageEntity
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

	response.ResetPaging()
	response.Data = all
	return
}

// ListAllImagesWithFilters iterates over all pages to return all images information
// at once based on a user defined filter
func (svc *EntitiesService) ListAllImagesWithFilters(filters SearchFilter) (
	response ImagesEntityResponse, err error,
) {
	response, err = svc.ListImagesWithFilters(filters)
	if err != nil {
		return
	}

	var (
		all    []ImageEntity
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

	response.ResetPaging()
	response.Data = all
	return
}

type ImagesEntityResponse struct {
	Data   []ImageEntity `json:"data"`
	Paging V2Pagination  `json:"paging"`

	v2PageMetadata `json:"-"`
}

// Fulfill Pageable interface (look at api/v2.go)
func (r ImagesEntityResponse) PageInfo() *V2Pagination {
	return &r.Paging
}
func (r *ImagesEntityResponse) ResetPaging() {
	r.Paging = V2Pagination{}
	r.Data = nil
}

type ImageEntity struct {
	ContainerType string    `json:"containerType"`
	CreatedTime   time.Time `json:"createdTime"`
	ImageID       string    `json:"imageId"`
	Mid           int       `json:"mid"`
	Repo          string    `json:"repo"`
	Size          int       `json:"size"`
	Tag           string    `json:"tag"`
}
