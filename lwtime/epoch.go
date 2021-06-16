package lwtime

import (
	"strconv"
	"strings"
	"time"
)

type EpochTime struct {
	time.Time
}

func (epoch *EpochTime) UnmarshalJSON(b []byte) (err error) {
	t := strings.Trim(string(b), `"`)
	seconds, err := strconv.ParseInt(t, 10, 64)
	epoch.Time = time.Unix(seconds, 0)
	return
}

func (epoch *EpochTime) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(epoch.Time.Unix(), 10)), nil
}
