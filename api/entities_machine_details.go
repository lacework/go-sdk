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

// ListMachineDetailsWithFilters returns a list of UserEntity based on a user defined filter
func (svc *EntitiesService) ListMachineDetailsWithFilters(filters SearchFilter) (response MachineDetailsEntityResponse, err error) {
	err = svc.Search(&response, filters)
	return
}

// ListAllMachineDetails iterates over all pages to return all machine details at once
func (svc *EntitiesService) ListAllMachineDetails() (response MachineDetailsEntityResponse, err error) {
	response, err = svc.ListMachineDetails()
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

// ListAllMachineDetailsWithFilters iterates over all pages to return all machine details at once based on a user defined filter
func (svc *EntitiesService) ListAllMachineDetailsWithFilters(filters SearchFilter) (response MachineDetailsEntityResponse, err error) {
	response, err = svc.ListMachineDetailsWithFilters(filters)
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

type MachineDetailsEntityResponse struct {
	Data   []MachineDetailEntity `json:"data"`
	Paging V2Pagination          `json:"paging"`

	v2PageMetadata `json:"-"`
}

// Fulfill Pageable interface (look at api/v2.go)
func (r MachineDetailsEntityResponse) PageInfo() *V2Pagination {
	return &r.Paging
}
func (r *MachineDetailsEntityResponse) ResetPaging() {
	r.Paging = V2Pagination{}
	r.Data = nil
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
		// Shared Tags
		Arch           string `json:"arch,omitempty"`
		ExternalIP     string `json:"ExternalIp,omitempty"`
		Hostname       string `json:"Hostname,omitempty"`
		InstanceID     string `json:"InstanceId,omitempty"`
		InternalIP     string `json:"InternalIp,omitempty"`
		LwTokenShort   string `json:"LwTokenShort,omitempty"`
		Os             string `json:"os,omitempty"`
		VMInstanceType string `json:"VmInstanceType,omitempty"`
		VMProvider     string `json:"VmProvider,omitempty"`
		Zone           string `json:"Zone,omitempty"`

		// AWS Tags
		Account  string `json:"Account,omitempty"`
		AmiID    string `json:"AmiId,omitempty"`
		Name     string `json:"Name,omitempty"`
		SubnetID string `json:"SubnetId,omitempty"`
		VpcID    string `json:"VpcId,omitempty"`

		// GCP Tags
		Cluster                 string `json:"Cluster,omitempty"`
		ClusterLocation         string `json:"cluster-location,omitempty"`
		ClusterName             string `json:"cluster-name,omitempty"`
		ClusterUID              string `json:"cluster-uid,omitempty"`
		CreatedBy               string `json:"created-by,omitempty"`
		EnableOSLogin           string `json:"enable-oslogin,omitempty"`
		Env                     string `json:"Env,omitempty"`
		GCEtags                 string `json:"GCEtags,omitempty"`
		GCIEnsureGKEDocker      string `json:"gci-ensure-gke-docker,omitempty"`
		GCIUpdateStrategy       string `json:"gci-update-strategy,omitempty"`
		GoogleComputeEnablePCID string `json:"google-compute-enable-pcid,omitempty"`
		InstanceName            string `json:"InstanceName,omitempty"`
		InstanceTemplate        string `json:"InstanceTemplate,omitempty"`
		KubeLabels              string `json:"kube-labels,omitempty"`
		LWKubernetesCluster     string `json:"lw_KubernetesCluster,omitempty"`
		NumericProjectID        string `json:"NumericProjectId,omitempty"`
		ProjectID               string `json:"ProjectId,omitempty"`
	} `json:"tags"`
}
