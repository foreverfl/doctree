package cmd

import (
	"os"
	"strconv"
	"syscall"
	"time"
)

// readPid reads a pid file. Returns (0, false) if missing or malformed.
func readPid(path string) (int, bool) {
	b, err := os.ReadFile(path)
	if err != nil {
		return 0, false
	}
	pid, err := strconv.Atoi(string(b))
	if err != nil || pid <= 0 {
		return 0, false
	}
	return pid, true
}

// processAlive returns true if a process with pid exists and we can signal it.
// kill(pid, 0) is the standard liveness probe on POSIX.
func processAlive(pid int) bool {
	return syscall.Kill(pid, 0) == nil
}

// waitExit polls processAlive until the process is gone or timeout elapses.
func waitExit(pid int, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if !processAlive(pid) {
			return true
		}
		time.Sleep(100 * time.Millisecond)
	}
	return !processAlive(pid)
}