//
// Author:: Darren Murray (<darren.murray@lacework.net>)
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
)

// ListMachines returns a list of MachineDetailEntity from the last 7 days
func (svc *EntitiesService) ListMachines() (response MachinesEntityResponse, err error) {
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

// ListMachinesWithFilters returns a list of UserEntity based on a user defined filter
func (svc *EntitiesService) ListMachinesWithFilters(filters SearchFilter) (response MachinesEntityResponse, err error) {
	err = svc.Search(&response, filters)
	return
}

// ListAllMachines iterates over all pages to return all machine details at once
func (svc *EntitiesService) ListAllMachines() (response MachinesEntityResponse, err error) {
	response, err = svc.ListMachines()
	if err != nil {
		return
	}

	var (
		all    []MachineDetailEntity
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

// ListAllMachinesWithFilters iterates over all pages to return all machine details at once based on a user defined filter
func (svc *EntitiesService) ListAllMachinesWithFilters(filters SearchFilter) (response MachinesEntityResponse, err error) {
	response, err = svc.ListMachinesWithFilters(filters)
	if err != nil {
		return
	}

	var (
		all    []MachineDetailEntity
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

type MachinesEntityResponse struct {
	Data   []MachineDetailEntity `json:"data"`
	Paging V2Pagination          `json:"paging"`

	v2PageMetadata `json:"-"`
}

// Fulfill Pageable interface (look at api/v2.go)
func (r MachinesEntityResponse) PageInfo() *V2Pagination {
	return &r.Paging
}
func (r *MachinesEntityResponse) ResetPaging() {
	r.Paging = V2Pagination{}
	r.Data = nil
}

type MachineEntity struct {
	AwsInstanceID string    `json:"awsInstanceId"`
	Hostname      string    `json:"hostname"`
	EntityType    string    `json:"entityType"`
	EndTime       time.Time `json:"endTime"`
	Mid           int       `json:"mid"`
	PrimaryIpAddr string    `json:"primaryIpAddr"`
	StartTime     time.Time `json:"startTime"`
	Tags          struct {
		// Shared Tags
		Cluster             string `json:"Cluster,omitempty"`
		Env                 string `json:"Env,omitempty"`
		Arch                string `json:"arch,omitempty"`
		ExternalIP          string `json:"ExternalIp,omitempty"`
		Hostname            string `json:"Hostname,omitempty"`
		InstanceID          string `json:"InstanceId,omitempty"`
		InternalIP          string `json:"InternalIp,omitempty"`
		LwTokenShort        string `json:"LwTokenShort,omitempty"`
		Os                  string `json:"os,omitempty"`
		VMInstanceType      string `json:"VmInstanceType,omitempty"`
		VMProvider          string `json:"VmProvider,omitempty"`
		Zone                string `json:"Zone,omitempty"`
		ClusterLocation     string `json:"cluster-location,omitempty"`
		ClusterName         string `json:"cluster-name,omitempty"`
		ClusterUid          string `json:"cluster-uid,omitempty"`
		CreatedBy           string `json:"created-by,omitempty"`
		LwKubernetesCluster string `json:"lw_KubernetesCluster,omitempty"`
		KubeLabels          string `json:"kube-labels,omitempty"`

		// AWS Tags
		Account  string `json:"Account,omitempty"`
		AmiId    string `json:"AmiId,omitempty"`
		SubnetId string `json:"SubnetId,omitempty"`
		VpcId    string `json:"VpcId,omitempty"`

		// GCP Tags
		GCEtags                 string `json:"GCEtags,omitempty"`
		InstanceName            string `json:"InstanceName,omitempty"`
		NumericProjectId        string `json:"NumericProjectId,omitempty"`
		ProjectId               string `json:"ProjectId,omitempty"`
		EnableOslogin           string `json:"enable-oslogin,omitempty"`
		GciEnsureGkeDocker      string `json:"gci-ensure-gke-docker,omitempty"`
		GciUpdateStrategy       string `json:"gci-update-strategy,omitempty"`
		GoogleComputeEnablePcid string `json:"google-compute-enable-pcid,omitempty"`
		InstanceTemplate        string `json:"instance-template,omitempty"`
	} `json:"machineTags"`
}
