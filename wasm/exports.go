package wasm

import (
	"encoding/json"
	"fmt"
	"github.com/bytecodealliance/wasmtime-go"
	"io/fs"
	"io/ioutil"
	"net/http"
	"strings"
)

func abort(a int32, b int32, c int32, d int32) {
	fmt.Println()
}

func readMemory(buf []byte, ptr int64, len int64) string {
	return string(buf[ptr : ptr+len])
}

func logging(store *wasmtime.Store, memory *wasmtime.Memory, ptr int64, length int64) {
	buf := memory.UnsafeData(store)

	msg := string(buf[ptr : ptr+length])

	fmt.Println(fmt.Sprintf("[log] %s", msg))
}

type HTTPResponse struct {
	ptr int64 `json:"ptr"`
}

type HTTPRequest struct {
	Verb     string                 `json:"verb"`
	URL      string                 `json:"url"`
	Headers  map[string]string      `json:"headers,omitempty"`
	Body     map[string]interface{} `json:"body,omitempty"`
	Response int64                  `json:"response,omitempty"`
}

func httpRequest(store *wasmtime.Store, memory *wasmtime.Memory, ptr int64, length int64) {
	buf := memory.UnsafeData(store)

	var r HTTPRequest
	err := json.Unmarshal(buf[ptr:ptr+length], &r)
	check(err)

	var req *http.Request

	body, err := json.Marshal(r.Body)
	check(err)

	req, err = http.NewRequest(r.Verb, r.URL, strings.NewReader(string(body)))
	check(err)

	for k, v := range r.Headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{}
	res, err := client.Do(req)

	rBody, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	check(err)

	for i := 0; i < len(rBody); i++ {
		buf[r.Response+int64(i)] = rBody[i]
	}
}

func writeFile(store *wasmtime.Store, memory *wasmtime.Memory, ptr int64, length int64, contentPtr int64, contentLen int64) {
	buf := memory.UnsafeData(store)

	path := string(buf[ptr : ptr+length])
	content := buf[contentPtr : contentPtr+contentLen]

	ioutil.WriteFile(path, content, fs.ModePerm)

	fmt.Println(path)
}

//func laceworkAPI(ptr int64, length int64) {
//    buf := memory.UnsafeData(store)
//
//    type LaceworkAPIRequest struct {
//        Function  string `json:"function"`
//        Arg1      string `json:"arg1"`
//        ResultPtr int64  `json:"result_ptr"`
//        ResultLen int64  `json:"result_len"`
//    }
//
//    var r LaceworkAPIRequest
//
//    err := json.Unmarshal(buf[ptr:ptr+length], &r)
//    check(err)
//
//    if err == nil {
//        cli.OutputJSON(&r)
//
//        var response api.TeamMemberResponse
//        err := cli.LwApi.V2.TeamMembers.Get("TECHALLY_DE894980E27BC66EEE46F65A585C4C588310B9CCDC531A9", &response)
//        check(err)
//
//        out, _ := json.Marshal(response.Data)
//
//        for i := 0; i < len(out); i++ {
//           buf[r.ResultPtr+int64(i)] = out[i]
//        }
//    }
//}
