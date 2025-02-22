package commands

import "github.com/spinozanilast/aseprite-assets-cli/aseprite"

type AssetCreateCommand struct {
	Ui         bool
	Width      int
	Height     int
	ColorMode  string `script:"color-mode"`
	OutputPath string `script:"output-path"`
}

func (a *AssetCreateCommand) GetArgs() []string {
	return aseprite.CreateArgsFromStruct(a)
}

func (a *AssetCreateCommand) GetScriptName() string {
	return "create-file.lua"
}
