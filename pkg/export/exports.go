package export

import (
	"unsafe"
)

//export proxy_on_host_hello
func ProxyHello(ctx unsafe.Pointer, ptr int32, len int32) {
	//return fmt.Sprintf("Welcome %s", name)

}
