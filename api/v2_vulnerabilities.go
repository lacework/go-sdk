//
// Author:: Salim Afiune Maya (<afiune@lacework.net>)
// Copyright:: Copyright 2021, Lacework Inc.
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

// v2VulnerabilitiesService is a service that interacts with the Vulnerabilities
// endpoints from the Lacework APIv2 Server
type v2VulnerabilitiesService struct {
	client     *Client
	Containers *v2ContainerVulnerabilityService
}

func NewV2VulnerabilitiesService(c *Client) *v2VulnerabilitiesService {
	return &v2VulnerabilitiesService{c,
		&v2ContainerVulnerabilityService{c},
	}
}

// v2ContainerVulnerabilityService is a service that interacts with the APIv2
// vulnerabilities endpoints for containers
type v2ContainerVulnerabilityService struct {
	client *Client
}

// SearchLastWeek returns a list of VulnerabilityContainer from the last 7 days
func (svc *v2ContainerVulnerabilityService) SearchLastWeek() (VulnerabilitiesContainersResponse, error) {
	var (
		now    = time.Now().UTC()
		before = now.AddDate(0, 0, -7) // 7 days from ago
	)

	return svc.Search(SearchFilter{
		TimeFilter: &TimeFilter{
			StartTime: &before,
			EndTime:   &now,
		},
	})
}

// Search returns a list of VulnerabilityContainer from the last 7 days
func (svc *v2ContainerVulnerabilityService) Search(filters SearchFilter) (
	response VulnerabilitiesContainersResponse, err error,
) {
	err = svc.client.RequestEncoderDecoder(
		"POST", apiV2VulnerabilitiesContainersSearch,
		filters, &response,
	)
	return
}

// SearchAllPages iterates over all pages and returns a list of VulnerabilityContainer
func (svc *v2ContainerVulnerabilityService) SearchAllPages(filters SearchFilter) (
	response VulnerabilitiesContainersResponse, err error,
) {
	response, err = svc.Search(filters)
	if err != nil {
		return
	}

	var (
		all    = []VulnerabilityContainer{}
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

type VulnerabilitiesContainersResponse struct {
	Data   []VulnerabilityContainer `json:"data"`
	Paging V2Pagination             `json:"paging"`
}

// Fulfill Pagination interface (look at api/v2.go)
func (r VulnerabilitiesContainersResponse) PageInfo() *V2Pagination {
	return &r.Paging
}
func (r *VulnerabilitiesContainersResponse) ResetPaging() {
	r.Paging = V2Pagination{}
}

type VulnerabilityContainer struct {
	EvalCtx struct {
		CveBatchInfo []struct {
			CveBatchID     string `json:"cve_batch_id"`
			CveCreatedTime string `json:"cve_created_time"`
		} `json:"cve_batch_info"`
		ExceptionProps []struct {
			Status string `json:"status"`
		} `json:"exception_props"`
		ImageInfo struct {
			CreatedTime int64    `json:"created_time"`
			Digest      string   `json:"digest"`
			ErrorMsg    []string `json:"error_msg"`
			ID          string   `json:"id"`
			Registry    string   `json:"registry"`
			Repo        string   `json:"repo"`
			Size        int      `json:"size"`
			Status      string   `json:"status"`
			Tags        []string `json:"tags"`
			Type        string   `json:"type"`
		} `json:"image_info"`
		IsDailyJob       string `json:"isDailyJob"`
		IsReeval         bool   `json:"is_reeval"`
		ScanBatchID      string `json:"scan_batch_id"`
		ScanCreatedTime  string `json:"scan_created_time"`
		ScanRequestProps struct {
			DataFormatVersion string `json:"data_format_version"`
			Environment       struct {
				DockerVersion struct {
					ErrorMessage string `json:"error_message"`
				} `json:"docker_version"`
			} `json:"environment"`
			Props struct {
				DataFormatVersion string `json:"data_format_version"`
				ScannerVersion    string `json:"scanner_version"`
			} `json:"props"`
			ScanCompletionUtcTime int    `json:"scanCompletionUtcTime"`
			ScanStartTime         int    `json:"scan_start_time"`
			ScannerVersion        string `json:"scanner_version"`
		} `json:"scan_request_props"`
		VulnBatchID     string `json:"vuln_batch_id"`
		VulnCreatedTime string `json:"vuln_created_time"`
	} `json:"evalCtx"`
	FeatureKey struct {
		Name      string `json:"name"`
		Namespace string `json:"namespace"`
		Version   string `json:"version"`
	} `json:"featureKey"`
	FixInfo struct {
		CompareResult int    `json:"compare_result"`
		FixAvailable  int    `json:"fix_available"`
		FixedVersion  string `json:"fixed_version"`
	} `json:"fixInfo"`
	ImageID   string    `json:"imageId"`
	Severity  string    `json:"severity"`
	StartTime time.Time `json:"startTime"`
	Status    string    `json:"status"`
	VulnID    string    `json:"vulnId"`
}
