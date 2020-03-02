package integrations

type AzureCfgIntegrationResponse struct {
	Data    []AzureCfgIntegrationData `json:"data"`
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}

type AzureCfgIntegrationData struct {
	CommonIntegrationData
	Data AzureCfg `json:"DATA"`
}

type AzureCfg struct {
	Credentials AzureCredentials `json:"CREDENTIALS"`
	IssueGrouping string `json:"ISSUE_GROUPING"`
	TenantId string `json:"TENANT_ID"`
}

type AzureAlIntegrationResponse struct {
	Data    []GcpAtIntegrationData `json:"data"`
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}

type AzureAlIntegrationData struct {
	CommonIntegrationData
	Data AzureAl `json:"DATA"`
}

type AzureAl struct {
	AzureCfg
	QueueUrl string `json:"QUEUE_URL"`
}

type AzureCredentials struct {
	ClientId string `json:"CLIENT_ID"`
	ClientSecret string `json:"CLIENT_SECRET"`
}

func GetAzureCfgInterface() *AzureCfgIntegrationResponse {
	return &AzureCfgIntegrationResponse{}
}

func GetAzureAlInterface() *AzureAlIntegrationResponse {
	return &AzureAlIntegrationResponse{}
}
