package cs

import (
	"syscall"
	"unsafe"
)

var (
	mod                          = syscall.NewLazyDLL("kernel32.dll")
	procReadProcessMemory        = mod.NewProc("ReadProcessMemory")
	procWriteProcessMemory       = mod.NewProc("WriteProcessMemory")
	procOpenProcess              = mod.NewProc("OpenProcess")
	procCreateToolhelp32Snapshot = mod.NewProc("CreateToolhelp32Snapshot")
	procCloseHandle              = mod.NewProc("CloseHandle")
	procModule32Next             = mod.NewProc("Module32NextW")
)

const (
	PROCESS_VM_READ    = 0x0010
	PROCESS_VM_WRITE   = 0x0020
	PROCESS_ALL_ACCESS = 2035711
)

const (
	TH32CS_SNAPHEAPLIST = 0x00000001
	TH32CS_SNAPPROCESS  = 0x00000002
	TH32CS_SNAPTHREAD   = 0x00000004
	TH32CS_SNAPMODULE   = 0x00000008
	TH32CS_SNAPMODULE32 = 0x00000010
	TH32CS_INHERIT      = 0x80000000
	TH32CS_SNAPALL      = TH32CS_SNAPHEAPLIST | TH32CS_SNAPMODULE | TH32CS_SNAPPROCESS | TH32CS_SNAPTHREAD
)

const (
	MAX_MODULE_NAME32 = 255
	MAX_PATH          = 260
)

type MODULEENTRY32 struct {
	Size         uint32
	ModuleID     uint32
	ProcessID    uint32
	GlblcntUsage uint32
	ProccntUsage uint32
	ModBaseAddr  *uint8
	ModBaseSize  uint32
	HModule      HMODULE
	SzModule     [MAX_MODULE_NAME32 + 1]uint16
	SzExePath    [MAX_PATH]uint16
}

type HANDLE uintptr
type HMODULE HANDLE

func OpenProcess(desiredAccess uint32, inheritHandle bool, processId uint32) (handle HANDLE, err error) {
	inherit := 0
	if inheritHandle {
		inherit = 1
	}

	ret, _, err := procOpenProcess.Call(
		uintptr(desiredAccess),
		uintptr(inherit),
		uintptr(processId))
	if err != nil && IsErrSuccess(err) {
		err = nil
	}
	handle = HANDLE(ret)
	return
}

func CreateToolhelp32Snapshot(flags, processId uint32) HANDLE {
	ret, _, _ := procCreateToolhelp32Snapshot.Call(
		uintptr(flags),
		uintptr(processId))

	if ret <= 0 {
		return HANDLE(0)
	}

	return HANDLE(ret)
}

func CloseHandle(object HANDLE) bool {
	ret, _, _ := procCloseHandle.Call(
		uintptr(object))
	return ret != 0
}

func Module32Next(snapshot HANDLE, me *MODULEENTRY32) bool {
	ret, _, _ := procModule32Next.Call(
		uintptr(snapshot),
		uintptr(unsafe.Pointer(me)))

	return ret != 0
}
