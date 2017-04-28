package sys

/*
#include <termios.h>
#include <sys/ioctl.h>

void getwinsize(int fd, int *row, int *col) {
	struct winsize wsz;
	ioctl(fd, TIOCGWINSZ, &wsz);
	*row = wsz.ws_row;
	*col = wsz.ws_col;
}
*/
// import "C"

import (
	"unsafe"

	"golang.org/x/sys/unix"
)

// GetWinsize queries the size of the terminal referenced by the given file
// descriptor.

type winSize struct {
	row    uint16
	col    uint16
	xpixel uint16 // unused
	Ypixel uint16 // unused
}

func GetWinsize(fd int) (row int, col int, err error) {
	ws := winSize{}

	if err := ioctl(uintptr(fd),
		unix.TIOCGWINSZ, unsafe.Pointer(&ws)); err != nil {
		return -1, -1, err
	}
	return int(ws.row), int(ws.col), nil
	//	var r, c C.int
	//	C.getwinsize(C.int(fd), &r, &c)
	//	return int(r), int(c)
}
