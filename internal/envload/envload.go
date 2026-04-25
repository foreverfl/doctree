package envload

// Load merges <composeDir>/.env.local then <composeDir>/.env.worktree (latter wins),
// returning the resulting key/value map. Both files are optional.
func Load(composeDir string) (map[string]string, error) {
	_ = composeDir
	// TODO: read .env.local (ignore if missing), parse KEY=VALUE
	// TODO: read .env.worktree (ignore if missing), overlay
	return map[string]string{}, nil
}

// WriteWorktreeEnv writes kv into <composeDir>/.env.worktree, replacing the file.
// This is what `aw add` calls after the daemon hands back allocated host ports.
func WriteWorktreeEnv(composeDir string, kv map[string]string) error {
	_ = composeDir
	_ = kv
	// TODO: deterministic key order, atomic write (tmp + rename)
	return nil
}
