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
	Use:     "config [args]",
	Aliases: []string{"cfg"},
	Short:   "Manage aseprite-assets-cli configuration",
	Long: `Manage aseprite-assets-cli configuration.

Available arguments:
  info  - Display the current configuration.
  edit  - Edit the configuration using a TUI.
  open  - Open the configuration file with the specified application.`,
	Args:      cobra.ExactArgs(1),
	ValidArgs: []string{"info", "edit", "open"},
	RunE:      runConfigCmd,
}

func runConfigCmd(cmd *cobra.Command, args []string) error {
	switch args[0] {
	case "info":
		config.ConfigInfo()
		return nil
	case "edit":
		return handleEditConfig(cmd)
	case "open":
		return handleOpenConfig(cmd)
	default:
		return fmt.Errorf("invalid argument: %s", args[0])
	}
}

func handleEditConfig(cmd *cobra.Command) error {
	scriptsDir, err := cmd.Flags().GetString("scripts-dir")
	if err != nil {
		return err
	}

	if scriptsDir != "" {
		return setScriptsDir(scriptsDir)
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}
	StartConfigInitializationTui(cfg)
	return nil
}

func setScriptsDir(scriptsDir string) error {
	if scriptsDir == "default" {
		config.SetDefaultScriptDirPath()
		return nil
	}

	info, err := os.Stat(scriptsDir)
	if err != nil || !info.IsDir() {
		return fmt.Errorf("scripts directory path is not a directory")
	}

	if err := config.SetScriptDirPath(scriptsDir); err != nil {
		return fmt.Errorf("failed to set scripts directory path: %w", err)
	}
	return nil
}

func handleOpenConfig(cmd *cobra.Command) error {
	appPath, err := cmd.Flags().GetString("app-path")
	if err != nil {
		return err
	}

	configFile := viper.ConfigFileUsed()
	if appPath != "" {
		fmt.Println(configFile)
		fmt.Println(appPath)
		return utils.OpenFileWith(configFile, appPath)
	}
	return utils.OpenFile(configFile)
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.Flags().StringP("app-path", "a", "", "app path for opening config file")
	configCmd.Flags().StringP("scripts-dir", "s", "", "scripts directory for opening config file")
}

func StartConfigInitializationTui(cfg *config.Config) {
	p := tea.NewProgram(configTui.InitialModel(cfg))
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Oopsy, something went wrong and cli is not executing properly: '%s'", err)
		os.Exit(1)
	}
}
