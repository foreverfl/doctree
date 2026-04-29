package cmd

import (
	"fmt"

	"github.com/foreverfl/gitt/internal/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Edit the gitt config file in $EDITOR",
	Long: "Opens ~/.gitt/config.toml in $VISUAL, $EDITOR, or vi.\n\n" +
		"On first run the file is created from the built-in defaults so\n" +
		"the editor always has something to open. Edits take effect on\n" +
		"the next command that reads config.",
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := config.EnsureFile()
		if err != nil {
			return fmt.Errorf("ensure config file: %w", err)
		}
		return config.OpenInEditor(cmd.Context(), path)
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
