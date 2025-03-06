package preview

import (
	"fmt"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/aseprite"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/aseprite/commands"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/utils"
)

var (
	ErrUnsupportedFileType = fmt.Errorf("unsupported file type")
)

type Generator struct {
	aseCli *aseprite.AsepriteCLI
}

func NewGenerator(aseCli *aseprite.AsepriteCLI) *Generator {
	return &Generator{aseCli: aseCli}
}

type GenerateParams struct {
	Filename         string
	ColorFormat      string
	ColorsPerRow     int
	IsPalettePreview bool
}

func (g *Generator) Generate(params GenerateParams) (string, error) {
	var cmd aseprite.Command

	switch {
	case params.IsPalettePreview:
		cmd = g.createPaletteCommand(params)
	case utils.СheckFileExtension(params.Filename, aseprite.SpritesExtensions()...):
		cmd = g.createSpriteCommand(params)
	case utils.СheckFileExtension(params.Filename, aseprite.AvailablePaletteExtensions()...):
		cmd = g.createPaletteCommand(params)
	default:
		return "", fmt.Errorf("%w: %s", ErrUnsupportedFileType, params.Filename)
	}

	output, err := g.aseCli.ExecuteCommandOutput(cmd)
	if err != nil {
		return "", fmt.Errorf("aseprite execution failed: %w", err)
	}

	return output, nil
}

func (g *Generator) createSpriteCommand(params GenerateParams) *commands.ShowSprite {
	return &commands.ShowSprite{
		BatchMode:      true,
		SpriteFilename: params.Filename,
	}
}

func (g *Generator) createPaletteCommand(params GenerateParams) *commands.ShowPalette {
	return &commands.ShowPalette{
		BatchMode:       true,
		PaletteFilename: params.Filename,
		OutputRowCount:  params.ColorsPerRow,
		ColorFormat:     params.ColorFormat,
	}
}
