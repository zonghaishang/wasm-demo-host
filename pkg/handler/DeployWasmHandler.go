package handler

import (
	"fmt"
	wasm "github.com/wasmerio/go-ext-wasm/wasmer"
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

	// Reads the WebAssembly module as bytes.
	bytes, err := wasm.ReadBytes(simpleWasmFile())
	if err != nil {
		w.Write([]byte(fmt.Sprintf("No wasm file found for path '%s'", simpleWasmFile())))
		return
	}

	// Instantiates the WebAssembly module.
	instance, err := wasm.NewInstance(bytes)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Failed to create wasm plugin instance '%s', err: %v", simpleWasmFile(), err)))
		return
	}

	// Close the WebAssembly instance later.
	defer instance.Close()

	// Gets the `sum` exported function from the WebAssembly instance.
	sum := instance.Exports["sum"]
	if err != nil {
		if err != nil {
			w.Write([]byte(fmt.Sprintf("Not found sum function for plugin '%s', err: %v", simpleWasmFile(), err)))
			return
		}
	}

	// _ := h.getPersonName(r)

	//memory, err := instance.Exports.GetMemory("memory")
	//if err != nil {
	//	w.Write([]byte(fmt.Sprintf("Failed to allocate memory for plugin '%s', err: %v", simpleWasmFile(), err)))
	//	return
	//}
	//buf := memory.Data()
	//// copy to instance buffer
	//copy(buf, []byte(name))

	// fmt.Fprintf(os.Stdout, "buf address: %p", &buf)

	//var (
	//	pos  int32
	//	size int32
	//)

	result, _ := sum(1, 1)

	fmt.Fprintf(os.Stdout, "method invoke complete")

	//ioBuf := memory.Data()
	//result := string(ioBuf[pos : pos+size])
	//fmt.Fprintf(os.Stdout, "ioBuf: %s", result)
	//
	//// output for response
	w.Write([]byte("1 + 1 = " + result.String()))
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
