package api

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	Foo string `json:"foo"`
	Bar int    `json:"bar"`
}

func TestJsonReader(t *testing.T) {
	var subject = testStruct{"foo", 1}

	reader, err := jsonReader(subject)
	if assert.Nil(t, err) {
		buf := new(bytes.Buffer)
		buf.ReadFrom(reader)
		assert.Equal(t,
			"{\"foo\":\"foo\",\"bar\":1}\n",
			buf.String(),
			"unexpected streaming encoder")
	}
}
