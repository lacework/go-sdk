package lwtime

import (
	"strings"
	"time"
)

type RFC1123Z time.Time

func (rfc *RFC1123Z) UnmarshalJSON(b []byte) (err error) {
	t := strings.Trim(string(b), `"`)
	res, _ := time.Parse(time.RFC1123Z, t)
	*rfc = RFC1123Z(res)
	return
}

func (rfc *RFC1123Z) MarshalJSON() ([]byte, error) {
	return rfc.ToTime().UTC().MarshalJSON()
}

func (rfc RFC1123Z) ToTime() time.Time {
	return time.Time(rfc)
}
func (rfc RFC1123Z) Format(s string) string {
	return rfc.ToTime().Format(s)
}
func (rfc RFC1123Z) UTC() time.Time {
	return rfc.ToTime().UTC()
}
