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

// AgentInfoService is the service that interacts with
// the AgentInfo schema from the Lacework APIv2 Server
type AgentInfoService struct {
	client *Client
}

func (svc *AgentInfoService) Search(response interface{}, filters SearchFilter) error {
	return svc.client.RequestEncoderDecoder("POST", apiV2AgentInfoSearch, filters, response)
}

type AgentInfoResponse struct {
	Data   []AgentInfo  `json:"data"`
	Paging V2Pagination `json:"paging"`
}

// Fulfill Pageable interface (look at api/v2.go)
func (r AgentInfoResponse) PageInfo() *V2Pagination {
	return &r.Paging
}
func (r *AgentInfoResponse) ResetPaging() {
	r.Paging = V2Pagination{}
}

type AgentInfo struct {
	AgentVersion string    `json:"agentVersion"`
	CreatedTime  time.Time `json:"createdTime"`
	Hostname     string    `json:"hostname"`
	IpAddr       string    `json:"ipAddr"`
	LastUpdate   time.Time `json:"lastUpdate"`
	Mid          int       `json:"mid"`
	Mode         string    `json:"mode"`
	Os           string    `json:"os"`
	Status       string    `json:"status"`
	Tags         struct {
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
