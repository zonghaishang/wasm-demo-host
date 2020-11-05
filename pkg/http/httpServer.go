package http

import (
	"github.com/zonghaishang/wasm-demo-host/pkg/handler"
	"log"
	"net/http"
)

func Server(host string) {
	log.Fatal(http.ListenAndServe(host, handler.NewWasmHandler()))
}
