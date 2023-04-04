//
// Author:: Darren Murray(<darren.murray@lacework.net>)
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

// GetGcpAtSes gets a single GcpAtSes integration matching the provided integration guid
func (svc *CloudAccountsService) GetGcpAtSes(guid string) (
	response GcpAtSesIntegrationResponse,
	err error,
) {
	err = svc.get(guid, &response)
	return
}

// UpdateGcpAtSes updates a single GcpAtSes integration on the Lacework Server
func (svc *CloudAccountsService) UpdateGcpAtSes(data CloudAccount) (
	response GcpAtSesIntegrationResponse,
	err error,
) {
	err = svc.update(data.ID(), data, &response)
	return
}

type GcpAtSesIntegrationResponse struct {
	Data V2GcpAtSesIntegration `json:"data"`
}

type V2GcpAtSesIntegration struct {
	v2CommonIntegrationData
	Data GcpAtSesData `json:"data"`
}

type GcpAtSesData struct {
	Credentials GcpAtSesCredentials `json:"credentials"`
	IDType      string              `json:"idType"`
	// Either the org id or project id
	ID               string `json:"id"`
	SubscriptionName string `json:"subscriptionName"`
}

type GcpAtSesCredentials struct {
	ClientID     string `json:"clientId"`
	ClientEmail  string `json:"clientEmail"`
	PrivateKeyID string `json:"privateKeyId,omitempty"`
	PrivateKey   string `json:"privateKey,omitempty"`
}
