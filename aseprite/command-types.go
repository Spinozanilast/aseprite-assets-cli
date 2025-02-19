package aseprite

type AsepriteAssetCreateCommand struct {
	Ui         bool
	Width      int
	Height     int
	ColorMode  string `script:"color-mode"`
	OutputPath string `script:"output-path"`
}

func (a *AsepriteAssetCreateCommand) GetArgs() []string {
	return createArgsFromStruct(a)
}

func (a *AsepriteAssetCreateCommand) GetScriptName() string {
	return "create-file.lua"
}
