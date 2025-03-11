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

func NewSpriteRemoveCmd(env *environment.Environment) *cobra.Command {
	var forceRemove bool

	var removeAction manager.Action
	removeAction = func(spriteFilename string, asePath string) error {
		if forceRemove {
			if err := files.RemoveFile(spriteFilename); err != nil {
				return err
			}
		}

		utils.PrintlnSuccess(fmt.Sprintf("🌌 Sprite was successfuly removed: %s", spriteFilename))
		return nil
	}

	managerParams := manager.Params{
		Env:                env,
		AssetAction:        removeAction,
		ValidExtensions:    aseprite.SpritesExtensions(),
		AssetsType:         consts.Sprite,
		ConfirmationNeeded: true,
	}

	managerCmd := manager.NewAssetManagerCommand(managerParams)

	cmd := &cobra.Command{
		Use:     "remove [arg]",
		Aliases: []string{"r", "rm"},
		Short:   "Remove sprite with help of autocompletion",
		Long: heredoc.Doc(`
Remove sprite with help of autocompletion of sprites from in-config sprites directories (recursive).`),
		RunE:              managerCmd.RunE,
		ValidArgsFunction: managerCmd.AutocompletionFunc,
	}

	cmd.Flags().BoolVarP(&forceRemove, "force", "f", false, "Force removal of autocompletion")

	return cmd
}
