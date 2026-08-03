//go:build windows

package main

import (
	"syscall"
	"unsafe"
)

type processMemoryCounters struct {
	CB                         uint32
	PageFaultCount             uint32
	PeakWorkingSetSize         uintptr
	WorkingSetSize             uintptr
	QuotaPeakPagedPoolUsage    uintptr
	QuotaPagedPoolUsage        uintptr
	QuotaPeakNonPagedPoolUsage uintptr
	QuotaNonPagedPoolUsage     uintptr
	PagefileUsage              uintptr
	PeakPagefileUsage          uintptr
}

var (
	psapi                    = syscall.NewLazyDLL("psapi.dll")
	getProcessMemoryInfoProc = psapi.NewProc("GetProcessMemoryInfo")
)

func currentRSSKB() uint64 {
	counters := processMemoryCounters{CB: uint32(unsafe.Sizeof(processMemoryCounters{}))}
	handle, err := syscall.GetCurrentProcess()
	if err != nil {
		return 0
	}
	r1, _, _ := getProcessMemoryInfoProc.Call(
		uintptr(handle),
		uintptr(unsafe.Pointer(&counters)),
		uintptr(counters.CB),
	)
	if r1 == 0 {
		return 0
	}
	return uint64(counters.WorkingSetSize) / 1024
}
