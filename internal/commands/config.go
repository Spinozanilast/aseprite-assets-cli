package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/utils"

	tea "github.com/charmbracelet/bubbletea"
	configTui "github.com/spinozanilast/aseprite-assets-cli/internal/tui/config"
	config "github.com/spinozanilast/aseprite-assets-cli/pkg/config"
)

var configCmd = &cobra.Command{
	Use:       "config [ARG]",
	Aliases:   []string{"cfg"},
	Short:     "Manage aseprite-assets-cli configuration",
	Args:      cobra.ExactArgs(1),
	ValidArgs: []string{"info", "edit", "open"},
	RunE: func(cmd *cobra.Command, args []string) error {
		switch args[0] {
		case "info":
			config.ConfigInfo()
		case "edit":
			scriptsDir, err := cmd.Flags().GetString("scripts-dir")
			if err == nil && scriptsDir != "" {
				if scriptsDir == "default" {
					config.SetDefaultScriptDirPath()
					return nil
				}
				info, err := os.Stat(scriptsDir)
				if info.IsDir() || err != nil {
					err = config.SetScriptDirPath(scriptsDir)
					if err != nil {
						return fmt.Errorf("failed to set scripts directory path: %w", err)
					}
					return nil
				}
				return fmt.Errorf("scripts directory path is not a directory")
			}

			config, err := config.LoadConfig()
			if err != nil {
				return err
			}

			if config.ScriptDirPath == "" {
				StartConfigInitializationTui(config)
			}

		case "open":
			appPath, err := cmd.Flags().GetString("app-path")
			if err == nil && appPath != "" {
				fmt.Println(viper.ConfigFileUsed())
				fmt.Println(appPath)
				err = utils.OpenFileWith(viper.ConfigFileUsed(), appPath)
			} else {
				err = utils.OpenFile(viper.ConfigFileUsed())
			}

			if err != nil {
				return err
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.Flags().StringP("app-path", "a", "", "app path for opening config file")
	configCmd.Flags().StringP("scripts-dir", "s", "", "scripts directory for opening config file")
}

func StartConfigInitializationTui(config *config.Config) {
	p := tea.NewProgram(configTui.InitialModel(config))
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Oopsy, something went wrong and cli is not executing properly: '%s'", err)
		os.Exit(1)
	}
}
