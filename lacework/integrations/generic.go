package integrations

type GenericIntegrationResponse struct {
	Data    []GenericIntegrationData `json:"data"`
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}

type CommonIntegrationData struct {
	IntgGuid string `json:"INTG_GUID"`
	Name string `json:"NAME"`
	CreatedOrUpdatedTime string `json:"CREATED_OR_UPDATED_TIME"`
	CreatedOrUpdatedBy string `json:"CREATED_OR_UPDATED_BY"`
	Type string `json:"TYPE"`
	Enabled int `json:"ENABLED"`
	State State `json:"STATE"`
	IsOrg int `json:"IS_ORG"`
	TypeName string `json:"TYPE_NAME"`
}

type State struct {
	Ok bool `json:"ok"`
	LastUpdatedTime string `json:"lastUpdatedTime"`
	LastSuccessfulTime string `json:"lastSuccessfulTime"`
}

type GenericIntegrationData struct {
	CommonIntegrationData
	Data map[string]interface{} `json:"DATA"`
}

func GetGenericInterface() *GenericIntegrationResponse {
	return &GenericIntegrationResponse{}
}