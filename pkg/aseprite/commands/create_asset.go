package commands

import (
	"github.com/spinozanilast/aseprite-assets-cli/pkg/aseprite"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/utils"
)

type CreateSprite struct {
	OpenAfterCreation bool `script:"ignore"`
	Width             int
	Height            int
	ColorMode         string `script:"color-mode"`
	OutputPath        string `script:"output-path" format:"quotes"`
}

func (c *CreateSprite) ScriptName() string {
	return "create-sprite.lua"
}

func (c *CreateSprite) Args() []string {
	return aseprite.CreateArgsFromStruct(c)
}

func (c *CreateSprite) ScriptCallback(asePath string) (openingCallback func()) {
	return func() {
		if c.OpenAfterCreation && utils.СheckFileExists(c.OutputPath, false) {
			err := utils.OpenFileWith(asePath, c.OutputPath)
			if err != nil {
				return
			}
		}
	}
}
