package open

import (
	"errors"
	"fmt"
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/aseprite"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/environment"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/utils"
)

func NewSpriteOpenCmd(env *environment.Environment) *cobra.Command {
	validExtensions := aseprite.SpritesExtensions()

	cmd := &cobra.Command{
		Use:     "open [arg]",
		Aliases: []string{"o"},
		Short:   "Open sprite with autocompletion",
		Long: heredoc.Doc(`
Open sprite with autocompletion of sprites from in-config sprites directories (recursive).`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := env.Config()
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
				if utils.СheckFileExists(arg, false) && utils.СheckFileExtension(arg, validExtensions...) {
					if err := utils.OpenFileWith(arg, cfg.AsepritePath); err != nil {
						return err
					}
				} else {
					return fmt.Errorf("invalid arg: %s", arg)
				}
			}

			return nil
		},
	}

	cmd.ValidArgsFunction = func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		cfg, err := env.Config()
		if err != nil {
			return []string{}, cobra.ShellCompDirectiveNoFileComp
		}

		var completions []string
		for _, dir := range cfg.AssetsFolderPaths {
			files, err := utils.FindFilesOfExtensionsRecursiveFlatten(dir, validExtensions...)
			if err == nil {
				completions = append(completions, files...)
			}
		}
		return completions, cobra.ShellCompDirectiveNoFileComp
	}

	return cmd
}
