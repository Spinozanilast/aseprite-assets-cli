package aseprite

type ColorMode string

const (
	ColorModeIndexed ColorMode = "indexed"
	ColorModeRGB     ColorMode = "rgb"
	ColorModeGray    ColorMode = "grayscale"
)

func (c ColorMode) String() string {
	return string(c)
}

func ColorModes() []string {
	return []string{ColorModeIndexed.String(), ColorModeRGB.String(), ColorModeGray.String()}
}
