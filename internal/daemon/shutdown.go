package daemon

import (
	"errors"
	"fmt"
	"io"
	"os"
	"syscall"
	"time"

	"github.com/foreverfl/doctree/internal/process"
)

const stopTimeout = 3 * time.Second

// Shutdown stops the daemon if it is running and cleans up the sock/pid files.
// Status messages are written to out, warnings to errw. Returns nil when the
// daemon wasn't running or was stopped successfully; non-nil only when the
// SIGTERM fallback also failed.
func Shutdown(sockpath, pidpath string, out, errw io.Writer) error {
	pid, hasPid := process.ReadPid(pidpath)
	if !hasPid || !process.Alive(pid) {
		_ = os.Remove(sockpath)
		_ = os.Remove(pidpath)
		fmt.Fprintln(out, "doctree daemon not running")
		return nil
	}

	_, callErr := Call(sockpath, Request{Op: OpShutdown})

	if process.WaitExit(pid, stopTimeout) {
		_ = os.Remove(sockpath)
		_ = os.Remove(pidpath)
		fmt.Fprintf(out, "doctree daemon stopped (pid=%d)\n", pid)
		return nil
	}

	if callErr != nil && !errors.Is(callErr, ErrNotRunning) {
		fmt.Fprintf(errw, "rpc shutdown failed: %v; sending SIGTERM\n", callErr)
	}
	if err := syscall.Kill(pid, syscall.SIGTERM); err != nil {
		return fmt.Errorf("sigterm pid=%d: %w", pid, err)
	}
	if !process.WaitExit(pid, stopTimeout) {
		return fmt.Errorf("daemon did not exit after SIGTERM (pid=%d)", pid)
	}
	_ = os.Remove(sockpath)
	_ = os.Remove(pidpath)
	fmt.Fprintf(out, "doctree daemon stopped via SIGTERM (pid=%d)\n", pid)
	return nil
}