package cs

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"syscall"
	"unsafe"
)

func WriteProcessMemory[T any](hProcess HANDLE, lpBaseAddress uint32, data T) (err error) {
	var buf bytes.Buffer

	err = binary.Write(&buf, binary.LittleEndian, data)

	if err != nil {
		fmt.Println(err)
		return
	}

	bData := buf.Bytes()

	size := unsafe.Sizeof(data)
	var numBytesRead uintptr
	_, _, err = procWriteProcessMemory.Call(uintptr(hProcess),
		uintptr(lpBaseAddress),
		uintptr(unsafe.Pointer(&bData[0])),
		size,
		uintptr(unsafe.Pointer(&numBytesRead)))

	if !IsErrSuccess(err) {
		return
	}
	err = nil

	return
}

func ReadProcessMemory[T any](hProcess HANDLE, lpBaseAddress uint32, receiver *T) (err error) {
	var numBytesRead uintptr
	size := uint(unsafe.Sizeof(*receiver))
	data := make([]byte, size)

	_, _, err = procReadProcessMemory.Call(uintptr(hProcess),
		uintptr(lpBaseAddress),
		uintptr(unsafe.Pointer(&data[0])),
		uintptr(size),
		uintptr(unsafe.Pointer(&numBytesRead)))
	if !IsErrSuccess(err) {
		return
	}

	buf := bytes.NewBuffer(data)
	err = binary.Read(buf, binary.LittleEndian, receiver)

	if err != nil {
		return err
	}

	return
}

func IsErrSuccess(err error) bool {
	if errno, ok := err.(syscall.Errno); ok {
		if errno == 0 {
			return true
		}
	}
	return false
}
