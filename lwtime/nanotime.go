package lwtime

import "time"

// time type to parse the returned time with nano format
//
// Example:
//
// "START_TIME":"2020-08-20T01:00:00+0000"
type NanoTime time.Time

func (self *NanoTime) UnmarshalJSON(b []byte) (err error) {
	s := string(b)
	t, err := time.Parse(time.RFC3339Nano, s[1:len(s)-1])
	if err != nil {
		t, err = time.Parse("2006-01-02T15:04:05.999999999Z0700", s[1:len(s)-1])
	}
	*self = NanoTime(t)
	return
}

func (self NanoTime) MarshalJSON() ([]byte, error) {
	// @afiune we might have problems changing the location :(
	return self.ToTime().UTC().MarshalJSON()
}

// A few format functions for printing and manipulating the custom date
func (self NanoTime) ToTime() time.Time {
	return time.Time(self)
}
func (self NanoTime) Format(s string) string {
	return self.ToTime().Format(s)
}
func (self NanoTime) UTC() time.Time {
	return self.ToTime().UTC()
}
