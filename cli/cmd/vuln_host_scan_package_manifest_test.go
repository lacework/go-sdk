package cmd

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/lacework/go-sdk/api"
	"github.com/stretchr/testify/assert"
)

func TestFilterHostScanPackagesVulnDetails(t *testing.T) {
	res := filterHostScanPackagesVulnDetails(mockVulnPackages)

	assert.Equal(t, len(res), 1)
	assert.Equal(t, res[0].FixInfo.EvalStatus, "VULNERABLE")
}

func TestFilterHostScanPackagesVulnDetailsFixable(t *testing.T) {
	vulCmdState.Fixable = true
	defer func() {
		vulCmdState.Fixable = false
	}()
	res := filterHostScanPackagesVulnDetails(mockVulnPackages)

	assert.Equal(t, len(res), 1)
	assert.Equal(t, res[0].FixInfo.EvalStatus, "VULNERABLE")
}

func TestHostScanPackagesFailOnSeverity(t *testing.T) {
	vulCmdState.FailOnSeverity = "critical"
	defer func() {
		vulCmdState.FailOnSeverity = ""
	}()
	response, err := mockVulnSoftwarePackagesResponse()
	if err != nil {
		log.Fatal("unable to unmarshall VulnerabilitySoftwarePackagesResponse")
	}
	var expectedCount int32 = 1
	var expectedTotal int32 = 3

	err = buildVulnHostScanPkgManifestReports(&response)
	assessmentCounts := response.VulnerabilityCounts()
	vulnPolicy := NewVulnerabilityPolicyError(
		&assessmentCounts,
		vulCmdState.FailOnSeverity,
		vulCmdState.FailOnFixable,
	)
	nonCompliant := vulnPolicy.NonCompliant()

	assert.NoError(t, err)
	assert.Equal(t, assessmentCounts.Critical, expectedCount)
	assert.Equal(t, assessmentCounts.Total, expectedTotal)
	assert.True(t, nonCompliant)
}

func mockVulnSoftwarePackagesResponse() (api.VulnerabilitySoftwarePackagesResponse, error) {
	var mock api.VulnerabilitySoftwarePackagesResponse
	err := json.Unmarshal([]byte(mockVulnSoftwareResponse), &mock)
	return mock, err
}

var mockVulnPackages = []api.VulnerabilitySoftwarePackage{{FixInfo: fixInfo{
	EvalStatus:       "VULNERABLE",
	FixAvailable:     1,
	FixedVersion:     "test-fix-1",
	VersionInstalled: "current-version-1",
}}, {FixInfo: fixInfo{
	EvalStatus:       "GOOD",
	FixAvailable:     0,
	FixedVersion:     "test-fix-2",
	VersionInstalled: "current-version-2",
}}, {FixInfo: fixInfo{
	EvalStatus:       "GOOD",
	FixAvailable:     1,
	FixedVersion:     "test-fix-3",
	VersionInstalled: "current-version-3",
}}}

type fixInfo struct {
	CompareResult               int    `json:"compareResult"`
	EvalStatus                  string `json:"evalStatus"`
	FixAvailable                int    `json:"fixAvailable"`
	FixedVersion                string `json:"fixedVersion"`
	FixedVersionComparisonInfos []struct {
		CurrFixVer                         string `json:"currFixVer"`
		IsCurrFixVerGreaterThanOtherFixVer string `json:"isCurrFixVerGreaterThanOtherFixVer"`
		OtherFixVer                        string `json:"otherFixVer"`
	} `json:"fixedVersionComparisonInfos"`
	FixedVersionComparisonScore int    `json:"fixedVersionComparisonScore"`
	MaxPrefixMatchingLenScore   int    `json:"maxPrefixMatchingLenScore"`
	VersionInstalled            string `json:"versionInstalled"`
}

var mockVulnSoftwareResponse = `
{"data":[
{"osPkgInfo": {"namespace":"amzn:2","os":"amzn","osVer":"2","pkg":"python-babel","pkgVer":"0:0.9.6-8.amzn2.0.1","versionFormat":"rpm"},
	"vulnId":"ALAS2-2023-2010","severity":"Critical","featureKey":
		{"affectedRange": {"end":{"inclusive":false,"value":"0.9.6-8.amzn2.0.2"},"fixVersion":"0.9.6-8.amzn2.0.2",
			"start":{"inclusive":false,"value":"#MINV#"}},
			"name":"python-babel","namespace":"amzn:2"},
			"cveProps":{
				"cveBatchId":"E61EE2ABF4A948E6A4E236F243B016DE",
				"description":"Example Description",
				"link":"https://alas.aws.amazon.com/AL2/ALAS-2023-2010.html",
				"metadata":{"nvd":{"cvssv2":{"publisheddatetime":"","score":0,"vectors":""},
				"cvssv3":{"exploitabilityscore":0,"impactscore":0,"score":0,"vectors":""}}}},
				"fixInfo":{"compareResult":1,"evalStatus":"VULNERABLE","fixAvailable":1,"fixedVersion":"0:0.9.6-8.amzn2.0.2",
				"fixedVersionComparisonInfos":[{"currFixVer":"0.9.6-8.amzn2.0.2","isCurrFixVerGreaterThanOtherFixVer":"0","otherFixVer":"0.9.6-8.amzn2.0.2"}],
				"fixedVersionComparisonScore":0,"maxPrefixMatchingLenScore":18,"versionInstalled":"0:0.9.6-8.amzn2.0.1"},
				"summary":{"evalCreatedTime":"Thu, 20 Apr 2023 06:33:25 -0700","evalStatus":"MATCH_VULN","numFixableVuln":1,
				"numFixableVulnBySeverity":{"1":0,"2":0,"3":0,"4":1,"5":0},"numTotal":1,"numVuln":1,"numVulnBySeverity":{"1":0,"2":0,"3":0,"4":1,"5":0}},
				"props":{"evalAlgo":"1001"}},
{"osPkgInfo":{"namespace":"amzn:2","os":"amzn","osVer":"2","pkg":"dbus","pkgVer":"1:1.10.24-7.amzn2.0.2","versionFormat":"rpm"},
	"vulnId":"ALAS2-2023-2006","severity":"High","featureKey":
		{"affectedRange":{"end":{"inclusive":false,"value":"1:1.10.24-7.amzn2.0.3"},
				"fixVersion":"1:1.10.24-7.amzn2.0.3","start":{"inclusive":false,"value":"#MINV#"}},"name":"dbus","namespace":"amzn:2"},
				"cveProps":{"cveBatchId":"E61EE2ABF4A948E6A4E236F243B016DE","description":"Example Description",
				"link":"https://alas.aws.amazon.com/AL2/ALAS-2023-2006.html","metadata":{"nvd":{"cvssv2":{"publisheddatetime":"","score":0,"vectors":""},"cvssv3":{"exploitabilityscore":0,"impactscore":0,"score":0,"vectors":""}}}},
				"fixInfo":{"compareResult":1,"evalStatus":"VULNERABLE","fixAvailable":1,"fixedVersion":"1:1.10.24-7.amzn2.0.3",
				"fixedVersionComparisonInfos":[{"currFixVer":"1:1.10.24-7.amzn2.0.3","isCurrFixVerGreaterThanOtherFixVer":"0","otherFixVer":"1:1.10.24-7.amzn2.0.3"}],
				"fixedVersionComparisonScore":0,"maxPrefixMatchingLenScore":20,"versionInstalled":"1:1.10.24-7.amzn2.0.2"},
				"summary":{"evalCreatedTime":"Thu, 20 Apr 2023 06:33:25 -0700","evalStatus":"MATCH_VULN","numFixableVuln":1,
				"numFixableVulnBySeverity":{"1":0,"2":0,"3":0,"4":1,"5":0},"numTotal":2,"numVuln":1,
				"numVulnBySeverity":{"1":0,"2":0,"3":0,"4":1,"5":0}},
				"props":{"evalAlgo":"1001"}},
{"osPkgInfo":{"namespace":"amzn:2","os":"amzn","osVer":"2","pkg":"dbus","pkgVer":"1:1.10.24-7.amzn2.0.2","versionFormat":"rpm"},
	"vulnId":"ALAS2-2023-2006","severity":"Critical","featureKey":
		{"affectedRange":{"end":{"inclusive":false,"value":"1:1.10.24-7.amzn2.0.3"},
				"fixVersion":"1:1.10.24-7.amzn2.0.3","start":{"inclusive":false,"value":"#MINV#"}},"name":"dbus","namespace":"amzn:2"},
				"cveProps":{"cveBatchId":"E61EE2ABF4A948E6A4E236F243B016DE","description":"Example Description",
				"link":"https://alas.aws.amazon.com/AL2/ALAS-2023-2006.html","metadata":{"nvd":{"cvssv2":{"publisheddatetime":"","score":0,"vectors":""},"cvssv3":{"exploitabilityscore":0,"impactscore":0,"score":0,"vectors":""}}}},
				"fixInfo":{"compareResult":1,"evalStatus":"GOOD","fixAvailable":1,"fixedVersion":"1:1.10.24-7.amzn2.0.3",
				"fixedVersionComparisonInfos":[{"currFixVer":"1:1.10.24-7.amzn2.0.3","isCurrFixVerGreaterThanOtherFixVer":"0","otherFixVer":"1:1.10.24-7.amzn2.0.3"}],
				"fixedVersionComparisonScore":0,"maxPrefixMatchingLenScore":20,"versionInstalled":"1:1.10.24-7.amzn2.0.2"},
				"summary":{"evalCreatedTime":"Thu, 20 Apr 2023 06:33:25 -0700","evalStatus":"MATCH_VULN","numFixableVuln":1,
				"numFixableVulnBySeverity":{"1":0,"2":0,"3":0,"4":1,"5":0},"numTotal":2,"numVuln":1,
				"numVulnBySeverity":{"1":0,"2":0,"3":0,"4":1,"5":0}},
				"props":{"evalAlgo":"1001"}},
{"osPkgInfo":{"namespace":"amzn:2","os":"amzn","osVer":"2","pkg":"vim-data","pkgVer":"2:9.0.1367-1.amzn2.0.1","versionFormat":"rpm"},
	"vulnId":"ALAS2-2023-2005","severity":"Medium","featureKey":
		{"affectedRange":{"end":{"inclusive":false,"value":"2:9.0.1403-1.amzn2.0.1"},
				"fixVersion":"2:9.0.1403-1.amzn2.0.1","start":{"inclusive":false,"value":"#MINV#"}},
				"name":"vim-data","namespace":"amzn:2"},"cveProps":{"cveBatchId":"E61EE2ABF4A948E6A4E236F243B016DE",
				"description":"Example Description.","link":"https://alas.aws.amazon.com/AL2/ALAS-2023-2005.html",
				"metadata":{"nvd":{"cvssv2":{"publisheddatetime":"","score":0,"vectors":""},"cvssv3":{"exploitabilityscore":0,"impactscore":0,"score":0,"vectors":""}}}},
				"fixInfo":{"compareResult":1,"evalStatus":"VULNERABLE","fixAvailable":1,"fixedVersion":"2:9.0.1403-1.amzn2.0.1",
				"fixedVersionComparisonInfos":[{"currFixVer":"2:9.0.1403-1.amzn2.0.1","isCurrFixVerGreaterThanOtherFixVer":"0","otherFixVer":"2:9.0.1403-1.amzn2.0.1"}],
				"fixedVersionComparisonScore":0,"maxPrefixMatchingLenScore":7,"versionInstalled":"2:9.0.1367-1.amzn2.0.1"},
				"summary":{"evalCreatedTime":"Thu, 20 Apr 2023 06:33:25 -0700","evalStatus":"MATCH_VULN","numFixableVuln":1,
				"numFixableVulnBySeverity":{"1":0,"2":0,"3":0,"4":1,"5":0},"numTotal":12,"numVuln":1,
				"numVulnBySeverity":{"1":0,"2":0,"3":0,"4":1,"5":0}},
				"props":{"evalAlgo":"1001"}}]}
`
