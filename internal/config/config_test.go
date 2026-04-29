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
	if cfg.UI.LogoEnabled {
		t.Errorf("ui.logo_enabled default should be false, got true")
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

func TestSave_RoundTripsLogoToggle(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load initial: %v", err)
	}
	cfg.UI.LogoEnabled = true
	if err := Save(cfg); err != nil {
		t.Fatalf("Save: %v", err)
	}

	reloaded, err := Load()
	if err != nil {
		t.Fatalf("Load after save: %v", err)
	}
	if !reloaded.UI.LogoEnabled {
		t.Errorf("ui.logo_enabled did not persist; got false")
	}
	if !slices.Contains(reloaded.Worktree.Copy, ".env") {
		t.Errorf("worktree.copy lost on save round-trip; got %v", reloaded.Worktree.Copy)
	}
}
