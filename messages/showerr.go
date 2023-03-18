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

const (
	errorICON = 0x00000010
	yesonly   = 0x00000000
)

func ShowError(title, msg string) (int, error) {
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
		uintptr(yesonly|errorICON),
	)
	return int(ret), nil
}
