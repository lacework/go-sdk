//
// Copyright:: Copyright 2026, Lacework Inc.
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

import "fmt"

// FortiDSPM-managed DSPM cloud accounts.
//
// They use the AwsDspm / AzureDspm cloud account types but carry only the
// account identity and the regions: FortiDSPM pushes scan results to Lacework
// itself, so there is no customer bucket and no cross-account credentials.
// On create, api-server asks FortiDSPM for a deployment and returns it as
// fortiDspmDeployment: one activation token and image reference per region,
// consumed by the terraform that builds the scan engines.

// GetAwsFortiDspm gets a FortiDSPM-managed AWS DSPM cloud account by guid.
func (svc *CloudAccountsService) GetAwsFortiDspm(guid string) (
	response AwsFortiDspmResponse,
	err error,
) {
	err = svc.get(guid, &response)
	return
}

// CreateAwsFortiDspm creates a FortiDSPM-managed AWS DSPM cloud account. The
// response carries the FortiDSPM deployment; it is issued once and is not
// returned by later reads.
func (svc *CloudAccountsService) CreateAwsFortiDspm(data CloudAccount) (
	response AwsFortiDspmResponse,
	err error,
) {
	err = svc.create(data, &response)
	return
}

// GetAzureFortiDspm gets a FortiDSPM-managed Azure DSPM cloud account by guid.
func (svc *CloudAccountsService) GetAzureFortiDspm(guid string) (
	response AzureFortiDspmResponse,
	err error,
) {
	err = svc.get(guid, &response)
	return
}

// CreateAzureFortiDspm creates a FortiDSPM-managed Azure DSPM cloud account.
func (svc *CloudAccountsService) CreateAzureFortiDspm(data CloudAccount) (
	response AzureFortiDspmResponse,
	err error,
) {
	err = svc.create(data, &response)
	return
}

type AwsFortiDspmResponse struct {
	Data AwsFortiDspm `json:"data"`
}

type AwsFortiDspm struct {
	v2CommonIntegrationData
	Data                AwsFortiDspmData     `json:"data"`
	FortiDspmDeployment *FortiDspmDeployment `json:"fortiDspmDeployment,omitempty"`
}

type AwsFortiDspmData struct {
	AccountID string   `json:"awsAccountId"`
	Regions   []string `json:"regions"`
}

type AzureFortiDspmResponse struct {
	Data AzureFortiDspm `json:"data"`
}

type AzureFortiDspm struct {
	v2CommonIntegrationData
	Data                AzureFortiDspmData   `json:"data"`
	FortiDspmDeployment *FortiDspmDeployment `json:"fortiDspmDeployment,omitempty"`
}

type AzureFortiDspmData struct {
	TenantID       string   `json:"tenantId"`
	SubscriptionID string   `json:"subscriptionId,omitempty"`
	Regions        []string `json:"regions"`
}

// FortiDspmDeployment is what FortiDSPM issued for a new cloud account.
type FortiDspmDeployment struct {
	DeploymentID      string            `json:"deploymentId"`
	DeploymentName    string            `json:"deploymentName"`
	EnvID             string            `json:"envId,omitempty"`
	TokenExpiresIn    int               `json:"tokenExpiresIn,omitempty"`
	ImageURLExpiresIn int               `json:"imageUrlExpiresIn,omitempty"`
	Regions           []FortiDspmRegion `json:"regions"`
}

// FortiDspmRegion is one region's scan engine input: a single-use activation
// token and the image to boot. AWS carries an AMI id; Azure carries a SAS URL
// to the image VHD plus its Hyper-V generation.
type FortiDspmRegion struct {
	Region           string `json:"region"`
	ActivationToken  string `json:"activationToken"`
	AmiID            string `json:"amiId,omitempty"`
	ImageID          string `json:"imageId,omitempty"`
	ImageURL         string `json:"imageUrl,omitempty"`
	HypervGeneration string `json:"hypervGeneration,omitempty"`
	ImageSizeBytes   int64  `json:"imageSizeBytes,omitempty"`
	ImageSha256      string `json:"imageSha256,omitempty"`
}

// ActivationTokens returns the tokens keyed by region.
func (d *FortiDspmDeployment) ActivationTokens() map[string]string {
	out := map[string]string{}
	if d == nil {
		return out
	}
	for _, r := range d.Regions {
		out[r.Region] = r.ActivationToken
	}
	return out
}

// Images returns the image reference keyed by region: the AMI id on AWS, the
// image URL on Azure (falling back to a gallery image id when that is what
// FortiDSPM sent).
func (d *FortiDspmDeployment) Images() map[string]string {
	out := map[string]string{}
	if d == nil {
		return out
	}
	for _, r := range d.Regions {
		switch {
		case r.AmiID != "":
			out[r.Region] = r.AmiID
		case r.ImageURL != "":
			out[r.Region] = r.ImageURL
		default:
			out[r.Region] = r.ImageID
		}
	}
	return out
}

// HypervGenerations returns the Azure Hyper-V generation keyed by region for
// the regions that carry one.
func (d *FortiDspmDeployment) HypervGenerations() map[string]string {
	out := map[string]string{}
	if d == nil {
		return out
	}
	for _, r := range d.Regions {
		if r.HypervGeneration != "" {
			out[r.Region] = r.HypervGeneration
		}
	}
	return out
}

// ReportFortiDspmDeploymentStatus tells FortiDSPM how terraform ended for a
// FortiDSPM-managed integration (PUT /api/v2/FortiDspm/deployments/{guid}/status).
// FortiDSPM merges regions by name, so a report per region is fine.
func (svc *CloudAccountsService) ReportFortiDspmDeploymentStatus(
	intgGuid string, report FortiDspmDeploymentStatus,
) error {
	var response map[string]interface{}
	return svc.client.RequestEncoderDecoder("PUT",
		fmt.Sprintf(apiV2FortiDspmDeploymentStatus, intgGuid), report, &response)
}

// FortiDspmDeploymentStatus is a terraform outcome report. Status is one of
// succeeded, failed, rolled_back, deleted_with_resources, deleted_orphaned.
type FortiDspmDeploymentStatus struct {
	DeploymentID string                    `json:"deploymentId"`
	Status       string                    `json:"status"`
	Regions      []FortiDspmRegionStatus   `json:"regions,omitempty"`
	Error        *FortiDspmDeploymentError `json:"error,omitempty"`
}

// FortiDspmRegionStatus is one region's outcome and the resource ids FortiDSPM
// shows for the scan engine.
type FortiDspmRegionStatus struct {
	Region              string `json:"region"`
	Status              string `json:"status"`
	InstanceID          string `json:"instanceId,omitempty"`
	PrivateIP           string `json:"privateIp,omitempty"`
	NatGatewayPublicIP  string `json:"natGatewayPublicIp,omitempty"`
	IamRoleArn          string `json:"iamRoleArn,omitempty"`
	IdentityPrincipalID string `json:"identityPrincipalId,omitempty"`
}

type FortiDspmDeploymentError struct {
	Phase   string `json:"phase,omitempty"`
	Message string `json:"message"`
}
