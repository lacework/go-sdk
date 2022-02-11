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

// ListUsers returns a list of UserEntity from the last 7 days
func (svc *EntitiesService) ListUsers() (response UsersEntityResponse, err error) {
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

// ListAllUsers iterates over all pages to return all user information at once
func (svc *EntitiesService) ListAllUsers() (response UsersEntityResponse, err error) {
	response, err = svc.ListUsers()
	if err != nil {
		return
	}

	var (
		all    []UserEntity
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

type UsersEntityResponse struct {
	Data   []UserEntity `json:"data"`
	Paging V2Pagination `json:"paging"`
}

// Fulfill Pagination interface (look at api/v2.go)
func (r UsersEntityResponse) PageInfo() *V2Pagination {
	return &r.Paging
}
func (r *UsersEntityResponse) ResetPaging() {
	r.Paging = V2Pagination{}
}

type UserEntity struct {
	CreatedTime      time.Time `json:"createdTime"`
	Mid              int       `json:"mid"`
	OtherGroupNames  []string  `json:"otherGroupNames"`
	PrimaryGroupName string    `json:"primaryGroupName"`
	UID              int       `json:"uid"`
	Username         string    `json:"username"`
}
