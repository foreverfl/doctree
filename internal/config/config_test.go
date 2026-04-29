package config

import (
	"slices"
	"testing"
)

func TestLoad_FallsBackToEmbeddedDefaults(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if !slices.Contains(cfg.Worktree.Copy, ".env") {
		t.Errorf("worktree.copy missing .env sentinel; got %v", cfg.Worktree.Copy)
	}
	if len(cfg.Worktree.Symlink) == 0 {
		t.Errorf("worktree.symlink should not be empty in defaults")
	}
	if len(cfg.Worktree.Ignore) == 0 {
		t.Errorf("worktree.ignore should not be empty in defaults")
	}
}
