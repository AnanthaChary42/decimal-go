//go:build windows

package main

import (
	"syscall"
	"unsafe"
)

var (
	kernel32                   = syscall.NewLazyDLL("kernel32.dll")
	queryPerformanceCounter   = kernel32.NewProc("QueryPerformanceCounter")
	queryPerformanceFrequency = kernel32.NewProc("QueryPerformanceFrequency")
	performanceFrequency       = queryCounterFrequency()
)

func queryCounterFrequency() int64 {
	var frequency int64
	r1, _, _ := queryPerformanceFrequency.Call(uintptr(unsafe.Pointer(&frequency)))
	if r1 == 0 || frequency <= 0 {
		panic("QueryPerformanceFrequency failed")
	}
	return frequency
}

// monotonicNowNS uses QueryPerformanceCounter rather than Windows wall-clock
// timestamps, whose effective resolution can produce zero-length samples.
func monotonicNowNS() int64 {
	var counter int64
	r1, _, _ := queryPerformanceCounter.Call(uintptr(unsafe.Pointer(&counter)))
	if r1 == 0 {
		panic("QueryPerformanceCounter failed")
	}
	seconds := counter / performanceFrequency
	nanoseconds := (counter % performanceFrequency) * 1e9 / performanceFrequency
	return seconds*1e9 + nanoseconds
}
