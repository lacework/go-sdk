package integration

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/lacework/go-sdk/api"
	"github.com/stretchr/testify/assert"
)

var testReportDefinitionID = ""

func TestReportDefintionsList(t *testing.T) {
	out, err, exitcode := LaceworkCLIWithTOMLConfig("report-definitions", "list")
	// assert response contains table headers
	assert.Contains(t, out.String(), "GUID")
	assert.Contains(t, out.String(), "NAME")
	assert.Contains(t, out.String(), "TYPE")
	assert.Contains(t, out.String(), "SUB-TYPE")

	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}

func TestReportDefintionsListWithSubtype(t *testing.T) {
	out, err, exitcode := LaceworkCLIWithTOMLConfig("report-definitions", "list", "--subtype", "AWS")
	// assert response contains table headers
	assert.Contains(t, out.String(), "GUID")
	assert.Contains(t, out.String(), "NAME")
	assert.Contains(t, out.String(), "TYPE")
	assert.Contains(t, out.String(), "SUB-TYPE")

	assert.Contains(t, out.String(), "AWS")
	assert.NotContains(t, out.String(), "GCP")
	assert.NotContains(t, out.String(), "AZURE")

	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}

func TestReportDefintionsListJson(t *testing.T) {
	out, err, exitcode := LaceworkCLIWithTOMLConfig("report-definitions", "list", "--json")
	// assert response contains json fields
	assert.Contains(t, out.String(), "\"data\"")
	assert.Contains(t, out.String(), "\"createdBy\"")
	assert.Contains(t, out.String(), "\"displayName\"")
	assert.Contains(t, out.String(), "\"reportDefinition\"")
	assert.Contains(t, out.String(), "\"category\"")
	assert.Contains(t, out.String(), "\"policies\"")

	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}

func TestReportDefintionsListJsonWithSubtype(t *testing.T) {
	out, err, exitcode := LaceworkCLIWithTOMLConfig("report-definitions", "list", "--json", "--subtype", "GCP")
	// assert response contains json fields
	assert.Contains(t, out.String(), "\"data\"")
	assert.Contains(t, out.String(), "\"createdBy\"")
	assert.Contains(t, out.String(), "\"displayName\"")
	assert.Contains(t, out.String(), "\"reportDefinition\"")
	assert.Contains(t, out.String(), "\"category\"")
	assert.Contains(t, out.String(), "\"policies\"")

	assert.Contains(t, out.String(), "\"GCP\"")
	assert.NotContains(t, out.String(), "\"Azure\"")
	assert.NotContains(t, out.String(), "\"AWS\"")

	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}

func TestReportDefintionsShow(t *testing.T) {
	if testReportDefinitionID == "" {
		t.Skip("skipping test. No report definition found")
	}

	out, err, exitcode := LaceworkCLIWithTOMLConfig("report-definitions", "show", testReportDefinitionID, "--json")
	// assert response contains table headers
	assert.Contains(t, out.String(), "GUID")
	assert.Contains(t, out.String(), "NAME")
	assert.Contains(t, out.String(), "TYPE")
	assert.Contains(t, out.String(), "SUB-TYPE")
	// assert response contains table data
	assert.Contains(t, out.String(), reportDefinitionDetailsOutput)
	assert.Contains(t, out.String(), "lacework-global-34")

	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}

func TestReportDefintionsShowJson(t *testing.T) {
	if testReportDefinitionID == "" {
		t.Skip("skipping test. No report definition found")
	}

	out, err, exitcode := LaceworkCLIWithTOMLConfig("report-definitions", "show", testReportDefinitionID, "--json")
	// assert response contains json fields
	assert.Contains(t, out.String(), "\"data\"")
	assert.Contains(t, out.String(), "\"createdBy\"")
	assert.Contains(t, out.String(), "\"displayName\"")
	assert.Contains(t, out.String(), "\"reportDefinition\"")
	assert.Contains(t, out.String(), "\"category\"")
	assert.Contains(t, out.String(), "\"policies\"")
	// assert response contains json data
	assert.Contains(t, out.String(), reportDefinitionDetailsJson)
	assert.Contains(t, out.String(), "lacework-global-34")
	assert.Contains(t, out.String(), fmt.Sprintf("\"%s\"", testReportDefinitionID))

	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}

func TestReportDefinitionsDiff(t *testing.T) {
	versions := fetchVersionedCustomReportDefinition()

	if len(versions.Data) == 0 {
		t.Skip("skipping test. No versions for custom report definition found")
	}
	id := versions.Data[0].ReportDefinitionGuid

	currentVersion := "2"
	lastVersion := "1"

	out, err, exitcode := LaceworkCLIWithTOMLConfig("report-definitions", "diff", id, currentVersion, lastVersion)

	assert.Contains(t, out.String(), reportDefinitionDiff)
	fmt.Println(out.String())
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")

	deleteErr := deleteCustomReportDefinition(id)
	assert.NoError(t, deleteErr)
}

func TestReportDefinitionsRevert(t *testing.T) {
	versions := fetchVersionedCustomReportDefinition()

	if len(versions.Data) == 0 {
		t.Skip("skipping test. No versions for custom report definition found")
	}
	id := versions.Data[0].ReportDefinitionGuid
	lastVersion := "1"

	out, err, exitcode := LaceworkCLIWithTOMLConfig("report-definitions", "revert", id, lastVersion)

	assert.Contains(t, out.String(), fmt.Sprintf("The report definition %s was reverted to version %s", id, lastVersion))
	fmt.Println(out.String())
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")

	deleteErr := deleteCustomReportDefinition(id)
	assert.NoError(t, deleteErr)
}

var reportDefinitionDetailsOutput = `              REPORT DEFINITION DETAILS              
-----------------------------------------------------
    ENGINE          lpp                              
    RELEASE LABEL                                    
    UPDATED BY      SYSTEM                           
    LAST UPDATED    2022-09-09 10:35:16 +0000 UTC    
                                                     
                               POLICIES                                
-----------------------------------------------------------------------
                TITLE                            POLICY                
  ---------------------------------+---------------------------------  
`

var reportDefinitionDetailsJson = `  "data": {
    "createdBy": "SYSTEM",
    "createdTime": "2022-09-09T10:35:16Z",
    "displayName": "AWS NIST 800-171 rev2 Report",
    "distributionType": "pdf",
    "enabled": 1,
    "frequency": "daily",
    "props": {
      "engine": "lpp"
    },`

var reportDefinitionDiff = `-reportName: Diff Test Versioning Updated
+reportName: Diff Test Versioning
 displayName: Diff Test Versioning
 reportType: COMPLIANCE
 subReportType: AWS
@@ -9,7 +9,7 @@
           title: Diff Test Report
           policies:
             - lacework-global-1
-version: 2
+version: 1`

func fetchVersionedCustomReportDefinition() api.ReportDefinitionsResponse {
	lacework, err := api.NewClient(os.Getenv("CI_ACCOUNT"),
		api.WithSubaccount(os.Getenv("CI_SUBACCOUNT")),
		api.WithApiKeys(os.Getenv("CI_API_KEY"), os.Getenv("CI_API_SECRET")),
		api.WithApiV2(),
	)
	if err != nil {
		log.Fatal(err)
	}

	rptCfg := api.ReportDefinitionConfig{
		ReportName:    "Diff Test Versioning",
		DisplayName:   "Diff Test Versioning",
		ReportType:    api.ReportDefinitionTypeCompliance.String(),
		SubReportType: api.ReportDefinitionSubTypeAws.String(),
		Sections: []api.ReportDefinitionSection{{
			Title:    "Diff Test Report",
			Policies: []string{"lacework-global-1"},
		}},
	}

	customReport := api.NewReportDefinition(rptCfg)

	report, err := lacework.V2.ReportDefinitions.Create(customReport)
	if err != nil {
		log.Fatal(err)
	}

	rptCfg.ReportName = "Diff Test Versioning Updated"

	report, err = lacework.V2.ReportDefinitions.Update(report.Data.ReportDefinitionGuid, api.NewReportDefinitionUpdate(rptCfg))
	if err != nil {
		log.Fatal(err)
	}

	reports, err := lacework.V2.ReportDefinitions.GetVersions(report.Data.ReportDefinitionGuid)
	if err != nil {
		log.Fatal(err)
	}

	return reports
}

func deleteCustomReportDefinition(id string) error {
	lacework, err := api.NewClient(os.Getenv("CI_ACCOUNT"),
		api.WithSubaccount(os.Getenv("CI_SUBACCOUNT")),
		api.WithApiKeys(os.Getenv("CI_API_KEY"), os.Getenv("CI_API_SECRET")),
		api.WithApiV2(),
	)
	if err != nil {
		log.Fatal(err)
	}

	return lacework.V2.ReportDefinitions.Delete(id)
}
