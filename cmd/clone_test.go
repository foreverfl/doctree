package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestCloneCmd(t *testing.T) {
	// Redirect HOME so paths.SockPath / RuntimeDir can't touch the dev machine's
	// real ~/.gitt (and the running daemon, if any).
	t.Setenv("HOME", t.TempDir())

	// Fixture: source repo on a non-default branch name ("trunk") so we verify
	// the default-branch detection actually reads HEAD instead of assuming "main".
	source := t.TempDir()
	runGit := func(dir string, args ...string) {
		t.Helper()
		c := exec.Command("git", args...)
		c.Dir = dir
		c.Env = append(os.Environ(),
			"GIT_AUTHOR_NAME=test",
			"GIT_AUTHOR_EMAIL=test@example.com",
			"GIT_COMMITTER_NAME=test",
			"GIT_COMMITTER_EMAIL=test@example.com",
		)
		if out, err := c.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
	runGit(source, "init", "-q", "-b", "trunk")
	runGit(source, "commit", "--allow-empty", "-q", "-m", "init")

	work := t.TempDir()
	t.Chdir(work)

	if err := cloneCmd.RunE(cloneCmd, []string{source, "myproj"}); err != nil {
		t.Fatalf("clone: %v", err)
	}

	project := filepath.Join(work, "myproj")

	if fi, err := os.Stat(filepath.Join(project, ".bare")); err != nil || !fi.IsDir() {
		t.Errorf("expected .bare directory, got err=%v isDir=%v", err, fi != nil && fi.IsDir())
	}

	pointer, err := os.ReadFile(filepath.Join(project, ".git"))
	if err != nil {
		t.Fatalf("read .git pointer: %v", err)
	}
	if got := strings.TrimSpace(string(pointer)); got != "gitdir: ./.bare" {
		t.Errorf(".git pointer = %q, want %q", got, "gitdir: ./.bare")
	}

	if fi, err := os.Stat(filepath.Join(project, ".worktrees", "trunk")); err != nil || !fi.IsDir() {
		t.Errorf("expected .worktrees/trunk directory, got err=%v", err)
	}
}