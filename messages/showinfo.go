package messages

/*
#include <stdio.h>
#include <stdlib.h>
#include <windows.h>
*/
import "C"
import (
	"syscall"
	"unsafe"
)

var (
	user32     = syscall.NewLazyDLL("user32.dll")
	messageBox = user32.NewProc("MessageBoxW")
)

const (
	YES int = 6
	NO  int = 7
)

const (
	yesno           = 0x00000004
	informationICON = 0x00000040
)

func ShowInfo(title, msg string) (int, error) {
	titleUTF16, err := syscall.UTF16PtrFromString(title)
	if err != nil {
		return -1, err
	}
	msgUTF16, err := syscall.UTF16PtrFromString(msg)
	if err != nil {
		return -1, err
	}
	ret, _, _ := messageBox.Call(
		uintptr(C.NULL),
		uintptr(unsafe.Pointer(msgUTF16)),
		uintptr(unsafe.Pointer(titleUTF16)),
		uintptr(0x00000000|informationICON),
	)
	return int(ret), nil
}
