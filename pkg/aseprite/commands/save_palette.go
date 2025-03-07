package commands

import "github.com/spinozanilast/aseprite-assets-cli/pkg/aseprite"

type SavePalette struct {
	PresetName      string `script:"preset-name" format:"quotes"`
	PaletteFilename string `script:"palette-filename"`
}

func (c *SavePalette) ScriptName() string {
	return "save-palette-as-preset.lua"
}

func (c *SavePalette) Args() []string {
	return aseprite.CreateArgsFromStruct(c)
}
