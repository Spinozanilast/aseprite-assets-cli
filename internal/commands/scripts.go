package commands

import (
	"github.com/spf13/cobra"
)

var scriptsCommand = &cobra.Command{
	Use:     "scripts",
	Aliases: []string{"s"},
	Short:   "List existing included scripts for aseprite cli",
	Long: `
	1) CreateAsset command creates a new asset with the specified options.
	2) SavePaletteAsPreset command with help of OpenAI API creates a new palette with the specified options and saves it as a preset or to a file (png or gpl).`,
}

func init() {
	rootCmd.AddCommand(scriptsCommand)
}
