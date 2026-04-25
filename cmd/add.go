package cmd

import (
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add <branch>",
	Short: "Create a new git worktree for <branch>",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireDaemon(); err != nil {
			return err
		}
		_ = args[0] // branch
		// TODO: daemon RPC: register worktree + allocate ports
		// TODO: git worktree add <repo>/../.worktrees/<repo>/<safe-branch> -b <branch>
		// TODO: write infra/docker/.env.worktree with allocated ports
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
