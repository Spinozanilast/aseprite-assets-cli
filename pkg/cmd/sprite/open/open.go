package open

import (
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"
	"github.com/spinozanilast/aseprite-assets-cli/internal/cmd/manager"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/aseprite"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/consts"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/environment"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/utils/files"
)

func NewSpriteOpenCmd(env *environment.Environment) *cobra.Command {

	var openAction manager.Action
	openAction = func(spriteFilename string, asePath string) error {
		if err := files.OpenFileWith(spriteFilename, asePath); err != nil {
			return err
		}

		return nil
	}

	managerParams := manager.Params{
		Env:             env,
		AssetAction:     openAction,
		ValidExtensions: aseprite.SpritesExtensions(),
		AssetsType:      consts.Sprite,
	}

	managerCmd := manager.NewAssetManagerCommand(managerParams)

	cmd := &cobra.Command{
		Use:     "open [arg]",
		Aliases: []string{"o"},
		Short:   "Open sprite with help of autocompletion",
		Long: heredoc.Doc(`
Open sprite with autocompletion of sprites from in-config sprites directories (recursive).`),
		RunE:              managerCmd.RunE,
		ValidArgsFunction: managerCmd.AutocompletionFunc,
	}

	return cmd
}
