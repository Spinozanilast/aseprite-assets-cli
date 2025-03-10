package manager

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/config"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/consts"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/environment"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/utils/files"
)

type Params struct {
	Env             *environment.Environment
	AssetAction     Action
	ValidExtensions []string
	AssetsType      consts.AssetsType
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

	cmd.AutocompletionFunc = func(_ *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		cfg, err := params.Env.Config()
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		var completions []string
		for _, dir := range chooseAssetsDirs(params.AssetsType, cfg) {
			fs, err := files.FindFilesOfExtensionsRecursiveFlatten(dir, validExtensions...)
			if err == nil {
				completions = append(completions, fs...)
			}
		}
		return completions, cobra.ShellCompDirectiveNoFileComp
	}

	return cmd
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
