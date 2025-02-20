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
	scriptsDir, _ := cmd.Flags().GetString("scripts-dir")
	openAiApiKey, _ := cmd.Flags().GetString("openai-key")
	openAiApiUrl, _ := cmd.Flags().GetString("openai-url")

	if scriptsDir != "" {
		if err := setScriptsDir(scriptsDir); err != nil {
			return err
		}
	}

	if openAiApiKey != "" || openAiApiUrl != "" {
		config.SetOpenAiConfig(openAiApiKey, openAiApiUrl)
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

	return config.SetScriptDirPath(scriptsDir)
}

func handleOpenConfig(cmd *cobra.Command) error {
	appPath, _ := cmd.Flags().GetString("app-path")
	configFile := viper.ConfigFileUsed()

	if appPath != "" {
		return utils.OpenFileWith(configFile, appPath)
	}
	return utils.OpenFile(configFile)
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.Flags().StringP("app-path", "a", "", "App path for opening config file")
	configCmd.Flags().StringP("scripts-dir", "s", "", "Scripts directory for opening config file")
	configCmd.Flags().StringP("openai-key", "k", "", "OpenAI API key for generating palettes")
	configCmd.Flags().StringP("openai-url", "u", "", "OpenAI API URL for generating palettes")
}

func StartConfigInitializationTui(cfg *config.Config) {
	p := tea.NewProgram(configTui.InitialModel(cfg))
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Oopsy, something went wrong and cli is not executing properly: '%s'", err)
		os.Exit(1)
	}
}
