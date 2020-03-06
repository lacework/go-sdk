package client

import (
	"bytes"
	"encoding/json"
	"io"
)

// jsonReader takes any arbitrary type and synthesizes a streaming encoder
func jsonReader(v interface{}) (r io.Reader, err error) {
	buf := new(bytes.Buffer)
	err = json.NewEncoder(buf).Encode(v)
	r = bytes.NewReader(buf.Bytes())
	return
}
