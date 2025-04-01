package helpers

import (
	"path/filepath"

	"github.com/spinozanilast/aseprite-assets-cli/pkg/aseprite"
)

type SpriteLayersNames struct {
	SpriteFilename string `script:"sprite-filename" format:"quotes"`
}

func (a *SpriteLayersNames) ScriptName() string {
	return filepath.Join(aseprite.ScriptsHelpersFolder, "sprite-layers-names.lua")
}

func (a *SpriteLayersNames) Args() []string {
	return aseprite.CreateArgsFromStruct(a)
}
