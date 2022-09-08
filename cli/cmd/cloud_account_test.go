package cmd

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/lacework/go-sdk/api"
	"github.com/stretchr/testify/assert"
)

func TestDetailsTable(t *testing.T) {
	detailsTableTest := []struct {
		Name     string
		Input    string
		Expected string
		RawType  api.V2RawType
	}{
		{
			Name:     "alertChannelRawResponse",
			Input:    mockAlertChannelRawResponse,
			Expected: mockAlertChannelTable,
			RawType:  api.AlertChannelRaw{},
		},
		{
			Name:     "cloudAccountRawResponse",
			Input:    mockCloudAccountRawResponse,
			Expected: mockCloudAccountTable,
			RawType:  api.CloudAccountRaw{},
		},
		{
			Name:     "cloudAccountRawResponse",
			Input:    mockContainerRegistryRawResponse,
			Expected: mockContainerRegistryTable,
			RawType:  api.CloudAccountRaw{},
		},
	}

	for _, test := range detailsTableTest {
		t.Run(test.Name, func(t *testing.T) {
			switch test.RawType {
			case api.AlertChannelRaw{}:
				test.RawType = new(api.AlertChannelRaw)
			case api.CloudAccountRaw{}:
				test.RawType = new(api.CloudAccountRaw)
			case api.ContainerRegistryRaw{}:
				test.RawType = new(api.ContainerRegistryRaw)
			}
			err := json.Unmarshal([]byte(test.Input), &test.RawType)
			if assert.NoError(t, err) {
				tableOut := buildDetailsTable(test.RawType)
				result := strings.Join(strings.Fields(tableOut), "")
				expected := strings.Join(strings.Fields(test.Expected), "")
				assert.Contains(t, result, expected)
			}
		})
	}
}

var mockAlertChannelRawResponse = `{
    "data": {
        "createdOrUpdatedBy": "darren.murray@lacework.net",
        "createdOrUpdatedTime": "2022-08-09T10:39:25.260Z",
        "enabled": 1,
        "intgGuid": "TECHALLY_E839836BC385C452E68B3CA7EB45BA0E7BDA39CCF65673A",
        "isOrg": 0,
        "name": "tech-ally-test",
        "state": {
            "ok": false,
            "lastUpdatedTime": 1662630616708,
            "lastSuccessfulTime": 0,
            "details": {
                "errorMessage": "403 403 Forbidden",
                "errorSubtitle": "Here is the response returned by AWS:",
                "errorTitle": "AWS Error",
                "message": "AWS Error: Error assuming role. INTG_GUID:TECHALLY_E839836BC385C452E68B3CA7EB45BA0E7BDA39CCF65673A [Error msg: 403 Forbidden]",
                "statusCode": "403"
            }
        },
        "type": "AwsS3",
        "data": {
            "s3CrossAccountCredentials": {
                "externalId": "12345678",
                "roleArn": "arn:aws:iam::1234567:role/lw-iam-abcdef",
                "bucketArn": "arn:aws:s3:::test"
            }
        }
    }
}`

var mockAlertChannelTable = `                                         DETAILS                                          
------------------------------------------------------------------------------------------
    CREATED OR UPDATED BY     darren.murray@lacework.net                                  
    CREATED OR UPDATED TIME   2022-08-09T10:39:25.260Z                                    
    ENABLED                   1                                                           
    INTG GUID                 TECHALLY_E839836BC385C452E68B3CA7EB45BA0E7BDA39CCF65673A    
    IS ORG                    0                                                           
    NAME                      tech-ally-test                                              
    TYPE                      AwsS3`

var mockCloudAccountRawResponse = `{
    "data": {
        "createdOrUpdatedBy": "salim.afiunemaya@lacework.net",
        "createdOrUpdatedTime": "2022-07-26T16:35:31.216Z",
        "enabled": 1,
        "intgGuid": "TECHALLY_B41875B7333E3858E1FA497DDC96BB51A7DED6F87F6A5D0",
        "isOrg": 0,
        "name": "TechAllyTest",
        "type": "GcpGkeAudit",
        "data": {
            "credentials": {
                "clientId": "12344566",
                "clientEmail": "tech-ally-test.iam.gserviceaccount.com",
                "privateKeyId": "",
                "privateKey": ""
            },
            "integrationType": "PROJECT",
            "projectId": "techally-test",
            "subscriptionName": "projects/techally--test"
        }
    }
}`

var mockCloudAccountTable = `                                         DETAILS                                          
------------------------------------------------------------------------------------------
    CREATED OR UPDATED BY     salim.afiunemaya@lacework.net                               
    CREATED OR UPDATED TIME   2022-07-26T16:35:31.216Z                                    
    ENABLED                   1                                                           
    INTEGRATION TYPE          PROJECT                                                     
    INTG GUID                 TECHALLY_B41875B7333E3858E1FA497DDC96BB51A7DED6F87F6A5D0    
    IS ORG                    0                                                           
    NAME                      TechAllyTest                                                
    PROJECT ID                techally-test                                               
    SUBSCRIPTION NAME         projects/techally--test                                     
    TYPE                      GcpGkeAudit`

var mockContainerRegistryRawResponse = `{
    "data": {
        "createdOrUpdatedBy": "salim.afiunemaya@lacework.net",
        "createdOrUpdatedTime": "2022-06-21T20:27:38.005Z",
        "enabled": 1,
        "intgGuid": "TECHALLY_EE34CCBAAEE52E1E6406DFA30231688DF5D44A5EE76BCD8",
        "isOrg": 0,
        "name": "Github Action With Policy",
        "props": {
            "policyEvaluation": {
                "policyGuids": [
                    "VULN_ABCD",
                    "VULN_EFGH"
                ],
                "evaluate": true
            },
            "tags": "INLINE_SCANNER"
        },
        "type": "ContVulnCfg",
        "data": {
            "registryType": "INLINE_SCANNER",
            "limitNumScan": "60",
            "identifierTag": []
        },
        "serverToken": {
            "serverToken": "abcd1234",
            "uri": "https://github.com/lacework/lacework-vulnerability-scanner/releases"
        }
    }
}`

var mockContainerRegistryTable = `                                               DETAILS                                               
-----------------------------------------------------------------------------------------------------
    CREATED OR UPDATED BY     salim.afiunemaya@lacework.net                                          
    CREATED OR UPDATED TIME   2022-06-21T20:27:38.005Z                                               
    ENABLED                   1                                                                      
    INTG GUID                 TECHALLY_EE34CCBAAEE52E1E6406DFA30231688DF5D44A5EE76BCD8               
    IS ORG                    0                                                                      
    LIMIT NUM SCAN            60                                                                     
    NAME                      Github Action With Policy                                              
    REGISTRY TYPE             INLINE_SCANNER                                                         
    SERVER TOKEN              abcd1234                                                               
    TAGS                      INLINE_SCANNER                                                         
    TYPE                      ContVulnCfg                                                            
    UPDATED AT                                                                                       
    UPDATED BY                                                                                       
    URI                       https://github.com/lacework/lacework-vulnerability-scanner/releases`
