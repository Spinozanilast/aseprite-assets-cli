package commands

import "github.com/spinozanilast/aseprite-assets-cli/pkg/aseprite"

type AssetCreate struct {
	BatchMode  bool
	Width      int
	Height     int
	ColorMode  string `script:"color-mode"`
	OutputPath string `script:"output-path" format:"quotes"`
}

func (c *AssetCreate) ScriptName() string {
	return "sprite-file.lua"
}

func (c *AssetCreate) Args() []string {
	return aseprite.CreateArgsFromStruct(c)
}
