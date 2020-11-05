package main

import "github.com/zonghaishang/wasm-demo-host/pkg/http"

func main() {
	http.Server(":8080")
}
