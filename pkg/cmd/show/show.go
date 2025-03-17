package show

import (
	"fmt"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/aseprite"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/aseprite/preview"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/config"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/environment"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/utils"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/utils/files"
)

type PreviewParams struct {
	Filename         string
	ColorFormat      string
	ColorsPerRow     int
	IsPalettePreview bool
}

func NewShowCmd(env *environment.Environment) *cobra.Command {
	params := &PreviewParams{}

	cmd := &cobra.Command{
		Use:     "show",
		Aliases: []string{"sh"},
		Example: heredoc.Doc(`
			# Example documentation
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := env.Config()
			if err != nil {
				return fmt.Errorf("config loading failed: %w", err)
			}

			if err := validateInput(params.Filename); err != nil {
				return err
			}

			generator := initializeGenerator(cfg)

			output, err := generator.Generate(preview.GenerateParams{
				Filename:         params.Filename,
				ColorFormat:      params.ColorFormat,
				ColorsPerRow:     params.ColorsPerRow,
				IsPalettePreview: params.IsPalettePreview,
			})

			return handleGenerationResult(output, err)
		},
	}

	cmd.Flags().StringVarP(&params.Filename, "Filename", "f", "", "asset Filename")
	cmd.Flags().StringVarP(&params.ColorFormat, "color-format", "c", "hex", "color format for palettes")
	cmd.Flags().IntVarP(&params.ColorsPerRow, "output-row-count", "r", 5, "colors per row for palettes")
	cmd.Flags().BoolVarP(&params.IsPalettePreview, "palette-preview", "p", false, "show palette preview")

	if err := cmd.MarkFlagRequired("Filename"); err != nil {
		return nil
	}

	return cmd
}

func validateInput(filename string) error {
	validExtension := files.CheckFileExtension(filename,
		append(aseprite.SpritesExtensions(), aseprite.AvailablePaletteExtensions()...)...,
	)

	if !files.CheckFileExists(filename, false) || !validExtension {
		return fmt.Errorf("invalid file: %s", filename)
	}
	return nil
}

func initializeGenerator(cfg *config.Config) *preview.Generator {
	aseCli := aseprite.NewCLI(cfg.AsepritePath, cfg.ScriptDirPath, cfg.FromSteam)
	return preview.NewGenerator(aseCli)
}

func handleGenerationResult(output string, err error) error {
	if err != nil {
		return err
	}

	if output == "" {
		utils.PrintFormatted("No preview available for this file\n")
		return nil
	}

	fmt.Println(output)
	return nil
}
