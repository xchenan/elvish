package sys

/*
#include <unistd.h>
*/
// import "C"

import (
	"github.com/mattn/go-isatty"
)

func IsATTY(fd int) bool {
	return isatty.IsTerminal(uintptr(fd)) ||
		isatty.IsCygwinTerminal(uintptr(fd))
	// return C.isatty(C.int(fd)) != 0
}
