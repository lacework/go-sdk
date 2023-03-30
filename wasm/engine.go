package wasm

import (
	"bytes"
	"github.com/bytecodealliance/wasmtime-go"
	"io/ioutil"
	"runtime"
	"strings"
)

type Engine struct {
	engine *wasmtime.Engine
	memory *wasmtime.Memory
	store  *wasmtime.Store
}

func NewEngine() *Engine {
	engine := wasmtime.NewEngine()
	store := wasmtime.NewStore(engine)

	return &Engine{engine: engine, store: store}
}

func (e *Engine) decorate(f func(store *wasmtime.Store, memory *wasmtime.Memory, x int64, y int64)) func(int64, int64) {
	return func(x int64, y int64) {
		f(e.store, e.memory, x, y)
	}
}

func (e *Engine) Commands(wasmFile string) []string {
	wasm, err := ioutil.ReadFile(wasmFile)
	check(err)

	module, err := wasmtime.NewModule(e.engine, wasm)
	check(err)

	var cmds []string

	for _, e := range module.Exports() {
		t := e.Type().FuncType()
		if t != nil {
			cmds = append(cmds, e.Name())
		}
	}

	return cmds
}

func (e *Engine) Run(wasmFile string, function string, args []string) {
	wasm, err := ioutil.ReadFile(wasmFile)
	check(err)

	module, err := wasmtime.NewModule(e.engine, wasm)
	check(err)

	extern := []wasmtime.AsExtern{
		wasmtime.WrapFunc(e.store, abort),
		wasmtime.WrapFunc(e.store, e.decorate(logging)),
		wasmtime.WrapFunc(e.store, e.decorate(cliOutput)),
		wasmtime.WrapFunc(e.store, e.decorate(httpRequest)),
	}

	instance, err := wasmtime.NewInstance(e.store, module, extern)
	check(err)

	run := instance.GetFunc(e.store, function)
	if run == nil {
		panic("not a function")
	}

	e.memory = instance.GetExport(e.store, "memory").Memory()

	// Hardcoded for now
	if function == "chat" {
		msg := strings.Join(args, " ")

		buf := bytes.NewBuffer(e.memory.UnsafeData(e.store))
		buf.Reset()
		buf.Write([]byte(msg))

		runtime.KeepAlive(buf)

		var ptr int64 = 0
		var length = int32(len(msg))

		_, err = run.Call(e.store, ptr, length)
	} else {
		_, err = run.Call(e.store)
	}
	check(err)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
