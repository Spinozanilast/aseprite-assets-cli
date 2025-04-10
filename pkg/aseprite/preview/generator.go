package preview

import (
	"fmt"

	"github.com/spinozanilast/aseprite-assets-cli/pkg/aseprite"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/aseprite/commands"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/utils/files"
)

var (
	ErrUnsupportedFileType = fmt.Errorf("unsupported file type")
)

type Generator struct {
	aseCli *aseprite.Cli
}

func NewGenerator(aseCli *aseprite.Cli) *Generator {
	return &Generator{aseCli: aseCli}
}

type GenerateParams struct {
	Filename         string
	ColorFormat      string
	ColorsPerRow     int
	IsPalettePreview bool
	Size             int
}

func (g *Generator) Generate(params GenerateParams) (string, error) {
	var cmd aseprite.Command

	switch {
	case params.IsPalettePreview:
		cmd = g.createPaletteCommand(params)
	case files.CheckFileExtension(params.Filename, aseprite.SpritesExtensions()...):
		cmd = g.createSpriteCommand(params)
	case files.CheckFileExtension(params.Filename, aseprite.AvailablePaletteExtensions()...):
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
		SpriteFilename: params.Filename,
		SpriteSize:     params.Size,
	}
}

func (g *Generator) createPaletteCommand(params GenerateParams) *commands.ShowPalette {
	return &commands.ShowPalette{
		PaletteFilename: params.Filename,
		OutputRowCount:  params.ColorsPerRow,
		ColorFormat:     params.ColorFormat,
	}
}
