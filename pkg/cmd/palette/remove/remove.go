package remove

import (
	"fmt"
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"
	"github.com/spinozanilast/aseprite-assets-cli/internal/cmd/manager"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/aseprite"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/consts"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/environment"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/utils"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/utils/files"
)

func NewPaletteRemoveCmd(env *environment.Environment) *cobra.Command {
	var forceRemove bool

	var removeAction manager.Action
	removeAction = func(paletteFilename string, asePath string) error {
		if forceRemove {
			if err := files.RemoveFile(paletteFilename); err != nil {
				return err
			}
		}

		utils.PrintlnSuccess(fmt.Sprintf("🌌 Palette was successfuly removed: %s", paletteFilename))
		return nil
	}

	managerParams := manager.Params{
		Env:                env,
		AssetAction:        removeAction,
		ValidExtensions:    aseprite.PaletteExtensions(),
		AssetsType:         consts.Palette,
		ConfirmationNeeded: true,
	}

	managerCmd := manager.NewAssetManagerCommand(managerParams)

	cmd := &cobra.Command{
		Use:     "remove [arg]",
		Aliases: []string{"r", "rm"},
		Short:   "Remove palette with help of autocompletion",
		Long: heredoc.Doc(`
Remove palette with help of autocompletion of sprites from in-config sprites directories (recursive).`),
		RunE:              managerCmd.RunE,
		ValidArgsFunction: managerCmd.AutocompletionFunc,
	}

	cmd.Flags().BoolVarP(&forceRemove, "force", "f", false, "Force removal of autocompletion or input at all")

	return cmd
}
