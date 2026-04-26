package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var helpCmd = &cobra.Command{
	Use:   "help",
	Short: "Show gitt help",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("gitt help — placeholder")
	},
}

func init() {
	rootCmd.SetHelpCommand(helpCmd)
}
