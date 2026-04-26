package cmd

import (
	"errors"
	"fmt"

	"github.com/foreverfl/gitt/internal/daemon"
	"github.com/foreverfl/gitt/internal/paths"
	"github.com/spf13/cobra"
)

var sqliteCmd = &cobra.Command{
	Use:   "sqlite",
	Short: "Run a SQLite self-test against the daemon's database",
	Long: "Asks the running gitt daemon to create a scratch table, insert a\n" +
		"row, read it back, and drop the table. Useful for confirming the\n" +
		"daemon's database connection is healthy.\n\n" +
		"Requires `gitt on` to be running.",
	RunE: func(cmd *cobra.Command, args []string) error {
		sockpath, err := paths.SockPath()
		if err != nil {
			return err
		}
		response, err := daemon.Call(sockpath, daemon.Request{Op: daemon.OpSqliteTest})
		if err != nil {
			if errors.Is(err, daemon.ErrNotRunning) {
				return fmt.Errorf("gitt daemon not running. start it with `gitt on`")
			}
			return err
		}
		if !response.OK {
			return fmt.Errorf("sqlite test failed: %s", response.Error)
		}
		message, _ := response.Data["message"].(string)
		fmt.Println(message)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(sqliteCmd)
}
