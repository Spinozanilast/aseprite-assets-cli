package show

import (
	"fmt"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/aseprite"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/aseprite/commands"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/config"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/environment"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/utils"
)

type showCmd struct {
	config           *config.Config
	asepriteCli      *aseprite.AsepriteCLI
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
		# Show asset preview based on file extension
		$ aseprite-assets-cli show -f <fullpath-to-file>

		# Show asset preview by filename
		$ aseprite-assets-cli show -f <fullpath-to-file>

		# Show palette preview by filename with custom colors format and palette size
		$ aseprite-assets-cli show -f <fullpath-to-file> -c rgb -r 10 -p

		# Show palette preview by filename
		$ aseprite-assets-cli show -f <fullpath-to-file> -p
	`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.LoadConfig()
			if err != nil {
				return err
			}

			colorFormat, _ := cmd.Flags().GetString("color-format")
			colorsPerRow, _ := cmd.Flags().GetInt("output-row-count")
			isPalettePreview, _ := cmd.Flags().GetBool("palette-preview")
			filename, err := cmd.Flags().GetString("filename")
			if err != nil {
				return fmt.Errorf("failed to get filename flag: %w", err)
			}

			hasAvailableExt := utils.СheckFileExtension(filename, aseprite.AvailablePaletteExtensions()...)
			if !utils.СheckFileExists(filename, false) || !hasAvailableExt {
				return fmt.Errorf("file does not exist")
			}

			showCmd := showCmd{
				config:           cfg,
				asepriteCli:      aseprite.NewCLI(cfg.AsepritePath, cfg.ScriptDirPath),
				filename:         filename,
				colorFormat:      string(utils.ColorFormatFromString(colorFormat)),
				colorsPerRow:     colorsPerRow,
				isPalettePreview: isPalettePreview,
			}

			err = showCmd.showAsset()
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringP("filename", "f", "", "asset filename")
	cmd.Flags().StringP("color-format", "c", "hex", "color format of values in palette (only for palettes preview)")
	cmd.Flags().IntP("output-row-count", "r", 5, "number of colors in one row (only for palettes preview)")
	cmd.Flags().BoolP("palette-preview", "p", false, "show palette preview (if you want to see current palette used in .ase or .aseprite files)")
	cmd.MarkFlagRequired("filename")

	return cmd
}

func (h *showCmd) showAsset() error {
	err := h.showAsepriteAsset()

	if err != nil {
		return err
	}

	return nil
}

func (h *showCmd) showAsepriteAsset() error {
	var err error
	var output string
	if utils.СheckFileExtension(h.filename, aseprite.SpritesExtensions()...) && !h.isPalettePreview {
		showSpriteCmd := &commands.ShowSprite{
			BatchMode:      true,
			SpriteFilename: h.filename,
		}
		output, err = h.asepriteCli.ExecuteCommandOutput(showSpriteCmd)
	} else {
		showPaletteCmd := &commands.ShowPalette{
			BatchMode:       true,
			PaletteFilename: h.filename,
			OutputRowCount:  h.colorsPerRow,
			ColorFormat:     h.colorFormat,
		}
		output, err = h.asepriteCli.ExecuteCommandOutput(showPaletteCmd)
	}

	if err != nil {
		return err
	}

	if output == "" {
		utils.PrintFormatted("Sorry, there is no preview for this file\n")
		return nil
	}

	fmt.Println(output)

	return nil
}
