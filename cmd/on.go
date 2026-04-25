package cmd

import (
	"github.com/spf13/cobra"
)

var onCmd = &cobra.Command{
	Use:   "on",
	Short: "Start the aw daemon",
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: check ~/.aw/aw.pid; if alive, noop with message
		// TODO: fork-exec self in daemon mode (or spawn detached)
		// TODO: write pid file, listen on ~/.aw/aw.sock, open ~/.aw/aw.db
		return nil
	},
}

func init() {
	rootCmd.AddCommand(onCmd)
}
