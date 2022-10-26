//
// Author:: Darren Murray (<dmurray-lacework@lacework.net>)
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
	"log"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/capturer"
)

func TestListCvesFilterSeverity(t *testing.T) {
	vulCmdState.Severity = "critical"
	defer clearVulnFilters()

	mockCves := map[string]VulnCveSummary{"TestVulnOne": mockCveOne, "TestVulnTwo": mockCveTwo}
	result, output := filterHostCVEsTable(mockCves)

	assert.Equal(t, len(result), 1)
	assert.Equal(t, output, "\n1 of 2 cve(s) showing\n")
}

func TestShowAssessmentFilterSeverityWithPackages(t *testing.T) {
	vulCmdState.Severity = "critical"
	vulCmdState.Packages = true
	defer clearVulnFilters()

	mockCves := map[string]VulnCveSummary{"TestVulnOne": mockCveOne, "TestVulnTwo": mockCveTwo}
	result, output := hostVulnPackagesTable(mockCves, true)

	assert.Equal(t, len(result), 1)
	assert.Equal(t, output, "1 of 2 package(s) showing\n")
}

func clearVulnFilters() {
	vulCmdState.Severity = ""
	vulCmdState.Packages = false
	vulCmdState.Active = false
}

var mockCveOne = VulnCveSummary{
	Host: api.VulnerabilityHost{
		Props:     api.VulnerabilityHostProps{},
		Mid:       1,
		Severity:  "High",
		StartTime: time.Time{},
		Status:    "",
		VulnID:    "TestVulnOne",
	},
}

var mockCveTwo = VulnCveSummary{
	Host: api.VulnerabilityHost{
		Props:     api.VulnerabilityHostProps{},
		Mid:       1,
		Severity:  "Critical",
		StartTime: time.Time{},
		Status:    "",
		VulnID:    "TestVulnTwo",
	},
}

func TestBuildVulnHostReportsNoVulnerabilities(t *testing.T) {
	cliOutput := capturer.CaptureOutput(func() {
		assert.Nil(t, buildVulnHostReports(api.VulnerabilitiesHostResponse{}))
	})
	assert.Contains(t, cliOutput, "Great news! This host has no vulnerabilities...")

	t.Run("test JSON output", func(t *testing.T) {
		cli.EnableJSONOutput()
		defer cli.EnableHumanOutput()
		cliJSONOutput := capturer.CaptureOutput(func() {
			assert.Nil(t, buildVulnHostReports(api.VulnerabilitiesHostResponse{}))
		})
		expectedJSON := `null
`
		assert.Equal(t, expectedJSON, cliJSONOutput)
	})
}

func TestBuildVulnHostReportsWithVulnerabilitiesSummaryOnlyAndNoFilters(t *testing.T) {
	cliOutput := capturer.CaptureOutput(func() {
		assert.Nil(t, buildVulnHostReports(mockHostVulnerabilityAssessment()))
	})
	// NOTE (@afiune): We purposely leave trailing spaces in this table, we need them!
	expectedTable := `
          HOST DETAILS                   VULNERABILITIES          
--------------------------------+---------------------------------
    Machine ID   51                 SEVERITY   COUNT   FIXABLE    
    Hostname     host-mock        -----------+-------+----------  
    External IP  mock               Critical       0         0    
    Internal IP  mock               High           1         0    
    Os           linux              Medium         3         2    
    Arch         arm64              Low            4         2    
    Namespace    ubuntu:18.04       Info          14         0    
    Provider     AWS                                              
    Instance ID  i-mock                                           
    AMI          ami-mock                                         
                                                                  

Try adding '--details' to increase details shown about the vulnerability assessment.
`
	assert.Equal(t, strings.TrimPrefix(expectedTable, "\n"), cliOutput)
}

func TestBuildVulnHostReportsWithVulnerabilitiesPackagesViewAndNoFilters(t *testing.T) {
	vulCmdState.Packages = true
	defer clearVulnFilters()

	cliOutput := capturer.CaptureOutput(func() {
		assert.Nil(t, buildVulnHostReports(mockHostVulnerabilityAssessment()))
	})
	// NOTE (@afiune): We purposly leave trailing spaces in this table, we need them!
	expectedTable := `
          HOST DETAILS                   VULNERABILITIES          
--------------------------------+---------------------------------
    Machine ID   51                 SEVERITY   COUNT   FIXABLE    
    Hostname     host-mock        -----------+-------+----------  
    External IP  mock               Critical       0         0    
    Internal IP  mock               High           1         0    
    Os           linux              Medium         3         2    
    Arch         arm64              Low            4         2    
    Namespace    ubuntu:18.04       Info          14         0    
    Provider     AWS                                              
    Instance ID  i-mock                                           
    AMI          ami-mock                                         
                                                                  
  CVE COUNT   SEVERITY     PACKAGE            CURRENT VERSION                  FIX VERSION            PKG STATUS  
------------+----------+--------------+------------------------------+------------------------------+-------------
  1           High       ubuntu:18.04   1:2017.3.23-2ubuntu0.18.04.2                                              
  1           Medium     ubuntu:18.04   1:2017.3.23-2ubuntu0.18.04.2                                              
  1           Medium     ubuntu:18.04   1:2017.3.23-2ubuntu0.18.04.2   7.58.0-2ubuntu3.18                         
  1           Medium     ubuntu:18.04   1:2017.3.23-2ubuntu0.18.04.2   1:2017.3.23-2ubuntu0.18.04.4               
  1           Low        ubuntu:18.04   1:2017.3.23-2ubuntu0.18.04.2   4.4.18-2ubuntu1.3                          
  2           Low        ubuntu:18.04   1:2017.3.23-2ubuntu0.18.04.2                                              
  1           Low        ubuntu:18.04   1:2017.3.23-2ubuntu0.18.04.2   3.6.9-1~18.04ubuntu1.8                     

Try adding '--active' to only show vulnerabilities of packages actively running.
`
	assert.Equal(t, strings.TrimPrefix(expectedTable, "\n"), cliOutput)
}

func TestBuildVulnHostReportsWithVulnerabilitiesPackagesViewWithFilters(t *testing.T) {
	vulCmdState.Packages = true
	vulCmdState.Severity = "high"
	defer clearVulnFilters()

	cliOutput := capturer.CaptureOutput(func() {
		assert.Nil(t, buildVulnHostReports(mockHostVulnerabilityAssessment()))
	})
	// NOTE (@afiune): We purposly leave trailing spaces in this table, we need them!
	expectedTable := `
          HOST DETAILS                   VULNERABILITIES          
--------------------------------+---------------------------------
    Machine ID   51                 SEVERITY   COUNT   FIXABLE    
    Hostname     host-mock        -----------+-------+----------  
    External IP  mock               Critical       0         0    
    Internal IP  mock               High           1         0    
    Os           linux              Medium         3         2    
    Arch         arm64              Low            4         2    
    Namespace    ubuntu:18.04       Info          14         0    
    Provider     AWS                                              
    Instance ID  i-mock                                           
    AMI          ami-mock                                         
                                                                  
  CVE COUNT   SEVERITY     PACKAGE            CURRENT VERSION          FIX VERSION   PKG STATUS  
------------+----------+--------------+------------------------------+-------------+-------------
  1           High       ubuntu:18.04   1:2017.3.23-2ubuntu0.18.04.2                             

Try adding '--active' to only show vulnerabilities of packages actively running.

1 of 8 cve(s) showing
`
	assert.Equal(t, strings.TrimPrefix(expectedTable, "\n"), cliOutput)
}

func TestBuildVulnHostReportsWithVulnerabilitiesWithFiltersSeverity(t *testing.T) {
	vulCmdState.Severity = "high"
	defer clearVulnFilters()

	cliOutput := capturer.CaptureOutput(func() {
		assert.Nil(t, buildVulnHostReports(mockHostVulnerabilityAssessment()))
	})
	// NOTE (@afiune): We purposly leave trailing spaces in this table, we need them!
	expectedTable := `
          HOST DETAILS                   VULNERABILITIES          
--------------------------------+---------------------------------
    Machine ID   51                 SEVERITY   COUNT   FIXABLE    
    Hostname     host-mock        -----------+-------+----------  
    External IP  mock               Critical       0         0    
    Internal IP  mock               High           1         0    
    Os           linux              Medium         3         2    
    Arch         arm64              Low            4         2    
    Namespace    ubuntu:18.04       Info          14         0    
    Provider     AWS                                              
    Instance ID  i-mock                                           
    AMI          ami-mock                                         
                                                                  
      CVE ID       SEVERITY   CVSSV2   CVSSV3    PACKAGE          CURRENT VERSION          FIX VERSION   PKG STATUS   VULN STATUS  
-----------------+----------+--------+--------+-----------+------------------------------+-------------+------------+--------------
  CVE-2022-33741   High       3.6      7.1      linux-aws   1:2017.3.23-2ubuntu0.18.04.2                              Reopened     

Try adding '--active' to only show vulnerabilities of packages actively running.

1 of 8 cve(s) showing
`
	assert.Equal(t, strings.TrimPrefix(expectedTable, "\n"), cliOutput)
}

func TestBuildVulnHostReportsWithVulnerabilitiesWithFiltersActive(t *testing.T) {
	vulCmdState.Active = true
	defer clearVulnFilters()

	cliOutput := capturer.CaptureOutput(func() {
		assert.Nil(t, buildVulnHostReports(mockHostVulnerabilityAssessment()))
	})
	// NOTE (@afiune): We purposely leave trailing spaces in this table, we need them!
	expectedTable := `
          HOST DETAILS                   VULNERABILITIES          
--------------------------------+---------------------------------
    Machine ID   51                 SEVERITY   COUNT   FIXABLE    
    Hostname     host-mock        -----------+-------+----------  
    External IP  mock               Critical       0         0    
    Internal IP  mock               High           1         0    
    Os           linux              Medium         3         2    
    Arch         arm64              Low            4         2    
    Namespace    ubuntu:18.04       Info          14         0    
    Provider     AWS                                              
    Instance ID  i-mock                                           
    AMI          ami-mock                                         
                                                                  
      CVE ID       SEVERITY   CVSSV2   CVSSV3    PACKAGE           CURRENT VERSION                  FIX VERSION            PKG STATUS   VULN STATUS  
-----------------+----------+--------+--------+------------+------------------------------+------------------------------+------------+--------------
  CVE-2022-30787   Medium     4.6      6.7      ntfs-3g      1:2017.3.23-2ubuntu0.18.04.2   1:2017.3.23-2ubuntu0.18.04.4                Active       
  CVE-2022-0351    Medium     4.6      7.8      vim          1:2017.3.23-2ubuntu0.18.04.2                                               Active       
  CVE-2022-27782   Medium     5.0      7.5      curl         1:2017.3.23-2ubuntu0.18.04.2   7.58.0-2ubuntu3.18                          Active       
  CVE-2019-18276   Low        7.2      7.8      bash         1:2017.3.23-2ubuntu0.18.04.2   4.4.18-2ubuntu1.3                           Active       
  CVE-2020-13988   Low        5.0      7.5      open-iscsi   1:2017.3.23-2ubuntu0.18.04.2                                               Active       
  CVE-2022-2129    Low        6.8      7.8      vim          1:2017.3.23-2ubuntu0.18.04.2                                               Active       
  CVE-2015-20107   Low        10.0     9.8      python3.6    1:2017.3.23-2ubuntu0.18.04.2   3.6.9-1~18.04ubuntu1.8                      Active       

Try adding '--fixable' to only show fixable vulnerabilities.

7 of 8 cve(s) showing
`
	assert.Equal(t, strings.TrimPrefix(expectedTable, "\n"), cliOutput)
}

func TestBuildCSVVulnHostReportsWithVulnerabilities(t *testing.T) {
	cli.EnableCSVOutput()
	vulCmdState.Details = true

	defer func() {
		cli.csvOutput = false
		vulCmdState.Details = false
	}()

	cliOutput := capturer.CaptureOutput(func() {
		assert.Nil(t, buildVulnHostReports(mockHostVulnerabilityAssessment()))
	})
	assert.Equal(t, strings.TrimPrefix(expectedCSVHostDetailsTable, "\n"), cliOutput)
}

func TestBuildVulnHostReportsWithVulnerabilitiesWithFiltersSeverityAndActiveReturnsNoVulns(t *testing.T) {
	vulCmdState.Severity = "high"
	vulCmdState.Active = true
	defer clearVulnFilters()

	cliOutput := capturer.CaptureOutput(func() {
		assert.Nil(t, buildVulnHostReports(mockHostVulnerabilityAssessment()))
	})
	// NOTE (@afiune): We purposly leave trailing spaces in this table, we need them!
	expectedTable := `
          HOST DETAILS                   VULNERABILITIES          
--------------------------------+---------------------------------
    Machine ID   51                 SEVERITY   COUNT   FIXABLE    
    Hostname     host-mock        -----------+-------+----------  
    External IP  mock               Critical       0         0    
    Internal IP  mock               High           1         0    
    Os           linux              Medium         3         2    
    Arch         arm64              Low            4         2    
    Namespace    ubuntu:18.04       Info          14         0    
    Provider     AWS                                              
    Instance ID  i-mock                                           
    AMI          ami-mock                                         
                                                                  
There are no high vulnerabilities of packages actively running in your environment.

Try adding '--fixable' to only show fixable vulnerabilities.

0 of 8 cve(s) showing
`
	assert.Equal(t, strings.TrimPrefix(expectedTable, "\n"), cliOutput)
}

func TestBuildVulnHostReportsWithVulnerabilitiesWithFilterReturnsNoVulns(t *testing.T) {
	vulCmdState.Severity = "critical"
	defer clearVulnFilters()

	cliOutput := capturer.CaptureOutput(func() {
		assert.Nil(t, buildVulnHostReports(mockHostVulnerabilityAssessment()))
	})
	// NOTE (@afiune): We purposly leave trailing spaces in this table, we need them!
	expectedTable := `
          HOST DETAILS                   VULNERABILITIES          
--------------------------------+---------------------------------
    Machine ID   51                 SEVERITY   COUNT   FIXABLE    
    Hostname     host-mock        -----------+-------+----------  
    External IP  mock               Critical       0         0    
    Internal IP  mock               High           1         0    
    Os           linux              Medium         3         2    
    Arch         arm64              Low            4         2    
    Namespace    ubuntu:18.04       Info          14         0    
    Provider     AWS                                              
    Instance ID  i-mock                                           
    AMI          ami-mock                                         
                                                                  
There are no critical vulnerabilities in your environment.

Try adding '--active' to only show vulnerabilities of packages actively running.

0 of 8 cve(s) showing
`
	assert.Equal(t, strings.TrimPrefix(expectedTable, "\n"), cliOutput)
}

func TestSeveritySummary(t *testing.T) {
	var (
		summary string
		hosts   = mockHostVulnerabilityAssessment().Data
	)
	hostSummary := hostsSummary(hosts)
	for _, sum := range hostSummary {
		summary = severitySummary(sum.severity, sum.fixable)
	}

	expected := " 1 High 1 Medium 1 Low 4 Fixable"
	assert.Equal(t, expected, summary)
}

func TestBuildCSVVulnHostReportsWithVulnerabilitiesPackagesViewAndNoFilters(t *testing.T) {
	vulCmdState.Packages = true
	cli.EnableCSVOutput()
	defer clearVulnFilters()
	defer func() { cli.csvOutput = false }()
	expected := `
CVE Count,Severity,Package,Current Version,Fix Version,Pkg Status
1,High,ubuntu:18.04,1:2017.3.23-2ubuntu0.18.04.2,,
1,Medium,ubuntu:18.04,1:2017.3.23-2ubuntu0.18.04.2,,
1,Medium,ubuntu:18.04,1:2017.3.23-2ubuntu0.18.04.2,7.58.0-2ubuntu3.18,
1,Medium,ubuntu:18.04,1:2017.3.23-2ubuntu0.18.04.2,1:2017.3.23-2ubuntu0.18.04.4,
1,Low,ubuntu:18.04,1:2017.3.23-2ubuntu0.18.04.2,4.4.18-2ubuntu1.3,
2,Low,ubuntu:18.04,1:2017.3.23-2ubuntu0.18.04.2,,
1,Low,ubuntu:18.04,1:2017.3.23-2ubuntu0.18.04.2,3.6.9-1~18.04ubuntu1.8,
`
	cliOutput := capturer.CaptureOutput(func() {
		assert.Nil(t, buildVulnHostReports(mockHostVulnerabilityAssessment()))
	})
	assert.Equal(t, strings.TrimPrefix(expected, "\n"), cliOutput)
}

func mockHostVulnerabilityAssessment() api.VulnerabilitiesHostResponse {
	assessment := api.VulnerabilitiesHostResponse{}
	err := json.Unmarshal([]byte(`{
    "paging": {
        "rows": 774,
        "totalRows": 774,
        "urls": {
            "nextPage": null
        }
    },
    "data": [
        {
            "cveProps": {
                "metadata": {}
            },
            "endTime": "2022-09-23T04:00:00.000Z",
            "evalCtx": {
                "exception_props": [],
                "Hostname": "host-mock",
                "mc_eval_guid": "mock"
            },
            "featureKey": {
                "name": "linux-aws-5.4-headers-5.4.0-1049",
                "namespace": "ubuntu:18.04",
                "package_active": 0,
                "package_path": "",
                "version_installed": "1:2017.3.23-2ubuntu0.18.04.2"
            },
            "fixInfo": {},
            "machineTags": {
                "VpcId": "1234567891011",
                "AmiId": "ami-mock",
                "ExternalIp": "mock",
                "Hostname": "host-mock",
                "InstanceId": "i-mock",
                "InternalIp": "mock",
                
                "SubnetId": "subnet-mock",
                "VmInstanceType": "t4g.nano",
                "VmProvider": "AWS",
                "VpcId": "vpc-mock",
                "Zone": "us-mock",
                "arch": "arm64",
                "os": "linux"
            },
            "mid": 51,
            "props": {
                "isDailyJob": 1
            },
            "startTime": "2022-09-23T03:00:00.000Z"
        },
        {
            "cveProps": {
                "metadata": {}
            },
            "endTime": "2022-09-23T04:00:00.000Z",
            "evalCtx": {
                "exception_props": [],
                "Hostname": "host-mock",
                "mc_eval_guid": "mock"
            },
            "featureKey": {
                "name": "lsb-release",
                "namespace": "ubuntu:18.04",
                "package_active": 0,
                "package_path": "",
                "version_installed": "1:2017.3.23-2ubuntu0.18.04.2"
            },
            "fixInfo": {},
            "machineTags": {
                "VpcId": "1234567891011",
                "AmiId": "ami-mock",
                "ExternalIp": "mock",
                "Hostname": "host-mock",
                "InstanceId": "i-mock",
                "InternalIp": "mock",
                
                "SubnetId": "subnet-mock",
                "VmInstanceType": "t4g.nano",
                "VmProvider": "AWS",
                "VpcId": "vpc-mock",
                "Zone": "us-mock",
                "arch": "arm64",
                "os": "linux"
            },
            "mid": 51,
            "props": {
                "isDailyJob": 1
            },
            "startTime": "2022-09-23T03:00:00.000Z"
        },
        {
            "cveProps": {
                "metadata": {}
            },
            "endTime": "2022-09-23T04:00:00.000Z",
            "evalCtx": {
                "exception_props": [],
                "Hostname": "host-mock",
                "mc_eval_guid": "mock"
            },
            "featureKey": {
                "name": "libnuma1",
                "namespace": "ubuntu:18.04",
                "package_active": 0,
                "package_path": "",
                "version_installed": "1:2017.3.23-2ubuntu0.18.04.2"
            },
            "fixInfo": {},
            "machineTags": {
                "VpcId": "1234567891011",
                "AmiId": "ami-mock",
                "ExternalIp": "mock",
                "Hostname": "host-mock",
                "InstanceId": "i-mock",
                "InternalIp": "mock",
                
                "SubnetId": "subnet-mock",
                "VmInstanceType": "t4g.nano",
                "VmProvider": "AWS",
                "VpcId": "vpc-mock",
                "Zone": "us-mock",
                "arch": "arm64",
                "os": "linux"
            },
            "mid": 51,
            "props": {
                "isDailyJob": 1
            },
            "startTime": "2022-09-23T03:00:00.000Z"
        },
        {
            "cveProps": {
                "metadata": {}
            },
            "endTime": "2022-09-23T04:00:00.000Z",
            "evalCtx": {
                "exception_props": [],
                "Hostname": "host-mock",
                "mc_eval_guid": "mock"
            },
            "featureKey": {
                "name": "libpython3.6-minimal",
                "namespace": "ubuntu:18.04",
                "package_active": 0,
                "package_path": "",
                "version_installed": "1:2017.3.23-2ubuntu0.18.04.2"
            },
            "fixInfo": {},
            "machineTags": {
                "VpcId": "1234567891011",
                "AmiId": "ami-mock",
                "ExternalIp": "mock",
                "Hostname": "host-mock",
                "InstanceId": "i-mock",
                "InternalIp": "mock",
                
                "SubnetId": "subnet-mock",
                "VmInstanceType": "t4g.nano",
                "VmProvider": "AWS",
                "VpcId": "vpc-mock",
                "Zone": "us-mock",
                "arch": "arm64",
                "os": "linux"
            },
            "mid": 51,
            "props": {
                "isDailyJob": 1
            },
            "startTime": "2022-09-23T03:00:00.000Z"
        },
        {
            "cveProps": {
                "cve_batch_id": "mock-id",
                "description": "mock-description",
                "link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2022-33741",
                "metadata": {
                    "NVD": {
                        "CVSSv2": {
                            "PublishedDateTime": "2022-07-05T13:15Z",
                            "Score": 3.6,
                            "Vectors": "AV:L/AC:L/Au:N/C:P/I:N/A:P"
                        },
                        "CVSSv3": {
                            "ExploitabilityScore": 1.8,
                            "ImpactScore": 5.2,
                            "Score": 7.1,
                            "Vectors": "CVSS:3.0/AV:L/AC:L/PR:L/UI:N/S:U/C:H/I:N/A:H"
                        }
                    }
                },
                "source": "lacework"
            },
            "endTime": "2022-09-23T04:00:00.000Z",
            "evalCtx": {
                "exception_props": [],
                "Hostname": "host-mock",
                "mc_eval_guid": "mock"
            },
            "featureKey": {
                "name": "linux-aws",
                "namespace": "ubuntu:18.04",
                "package_active": 0,
                "package_path": "",
                "version_installed": "1:2017.3.23-2ubuntu0.18.04.2"
            },
            "fixInfo": {
                "fix_available": "0"
            },
            "machineTags": {
                "VpcId": "1234567891011",
                "AmiId": "ami-mock",
                "ExternalIp": "mock",
                "Hostname": "host-mock",
                "InstanceId": "i-mock",
                "InternalIp": "mock",
                
                "SubnetId": "subnet-mock",
                "VmInstanceType": "t4g.nano",
                "VmProvider": "AWS",
                "VpcId": "vpc-mock",
                "Zone": "us-mock",
                "arch": "arm64",
                "os": "linux"
            },
            "mid": 51,
            "props": {
                "first_time_seen": "2022-09-16T03:00:00.000Z",
                "isDailyJob": 1,
                "last_updated_time": "2022-09-17T03:00:00.000Z"
            },
            "severity": "High",
            "startTime": "2022-09-23T03:00:00.000Z",
            "status": "Reopened",
            "vulnId": "CVE-2022-33741"
        },
        {
            "cveProps": {
                "cve_batch_id": "mock-id",
                "description": "mock-description",
                "link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2019-18276",
                "metadata": {
                    "NVD": {
                        "CVSSv2": {
                            "PublishedDateTime": "2019-11-28T01:15Z",
                            "Score": 7.2,
                            "Vectors": "AV:L/AC:L/Au:N/C:C/I:C/A:C"
                        },
                        "CVSSv3": {
                            "ExploitabilityScore": 1.8,
                            "ImpactScore": 5.9,
                            "Score": 7.8,
                            "Vectors": "CVSS:3.0/AV:L/AC:L/PR:L/UI:N/S:U/C:H/I:H/A:H"
                        }
                    }
                },
                "source": "lacework"
            },
            "endTime": "2022-09-23T04:00:00.000Z",
            "evalCtx": {
                "exception_props": [],
                "Hostname": "host-mock",
                "mc_eval_guid": "mock"
            },
            "featureKey": {
                "name": "bash",
                "namespace": "ubuntu:18.04",
                "package_active": 0,
                "package_path": "",
                "version_installed": "1:2017.3.23-2ubuntu0.18.04.2"
            },
            "fixInfo": {
                "fix_available": "1",
                "fixed_version": "4.4.18-2ubuntu1.3"
            },
            "machineTags": {
                "VpcId": "1234567891011",
                "AmiId": "ami-mock",
                "ExternalIp": "mock",
                "Hostname": "host-mock",
                "InstanceId": "i-mock",
                "InternalIp": "mock",
                
                "SubnetId": "subnet-mock",
                "VmInstanceType": "t4g.nano",
                "VmProvider": "AWS",
                "VpcId": "vpc-mock",
                "Zone": "us-mock",
                "arch": "arm64",
                "os": "linux"
            },
            "mid": 51,
            "props": {
                "first_time_seen": "2021-01-05T11:00:00.000Z",
                "isDailyJob": 1,
                "last_updated_time": "2021-01-06T03:00:00.000Z"
            },
            "severity": "Low",
            "startTime": "2022-09-23T03:00:00.000Z",
            "status": "Active",
            "vulnId": "CVE-2019-18276"
        },
        {
            "cveProps": {
                "cve_batch_id": "mock-id",
                "description": "mock-description",
                "link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2022-27782",
                "metadata": {
                    "NVD": {
                        "CVSSv2": {
                            "PublishedDateTime": "2022-06-02T14:15Z",
                            "Score": 5,
                            "Vectors": "AV:N/AC:L/Au:N/C:N/I:P/A:N"
                        },
                        "CVSSv3": {
                            "ExploitabilityScore": 3.9,
                            "ImpactScore": 3.6,
                            "Score": 7.5,
                            "Vectors": "CVSS:3.0/AV:N/AC:L/PR:N/UI:N/S:U/C:N/I:H/A:N"
                        }
                    }
                },
                "source": "lacework"
            },
            "endTime": "2022-09-23T04:00:00.000Z",
            "evalCtx": {
                "exception_props": [],
                "Hostname": "host-mock",
                "mc_eval_guid": "mock"
            },
            "featureKey": {
                "name": "curl",
                "namespace": "ubuntu:18.04",
                "package_active": 0,
                "package_path": "",
                "version_installed": "1:2017.3.23-2ubuntu0.18.04.2"
            },
            "fixInfo": {
                "fix_available": "1",
                "fixed_version": "7.58.0-2ubuntu3.18"
            },
            "machineTags": {
                "VpcId": "1234567891011",
                "AmiId": "ami-mock",
                "ExternalIp": "mock",
                "Hostname": "host-mock",
                "InstanceId": "i-mock",
                "InternalIp": "mock",
                
                "SubnetId": "subnet-mock",
                "VmInstanceType": "t4g.nano",
                "VmProvider": "AWS",
                "VpcId": "vpc-mock",
                "Zone": "us-mock",
                "arch": "arm64",
                "os": "linux"
            },
            "mid": 51,
            "props": {
                "first_time_seen": "2022-05-12T03:00:00.000Z",
                "isDailyJob": 1,
                "last_updated_time": "2022-05-13T03:00:00.000Z"
            },
            "severity": "Medium",
            "startTime": "2022-09-23T03:00:00.000Z",
            "status": "Active",
            "vulnId": "CVE-2022-27782"
        },
        {
            "cveProps": {
                "cve_batch_id": "mock-id",
                "description": "mock-description",
                "link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2022-0351",
                "metadata": {
                    "NVD": {
                        "CVSSv2": {
                            "PublishedDateTime": "2022-01-25T18:15Z",
                            "Score": 4.6,
                            "Vectors": "AV:L/AC:L/Au:N/C:P/I:P/A:P"
                        },
                        "CVSSv3": {
                            "ExploitabilityScore": 1.8,
                            "ImpactScore": 5.9,
                            "Score": 7.8,
                            "Vectors": "CVSS:3.0/AV:L/AC:L/PR:L/UI:N/S:U/C:H/I:H/A:H"
                        }
                    }
                },
                "source": "lacework"
            },
            "endTime": "2022-09-23T04:00:00.000Z",
            "evalCtx": {
                "exception_props": [],
                "Hostname": "host-mock",
                "mc_eval_guid": "mock"
            },
            "featureKey": {
                "name": "vim",
                "namespace": "ubuntu:18.04",
                "package_active": 0,
                "package_path": "",
                "version_installed": "1:2017.3.23-2ubuntu0.18.04.2"
            },
            "fixInfo": {
                "fix_available": "0"
            },
            "machineTags": {
                "VpcId": "1234567891011",
                "AmiId": "ami-mock",
                "ExternalIp": "mock",
                "Hostname": "host-mock",
                "InstanceId": "i-mock",
                "InternalIp": "mock",
                
                "SubnetId": "subnet-mock",
                "VmInstanceType": "t4g.nano",
                "VmProvider": "AWS",
                "VpcId": "vpc-mock",
                "Zone": "us-mock",
                "arch": "arm64",
                "os": "linux"
            },
            "mid": 51,
            "props": {
                "first_time_seen": "2022-01-30T03:00:00.000Z",
                "isDailyJob": 1,
                "last_updated_time": "2022-01-31T03:00:00.000Z"
            },
            "severity": "Medium",
            "startTime": "2022-09-23T03:00:00.000Z",
            "status": "Active",
            "vulnId": "CVE-2022-0351"
        },
        {
            "cveProps": {
                "metadata": {}
            },
            "endTime": "2022-09-23T04:00:00.000Z",
            "evalCtx": {
                "exception_props": [],
                "Hostname": "host-mock",
                "mc_eval_guid": "mock"
            },
            "featureKey": {
                "name": "linux-aws-5.4-headers-5.4.0-1039",
                "namespace": "ubuntu:18.04",
                "package_active": 0,
                "package_path": "",
                "version_installed": "1:2017.3.23-2ubuntu0.18.04.2"
            },
            "fixInfo": {},
            "machineTags": {
                "VpcId": "1234567891011",
                "AmiId": "ami-mock",
                "ExternalIp": "mock",
                "Hostname": "host-mock",
                "InstanceId": "i-mock",
                "InternalIp": "mock",
                
                "SubnetId": "subnet-mock",
                "VmInstanceType": "t4g.nano",
                "VmProvider": "AWS",
                "VpcId": "vpc-mock",
                "Zone": "us-mock",
                "arch": "arm64",
                "os": "linux"
            },
            "mid": 51,
            "props": {
                "isDailyJob": 1
            },
            "startTime": "2022-09-23T03:00:00.000Z"
        },
        {
            "cveProps": {
                "metadata": {}
            },
            "endTime": "2022-09-23T04:00:00.000Z",
            "evalCtx": {
                "exception_props": [],
                "Hostname": "host-mock",
                "mc_eval_guid": "mock"
            },
            "featureKey": {
                "name": "libcurl4",
                "namespace": "ubuntu:18.04",
                "package_active": 0,
                "package_path": "",
                "version_installed": "1:2017.3.23-2ubuntu0.18.04.2"
            },
            "fixInfo": {},
            "machineTags": {
                "VpcId": "1234567891011",
                "AmiId": "ami-mock",
                "ExternalIp": "mock",
                "Hostname": "host-mock",
                "InstanceId": "i-mock",
                "InternalIp": "mock",
                
                "SubnetId": "subnet-mock",
                "VmInstanceType": "t4g.nano",
                "VmProvider": "AWS",
                "VpcId": "vpc-mock",
                "Zone": "us-mock",
                "arch": "arm64",
                "os": "linux"
            },
            "mid": 51,
            "props": {
                "isDailyJob": 1
            },
            "startTime": "2022-09-23T03:00:00.000Z"
        },
        {
            "cveProps": {
                "metadata": {}
            },
            "endTime": "2022-09-23T04:00:00.000Z",
            "evalCtx": {
                "exception_props": [],
                "Hostname": "host-mock",
                "mc_eval_guid": "mock"
            },
            "featureKey": {
                "name": "htop",
                "namespace": "ubuntu:18.04",
                "package_active": 0,
                "package_path": "",
                "version_installed": "1:2017.3.23-2ubuntu0.18.04.2"
            },
            "fixInfo": {},
            "machineTags": {
                "VpcId": "1234567891011",
                "AmiId": "ami-mock",
                "ExternalIp": "mock",
                "Hostname": "host-mock",
                "InstanceId": "i-mock",
                "InternalIp": "mock",
                
                "SubnetId": "subnet-mock",
                "VmInstanceType": "t4g.nano",
                "VmProvider": "AWS",
                "VpcId": "vpc-mock",
                "Zone": "us-mock",
                "arch": "arm64",
                "os": "linux"
            },
            "mid": 51,
            "props": {
                "isDailyJob": 1
            },
            "startTime": "2022-09-23T03:00:00.000Z"
        },
        {
            "cveProps": {
                "cve_batch_id": "mock-id",
                "description": "mock-description",
                "link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2022-2129",
                "metadata": {
                    "NVD": {
                        "CVSSv2": {
                            "PublishedDateTime": "2022-06-19T19:15Z",
                            "Score": 6.8,
                            "Vectors": "AV:N/AC:M/Au:N/C:P/I:P/A:P"
                        },
                        "CVSSv3": {
                            "ExploitabilityScore": 1.8,
                            "ImpactScore": 5.9,
                            "Score": 7.8,
                            "Vectors": "CVSS:3.0/AV:L/AC:L/PR:N/UI:R/S:U/C:H/I:H/A:H"
                        }
                    }
                },
                "source": "lacework"
            },
            "endTime": "2022-09-23T04:00:00.000Z",
            "evalCtx": {
                "exception_props": [],
                "Hostname": "host-mock",
                "mc_eval_guid": "mock"
            },
            "featureKey": {
                "name": "vim",
                "namespace": "ubuntu:18.04",
                "package_active": 0,
                "package_path": "",
                "version_installed": "1:2017.3.23-2ubuntu0.18.04.2"
            },
            "fixInfo": {
                "fix_available": "0"
            },
            "machineTags": {
                "VpcId": "1234567891011",
                "AmiId": "ami-mock",
                "ExternalIp": "mock",
                "Hostname": "host-mock",
                "InstanceId": "i-mock",
                "InternalIp": "mock",
                
                "SubnetId": "subnet-mock",
                "VmInstanceType": "t4g.nano",
                "VmProvider": "AWS",
                "VpcId": "vpc-mock",
                "Zone": "us-mock",
                "arch": "arm64",
                "os": "linux"
            },
            "mid": 51,
            "props": {
                "first_time_seen": "2022-07-20T03:00:00.000Z",
                "isDailyJob": 1,
                "last_updated_time": "2022-07-21T03:00:00.000Z"
            },
            "severity": "Low",
            "startTime": "2022-09-23T03:00:00.000Z",
            "status": "Active",
            "vulnId": "CVE-2022-2129"
        },
        {
            "cveProps": {
                "metadata": {}
            },
            "endTime": "2022-09-23T04:00:00.000Z",
            "evalCtx": {
                "exception_props": [],
                "Hostname": "host-mock",
                "mc_eval_guid": "mock"
            },
            "featureKey": {
                "name": "telnet",
                "namespace": "ubuntu:18.04",
                "package_active": 0,
                "package_path": "",
                "version_installed": "1:2017.3.23-2ubuntu0.18.04.2"
            },
            "fixInfo": {},
            "machineTags": {
                "VpcId": "1234567891011",
                "AmiId": "ami-mock",
                "ExternalIp": "mock",
                "Hostname": "host-mock",
                "InstanceId": "i-mock",
                "InternalIp": "mock",
                
                "SubnetId": "subnet-mock",
                "VmInstanceType": "t4g.nano",
                "VmProvider": "AWS",
                "VpcId": "vpc-mock",
                "Zone": "us-mock",
                "arch": "arm64",
                "os": "linux"
            },
            "mid": 51,
            "props": {
                "isDailyJob": 1
            },
            "startTime": "2022-09-23T03:00:00.000Z"
        },
        {
            "cveProps": {
                "metadata": {}
            },
            "endTime": "2022-09-23T04:00:00.000Z",
            "evalCtx": {
                "exception_props": [],
                "Hostname": "host-mock",
                "mc_eval_guid": "mock"
            },
            "featureKey": {
                "name": "libpython3.6",
                "namespace": "ubuntu:18.04",
                "package_active": 0,
                "package_path": "",
                "version_installed": "1:2017.3.23-2ubuntu0.18.04.2"
            },
            "fixInfo": {},
            "machineTags": {
                "VpcId": "1234567891011",
                "AmiId": "ami-mock",
                "ExternalIp": "mock",
                "Hostname": "host-mock",
                "InstanceId": "i-mock",
                "InternalIp": "mock",
                
                "SubnetId": "subnet-mock",
                "VmInstanceType": "t4g.nano",
                "VmProvider": "AWS",
                "VpcId": "vpc-mock",
                "Zone": "us-mock",
                "arch": "arm64",
                "os": "linux"
            },
            "mid": 51,
            "props": {
                "isDailyJob": 1
            },
            "startTime": "2022-09-23T03:00:00.000Z"
        },
        {
            "cveProps": {
                "metadata": {}
            },
            "endTime": "2022-09-23T04:00:00.000Z",
            "evalCtx": {
                "exception_props": [],
                "Hostname": "host-mock",
                "mc_eval_guid": "mock"
            },
            "featureKey": {
                "name": "vim",
                "namespace": "ubuntu:18.04",
                "package_active": 0,
                "package_path": "",
                "version_installed": "1:2017.3.23-2ubuntu0.18.04.2"
            },
            "fixInfo": {},
            "machineTags": {
                "VpcId": "1234567891011",
                "AmiId": "ami-mock",
                "ExternalIp": "mock",
                "Hostname": "host-mock",
                "InstanceId": "i-mock",
                "InternalIp": "mock",
                
                "SubnetId": "subnet-mock",
                "VmInstanceType": "t4g.nano",
                "VmProvider": "AWS",
                "VpcId": "vpc-mock",
                "Zone": "us-mock",
                "arch": "arm64",
                "os": "linux"
            },
            "mid": 51,
            "props": {
                "isDailyJob": 1
            },
            "startTime": "2022-09-23T03:00:00.000Z"
        },
        {
            "cveProps": {
                "metadata": {}
            },
            "endTime": "2022-09-23T04:00:00.000Z",
            "evalCtx": {
                "exception_props": [],
                "Hostname": "host-mock",
                "mc_eval_guid": "mock"
            },
            "featureKey": {
                "name": "publicsuffix",
                "namespace": "ubuntu:18.04",
                "package_active": 0,
                "package_path": "",
                "version_installed": "1:2017.3.23-2ubuntu0.18.04.2"
            },
            "fixInfo": {},
            "machineTags": {
                "VpcId": "1234567891011",
                "AmiId": "ami-mock",
                "ExternalIp": "mock",
                "Hostname": "host-mock",
                "InstanceId": "i-mock",
                "InternalIp": "mock",
                
                "SubnetId": "subnet-mock",
                "VmInstanceType": "t4g.nano",
                "VmProvider": "AWS",
                "VpcId": "vpc-mock",
                "Zone": "us-mock",
                "arch": "arm64",
                "os": "linux"
            },
            "mid": 51,
            "props": {
                "isDailyJob": 1
            },
            "startTime": "2022-09-23T03:00:00.000Z"
        },
        {
            "cveProps": {
                "cve_batch_id": "mock-id",
                "description": "mock-description",
                "link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2015-20107",
                "metadata": {
                    "NVD": {
                        "CVSSv2": {
                            "PublishedDateTime": "2022-04-13T16:15Z",
                            "Score": 10,
                            "Vectors": "AV:N/AC:L/Au:N/C:C/I:C/A:C"
                        },
                        "CVSSv3": {
                            "ExploitabilityScore": 3.9,
                            "ImpactScore": 5.9,
                            "Score": 9.8,
                            "Vectors": "CVSS:3.0/AV:N/AC:L/PR:N/UI:N/S:U/C:H/I:H/A:H"
                        }
                    }
                },
                "source": "lacework"
            },
            "endTime": "2022-09-23T04:00:00.000Z",
            "evalCtx": {
                "exception_props": [],
                "Hostname": "host-mock",
                "mc_eval_guid": "mock"
            },
            "featureKey": {
                "name": "python3.6",
                "namespace": "ubuntu:18.04",
                "package_active": 0,
                "package_path": "",
                "version_installed": "1:2017.3.23-2ubuntu0.18.04.2"
            },
            "fixInfo": {
                "fix_available": "1",
                "fixed_version": "3.6.9-1~18.04ubuntu1.8"
            },
            "machineTags": {
                "VpcId": "1234567891011",
                "AmiId": "ami-mock",
                "ExternalIp": "mock",
                "Hostname": "host-mock",
                "InstanceId": "i-mock",
                "InternalIp": "mock",
                
                "SubnetId": "subnet-mock",
                "VmInstanceType": "t4g.nano",
                "VmProvider": "AWS",
                "VpcId": "vpc-mock",
                "Zone": "us-mock",
                "arch": "arm64",
                "os": "linux"
            },
            "mid": 51,
            "props": {
                "first_time_seen": "2022-04-16T03:00:00.000Z",
                "isDailyJob": 1,
                "last_updated_time": "2022-04-17T03:00:00.000Z"
            },
            "severity": "Low",
            "startTime": "2022-09-23T03:00:00.000Z",
            "status": "Active",
            "vulnId": "CVE-2015-20107"
        },
        {
            "cveProps": {
                "cve_batch_id": "mock-id",
                "description": "mock-description",
                "link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2020-13988",
                "metadata": {
                    "NVD": {
                        "CVSSv2": {
                            "PublishedDateTime": "2020-12-11T22:15Z",
                            "Score": 5,
                            "Vectors": "AV:N/AC:L/Au:N/C:N/I:N/A:P"
                        },
                        "CVSSv3": {
                            "ExploitabilityScore": 3.9,
                            "ImpactScore": 3.6,
                            "Score": 7.5,
                            "Vectors": "CVSS:3.0/AV:N/AC:L/PR:N/UI:N/S:U/C:N/I:N/A:H"
                        }
                    }
                },
                "source": "lacework"
            },
            "endTime": "2022-09-23T04:00:00.000Z",
            "evalCtx": {
                "exception_props": [],
                "Hostname": "host-mock",
                "mc_eval_guid": "mock"
            },
            "featureKey": {
                "name": "open-iscsi",
                "namespace": "ubuntu:18.04",
                "package_active": 0,
                "package_path": "",
                "version_installed": "1:2017.3.23-2ubuntu0.18.04.2"
            },
            "fixInfo": {
                "fix_available": "0"
            },
            "machineTags": {
                "VpcId": "1234567891011",
                "AmiId": "ami-mock",
                "ExternalIp": "mock",
                "Hostname": "host-mock",
                "InstanceId": "i-mock",
                "InternalIp": "mock",
                
                "SubnetId": "subnet-mock",
                "VmInstanceType": "t4g.nano",
                "VmProvider": "AWS",
                "VpcId": "vpc-mock",
                "Zone": "us-mock",
                "arch": "arm64",
                "os": "linux"
            },
            "mid": 51,
            "props": {
                "first_time_seen": "2021-01-23T03:00:00.000Z",
                "isDailyJob": 1,
                "last_updated_time": "2021-01-24T03:00:00.000Z"
            },
            "severity": "Low",
            "startTime": "2022-09-23T03:00:00.000Z",
            "status": "Active",
            "vulnId": "CVE-2020-13988"
        },
        {
            "cveProps": {
                "metadata": {}
            },
            "endTime": "2022-09-23T04:00:00.000Z",
            "evalCtx": {
                "exception_props": [],
                "Hostname": "host-mock",
                "mc_eval_guid": "mock"
            },
            "featureKey": {
                "name": "libpsl5",
                "namespace": "ubuntu:18.04",
                "package_active": 0,
                "package_path": "",
                "version_installed": "1:2017.3.23-2ubuntu0.18.04.2"
            },
            "fixInfo": {},
            "machineTags": {
                "VpcId": "1234567891011",
                "AmiId": "ami-mock",
                "ExternalIp": "mock",
                "Hostname": "host-mock",
                "InstanceId": "i-mock",
                "InternalIp": "mock",
                
                "SubnetId": "subnet-mock",
                "VmInstanceType": "t4g.nano",
                "VmProvider": "AWS",
                "VpcId": "vpc-mock",
                "Zone": "us-mock",
                "arch": "arm64",
                "os": "linux"
            },
            "mid": 51,
            "props": {
                "isDailyJob": 1
            },
            "startTime": "2022-09-23T03:00:00.000Z"
        },
        {
            "cveProps": {
                "metadata": {}
            },
            "endTime": "2022-09-23T04:00:00.000Z",
            "evalCtx": {
                "exception_props": [],
                "Hostname": "host-mock",
                "mc_eval_guid": "mock"
            },
            "featureKey": {
                "name": "dirmngr",
                "namespace": "ubuntu:18.04",
                "package_active": 0,
                "package_path": "",
                "version_installed": "1:2017.3.23-2ubuntu0.18.04.2"
            },
            "fixInfo": {},
            "machineTags": {
                "VpcId": "1234567891011",
                "AmiId": "ami-mock",
                "ExternalIp": "mock",
                "Hostname": "host-mock",
                "InstanceId": "i-mock",
                "InternalIp": "mock",
                
                "SubnetId": "subnet-mock",
                "VmInstanceType": "t4g.nano",
                "VmProvider": "AWS",
                "VpcId": "vpc-mock",
                "Zone": "us-mock",
                "arch": "arm64",
                "os": "linux"
            },
            "mid": 51,
            "props": {
                "isDailyJob": 1
            },
            "startTime": "2022-09-23T03:00:00.000Z"
        },
        {
            "cveProps": {
                "metadata": {}
            },
            "endTime": "2022-09-23T04:00:00.000Z",
            "evalCtx": {
                "exception_props": [],
                "Hostname": "host-mock",
                "mc_eval_guid": "mock"
            },
            "featureKey": {
                "name": "libfribidi0",
                "namespace": "ubuntu:18.04",
                "package_active": 0,
                "package_path": "",
                "version_installed": "1:2017.3.23-2ubuntu0.18.04.2"
            },
            "fixInfo": {},
            "machineTags": {
                "VpcId": "1234567891011",
                "AmiId": "ami-mock",
                "ExternalIp": "mock",
                "Hostname": "host-mock",
                "InstanceId": "i-mock",
                "InternalIp": "mock",
                
                "SubnetId": "subnet-mock",
                "VmInstanceType": "t4g.nano",
                "VmProvider": "AWS",
                "VpcId": "vpc-mock",
                "Zone": "us-mock",
                "arch": "arm64",
                "os": "linux"
            },
            "mid": 51,
            "props": {
                "isDailyJob": 1
            },
            "startTime": "2022-09-23T03:00:00.000Z"
        },
        {
            "cveProps": {
                "metadata": {}
            },
            "endTime": "2022-09-23T04:00:00.000Z",
            "evalCtx": {
                "exception_props": [],
                "Hostname": "host-mock",
                "mc_eval_guid": "mock"
            },
            "featureKey": {
                "name": "linux-aws",
                "namespace": "ubuntu:18.04",
                "package_active": 0,
                "package_path": "",
                "version_installed": "1:2017.3.23-2ubuntu0.18.04.2"
            },
            "fixInfo": {
                "fix_available": "0"
            },
            "machineTags": {
                "VpcId": "1234567891011",
                "AmiId": "ami-mock",
                "ExternalIp": "mock",
                "Hostname": "host-mock",
                "InstanceId": "i-mock",
                "InternalIp": "mock",
                
                "SubnetId": "subnet-mock",
                "VmInstanceType": "t4g.nano",
                "VmProvider": "AWS",
                "VpcId": "vpc-mock",
                "Zone": "us-mock",
                "arch": "arm64",
                "os": "linux"
            },
            "mid": 51,
            "props": {
                "first_time_seen": "2022-08-18T20:00:00.000Z",
                "isDailyJob": 1,
                "last_updated_time": "2022-09-15T03:00:00.000Z"
            },
            "severity": "Medium",
            "startTime": "2022-09-23T03:00:00.000Z",
            "status": "Fixed",
            "vulnId": "CVE-2021-33655"
        },
        {
            "cveProps": {
                "cve_batch_id": "mock-id",
                "description": "mock-description",
                "link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2022-30787",
                "metadata": {
                    "NVD": {
                        "CVSSv2": {
                            "PublishedDateTime": "2022-05-26T16:15Z",
                            "Score": 4.6,
                            "Vectors": "AV:L/AC:L/Au:N/C:P/I:P/A:P"
                        },
                        "CVSSv3": {
                            "ExploitabilityScore": 0.8,
                            "ImpactScore": 5.9,
                            "Score": 6.7,
                            "Vectors": "CVSS:3.0/AV:L/AC:L/PR:H/UI:N/S:U/C:H/I:H/A:H"
                        }
                    }
                },
                "source": "lacework"
            },
            "endTime": "2022-09-23T04:00:00.000Z",
            "evalCtx": {
                "exception_props": [],
                "Hostname": "host-mock",
                "mc_eval_guid": "mock"
            },
            "featureKey": {
                "name": "ntfs-3g",
                "namespace": "ubuntu:18.04",
                "package_active": 0,
                "package_path": "",
                "version_installed": "1:2017.3.23-2ubuntu0.18.04.2"
            },
            "fixInfo": {
                "fix_available": "1",
                "fixed_version": "1:2017.3.23-2ubuntu0.18.04.4"
            },
            "machineTags": {
                "VpcId": "1234567891011",
                "AmiId": "ami-mock",
                "ExternalIp": "mock",
                "Hostname": "host-mock",
                "InstanceId": "i-mock",
                "InternalIp": "mock",
                "SubnetId": "subnet-mock",
                "VmInstanceType": "t4g.nano",
                "VmProvider": "AWS",
                "VpcId": "vpc-mock",
                "Zone": "us-mock",
                "arch": "arm64",
                "os": "linux"
            },
            "mid": 51,
            "props": {
                "first_time_seen": "2022-06-08T03:00:00.000Z",
                "isDailyJob": 1,
                "last_updated_time": "2022-06-09T03:00:00.000Z"
            },
            "severity": "Medium",
            "startTime": "2022-09-23T03:00:00.000Z",
            "status": "Active",
            "vulnId": "CVE-2022-30787"
        }
	]
}`),
		&assessment)
	if err != nil {
		log.Fatal(err)
	}
	return assessment
}

var expectedCSVHostDetailsTable = `
CVE ID,Severity,Score,Package,Package Namespace,Current Version,Fix Version,Pkg Status,First Seen,Last Status Update,Vuln Status
CVE-2022-33741,High,3.6,7.1,linux-aws,ubuntu:18.04,1:2017.3.23-2ubuntu0.18.04.2,,,2022-09-16 03:00:00 +0000 UTC,2022-09-17 03:00:00 +0000 UTC,
CVE-2022-30787,Medium,4.6,6.7,ntfs-3g,ubuntu:18.04,1:2017.3.23-2ubuntu0.18.04.2,1:2017.3.23-2ubuntu0.18.04.4,,2022-06-08 03:00:00 +0000 UTC,2022-06-09 03:00:00 +0000 UTC,
CVE-2022-0351,Medium,4.6,7.8,vim,ubuntu:18.04,1:2017.3.23-2ubuntu0.18.04.2,,,2022-01-30 03:00:00 +0000 UTC,2022-01-31 03:00:00 +0000 UTC,
CVE-2022-27782,Medium,5.0,7.5,curl,ubuntu:18.04,1:2017.3.23-2ubuntu0.18.04.2,7.58.0-2ubuntu3.18,,2022-05-12 03:00:00 +0000 UTC,2022-05-13 03:00:00 +0000 UTC,
CVE-2020-13988,Low,5.0,7.5,open-iscsi,ubuntu:18.04,1:2017.3.23-2ubuntu0.18.04.2,,,2021-01-23 03:00:00 +0000 UTC,2021-01-24 03:00:00 +0000 UTC,
CVE-2022-2129,Low,6.8,7.8,vim,ubuntu:18.04,1:2017.3.23-2ubuntu0.18.04.2,,,2022-07-20 03:00:00 +0000 UTC,2022-07-21 03:00:00 +0000 UTC,
CVE-2015-20107,Low,10.0,9.8,python3.6,ubuntu:18.04,1:2017.3.23-2ubuntu0.18.04.2,3.6.9-1~18.04ubuntu1.8,,2022-04-16 03:00:00 +0000 UTC,2022-04-17 03:00:00 +0000 UTC,
CVE-2019-18276,Low,7.2,7.8,bash,ubuntu:18.04,1:2017.3.23-2ubuntu0.18.04.2,4.4.18-2ubuntu1.3,,2021-01-05 11:00:00 +0000 UTC,2021-01-06 03:00:00 +0000 UTC,
`
