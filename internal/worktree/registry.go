package worktree

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/foreverfl/gitt/internal/daemon"
	"github.com/foreverfl/gitt/internal/paths"
)

// Register tells the running daemon about a worktree. Returns an error if the
// daemon is unreachable or rejects the request.
func Register(mainRoot, branch, target string) error {
	sockpath, err := paths.SockPath()
	if err != nil {
		return err
	}
	response, err := daemon.Call(sockpath, daemon.Request{
		Op: daemon.OpRegisterWorktree,
		Args: map[string]any{
			"repo_root":        mainRoot,
			"repo_name":        filepath.Base(mainRoot),
			"branch_name":      branch,
			"safe_branch_name": SafeBranch(branch),
			"worktree_path":    target,
		},
	})
	if err != nil {
		return err
	}
	if !response.OK {
		return fmt.Errorf("%s", response.Error)
	}
	return nil
}

// TryRegister is the best-effort variant: if the daemon isn't running, it
// silently returns nil. Used by bootstrap commands like `gitt clone` that
// must work before the user has ever invoked `gitt on`.
func TryRegister(mainRoot, branch, target string) error {
	sockpath, err := paths.SockPath()
	if err != nil {
		return err
	}
	if err := daemon.Ping(sockpath); err != nil {
		if errors.Is(err, daemon.ErrNotRunning) {
			return nil
		}
		return err
	}
	return Register(mainRoot, branch, target)
}