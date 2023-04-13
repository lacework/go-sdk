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

package api_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/lacework"
)

func TestV2Vulnerabilities_Containers_SearchAllPages_EmptyData(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.UseApiV2()
	fakeServer.MockToken("TOKEN")
	fakeServer.MockAPI("Vulnerabilities/Containers/search",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method, "Search() should be a POST method")
			fmt.Fprintf(w, mockPaginationEmptyResponse())
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithApiV2(),
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.NoError(t, err)

	response, err := c.V2.Vulnerabilities.Containers.SearchAllPages(api.SearchFilter{})
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 0, len(response.Data))
}

func TestV2Vulnerabilities_Hosts_SearchAllPages_EmptyData(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.UseApiV2()
	fakeServer.MockToken("TOKEN")
	fakeServer.MockAPI("Vulnerabilities/Hosts/search",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method, "Search() should be a POST method")
			fmt.Fprintf(w, mockPaginationEmptyResponse())
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithApiV2(),
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.NoError(t, err)

	response, err := c.V2.Vulnerabilities.Hosts.SearchAllPages(api.SearchFilter{})
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 0, len(response.Data))
}

func TestV2VulnerabilitiesFilterSingleVulnID(t *testing.T) {
	var (
		response api.VulnerabilitiesContainersResponse
		vulnID   = "CVE-2017-12670"
	)

	json.Unmarshal([]byte(mockVulnerabilitiesContainersResponse()), &response)
	response.FilterSingleVulnIDData(vulnID)

	assert.Equal(t, len(response.Data), 1)
	assert.Equal(t, response.Data[0].VulnID, vulnID)
}
func TestV2Vulnerabilities_Containers_Search(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.UseApiV2()
	fakeServer.MockToken("TOKEN")
	fakeServer.MockAPI("Vulnerabilities/Containers/search",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method, "Search() should be a POST method")
			fmt.Fprintf(w, mockVulnerabilitiesContainersResponse())
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithApiV2(),
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.NoError(t, err)

	response, err := c.V2.Vulnerabilities.Containers.Search(api.SearchFilter{})
	assert.NoError(t, err)
	assert.NotNil(t, response)
	if assert.Equal(t, 2, len(response.Data)) {
		assert.Equal(t, "sha256:123472164cb78c4d04f57bd66201c775e2dab08fce394806a03a933c5daf9e48",
			response.Data[0].ImageID)
	}
}

func TestV2Vulnerabilities_Containers_SearchLastWeek(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.UseApiV2()
	fakeServer.MockToken("TOKEN")
	fakeServer.MockAPI("Vulnerabilities/Containers/search",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method, "Search() should be a POST method")
			fmt.Fprintf(w, mockVulnerabilitiesContainersResponse())
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithApiV2(),
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.NoError(t, err)

	response, err := c.V2.Vulnerabilities.Containers.SearchLastWeek()
	assert.NoError(t, err)
	assert.NotNil(t, response)
	if assert.Equal(t, 2, len(response.Data)) {
		assert.Equal(t, "EXCEPTION", response.Data[0].Status)
		assert.Equal(t, "sha256:123472164cb78c4d04f57bd66201c775e2dab08fce394806a03a933c5daf9e48",
			response.Data[0].ImageID)
		assert.Empty(t, response.Paging.Urls.NextPage)
	}
}

func TestV2Vulnerabilities_Containers_AllPages(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.UseApiV2()
	fakeServer.MockToken("TOKEN")
	fakeServer.MockAPI("Vulnerabilities/Containers/search",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method, "Search() should be a POST method")
			fmt.Fprintf(w, mockVulnerabilitiesContainersResponse())
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithApiV2(),
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.NoError(t, err)

	response, err := c.V2.Vulnerabilities.Containers.SearchAllPages(api.SearchFilter{})
	assert.NoError(t, err)
	assert.NotNil(t, response)
	if assert.Equal(t, 2, len(response.Data)) {
		assert.Equal(t, "Low", response.Data[0].Severity)
		assert.Empty(t, response.Paging, "paging should be empty")
	}
}

func mockPaginationEmptyResponse() string {
	return `
{
  "data": [],
  "paging": {
    "rows": 0,
    "totalRows": 0,
    "urls": {
      "nextPage": null
    }
  }
}`
}
func mockVulnerabilitiesContainersResponse() string {
	return `
{
  "data": [
	    {
      "evalCtx": {
        "cve_batch_info": [
          {
            "cve_batch_id": "12341234FA904FC1986CD8E5387ED053",
            "cve_created_time": "2022-02-10 00:06:46.292000000"
          }
        ],
        "exception_props": [
          {
            "exception_guid": "VULN_1234F2CBE09F0E705565BEA1A0C1D2A5D1534857F2C7CDF8381",
            "exception_name": "registry index.docker.io severity Low",
            "exception_reason": "Accepted Risk"
          }
        ],
        "image_info": {
          "created_time": 1605140985874,
          "digest": "sha256:1234d2246518044ef95e3dbd029e51dd477788e5bf8e278e418685aabc3fe28a",
          "error_msg": [],
          "id": "sha256:123472164cb78c4d04f57bd66201c775e2dab08fce394806a03a933c5daf9e48",
          "registry": "index.docker.io",
          "repo": "techallylw/test-cli-dirty",
          "scan_created_time": "2020-11-12T00:29:45.874+00:00",
          "size": 360608563,
          "status": "Success",
          "tags": [
            "latest"
          ],
          "type": "Docker"
        },
        "integration_props": {
          "INTG_GUID": "TECHALLY_12341E30911A91A63494EC993FD4791F506D6561C55D98F",
          "NAME": "TF tech-ally docker",
          "REGISTRY_TYPE": "DOCKERHUB"
        },
        "is_reeval": false,
        "request_source": "PLATFORM_SCANNER",
        "scan_batch_id": "1234a960-992e-44b6-9d0b-42c91ec66fbc-1644487510402549454",
        "scan_request_props": {
          "data_format_version": "1.0",
          "props": {
            "data_format_version": "1.0",
            "scanner_version": "0.2.8"
          },
          "reqId": "1234367f-707b-44a5-abf9-a20cbf8d8369",
          "scanCompletionUtcTime": 1644487510,
          "scan_start_time": 1644487502,
          "scanner_version": "0.2.8"
        },
        "vuln_batch_id": "1234D974FA904FC1986CD8E5387ED053",
        "vuln_created_time": "2022-02-10 00:06:46.292000000"
      },
      "featureKey": {
        "name": "imagemagick",
        "namespace": "debian:9",
        "version": "8:6.9.7.4+dfsg-11+deb9u10"
      },
      "fixInfo": {
        "compare_result": 1,
        "fix_available": 0,
        "fixed_version": ""
      },
      "imageId": "sha256:123472164cb78c4d04f57bd66201c775e2dab08fce394806a03a933c5daf9e48",
      "severity": "Low",
      "startTime": "2022-02-10T10:05:11.418Z",
      "status": "EXCEPTION",
      "vulnId": "CVE-2017-12670"
    },
    {
      "evalCtx": {
        "cve_batch_info": [
          {
            "cve_batch_id": "1234D974FA904FC1986CD8E5387ED053",
            "cve_created_time": "2022-02-10 00:06:46.292000000"
          }
        ],
        "exception_props": [
          {
            "exception_guid": "VULN_1234F2CBE09F0E705565BEA1A0C1D2A5D1534857F2C7CDF8381",
            "exception_name": "registry index.docker.io severity Low",
            "exception_reason": "Accepted Risk"
          }
        ],
        "image_info": {
          "created_time": 1605140985874,
          "digest": "sha256:1234d2246518044ef95e3dbd029e51dd477788e5bf8e278e418685aabc3fe28a",
          "error_msg": [],
          "id": "sha256:123472164cb78c4d04f57bd66201c775e2dab08fce394806a03a933c5daf9e48",
          "registry": "index.docker.io",
          "repo": "techallylw/test-cli-dirty",
          "scan_created_time": "2020-11-12T00:29:45.874+00:00",
          "size": 360608563,
          "status": "Success",
          "tags": [
            "latest"
          ],
          "type": "Docker"
        },
        "integration_props": {
          "INTG_GUID": "TECHALLY_12341E30911A91A63494EC993FD4791F506D6561C55D98F",
          "NAME": "TF tech-ally docker",
          "REGISTRY_TYPE": "DOCKERHUB"
        },
        "is_reeval": false,
        "request_source": "PLATFORM_SCANNER",
        "scan_batch_id": "1234a960-992e-44b6-9d0b-42c91ec66fbc-1644487510402549454",
        "scan_request_props": {
          "data_format_version": "1.0",
          "props": {
            "data_format_version": "1.0",
            "scanner_version": "0.2.8"
          },
          "reqId": "1234367f-707b-44a5-abf9-a20cbf8d8369",
          "scanCompletionUtcTime": 1644487510,
          "scan_start_time": 1644487502,
          "scanner_version": "0.2.8"
        },
        "vuln_batch_id": "1234D974FA904FC1986CD8E5387ED053",
        "vuln_created_time": "2022-02-10 00:06:46.292000000"
      },
      "featureKey": {
        "name": "imagemagick",
        "namespace": "debian:9",
        "version": "8:6.9.7.4+dfsg-11+deb9u10"
      },
      "fixInfo": {
        "compare_result": 1,
        "fix_available": 0,
        "fixed_version": ""
      },
      "imageId": "sha256:123472164cb78c4d04f57bd66201c775e2dab08fce394806a03a933c5daf9e48",
      "severity": "Info",
      "startTime": "2022-02-10T10:05:11.418Z",
      "status": "VULNERABLE",
      "vulnId": "CVE-2017-14342"
    }
  ],
  "paging": {
    "rows": 2,
    "totalRows": 2,
    "urls": {
      "nextPage": null
    }
  }
}
	`
}

func TestV2Vulnerabilities_Hosts_Search(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.UseApiV2()
	fakeServer.MockToken("TOKEN")
	fakeServer.MockAPI("Vulnerabilities/Hosts/search",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method, "Search() should be a POST method")
			fmt.Fprintf(w, mockVulnerabilitiesHostsResponse())
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithApiV2(),
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.NoError(t, err)

	response, err := c.V2.Vulnerabilities.Hosts.Search(api.SearchFilter{})
	assert.NoError(t, err)
	assert.NotNil(t, response)
	if assert.Equal(t, 2, len(response.Data)) {
		assert.Equal(t, "Medium", response.Data[0].Severity)
	}
}

func TestV2Vulnerabilities_Hosts_SearchLastWeek(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.UseApiV2()
	fakeServer.MockToken("TOKEN")
	fakeServer.MockAPI("Vulnerabilities/Hosts/search",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method, "Search() should be a POST method")
			fmt.Fprintf(w, mockVulnerabilitiesHostsResponse())
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithApiV2(),
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.NoError(t, err)

	response, err := c.V2.Vulnerabilities.Hosts.SearchLastWeek()
	assert.NoError(t, err)
	assert.NotNil(t, response)
	if assert.Equal(t, 2, len(response.Data)) {
		assert.Equal(t, "FixedOnDiscovery", response.Data[0].Status)
		assert.Equal(t, "ALAS2-2019-1340", response.Data[0].VulnID)
		assert.Empty(t, response.Paging.Urls.NextPage)
	}
}

func TestV2Vulnerabilities_Hosts_AllPages(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.UseApiV2()
	fakeServer.MockToken("TOKEN")
	fakeServer.MockAPI("Vulnerabilities/Hosts/search",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method, "Search() should be a POST method")
			fmt.Fprintf(w, mockVulnerabilitiesHostsResponse())
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithApiV2(),
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.NoError(t, err)

	response, err := c.V2.Vulnerabilities.Hosts.SearchAllPages(api.SearchFilter{})
	assert.NoError(t, err)
	assert.NotNil(t, response)
	if assert.Equal(t, 2, len(response.Data)) {
		assert.Equal(t, "Medium", response.Data[0].Severity)
		assert.Empty(t, response.Paging, "paging should be empty")
	}
}

func TestV2Vulnerabilities_HostGetAwsMachineTags(t *testing.T) {
	var mockHostResponse api.VulnerabilitiesHostResponse
	err := json.Unmarshal([]byte(mockVulnerabilitiesHostsResponseSetTags(vulnerabilityHostAwsMachineTags)), &mockHostResponse)
	assert.NoError(t, err)

	tags, err := mockHostResponse.Data[0].GetMachineTags()
	assert.NoError(t, err)
	assert.Equal(t, tags.Account, "123456789038")
	assert.Equal(t, tags.AmiID, "ami-1234567890540c038")
	assert.Equal(t, tags.ExternalIP, "1.5.8.9")
	assert.Equal(t, tags.Hostname, "ip-192-168-28-69.us-east-2.compute.internal")
}

func TestV2Vulnerabilities_HostGetGcpMachineTags(t *testing.T) {
	var mockHostResponse api.VulnerabilitiesHostResponse
	err := json.Unmarshal([]byte(mockVulnerabilitiesHostsResponseSetTags(vulnerabilityHostGcpMachineTags)), &mockHostResponse)
	assert.NoError(t, err)

	tags, err := mockHostResponse.Data[0].GetMachineTags()
	assert.NoError(t, err)
	assert.Equal(t, tags.ProjectId, "tech-ally-test")
	assert.Equal(t, tags.NumericProjectId, "123456789012")
	assert.Equal(t, tags.GCEtags, "test")
	assert.Equal(t, tags.InstanceName, "tech-ally-test")
	assert.Equal(t, tags.VMInstanceType, "e2-small")
	assert.Equal(t, tags.VMProvider, "GCE")
}

func TestV2Vulnerabilities_HostGetEmptyMachineTags(t *testing.T) {
	var mockHostResponse api.VulnerabilitiesHostResponse
	err := json.Unmarshal([]byte(mockVulnerabilitiesHostsResponseSetTags("{}")), &mockHostResponse)
	assert.NoError(t, err)

	tags, err := mockHostResponse.Data[0].GetMachineTags()
	assert.NoError(t, err)
	assert.NotNil(t, tags)
	assert.Equal(t, tags.Account, "")
	assert.Equal(t, tags.AmiID, "")
	assert.Equal(t, tags.ExternalIP, "")
	assert.Equal(t, tags.Hostname, "")
	assert.Equal(t, tags.ProjectId, "")
	assert.Equal(t, tags.NumericProjectId, "")
}

func mockVulnerabilitiesHostsResponse() string {
	return `
{
  "data": [
	  {
      "cveProps": {
        "cve_batch_id": "1234567890904FC1986CD8E5387ED053",
        "description": "Package updates are available for Amazon Linux 2 that fix the following vulnerabilities: CVE-2019-5482: Heap buffer overflow in the TFTP protocol handler in cURL 7.19.4 to 7.65.3. 99999: CVE-2019-5482 curl: heap buffer overflow in function tftp_receive_packet() CVE-2019-5481: Double-free vulnerability in the FTP-kerberos code in cURL 7.52.0 to 7.65.3. 99999: CVE-2019-5481 curl: double free due to subsequent call of realloc()",
        "link": "https://alas.aws.amazon.com/AL2/ALAS-2019-1340.html"
      },
      "endTime": "2022-02-10T04:00:00.000Z",
      "evalCtx": {
        "exception_props": [],
        "hostname": "ip-192-168-28-69.us-east-2.compute.internal",
        "mc_eval_guid": "1234567890736f2a56c224175a9c06c4"
      },
      "featureKey": {
        "name": "curl",
        "namespace": "amzn:2",
        "package_active": 0,
        "version_installed": "0:7.61.1-12.amzn2.0.2"
      },
      "fixInfo": {
        "compare_result": "-1",
        "eval_status": "GOOD",
        "fix_available": "0",
        "fixed_version": "7.61.1-12.amzn2.0.1",
        "fixed_version_comparison_infos": [
          {
            "curr_fix_ver": "7.61.1-12.amzn2.0.1",
            "is_curr_fix_ver_greater_than_other_fix_ver": "0",
            "other_fix_ver": "7.61.1-12.amzn2.0.1"
          }
        ],
        "fixed_version_comparison_score": 0,
        "version_installed": "0:7.61.1-12.amzn2.0.2"
      },
      "machineTags": {
        "Account": "123456789038",
        "AmiId": "ami-1234567890540c038",
        "Env": "k8s",
        "ExternalIp": "1.5.8.9",
        "Hostname": "ip-192-168-28-69.us-east-2.compute.internal",
        "InstanceId": "i-12345678903bd1a6c",
        "InternalIp": "192.168.28.69",
        "LwTokenShort": "12345678904c316d1ed18fbd15f168",
        "Name": "techally-sandbox-standard-workers-Node",
        "SubnetId": "subnet-1234567890bba6219",
        "VmInstanceType": "t3.small",
        "VmProvider": "AWS",
        "VpcId": "vpc-1234567890252f137",
        "Zone": "us-east-2c",
        "alpha.eksctl.io/nodegroup-name": "standard-workers",
        "alpha.eksctl.io/nodegroup-type": "managed",
        "arch": "amd64",
        "aws:autoscaling:groupName": "eks-123456789019e-1234567890156a8d5dc982",
        "aws:ec2:fleet-id": "fleet-1234567890e80-44f6-123456789068ef762",
        "aws:ec2launchtemplate:id": "lt-1234567890c17103c",
        "aws:ec2launchtemplate:version": "1",
        "eks:cluster-name": "techally-sandbox",
        "eks:nodegroup-name": "standard-workers",
        "k8s.io/cluster-autoscaler/enabled": 1,
        "k8s.io/cluster-autoscaler/techally-sandbox": "owned",
        "kubernetes.io/cluster/techally-sandbox": "owned",
        "lw_KubernetesCluster": "techally-sandbox",
        "os": "linux"
      },
      "mid": 38,
      "severity": "Medium",
      "startTime": "2022-02-10T03:00:00.000Z",
      "status": "FixedOnDiscovery",
      "vulnId": "ALAS2-2019-1340"
    },
    {
      "cveProps": {
        "cve_batch_id": "1234567890904FC1986CD8E5387ED053",
        "description": "Package updates are available for Amazon Linux 2 that fix the following vulnerabilities: CVE-2021-40490: 2001951: CVE-2021-40490 kernel: race condition was discovered in ext4_write_inline_data_end in fs/ext4/inline.c in the ext4 subsystem A flaw was found in the Linux kernel. A race condition was discovered in the ext4 subsystem. The highest threat from this vulnerability is to data confidentiality and integrity as well as system availability. CVE-2021-38198: A flaw was found in the Linux kernel, where it incorrectly computes the access permissions of a shadow page. This issue leads to a missing guest protection page fault. 1992264: CVE-2021-38198 kernel: arch/x86/kvm/mmu/paging_tmpl.h incorrectly computes the access permissions of a shadow page CVE-2021-3753: 99999:Linux Kernel could allow a local attacker to obtain sensitive information, caused by an out-of-bounds read flaw in VT. By using a specially-crafted vc_visible_origin setting, an attacker could exploit this vulnerability to obtain sensitive information, or cause a denial of service condition. CVE-2021-3732: 1995249: CVE-2021-3732 kernel: overlayfs: Mounting overlayfs inside an unprivileged user namespace can reveal files A flaw was found in the Linux kernel's OverlayFS subsystem in the way the user mounts the TmpFS filesystem with OverlayFS. This flaw allows a local user to gain access to hidden files that should not be accessible. CVE-2021-3656: 1983988: CVE-2021-3656 kernel: SVM nested virtualization issue in KVM (VMLOAD/VMSAVE) A flaw was found in the KVM's AMD code for supporting SVM nested virtualization. The flaw occurs when processing the VMCB (virtual machine control block) provided by the L1 guest to spawn/handle a nested guest (L2). Due to improper validation of the \\\\virt_ext\\\\ field, this issue could allow a malicious L1 to disable both VMLOAD/VMSAVE intercepts and VLS (Virtual VMLOAD/VMSAVE) for the L2 guest. As a result, the L2 guest would be allowed to read/write physical pages of the host, resulting in a crash of the entire system, leak of sensitive data or potential guest-to-host escape. CVE-2021-3653: A flaw was found in the KVM's AMD code for supporting SVM nested virtualization. The flaw occurs when processing the VMCB (virtual machine control block) provided by the L1 guest to spawn/handle a nested guest (L2). Due to improper validation of the \\\\int_ctl\\\\ field, this issue could allow a malicious L1 to enable AVIC support (Advanced Virtual Interrupt Controller) for the L2 guest. As a result, the L2 guest would be allowed to read/write physical pages of the host, resulting in a crash of the entire system, leak of sensitive data or potential guest-to-host escape. 1983686: CVE-2021-3653 kernel: SVM nested virtualization issue in KVM (AVIC support)",
        "link": "https://alas.aws.amazon.com/AL2/ALAS-2021-1704.html"
      },
      "endTime": "2022-02-10T04:00:00.000Z",
      "evalCtx": {
        "exception_props": [],
        "hostname": "ip-10-0-1-46.us-west-2.compute.internal",
        "mc_eval_guid": "1234567890736f2a56c224175a9c06c4"
      },
      "featureKey": {
        "name": "kernel-tools",
        "namespace": "amzn:2",
        "package_active": 0,
        "version_installed": "0:4.14.209-160.339.amzn2"
      },
      "fixInfo": {
        "compare_result": "1",
        "eval_status": "VULNERABLE",
        "fix_available": "1",
        "fixed_version": "4.14.246-187.474.amzn2",
        "fixed_version_comparison_infos": [
          {
            "curr_fix_ver": "4.14.246-187.474.amzn2",
            "is_curr_fix_ver_greater_than_other_fix_ver": "0",
            "other_fix_ver": "4.14.246-187.474.amzn2"
          }
        ],
        "fixed_version_comparison_score": 0,
        "version_installed": "0:4.14.209-160.339.amzn2"
      },
      "machineTags": {
        "Account": "123456789038",
        "AmiId": "ami-1234567890395e172",
        "ExternalIp": "5.2.1.1",
        "Hostname": "ip-10-0-1-46.us-west-2.compute.internal",
        "InstanceId": "i-1234567890c3aa2ce",
        "InternalIp": "10.0.1.46",
        "LwTokenShort": "1234567890775231a4bd605be8bc9b",
        "SubnetId": "subnet-1234567890a99009c",
        "VmInstanceType": "t2.micro",
        "VmProvider": "AWS",
        "VpcId": "vpc-123456789057b1a1d",
        "Zone": "us-west-2a",
        "arch": "amd64",
        "os": "linux"
      },
      "mid": 52,
      "severity": "Medium",
      "startTime": "2022-02-10T03:00:00.000Z",
      "status": "Reopened",
      "vulnId": "ALAS2-2021-1704"
    }
  ],
  "paging": {
    "rows": 2,
    "totalRows": 2,
    "urls": {
      "nextPage": null
    }
  }
}
	`
}

func mockVulnerabilitiesHostsResponseSetTags(machineTagsJson string) string {
	return fmt.Sprintf(`
{
  "data": [
	  {
      "cveProps": {
        "cve_batch_id": "1234567890904FC1986CD8E5387ED053",
        "description": "Package updates are available for Amazon Linux 2 that fix the following vulnerabilities: CVE-2019-5482: Heap buffer overflow in the TFTP protocol handler in cURL 7.19.4 to 7.65.3. 99999: CVE-2019-5482 curl: heap buffer overflow in function tftp_receive_packet() CVE-2019-5481: Double-free vulnerability in the FTP-kerberos code in cURL 7.52.0 to 7.65.3. 99999: CVE-2019-5481 curl: double free due to subsequent call of realloc()",
        "link": "https://alas.aws.amazon.com/AL2/ALAS-2019-1340.html"
      },
      "endTime": "2022-02-10T04:00:00.000Z",
      "evalCtx": {
        "exception_props": [],
        "hostname": "ip-192-168-28-69.us-east-2.compute.internal",
        "mc_eval_guid": "1234567890736f2a56c224175a9c06c4"
      },
      "featureKey": {
        "name": "curl",
        "namespace": "amzn:2",
        "package_active": 0,
        "version_installed": "0:7.61.1-12.amzn2.0.2"
      },
      "fixInfo": {
        "compare_result": "-1",
        "eval_status": "GOOD",
        "fix_available": "0",
        "fixed_version": "7.61.1-12.amzn2.0.1",
        "fixed_version_comparison_infos": [
          {
            "curr_fix_ver": "7.61.1-12.amzn2.0.1",
            "is_curr_fix_ver_greater_than_other_fix_ver": "0",
            "other_fix_ver": "7.61.1-12.amzn2.0.1"
          }
        ],
        "fixed_version_comparison_score": 0,
        "version_installed": "0:7.61.1-12.amzn2.0.2"
      },
      "machineTags": %s,
      "mid": 38,
      "severity": "Medium",
      "startTime": "2022-02-10T03:00:00.000Z",
      "status": "FixedOnDiscovery",
      "vulnId": "ALAS2-2019-1340"
    }
  ],
  "paging": {
    "rows": 2,
    "totalRows": 2,
    "urls": {
      "nextPage": null
    }
  }
}
	`, machineTagsJson)
}

var vulnerabilityHostAwsMachineTags = `{
        "Account": "123456789038",
        "AmiId": "ami-1234567890540c038",
        "Env": "k8s",
        "ExternalIp": "1.5.8.9",
        "Hostname": "ip-192-168-28-69.us-east-2.compute.internal",
        "InstanceId": "i-12345678903bd1a6c",
        "InternalIp": "192.168.28.69",
        "LwTokenShort": "12345678904c316d1ed18fbd15f168",
        "Name": "techally-sandbox-standard-workers-Node",
        "SubnetId": "subnet-1234567890bba6219",
        "VmInstanceType": "t3.small",
        "VmProvider": "AWS",
        "VpcId": "vpc-1234567890252f137",
        "Zone": "us-east-2c",
        "alpha.eksctl.io/nodegroup-name": "standard-workers",
        "alpha.eksctl.io/nodegroup-type": "managed",
        "arch": "amd64",
        "aws:autoscaling:groupName": "eks-123456789019e-1234567890156a8d5dc982",
        "aws:ec2:fleet-id": "fleet-1234567890e80-44f6-123456789068ef762",
        "aws:ec2launchtemplate:id": "lt-1234567890c17103c",
        "aws:ec2launchtemplate:version": "1",
        "eks:cluster-name": "techally-sandbox",
        "eks:nodegroup-name": "standard-workers",
        "k8s.io/cluster-autoscaler/enabled": 1,
        "k8s.io/cluster-autoscaler/techally-sandbox": "owned",
        "kubernetes.io/cluster/techally-sandbox": "owned",
        "lw_KubernetesCluster": "techally-sandbox",
        "os": "linux"
      }`

var vulnerabilityHostGcpMachineTags = `{
                "ExternalIp": "1.123.123.123",
                "GCEtags": "test",
                "Hostname": "test.c.test.internal",
                "InstanceId": "1234567892345678901",
                "InstanceName": "tech-ally-test",
                "InternalIp": "10.123.1.2",
                "LwTokenShort": "123456abcde12a1234a112345abcd1",
                "NumericProjectId": "123456789012",
                "ProjectId": "tech-ally-test",
                "VmInstanceType": "e2-small",
                "VmProvider": "GCE",
                "Zone": "us-central1-a",
                "arch": "amd64",
                "lw_InternetExposure": "Unknown",
                "os": "linux"
      }`
