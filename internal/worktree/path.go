// Package worktree owns gitt's per-branch worktree layout convention:
// <repo-parent>/.worktrees/<repo-name>/<safe-branch>.
package worktree

import (
	"path/filepath"
	"strings"
)

// SafeBranch turns a git branch name into a single path segment by replacing
// directory separators with dashes. e.g. "feature/foo" -> "feature-foo".
func SafeBranch(branch string) string {
	return strings.NewReplacer("/", "-", "\\", "-").Replace(branch)
}

// Path returns the directory where the worktree for branch should live,
// following gitt's layout: <repo-parent>/.worktrees/<repo-name>/<safe-branch>.
func Path(repoRoot, branch string) string {
	return filepath.Join(
		filepath.Dir(repoRoot),
		".worktrees",
		filepath.Base(repoRoot),
		SafeBranch(branch),
	)
}
