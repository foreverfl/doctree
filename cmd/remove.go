package cmd

import (
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Stop containers and remove the current worktree",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireDaemon(); err != nil {
			return err
		}
		// TODO: docker compose -p <project> down
		// TODO: daemon RPC: release ports + unregister
		// TODO: git worktree remove <pwd> (and rm -rf folder)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
