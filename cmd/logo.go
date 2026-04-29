package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/foreverfl/gitt/internal/config"
	"github.com/foreverfl/gitt/internal/ui"
	"github.com/spf13/cobra"
)

var logoCmd = &cobra.Command{
	Use:   "logo",
	Short: "Print the logo and toggle whether it appears on `gitt on` startup",
	Long: "Prints the gitt logo art, then prompts:\n\n" +
		"  Show this logo on `gitt on` startup? [y/N]\n\n" +
		"The default shown in brackets reflects the current value. Answering\n" +
		"yes or no persists the choice to ~/.gitt/config.toml as\n" +
		"`[ui] logo_enabled`. Logo display is disabled by default.\n\n" +
		"Requires an interactive terminal (stdin must be a TTY).",
	RunE: func(cmd *cobra.Command, args []string) error {
		ui.Logo(os.Stdout)

		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}

		enable, err := ui.Confirm("Show this logo on `gitt on` startup?", cfg.UI.LogoEnabled)
		if err != nil {
			if errors.Is(err, ui.ErrNoTTY) {
				return fmt.Errorf("`gitt logo` needs an interactive terminal")
			}
			return err
		}

		if enable == cfg.UI.LogoEnabled {
			fmt.Printf("logo_enabled stays %v\n", enable)
			return nil
		}
		cfg.UI.LogoEnabled = enable
		if err := config.Save(cfg); err != nil {
			return err
		}
		fmt.Printf("logo_enabled = %v\n", enable)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(logoCmd)
}
