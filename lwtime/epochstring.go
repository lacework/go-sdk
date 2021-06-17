package lwtime

import (
	"strconv"
	"strings"
	"time"
)

// EpochString time type to parse the returned 13 digit time in milliseconds
// Used instead of Epoch type when unmarshalling a json response where epoch time is a string
type EpochString time.Time

func (epoch *EpochString) UnmarshalJSON(b []byte) error {
	t := strings.Trim(string(b), `"`)
	seconds, _ := strconv.ParseInt(t, 10, 64)
	*epoch = EpochString(time.Unix(seconds, 0))
	return nil
}

func (epoch *EpochString) MarshalJSON() ([]byte, error) {
	return epoch.ToTime().UTC().MarshalJSON()
}

func (epoch EpochString) ToTime() time.Time {
	return time.Time(epoch)
}
func (epoch EpochString) Format(s string) string {
	return epoch.ToTime().Format(s)
}
func (epoch EpochString) UTC() time.Time {
	return epoch.ToTime().UTC()
}
