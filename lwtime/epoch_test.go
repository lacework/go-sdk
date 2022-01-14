package lwtime

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockLastUpdatedResponse struct {
	LastUpdatedTime Epoch `json:"last_updated"`
}

type mockLastUpdatedIntResponse struct {
	LastUpdatedTime int `json:"last_updated"`
}

func TestMarshallEpoch(t *testing.T) {
	var res mockLastUpdatedResponse
	lastUpdated := mockLastUpdatedIntResponse{LastUpdatedTime: 1635604492078}
	lastUpdatedString, _ := json.Marshal(lastUpdated)
	json.Unmarshal(lastUpdatedString, &res)
	expected, _ := json.Marshal(lastUpdated)
	assert.Equal(t, "{\"last_updated\":1635604492078}", string(expected))
}
