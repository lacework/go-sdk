package integrations

type GcpCfgIntegrationResponse struct {
	Data    []GcpCfgIntegrationData `json:"data"`
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}

type GcpCfgIntegrationData struct {
	CommonIntegrationData
	Data GcpCfg `json:"DATA"`
}

type GcpCfg struct {
	Credentials GcpCredentials `json:"CREDENTIALS"`
	IssueGrouping string `json:"ISSUE_GROUPING"`
	IdType string `json:"ID_TYPE"`
	Id string `json:"ID"`
}

type GcpAtIntegrationResponse struct {
	Data    []GcpAtIntegrationData `json:"data"`
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}

type GcpAtIntegrationData struct {
	CommonIntegrationData
	Data GcpAt `json:"DATA"`
}

type GcpAt struct {
	GcpCfg
	SubscriptionName string `json:"SUBSCRIPTION_NAME"`
}

type GcpCredentials struct {
	ClientId string `json:"CLIENT_ID"`
	ClientEmail string `json:"CLIENT_EMAIL"`
	PrivateKeyId string `json:"PRIVATE_KEY_ID"`
	PrivateKey string `json:"PRIVATE_KEY"`
}

func GetGcpCfgInterface() *GcpCfgIntegrationResponse {
	return &GcpCfgIntegrationResponse{}
}

func GetGcpAtInterface() *GcpAtIntegrationResponse {
	return &GcpAtIntegrationResponse{}
}
