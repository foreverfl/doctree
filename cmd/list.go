package cmd

import (
	"errors"
	"fmt"
	"text/tabwriter"

	"github.com/foreverfl/gitt/internal/daemon"
	"github.com/foreverfl/gitt/internal/paths"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all gitt-managed worktrees from the daemon's database",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		sockpath, err := paths.SockPath()
		if err != nil {
			return err
		}
		response, err := daemon.Call(sockpath, daemon.Request{Op: daemon.OpListWorktrees})
		if err != nil {
			if errors.Is(err, daemon.ErrNotRunning) {
				return fmt.Errorf("gitt daemon is not running. start it first: gitt on")
			}
			return err
		}
		if !response.OK {
			return fmt.Errorf("list worktrees failed: %s", response.Error)
		}

		// Daemon returns []store.Worktree as JSON; the decoder lands it as
		// []any of map[string]any, so we walk it generically.
		raw, _ := response.Data["worktrees"].([]any)
		if len(raw) == 0 {
			fmt.Println("(no worktrees registered)")
			return nil
		}

		writer := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
		fmt.Fprintln(writer, "REPO\tBRANCH\tSTATUS\tPATH")
		for _, item := range raw {
			row, ok := item.(map[string]any)
			if !ok {
				continue
			}
			fmt.Fprintf(writer, "%s\t%s\t%s\t%s\n",
				stringField(row, "repo_name"),
				stringField(row, "branch_name"),
				stringField(row, "status"),
				stringField(row, "worktree_path"),
			)
		}
		return writer.Flush()
	},
}

func stringField(row map[string]any, key string) string {
	value, _ := row[key].(string)
	return value
}

func init() {
	rootCmd.AddCommand(listCmd)
}
