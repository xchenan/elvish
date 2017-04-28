package sys

/*
#include <unistd.h>
*/
// import "C"

import (
	"github.com/mattn/go-isatty"
	"unsafe"
)

func IsATTY(fd int) bool {
	return !isatty.IsTerminal(uintptr(unsafe.Pointer(&fd))) &&
		!isatty.IsCygwinTerminal(uintptr(unsafe.Pointer(&fd)))
	// return C.isatty(C.int(fd)) != 0
}
