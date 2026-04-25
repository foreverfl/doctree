package cmd

import (
	"github.com/spf13/cobra"
)

var offCmd = &cobra.Command{
	Use:   "off",
	Short: "Stop the aw daemon",
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: rpc shutdown to daemon (or signal pid)
		// TODO: clean up sock and pid file
		return nil
	},
}

func init() {
	rootCmd.AddCommand(offCmd)
}
