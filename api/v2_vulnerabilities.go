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
	"strconv"
	"time"
)

// v2VulnerabilitiesService is a service that interacts with the Vulnerabilities
// endpoints from the Lacework APIv2 Server
type v2VulnerabilitiesService struct {
	client           *Client
	Hosts            *v2HostVulnerabilityService
	Containers       *v2ContainerVulnerabilityService
	SoftwarePackages *v2SoftwarePackagesVulnerabilityService
}

func NewV2VulnerabilitiesService(c *Client) *v2VulnerabilitiesService {
	return &v2VulnerabilitiesService{c,
		&v2HostVulnerabilityService{c},
		&v2ContainerVulnerabilityService{c},
		&v2SoftwarePackagesVulnerabilityService{c},
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
		all    []VulnerabilityContainer
		pageOk bool
	)
	for {
		all = append(all, response.Data...)

		newResponse := VulnerabilitiesContainersResponse{
			Paging: response.Paging,
		}
		pageOk, err = svc.client.NextPage(&newResponse)
		if err == nil && pageOk {
			response = newResponse
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

// v2HostVulnerabilityService is a service that interacts with the APIv2
// vulnerabilities endpoints for hosts
type v2HostVulnerabilityService struct {
	client *Client
}

// SearchLastWeek returns a list of VulnerabilityHost from the last 7 days
func (svc *v2HostVulnerabilityService) SearchLastWeek() (VulnerabilitiesHostResponse, error) {
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

// Search returns a list of VulnerabilityHost from the last 7 days
func (svc *v2HostVulnerabilityService) Search(filters SearchFilter) (
	response VulnerabilitiesHostResponse, err error,
) {
	err = svc.client.RequestEncoderDecoder(
		"POST", apiV2VulnerabilitiesHostsSearch,
		filters, &response,
	)
	return
}

// SearchAllPages iterates over all pages and returns a list of VulnerabilityHost
func (svc *v2HostVulnerabilityService) SearchAllPages(filters SearchFilter) (
	response VulnerabilitiesHostResponse, err error,
) {
	response, err = svc.Search(filters)
	if err != nil {
		return
	}

	var (
		all    []VulnerabilityHost
		pageOk bool
	)
	for {
		all = append(all, response.Data...)

		newResponse := VulnerabilitiesHostResponse{
			Paging: response.Paging,
		}
		pageOk, err = svc.client.NextPage(&newResponse)
		if err == nil && pageOk {
			response = newResponse
			continue
		}
		break
	}

	response.Data = all
	response.ResetPaging()
	return
}

type VulnerabilitiesHostResponse struct {
	Data   []VulnerabilityHost `json:"data"`
	Paging V2Pagination        `json:"paging"`
}

// Fulfill Pagination interface (look at api/v2.go)
func (r VulnerabilitiesHostResponse) PageInfo() *V2Pagination {
	return &r.Paging
}
func (r *VulnerabilitiesHostResponse) ResetPaging() {
	r.Paging = V2Pagination{}
}

type VulnerabilityHost struct {
	CveProps struct {
		CveBatchID  string                     `json:"cve_batch_id"`
		Description string                     `json:"description"`
		Link        string                     `json:"link"`
		Metadata    *VulnerabilityHostMetadata `json:"metadata,omitempty"`
	} `json:"cveProps"`
	EndTime time.Time `json:"endTime"`
	EvalCtx struct {
		ExceptionProps []interface{} `json:"exception_props"`
		Hostname       string        `json:"hostname"`
		McEvalGUID     string        `json:"mc_eval_guid"`
	} `json:"evalCtx"`
	FeatureKey struct {
		Name             string `json:"name"`
		Namespace        string `json:"namespace"`
		PackageActive    int    `json:"package_active"`
		VersionInstalled string `json:"version_installed"`
	} `json:"featureKey"`
	FixInfo struct {
		CompareResult               string `json:"compare_result"`
		EvalStatus                  string `json:"eval_status"`
		FixAvailable                string `json:"fix_available"`
		FixedVersion                string `json:"fixed_version"`
		FixedVersionComparisonInfos []struct {
			CurrFixVer                         string `json:"curr_fix_ver"`
			IsCurrFixVerGreaterThanOtherFixVer string `json:"is_curr_fix_ver_greater_than_other_fix_ver"`
			OtherFixVer                        string `json:"other_fix_ver"`
		} `json:"fixed_version_comparison_infos"`
		FixedVersionComparisonScore int    `json:"fixed_version_comparison_score"`
		VersionInstalled            string `json:"version_installed"`
	} `json:"fixInfo"`
	MachineTags struct {
		Account                               string `json:"Account"`
		AmiID                                 string `json:"AmiId"`
		Env                                   string `json:"Env"`
		ExternalIP                            string `json:"ExternalIp"`
		Hostname                              string `json:"Hostname"`
		InstanceID                            string `json:"InstanceId"`
		InternalIP                            string `json:"InternalIp"`
		LwTokenShort                          string `json:"LwTokenShort"`
		Name                                  string `json:"Name"`
		SubnetID                              string `json:"SubnetId"`
		VMInstanceType                        string `json:"VmInstanceType"`
		VMProvider                            string `json:"VmProvider"`
		VpcID                                 string `json:"VpcId"`
		Zone                                  string `json:"Zone"`
		AlphaEksctlIoNodegroupName            string `json:"alpha.eksctl.io/nodegroup-name"`
		AlphaEksctlIoNodegroupType            string `json:"alpha.eksctl.io/nodegroup-type"`
		Arch                                  string `json:"arch"`
		AwsAutoscalingGroupName               string `json:"aws:autoscaling:groupName"`
		AwsEc2FleetID                         string `json:"aws:ec2:fleet-id"`
		AwsEc2LaunchtemplateID                string `json:"aws:ec2launchtemplate:id"`
		AwsEc2LaunchtemplateVersion           string `json:"aws:ec2launchtemplate:version"`
		EksClusterName                        string `json:"eks:cluster-name"`
		EksNodegroupName                      string `json:"eks:nodegroup-name"`
		K8SIoClusterAutoscalerEnabled         int    `json:"k8s.io/cluster-autoscaler/enabled"`
		K8SIoClusterAutoscalerTechallySandbox string `json:"k8s.io/cluster-autoscaler/techally-sandbox"`
		KubernetesIoClusterTechallySandbox    string `json:"kubernetes.io/cluster/techally-sandbox"`
		LwKubernetesCluster                   string `json:"lw_KubernetesCluster"`
		Os                                    string `json:"os"`
	} `json:"machineTags"`
	Props     VulnerabilityHostProps `json:"props"`
	Mid       int                    `json:"mid"`
	Severity  string                 `json:"severity"`
	StartTime time.Time              `json:"startTime"`
	Status    string                 `json:"status"`
	VulnID    string                 `json:"vulnId"`
}

func (v *VulnerabilityHost) PackageActive() string {
	if v.FeatureKey.PackageActive == 0 {
		return ""
	}
	return "ACTIVE"
}

func (v *VulnerabilityHost) CvssV2() (cvssV2Score string) {
	if v.CveProps.Metadata != nil {
		score := v.CveProps.Metadata.NVD.CVSSv2.Score
		cvssV2Score = strconv.FormatFloat(score, 'f', 1, 64)
	}
	return
}

func (v *VulnerabilityHost) CvssV3() (cvssV3Score string) {
	if v.CveProps.Metadata != nil {
		score := v.CveProps.Metadata.NVD.CVSSv3.Score
		cvssV3Score = strconv.FormatFloat(score, 'f', 1, 64)
	}
	return
}

type VulnerabilityHostMetadata struct {
	NVD struct {
		CVSSv2 struct {
			PublishedDateTime string  `json:"PublishedDateTime"`
			Score             float64 `json:"Score"`
			Vectors           string  `json:"Vectors"`
		} `json:"CVSSv2"`
		CVSSv3 struct {
			ExploitabilityScore float64 `json:"ExploitabilityScore"`
			ImpactScore         float64 `json:"ImpactScore"`
			Score               float64 `json:"Score"`
			Vectors             string  `json:"Vectors"`
		} `json:"CVSSv3"`
	} `json:"NVD"`
}

type VulnerabilityHostProps struct {
	FirstTimeSeen   *time.Time `json:"first_time_seen,omitempty"`
	IsDailyJob      int        `json:"isDailyJob,omitempty"`
	LastUpdatedTime *time.Time `json:"last_updated_time,omitempty"`
}

func (v *VulnerabilityHost) HasFix() bool {
	return v.FixInfo.FixAvailable == "1"
}

func (hosts *VulnerabilitiesHostResponse) VulnerabilityCounts() HostVulnCounts {
	var hostCounts = HostVulnCounts{}

	// remove duplicates before count.
	for _, h := range hosts.Data {
		switch h.Severity {
		case "Critical":
			hostCounts.Critical++
			hostCounts.Total++
			if h.HasFix() {
				hostCounts.CritFixable++
				hostCounts.TotalFixable++
			}
		case "High":
			hostCounts.High++
			hostCounts.Total++
			if h.HasFix() {
				hostCounts.HighFixable++
				hostCounts.TotalFixable++
			}
		case "Medium":
			hostCounts.Medium++
			hostCounts.Total++
			if h.HasFix() {
				hostCounts.MedFixable++
				hostCounts.TotalFixable++
			}
		case "Low":
			hostCounts.Low++
			hostCounts.Total++
			if h.HasFix() {
				hostCounts.LowFixable++
				hostCounts.TotalFixable++
			}
		default:
			hostCounts.Info++
			hostCounts.Total++
			if h.HasFix() {
				hostCounts.InfoFixable++
				hostCounts.TotalFixable++
			}
		}
	}

	return hostCounts
}

type VulnCveSummary struct {
	Host      VulnerabilityHost
	Count     int
	Hostnames []string
}
