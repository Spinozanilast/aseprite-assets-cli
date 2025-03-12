package manager

import (
	"errors"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	autocomp "github.com/spinozanilast/aseprite-assets-cli/internal/cmd"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/config"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/consts"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/environment"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/utils/files"
)

type Params struct {
	Env                *environment.Environment
	AssetAction        Action
	ValidExtensions    []string
	AssetsType         consts.AssetsType
	ConfirmationNeeded bool
}

type Command struct {
	RunE               func(cmd *cobra.Command, args []string) error
	AutocompletionFunc func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective)
}

type Action func(filename string, asepritePath string) error

func NewAssetManagerCommand(params Params) *Command {
	cmd := &Command{}
	validExtensions := params.ValidExtensions

	cmd.RunE = func(_ *cobra.Command, args []string) error {
		cfg, err := params.Env.Config()
		if err != nil {
			return err
		}

		if len(cfg.AsepritePath) == 0 {
			return errors.New("no aseprite path specified")
		}

		if len(args) == 0 {
			return errors.New("no arg given")
		}

		if params.ConfirmationNeeded && !ConfirmActions() {
			return nil
		}

		for _, arg := range args {
			if files.CheckFileExists(arg, false) && files.CheckFileExtension(arg, validExtensions...) {
				if err := params.AssetAction(arg, cfg.AsepritePath); err != nil {
					return err
				}
			} else {
				return fmt.Errorf("invalid arg: %s", arg)
			}
		}

		return nil
	}

	cmd.AutocompletionFunc = func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		cfg, err := params.Env.Config()

		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		autocomplete := autocomp.GenerateFilesAutoCompletions(chooseAssetsDirs(params.AssetsType, cfg), validExtensions)
		return autocomplete(cmd, args, toComplete)
	}

	return cmd
}

func ConfirmActions() bool {
	question := &survey.Confirm{
		Message: "Are you sure you want to proceed?",
		Default: false,
	}

	var confirm bool

	if err := survey.AskOne(question, &confirm); err != nil {
		return false
	} else {
		return confirm
	}
}

func chooseAssetsDirs(assetsType consts.AssetsType, cfg *config.Config) []string {
	switch assetsType {
	case consts.Sprite:
		return cfg.SpriteFolderPaths
	case consts.Palette:
		return cfg.PalettesFolderPaths
	default:
		return nil
	}
}
