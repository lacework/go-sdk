package cmd

import (
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
