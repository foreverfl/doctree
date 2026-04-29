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
	Short: "Show the logo and toggle whether `gitt on` prints it on startup",
	Long: "Prints the gitt logo art and asks whether to enable it on\n" +
		"`gitt on` startup. The choice is persisted to ~/.gitt/config.toml\n" +
		"as `[ui] logo_enabled`. Default is disabled.",
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
