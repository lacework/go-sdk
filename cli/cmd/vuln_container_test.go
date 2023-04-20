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

package cmd

import (
	"encoding/json"
	"regexp"
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/capturer"
)

func TestUserFriendlyErrorFromOnDemandCtrVulnScanRepositoryNotFound(t *testing.T) {
	err := userFriendlyErrorForOnDemandCtrVulnScan(
		errors.New("Could not successfully send scan request to available integrations for given repo and label"),
		"my-registry.example.com", "image", "tag",
	)
	if assert.NotNil(t, err) {
		assert.Contains(t,
			err.Error(),
			"container image 'image@tag' not found in registry 'my-registry.example.com'")
		assert.Contains(t,
			err.Error(),
			"To view all container registries configured in your account use the command:")
		assert.Contains(t,
			err.Error(),
			"lacework vulnerability container list-registries")
	}
}

func TestBuildCSVVulnCtrReportVulnerabilitiesListing(t *testing.T) {
	cli.EnableCSVOutput()
	defer func() { cli.csvOutput = false }()

	var data api.VulnerabilitiesContainersResponse
	if err := json.Unmarshal([]byte(rawListAssessments), &data); err != nil {
		panic(err)
	}

	headers := []string{"Registry", "Repository", "Last Scan", "Status", "Containers", "Vulnerabilities", "Image Digest"}
	assessments := buildVulnCtrAssessmentSummary(data.Data, api.ContainersEntityResponse{
		Data: []api.ContainerEntity{
			api.ContainerEntity{
				ImageID: "sha256:7652596622b05043763f962cff30edf01f6ea1ba29374f1703dda759dc9ff3a1",
				Mid:     1,
			},
			api.ContainerEntity{
				ImageID: "sha256:7652596622b05043763f962cff30edf01f6ea1ba29374f1703dda759dc9ff3a1",
				Mid:     2,
			},
			api.ContainerEntity{
				ImageID: "sha256:1252596622b05043763f962gff30adf01f6ea1ba29374f1703dda759dc9ab3a1",
				Mid:     3,
			},
		},
	})
	filteredAssessments := applyVulnCtrFilters(assessments)
	assessmentOutput := assessmentSummaryToOutputFormat(filteredAssessments)
	rows := vulAssessmentsToTable(assessmentOutput)
	csv, err := renderAsCSV(headers, rows)
	if err != nil {
		panic(err)
	}

	expected := `
Registry,Repository,Last Scan,Status,Containers,Vulnerabilities,Image Digest
gcr.io,techally-test-2/exservice,2022-11-21T19:21:57Z,Success,2,1 Critical,sha256:12b072fd2ce1732e4c2f0f601c2c12ea2ea657c9572d9ba477b1174d9159e123
gcr.io,techally-test-4/exservice,2022-11-21T19:21:57Z,Success,1,1 High 1 Fixable,sha256:15b072fd2ce1732e4c2f0f601c2c12ea2ea657c9572d9ba477b1174d9159e123
index.docker.io,techally-test/test-cli,2022-11-21T18:33:28Z,Success,0,1 Medium 1 Fixable,sha256:77b2d2246518044ef95e3dbd029e51dd477788e5bf8e278e418685aabc3fe28a
`

	assert.Equal(t, strings.TrimPrefix(expected, "\n"), csv)
}

func TestBuildCSVVulnCtrReportWithVulnerabilities(t *testing.T) {
	cli.EnableCSVOutput()
	vulCmdState.Details = true
	defer func() {
		cli.csvOutput = false
		vulCmdState.Details = false
	}()

	var response api.VulnerabilitiesContainersResponse
	if err := json.Unmarshal([]byte(rawListAssessments), &response); err != nil {
		panic(err)
	}
	cliOutput := capturer.CaptureOutput(func() {
		assert.Nil(t, buildVulnContainerAssessmentReports(response))
	})

	expected := `CVE ID,Severity,CVSSv2,CVSSv3,Package,Current Version,Fix Version,Version Format,Feed,Src,Start Time,Status,Namespace,Image Digest,Image ID,Image Repo,Image Registry,Image Size,Introduced in Layer
CVE-2020-12345,Critical,0.0,0.0,example-2,1.2.0,,apk,n/a,,2022-11-21T19:21:57Z,alpine:v3.11,sha256:7652596622b05043763f962cff30edf01f6ea1ba29374f1703dda759dc9ff3a1,sha256:12b072fd2ce1732e4c2f0f601c2c12ea2ea657c9572d9ba477b1174d9159e123,techally-test-2/exservice,gcr.io,14933503,apk add --no-cache ca-certificates
CVE-2020-12345,High,0.0,0.0,example-4,1.0.0,1.31.1-r11,apk,lacework,,2022-11-21T19:21:57Z,alpine:v3.11,sha256:1252596622b05043763f962gff30adf01f6ea1ba29374f1703dda759dc9ab3a1,sha256:15b072fd2ce1732e4c2f0f601c2c12ea2ea657c9572d9ba477b1174d9159e123,techally-test-4/exservice,gcr.io,14933503,apk add --no-cache ca-certificates
CVE-2029-21234,Medium,0.0,0.0,example-1,1.0.0,2.2.0-11+deb9u4,dpkg,lacework,var/lib/dpkg/status,2022-11-21T18:33:28Z,debian:9,sha256:a65572164cb78c4d04f57bd66201c775e2dab08fce394806a03a933c5daf9e48,sha256:77b2d2246518044ef95e3dbd029e51dd477788e5bf8e278e418685aabc3fe28a,techally-test/test-cli,index.docker.io,360608563,example introduced in layer
`
	assert.Equal(t, strings.TrimPrefix(expected, "\n"), cliOutput)
}

func TestBuildVulnCtrReportWithIntroducedInLayerCSV(t *testing.T) {
	cli.EnableCSVOutput()
	vulCmdState.Details = true
	defer func() {
		cli.csvOutput = false
		vulCmdState.Details = false
	}()

	var response api.VulnerabilitiesContainersResponse
	if err := json.Unmarshal([]byte(mockIntroducedInLayerResponse), &response); err != nil {
		panic(err)
	}
	cliOutput := capturer.CaptureOutput(func() {
		assert.Nil(t, buildVulnContainerAssessmentReports(response))
	})

	expected := `CVE ID,Severity,CVSSv2,CVSSv3,Package,Current Version,Fix Version,Version Format,Feed,Src,Start Time,Status,Namespace,Image Digest,Image ID,Image Repo,Image Registry,Image Size,Introduced in Layer
CVE-2029-21234,Medium,0.0,0.0,example-1,1.0.0,2.2.0-11+deb9u4,dpkg,lacework,var/lib/dpkg/status,2022-11-21T18:33:28Z,debian:9,sha256:a65572164cb78c4d04f57bd66201c775e2dab08fce394806a03a933c5daf9e48,sha256:77b2d2246518044ef95e3dbd029e51dd477788e5bf8e278e418685aabc3fe28a,techally-test/test-cli,index.docker.io,360608563,example introduced in layer 1
CVE-2029-21234,Medium,0.0,0.0,example-1,1.0.0,2.2.0-11+deb9u4,dpkg,lacework,var/lib/dpkg/status,2022-11-21T18:33:28Z,debian:9,sha256:a65572164cb78c4d04f57bd66201c775e2dab08fce394806a03a933c5daf9e48,sha256:77b2d2246518044ef95e3dbd029e51dd477788e5bf8e278e418685aabc3fe28a,techally-test/test-cli,index.docker.io,360608563,example introduced in layer 2
`
	assert.Equal(t, strings.TrimPrefix(expected, "\n"), cliOutput)
}

func TestBuildVulnCtrReportWithAggregatedIntroducedInLayer(t *testing.T) {
	vulCmdState.Details = true
	defer func() {
		vulCmdState.Details = false
	}()

	var response api.VulnerabilitiesContainersResponse
	if err := json.Unmarshal([]byte(mockIntroducedInLayerResponse), &response); err != nil {
		panic(err)
	}
	cliOutput := capturer.CaptureOutput(func() {
		assert.Nil(t, buildVulnContainerAssessmentReports(response))
	})

	assert.Contains(t, cliOutput, "introduced in 2 layers...")
}

func TestBuildVulnCtrReportAndJsonCount(t *testing.T) {
	cli.EnableCSVOutput()
	vulCmdState.Details = true
	defer func() {
		vulCmdState.Details = false
		cli.jsonOutput = false
		cli.csvOutput = false
	}()

	var response api.VulnerabilitiesContainersResponse
	if err := json.Unmarshal([]byte(mockIntroducedInLayerResponse), &response); err != nil {
		panic(err)
	}

	cliCSVOutput := capturer.CaptureOutput(func() {
		assert.Nil(t, buildVulnContainerAssessmentReports(response))
	})

	cli.csvOutput = false
	cli.EnableJSONOutput()

	cliJsonOutput := capturer.CaptureOutput(func() {
		assert.Nil(t, buildVulnContainerAssessmentReports(response))
	})

	jsonCount := len(strings.Split(cliJsonOutput, "CVE-"))
	csvCount := len(strings.Split(cliCSVOutput, "CVE-"))

	assert.Equal(t, csvCount, jsonCount)
}

func TestVulnCtrIntroducedInRegex(t *testing.T) {
	tabsRegexTableTest := []struct {
		Name     string
		Input    string
		Expected string
	}{
		{
			Name:     "single tab",
			Input:    "My\tString",
			Expected: "My\nString",
		},
		{
			Name:     "many tabs",
			Input:    "My\t\tString\t",
			Expected: "My\nString\n",
		},
		{
			Name:     "many tabs",
			Input:    "This\t\t\t\tis\t\t\tan\t\tExample\tString\t\t\t",
			Expected: "This\nis\nan\nExample\nString\n",
		},
	}

	regex := regexp.MustCompile(regexAllTabs)

	for _, test := range tabsRegexTableTest {
		t.Run(test.Name, func(t *testing.T) {
			result := regex.ReplaceAllString(test.Input, "\n")
			assert.Contains(t, result, test.Expected)
		})
	}
}

func TestVulnCtrCountPackages(t *testing.T) {
	var (
		response api.VulnerabilitiesContainersResponse
		expected = 3
	)
	if err := json.Unmarshal([]byte(rawListAssessments), &response); err != nil {
		panic(err)
	}

	totalPackages := countVulnContainerImagePackages(response.Data)
	assert.Equal(t, totalPackages, expected)
}

var rawListAssessments = `
{
    "paging": {
        "rows": 5000,
        "totalRows": 6419,
        "urls": {
            "nextPage": "https://example.lacework.net/api/v2/Vulnerabilities/Containers/"
        }
    },
"data": [
{
            "evalCtx": {
                "cve_batch_info": [
                    {
                        "cve_created_time": "2022-11-21 00:21:41.678000000"
                    }
                ],
                "exception_props": [
                    {
                        "exception_guid": "VULN_C44BF2CBE09F0E705565BEA1A0C1D2A5D1534857F2C7CDF8381",
                        "exception_name": "registry index.docker.io severity Low",
                        "exception_reason": "Accepted Risk"
                    }
                ],
                "image_info": {
                    "created_time": 1605140985874,
                    "digest": "sha256:77b2d2246518044ef95e3dbd029e51dd477788e5bf8e278e418685aabc3fe28a",
                    "id": "sha256:a65572164cb78c4d04f57bd66201c775e2dab08fce394806a03a933c5daf9e48",
                    "registry": "index.docker.io",
                    "repo": "techally-test/test-cli",
                    "scan_created_time": 1669055600,
                    "size": 360608563,
                    "status": "Success",
                    "tags": [
                        "latest"
                    ],
                    "type": "Docker"
                },
                "integration_props": {
                    "INTG_GUID": "TECHALLY_FC5485B5ACFF3DAFE77E8C8A734C6C2FAD7CAAC9F01313C",
                    "NAME": "Terraform-Dockerhub",
                    "REGISTRY_TYPE": "DOCKERHUB"
                },
                "is_reeval": false,
                "request_source": "PLATFORM_SCANNER",
                "scan_batch_id": "467a274c-f847-456b-b62d-13f9d88988cc-1669055607923432004",
                "scan_request_props": {
                    "data_format_version": "1.0",
                    "props": {
                        "data_format_version": "1.0",
                        "scanner_version": "10.0.155"
                    },
                    "reqId": "2ac494a9-b7be-453a-81b9-7a2f1f9e2113",
                    "reqSource": "ondemand",
                    "scanCompletionUtcTime": 1669055607,
                    "scan_start_time": 1669055600,
                    "scanner_version": "10.0.155"
                },
                "vuln_batch_id": "7B2EDDD2D2D140ECA6B85001FC62AE45",
                "vuln_created_time": "2022-11-21 00:21:41.678000000"
            },
            "evalGuid": "781865fdff984def2587b5f05065f0db",
            "featureKey": {
                "name": "example-1",
                "namespace": "debian:9",
                "version": "1.0.0"
            },
            "featureProps": {
                "feed": "lacework",
                "introduced_in": "example introduced in layer",
                "layer": "sha256:sha256:572866ab72a68759e23b071fbbdce6341137c9606936b4fff9846f74997bbaac",
                "src": "var/lib/dpkg/status",
                "version_format": "dpkg"
            },
            "fixInfo": {
                "fix_available": 1,
                "fixed_version": "2.2.0-11+deb9u4"
            },
            "imageId": "sha256:a65572164cb78c4d04f57bd66201c775e2dab08fce394806a03a933c5daf9e48",
            "severity": "Medium",
            "startTime": "2022-11-21T18:33:28.076Z",
            "status": "VULNERABLE",
            "vulnId": "CVE-2029-21234"
        },
{
            "evalCtx": {
                "cve_batch_info": [
                    {
                        "cve_created_time": "2022-11-21 19:05:48.075000000"
                    }
                ],
                "image_info": {
                    "created_time": 1588284823675,
                    "digest": "sha256:12b072fd2ce1732e4c2f0f601c2c12ea2ea657c9572d9ba477b1174d9159e123",
                    "id": "sha256:7652596622b05043763f962cff30edf01f6ea1ba29374f1703dda759dc9ff3a1",
                    "registry": "gcr.io",
                    "repo": "techally-test-2/exservice",
                    "scan_created_time": 1636768856,
                    "size": 14933503,
                    "status": "Success",
                    "tags": [
                        "v1.0.0"
                    ],
                    "type": "Docker"
                },
                "integration_props": {},
                "is_reeval": true,
                "request_source": "PLATFORM_SCANNER",
                "scan_batch_id": "cd1d57ca-c018-4ffd-ac07-6664bc7c7a85-1636768857524097264",
                "scan_request_props": {
                    "data_format_version": "1.0",
                    "props": {
                        "data_format_version": "1.0",
                        "scanner_version": "0.2.2"
                    },
                    "scanCompletionUtcTime": 1636768857,
                    "scan_start_time": 1636768856,
                    "scanner_version": "0.2.2"
                },
                "vuln_batch_id": "E1BA1053AB374E4C968C689F0F013C9A",
                "vuln_created_time": "2022-11-21 19:05:48.075000000"
            },
            "evalGuid": "097464827bb2d34b6f62c5ebbdab3385",
            "featureKey": {
                "name": "example-2",
                "namespace": "alpine:v3.11",
                "version": "1.2.0"
            },
            "featureProps": {
                "feed": "n/a",
                "introduced_in": "apk add --no-cache ca-certificates",
                "layer": "sha256:sha256:e3693d3358098cb60aed2d9820d583add96dec7313befcf714ffc4d9c464a275",
                "src": "",
                "version_format": "apk"
            },
            "fixInfo": {
                "fix_available": 0,
                "fixed_version": ""
            },
            "imageId": "sha256:7652596622b05043763f962cff30edf01f6ea1ba29374f1703dda759dc9ff3a1",
            "startTime": "2022-11-21T19:21:57.765Z",
            "status": "VULNERABLE",
            "severity": "Critical",
            "vulnId": "CVE-2020-12345"
        },
        {
            "evalCtx": {
                "cve_batch_info": [
                    {
                        "cve_created_time": "2022-11-21 19:05:48.075000000"
                    }
                ],
                "image_info": {
                    "created_time": 1588284823675,
                    "digest": "sha256:15b072fd2ce1732e4c2f0f601c2c12ea2ea657c9572d9ba477b1174d9159e123",
                    "id": "sha256:1252596622b05043763f962gff30adf01f6ea1ba29374f1703dda759dc9ab3a1",
                    "registry": "gcr.io",
                    "repo": "techally-test-4/exservice",
                    "scan_created_time": 1636768856,
                    "size": 14933503,
                    "status": "Success",
                    "tags": [
                        "v1.0.0"
                    ],
                    "type": "Docker"
                },
                "integration_props": {},
                "is_reeval": true,
                "request_source": "PLATFORM_SCANNER",
                "scan_batch_id": "cd1d57ca-c018-4ffd-ac07-6664bc7c7a85-1636768857524097264",
                "scan_request_props": {
                    "data_format_version": "1.0",
                    "props": {
                        "data_format_version": "1.0",
                        "scanner_version": "0.2.2"
                    },
                    "scanCompletionUtcTime": 1636768857,
                    "scan_start_time": 1636768856,
                    "scanner_version": "0.2.2"
                },
                "vuln_batch_id": "E1BA1053AB374E4C968C689F0F013C9A",
                "vuln_created_time": "2022-11-21 19:05:48.075000000"
            },
            "evalGuid": "097464827bb2d34b6f62c5ebbdab3385",
            "featureKey": {
                "name": "example-4",
                "namespace": "alpine:v3.11",
                "version": "1.0.0"
            },
            "featureProps": {
                "feed": "lacework",
                "introduced_in": "apk add --no-cache ca-certificates",
                "layer": "sha256:sha256:e3693d3358098cb60aed2d9820d583add96dec7313befcf714ffc4d9c464a275",
                "src": "",
                "version_format": "apk"
            },
            "fixInfo": {
                "fix_available": 1,
                "fixed_version": "1.31.1-r11"
            },
            "imageId": "sha256:1252596622b05043763f962gff30adf01f6ea1ba29374f1703dda759dc9ab3a1",
            "severity": "High",
            "startTime": "2022-11-21T19:21:57.765Z",
            "status": "VULNERABLE",
            "vulnId": "CVE-2020-12345"
        }
]
}`

var mockIntroducedInLayerResponse = `
{
    "paging": {
        "rows": 5000,
        "totalRows": 6419,
        "urls": {
            "nextPage": "https://example.lacework.net/api/v2/Vulnerabilities/Containers/"
        }
    },
"data": [
{
            "evalCtx": {
                "cve_batch_info": [
                    {
                        "cve_created_time": "2022-11-21 00:21:41.678000000"
                    }
                ],
                "exception_props": [
                    {
                        "exception_guid": "VULN_C44BF2CBE09F0E705565BEA1A0C1D2A5D1534857F2C7CDF8381",
                        "exception_name": "registry index.docker.io severity Low",
                        "exception_reason": "Accepted Risk"
                    }
                ],
                "image_info": {
                    "created_time": 1605140985874,
                    "digest": "sha256:77b2d2246518044ef95e3dbd029e51dd477788e5bf8e278e418685aabc3fe28a",
                    "id": "sha256:a65572164cb78c4d04f57bd66201c775e2dab08fce394806a03a933c5daf9e48",
                    "registry": "index.docker.io",
                    "repo": "techally-test/test-cli",
                    "scan_created_time": 1669055600,
                    "size": 360608563,
                    "status": "Success",
                    "tags": [
                        "latest"
                    ],
                    "type": "Docker"
                },
                "integration_props": {
                    "INTG_GUID": "TECHALLY_FC5485B5ACFF3DAFE77E8C8A734C6C2FAD7CAAC9F01313C",
                    "NAME": "Terraform-Dockerhub",
                    "REGISTRY_TYPE": "DOCKERHUB"
                },
                "is_reeval": false,
                "request_source": "PLATFORM_SCANNER",
                "scan_batch_id": "467a274c-f847-456b-b62d-13f9d88988cc-1669055607923432004",
                "scan_request_props": {
                    "data_format_version": "1.0",
                    "props": {
                        "data_format_version": "1.0",
                        "scanner_version": "10.0.155"
                    },
                    "reqId": "2ac494a9-b7be-453a-81b9-7a2f1f9e2113",
                    "reqSource": "ondemand",
                    "scanCompletionUtcTime": 1669055607,
                    "scan_start_time": 1669055600,
                    "scanner_version": "10.0.155"
                },
                "vuln_batch_id": "7B2EDDD2D2D140ECA6B85001FC62AE45",
                "vuln_created_time": "2022-11-21 00:21:41.678000000"
            },
            "evalGuid": "781865fdff984def2587b5f05065f0db",
            "featureKey": {
                "name": "example-1",
                "namespace": "debian:9",
                "version": "1.0.0"
            },
            "featureProps": {
                "feed": "lacework",
                "introduced_in": "example introduced in layer 1",
                "layer": "sha256:sha256:572866ab72a68759e23b071fbbdce6341137c9606936b4fff9846f74997bbaac",
                "src": "var/lib/dpkg/status",
                "version_format": "dpkg"
            },
            "fixInfo": {
                "fix_available": 1,
                "fixed_version": "2.2.0-11+deb9u4"
            },
            "imageId": "sha256:a65572164cb78c4d04f57bd66201c775e2dab08fce394806a03a933c5daf9e48",
            "severity": "Medium",
            "startTime": "2022-11-21T18:33:28.076Z",
            "status": "VULNERABLE",
            "vulnId": "CVE-2029-21234"
        },{
            "evalCtx": {
                "cve_batch_info": [
                    {
                        "cve_created_time": "2022-11-21 00:21:41.678000000"
                    }
                ],
                "exception_props": [
                    {
                        "exception_guid": "VULN_C44BF2CBE09F0E705565BEA1A0C1D2A5D1534857F2C7CDF8381",
                        "exception_name": "registry index.docker.io severity Low",
                        "exception_reason": "Accepted Risk"
                    }
                ],
                "image_info": {
                    "created_time": 1605140985874,
                    "digest": "sha256:77b2d2246518044ef95e3dbd029e51dd477788e5bf8e278e418685aabc3fe28a",
                    "id": "sha256:a65572164cb78c4d04f57bd66201c775e2dab08fce394806a03a933c5daf9e48",
                    "registry": "index.docker.io",
                    "repo": "techally-test/test-cli",
                    "scan_created_time": 1669055600,
                    "size": 360608563,
                    "status": "Success",
                    "tags": [
                        "latest"
                    ],
                    "type": "Docker"
                },
                "integration_props": {
                    "INTG_GUID": "TECHALLY_FC5485B5ACFF3DAFE77E8C8A734C6C2FAD7CAAC9F01313C",
                    "NAME": "Terraform-Dockerhub",
                    "REGISTRY_TYPE": "DOCKERHUB"
                },
                "is_reeval": false,
                "request_source": "PLATFORM_SCANNER",
                "scan_batch_id": "467a274c-f847-456b-b62d-13f9d88988cc-1669055607923432004",
                "scan_request_props": {
                    "data_format_version": "1.0",
                    "props": {
                        "data_format_version": "1.0",
                        "scanner_version": "10.0.155"
                    },
                    "reqId": "2ac494a9-b7be-453a-81b9-7a2f1f9e2113",
                    "reqSource": "ondemand",
                    "scanCompletionUtcTime": 1669055607,
                    "scan_start_time": 1669055600,
                    "scanner_version": "10.0.155"
                },
                "vuln_batch_id": "7B2EDDD2D2D140ECA6B85001FC62AE45",
                "vuln_created_time": "2022-11-21 00:21:41.678000000"
            },
            "evalGuid": "781865fdff984def2587b5f05065f0db",
            "featureKey": {
                "name": "example-1",
                "namespace": "debian:9",
                "version": "1.0.0"
            },
            "featureProps": {
                "feed": "lacework",
                "introduced_in": "example introduced in layer 2",
                "layer": "sha256:sha256:572866ab72a68759e23b071fbbdce6341137c9606936b4fff9846f74997bbaac",
                "src": "var/lib/dpkg/status",
                "version_format": "dpkg"
            },
            "fixInfo": {
                "fix_available": 1,
                "fixed_version": "2.2.0-11+deb9u4"
            },
            "imageId": "sha256:a65572164cb78c4d04f57bd66201c775e2dab08fce394806a03a933c5daf9e48",
            "severity": "Medium",
            "startTime": "2022-11-21T18:33:28.076Z",
            "status": "VULNERABLE",
            "vulnId": "CVE-2029-21234"
        }

]
}`
