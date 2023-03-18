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

func AskYesNo(title, msg string) (bool, error) {
	titleUTF16, err := syscall.UTF16PtrFromString(title)
	if err != nil {
		return false, err
	}
	msgUTF16, err := syscall.UTF16PtrFromString(msg)
	if err != nil {
		return false, err
	}
	ret, _, _ := messageBox.Call(
		uintptr(C.NULL),
		uintptr(unsafe.Pointer(msgUTF16)),
		uintptr(unsafe.Pointer(titleUTF16)),
		uintptr(yesno|0x00000030),
	)
	return int(ret) == 6, nil
}
