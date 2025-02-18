package aseprite

type AsepriteAssetCreateCommand struct {
	BaseCommand
	Ui         bool
	Width      int
	Height     int
	ColorMode  string `script:"color-mode"`
	OutputPath string `script:"output-path"`
}
