package lwtime

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type mockEvalTimeResponse struct {
	LastEvaluationTime EpochString `json:"last_evaluation_time"`
}

func TestUnmarshallEpochString(t *testing.T) {
	timeString := "1623812400000"
	jsonString := fmt.Sprintf(`{"last_evaluation_time": "%s"}`, timeString)
	res := mockEvalTimeResponse{}
	json.Unmarshal([]byte(jsonString), &res)
	assert.Contains(t, timeString, fmt.Sprint(res.LastEvaluationTime.ToTime().Unix()*1000), "failed to parse Epoch time")
}

func TestMarshallEpochString(t *testing.T) {
	timeEpoch := time.Unix(1623898800000/1000, 0)
	expectedJson := "{\"last_evaluation_time\":\"2021-06-17T03:00:00Z\"}"
	res := mockEvalTimeResponse{LastEvaluationTime: EpochString(timeEpoch)}

	marshalled, _ := json.Marshal(&res)
	assert.Equal(t, expectedJson, string(marshalled), "first_seen_time output is not RFC3339")
}
