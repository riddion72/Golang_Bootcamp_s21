package wind

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa
#include "application.h"
#include "window.h"
*/
import "C"
import "unsafe"

func InitApplication() {
	C.InitApplication()
}

func RunApplication() {
	C.RunApplication()
}

func CreateWindow(x, y, width, height int, title string) unsafe.Pointer {
	cSrting := C.CString(title)
	defer C.free(unsafe.Pointer(cSrting))

	return C.Window_Create(C.int(x), C.int(y), C.int(width), C.int(height), cSrting)
}

func MakeKeyAndOrderFront(window unsafe.Pointer) {
	C.Window_MakeKeyAndOrderFront(window)
}
