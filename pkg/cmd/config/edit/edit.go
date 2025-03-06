package edit

import (
	"fmt"
	"os"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/config"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/environment"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/tui"
)

func NewConfigEditCmd(env *environment.Environment) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edit",
		Short: "Modify configuration settings",
		Long:  `Configure settings through flags or interactive TUI. Without flags, launches interactive configuration interface.`,
		Example: heredoc.Doc(`aseprite-assets config edit --scripts-dir ./scripts
aseprite-assets config edit`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if scriptsDir, _ := cmd.Flags().GetString("scripts-dir"); scriptsDir != "" {
				if scriptsDir == "default" {
					err := config.SetDefaultScriptDirPath()
					if err != nil {
						return err
					}
					return nil
				}

				if info, err := os.Stat(scriptsDir); err != nil || !info.IsDir() {
					return fmt.Errorf("invalid scripts directory: %w", err)
				}

				if err := config.SetScriptDirPath(scriptsDir); err != nil {
					return fmt.Errorf("failed to set scripts directory: %w", err)
				}
				return nil
			}

			cfg, err := config.LoadConfig()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			if err = tui.StartConfigTui(cfg); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringP("scripts-dir", "s", "",
		"Set custom scripts directory path")

	return cmd
}
