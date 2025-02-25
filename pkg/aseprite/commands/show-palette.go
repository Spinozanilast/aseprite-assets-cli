package commands

import "github.com/spinozanilast/aseprite-assets-cli/pkg/aseprite"

type ShowPalette struct {
	BatchMode       bool
	PaletteFilename string `script:"palette-filename" format:"quotes"`
	OutputRowCount  int    `script:"output-row-count"`
	ColorFormat     string `script:"color-format"`
}

func (a *ShowPalette) ScriptName() string {
	return "palette-preview.lua"
}

func (a *ShowPalette) Args() []string {
	return aseprite.CreateArgsFromStruct(a)
}
