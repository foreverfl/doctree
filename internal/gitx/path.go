package gitx

import (
	"path/filepath"
	"strings"
)

// SafeBranch turns a git branch name into a single path segment by replacing
// directory separators with dashes. e.g. "feature/foo" -> "feature-foo".
func SafeBranch(branch string) string {
	return strings.NewReplacer("/", "-", "\\", "-").Replace(branch)
}

// WorktreePath returns the directory where the worktree for branch should
// live, following gitt's layout: <mainRoot>/.worktrees/<safe-branch>.
// mainRoot must be the main repository's top-level directory (see MainRepoRoot).
func WorktreePath(mainRoot, branch string) string {
	return filepath.Join(mainRoot, ".worktrees", SafeBranch(branch))
}
