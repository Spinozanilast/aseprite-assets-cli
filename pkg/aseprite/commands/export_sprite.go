package commands

import "github.com/spinozanilast/aseprite-assets-cli/pkg/aseprite"

type ExportSprite struct {
	BatchMode      bool
	SpriteFilename string `script:"sprite-filename" format:"quotes"`
	OutputFilename string `script:"output-filename" format:"quotes"`
	Format         string
	Sizes          string
	Scales         string
}

func (c *ExportSprite) ScriptName() string {
	return "export-sprite.lua"
}

func (c *ExportSprite) Args() []string {
	return aseprite.CreateArgsFromStruct(c)
}
