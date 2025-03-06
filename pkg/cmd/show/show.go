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
)

type showCmd struct {
	config           *config.Config
	previewGenerator *preview.Generator
	filename         string
	colorFormat      string
	colorsPerRow     int
	isPalettePreview bool
}

func NewShowCmd(env *environment.Environment) *cobra.Command {
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

			params, err := parseCommandParams(cmd)
			if err != nil {
				return err
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

	cmd.Flags().StringP("filename", "f", "", "asset filename")
	cmd.Flags().StringP("color-format", "c", "hex", "color format for palettes")
	cmd.Flags().IntP("output-row-count", "r", 5, "colors per row for palettes")
	cmd.Flags().BoolP("palette-preview", "p", false, "show palette preview")
	cmd.MarkFlagRequired("filename")

	return cmd
}

func parseCommandParams(cmd *cobra.Command) (*struct {
	Filename         string
	ColorFormat      string
	ColorsPerRow     int
	IsPalettePreview bool
}, error) {
	filename, err := cmd.Flags().GetString("filename")
	if err != nil {
		return nil, fmt.Errorf("filename parameter error: %w", err)
	}

	colorFormat, _ := cmd.Flags().GetString("color-format")
	colorsPerRow, _ := cmd.Flags().GetInt("output-row-count")
	isPalettePreview, _ := cmd.Flags().GetBool("palette-preview")

	return &struct {
		Filename         string
		ColorFormat      string
		ColorsPerRow     int
		IsPalettePreview bool
	}{
		Filename:         filename,
		ColorFormat:      string(utils.ColorFormatFromString(colorFormat)),
		ColorsPerRow:     colorsPerRow,
		IsPalettePreview: isPalettePreview,
	}, nil
}

func validateInput(filename string) error {
	validExtension := utils.СheckFileExtension(filename,
		append(aseprite.SpritesExtensions(), aseprite.AvailablePaletteExtensions()...)...,
	)

	if !utils.СheckFileExists(filename, false) || !validExtension {
		return fmt.Errorf("invalid file: %s", filename)
	}
	return nil
}

func initializeGenerator(cfg *config.Config) *preview.Generator {
	aseCli := aseprite.NewCLI(cfg.AsepritePath, cfg.ScriptDirPath)
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
