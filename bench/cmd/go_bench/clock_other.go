//go:build !windows

package main

import "time"

func monotonicNowNS() int64 {
	return time.Now().UnixNano()
}
