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

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/lwtime"
	"github.com/stretchr/testify/assert"
)

func TestListCvesFilterSeverity(t *testing.T) {
	vulCmdState.Severity = "critical"
	defer clearVulnFilters()

	mockCves := []api.HostVulnCVE{mockCveOne}
	result, output := filterHostCVEsTable(mockCves)

	assert.Equal(t, len(result), 1)
	assert.Equal(t, output, "\n1 of 2 cve(s) showing\n")
}

func TestShowAssessmentFilterSeverity(t *testing.T) {
	vulCmdState.Severity = "critical"
	defer clearVulnFilters()

	mockCves := []api.HostVulnCVE{mockCveOne}
	result, output := filterHostCVEsTable(mockCves)

	assert.Equal(t, len(result), 1)
	assert.Equal(t, output, "\n1 of 2 cve(s) showing\n")
}

func TestShowAssessmentFilterSeverityWithPackages(t *testing.T) {
	vulCmdState.Severity = "critical"
	vulCmdState.Packages = true
	defer clearVulnFilters()

	mockCves := []api.HostVulnCVE{mockCveOne}
	result, output := hostVulnPackagesTable(mockCves, true)

	assert.Equal(t, len(result), 1)
	assert.Equal(t, output, "1 of 2 package(s) showing\n")
}

func clearVulnFilters() {
	vulCmdState.Severity = ""
	vulCmdState.Packages = false
	vulCmdState.Active = false
}

var mockCveOne = api.HostVulnCVE{
	ID:       "TestID",
	Packages: []api.HostVulnPackage{mockPackageOne, mockPackageTwo},
	Summary: api.HostVulnCveSummary{
		Severity: api.HostVulnSeverityCounts{
			Critical: &api.HostVulnSeverityCountsDetails{
				Fixable:         1,
				Vulnerabilities: 1,
			},
			High: &api.HostVulnSeverityCountsDetails{
				Fixable:         1,
				Vulnerabilities: 1,
			},
		},
		TotalVulnerabilities: 2,
		LastEvaluationTime:   lwtime.EpochString{},
	},
}

var mockPackageOne = api.HostVulnPackage{
	Name:         "test",
	Namespace:    "rhel:8",
	Severity:     "High",
	Status:       "Active",
	HostCount:    "1",
	FixAvailable: "1",
}

var mockPackageTwo = api.HostVulnPackage{
	Name:         "test2",
	Namespace:    "rhel:8",
	Severity:     "Critical",
	Status:       "Active",
	HostCount:    "1",
	FixAvailable: "1",
}

func TestBuildVulnHostReportsNoVulnerabilities(t *testing.T) {
	cliOutput := captureOutput(func() {
		assert.Nil(t, buildVulnHostReports(api.HostVulnHostAssessment{}))
	})
	assert.Contains(t, cliOutput, "Great news! This host has no vulnerabilities...")

	t.Run("test JSON output", func(t *testing.T) {
		cli.EnableJSONOutput()
		defer cli.EnableHumanOutput()
		cliJSONOutput := captureOutput(func() {
			assert.Nil(t, buildVulnHostReports(api.HostVulnHostAssessment{}))
		})
		expectedJSON := `{
  "host": {
    "hostname": "",
    "machine_id": "",
    "tags": {
      "Account": "",
      "AmiId": "",
      "ExternalIp": "",
      "Hostname": "",
      "InstanceId": "",
      "InternalIp": "",
      "LwTokenShort": "",
      "SubnetId": "",
      "VmInstanceType": "",
      "VmProvider": "",
      "VpcId": "",
      "Zone": "",
      "arch": "",
      "os": ""
    }
  },
  "vulnerabilities": []
}
`
		assert.Equal(t, expectedJSON, cliJSONOutput)
	})
}

func TestBuildVulnHostReportsWithVulnerabilitiesSummaryOnlyAndNoFilters(t *testing.T) {
	cliOutput := captureOutput(func() {
		assert.Nil(t, buildVulnHostReports(mockHostVulnerabilityAssessment()))
	})
	// NOTE (@afiune): We purposly leave trailing spaces in this table, we need them!
	expectedTable := `
          HOST DETAILS                   VULNERABILITIES          
--------------------------------+---------------------------------
    Machine ID   51                 SEVERITY   COUNT   FIXABLE    
    Hostname     ip-10-0-1-48     -----------+-------+----------  
    External IP  1.2.3.4            Critical       0         0    
    Internal IP  10.0.1.1           High           1         1    
    Os           linux              Medium        42        38    
    Arch         arm64              Low           34         5    
    Namespace    ubuntu:18.04       Info           3         0    
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

	cliOutput := captureOutput(func() {
		assert.Nil(t, buildVulnHostReports(mockHostVulnerabilityAssessment()))
	})
	// NOTE (@afiune): We purposly leave trailing spaces in this table, we need them!
	expectedTable := `
          HOST DETAILS                   VULNERABILITIES          
--------------------------------+---------------------------------
    Machine ID   51                 SEVERITY   COUNT   FIXABLE    
    Hostname     ip-10-0-1-48     -----------+-------+----------  
    External IP  1.2.3.4            Critical       0         0    
    Internal IP  10.0.1.1           High           1         1    
    Os           linux              Medium        42        38    
    Arch         arm64              Low           34         5    
    Namespace    ubuntu:18.04       Info           3         0    
    Provider     AWS                                              
    Instance ID  i-mock                                           
    AMI          ami-mock                                         
                                                                  
  CVE COUNT   SEVERITY       PACKAGE       CURRENT VERSION              FIX VERSION               PKG STATUS  
------------+----------+-----------------+-----------------+------------------------------------+-------------
  1           High       openssl                             1.1.1-1ubuntu2.1~18.04.13                        
  2           Medium     vim                                 2:8.0.1453-1ubuntu1.6                            
  1           Medium     wget                                                                     ACTIVE      
  2           Medium     apport                              2.20.9-0ubuntu7.26                               
  1           Medium     libgcrypt20                         1.8.1-4ubuntu1.3                                 
  1           Medium     cloud-init                          21.1-19-gbad84ad4-0ubuntu1~18.04.1               
  2           Medium     curl                                7.58.0-2ubuntu3.15                               
  1           Medium     squashfs-tools                      1:4.3-6ubuntu0.18.04.4                           
  1           Medium     apparmor                                                                             
  1           Medium     cpio                                2.12+dfsg-6ubuntu0.18.04.4                       
  2           Medium     vim                                 2:8.0.1453-1ubuntu1.7                            
  2           Medium     curl                                7.58.0-2ubuntu3.14                               
  1           Medium     squashfs-tools                      1:4.3-6ubuntu0.18.04.3                           
  1           Medium     openssl                             1.1.1-1ubuntu2.1~18.04.13                        
  21          Medium     ntfs-3g                             1:2017.3.23-2ubuntu0.18.04.3                     
  2           Medium     snapd                                                                                
  1           Medium     git                                 1:2.17.1-1ubuntu0.9                              
  1           Low        tcpdump                                                                              
  1           Low        vim                                 2:8.0.1453-1ubuntu1.7                            
  1           Low        git                                                                                  
  3           Low        open-iscsi                                                                           
  1           Low        fuse                                                                                 
  1           Low        dbus                                                                     ACTIVE      
  2           Low        binutils                            2.30-21ubuntu1~18.04.7                           
  1           Low        curl                                7.58.0-2ubuntu3.14                               
  1           Low        cron                                                                                 
  1           Low        systemd                                                                  ACTIVE      
  1           Low        libgcrypt20                         1.8.1-4ubuntu1.3                                 
  6           Low        binutils                                                                             
  1           Low        iptables                                                                             
  1           Low        vim                                                                                  
  1           Low        byobu                                                                                
  1           Low        xdg-user-dirs                                                                        
  3           Low        python3.6                                                                            
  2           Low        rsyslog                                                                  ACTIVE      
  1           Low        snapd                                                                                
  1           Low        bash                                                                                 
  1           Low        policykit-1                                                                          
  1           Low        coreutils                                                                            
  1           Low        accountsservice                                                                      
  1           Info       binutils                                                                             
  1           Info       patch                                                                                
  1           Info       libtasn1-6                                                                           

Try adding '--active' to only show vulnerabilities of packages actively running.
`
	assert.Equal(t, strings.TrimPrefix(expectedTable, "\n"), cliOutput)
}

func TestBuildVulnHostReportsWithVulnerabilitiesPackagesViewWithFilters(t *testing.T) {
	vulCmdState.Packages = true
	vulCmdState.Severity = "high"
	defer clearVulnFilters()

	cliOutput := captureOutput(func() {
		assert.Nil(t, buildVulnHostReports(mockHostVulnerabilityAssessment()))
	})
	// NOTE (@afiune): We purposly leave trailing spaces in this table, we need them!
	expectedTable := `
          HOST DETAILS                   VULNERABILITIES          
--------------------------------+---------------------------------
    Machine ID   51                 SEVERITY   COUNT   FIXABLE    
    Hostname     ip-10-0-1-48     -----------+-------+----------  
    External IP  1.2.3.4            Critical       0         0    
    Internal IP  10.0.1.1           High           1         1    
    Os           linux              Medium        42        38    
    Arch         arm64              Low           34         5    
    Namespace    ubuntu:18.04       Info           3         0    
    Provider     AWS                                              
    Instance ID  i-mock                                           
    AMI          ami-mock                                         
                                                                  
  CVE COUNT   SEVERITY   PACKAGE   CURRENT VERSION          FIX VERSION          PKG STATUS  
------------+----------+---------+-----------------+---------------------------+-------------
  1           High       openssl                     1.1.1-1ubuntu2.1~18.04.13               

Try adding '--active' to only show vulnerabilities of packages actively running.

1 of 80 cve(s) showing
`
	assert.Equal(t, strings.TrimPrefix(expectedTable, "\n"), cliOutput)
}

func TestBuildVulnHostReportsWithVulnerabilitiesWithFiltersSeverity(t *testing.T) {
	vulCmdState.Severity = "high"
	defer clearVulnFilters()

	cliOutput := captureOutput(func() {
		assert.Nil(t, buildVulnHostReports(mockHostVulnerabilityAssessment()))
	})
	// NOTE (@afiune): We purposly leave trailing spaces in this table, we need them!
	expectedTable := `
          HOST DETAILS                   VULNERABILITIES          
--------------------------------+---------------------------------
    Machine ID   51                 SEVERITY   COUNT   FIXABLE    
    Hostname     ip-10-0-1-48     -----------+-------+----------  
    External IP  1.2.3.4            Critical       0         0    
    Internal IP  10.0.1.1           High           1         1    
    Os           linux              Medium        42        38    
    Arch         arm64              Low           34         5    
    Namespace    ubuntu:18.04       Info           3         0    
    Provider     AWS                                              
    Instance ID  i-mock                                           
    AMI          ami-mock                                         
                                                                  
     CVE ID       SEVERITY   SCORE   PACKAGE   CURRENT VERSION          FIX VERSION          PKG STATUS   VULN STATUS  
----------------+----------+-------+---------+-----------------+---------------------------+------------+--------------
  CVE-2021-3711   High       9.8     openssl                     1.1.1-1ubuntu2.1~18.04.13                Active       

Try adding '--active' to only show vulnerabilities of packages actively running.

1 of 80 cve(s) showing
`
	assert.Equal(t, strings.TrimPrefix(expectedTable, "\n"), cliOutput)
}

func TestBuildVulnHostReportsWithVulnerabilitiesWithFiltersActive(t *testing.T) {
	vulCmdState.Active = true
	defer clearVulnFilters()

	cliOutput := captureOutput(func() {
		assert.Nil(t, buildVulnHostReports(mockHostVulnerabilityAssessment()))
	})
	// NOTE (@afiune): We purposly leave trailing spaces in this table, we need them!
	expectedTable := `
          HOST DETAILS                   VULNERABILITIES          
--------------------------------+---------------------------------
    Machine ID   51                 SEVERITY   COUNT   FIXABLE    
    Hostname     ip-10-0-1-48     -----------+-------+----------  
    External IP  1.2.3.4            Critical       0         0    
    Internal IP  10.0.1.1           High           1         1    
    Os           linux              Medium        42        38    
    Arch         arm64              Low           34         5    
    Namespace    ubuntu:18.04       Info           3         0    
    Provider     AWS                                              
    Instance ID  i-mock                                           
    AMI          ami-mock                                         
                                                                  
      CVE ID       SEVERITY   SCORE   PACKAGE   CURRENT VERSION   FIX VERSION   PKG STATUS   VULN STATUS  
-----------------+----------+-------+---------+-----------------+-------------+------------+--------------
  CVE-2021-31879   Medium     6.1     wget                                      ACTIVE       Active       
  CVE-2019-17042   Low        9.8     rsyslog                                   ACTIVE       Active       
  CVE-2020-13776   Low        6.7     systemd                                   ACTIVE       Active       
  CVE-2019-17041   Low        9.8     rsyslog                                   ACTIVE       Active       
  CVE-2020-35512   Low        7.8     dbus                                      ACTIVE       Active       

Try adding '--fixable' to only show fixable vulnerabilities.

5 of 80 cve(s) showing
`
	assert.Equal(t, strings.TrimPrefix(expectedTable, "\n"), cliOutput)
}

func TestBuildVulnHostReportsWithVulnerabilitiesWithFiltersSeverityAndActiveReturnsNoVulns(t *testing.T) {
	vulCmdState.Severity = "high"
	vulCmdState.Active = true
	defer clearVulnFilters()

	cliOutput := captureOutput(func() {
		assert.Nil(t, buildVulnHostReports(mockHostVulnerabilityAssessment()))
	})
	// NOTE (@afiune): We purposly leave trailing spaces in this table, we need them!
	expectedTable := `
          HOST DETAILS                   VULNERABILITIES          
--------------------------------+---------------------------------
    Machine ID   51                 SEVERITY   COUNT   FIXABLE    
    Hostname     ip-10-0-1-48     -----------+-------+----------  
    External IP  1.2.3.4            Critical       0         0    
    Internal IP  10.0.1.1           High           1         1    
    Os           linux              Medium        42        38    
    Arch         arm64              Low           34         5    
    Namespace    ubuntu:18.04       Info           3         0    
    Provider     AWS                                              
    Instance ID  i-mock                                           
    AMI          ami-mock                                         
                                                                  
There are no high vulnerabilities of packages actively running in your environment.

Try adding '--fixable' to only show fixable vulnerabilities.

0 of 80 cve(s) showing
`
	assert.Equal(t, strings.TrimPrefix(expectedTable, "\n"), cliOutput)
}

func TestBuildVulnHostReportsWithVulnerabilitiesWithFilterReturnsNoVulns(t *testing.T) {
	vulCmdState.Severity = "critical"
	defer clearVulnFilters()

	cliOutput := captureOutput(func() {
		assert.Nil(t, buildVulnHostReports(mockHostVulnerabilityAssessment()))
	})
	// NOTE (@afiune): We purposly leave trailing spaces in this table, we need them!
	expectedTable := `
          HOST DETAILS                   VULNERABILITIES          
--------------------------------+---------------------------------
    Machine ID   51                 SEVERITY   COUNT   FIXABLE    
    Hostname     ip-10-0-1-48     -----------+-------+----------  
    External IP  1.2.3.4            Critical       0         0    
    Internal IP  10.0.1.1           High           1         1    
    Os           linux              Medium        42        38    
    Arch         arm64              Low           34         5    
    Namespace    ubuntu:18.04       Info           3         0    
    Provider     AWS                                              
    Instance ID  i-mock                                           
    AMI          ami-mock                                         
                                                                  
There are no critical vulnerabilities in your environment.

Try adding '--active' to only show vulnerabilities of packages actively running.

0 of 80 cve(s) showing
`
	assert.Equal(t, strings.TrimPrefix(expectedTable, "\n"), cliOutput)
}

func mockHostVulnerabilityAssessment() api.HostVulnHostAssessment {
	assessment := api.HostVulnHostAssessment{}
	err := json.Unmarshal([]byte(`{
  "host": {
    "hostname": "ip-10-0-1-48",
    "machine_id": "51",
    "tags": {
      "Account": "mocked-account",
      "AmiId": "ami-mock",
      "ExternalIp": "1.2.3.4",
      "Hostname": "ip-10-0-1-48.us-west-2.compute.internal",
      "InstanceId": "i-mock",
      "InternalIp": "10.0.1.1",
      "LwTokenShort": "abc123",
      "SubnetId": "subnet-mock",
      "VmInstanceType": "t4g.nano",
      "VmProvider": "AWS",
      "VpcId": "vpc-mock",
      "Zone": "us-west-2",
      "arch": "arm64",
      "os": "linux"
    }
  },
  "vulnerabilities": [
    {
      "cve_id": "CVE-2021-3875",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2021-3875",
          "cvss_score": "5.5",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "vim is vulnerable to Heap-based Buffer Overflow",
          "first_seen_time": "2021-10-24T03:00:00Z",
          "fix_available": "0",
          "fixed_version": "",
          "host_count": "1",
          "name": "vim",
          "namespace": "ubuntu:18.04",
          "package_status": "",
          "severity": "Medium",
          "version": "",
          "vulnerability_status": "Fixed"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": null,
          "Medium": {
            "fixable": 0,
            "vulnerabilities": 1
          }
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2021-33287",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2021-33287",
          "cvss_score": "7.8",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "In NTFS-3G versions < 2021.8.22, when specially crafted NTFS attributes are read in the function ntfs_attr_pread_i, a heap buffer overflow can occur and allow for writing to arbitrary memory or denial of service of the application.",
          "first_seen_time": "2021-09-09T03:00:00Z",
          "fix_available": "1",
          "fixed_version": "1:2017.3.23-2ubuntu0.18.04.3",
          "host_count": "1",
          "name": "ntfs-3g",
          "namespace": "ubuntu:18.04",
          "package_status": "",
          "severity": "Medium",
          "version": "",
          "vulnerability_status": "Active"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": null,
          "Medium": {
            "fixable": 1,
            "vulnerabilities": 1
          }
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2021-35266",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2021-35266",
          "cvss_score": "7.8",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "In NTFS-3G versions < 2021.8.22, when a specially crafted NTFS inode pathname is supplied in an NTFS image a heap buffer overflow can occur resulting in memory disclosure, denial of service and even code execution.",
          "first_seen_time": "2021-09-09T03:00:00Z",
          "fix_available": "1",
          "fixed_version": "1:2017.3.23-2ubuntu0.18.04.3",
          "host_count": "1",
          "name": "ntfs-3g",
          "namespace": "ubuntu:18.04",
          "package_status": "",
          "severity": "Medium",
          "version": "",
          "vulnerability_status": "Active"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": null,
          "Medium": {
            "fixable": 1,
            "vulnerabilities": 1
          }
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2019-18276",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2019-18276",
          "cvss_score": "7.8",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "An issue was discovered in disable_priv_mode in shell.c in GNU Bash through 5.0 patch 11. By default, if Bash is run with its effective UID not equal to its real UID, it will drop privileges by setting its effective UID to its real UID. However, it does so incorrectly. On Linux and other systems that support \\\\saved UID\\\\ functionality, the saved UID is not dropped. An attacker with command execution in the shell can use \\\\enable -f\\\\ for runtime loading of a new builtin, which can be a shared object that calls setuid() and therefore regains privileges. However, binaries running with an effective UID of 0 are unaffected.",
          "first_seen_time": "2021-01-05T11:00:00Z",
          "fix_available": "0",
          "fixed_version": "",
          "host_count": "1",
          "name": "bash",
          "namespace": "ubuntu:18.04",
          "package_status": "",
          "severity": "Low",
          "version": "",
          "vulnerability_status": "Active"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": {
            "fixable": 0,
            "vulnerabilities": 1
          },
          "Medium": null
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2021-39254",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2021-39254",
          "cvss_score": "7.8",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "A crafted NTFS image can cause an integer overflow in memmove, leading to a heap-based buffer overflow in the function ntfs_attr_record_resize, in NTFS-3G < 2021.8.22.",
          "first_seen_time": "2021-09-09T03:00:00Z",
          "fix_available": "1",
          "fixed_version": "1:2017.3.23-2ubuntu0.18.04.3",
          "host_count": "1",
          "name": "ntfs-3g",
          "namespace": "ubuntu:18.04",
          "package_status": "",
          "severity": "Medium",
          "version": "",
          "vulnerability_status": "Active"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": null,
          "Medium": {
            "fixable": 1,
            "vulnerabilities": 1
          }
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2017-3204",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2017-3204",
          "cvss_score": "8.1",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "The Go SSH library (x/crypto/ssh) by default does not verify host keys, facilitating man-in-the-middle attacks. Default behavior changed in commit e4e2799 to require explicitly registering a hostkey verification mechanism.",
          "first_seen_time": "2021-02-11T03:00:00Z",
          "fix_available": "0",
          "fixed_version": "",
          "host_count": "1",
          "name": "snapd",
          "namespace": "ubuntu:18.04",
          "package_status": "",
          "severity": "Low",
          "version": "",
          "vulnerability_status": "Active"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": {
            "fixable": 0,
            "vulnerabilities": 1
          },
          "Medium": null
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2021-3710",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2021-3710",
          "cvss_score": "5.5",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "An information disclosure via path traversal was discovered in apport/hookutils.py function read_file(). This issue affects: apport 2.14.1 versions prior to 2.14.1-0ubuntu3.29+esm8; 2.20.1 versions prior to 2.20.1-0ubuntu2.30+esm2; 2.20.9 versions prior to 2.20.9-0ubuntu7.26; 2.20.11 versions prior to 2.20.11-0ubuntu27.20; 2.20.11 versions prior to 2.20.11-0ubuntu65.3;",
          "first_seen_time": "2021-09-15T03:00:00Z",
          "fix_available": "1",
          "fixed_version": "2.20.9-0ubuntu7.26",
          "host_count": "1",
          "name": "apport",
          "namespace": "ubuntu:18.04",
          "package_status": "",
          "severity": "Medium",
          "version": "",
          "vulnerability_status": "Active"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": null,
          "Medium": {
            "fixable": 1,
            "vulnerabilities": 1
          }
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2019-17042",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2019-17042",
          "cvss_score": "9.8",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "An issue was discovered in Rsyslog v8.1908.0. contrib/pmcisconames/pmcisconames.c has a heap overflow in the parser for Cisco log messages. The parser tries to locate a log message delimiter (in this case, a space or a colon), but fails to account for strings that do not satisfy this constraint. If the string does not match, then the variable lenMsg will reach the value zero and will skip the sanity check that detects invalid log messages. The message will then be considered valid, and the parser will eat up the nonexistent colon delimiter. In doing so, it will decrement lenMsg, a signed integer, whose value was zero and now becomes minus one. The following step in the parser is to shift left the contents of the message. To do this, it will call memmove with the right pointers to the target and destination strings, but the lenMsg will now be interpreted as a huge value, causing a heap overflow.",
          "first_seen_time": "2021-01-05T11:00:00Z",
          "fix_available": "0",
          "fixed_version": "",
          "host_count": "1",
          "name": "rsyslog",
          "namespace": "ubuntu:18.04",
          "package_status": "ACTIVE",
          "severity": "Low",
          "version": "",
          "vulnerability_status": "Active"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": {
            "fixable": 0,
            "vulnerabilities": 1
          },
          "Medium": null
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2021-3429",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2021-3429",
          "cvss_score": "0",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "When instructing cloud-init to set a random password for a new user account, versions before 21.1.19 would write that password to the world-readable log file /var/log/cloud-init-output.log. This could allow a local user to log in as another user.",
          "first_seen_time": "2021-03-28T03:00:00Z",
          "fix_available": "1",
          "fixed_version": "21.1-19-gbad84ad4-0ubuntu1~18.04.1",
          "host_count": "1",
          "name": "cloud-init",
          "namespace": "ubuntu:18.04",
          "package_status": "",
          "severity": "Medium",
          "version": "",
          "vulnerability_status": "Reopened"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": null,
          "Medium": {
            "fixable": 1,
            "vulnerabilities": 1
          }
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2021-35269",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2021-35269",
          "cvss_score": "7.8",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "NTFS-3G versions < 2021.8.22, when a specially crafted NTFS attribute from the MFT is setup in the function ntfs_attr_setup_flag, a heap buffer overflow can occur allowing for code execution and escalation of privileges.",
          "first_seen_time": "2021-09-09T03:00:00Z",
          "fix_available": "1",
          "fixed_version": "1:2017.3.23-2ubuntu0.18.04.3",
          "host_count": "1",
          "name": "ntfs-3g",
          "namespace": "ubuntu:18.04",
          "package_status": "",
          "severity": "Medium",
          "version": "",
          "vulnerability_status": "Active"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": null,
          "Medium": {
            "fixable": 1,
            "vulnerabilities": 1
          }
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2021-23336",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2021-23336",
          "cvss_score": "5.9",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "The package python/cpython from 0 and before 3.6.13, from 3.7.0 and before 3.7.10, from 3.8.0 and before 3.8.8, from 3.9.0 and before 3.9.2 are vulnerable to Web Cache Poisoning via urllib.parse.parse_qsl and urllib.parse.parse_qs by using a vector called parameter cloaking. When the attacker can separate query parameters using a semicolon (;), they can cause a difference in the interpretation of the request between the proxy (running with default configuration) and the server. This can result in malicious requests being cached as completely safe ones, as the proxy would usually not see the semicolon as a separator, and therefore would not include it in a cache key of an unkeyed parameter.",
          "first_seen_time": "2021-03-03T03:00:00Z",
          "fix_available": "0",
          "fixed_version": "",
          "host_count": "1",
          "name": "python3.6",
          "namespace": "ubuntu:18.04",
          "package_status": "",
          "severity": "Low",
          "version": "",
          "vulnerability_status": "Active"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": {
            "fixable": 0,
            "vulnerabilities": 1
          },
          "Medium": null
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2021-20197",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2021-20197",
          "cvss_score": "6.3",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "There is an open race window when writing output in the following utilities in GNU binutils version 2.35 and earlier:ar, objcopy, strip, ranlib. When these utilities are run as a privileged user (presumably as part of a script updating binaries across different users), an unprivileged user can trick these utilities into getting ownership of arbitrary files through a symlink.",
          "first_seen_time": "2021-09-30T03:00:00Z",
          "fix_available": "0",
          "fixed_version": "",
          "host_count": "1",
          "name": "binutils",
          "namespace": "ubuntu:18.04",
          "package_status": "",
          "severity": "Low",
          "version": "",
          "vulnerability_status": "Active"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": {
            "fixable": 0,
            "vulnerabilities": 1
          },
          "Medium": null
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2021-35268",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2021-35268",
          "cvss_score": "7.8",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "In NTFS-3G versions < 2021.8.22, when a specially crafted NTFS inode is loaded in the function ntfs_inode_real_open, a heap buffer overflow can occur allowing for code execution and escalation of privileges.",
          "first_seen_time": "2021-09-09T03:00:00Z",
          "fix_available": "1",
          "fixed_version": "1:2017.3.23-2ubuntu0.18.04.3",
          "host_count": "1",
          "name": "ntfs-3g",
          "namespace": "ubuntu:18.04",
          "package_status": "",
          "severity": "Medium",
          "version": "",
          "vulnerability_status": "Active"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": null,
          "Medium": {
            "fixable": 1,
            "vulnerabilities": 1
          }
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2020-13776",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2020-13776",
          "cvss_score": "6.7",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "systemd through v245 mishandles numerical usernames such as ones composed of decimal digits or 0x followed by hex digits, as demonstrated by use of root privileges when privileges of the 0x0 user account were intended. NOTE: this issue exists because of an incomplete fix for CVE-2017-1000082.",
          "first_seen_time": "2021-07-22T03:00:00Z",
          "fix_available": "0",
          "fixed_version": "",
          "host_count": "1",
          "name": "systemd",
          "namespace": "ubuntu:18.04",
          "package_status": "ACTIVE",
          "severity": "Low",
          "version": "",
          "vulnerability_status": "Active"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": {
            "fixable": 0,
            "vulnerabilities": 1
          },
          "Medium": null
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2021-39260",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2021-39260",
          "cvss_score": "7.8",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "A crafted NTFS image can cause an out-of-bounds access in ntfs_inode_sync_standard_information in NTFS-3G < 2021.8.22.",
          "first_seen_time": "2021-09-09T03:00:00Z",
          "fix_available": "1",
          "fixed_version": "1:2017.3.23-2ubuntu0.18.04.3",
          "host_count": "1",
          "name": "ntfs-3g",
          "namespace": "ubuntu:18.04",
          "package_status": "",
          "severity": "Medium",
          "version": "",
          "vulnerability_status": "Active"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": null,
          "Medium": {
            "fixable": 1,
            "vulnerabilities": 1
          }
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2021-38185",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2021-38185",
          "cvss_score": "7.8",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "GNU cpio through 2.13 allows attackers to execute arbitrary code via a crafted pattern file, because of a dstring.c ds_fgetstr integer overflow that triggers an out-of-bounds heap write. NOTE: it is unclear whether there are common cases where the pattern file, associated with the -E option, is untrusted data.",
          "first_seen_time": "2021-08-20T03:00:00Z",
          "fix_available": "1",
          "fixed_version": "2.12+dfsg-6ubuntu0.18.04.4",
          "host_count": "1",
          "name": "cpio",
          "namespace": "ubuntu:18.04",
          "package_status": "",
          "severity": "Medium",
          "version": "",
          "vulnerability_status": "Active"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": null,
          "Medium": {
            "fixable": 1,
            "vulnerabilities": 1
          }
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2021-33285",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2021-33285",
          "cvss_score": "7.8",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "In NTFS-3G versions < 2021.8.22, when a specially crafted NTFS attribute is supplied to the function ntfs_get_attribute_value, a heap buffer overflow can occur allowing for memory disclosure or denial of service. The vulnerability is caused by an out-of-bound buffer access which can be triggered by mounting a crafted ntfs partition. The root cause is a missing consistency check after reading an MFT record : the \\\\bytes_in_use\\\\ field should be less than the \\\\bytes_allocated\\\\ field. When it is not, the parsing of the records proceeds into the wild.",
          "first_seen_time": "2021-09-09T03:00:00Z",
          "fix_available": "1",
          "fixed_version": "1:2017.3.23-2ubuntu0.18.04.3",
          "host_count": "1",
          "name": "ntfs-3g",
          "namespace": "ubuntu:18.04",
          "package_status": "",
          "severity": "Medium",
          "version": "",
          "vulnerability_status": "Active"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": null,
          "Medium": {
            "fixable": 1,
            "vulnerabilities": 1
          }
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2020-24977",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2020-24977",
          "cvss_score": "6.5",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "GNOME project libxml2 v2.9.10 has a global buffer over-read vulnerability in xmlEncodeEntitiesInternal at libxml2/entities.c. The issue has been fixed in commit 50f06b3e.",
          "first_seen_time": "2021-06-20T03:00:00Z",
          "fix_available": "0",
          "fixed_version": "",
          "host_count": "1",
          "name": "libxml2",
          "namespace": "ubuntu:18.04",
          "package_status": "",
          "severity": "Low",
          "version": "",
          "vulnerability_status": "Fixed"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": {
            "fixable": 0,
            "vulnerabilities": 1
          },
          "Medium": null
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2021-3927",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2021-3927",
          "cvss_score": "7.8",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "vim is vulnerable to Heap-based Buffer Overflow",
          "first_seen_time": "2021-11-12T03:00:00Z",
          "fix_available": "1",
          "fixed_version": "2:8.0.1453-1ubuntu1.7",
          "host_count": "1",
          "name": "vim",
          "namespace": "ubuntu:18.04",
          "package_status": "",
          "severity": "Medium",
          "version": "",
          "vulnerability_status": "Active"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": null,
          "Medium": {
            "fixable": 1,
            "vulnerabilities": 1
          }
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2017-9525",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2017-9525",
          "cvss_score": "6.7",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "In the cron package through 3.0pl1-128 on Debian, and through 3.0pl1-128ubuntu2 on Ubuntu, the postinst maintainer script allows for group-crontab-to-root privilege escalation via symlink attacks against unsafe usage of the chown and chmod programs.",
          "first_seen_time": "2021-01-05T11:00:00Z",
          "fix_available": "0",
          "fixed_version": "",
          "host_count": "1",
          "name": "cron",
          "namespace": "ubuntu:18.04",
          "package_status": "",
          "severity": "Low",
          "version": "",
          "vulnerability_status": "Active"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": {
            "fixable": 0,
            "vulnerabilities": 1
          },
          "Medium": null
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2019-17041",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2019-17041",
          "cvss_score": "9.8",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "An issue was discovered in Rsyslog v8.1908.0. contrib/pmaixforwardedfrom/pmaixforwardedfrom.c has a heap overflow in the parser for AIX log messages. The parser tries to locate a log message delimiter (in this case, a space or a colon) but fails to account for strings that do not satisfy this constraint. If the string does not match, then the variable lenMsg will reach the value zero and will skip the sanity check that detects invalid log messages. The message will then be considered valid, and the parser will eat up the nonexistent colon delimiter. In doing so, it will decrement lenMsg, a signed integer, whose value was zero and now becomes minus one. The following step in the parser is to shift left the contents of the message. To do this, it will call memmove with the right pointers to the target and destination strings, but the lenMsg will now be interpreted as a huge value, causing a heap overflow.",
          "first_seen_time": "2021-01-05T11:00:00Z",
          "fix_available": "0",
          "fixed_version": "",
          "host_count": "1",
          "name": "rsyslog",
          "namespace": "ubuntu:18.04",
          "package_status": "ACTIVE",
          "severity": "Low",
          "version": "",
          "vulnerability_status": "Active"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": {
            "fixable": 0,
            "vulnerabilities": 1
          },
          "Medium": null
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2021-39258",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2021-39258",
          "cvss_score": "7.8",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "A crafted NTFS image can cause out-of-bounds reads in ntfs_attr_find and ntfs_external_attr_find in NTFS-3G < 2021.8.22.",
          "first_seen_time": "2021-09-09T03:00:00Z",
          "fix_available": "1",
          "fixed_version": "1:2017.3.23-2ubuntu0.18.04.3",
          "host_count": "1",
          "name": "ntfs-3g",
          "namespace": "ubuntu:18.04",
          "package_status": "",
          "severity": "Medium",
          "version": "",
          "vulnerability_status": "Active"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": null,
          "Medium": {
            "fixable": 1,
            "vulnerabilities": 1
          }
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2021-39255",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2021-39255",
          "cvss_score": "7.8",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "A crafted NTFS image can trigger an out-of-bounds read, caused by an invalid attribute in ntfs_attr_find_in_attrdef, in NTFS-3G < 2021.8.22.",
          "first_seen_time": "2021-09-09T03:00:00Z",
          "fix_available": "1",
          "fixed_version": "1:2017.3.23-2ubuntu0.18.04.3",
          "host_count": "1",
          "name": "ntfs-3g",
          "namespace": "ubuntu:18.04",
          "package_status": "",
          "severity": "Medium",
          "version": "",
          "vulnerability_status": "Active"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": null,
          "Medium": {
            "fixable": 1,
            "vulnerabilities": 1
          }
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2019-0816",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2019-0816",
          "cvss_score": "5.1",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "A security feature bypass exists in Azure SSH Keypairs, due to a change in the provisioning logic for some Linux images that use cloud-init, aka 'Azure SSH Keypairs Security Feature Bypass Vulnerability'.",
          "first_seen_time": "2021-01-05T11:00:00Z",
          "fix_available": "1",
          "fixed_version": "18.5-45-g3554ffe8-0ubuntu1~18.04.1",
          "host_count": "1",
          "name": "cloud-init",
          "namespace": "ubuntu:18.04",
          "package_status": "",
          "severity": "Medium",
          "version": "",
          "vulnerability_status": "Fixed"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": null,
          "Medium": {
            "fixable": 1,
            "vulnerabilities": 1
          }
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2021-3487",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2021-3487",
          "cvss_score": "6.5",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "There's a flaw in the BFD library of binutils in versions before 2.36. An attacker who supplies a crafted file to an application linked with BFD, and using the DWARF functionality, could cause an impact to system availability by way of excessive memory consumption.",
          "first_seen_time": "2021-05-05T03:00:00Z",
          "fix_available": "1",
          "fixed_version": "2.30-21ubuntu1~18.04.7",
          "host_count": "1",
          "name": "binutils",
          "namespace": "ubuntu:18.04",
          "package_status": "",
          "severity": "Low",
          "version": "",
          "vulnerability_status": "Active"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": {
            "fixable": 1,
            "vulnerabilities": 1
          },
          "Medium": null
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2021-3537",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2021-3537",
          "cvss_score": "5.9",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "A vulnerability found in libxml2 in versions before 2.9.11 shows that it did not propagate errors while parsing XML mixed content, causing a NULL dereference. If an untrusted XML document was parsed in recovery mode and post-validated, the flaw could be used to crash the application. The highest threat from this vulnerability is to system availability.",
          "first_seen_time": "2021-06-20T03:00:00Z",
          "fix_available": "0",
          "fixed_version": "",
          "host_count": "1",
          "name": "libxml2",
          "namespace": "ubuntu:18.04",
          "package_status": "",
          "severity": "Medium",
          "version": "",
          "vulnerability_status": "Fixed"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": null,
          "Medium": {
            "fixable": 0,
            "vulnerabilities": 1
          }
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2015-4625",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2015-4625",
          "cvss_score": "4.6",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "Integer overflow in the authentication_agent_new_cookie function in PolicyKit (aka polkit) before 0.113 allows local users to gain privileges by creating a large number of connections, which triggers the issuance of a duplicate cookie value.",
          "first_seen_time": "2021-01-05T11:00:00Z",
          "fix_available": "1",
          "fixed_version": "0.105-11ubuntu1",
          "host_count": "1",
          "name": "policykit-1",
          "namespace": "ubuntu:18.04",
          "package_status": "",
          "severity": "Low",
          "version": "",
          "vulnerability_status": "Fixed"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": {
            "fixable": 1,
            "vulnerabilities": 1
          },
          "Medium": null
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2021-3516",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2021-3516",
          "cvss_score": "7.8",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "There's a flaw in libxml2's xmllint in versions before 2.9.11. An attacker who is able to submit a crafted file to be processed by xmllint could trigger a use-after-free. The greatest impact of this flaw is to confidentiality, integrity, and availability.",
          "first_seen_time": "2021-06-20T03:00:00Z",
          "fix_available": "0",
          "fixed_version": "",
          "host_count": "1",
          "name": "libxml2",
          "namespace": "ubuntu:18.04",
          "package_status": "",
          "severity": "Medium",
          "version": "",
          "vulnerability_status": "Fixed"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": null,
          "Medium": {
            "fixable": 0,
            "vulnerabilities": 1
          }
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2021-3712",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2021-3712",
          "cvss_score": "7.4",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "ASN.1 strings are represented internally within OpenSSL as an ASN1_STRING structure which contains a buffer holding the string data and a field holding the buffer length. This contrasts with normal C strings which are repesented as a buffer for the string data which is terminated with a NUL (0) byte. Although not a strict requirement, ASN.1 strings that are parsed using OpenSSL's own \\\\d2i\\\\ functions (and other similar parsing functions) as well as any string whose value has been set with the ASN1_STRING_set() function will additionally NUL terminate the byte array in the ASN1_STRING structure. However, it is possible for applications to directly construct valid ASN1_STRING structures which do not NUL terminate the byte array by directly setting the \\\\data\\\\ and \\\\length\\\\ fields in the ASN1_STRING array. This can also happen by using the ASN1_STRING_set0() function. Numerous OpenSSL functions that print ASN.1 data have been found to assume that the ASN1_STRING byte array will be NUL terminated, even though this is not guaranteed for strings that have been directly constructed. Where an application requests an ASN.1 structure to be printed, and where that ASN.1 structure contains ASN1_STRINGs that have been directly constructed by the application without NUL terminating the \\\\data\\\\ field, then a read buffer overrun can occur. The same thing can also occur during name constraints processing of certificates (for example if a certificate has been directly constructed by the application instead of loading it via the OpenSSL parsing functions, and the certificate contains non NUL terminated ASN1_STRING structures). It can also occur in the X509_get1_email(), X509_REQ_get1_email() and X509_get1_ocsp() functions. If a malicious actor can cause an application to directly construct an ASN1_STRING and then process it through one of the affected OpenSSL functions then this issue could be hit. This might result in a crash (causing a Denial of Service attack). It could also result in the disclosure of private memory contents (such as private keys, or sensitive plaintext). Fixed in OpenSSL 1.1.1l (Affected 1.1.1-1.1.1k). Fixed in OpenSSL 1.0.2za (Affected 1.0.2-1.0.2y).",
          "first_seen_time": "2021-09-01T03:00:00Z",
          "fix_available": "1",
          "fixed_version": "1.1.1-1ubuntu2.1~18.04.13",
          "host_count": "1",
          "name": "openssl",
          "namespace": "ubuntu:18.04",
          "package_status": "",
          "severity": "Medium",
          "version": "",
          "vulnerability_status": "Active"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": null,
          "Medium": {
            "fixable": 1,
            "vulnerabilities": 1
          }
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2021-39261",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2021-39261",
          "cvss_score": "7.8",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "A crafted NTFS image can cause a heap-based buffer overflow in ntfs_compressed_pwrite in NTFS-3G < 2021.8.22.",
          "first_seen_time": "2021-09-09T03:00:00Z",
          "fix_available": "1",
          "fixed_version": "1:2017.3.23-2ubuntu0.18.04.3",
          "host_count": "1",
          "name": "ntfs-3g",
          "namespace": "ubuntu:18.04",
          "package_status": "",
          "severity": "Medium",
          "version": "",
          "vulnerability_status": "Active"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": null,
          "Medium": {
            "fixable": 1,
            "vulnerabilities": 1
          }
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2018-10906",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2018-10906",
          "cvss_score": "7.8",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "In fuse before versions 2.9.8 and 3.x before 3.2.5, fusermount is vulnerable to a restriction bypass when SELinux is active. This allows non-root users to mount a FUSE file system with the 'allow_other' mount option regardless of whether 'user_allow_other' is set in the fuse configuration. An attacker may use this flaw to mount a FUSE file system, accessible by other users, and trick them into accessing files on that file system, possibly causing Denial of Service or other unspecified effects.",
          "first_seen_time": "2021-06-12T03:00:00Z",
          "fix_available": "0",
          "fixed_version": "",
          "host_count": "1",
          "name": "fuse",
          "namespace": "ubuntu:18.04",
          "package_status": "",
          "severity": "Low",
          "version": "",
          "vulnerability_status": "Active"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": {
            "fixable": 0,
            "vulnerabilities": 1
          },
          "Medium": null
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2021-3709",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2021-3709",
          "cvss_score": "5.5",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "Function check_attachment_for_errors() in file data/general-hooks/ubuntu.py could be tricked into exposing private data via a constructed crash file. This issue affects: apport 2.14.1 versions prior to 2.14.1-0ubuntu3.29+esm8; 2.20.1 versions prior to 2.20.1-0ubuntu2.30+esm2; 2.20.9 versions prior to 2.20.9-0ubuntu7.26; 2.20.11 versions prior to 2.20.11-0ubuntu27.20; 2.20.11 versions prior to 2.20.11-0ubuntu65.3;",
          "first_seen_time": "2021-09-15T03:00:00Z",
          "fix_available": "1",
          "fixed_version": "2.20.9-0ubuntu7.26",
          "host_count": "1",
          "name": "apport",
          "namespace": "ubuntu:18.04",
          "package_status": "",
          "severity": "Medium",
          "version": "",
          "vulnerability_status": "Active"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": null,
          "Medium": {
            "fixable": 1,
            "vulnerabilities": 1
          }
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2021-39259",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2021-39259",
          "cvss_score": "7.8",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "A crafted NTFS image can trigger an out-of-bounds access, caused by an unsanitized attribute length in ntfs_inode_lookup_by_name, in NTFS-3G < 2021.8.22.",
          "first_seen_time": "2021-09-09T03:00:00Z",
          "fix_available": "1",
          "fixed_version": "1:2017.3.23-2ubuntu0.18.04.3",
          "host_count": "1",
          "name": "ntfs-3g",
          "namespace": "ubuntu:18.04",
          "package_status": "",
          "severity": "Medium",
          "version": "",
          "vulnerability_status": "Active"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": null,
          "Medium": {
            "fixable": 1,
            "vulnerabilities": 1
          }
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2021-3426",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2021-3426",
          "cvss_score": "5.7",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "There's a flaw in Python 3's pydoc. A local or adjacent attacker who discovers or is able to convince another local or adjacent user to start a pydoc server could access the server and use it to disclose sensitive information belonging to the other user that they would not normally be able to access. The highest risk of this flaw is to data confidentiality. This flaw affects Python versions before 3.8.9, Python versions before 3.9.3 and Python versions before 3.10.0a7.",
          "first_seen_time": "2021-09-21T03:00:00Z",
          "fix_available": "0",
          "fixed_version": "",
          "host_count": "1",
          "name": "python3.6",
          "namespace": "ubuntu:18.04",
          "package_status": "",
          "severity": "Low",
          "version": "",
          "vulnerability_status": "Active"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": {
            "fixable": 0,
            "vulnerabilities": 1
          },
          "Medium": null
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2019-11840",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2019-11840",
          "cvss_score": "5.9",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "An issue was discovered in supplementary Go cryptography libraries, aka golang-googlecode-go-crypto, before 2019-03-20. A flaw was found in the amd64 implementation of golang.org/x/crypto/salsa20 and golang.org/x/crypto/salsa20/salsa. If more than 256 GiB of keystream is generated, or if the counter otherwise grows greater than 32 bits, the amd64 implementation will first generate incorrect output, and then cycle back to previously generated keystream. Repeated keystream bytes can lead to loss of confidentiality in encryption applications, or to predictability in CSPRNG applications.",
          "first_seen_time": "2021-02-11T03:00:00Z",
          "fix_available": "0",
          "fixed_version": "",
          "host_count": "1",
          "name": "snapd",
          "namespace": "ubuntu:18.04",
          "package_status": "",
          "severity": "Medium",
          "version": "",
          "vulnerability_status": "Active"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": null,
          "Medium": {
            "fixable": 0,
            "vulnerabilities": 1
          }
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2020-8037",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2020-8037",
          "cvss_score": "7.5",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "The ppp decapsulator in tcpdump 4.9.3 can be convinced to allocate a large amount of memory.",
          "first_seen_time": "2021-07-09T03:00:00Z",
          "fix_available": "0",
          "fixed_version": "",
          "host_count": "1",
          "name": "tcpdump",
          "namespace": "ubuntu:18.04",
          "package_status": "",
          "severity": "Low",
          "version": "",
          "vulnerability_status": "Active"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": {
            "fixable": 0,
            "vulnerabilities": 1
          },
          "Medium": null
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2021-3928",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2021-3928",
          "cvss_score": "7.8",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "vim is vulnerable to Stack-based Buffer Overflow",
          "first_seen_time": "2021-11-12T03:00:00Z",
          "fix_available": "1",
          "fixed_version": "2:8.0.1453-1ubuntu1.7",
          "host_count": "1",
          "name": "vim",
          "namespace": "ubuntu:18.04",
          "package_status": "",
          "severity": "Medium",
          "version": "",
          "vulnerability_status": "Active"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": null,
          "Medium": {
            "fixable": 1,
            "vulnerabilities": 1
          }
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2019-1010204",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2019-1010204",
          "cvss_score": "5.5",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "GNU binutils gold gold v1.11-v1.16 (GNU binutils v2.21-v2.31.1) is affected by: Improper Input Validation, Signed/Unsigned Comparison, Out-of-bounds Read. The impact is: Denial of service. The component is: gold/fileread.cc:497, elfcpp/elfcpp_file.h:644. The attack vector is: An ELF file with an invalid e_shoff header field must be opened.",
          "first_seen_time": "2021-03-04T03:00:00Z",
          "fix_available": "0",
          "fixed_version": "",
          "host_count": "1",
          "name": "binutils",
          "namespace": "ubuntu:18.04",
          "package_status": "",
          "severity": "Low",
          "version": "",
          "vulnerability_status": "Active"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": {
            "fixable": 0,
            "vulnerabilities": 1
          },
          "Medium": null
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2021-3778",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2021-3778",
          "cvss_score": "7.8",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "vim is vulnerable to Heap-based Buffer Overflow",
          "first_seen_time": "2021-09-18T03:00:00Z",
          "fix_available": "1",
          "fixed_version": "2:8.0.1453-1ubuntu1.6",
          "host_count": "1",
          "name": "vim",
          "namespace": "ubuntu:18.04",
          "package_status": "",
          "severity": "Medium",
          "version": "",
          "vulnerability_status": "Active"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": null,
          "Medium": {
            "fixable": 1,
            "vulnerabilities": 1
          }
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2021-3903",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2021-3903",
          "cvss_score": "7.8",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "vim is vulnerable to Heap-based Buffer Overflow",
          "first_seen_time": "2021-11-12T03:00:00Z",
          "fix_available": "1",
          "fixed_version": "2:8.0.1453-1ubuntu1.7",
          "host_count": "1",
          "name": "vim",
          "namespace": "ubuntu:18.04",
          "package_status": "",
          "severity": "Low",
          "version": "",
          "vulnerability_status": "Active"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": {
            "fixable": 1,
            "vulnerabilities": 1
          },
          "Medium": null
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2020-16592",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2020-16592",
          "cvss_score": "5.5",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "A use after free issue exists in the Binary File Descriptor (BFD) library (aka libbfd) in GNU Binutils 2.34 in bfd_hash_lookup, as demonstrated in nm-new, that can cause a denial of service via a crafted file.",
          "first_seen_time": "2021-10-27T03:00:00Z",
          "fix_available": "1",
          "fixed_version": "2.30-21ubuntu1~18.04.7",
          "host_count": "1",
          "name": "binutils",
          "namespace": "ubuntu:18.04",
          "package_status": "",
          "severity": "Low",
          "version": "",
          "vulnerability_status": "Active"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": {
            "fixable": 1,
            "vulnerabilities": 1
          },
          "Medium": null
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2021-33286",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2021-33286",
          "cvss_score": "7.8",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "In NTFS-3G versions < 2021.8.22, when a specially crafted unicode string is supplied in an NTFS image a heap buffer overflow can occur and allow for code execution.",
          "first_seen_time": "2021-09-09T03:00:00Z",
          "fix_available": "1",
          "fixed_version": "1:2017.3.23-2ubuntu0.18.04.3",
          "host_count": "1",
          "name": "ntfs-3g",
          "namespace": "ubuntu:18.04",
          "package_status": "",
          "severity": "Medium",
          "version": "",
          "vulnerability_status": "Active"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": null,
          "Medium": {
            "fixable": 1,
            "vulnerabilities": 1
          }
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2018-1000021",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2018-1000021",
          "cvss_score": "8.8",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "GIT version 2.15.1 and earlier contains a Input Validation Error vulnerability in Client that can result in problems including messing up terminal configuration to RCE. This attack appear to be exploitable via The user must interact with a malicious git server, (or have their traffic modified in a MITM attack).",
          "first_seen_time": "2021-03-11T03:00:00Z",
          "fix_available": "0",
          "fixed_version": "",
          "host_count": "1",
          "name": "git",
          "namespace": "ubuntu:18.04",
          "package_status": "",
          "severity": "Low",
          "version": "",
          "vulnerability_status": "Active"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": {
            "fixable": 0,
            "vulnerabilities": 1
          },
          "Medium": null
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2021-39256",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2021-39256",
          "cvss_score": "7.8",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "A crafted NTFS image can cause a heap-based buffer overflow in ntfs_inode_lookup_by_name in NTFS-3G < 2021.8.22.",
          "first_seen_time": "2021-09-09T03:00:00Z",
          "fix_available": "1",
          "fixed_version": "1:2017.3.23-2ubuntu0.18.04.3",
          "host_count": "1",
          "name": "ntfs-3g",
          "namespace": "ubuntu:18.04",
          "package_status": "",
          "severity": "Medium",
          "version": "",
          "vulnerability_status": "Active"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": null,
          "Medium": {
            "fixable": 1,
            "vulnerabilities": 1
          }
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2020-13987",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2020-13987",
          "cvss_score": "7.5",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "An issue was discovered in Contiki through 3.0. An Out-of-Bounds Read vulnerability exists in the uIP TCP/IP Stack component when calculating the checksums for IP packets in upper_layer_chksum in net/ipv4/uip.c.",
          "first_seen_time": "2021-01-23T03:00:00Z",
          "fix_available": "0",
          "fixed_version": "",
          "host_count": "1",
          "name": "open-iscsi",
          "namespace": "ubuntu:18.04",
          "package_status": "",
          "severity": "Low",
          "version": "",
          "vulnerability_status": "Active"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": {
            "fixable": 0,
            "vulnerabilities": 1
          },
          "Medium": null
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2019-20388",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2019-20388",
          "cvss_score": "7.5",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "xmlSchemaPreRun in xmlschemas.c in libxml2 2.9.10 allows an xmlSchemaValidateStream memory leak.",
          "first_seen_time": "2021-06-20T03:00:00Z",
          "fix_available": "0",
          "fixed_version": "",
          "host_count": "1",
          "name": "libxml2",
          "namespace": "ubuntu:18.04",
          "package_status": "",
          "severity": "Low",
          "version": "",
          "vulnerability_status": "Fixed"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": {
            "fixable": 0,
            "vulnerabilities": 1
          },
          "Medium": null
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2021-3518",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2021-3518",
          "cvss_score": "8.8",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "There's a flaw in libxml2 in versions before 2.9.11. An attacker who is able to submit a crafted file to be processed by an application linked with libxml2 could trigger a use-after-free. The greatest impact from this flaw is to confidentiality, integrity, and availability.",
          "first_seen_time": "2021-06-20T03:00:00Z",
          "fix_available": "0",
          "fixed_version": "",
          "host_count": "1",
          "name": "libxml2",
          "namespace": "ubuntu:18.04",
          "package_status": "",
          "severity": "Medium",
          "version": "",
          "vulnerability_status": "Fixed"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": null,
          "Medium": {
            "fixable": 0,
            "vulnerabilities": 1
          }
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2018-1000654",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2018-1000654",
          "cvss_score": "5.5",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "GNU Libtasn1-4.13 libtasn1-4.13 version libtasn1-4.13, libtasn1-4.12 contains a DoS, specifically CPU usage will reach 100% when running asn1Paser against the POC due to an issue in _asn1_expand_object_id(p_tree), after a long time, the program will be killed. This attack appears to be exploitable via parsing a crafted file.",
          "first_seen_time": "2021-01-05T11:00:00Z",
          "fix_available": "0",
          "fixed_version": "",
          "host_count": "1",
          "name": "libtasn1-6",
          "namespace": "ubuntu:18.04",
          "package_status": "",
          "severity": "Info",
          "version": "",
          "vulnerability_status": "Active"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": {
            "fixable": 0,
            "vulnerabilities": 1
          },
          "Low": null,
          "Medium": null
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2020-35512",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2020-35512",
          "cvss_score": "7.8",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "A use-after-free flaw was found in D-Bus Development branch <= 1.13.16, dbus-1.12.x stable branch <= 1.12.18, and dbus-1.10.x and older branches <= 1.10.30 when a system has multiple usernames sharing the same UID. When a set of policy rules references these usernames, D-Bus may free some memory in the heap, which is still used by data structures necessary for the other usernames sharing the UID, possibly leading to a crash or other undefined behaviors",
          "first_seen_time": "2021-07-09T03:00:00Z",
          "fix_available": "0",
          "fixed_version": "",
          "host_count": "1",
          "name": "dbus",
          "namespace": "ubuntu:18.04",
          "package_status": "ACTIVE",
          "severity": "Low",
          "version": "",
          "vulnerability_status": "Active"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": {
            "fixable": 0,
            "vulnerabilities": 1
          },
          "Medium": null
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2021-39263",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2021-39263",
          "cvss_score": "7.8",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "A crafted NTFS image can trigger a heap-based buffer overflow, caused by an unsanitized attribute in ntfs_get_attribute_value, in NTFS-3G < 2021.8.22.",
          "first_seen_time": "2021-09-09T03:00:00Z",
          "fix_available": "1",
          "fixed_version": "1:2017.3.23-2ubuntu0.18.04.3",
          "host_count": "1",
          "name": "ntfs-3g",
          "namespace": "ubuntu:18.04",
          "package_status": "",
          "severity": "Medium",
          "version": "",
          "vulnerability_status": "Active"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": null,
          "Medium": {
            "fixable": 1,
            "vulnerabilities": 1
          }
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2021-33289",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2021-33289",
          "cvss_score": "7.8",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "In NTFS-3G versions < 2021.8.22, when a specially crafted MFT section is supplied in an NTFS image a heap buffer overflow can occur and allow for code execution.",
          "first_seen_time": "2021-09-09T03:00:00Z",
          "fix_available": "1",
          "fixed_version": "1:2017.3.23-2ubuntu0.18.04.3",
          "host_count": "1",
          "name": "ntfs-3g",
          "namespace": "ubuntu:18.04",
          "package_status": "",
          "severity": "Medium",
          "version": "",
          "vulnerability_status": "Active"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": null,
          "Medium": {
            "fixable": 1,
            "vulnerabilities": 1
          }
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2018-20482",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2018-20482",
          "cvss_score": "4.7",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "GNU Tar through 1.30, when --sparse is used, mishandles file shrinkage during read access, which allows local users to cause a denial of service (infinite read loop in sparse_dump_region in sparse.c) by modifying a file that is supposed to be archived by a different user's process (e.g., a system backup running as root).",
          "first_seen_time": "2021-01-15T03:00:00Z",
          "fix_available": "0",
          "fixed_version": "",
          "host_count": "1",
          "name": "tar",
          "namespace": "ubuntu:18.04",
          "package_status": "",
          "severity": "Low",
          "version": "",
          "vulnerability_status": "Fixed"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": {
            "fixable": 0,
            "vulnerabilities": 1
          },
          "Medium": null
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2019-9072",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2019-9072",
          "cvss_score": "5.5",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "An issue was discovered in the Binary File Descriptor (BFD) library (aka libbfd), as distributed in GNU Binutils 2.32. It is an attempted excessive memory allocation in setup_group in elf.c.",
          "first_seen_time": "2021-03-04T03:00:00Z",
          "fix_available": "0",
          "fixed_version": "",
          "host_count": "1",
          "name": "binutils",
          "namespace": "ubuntu:18.04",
          "package_status": "",
          "severity": "Low",
          "version": "",
          "vulnerability_status": "Active"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": {
            "fixable": 0,
            "vulnerabilities": 1
          },
          "Medium": null
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2021-40153",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2021-40153",
          "cvss_score": "8.1",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "squashfs_opendir in unsquash-1.c in Squashfs-Tools 4.5 stores the filename in the directory entry; this is then used by unsquashfs to create the new file during the unsquash. The filename is not validated for traversal outside of the destination directory, and thus allows writing to locations outside of the destination.",
          "first_seen_time": "2021-09-01T03:00:00Z",
          "fix_available": "1",
          "fixed_version": "1:4.3-6ubuntu0.18.04.3",
          "host_count": "1",
          "name": "squashfs-tools",
          "namespace": "ubuntu:18.04",
          "package_status": "",
          "severity": "Medium",
          "version": "",
          "vulnerability_status": "Active"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": null,
          "Medium": {
            "fixable": 1,
            "vulnerabilities": 1
          }
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2021-22898",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2021-22898",
          "cvss_score": "3.1",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "curl 7.7 through 7.76.1 suffers from an information disclosure when the '-t' command line option, known as 'CURLOPT_TELNETOPTIONS' in libcurl, is used to send variable=content pairs to TELNET servers. Due to a flaw in the option parser for sending NEW_ENV variables, libcurl could be made to pass on uninitialized data from a stack based buffer to the server, resulting in potentially revealing sensitive internal information to the server using a clear-text network protocol.",
          "first_seen_time": "2021-06-12T03:00:00Z",
          "fix_available": "1",
          "fixed_version": "7.58.0-2ubuntu3.14",
          "host_count": "1",
          "name": "curl",
          "namespace": "ubuntu:18.04",
          "package_status": "",
          "severity": "Low",
          "version": "",
          "vulnerability_status": "Active"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": {
            "fixable": 1,
            "vulnerabilities": 1
          },
          "Medium": null
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2021-22925",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2021-22925",
          "cvss_score": "5.3",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "curl supports the '-t' command line option, known as 'CURLOPT_TELNETOPTIONS'in libcurl. This rarely used option is used to send variable=content pairs toTELNET servers.Due to flaw in the option parser for sending 'NEW_ENV' variables, libcurlcould be made to pass on uninitialized data from a stack based buffer to theserver. Therefore potentially revealing sensitive internal information to theserver using a clear-text network protocol.This could happen because curl did not call and use sscanf() correctly whenparsing the string provided by the application.",
          "first_seen_time": "2021-08-04T03:00:00Z",
          "fix_available": "1",
          "fixed_version": "7.58.0-2ubuntu3.14",
          "host_count": "1",
          "name": "curl",
          "namespace": "ubuntu:18.04",
          "package_status": "",
          "severity": "Medium",
          "version": "",
          "vulnerability_status": "Active"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": null,
          "Medium": {
            "fixable": 1,
            "vulnerabilities": 1
          }
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2021-40330",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2021-40330",
          "cvss_score": "7.5",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "git_connect_git in connect.c in Git before 2.30.1 allows a repository path to contain a newline character, which may result in unexpected cross-protocol requests, as demonstrated by the git://localhost:1234/%0d%0a%0d%0aGET%20/%20HTTP/1.1 substring.",
          "first_seen_time": "2021-09-09T03:00:00Z",
          "fix_available": "1",
          "fixed_version": "1:2.17.1-1ubuntu0.9",
          "host_count": "1",
          "name": "git",
          "namespace": "ubuntu:18.04",
          "package_status": "",
          "severity": "Medium",
          "version": "",
          "vulnerability_status": "Active"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": null,
          "Medium": {
            "fixable": 1,
            "vulnerabilities": 1
          }
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2021-3155",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2021-3155",
          "cvss_score": "0",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "snapd does not enforce that the per-snap user data directory ~/snap/<snap-name> is private. This could expose sensitive secrets or tokens to other local users",
          "first_seen_time": "2021-06-12T03:00:00Z",
          "fix_available": "0",
          "fixed_version": "",
          "host_count": "1",
          "name": "snapd",
          "namespace": "ubuntu:18.04",
          "package_status": "",
          "severity": "Medium",
          "version": "",
          "vulnerability_status": "Active"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": null,
          "Medium": {
            "fixable": 0,
            "vulnerabilities": 1
          }
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2021-39252",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2021-39252",
          "cvss_score": "7.8",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "A crafted NTFS image can cause an out-of-bounds read in ntfs_ie_lookup in NTFS-3G < 2021.8.22.",
          "first_seen_time": "2021-09-09T03:00:00Z",
          "fix_available": "1",
          "fixed_version": "1:2017.3.23-2ubuntu0.18.04.3",
          "host_count": "1",
          "name": "ntfs-3g",
          "namespace": "ubuntu:18.04",
          "package_status": "",
          "severity": "Medium",
          "version": "",
          "vulnerability_status": "Active"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": null,
          "Medium": {
            "fixable": 1,
            "vulnerabilities": 1
          }
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2021-33560",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2021-33560",
          "cvss_score": "7.5",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "Libgcrypt before 1.8.8 and 1.9.x before 1.9.3 mishandles ElGamal encryption because it lacks exponent blinding to address a side-channel attack against mpi_powm, and the window size is not chosen appropriately. This, for example, affects use of ElGamal in OpenPGP.",
          "first_seen_time": "2021-06-12T03:00:00Z",
          "fix_available": "1",
          "fixed_version": "1.8.1-4ubuntu1.3",
          "host_count": "1",
          "name": "libgcrypt20",
          "namespace": "ubuntu:18.04",
          "package_status": "",
          "severity": "Low",
          "version": "",
          "vulnerability_status": "Active"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": {
            "fixable": 1,
            "vulnerabilities": 1
          },
          "Medium": null
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2018-6952",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2018-6952",
          "cvss_score": "7.5",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "A double free exists in the another_hunk function in pch.c in GNU patch through 2.7.6.",
          "first_seen_time": "2021-01-05T11:00:00Z",
          "fix_available": "0",
          "fixed_version": "",
          "host_count": "1",
          "name": "patch",
          "namespace": "ubuntu:18.04",
          "package_status": "",
          "severity": "Info",
          "version": "",
          "vulnerability_status": "Active"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": {
            "fixable": 0,
            "vulnerabilities": 1
          },
          "Low": null,
          "Medium": null
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2021-39262",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2021-39262",
          "cvss_score": "7.8",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "A crafted NTFS image can cause an out-of-bounds access in ntfs_decompress in NTFS-3G < 2021.8.22.",
          "first_seen_time": "2021-09-09T03:00:00Z",
          "fix_available": "1",
          "fixed_version": "1:2017.3.23-2ubuntu0.18.04.3",
          "host_count": "1",
          "name": "ntfs-3g",
          "namespace": "ubuntu:18.04",
          "package_status": "",
          "severity": "Medium",
          "version": "",
          "vulnerability_status": "Active"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": null,
          "Medium": {
            "fixable": 1,
            "vulnerabilities": 1
          }
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2021-39257",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2021-39257",
          "cvss_score": "5.5",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "A crafted NTFS image with an unallocated bitmap can lead to a endless recursive function call chain (starting from ntfs_attr_pwrite), causing stack consumption in NTFS-3G < 2021.8.22.",
          "first_seen_time": "2021-09-09T03:00:00Z",
          "fix_available": "1",
          "fixed_version": "1:2017.3.23-2ubuntu0.18.04.3",
          "host_count": "1",
          "name": "ntfs-3g",
          "namespace": "ubuntu:18.04",
          "package_status": "",
          "severity": "Medium",
          "version": "",
          "vulnerability_status": "Active"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": null,
          "Medium": {
            "fixable": 1,
            "vulnerabilities": 1
          }
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2018-20657",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2018-20657",
          "cvss_score": "7.5",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "The demangle_template function in cplus-dem.c in GNU libiberty, as distributed in GNU Binutils 2.31.1, has a memory leak via a crafted string, leading to a denial of service (memory consumption), as demonstrated by cxxfilt, a related issue to CVE-2018-12698.",
          "first_seen_time": "2021-03-04T03:00:00Z",
          "fix_available": "0",
          "fixed_version": "",
          "host_count": "1",
          "name": "binutils",
          "namespace": "ubuntu:18.04",
          "package_status": "",
          "severity": "Info",
          "version": "",
          "vulnerability_status": "Active"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": {
            "fixable": 0,
            "vulnerabilities": 1
          },
          "Low": null,
          "Medium": null
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2017-1000382",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2017-1000382",
          "cvss_score": "5.5",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "VIM version 8.0.1187 (and other versions most likely) ignores umask when creating a swap file (\\\\[ORIGINAL_FILENAME].swp\\\\) resulting in files that may be world readable or otherwise accessible in ways not intended by the user running the vi binary.",
          "first_seen_time": "2021-01-05T11:00:00Z",
          "fix_available": "0",
          "fixed_version": "",
          "host_count": "1",
          "name": "vim",
          "namespace": "ubuntu:18.04",
          "package_status": "",
          "severity": "Low",
          "version": "",
          "vulnerability_status": "Active"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": {
            "fixable": 0,
            "vulnerabilities": 1
          },
          "Medium": null
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2020-17437",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2020-17437",
          "cvss_score": "8.2",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "An issue was discovered in uIP 1.0, as used in Contiki 3.0 and other products. When the Urgent flag is set in a TCP packet, and the stack is configured to ignore the urgent data, the stack attempts to use the value of the Urgent pointer bytes to separate the Urgent data from the normal data, by calculating the offset at which the normal data should be present in the global buffer. However, the length of this offset is not checked; therefore, for large values of the Urgent pointer bytes, the data pointer can point to memory that is way beyond the data buffer in uip_process in uip.c.",
          "first_seen_time": "2021-01-23T03:00:00Z",
          "fix_available": "0",
          "fixed_version": "",
          "host_count": "1",
          "name": "open-iscsi",
          "namespace": "ubuntu:18.04",
          "package_status": "",
          "severity": "Low",
          "version": "",
          "vulnerability_status": "Active"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": {
            "fixable": 0,
            "vulnerabilities": 1
          },
          "Medium": null
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2019-7306",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2019-7306",
          "cvss_score": "7.5",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "Byobu Apport hook may disclose sensitive information since it automatically uploads the local user's .screenrc which may contain private hostnames, usernames and passwords. This issue affects: byobu",
          "first_seen_time": "2021-01-05T11:00:00Z",
          "fix_available": "0",
          "fixed_version": "",
          "host_count": "1",
          "name": "byobu",
          "namespace": "ubuntu:18.04",
          "package_status": "",
          "severity": "Low",
          "version": "",
          "vulnerability_status": "Active"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": {
            "fixable": 0,
            "vulnerabilities": 1
          },
          "Medium": null
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2021-3796",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2021-3796",
          "cvss_score": "7.3",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "vim is vulnerable to Use After Free",
          "first_seen_time": "2021-09-18T03:00:00Z",
          "fix_available": "1",
          "fixed_version": "2:8.0.1453-1ubuntu1.6",
          "host_count": "1",
          "name": "vim",
          "namespace": "ubuntu:18.04",
          "package_status": "",
          "severity": "Medium",
          "version": "",
          "vulnerability_status": "Active"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": null,
          "Medium": {
            "fixable": 1,
            "vulnerabilities": 1
          }
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2017-15131",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2017-15131",
          "cvss_score": "7.8",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "It was found that system umask policy is not being honored when creating XDG user directories, since Xsession sources xdg-user-dirs.sh before setting umask policy. This only affects xdg-user-dirs before 0.15.5 as shipped with Red Hat Enterprise Linux.",
          "first_seen_time": "2021-01-05T11:00:00Z",
          "fix_available": "0",
          "fixed_version": "",
          "host_count": "1",
          "name": "xdg-user-dirs",
          "namespace": "ubuntu:18.04",
          "package_status": "",
          "severity": "Low",
          "version": "",
          "vulnerability_status": "Active"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": {
            "fixable": 0,
            "vulnerabilities": 1
          },
          "Medium": null
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2021-22924",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2021-22924",
          "cvss_score": "3.7",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "libcurl keeps previously used connections in a connection pool for subsequenttransfers to reuse, if one of them matches the setup.Due to errors in the logic, the config matching function did not take 'issuercert' into account and it compared the involved paths *case insensitively*,which could lead to libcurl reusing wrong connections.File paths are, or can be, case sensitive on many systems but not all, and caneven vary depending on used file systems.The comparison also didn't include the 'issuer cert' which a transfer can setto qualify how to verify the server certificate.",
          "first_seen_time": "2021-08-04T03:00:00Z",
          "fix_available": "1",
          "fixed_version": "7.58.0-2ubuntu3.14",
          "host_count": "1",
          "name": "curl",
          "namespace": "ubuntu:18.04",
          "package_status": "",
          "severity": "Medium",
          "version": "",
          "vulnerability_status": "Active"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": null,
          "Medium": {
            "fixable": 1,
            "vulnerabilities": 1
          }
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2021-22947",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2021-22947",
          "cvss_score": "5.9",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "When curl >= 7.20.0 and <= 7.78.0 connects to an IMAP or POP3 server to retrieve data using STARTTLS to upgrade to TLS security, the server can respond and send back multiple responses at once that curl caches. curl would then upgrade to TLS but not flush the in-queue of cached responses but instead continue using and trustingthe responses it got *before* the TLS handshake as if they were authenticated.Using this flaw, it allows a Man-In-The-Middle attacker to first inject the fake responses, then pass-through the TLS traffic from the legitimate server and trick curl into sending data back to the user thinking the attacker's injected data comes from the TLS-protected server.",
          "first_seen_time": "2021-09-16T03:00:00Z",
          "fix_available": "1",
          "fixed_version": "7.58.0-2ubuntu3.15",
          "host_count": "1",
          "name": "curl",
          "namespace": "ubuntu:18.04",
          "package_status": "",
          "severity": "Medium",
          "version": "",
          "vulnerability_status": "Active"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": null,
          "Medium": {
            "fixable": 1,
            "vulnerabilities": 1
          }
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2021-39251",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2021-39251",
          "cvss_score": "7.8",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "A crafted NTFS image can cause a NULL pointer dereference in ntfs_extent_inode_open in NTFS-3G < 2021.8.22.",
          "first_seen_time": "2021-09-09T03:00:00Z",
          "fix_available": "1",
          "fixed_version": "1:2017.3.23-2ubuntu0.18.04.3",
          "host_count": "1",
          "name": "ntfs-3g",
          "namespace": "ubuntu:18.04",
          "package_status": "",
          "severity": "Medium",
          "version": "",
          "vulnerability_status": "Active"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": null,
          "Medium": {
            "fixable": 1,
            "vulnerabilities": 1
          }
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2019-9923",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2019-9923",
          "cvss_score": "7.5",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "pax_decode_header in sparse.c in GNU Tar before 1.32 had a NULL pointer dereference when parsing certain archives that have malformed extended headers.",
          "first_seen_time": "2021-01-15T03:00:00Z",
          "fix_available": "0",
          "fixed_version": "",
          "host_count": "1",
          "name": "tar",
          "namespace": "ubuntu:18.04",
          "package_status": "",
          "severity": "Low",
          "version": "",
          "vulnerability_status": "Fixed"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": {
            "fixable": 0,
            "vulnerabilities": 1
          },
          "Medium": null
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2018-20673",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2018-20673",
          "cvss_score": "5.5",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "The demangle_template function in cplus-dem.c in GNU libiberty, as distributed in GNU Binutils 2.31.1, contains an integer overflow vulnerability (for \\\\Create an array for saving the template argument values\\\\) that can trigger a heap-based buffer overflow, as demonstrated by nm.",
          "first_seen_time": "2021-03-04T03:00:00Z",
          "fix_available": "0",
          "fixed_version": "",
          "host_count": "1",
          "name": "binutils",
          "namespace": "ubuntu:18.04",
          "package_status": "",
          "severity": "Low",
          "version": "",
          "vulnerability_status": "Active"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": {
            "fixable": 0,
            "vulnerabilities": 1
          },
          "Medium": null
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2021-22946",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2021-22946",
          "cvss_score": "7.5",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "A user can tell curl >= 7.20.0 and <= 7.78.0 to require a successful upgrade to TLS when speaking to an IMAP, POP3 or FTP server ('--ssl-reqd' on the command line or'CURLOPT_USE_SSL' set to 'CURLUSESSL_CONTROL' or 'CURLUSESSL_ALL' withlibcurl). This requirement could be bypassed if the server would return a properly crafted but perfectly legitimate response.This flaw would then make curl silently continue its operations **withoutTLS** contrary to the instructions and expectations, exposing possibly sensitive data in clear text over the network.",
          "first_seen_time": "2021-09-16T03:00:00Z",
          "fix_available": "1",
          "fixed_version": "7.58.0-2ubuntu3.15",
          "host_count": "1",
          "name": "curl",
          "namespace": "ubuntu:18.04",
          "package_status": "",
          "severity": "Medium",
          "version": "",
          "vulnerability_status": "Active"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": null,
          "Medium": {
            "fixable": 1,
            "vulnerabilities": 1
          }
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2015-3255",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2015-3255",
          "cvss_score": "4.6",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "The polkit_backend_action_pool_init function in polkitbackend/polkitbackendactionpool.c in PolicyKit (aka polkit) before 0.113 might allow local users to gain privileges via duplicate action IDs in action descriptions.",
          "first_seen_time": "2021-01-05T11:00:00Z",
          "fix_available": "1",
          "fixed_version": "0.105-11ubuntu1",
          "host_count": "1",
          "name": "policykit-1",
          "namespace": "ubuntu:18.04",
          "package_status": "",
          "severity": "Low",
          "version": "",
          "vulnerability_status": "Fixed"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": {
            "fixable": 1,
            "vulnerabilities": 1
          },
          "Medium": null
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2021-40528",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2021-40528",
          "cvss_score": "5.9",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "The ElGamal implementation in Libgcrypt before 1.9.4 allows plaintext recovery because, during interaction between two cryptographic libraries, a certain dangerous combination of the prime defined by the receiver's public key, the generator defined by the receiver's public key, and the sender's ephemeral exponents can lead to a cross-configuration attack against OpenPGP.",
          "first_seen_time": "2021-09-14T03:00:00Z",
          "fix_available": "1",
          "fixed_version": "1.8.1-4ubuntu1.3",
          "host_count": "1",
          "name": "libgcrypt20",
          "namespace": "ubuntu:18.04",
          "package_status": "",
          "severity": "Medium",
          "version": "",
          "vulnerability_status": "Active"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": null,
          "Medium": {
            "fixable": 1,
            "vulnerabilities": 1
          }
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2017-18207",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2017-18207",
          "cvss_score": "6.5",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "** DISPUTED ** The Wave_read._read_fmt_chunk function in Lib/wave.py in Python through 3.6.4 does not ensure a nonzero channel value, which allows attackers to cause a denial of service (divide-by-zero and exception) via a crafted wav format audio file. NOTE: the vendor disputes this issue because Python applications \\\\need to be prepared to handle a wide variety of exceptions.\\\\",
          "first_seen_time": "2021-02-28T03:00:00Z",
          "fix_available": "0",
          "fixed_version": "",
          "host_count": "1",
          "name": "python3.6",
          "namespace": "ubuntu:18.04",
          "package_status": "",
          "severity": "Low",
          "version": "",
          "vulnerability_status": "Active"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": {
            "fixable": 0,
            "vulnerabilities": 1
          },
          "Medium": null
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2021-31879",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2021-31879",
          "cvss_score": "6.1",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "GNU Wget through 1.21.1 does not omit the Authorization header upon a redirect to a different origin, a related issue to CVE-2018-1000007.",
          "first_seen_time": "2021-05-05T03:00:00Z",
          "fix_available": "0",
          "fixed_version": "",
          "host_count": "1",
          "name": "wget",
          "namespace": "ubuntu:18.04",
          "package_status": "ACTIVE",
          "severity": "Medium",
          "version": "",
          "vulnerability_status": "Active"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": null,
          "Medium": {
            "fixable": 0,
            "vulnerabilities": 1
          }
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2021-3517",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2021-3517",
          "cvss_score": "8.6",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "There is a flaw in the xml entity encoding functionality of libxml2 in versions before 2.9.11. An attacker who is able to supply a crafted file to be processed by an application linked with the affected functionality of libxml2 could trigger an out-of-bounds read. The most likely impact of this flaw is to application availability, with some potential impact to confidentiality and integrity if an attacker is able to use memory information to further exploit the application.",
          "first_seen_time": "2021-06-20T03:00:00Z",
          "fix_available": "0",
          "fixed_version": "",
          "host_count": "1",
          "name": "libxml2",
          "namespace": "ubuntu:18.04",
          "package_status": "",
          "severity": "Medium",
          "version": "",
          "vulnerability_status": "Fixed"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": null,
          "Medium": {
            "fixable": 0,
            "vulnerabilities": 1
          }
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2021-41072",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2021-41072",
          "cvss_score": "8.1",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "squashfs_opendir in unsquash-2.c in Squashfs-Tools 4.5 allows Directory Traversal, a different vulnerability than CVE-2021-40153. A squashfs filesystem that has been crafted to include a symbolic link and then contents under the same filename in a filesystem can cause unsquashfs to first create the symbolic link pointing outside the expected directory, and then the subsequent write operation will cause the unsquashfs process to write through the symbolic link elsewhere in the filesystem.",
          "first_seen_time": "2021-09-16T03:00:00Z",
          "fix_available": "1",
          "fixed_version": "1:4.3-6ubuntu0.18.04.4",
          "host_count": "1",
          "name": "squashfs-tools",
          "namespace": "ubuntu:18.04",
          "package_status": "",
          "severity": "Medium",
          "version": "",
          "vulnerability_status": "Active"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": null,
          "Medium": {
            "fixable": 1,
            "vulnerabilities": 1
          }
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2021-35267",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2021-35267",
          "cvss_score": "7.8",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "NTFS-3G versions < 2021.8.22, a stack buffer overflow can occur when correcting differences in the MFT and MFTMirror allowing for code execution or escalation of privileges when setuid-root.",
          "first_seen_time": "2021-09-09T03:00:00Z",
          "fix_available": "1",
          "fixed_version": "1:2017.3.23-2ubuntu0.18.04.3",
          "host_count": "1",
          "name": "ntfs-3g",
          "namespace": "ubuntu:18.04",
          "package_status": "",
          "severity": "Medium",
          "version": "",
          "vulnerability_status": "Active"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": null,
          "Medium": {
            "fixable": 1,
            "vulnerabilities": 1
          }
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2016-2568",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2016-2568",
          "cvss_score": "7.8",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "pkexec, when used with --user nonpriv, allows local users to escape to the parent session via a crafted TIOCSTI ioctl call, which pushes characters to the terminal's input buffer.",
          "first_seen_time": "2021-01-05T11:00:00Z",
          "fix_available": "0",
          "fixed_version": "",
          "host_count": "1",
          "name": "policykit-1",
          "namespace": "ubuntu:18.04",
          "package_status": "",
          "severity": "Low",
          "version": "",
          "vulnerability_status": "Active"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": {
            "fixable": 0,
            "vulnerabilities": 1
          },
          "Medium": null
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2021-39253",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2021-39253",
          "cvss_score": "7.8",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "A crafted NTFS image can cause an out-of-bounds read in ntfs_runlists_merge_i in NTFS-3G < 2021.8.22.",
          "first_seen_time": "2021-09-09T03:00:00Z",
          "fix_available": "1",
          "fixed_version": "1:2017.3.23-2ubuntu0.18.04.3",
          "host_count": "1",
          "name": "ntfs-3g",
          "namespace": "ubuntu:18.04",
          "package_status": "",
          "severity": "Medium",
          "version": "",
          "vulnerability_status": "Active"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": null,
          "Medium": {
            "fixable": 1,
            "vulnerabilities": 1
          }
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2016-2781",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2016-2781",
          "cvss_score": "6.5",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "chroot in GNU coreutils, when used with --userspec, allows local users to escape to the parent session via a crafted TIOCSTI ioctl call, which pushes characters to the terminal's input buffer.",
          "first_seen_time": "2021-01-05T11:00:00Z",
          "fix_available": "0",
          "fixed_version": "",
          "host_count": "1",
          "name": "coreutils",
          "namespace": "ubuntu:18.04",
          "package_status": "",
          "severity": "Low",
          "version": "",
          "vulnerability_status": "Active"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": {
            "fixable": 0,
            "vulnerabilities": 1
          },
          "Medium": null
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2012-2663",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2012-2663",
          "cvss_score": "7.5",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "extensions/libxt_tcp.c in iptables through 1.4.21 does not match TCP SYN+FIN packets in --syn rules, which might allow remote attackers to bypass intended firewall restrictions via crafted packets.  NOTE: the CVE-2012-6638 fix makes this issue less relevant.",
          "first_seen_time": "2021-01-05T11:00:00Z",
          "fix_available": "0",
          "fixed_version": "",
          "host_count": "1",
          "name": "iptables",
          "namespace": "ubuntu:18.04",
          "package_status": "",
          "severity": "Low",
          "version": "",
          "vulnerability_status": "Active"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": {
            "fixable": 0,
            "vulnerabilities": 1
          },
          "Medium": null
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2017-13716",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2017-13716",
          "cvss_score": "5.5",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "The C++ symbol demangler routine in cplus-dem.c in libiberty, as distributed in GNU Binutils 2.29, allows remote attackers to cause a denial of service (excessive memory allocation and application crash) via a crafted file, as demonstrated by a call from the Binary File Descriptor (BFD) library (aka libbfd).",
          "first_seen_time": "2021-03-04T03:00:00Z",
          "fix_available": "0",
          "fixed_version": "",
          "host_count": "1",
          "name": "binutils",
          "namespace": "ubuntu:18.04",
          "package_status": "",
          "severity": "Low",
          "version": "",
          "vulnerability_status": "Active"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": {
            "fixable": 0,
            "vulnerabilities": 1
          },
          "Medium": null
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2020-13988",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2020-13988",
          "cvss_score": "7.5",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "An issue was discovered in Contiki through 3.0. An Integer Overflow exists in the uIP TCP/IP Stack component when parsing TCP MSS options of IPv4 network packets in uip_process in net/ipv4/uip.c.",
          "first_seen_time": "2021-01-23T03:00:00Z",
          "fix_available": "0",
          "fixed_version": "",
          "host_count": "1",
          "name": "open-iscsi",
          "namespace": "ubuntu:18.04",
          "package_status": "",
          "severity": "Low",
          "version": "",
          "vulnerability_status": "Active"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": {
            "fixable": 0,
            "vulnerabilities": 1
          },
          "Medium": null
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2021-3711",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2021-3711",
          "cvss_score": "9.8",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "In order to decrypt SM2 encrypted data an application is expected to call the API function EVP_PKEY_decrypt(). Typically an application will call this function twice. The first time, on entry, the \\\\out\\\\ parameter can be NULL and, on exit, the \\\\outlen\\\\ parameter is populated with the buffer size required to hold the decrypted plaintext. The application can then allocate a sufficiently sized buffer and call EVP_PKEY_decrypt() again, but this time passing a non-NULL value for the \\\\out\\\\ parameter. A bug in the implementation of the SM2 decryption code means that the calculation of the buffer size required to hold the plaintext returned by the first call to EVP_PKEY_decrypt() can be smaller than the actual size required by the second call. This can lead to a buffer overflow when EVP_PKEY_decrypt() is called by the application a second time with a buffer that is too small. A malicious attacker who is able present SM2 content for decryption to an application could cause attacker chosen data to overflow the buffer by up to a maximum of 62 bytes altering the contents of other data held after the buffer, possibly changing application behaviour or causing the application to crash. The location of the buffer is application dependent but is typically heap allocated. Fixed in OpenSSL 1.1.1l (Affected 1.1.1-1.1.1k).",
          "first_seen_time": "2021-09-01T03:00:00Z",
          "fix_available": "1",
          "fixed_version": "1.1.1-1ubuntu2.1~18.04.13",
          "host_count": "1",
          "name": "openssl",
          "namespace": "ubuntu:18.04",
          "package_status": "",
          "severity": "High",
          "version": "",
          "vulnerability_status": "Active"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": {
            "fixable": 1,
            "vulnerabilities": 1
          },
          "Info": null,
          "Low": null,
          "Medium": null
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2018-6557",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2018-6557",
          "cvss_score": "7",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "The MOTD update script in the base-files package in Ubuntu 18.04 LTS before 10.1ubuntu2.2, and Ubuntu 18.10 before 10.1ubuntu6 incorrectly handled temporary files. A local attacker could use this issue to cause a denial of service, or possibly escalate privileges if kernel symlink restrictions were disabled.",
          "first_seen_time": "2021-01-05T11:00:00Z",
          "fix_available": "1",
          "fixed_version": "10.1ubuntu2.2",
          "host_count": "1",
          "name": "base-files",
          "namespace": "ubuntu:18.04",
          "package_status": "",
          "severity": "Low",
          "version": "",
          "vulnerability_status": "Fixed"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": {
            "fixable": 1,
            "vulnerabilities": 1
          },
          "Medium": null
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2016-1585",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2016-1585",
          "cvss_score": "9.8",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "In all versions of AppArmor mount rules are accidentally widened when compiled.",
          "first_seen_time": "2021-01-05T11:00:00Z",
          "fix_available": "0",
          "fixed_version": "",
          "host_count": "1",
          "name": "apparmor",
          "namespace": "ubuntu:18.04",
          "package_status": "",
          "severity": "Medium",
          "version": "",
          "vulnerability_status": "Active"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": null,
          "Medium": {
            "fixable": 0,
            "vulnerabilities": 1
          }
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2012-6655",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2012-6655",
          "cvss_score": "3.3",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "An issue exists AccountService 0.6.37 in the user_change_password_authorized_cb() function in user.c which could let a local users obtain encrypted passwords.",
          "first_seen_time": "2021-01-05T11:00:00Z",
          "fix_available": "0",
          "fixed_version": "",
          "host_count": "1",
          "name": "accountsservice",
          "namespace": "ubuntu:18.04",
          "package_status": "",
          "severity": "Low",
          "version": "",
          "vulnerability_status": "Active"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": {
            "fixable": 0,
            "vulnerabilities": 1
          },
          "Medium": null
        },
        "total_vulnerabilities": 1
      }
    },
    {
      "cve_id": "CVE-2019-9076",
      "packages": [
        {
          "cve_link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2019-9076",
          "cvss_score": "5.5",
          "cvss_v_2_score": "",
          "cvss_v_3_score": "",
          "description": "An issue was discovered in the Binary File Descriptor (BFD) library (aka libbfd), as distributed in GNU Binutils 2.32. It is an attempted excessive memory allocation in elf_read_notes in elf.c.",
          "first_seen_time": "2021-03-04T03:00:00Z",
          "fix_available": "0",
          "fixed_version": "",
          "host_count": "1",
          "name": "binutils",
          "namespace": "ubuntu:18.04",
          "package_status": "",
          "severity": "Low",
          "version": "",
          "vulnerability_status": "Active"
        }
      ],
      "summary": {
        "last_evaluation_time": "2021-11-23T03:00:00Z",
        "severity": {
          "Critical": null,
          "High": null,
          "Info": null,
          "Low": {
            "fixable": 0,
            "vulnerabilities": 1
          },
          "Medium": null
        },
        "total_vulnerabilities": 1
      }
    }
  ]
}`), &assessment)
	if err != nil {
		log.Fatal(err)
	}
	return assessment
}
