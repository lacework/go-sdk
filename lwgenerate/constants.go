package lwgenerate

const (
	LaceworkProviderSource  = "lacework/lacework"
	LaceworkProviderVersion = "~> 1.0"

	AwsConfigSource      = "lacework/config/aws"
	AwsConfigVersion     = "~> 0.5"
	AwsCloudTrailSource  = "lacework/cloudtrail/aws"
	AwsCloudTrailVersion = "~> 2.0"
	AwsEksAuditSource    = "lacework/eks-audit-log/aws"
	AwsEksAuditVersion   = "~> 0.4"

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

	GcpConfigSource          = "lacework/config/gcp"
	GcpConfigVersion         = "~> 2.3"
	GcpAuditLogSource        = "lacework/audit-log/gcp"
	GcpAuditLogVersion       = "~> 3.0"
	GcpGKEAuditLogSource     = "lacework/gke-audit-log/gcp"
	GcpGKEAuditLogVersion    = "~> 0.3"
	GcpPubSubAuditLog        = "lacework/pub-sub-audit-log/gcp"
	GcpPubSubAuditLogVersion = "~> 0.2"
)
