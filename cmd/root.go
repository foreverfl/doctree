package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/foreverfl/agent-worktree/internal/daemon"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "aw",
	Short: "agent-worktree — git worktree + docker compose orchestrator",
	Long:  "Coordinates per-branch git worktrees and their docker compose stacks via a small SQLite-backed daemon.",
}

func Execute() error {
	return rootCmd.Execute()
}

// runtimeDir returns ~/.aw, creating it on demand.
func runtimeDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, ".aw")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	return dir, nil
}

func sockPath() (string, error) {
	dir, err := runtimeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "aw.sock"), nil
}

// requireDaemon errors out with an init hint when the daemon isn't reachable.
func requireDaemon() error {
	sp, err := sockPath()
	if err != nil {
		return err
	}
	if err := daemon.Ping(sp); err != nil {
		if errors.Is(err, daemon.ErrNotRunning) {
			return fmt.Errorf("aw daemon is not running. start it first: aw on")
		}
		return err
	}
	return nil
}
