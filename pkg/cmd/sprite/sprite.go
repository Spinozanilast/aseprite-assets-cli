package sprite

import (
	"github.com/spf13/cobra"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/cmd/sprite/create"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/environment"
)

func NewSpriteCmd(env *environment.Environment) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "sprite [command]",
		Aliases: []string{"s"},
		Short:   "Sprite commands to manage sprites",
		Long: `
Subcommands allow you to:
- Create sprite (create)`,
	}

	cmd.AddCommand(create.NewSpriteCreateCmd(env))

	return cmd
}
