package sys

/*
#include <termios.h>
*/
// import "C"
import (
	"unsafe"

	"golang.org/x/sys/unix"
)

// Termios represents terminal attributes.
// type Termios C.struct_termios
type Termios unix.Termios

func ioctl(fd, cmd uintptr, arg unsafe.Pointer) error {
	return ioctlu(fd, cmd, uintptr(arg))
}

func ioctlu(fd, cmd, arg uintptr) error {
	_, _, errno := unix.Syscall(unix.SYS_IOCTL, fd, cmd, arg)
	if errno == 0 {
		return nil
	}
	return errno
}

// NewTermiosFromFd extracts the terminal attribute of the given file
// descriptor.
func NewTermiosFromFd(fd int) (*Termios, error) {
	var tio Termios
	if err := ioctl(uintptr(fd),
		getAttrIOCTL, unsafe.Pointer(&tio)); err != nil {
		return nil, err
	}

	return &tio, nil

	//	term := new(Termios)
	//	err := term.FromFd(fd)
	//	if err != nil {
	//		return nil, err
	//	}
	//	return term, nil
}

//func (term *Termios) c() *C.struct_termios {
//	return (*C.struct_termios)(term)
//}
//
//// FromFd fills term with the terminal attribute of the given file descriptor.
//func (term *Termios) FromFd(fd int) error {
//	_, err := C.tcgetattr((C.int)(fd), term.c())
//	return err
//}

// ApplyToFd applies term to the given file descriptor.
func (term *Termios) ApplyToFd(fd int) error {
	return ioctl(uintptr(fd), setAttrNowIOCTL, unsafe.Pointer(&term))
	//	_, err := C.tcsetattr((C.int)(fd), 0, term.c())
	//	return err
}

// Copy returns a copy of term.
func (term *Termios) Copy() *Termios {
	v := *term
	return &v
}

// SetVTime sets the timeout in deciseconds for noncanonical read.
func (term *Termios) SetVTime(v uint8) {
	term.Cc[unix.VTIME] = v
}

// SetVMin sets the minimal number of characters for noncanonical read.
func (term *Termios) SetVMin(v uint8) {
	term.Cc[unix.VMIN] = v
}

func setFlag(flag *uint32, mask uint32, v bool) {
	if v {
		*flag |= mask
	} else {
		*flag &= ^mask
	}
}

// SetICanon sets the canonical flag.
func (term *Termios) SetICanon(v bool) {
	setFlag(&term.Lflag, unix.ICANON, v)
}

// SetEcho sets the echo flag.
func (term *Termios) SetEcho(v bool) {
	setFlag(&term.Lflag, unix.ECHO, v)
}

type Queue uint8

const (
	AnyQueue Queue = iota
	InputQueue
	OutputQueue
)

// Flush discards data written to the terminal but not transmitted, or
// received from the terminal but not read, depending on the queue.
//
// InputQueue   received but not read.
// OutputQueue  written but not transmitted.
// AnyQueue     both

// FlushInput discards data written to a file descriptor but not read.
//func FlushInput(fd int) error {
//	return ioctlu(uintptr(unsafe.Pointer(&fd)), flushIOCTL, queue.bits())
//}
