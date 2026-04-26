package cmd

import (
	"fmt"
	"os"

	"github.com/foreverfl/doctree/internal/gitx"
	"github.com/foreverfl/doctree/internal/worktree"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove <branch>",
	Short: "Remove the git worktree for <branch>",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireDaemon(); err != nil {
			return err
		}
		branch := args[0]

		repoRoot, err := gitx.RepoRoot()
		if err != nil {
			return err
		}
		target := worktree.Path(repoRoot, branch)

		if _, err := os.Stat(target); err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("no worktree at %s", target)
			}
			return fmt.Errorf("stat worktree: %w", err)
		}

		if err := gitx.WorktreeRemove(target); err != nil {
			fmt.Fprintln(os.Stderr, "tip: if the worktree has uncommitted or untracked changes, commit or stash them first.")
			return err
		}

		fmt.Printf("removed worktree\n  path:   %s\n  branch: %s\n", target, branch)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
