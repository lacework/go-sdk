package client

const defaultTokenExpiryTime = 3600

// authConfig representing information like key_id, secret and token
// used for authenticating requests
type authConfig struct {
	keyID      string
	secret     string
	token      string
	expiration int
}

// WithApiKeys sets the key_id and secret used to generate API access tokens
func WithApiKeys(id, secret string) Option {
	return clientFunc(func(c *client) {
		if c.auth == nil {
			c.auth = &authConfig{}
		}
		c.auth.keyID = id
		c.auth.secret = secret
	})
}

// WithToken sets the token used to authenticate the API requests
func WithToken(token string) Option {
	return clientFunc(func(c *client) {
		c.auth.token = token
	})
}

// WithExpirationTime configures the token expiration time
func WithExpirationTime(t int) Option {
	return clientFunc(func(c *client) {
		c.auth.expiration = t
	})
}

// GenerateToken generates a new access token
func (c *client) GenerateToken(keyID, secretKey string) (response tokenResponse, err error) {
	c.auth.keyID = keyID
	c.auth.secret = secretKey
	body, err := jsonReader(tokenRequest{keyID, c.auth.expiration})
	if err != nil {
		return
	}

	err = c.requestDecoder("POST", apiTokens, body, &response)
	if err != nil {
		return
	}

	c.auth.token = response.Data[0].Token

	return
}

type tokenResponse struct {
	Data    []tokenData `json:"data"`
	Ok      bool        `json:"ok"`
	Message string      `json:"message"`
}

type tokenData struct {
	ExpiresAt string `json:"expiresAt"`
	Token     string `json:"token"`
}

type tokenRequest struct {
	KeyId      string `json:"keyId"`
	ExpiryTime int    `json:"expiryTime"`
}
