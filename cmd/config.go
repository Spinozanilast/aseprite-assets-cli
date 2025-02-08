package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	config "github.com/spinozanilast/aseprite-assets-cli/config"
	configTui "github.com/spinozanilast/aseprite-assets-cli/tui/config"
	utils "github.com/spinozanilast/aseprite-assets-cli/util"
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
	rootCmd.AddCommand(configCmd)
	configCmd.Flags().StringP("app-path", "a", "", "app path for opening config file")
}

func StartConfigInitializationTui(config *config.Config) {
	p := tea.NewProgram(configTui.InitialModel(config))
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Oopsy, something went wrong and cli is not executing properly: '%s'", err)
		os.Exit(1)
	}
}
