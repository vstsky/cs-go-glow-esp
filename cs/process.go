package cs

import (
	"errors"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"unsafe"
)

var processName = "csgo.exe"
var moduleName = "client.dll"

var handle HANDLE
var clientDllModuleAddress uintptr

func getModuleBaseAddress(pid uint32, moduleName string) (uintptr, error) {
	snap := CreateToolhelp32Snapshot(TH32CS_SNAPMODULE, pid)
	if snap == 0 {
		return 0, errors.New("snapshot could not be created")
	}
	defer CloseHandle(snap)

	var me32 MODULEENTRY32
	me32.Size = uint32(unsafe.Sizeof(me32))

	for Module32Next(snap, &me32) {
		if moduleName == syscall.UTF16ToString(me32.SzModule[:]) {
			return uintptr(unsafe.Pointer(me32.ModBaseAddr)), nil
		}
	}

	return 0, errors.New("couldn't retrieve module. Check if GOARCH = x86")
}

func findProcessPid(pn string) (uint32, error) {
	plb, _ := exec.Command("tasklist").CombinedOutput()
	pls := string(plb)
	pl := strings.Split(pls, "\r\n")

	for _, process := range pl {
		processParts := strings.Fields(process)
		if len(processParts) > 0 && processParts[0] == pn {
			u64, _ := strconv.ParseUint(processParts[1], 10, 32)
			return uint32(u64), nil
		}
	}
	return 0, errors.New("Unable to retrieve process PID")
}
