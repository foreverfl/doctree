package cmd

import (
	"github.com/foreverfl/doctree/internal/daemon"
	"github.com/foreverfl/doctree/internal/paths"
	"github.com/spf13/cobra"
)

// daemonRunCmd is the in-process entrypoint that `doctree on` fork-execs into.
// Hidden from --help: end users don't run this directly.
var daemonRunCmd = &cobra.Command{
	Use:    "daemon-run",
	Short:  "Internal: run the doctree daemon in foreground",
	Hidden: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		sockpath, err := paths.SockPath()
		if err != nil {
			return err
		}
		dbpath, err := paths.DBPath()
		if err != nil {
			return err
		}
		return daemon.Run(sockpath, dbpath)
	},
}

func init() {
	rootCmd.AddCommand(daemonRunCmd)
}