package handler

import (
	"fmt"
	wasm "github.com/wasmerio/wasmer-go/wasmer"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"runtime"
)

type WasmHandler struct {
}

func NewWasmHandler() *WasmHandler {
	return &WasmHandler{}
}

func (h *WasmHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	bytes, err := ioutil.ReadFile(simpleWasmFile())
	if err != nil {
		w.Write([]byte(fmt.Sprintf("No wasm file found for path '%s'", simpleWasmFile())))
		return
	}

	store := wasm.NewStore(wasm.NewEngine())
	module, err := wasm.NewModule(store, bytes)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Failed to init wasm plugin '%s', err: %v", simpleWasmFile(), err)))
		return
	}

	instance, err := wasm.NewInstance(module, wasm.NewImportObject())
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Failed to create wasm plugin instance '%s', err: %v", simpleWasmFile(), err)))
		return
	}

	// Gets the `proxy_on_plugin_hello` exported function from the WebAssembly instance.
	helloFunc, err := instance.Exports.GetFunction("proxy_on_plugin_hello")
	if err != nil {
		if err != nil {
			w.Write([]byte(fmt.Sprintf("Not found proxy_on_plugin_hello function for plugin '%s', err: %v", simpleWasmFile(), err)))
			return
		}
	}

	name := h.getPersonName(r)

	memory, err := instance.Exports.GetMemory("memory")
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Failed to allocate memory for plugin '%s', err: %v", simpleWasmFile(), err)))
		return
	}
	buf := memory.Data()
	// copy to instance buffer
	copy(buf, []byte(name))

	fmt.Fprintf(os.Stdout, "buf address: %p", &buf)

	var (
		pos  int32
		size int32
	)

	helloFunc(&buf, len(name), &pos, &size)

	ioBuf := memory.Data()
	result := string(ioBuf[pos : pos+size])
	fmt.Fprintf(os.Stdout, "ioBuf: %s", result)

	// output for response
	w.Write([]byte(result))
}

func (h *WasmHandler) getPersonName(r *http.Request) string {
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "诣极"
	}
	return name
}

func simpleWasmFile() string {
	_, filename, _, _ := runtime.Caller(0)
	return path.Join(path.Dir(filename), "wasm", "simple.wasm")
}
