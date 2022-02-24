package lwgenerate

const (
	LaceworkProviderSource  = "lacework/lacework"
	LaceworkProviderVersion = "~> 0.12.2"

	AwsConfigSource      = "lacework/config/aws"
	AwsConfigVersion     = "~> 0.1"
	AwsCloudTrailSource  = "lacework/cloudtrail/aws"
	AwsCloudTrailVersion = "~> 0.1"

	LWAzureConfigSource       = "lacework/config/azure"
	LWAzureConfigVersion      = "~> 1.0"
	LWAzureActivityLogSource  = "lacework/activity-log/azure"
	LWAzureActivityLogVersion = "~> 1.0"
	LWAzureADSource           = "lacework/ad-application/azure"
	LWAzureADVersion          = "~> 1.0"

	HashAzureADProviderSource  = "hashicorp/azuread"
	HashAzureADProviderVersion = "~> 2.16"
	HashAzureRMProviderSource  = "hashicorp/azurerm"
	HashAzureRMProviderVersion = "~> 2.91.0"

	GcpConfigSource    = "lacework/config/gcp"
	GcpConfigVersion   = "~> 1.0"
	GcpAuditLogSource  = "lacework/audit-log/gcp"
	GcpAuditLogVersion = "~> 2.0"
)
