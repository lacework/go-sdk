package lwtime

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type mockFirstSeenResponse struct {
	FirstTimeSeen RFC1123Z `json:"first_seen_time"`
}

func TestUnmarshallRFC1123Z(t *testing.T) {
	timeString := "Mon, 12 Apr 2021 20:00:00 -0700"
	jsonString := fmt.Sprintf(`{"first_seen_time": "%s"}`, timeString)
	res := mockFirstSeenResponse{}
	json.Unmarshal([]byte(jsonString), &res)
	assert.Equal(t, timeString, res.FirstTimeSeen.Format(time.RFC1123Z), "failed to parse RFC1123Z time")
}

func TestMarshallRFC1123Z(t *testing.T) {
	timeString, _ := time.Parse(time.RFC1123Z, "Mon, 12 Apr 2021 20:00:00 -0700")
	expectedJson := "{\"first_seen_time\":\"2021-04-13T03:00:00Z\"}"
	res := mockFirstSeenResponse{FirstTimeSeen: RFC1123Z(timeString)}

	marshalled, _ := json.Marshal(&res)
	assert.Equal(t, expectedJson, string(marshalled), "first_seen_time output is not RFC3339")
}
