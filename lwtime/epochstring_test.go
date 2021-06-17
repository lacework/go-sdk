package lwtime

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockEvalTimeResponse struct {
	LastEvaluationTime EpochString `json:"last_evaluation_time"`
}

func TestUnmarshallEpoch(t *testing.T) {
	timeString := "1623812400000"
	jsonString := fmt.Sprintf(`{"last_evaluation_time": "%s"}`, timeString)
	res := mockEvalTimeResponse{}
	json.Unmarshal([]byte(jsonString), &res)
	assert.Equal(t, timeString, fmt.Sprint(res.LastEvaluationTime.ToTime().Unix()), "failed to parse Epoch time")
}
