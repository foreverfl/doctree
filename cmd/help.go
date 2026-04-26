package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var helpCmd = &cobra.Command{
	Use:   "help",
	Short: "Show doctree help",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("doctree help — placeholder")
	},
}

func init() {
	rootCmd.SetHelpCommand(helpCmd)
}
