package open

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/environment"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/utils/files"
)

func NewConfigOpenCmd(env *environment.Environment) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "open",
		Short: "Open configuration file in editor or via app",
		Example: `  aseprite-assets config open
  aseprite-assets config open --app-path /usr/bin/code`,
		RunE: func(cmd *cobra.Command, args []string) error {
			appPath, _ := cmd.Flags().GetString("app-path")

			if appPath != "" {
				if err := files.OpenFileWith(viper.ConfigFileUsed(), appPath); err != nil {
					return fmt.Errorf("failed to open with %s: %w", appPath, err)
				}
				return nil
			}

			if err := files.OpenFile(viper.ConfigFileUsed()); err != nil {
				return fmt.Errorf("failed to open config file: %w", err)
			}
			return nil
		},
	}

	cmd.Flags().StringP("app-path", "a", "",
		"Specify application to open config file")

	return cmd
}
