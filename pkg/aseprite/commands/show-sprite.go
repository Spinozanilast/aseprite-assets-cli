package commands

import "github.com/spinozanilast/aseprite-assets-cli/pkg/aseprite"

type ShowSprite struct {
	SpriteFilename string `script:"sprite-filename" format:"quotes"`
	SpriteSize     int    `script:"size"`
}

func (a *ShowSprite) ScriptName() string {
	return "sprite-preview.lua"
}

func (a *ShowSprite) Args() []string {
	return aseprite.CreateArgsFromStruct(a)
}
