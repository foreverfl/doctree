package paths

import (
	"os"
	"strings"
)

// InstalledVersion reads ~/.gitt/VERSION (written by install.sh) and returns
// the version string, or "" when the file is missing or unreadable.
func InstalledVersion() string {
	path, err := VersionPath()
	if err != nil {
		return ""
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}
