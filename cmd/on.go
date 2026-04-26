package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"syscall"
	"time"

	"github.com/foreverfl/doctree/internal/daemon"
	"github.com/spf13/cobra"
)

const (
	onReadyTimeout = 5 * time.Second
	onPollInterval = 100 * time.Millisecond
)

var onCmd = &cobra.Command{
	Use:   "on",
	Short: "Start the doctree daemon",
	RunE: func(cmd *cobra.Command, args []string) error {
		sockpath, err := sockPath()
		if err != nil {
			return err
		}
		pidpath, err := pidPath()
		if err != nil {
			return err
		}
		logpath, err := logPath()
		if err != nil {
			return err
		}

		// Already running? Skip.
		if pid, ok := readPid(pidpath); ok && processAlive(pid) {
			if err := daemon.Ping(sockpath); err == nil {
				fmt.Printf("doctree daemon already running (pid=%d)\n", pid)
				return nil
			}
		}

		// Clean stale state from a previous crash.
		_ = os.Remove(pidpath)
		_ = os.Remove(sockpath)

		self, err := os.Executable()
		if err != nil {
			return fmt.Errorf("locate self: %w", err)
		}

		logFile, err := os.OpenFile(logpath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
		if err != nil {
			return fmt.Errorf("open log: %w", err)
		}

		c := exec.Command(self, "daemon-run")
		c.Stdout = logFile
		c.Stderr = logFile
		// Setsid detaches the daemon into its own session so it survives the
		// parent shell.
		c.SysProcAttr = &syscall.SysProcAttr{Setsid: true}

		if err := c.Start(); err != nil {
			_ = logFile.Close()
			return fmt.Errorf("start daemon: %w", err)
		}
		// Parent owns the fd until Start succeeds; once handed off, the child
		// keeps it open via dup.
		_ = logFile.Close()

		pid := c.Process.Pid
		if err := os.WriteFile(pidpath, []byte(strconv.Itoa(pid)), 0o644); err != nil {
			return fmt.Errorf("write pid: %w", err)
		}
		// Release so the OS reaps the child instead of us.
		if err := c.Process.Release(); err != nil {
			return fmt.Errorf("release: %w", err)
		}

		// Poll for sock readiness.
		deadline := time.Now().Add(onReadyTimeout)
		for time.Now().Before(deadline) {
			if err := daemon.Ping(sockpath); err == nil {
				fmt.Printf("doctree daemon started (pid=%d)\n", pid)
				return nil
			} else if !errors.Is(err, daemon.ErrNotRunning) {
				return err
			}
			time.Sleep(onPollInterval)
		}
		return fmt.Errorf("daemon failed to become ready within %s. see %s", onReadyTimeout, logpath)
	},
}

func init() {
	rootCmd.AddCommand(onCmd)
}

