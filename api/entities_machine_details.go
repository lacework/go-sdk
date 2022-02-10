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
	"time"
)

// ListMachineDetails returns a list of MachineDetailEntity from the last 7 days
func (svc *EntitiesService) ListMachineDetails() (response MachineDetailsEntityResponse, err error) {
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

// ListAllMachineDetails iterates over all pages to return all machine details at once
func (svc *EntitiesService) ListAllMachineDetails() (response MachineDetailsEntityResponse, err error) {
	response, err = svc.ListMachineDetails()
	if err != nil {
		return
	}

	all := []MachineDetailEntity{}
	for {
		all = append(all, response.Data...)

		pageOk, err := svc.client.NextPage(&response)
		if err == nil && pageOk {
			continue
		}
		break
	}

	response.Data = all
	response.ResetPaging()
	return
}

type MachineDetailsEntityResponse struct {
	Data   []MachineDetailEntity `json:"data"`
	Paging V2Pagination          `json:"paging"`
}

// Fulfill Pageable interface (look at api/v2.go)
func (r MachineDetailsEntityResponse) PageInfo() *V2Pagination {
	return &r.Paging
}
func (r *MachineDetailsEntityResponse) ResetPaging() {
	r.Paging = V2Pagination{}
}

type MachineDetailEntity struct {
	AwsInstanceID string    `json:"awsInstanceId"`
	AwsZone       string    `json:"awsZone"`
	CreatedTime   time.Time `json:"createdTime"`
	Domain        string    `json:"domain"`
	Hostname      string    `json:"hostname"`
	Kernel        string    `json:"kernel"`
	KernelRelease string    `json:"kernelRelease"`
	KernelVersion string    `json:"kernelVersion"`
	Mid           int       `json:"mid"`
	Os            string    `json:"os"`
	OsVersion     string    `json:"osVersion"`
	Tags          struct {
		Account        string `json:"Account"`
		AmiID          string `json:"AmiId"`
		ExternalIP     string `json:"ExternalIp"`
		Hostname       string `json:"Hostname"`
		Name           string `json:"Name"`
		InstanceID     string `json:"InstanceId"`
		InternalIP     string `json:"InternalIp"`
		LwTokenShort   string `json:"LwTokenShort"`
		SubnetID       string `json:"SubnetId"`
		VMInstanceType string `json:"VmInstanceType"`
		VMProvider     string `json:"VmProvider"`
		VpcID          string `json:"VpcId"`
		Zone           string `json:"Zone"`
		Arch           string `json:"arch"`
		Os             string `json:"os"`
	} `json:"tags"`
}
