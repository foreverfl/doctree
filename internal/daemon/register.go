package daemon

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/foreverfl/gitt/internal/paths"
	"github.com/foreverfl/gitt/internal/worktree"
)

// RegisterWorktree tells the running daemon about a worktree. Returns an error
// if the daemon is unreachable or rejects the request.
func RegisterWorktree(mainRoot, branch, target string) error {
	sockpath, err := paths.SockPath()
	if err != nil {
		return err
	}
	response, err := Call(sockpath, Request{
		Op: OpRegisterWorktree,
		Args: map[string]any{
			"repo_root":        mainRoot,
			"repo_name":        filepath.Base(mainRoot),
			"branch_name":      branch,
			"safe_branch_name": worktree.SafeBranch(branch),
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

// TryRegisterWorktree is the best-effort variant: if the daemon isn't running,
// it silently returns nil. Used by bootstrap commands like `gitt clone` that
// must work before the user has ever invoked `gitt on`.
func TryRegisterWorktree(mainRoot, branch, target string) error {
	sockpath, err := paths.SockPath()
	if err != nil {
		return err
	}
	if err := Ping(sockpath); err != nil {
		if errors.Is(err, ErrNotRunning) {
			return nil
		}
		return err
	}
	return RegisterWorktree(mainRoot, branch, target)
}

// ReleaseWorktree tells the daemon to drop the worktree row identified by
// (mainRoot, branch). cmd/remove calls this after `git worktree remove`
// succeeds so the daemon's view stays in sync with the filesystem.
func ReleaseWorktree(mainRoot, branch string) error {
	sockpath, err := paths.SockPath()
	if err != nil {
		return err
	}
	response, err := Call(sockpath, Request{
		Op: OpRelease,
		Args: map[string]any{
			"repo_root":   mainRoot,
			"branch_name": branch,
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