//
// Author:: Ammar Ekbote(<ammar.ekbote@lacework.net>)
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
	"encoding/base64"
	"fmt"
	"strings"
)

// GetGcpSidekick gets a single GcpSidekick integration matching the provided integration guid
func (svc *CloudAccountsService) GetGcpSidekick(guid string) (
	response GcpSidekickIntegrationResponse,
	err error,
) {
	err = svc.get(guid, &response)
	return
}

// CreateGcpSidekick creates an GcpSidekick Cloud Account integration
func (svc *CloudAccountsService) CreateGcpSidekick(data CloudAccount) (
	response GcpSidekickIntegrationResponse,
	err error,
) {
	err = svc.create(data, &response)
	return
}

// UpdateGcpSidekick updates a single GcpSidekick integration on the Lacework Server
func (svc *CloudAccountsService) UpdateGcpSidekick(data CloudAccount) (
	response GcpSidekickIntegrationResponse,
	err error,
) {
	err = svc.update(data.ID(), data, &response)
	return
}

type GcpSidekickIntegrationResponse struct {
	Data V2GcpSidekickIntegration `json:"data"`
}

type GcpSidekickToken struct {
	ServerToken string `json:"serverToken"`
	Uri         string `json:"uri"`
}

type V2GcpSidekickIntegration struct {
	v2CommonIntegrationData
	GcpSidekickToken `json:"serverToken"`
	Data             GcpSidekickData `json:"data"`
}

type GcpSidekickData struct {
	Credentials GcpSidekickCredentials `json:"credentials"`
	IDType      string                 `json:"idType"`
	// Either the org id or project id
	ID                string `json:"id"`
	ScanningProjectId string `json:"scanningProjectId"`
	SharedBucket      string `json:"sharedBucketName"`
	FilterList        string `json:"filterList,omitempty"`
	QueryText         string `json:"queryText,omitempty"`
	//ScanFrequency in hours, 24 == 24 hours
	ScanFrequency           int    `json:"scanFrequency"`
	ScanContainers          bool   `json:"scanContainers"`
	ScanHostVulnerabilities bool   `json:"scanHostVulnerabilities"`
	ScanMultiVolume         bool   `json:"scanMultiVolume"`
	AccountMappingFile      string `json:"accountMappingFile,omitempty"`
}

type GcpSidekickCredentials struct {
	ClientID     string `json:"clientId"`
	ClientEmail  string `json:"clientEmail"`
	PrivateKeyID string `json:"privateKeyId"`
	PrivateKey   string `json:"privateKey,omitempty"`
	TokenUri     string `json:"tokenUri,omitempty"`
}

func (gcp *GcpSidekickData) EncodeAccountMappingFile(mapping []byte) {
	encodedMappings := base64.StdEncoding.EncodeToString(mapping)
	gcp.AccountMappingFile = fmt.Sprintf("data:application/json;name=i.json;base64,%s", encodedMappings)
}

func (gcp *GcpSidekickData) DecodeAccountMappingFile() ([]byte, error) {
	if len(gcp.AccountMappingFile) == 0 {
		return []byte{}, nil
	}

	var (
		b64      = strings.Split(gcp.AccountMappingFile, ",")
		raw, err = base64.StdEncoding.DecodeString(b64[1])
	)
	if err != nil {
		return []byte{}, err
	}

	return raw, nil
}
