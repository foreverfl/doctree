package paths

import "path/filepath"

// ConfigPath returns ~/.gitt/config.toml. The file is created on demand
// by config.EnsureFile (typically the first time `gitt config` runs);
// callers that only need to read it should treat a missing file as
// "use built-in defaults".
func ConfigPath() (string, error) {
	dir, err := RuntimeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.toml"), nil
}
