package api_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/lacework"
	"github.com/stretchr/testify/assert"
)

func TestV2ReportsAwsGet(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.MockToken("TOKEN")
	fakeServer.MockAPI("Reports",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method, "Get() should be a GET method")
			fmt.Fprintf(w, mockAwsReport)
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.NoError(t, err)

	response, err := c.V2.Reports.Aws.Get(api.AwsReportConfig{AccountID: "123456789", Value: api.AWS_CIS_14.String(), Parameter: api.ReportFilterType})
	assert.NoError(t, err)
	assert.NotNil(t, response)
	if assert.Equal(t, 1, len(response.Data)) {
		assert.Equal(t, "123456789", response.Data[0].AccountID)
		assert.Equal(t, 44, response.Data[0].Summary[0].NumCompliant)
		assert.Equal(t, "AWS PCI DSS Report", response.Data[0].ReportTitle)
		assert.Equal(t, "test", response.Data[0].AccountAlias)
		assert.Equal(t, "example-region", response.Data[0].Recommendations[1].Violations[0].Region)
	}
}

func TestV2ReportsAwsGetByName(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.MockToken("TOKEN")
	fakeServer.MockAPI("Reports",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method, "Get() should be a GET method")
			fmt.Fprintf(w, mockAwsReport)
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.NoError(t, err)

	response, err := c.V2.Reports.Aws.Get(api.AwsReportConfig{AccountID: "123456789", Value: "AWS CIS Benchmark and S3", Parameter: api.ReportFilterName})
	assert.NoError(t, err)
	assert.NotNil(t, response)
	if assert.Equal(t, 1, len(response.Data)) {
		assert.Equal(t, "123456789", response.Data[0].AccountID)
		assert.Equal(t, 44, response.Data[0].Summary[0].NumCompliant)
		assert.Equal(t, "AWS PCI DSS Report", response.Data[0].ReportTitle)
		assert.Equal(t, "test", response.Data[0].AccountAlias)
		assert.Equal(t, "example-region", response.Data[0].Recommendations[1].Violations[0].Region)
	}
}

func TestV2ReportsAzureGet(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.MockToken("TOKEN")
	fakeServer.MockAPI("Reports",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method, "Get() should be a GET method")
			fmt.Fprintf(w, mockAzureReport)
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.NoError(t, err)

	response, err := c.V2.Reports.Azure.Get(api.AzureReportConfig{TenantID: "example-tenant", SubscriptionID: "example-subscription", Value: api.AZURE_CIS.String(), Parameter: api.ReportFilterType})
	assert.NoError(t, err)
	assert.NotNil(t, response)
	if assert.Equal(t, 1, len(response.Data)) {
		assert.Equal(t, "abcdefg-1234", response.Data[0].TenantID)
		assert.Equal(t, 13, response.Data[0].Summary[0].NumCompliant)
		assert.Equal(t, "Azure CIS Benchmark", response.Data[0].ReportTitle)
		assert.Equal(t, "123456-123456", response.Data[0].SubscriptionID)
		assert.Equal(t, "my-invalid-custom-role-12345", response.Data[0].Recommendations[1].Violations[0].Resource)
	}
}

func TestV2ReportsAzureGetByName(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.MockToken("TOKEN")
	fakeServer.MockAPI("Reports",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method, "Get() should be a GET method")
			fmt.Fprintf(w, mockAzureReport)
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.NoError(t, err)

	response, err := c.V2.Reports.Azure.Get(api.AzureReportConfig{TenantID: "example-tenant", SubscriptionID: "example-subscription", Value: "Azure CIS 1.3.1 Report", Parameter: api.ReportFilterName})
	assert.NoError(t, err)
	assert.NotNil(t, response)
	if assert.Equal(t, 1, len(response.Data)) {
		assert.Equal(t, "abcdefg-1234", response.Data[0].TenantID)
		assert.Equal(t, 13, response.Data[0].Summary[0].NumCompliant)
		assert.Equal(t, "Azure CIS Benchmark", response.Data[0].ReportTitle)
		assert.Equal(t, "123456-123456", response.Data[0].SubscriptionID)
		assert.Equal(t, "my-invalid-custom-role-12345", response.Data[0].Recommendations[1].Violations[0].Resource)
	}
}

func TestV2ReportsGcpGet(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.MockToken("TOKEN")
	fakeServer.MockAPI("Reports",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method, "Get() should be a GET method")
			fmt.Fprintf(w, mockGcpReport)
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.NoError(t, err)

	response, err := c.V2.Reports.Gcp.Get(api.GcpReportConfig{OrganizationID: "example-org", ProjectID: "example-project", Value: api.GCP_CIS.String(), Parameter: api.ReportFilterType})
	assert.NoError(t, err)
	assert.NotNil(t, response)
	if assert.Equal(t, 1, len(response.Data)) {
		assert.Equal(t, "123456789", response.Data[0].OrganizationID)
		assert.Equal(t, 36, response.Data[0].Summary[0].NumCompliant)
		assert.Equal(t, "GCP CIS Benchmark", response.Data[0].ReportTitle)
		assert.Equal(t, "test", response.Data[0].ProjectID)
		assert.Equal(t, "ServiceAccountHasAdminPrivileges", response.Data[0].Recommendations[1].Violations[0].Reasons[0])
	}
}

func TestV2ReportsGcpGetByName(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.MockToken("TOKEN")
	fakeServer.MockAPI("Reports",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method, "Get() should be a GET method")
			fmt.Fprintf(w, mockGcpReport)
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.NoError(t, err)

	response, err := c.V2.Reports.Gcp.Get(api.GcpReportConfig{OrganizationID: "example-org", ProjectID: "example-project", Value: "GCP CIS Benchmark 1.3", Parameter: api.ReportFilterName})
	assert.NoError(t, err)
	assert.NotNil(t, response)
	if assert.Equal(t, 1, len(response.Data)) {
		assert.Equal(t, "123456789", response.Data[0].OrganizationID)
		assert.Equal(t, 36, response.Data[0].Summary[0].NumCompliant)
		assert.Equal(t, "GCP CIS Benchmark", response.Data[0].ReportTitle)
		assert.Equal(t, "test", response.Data[0].ProjectID)
		assert.Equal(t, "ServiceAccountHasAdminPrivileges", response.Data[0].Recommendations[1].Violations[0].Reasons[0])
	}
}

var (
	mockAwsReport = `{
    "data": [
        {
            "reportType": "AWS PCI DSS Report",
            "reportTitle": "AWS PCI DSS Report",
            "recommendations": [
                {
                    "ACCOUNT_ID": "123456789",
                    "ACCOUNT_ALIAS": "test",
                    "START_TIME": 1665418788940,
                    "SUPPRESSIONS": [],
                    "INFO_LINK": "https://api.lacework.net/ui/documents/AWS_CIS_Foundations_Benchmark.pdf#page=114",
                    "ASSESSED_RESOURCE_COUNT": 0,
                    "STATUS": "RequiresManualAssessment",
                    "REC_ID": "AWS_CIS_3_10",
                    "CATEGORY": "1. Install and maintain a firewall configuration to protect cardholder data",
                    "SERVICE": "lw:cloudtrail",
                    "TITLE": "Ensure a log metric filter and alarm exist for security group changes",
                    "VIOLATIONS": [],
                    "RESOURCE_COUNT": 0,
                    "SEVERITY": 3
                },
                {
                    "ACCOUNT_ID": "123456789",
                    "ACCOUNT_ALIAS": "lacework-devtest",
                    "START_TIME": 1665418788940,
                    "SUPPRESSIONS": [],
                    "INFO_LINK": "https://api.lacework.net/ui/documents/Lacework_SecurityAudit_Descriptions.pdf#page=84",
                    "ASSESSED_RESOURCE_COUNT": 1,
                    "STATUS": "NonCompliant",
                    "REC_ID": "LW_AWS_ELASTICSEARCH_1",
                    "CATEGORY": "8. Identify and authenticate access to system components",
                    "SERVICE": "aws:es",
                    "TITLE": "This is an example reason",
                    "VIOLATIONS": [
                        {
                            "region": "example-region",
                            "reasons": [
                                "example-reason"
                            ],
                            "resource": "arn:aws:test:123456789:domain/test123"
                        }
                    ],
                    "RESOURCE_COUNT": 1,
                    "SEVERITY": 2
                }
            ],
            "summary": [
                {
                    "NUM_RECOMMENDATIONS": 86,
                    "NUM_SEVERITY_2_NON_COMPLIANCE": 14,
                    "NUM_SEVERITY_4_NON_COMPLIANCE": 2,
                    "NUM_SEVERITY_1_NON_COMPLIANCE": 9,
                    "NUM_COMPLIANT": 44,
                    "NUM_SEVERITY_3_NON_COMPLIANCE": 9,
                    "ASSESSED_RESOURCE_COUNT": 100,
                    "NUM_SUPPRESSED": 0,
                    "NUM_SEVERITY_5_NON_COMPLIANCE": 0,
                    "NUM_NOT_COMPLIANT": 34,
                    "VIOLATED_RESOURCE_COUNT": 100,
                    "SUPPRESSED_RESOURCE_COUNT": 0
                }
            ],
            "accountId": "123456789",
            "accountAlias": "test",
            "reportTime": "2022-10-10T16:19:48.940Z"
        }
    ],
    "ok": true,
    "message": "SUCCESS"
}`
	mockAzureReport = `{
    "data": [
        {
            "reportType": "Azure CIS Benchmark",
            "reportTitle": "Azure CIS Benchmark",
            "recommendations": [
                {
                    "TENANT_ID": "abcdefg-1234",
                    "TENANT_NAME": "abcdefg-1234",
                    "SUBSCRIPTION_ID": "123456-123456",
                    "SUBSCRIPTION_NAME": "Test",
                    "START_TIME": 1665565713798,
                    "INFO_LINK": "https://api.lacework.net/ui/documents/Azure_CIS_Foundations_Benchmark_v1.0.0.pdf#page=12",
                    "ASSESSED_RESOURCE_COUNT": 0,
                    "STATUS": "RequiresManualAssessment",
                    "REC_ID": "Azure_CIS_1_1",
                    "CATEGORY": "Identity and Access Management",
                    "SERVICE": "azure:ms:ad",
                    "TITLE": "Ensure that multi-factor authentication is enabled for all privileged users",
                    "RESOURCE_COUNT": 0,
                    "SEVERITY": 1
                },
                {
                    "TENANT_ID": "abcderg-12345",
                    "TENANT_NAME": "abcdefg-1234",
                    "SUBSCRIPTION_ID": "123456-123456",
                    "SUBSCRIPTION_NAME": "Test",
                    "START_TIME": 1665565713798,
                    "SUPPRESSIONS": [],
                    "INFO_LINK": "https://api.lacework.net/ui/documents/Azure_CIS_Foundations_Benchmark_v1.0.0.pdf#page=58",
                    "ASSESSED_RESOURCE_COUNT": 376,
                    "STATUS": "NonCompliant",
                    "REC_ID": "Azure_CIS_1_23",
                    "CATEGORY": "Identity and Access Management",
                    "SERVICE": "azure:ms:authority",
                    "TITLE": "Ensure that no custom subscription owner roles are created",
                    "VIOLATIONS": [
                        {
                            "reasons": [
                                "CustomRoleWithSubscriptionOwnership"
                            ],
                            "resource": "my-invalid-custom-role-12345"
                        }
                    ],
                    "RESOURCE_COUNT": 376,
                    "SEVERITY": 3
                }
            ],
            "summary": [
                {
                    "NUM_RECOMMENDATIONS": 60,
                    "NUM_SEVERITY_2_NON_COMPLIANCE": 37,
                    "NUM_SEVERITY_4_NON_COMPLIANCE": 0,
                    "NUM_SEVERITY_1_NON_COMPLIANCE": 0,
                    "NUM_COMPLIANT": 13,
                    "NUM_SEVERITY_3_NON_COMPLIANCE": 2,
                    "ASSESSED_RESOURCE_COUNT": 532,
                    "NUM_SUPPRESSED": 0,
                    "NUM_SEVERITY_5_NON_COMPLIANCE": 0,
                    "NUM_NOT_COMPLIANT": 39,
                    "VIOLATED_RESOURCE_COUNT": 72,
                    "SUPPRESSED_RESOURCE_COUNT": 0
                }
            ],
            "tenantId": "abcdefg-1234",
            "tenantName": "abcdefg-1234",
            "subscriptionId": "123456-123456",
            "subscriptionName": "Test",
            "reportTime": "2022-10-12T09:08:33.798Z"
        }
    ],
    "ok": true,
    "message": "SUCCESS"
}`

	mockGcpReport = `{
    "data": [
        {
            "reportType": "GCP CIS Benchmark",
            "reportTitle": "GCP CIS Benchmark",
            "recommendations": [
                {
                    "PROJECT_ID": "test",
                    "PROJECT_NAME": "test",
                    "ORGANIZATION_ID": "123456789",
                    "ORGANIZATION_NAME": "n/a",
                    "START_TIME": 1665655783846,
                    "SUPPRESSIONS": [],
                    "INFO_LINK": "https://api.lacework.net/ui/documents/GCP_CIS_Foundations_Benchmark.pdf#page=12",
                    "ASSESSED_RESOURCE_COUNT": 4,
                    "STATUS": "Compliant",
                    "REC_ID": "GCP_CIS_1_1",
                    "CATEGORY": "Identity and Access Management",
                    "SERVICE": "gcp:crm:projectIamPolicy",
                    "TITLE": "Ensure that corporate login credentials are used instead of Gmail accounts",
                    "VIOLATIONS": [],
                    "RESOURCE_COUNT": 4,
                    "SEVERITY": 3
                },
                {
                    "PROJECT_ID": "test",
                    "PROJECT_NAME": "test",
                    "ORGANIZATION_ID": "123456789",
                    "ORGANIZATION_NAME": "n/a",
                    "START_TIME": 1665655783846,
                    "SUPPRESSIONS": [],
                    "INFO_LINK": "https://api.lacework.net/ui/documents/GCP_CIS_Foundations_Benchmark.pdf#page=18",
                    "ASSESSED_RESOURCE_COUNT": 5,
                    "STATUS": "NonCompliant",
                    "REC_ID": "GCP_CIS_1_4",
                    "CATEGORY": "Identity and Access Management",
                    "SERVICE": "gcp:crm:projectIamPolicy",
                    "TITLE": "Ensure that ServiceAccount has no Admin privileges",
                    "VIOLATIONS": [
                        {
                            "reasons": [
                                "ServiceAccountHasAdminPrivileges"
                            ],
                            "resource": "serviceAccount:test@-test.iam.gserviceaccount.com"
                        }
                    ],
                    "RESOURCE_COUNT": 5,
                    "SEVERITY": 3
                }
            ],
            "summary": [
                {
                    "NUM_RECOMMENDATIONS": 62,
                    "NUM_SEVERITY_2_NON_COMPLIANCE": 0,
                    "NUM_SEVERITY_4_NON_COMPLIANCE": 0,
                    "NUM_SEVERITY_1_NON_COMPLIANCE": 3,
                    "NUM_COMPLIANT": 36,
                    "NUM_SEVERITY_3_NON_COMPLIANCE": 8,
                    "ASSESSED_RESOURCE_COUNT": 179,
                    "NUM_SUPPRESSED": 0,
                    "NUM_SEVERITY_5_NON_COMPLIANCE": 0,
                    "NUM_NOT_COMPLIANT": 11,
                    "VIOLATED_RESOURCE_COUNT": 71,
                    "SUPPRESSED_RESOURCE_COUNT": 0
                }
            ],
            "projectId": "test",
            "projectName": "test",
            "organizationId": "123456789",
            "organizationName": "n/a",
            "reportTime": "2022-10-13T10:09:43.846Z"
        }
    ],
    "ok": true,
    "message": "SUCCESS"
}
`
)
