// Package config loads gitt's user configuration from ~/.gitt/config.toml.
// A missing file is not an error: callers transparently get the built-in
// defaults embedded at build time. The first time the user runs
// `gitt config`, EnsureFile materialises those defaults to disk so they
// can be edited.
package config

import (
	_ "embed"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"

	"github.com/foreverfl/gitt/internal/paths"
)

//go:embed default.toml
var defaultTOML []byte

// Config mirrors config.toml. Sections are flat for now; add more as
// features land.
type Config struct {
	Worktree WorktreeSection `toml:"worktree"`
}

type WorktreeSection struct {
	Copy    []string `toml:"copy"`
	Symlink []string `toml:"symlink"`
	Ignore  []string `toml:"ignore"`
}

// Load returns the resolved config. If ~/.gitt/config.toml exists it is
// decoded and returned; otherwise the embedded defaults are returned.
// Callers should not assume which source they received — both shapes are
// fully populated.
func Load() (*Config, error) {
	path, err := paths.ConfigPath()
	if err != nil {
		return nil, err
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return decode(defaultTOML)
		}
		return nil, fmt.Errorf("read config %s: %w", path, err)
	}
	return decode(raw)
}

// EnsureFile creates ~/.gitt/config.toml from the embedded defaults if it
// does not yet exist, and returns the path either way. Used by
// `gitt config` so the editor always has something to open.
func EnsureFile() (string, error) {
	path, err := paths.ConfigPath()
	if err != nil {
		return "", err
	}
	if _, err := os.Stat(path); err == nil {
		return path, nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return "", fmt.Errorf("stat config %s: %w", path, err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return "", fmt.Errorf("mkdir config dir: %w", err)
	}
	if err := os.WriteFile(path, defaultTOML, 0o644); err != nil {
		return "", fmt.Errorf("write config %s: %w", path, err)
	}
	return path, nil
}

func decode(raw []byte) (*Config, error) {
	var cfg Config
	if err := toml.Unmarshal(raw, &cfg); err != nil {
		return nil, fmt.Errorf("decode config: %w", err)
	}
	return &cfg, nil
}
