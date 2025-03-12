package palette

import (
	"github.com/spf13/cobra"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/cmd/palette/create"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/cmd/palette/lospec"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/cmd/palette/remove"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/environment"
)

func NewPaletteCmd(env *environment.Environment) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "palette [command]",
		Aliases: []string{"p"},
		Short:   "Palette command to manage palettes",
		Long: `
Subcommands allow you to:
- Create sprite (create)`,
	}

	cmd.AddCommand(create.NewPaletteCreateCmd(env))
	cmd.AddCommand(remove.NewPaletteRemoveCmd(env))
	cmd.AddCommand(lospec.NewPaletteLospecCmd(env))

	return cmd
}
