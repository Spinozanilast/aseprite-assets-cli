package commands

import "github.com/spinozanilast/aseprite-assets-cli/aseprite"

type SavePaletteAsPresetCommand struct {
	Ui              bool
	PresetName      string `script:"preset-name"`
	PaletteFilename string `script:"palette-filename"`
}

func (a *SavePaletteAsPresetCommand) GetArgs() []string {
	return aseprite.CreateArgsFromStruct(a)
}

func (a *SavePaletteAsPresetCommand) GetScriptName() string {
	return "save-palette-as-preset.lua"
}
