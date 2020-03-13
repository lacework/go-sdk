package api_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/lacework/go-sdk/api"
	"github.com/stretchr/testify/assert"
)

func TestGetIntegrations(t *testing.T) {
	// TODO @afiune implement a mocked Lacework API server
}

func TestCreateGCPConfigIntegration(t *testing.T) {
	intgGUID := "12345"

	fakeAPI := NewLaceworkServer()
	fakeAPI.MockAPI("external/integrations", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `
			{
				"data": [
					{
						"INTG_GUID": "`+intgGUID+`",
						"NAME": "integration_name",
						"CREATED_OR_UPDATED_TIME": "2020-Mar-10 01:00:00 UTC",
						"CREATED_OR_UPDATED_BY": "user@email.com",
						"TYPE": "GCP_CFG",
						"ENABLED": 1,
						"STATE": {
							"ok": true,
							"lastUpdatedTime": "2020-Mar-10 01:00:00 UTC",
							"lastSuccessfulTime": "2020-Mar-10 01:00:00 UTC"
						},
						"IS_ORG": 0,
						"DATA": {
							"CREDENTIALS": {
								"CLIENT_ID": "xxxxxxxxx",
								"CLIENT_EMAIL": "xxxxxx@xxxxx.iam.gserviceaccount.com",
								"PRIVATE_KEY_ID": "xxxxxxxxxxxxxxxx"
							},
							"ID_TYPE": "PROJECT",
							"ID": "xxxxxxxxxx"
						},
						"TYPE_NAME": "GCP Compliance"
					}
				],
				"ok": true,
				"message": "SUCCESS"
			}
		`)
	})
	defer fakeAPI.Close()

	c, err := api.NewClient("test", api.WithToken("xxxxxx"), api.WithURL(fakeAPI.URL()))
	assert.Nil(t, err)

	data := api.NewGCPIntegrationData("integration_name", api.GcpProject)
	assert.Equal(t, "GCP_CFG", data.Type, "a new GCP integration should match its type")
	data.Data.ID = "xxxxxxxxxx"
	data.Data.Credentials.ClientId = "xxxxxxxxx"
	data.Data.Credentials.ClientEmail = "xxxxxx@xxxxx.iam.gserviceaccount.com"
	data.Data.Credentials.PrivateKeyId = "xxxxxxxxxxxxxxxx"

	response, err := c.CreateGCPConfigIntegration(data)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.True(t, response.Ok)
	assert.Equal(t, 1, len(response.Data))
	assert.Equal(t, intgGUID, response.Data[0].IntgGuid)
}

func TestGetGCPConfigIntegration(t *testing.T) {
	intgGUID := "12345"
	apiPath := fmt.Sprintf("external/integrations/%s", intgGUID)

	fakeAPI := NewLaceworkServer()
	fakeAPI.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `
			{
				"data": [
					{
						"INTG_GUID": "`+intgGUID+`",
						"NAME": "integration_name",
						"CREATED_OR_UPDATED_TIME": "2020-Mar-10 01:00:00 UTC",
						"CREATED_OR_UPDATED_BY": "user@email.com",
						"TYPE": "GCP_CFG",
						"ENABLED": 1,
						"STATE": {
							"ok": true,
							"lastUpdatedTime": "2020-Mar-10 01:00:00 UTC",
							"lastSuccessfulTime": "2020-Mar-10 01:00:00 UTC"
						},
						"IS_ORG": 0,
						"DATA": {
							"CREDENTIALS": {
								"CLIENT_ID": "xxxxxxxxx",
								"CLIENT_EMAIL": "xxxxxx@xxxxx.iam.gserviceaccount.com",
								"PRIVATE_KEY_ID": "xxxxxxxxxxxxxxxx"
							},
							"ID_TYPE": "PROJECT",
							"ID": "xxxxxxxxxx"
						},
						"TYPE_NAME": "GCP Compliance"
					}
				],
				"ok": true,
				"message": "SUCCESS"
			}
		`)
	})
	defer fakeAPI.Close()

	c, err := api.NewClient("test", api.WithToken("xxxxxx"), api.WithURL(fakeAPI.URL()))
	assert.Nil(t, err)

	response, err := c.GetGCPConfigIntegration(intgGUID)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.True(t, response.Ok)
	assert.Equal(t, 1, len(response.Data))
	assert.Equal(t, intgGUID, response.Data[0].IntgGuid)
}
