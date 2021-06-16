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
	assert.Equal(t, timeString, res.FirstTimeSeen.Time.Format(time.RFC1123Z), "failed to parse RFC1123Z time")
}
