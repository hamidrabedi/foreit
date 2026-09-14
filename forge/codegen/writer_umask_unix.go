//go:build unix

package generator

import (
	"sync"
	"syscall"
)

var umaskMu sync.Mutex

// processUmask reads and restores the process-wide file creation mask.
func processUmask() int {
	umaskMu.Lock()
	defer umaskMu.Unlock()
	mask := syscall.Umask(0)
	syscall.Umask(mask)
	return mask
}
