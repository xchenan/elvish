package sys

/*
#include <unistd.h>
#include <errno.h>

pid_t get(int fd) {
	return tcgetpgrp(fd);
}

int set(int fd, pid_t pid) {
	return tcsetpgrp(fd, pid);
}

int e() {
	return errno;
}
*/
// import "C"
import (
	"syscall"
	"unsafe"
)

func Tcgetpgrp(fd int) (int, error) {
	var pid int
	_, _, err := syscall.RawSyscall(syscall.SYS_IOCTL, uintptr(unsafe.Pointer(&fd)),
		uintptr(syscall.TIOCGPGRP), uintptr(unsafe.Pointer(&pid)))
	return pid, err
	//	i := C.get(C.int(fd))
	//	if i == -1 {
	//		return -1, syscall.Errno(C.e())
	//	}
	//	return int(i), nil
}

func Tcsetpgrp(fd int, pid int) error {
	_, _, err := syscall.RawSyscall(syscall.SYS_IOCTL, uintptr(unsafe.Pointer(&fd)),
		uintptr(syscall.TIOCSPGRP), uintptr(unsafe.Pointer(&pid)))
	//	i := C.set(C.int(fd), C.pid_t(pid))
	//	if i != 0 {
	//		return syscall.Errno(C.e())
	//	}
	return err
}
