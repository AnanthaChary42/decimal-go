//go:build !windows

package main

import "runtime"

// Go's standard library has no portable resident-set API. This fallback keeps
// the benchmark runnable elsewhere while Windows reports the true working set.
func currentRSSKB() uint64 {
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)
	return stats.Sys / 1024
}
