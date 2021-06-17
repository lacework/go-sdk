package lwtime

import (
	"strconv"
	"time"
)

// Epoch time type to parse the returned 13 digit time in milliseconds
type Epoch time.Time

// implement Marshal and Unmarshal interfaces
func (epoch *Epoch) UnmarshalJSON(b []byte) error {
	ms, _ := strconv.Atoi(string(b))
	t := time.Unix(0, int64(ms)*int64(time.Millisecond))
	*epoch = Epoch(t)
	return nil
}

func (epoch *Epoch) MarshalJSON() ([]byte, error) {
	// @afiune we might have problems changing the location :(
	return epoch.ToTime().UTC().MarshalJSON()
}

// A few format functions for printing and manipulating the custom date
func (epoch Epoch) ToTime() time.Time {
	return time.Time(epoch)
}
func (epoch Epoch) Format(s string) string {
	return epoch.ToTime().Format(s)
}
func (epoch Epoch) UTC() time.Time {
	return epoch.ToTime().UTC()
}
