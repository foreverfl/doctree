package compose

// ProjectName mirrors the original worktree-build naming:
// "<repo>_<worktree-dir>" so each worktree gets its own compose namespace.
func ProjectName(repoName, worktreeDir string) string {
	return repoName + "_" + worktreeDir
}

// File returns the conventional compose file path for a repo.
// TODO: make configurable via flag / env when more layouts appear.
func File(repoRoot string) string {
	return repoRoot + "/infra/docker/compose.local.yml"
}
