package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/foreverfl/gitt/internal/daemon"
	"github.com/foreverfl/gitt/internal/paths"
	"github.com/foreverfl/gitt/internal/process"
	"github.com/foreverfl/gitt/internal/prompt"
	"github.com/foreverfl/gitt/internal/release"
	"github.com/foreverfl/gitt/internal/version"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update gitt to the latest release",
	RunE: func(cmd *cobra.Command, args []string) error {
		current := version.Installed()

		fmt.Println("checking latest release...")
		latest, err := release.LatestTag()
		if err != nil {
			return err
		}

		if current != "" && current == latest {
			fmt.Printf("already at latest (%s)\n", current)
			return nil
		}
		if current == "" {
			fmt.Printf("updating to %s\n", latest)
		} else {
			fmt.Printf("updating %s -> %s\n", current, latest)
		}

		runtime, err := paths.RuntimeDir()
		if err != nil {
			return err
		}
		sockpath, err := paths.SockPath()
		if err != nil {
			return err
		}
		pidpath, err := paths.PidPath()
		if err != nil {
			return err
		}
		logpath, err := paths.LogPath()
		if err != nil {
			return err
		}

		daemonRunning := false
		if pid, ok := process.ReadPid(pidpath); ok && process.Alive(pid) {
			if err := daemon.Ping(sockpath); err == nil {
				daemonRunning = true
			}
		}

		fmt.Println("update will reset gitt runtime data:")
		fmt.Printf("  - %s (db, registered worktrees, sock, pid, log)\n", runtime)
		if daemonRunning {
			fmt.Println("  - the running daemon will be stopped and restarted")
		}
		fmt.Println()

		yes, _ := cmd.Flags().GetBool("yes")
		if !yes {
			ok, err := prompt.Confirm("proceed?", false)
			if err != nil {
				if errors.Is(err, prompt.ErrNoTTY) {
					return fmt.Errorf("non-interactive shell — pass --yes to confirm")
				}
				return err
			}
			if !ok {
				fmt.Println("aborted.")
				return nil
			}
		}

		selfPath, err := os.Executable()
		if err != nil {
			return fmt.Errorf("locate self: %w", err)
		}
		newPath := selfPath + ".new"

		fmt.Println("downloading...")
		if err := release.Download(latest, newPath); err != nil {
			_ = os.Remove(newPath)
			return err
		}

		if daemonRunning {
			if err := daemon.Shutdown(sockpath, pidpath, os.Stdout, os.Stderr); err != nil {
				_ = os.Remove(newPath)
				return fmt.Errorf("stop daemon: %w", err)
			}
		}

		if err := os.RemoveAll(runtime); err != nil {
			_ = os.Remove(newPath)
			return fmt.Errorf("reset %s: %w", runtime, err)
		}

		if err := os.Rename(newPath, selfPath); err != nil {
			_ = os.Remove(newPath)
			return fmt.Errorf("replace binary: %w", err)
		}

		if vpath, verr := paths.VersionPath(); verr == nil {
			_ = os.WriteFile(vpath, []byte(latest+"\n"), 0o644)
		}

		fmt.Printf("updated to %s\n", latest)

		if daemonRunning {
			fmt.Println("restarting daemon...")
			pid, err := daemon.Spawn(selfPath, sockpath, pidpath, logpath)
			if err != nil {
				return fmt.Errorf("restart daemon: %w", err)
			}
			fmt.Printf("gitt daemon started (pid=%d)\n", pid)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
