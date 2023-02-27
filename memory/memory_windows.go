//go:build windows

package memory

import (
	"golang.org/x/sys/windows"
	"unsafe"
)

type memoryStatusEx struct {
	cbSize                  uint32
	dwMemoryLoad            uint32
	ullTotalPhys            uint64 // in bytes
	ullAvailPhys            uint64
	ullTotalPageFile        uint64
	ullAvailPageFile        uint64
	ullTotalVirtual         uint64
	ullAvailVirtual         uint64
	ullAvailExtendedVirtual uint64
}

func Stats() (mem *Memory, err error) {
	procGlobalMemoryStatusEx := windows.NewLazySystemDLL("kernel32.dll").NewProc("GlobalMemoryStatusEx")
	var memInfo memoryStatusEx
	memInfo.cbSize = uint32(unsafe.Sizeof(memInfo))
	v, _, _ := procGlobalMemoryStatusEx.Call(uintptr(unsafe.Pointer(&memInfo)))
	if v == 0 {
		err = windows.GetLastError()
		return
	}
	mem = &Memory{
		Total:       memInfo.ullTotalPhys,
		Free:        memInfo.ullAvailPhys,
		Available:   memInfo.ullAvailPhys,
		Used:        memInfo.ullTotalPhys - memInfo.ullAvailPhys,
		UsedPercent: float64(memInfo.dwMemoryLoad),
	}
	fillRuntime(mem)
	return
}
