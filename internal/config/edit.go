package config

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// pickEditor resolves which editor to spawn for `gitt config`. Order
// follows the POSIX/UNIX convention: $VISUAL takes precedence over
// $EDITOR, with `vi` as the universal fallback. The returned slice is
// the argv to exec — splitting on whitespace is enough for common cases
// like `code --wait`; users with stranger commands should point the env
// var at a wrapper script.
func pickEditor() []string {
	for _, env := range []string{"VISUAL", "EDITOR"} {
		if value := strings.TrimSpace(os.Getenv(env)); value != "" {
			return strings.Fields(value)
		}
	}
	return []string{"vi"}
}

// OpenInEditor blocks until the user's editor exits, wiring the editor's
// stdio straight to the current terminal so it can run interactively.
func OpenInEditor(ctx context.Context, path string) error {
	argv := pickEditor()
	argv = append(argv, path)

	cmd := exec.CommandContext(ctx, argv[0], argv[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("editor %s: %w", argv[0], err)
	}
	return nil
}
