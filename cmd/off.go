package cmd

import (
	"errors"
	"fmt"
	"os"
	"syscall"
	"time"

	"github.com/foreverfl/doctree/internal/daemon"
	"github.com/spf13/cobra"
)

const offWaitTimeout = 3 * time.Second

var offCmd = &cobra.Command{
	Use:   "off",
	Short: "Stop the doctree daemon",
	RunE: func(cmd *cobra.Command, args []string) error {
		sockpath, err := sockPath()
		if err != nil {
			return err
		}
		pidpath, err := pidPath()
		if err != nil {
			return err
		}

		pid, hasPid := readPid(pidpath)
		if !hasPid || !processAlive(pid) {
			_ = os.Remove(sockpath)
			_ = os.Remove(pidpath)
			fmt.Println("doctree daemon not running")
			return nil
		}

		// Try graceful shutdown via RPC. If the call fails we still fall
		// through to wait/SIGTERM — the goal is to leave no daemon behind.
		_, callErr := daemon.Call(sockpath, daemon.Request{Op: daemon.OpShutdown})

		if waitExit(pid, offWaitTimeout) {
			_ = os.Remove(sockpath)
			_ = os.Remove(pidpath)
			fmt.Printf("doctree daemon stopped (pid=%d)\n", pid)
			return nil
		}

		if callErr != nil && !errors.Is(callErr, daemon.ErrNotRunning) {
			fmt.Fprintf(os.Stderr, "rpc shutdown failed: %v; sending SIGTERM\n", callErr)
		}
		if err := syscall.Kill(pid, syscall.SIGTERM); err != nil {
			return fmt.Errorf("sigterm pid=%d: %w", pid, err)
		}
		if !waitExit(pid, offWaitTimeout) {
			return fmt.Errorf("daemon did not exit after SIGTERM (pid=%d)", pid)
		}
		_ = os.Remove(sockpath)
		_ = os.Remove(pidpath)
		fmt.Printf("doctree daemon stopped via SIGTERM (pid=%d)\n", pid)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(offCmd)
}
