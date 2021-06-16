package lwtime

import (
	"strings"
	"time"
)

type RFC1123Z struct {
	time.Time
}

func (rfc *RFC1123Z) UnmarshalJSON(b []byte) (err error) {
	t := strings.Trim(string(b), `"`)
	rfc.Time, err = time.Parse(time.RFC1123Z, t)
	return
}

func (rfc *RFC1123Z) MarshalJSON() ([]byte, error) {
	return []byte(`"` + rfc.Time.Format(time.RFC1123Z) + `"`), nil
}
