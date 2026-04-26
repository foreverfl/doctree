package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/foreverfl/doctree/internal/paths"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the installed doctree version",
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := paths.VersionPath()
		if err != nil {
			return err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Println("unknown (not installed via install.sh)")
				return nil
			}
			return err
		}
		fmt.Println(strings.TrimSpace(string(data)))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
