package integrations

import "encoding/json"

type Response struct {
	Data    []CommonData `json:"data"`
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}

type CommonData struct {
	IntgGuid string `json:"INTG_GUID,omitempty"`
	Name string `json:"NAME"`
	CreatedOrUpdatedTime string `json:"CREATED_OR_UPDATED_TIME,omitempty"`
	CreatedOrUpdatedBy string `json:"CREATED_OR_UPDATED_BY,omitempty"`
	Type string `json:"TYPE"`
	Enabled int `json:"ENABLED"`
	State State `json:"STATE,omitempty"`
	IsOrg int `json:"IS_ORG,omitempty"`
	TypeName string `json:"TYPE_NAME,omitempty"`
	data map[string]interface{}
}

type State struct {
	Ok bool `json:"ok"`
	LastUpdatedTime string `json:"lastUpdatedTime"`
	LastSuccessfulTime string `json:"lastSuccessfulTime"`
}

type Data struct {
	CommonData
	Data map[string]interface{} `json:"DATA"`
}

func (commonData *Data) GetData() (interface{}, error) {
	var v interface{}
	switch commonData.Type {
	case "GCP_CFG", "GCP_AT_SES":
		v = &GcpInput{}
	case "AZURE_CFG", "AZURE_AL_SEQ":
		v = &AzureInput{}
	default:
		return commonData.Data, nil
	}
	jsonString, _ := json.Marshal(commonData.Data)
	err := json.Unmarshal(jsonString, v)
	if err != nil {
		return nil, err
	}
	return v, err
}

func CommonResponse() *Response {
	return &Response{}
}