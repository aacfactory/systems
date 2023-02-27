//go:build windows

package stat

import (
	"fmt"
	"golang.org/x/sys/windows"
	"unsafe"
)

type win32_SystemProcessorPerformanceInformation struct {
	IdleTime       int64
	KernelTime     int64
	UserTime       int64
	DpcTime        int64
	InterruptTime  int64
	InterruptCount uint32
}

const (
	ClocksPerSec                                     = 10000000.0
	win32_SystemProcessorPerformanceInformationClass = 8
	win32_SystemProcessorPerformanceInfoSize         = uint32(unsafe.Sizeof(win32_SystemProcessorPerformanceInformation{}))
)

func perfInfo() (infos []win32_SystemProcessorPerformanceInformation, err error) {
	maxBuffer := 2056
	resultBuffer := make([]win32_SystemProcessorPerformanceInformation, maxBuffer)
	bufferSize := uintptr(win32_SystemProcessorPerformanceInfoSize) * uintptr(maxBuffer)
	var retSize uint32
	var retCode uintptr
	retCode, _, err = windows.NewLazySystemDLL("ntdll.dll").NewProc("NtQuerySystemInformation").Call(
		win32_SystemProcessorPerformanceInformationClass,
		uintptr(unsafe.Pointer(&resultBuffer[0])),
		bufferSize,
		uintptr(unsafe.Pointer(&retSize)),
	)
	if retCode != 0 {
		err = fmt.Errorf("call to NtQuerySystemInformation returned %d. err: %s", retSize, err.Error())
		return
	}
	numReturnedElements := retSize / win32_SystemProcessorPerformanceInfoSize
	resultBuffer = resultBuffer[:numReturnedElements]
	return resultBuffer, nil
}

func Times() (times []TimesStat, err error) {
	stats, infoErr := perfInfo()
	if infoErr != nil {
		err = infoErr
		return
	}
	times = make([]TimesStat, 0, 1)
	for core, v := range stats {
		c := TimesStat{
			CPU:    fmt.Sprintf("%d", core),
			User:   float64(v.UserTime) / ClocksPerSec,
			System: float64(v.KernelTime-v.IdleTime) / ClocksPerSec,
			Idle:   float64(v.IdleTime) / ClocksPerSec,
			Irq:    float64(v.InterruptTime) / ClocksPerSec,
		}
		times = append(times, c)
	}
	return
}
