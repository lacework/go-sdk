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

package cmd

import (
	"encoding/json"
	"io/ioutil"
)

// @afiune POC - This depends on LQL
func loadAgents() (AgentsResponse, error) {
	data := AgentsResponse{}
	file, err := ioutil.ReadFile("/tmp/Card70.json")
	if err != nil {
		return data, err
	}
	err = json.Unmarshal([]byte(file), &data)
	return data, err
}

type AgentsResponse struct {
	Data []AgentHost `json:"data"`
	Ok   bool        `json:"ok"`
}

type AgentHost struct {
	MachineHostname string `json:"MACHINE_HOSTNAME"`
	Name            string `json:"NAME"`
	MachineIP       string `json:"MACHINE_IP_ADDR"`
	Status          string `json:"STATUS"`
	//CreatedTime     int64       `json:"CREATED_TIME"`
	AgentVersion string      `json:"AGENT_VERSION"`
	Mode         string      `json:"MODE"`
	Tags         machineTags `json:"TAGS"`
}
type machineTags struct {
	Account        string `json:"Account"`
	AmiID          string `json:"AmiId"`
	ExternalIP     string `json:"ExternalIp"`
	Hostname       string `json:"Hostname"`
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
}
