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

			StartConfigInitializationTui(config)
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
	desc := &CommandDescription{
		Title: "Configure aseprite-assets CLI",
		Short: "Manage configuration settings",
		Long:  "Configure various settings for the aseprite-assets CLI tool.",
		Sections: []Section{
			{
				Title: "Configuration Options",
				Text:  "You can configure the following settings:\n- Application path\n- Scripts directory\n- Assets directories\n- Palettes directories\n- OpenAI API settings",
			},
			{
				Title: "Interactive Configuration (TUI)",
				Text: "When running 'config edit' without parameters, an interactive Terminal User Interface will be launched allowing you to:\n" +
					"- Set Aseprite executable path\n" +
					"- Configure assets directories\n" +
					"- Configure palettes directories\n" +
					"- Set up OpenAI API configuration",
			},
			{
				Title: "Configuration File",
				Text:  fmt.Sprintf("Configuration is stored in JSON format at %s\n", viper.ConfigFileUsed()) + "You can directly edit this file using 'config open' command.",
			},
		},
		Examples: []string{
			"aseprite-assets config edit --app-path \"C:\\Program Files\\Aseprite\\Aseprite.exe\"",
			"aseprite-assets config edit --scripts-dir \"./scripts\"",
			"aseprite-assets config edit              # Opens interactive TUI",
			"aseprite-assets config open              # Opens config file in default editor",
			"aseprite-assets config info              # Shows current configuration",
		},
	}

	desc.ApplyToCommand(configCmd)

	configCmd.Flags().StringP("app-path", "a", "", "app path for opening config file")
	configCmd.Flags().StringP("scripts-dir", "s", "", "scripts directory for opening config file")
	rootCmd.AddCommand(configCmd)
}

func StartConfigInitializationTui(config *config.Config) {
	p := tea.NewProgram(configTui.InitialModel(config))
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Oopsy, something went wrong and cli is not executing properly: '%s'", err)
		os.Exit(1)
	}
}
