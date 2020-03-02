package integrations

type GcpResponse struct {
	Data    []GcpData `json:"data"`
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}

type GcpData struct {
	CommonData
	Data GcpInput `json:"DATA"`
}

type GcpInput struct {
	Credentials GcpCredentials `json:"CREDENTIALS"`
	IssueGrouping string `json:"ISSUE_GROUPING,omitempty"`
	IdType string `json:"ID_TYPE"`
	Id string `json:"ID"`
	SubscriptionName string `json:"SUBSCRIPTION_NAME,omitempty"`
}

type GcpCredentials struct {
	ClientId string `json:"CLIENT_ID"`
	ClientEmail string `json:"CLIENT_EMAIL"`
	PrivateKeyId string `json:"PRIVATE_KEY_ID"`
	PrivateKey string `json:"PRIVATE_KEY"`
}

func Gcp() *GcpResponse {
	return &GcpResponse{}
}
