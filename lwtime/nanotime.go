package lwtime

import "time"

// NanoTime time type to parse the returned time with nano format
//
// Example: "2020-08-20T01:00:00+0000"
type NanoTime time.Time

func (nano *NanoTime) UnmarshalJSON(b []byte) (err error) {
	s := string(b)
	t, err := time.Parse(time.RFC3339Nano, s[1:len(s)-1])
	if err != nil {
		t, err = time.Parse("2006-01-02T15:04:05.999999999Z0700", s[1:len(s)-1])
	}
	*nano = NanoTime(t)
	return
}

func (nano NanoTime) MarshalJSON() ([]byte, error) {
	// @afiune we might have problems changing the location :(
	return nano.ToTime().UTC().MarshalJSON()
}

func (nano NanoTime) ToTime() time.Time {
	return time.Time(nano)
}
func (nano NanoTime) Format(s string) string {
	return nano.ToTime().Format(s)
}
func (nano NanoTime) UTC() time.Time {
	return nano.ToTime().UTC()
}
