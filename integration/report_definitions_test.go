package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReportDefintionsList(t *testing.T) {
	out, err, exitcode := LaceworkCLI("report-definitions", "list")
	// assert response contains table headers
	assert.Contains(t, out.String(), "GUID")
	assert.Contains(t, out.String(), "NAME")
	assert.Contains(t, out.String(), "TYPE")
	assert.Contains(t, out.String(), "SUB-TYPE")

	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}

func TestReportDefintionsListJson(t *testing.T) {
	out, err, exitcode := LaceworkCLI("report-definitions", "list", "--json")
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

func TestReportDefintionsShow(t *testing.T) {
	out, err, exitcode := LaceworkCLI("report-definitions", "show", "FC86105C0B6C8549F55104CBEE2FAAB6EF056C313DE87351BDBC3C79", "--json")
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
	out, err, exitcode := LaceworkCLI("report-definitions", "show", "FC86105C0B6C8549F55104CBEE2FAAB6EF056C313DE87351BDBC3C79", "--json")
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
	assert.Contains(t, out.String(), "\"FC86105C0B6C8549F55104CBEE2FAAB6EF056C313DE87351BDBC3C79\"")

	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
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
