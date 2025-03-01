package scripts

import (
	"github.com/spf13/cobra"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/environment"
)

func NewScriptsCmd(env *environment.Environment) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "scripts",
		Aliases: []string{"sp"},
		Short:   "List existing included scripts for aseprite cli",
		Long: `
	1) CreateAsset command creates a new asset with the specified options.
	2) SavePaletteAsPreset command with help of OpenAI API creates a new palette with the specified options and saves it as a preset or to a file (png or gpl).`,
	}

	return cmd
}
