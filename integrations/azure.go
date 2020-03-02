package integrations

type AzureResponse struct {
	Data    []AzureData `json:"data"`
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}

type AzureData struct {
	CommonData
	Data AzureInput `json:"DATA"`
}

type AzureInput struct {
	Credentials AzureCredentials `json:"CREDENTIALS"`
	IssueGrouping string `json:"ISSUE_GROUPING"`
	TenantId string `json:"TENANT_ID"`
	QueueUrl string `json:"QUEUE_URL"`
}

type AzureCredentials struct {
	ClientId string `json:"CLIENT_ID"`
	ClientSecret string `json:"CLIENT_SECRET"`
}

func Azure() *AzureResponse {
	return &AzureResponse{}
}
