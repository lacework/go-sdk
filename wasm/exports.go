package wasm

import (
	"fmt"
	"github.com/bytecodealliance/wasmtime-go"
)

func abort(a int32, b int32, c int32, d int32) {
	fmt.Println()
}

func logging(store *wasmtime.Store, memory *wasmtime.Memory, ptr int64, length int64) {
	buf := memory.UnsafeData(store)

	msg := string(buf[ptr : ptr+length])

	fmt.Println(fmt.Sprintf("[log] %s", msg))
}

func cliOutput(store *wasmtime.Store, memory *wasmtime.Memory, ptr int64, length int64) {
	buf := memory.UnsafeData(store)

	msg := string(buf[ptr : ptr+length])

	fmt.Println(msg)
}

func httpRequest(store *wasmtime.Store, memory *wasmtime.Memory, ptr int64, length int64) {
	buf := memory.UnsafeData(store)

	url := string(buf[ptr : ptr+length])

	fmt.Println(url)
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
