package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const (
	path = "/api/v1/external/integrations"
	tokenPath = "/api/v1/access/tokens"
	defaultTimeout = 60 * time.Second
	expiryTime = 3600
)

type GenerateTokenResponse struct {
	Data    []TokenData `json:"data"`
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}

type TokenData struct {
	ExpiresAt string `json:"expiresAt"`
	Token string `json:"token"`
}

type TokenBody struct {
	KeyId      string `json:"keyId"`
	ExpiryTime int    `json:"expiryTime"`
}

type Client struct {
	BaseURL *url.URL
	account, AuthToken string
	httpClient *http.Client
	//GcpService, AwsService, AzureService Service
}

func NewClient(account string, keyId, secretKey, authToken string) (*Client, error) {
	baseUrl, err := url.Parse("https://" + account + ".lacework.net")
	if err != nil {
		return nil, err
	}
	httpClient := &http.Client{Timeout: defaultTimeout}

	client := &Client{
		httpClient: httpClient,
		account: account,
		BaseURL: baseUrl,
	}

	if authToken == "" {
		if keyId != "" && secretKey != "" {
			v := &GenerateTokenResponse{}

			_, response, err := client.GenerateToken(keyId, secretKey, v)
			if err != nil {
				fmt.Printf("%s\n", err)
				return nil, err
			} else {
				fmt.Println(response.Status)
				fmt.Printf("Auth token generated: %s\n", v.Data[0].Token)
				client.AuthToken = v.Data[0].Token
			}
		}
	} else {
		client.AuthToken = authToken
	}

	return client, nil
}

func (client *Client) GenerateToken(keyId, secretKey string, v interface{}) (string, *http.Response, error) {
	body := &TokenBody{
		KeyId:      keyId,
		ExpiryTime: expiryTime,
	}

	headers := make(map[string]string)
	headers["X-LW-UAKS"] = secretKey

	request, err := client.newRequest("POST", tokenPath, headers, body)
	if err != nil {
		return "", nil, err
	}

	data, response, err := client.do(request, v)
	if err != nil {
		return "", nil, err
	}
	return data, response, err
}

func (client *Client) GetIntegrations(requestPath string, v interface{}) (string, *http.Response, error) {
	headers := make(map[string]string)
	headers["Accept"] = "application/json"
	headers["Authorization"] = client.AuthToken

	if requestPath == "" {
		requestPath = path
	}

	request, err := client.newRequest("GET", requestPath, headers, nil)
	if err != nil {
		return "", nil, err
	}

	data, response, err := client.do(request, v)
	if err != nil {
		return "", nil, err
	}
	return data, response, err
}

func (client *Client) GetIntegrationOfType(integrationType string, v interface{}) (string, *http.Response, error) {
	return client.GetIntegrations(path + "/type/" + integrationType, v)
}

func (client *Client) CreateIntegration(body interface{}, v interface{}) (string, *http.Response, error) {
	headers := make(map[string]string)
	headers["Accept"] = "application/json"
	headers["Authorization"] = client.AuthToken

	request, err := client.newRequest("POST", path, headers, body)
	if err != nil {
		return "", nil, err
	}

	data, response, err := client.do(request, v)
	if err != nil {
		return "", nil, err
	}
	return data, response, err
}

func (client *Client) newRequest(method, path string, headers map[string]string, body interface{}) (*http.Request, error) {
	relativeUrl := &url.URL{Path: path}
	resolvedUrl := client.BaseURL.ResolveReference(relativeUrl)
	var buffer io.ReadWriter
	if body != nil {
		buffer = new(bytes.Buffer)
		err := json.NewEncoder(buffer).Encode(body)
		if err != nil {
			return nil, err
		}
	}
	request, err := http.NewRequest(method, resolvedUrl.String(), buffer)
	if err != nil {
		return nil, err
	}
	if body != nil {
		request.Header.Set("Content-Type", "application/json")
	}

	for k, v := range headers {
		request.Header.Set(k, v)
	}

	return request, nil
}

func (client *Client) do(req *http.Request, v interface{}) (string, *http.Response, error) {
	response, err := client.httpClient.Do(req)
	if err != nil {
		return "", nil, err
	}

	if err := checkErrorInResponse(response); err != nil {
		return "", response, err
	}

	defer func() {
		if responseErr := response.Body.Close(); err == nil {
			err = responseErr
		}
	}()

	if v != nil {
		if err := decodeJSON(response, v); err != nil {
			return "", response, err
		}
	}

	data, _ := ioutil.ReadAll(response.Body)
	return string(data), response, err
}

func decodeJSON(res *http.Response, v interface{}) error {
	return json.NewDecoder(res.Body).Decode(v)
}
