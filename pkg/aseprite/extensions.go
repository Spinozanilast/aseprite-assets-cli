package aseprite

type PaletteExtension string
type SpriteExtension string

const (
	GPL PaletteExtension = ".gpl"
	PNG PaletteExtension = ".png"
)

const (
	Aseprite SpriteExtension = ".aseprite"
	Ase      SpriteExtension = ".ase"
)

func (f PaletteExtension) String() string {
	return string(f)
}

func (f SpriteExtension) String() string {
	return string(f)
}

func PaletteExtensions() []string {
	return []string{GPL.String(), PNG.String()}
}

func SpritesExtensions() []string {
	return []string{Aseprite.String(), Ase.String()}
}

func AvailableSupportedExtensions() []string {
	var extensions []string
	extensions = append(extensions, SpritesExtensions()...)
	extensions = append(extensions, PaletteExtensions()...)
	return extensions
}

// AvailablePaletteExtensions consists of all available palette extensions
// Also Aseprite and Ase extensions (from aseprite load palettes menu)
func AvailablePaletteExtensions() []string {
	return []string{
		".ase",
		".aseprite",
		".bmp",
		".flc",
		".fli",
		".gif",
		".ico",
		".jpeg",
		".jpg",
		".pcc",
		".png",
		".qoi",
		".tga",
		".webp",
		".act",
		".col",
		".gpl",
		".hex",
		".pal",
		".*",
	}
}

func AvailableExportExtensions() []string {
	extensions := []string{}
	extensions = append(extensions, SpritesExtensions()...)
	extensions = append(extensions,
		".bmp",
		".css",
		".flc",
		".fli",
		".gif",
		".ico",
		".jpeg",
		".jpg",
		".pcx",
		".pcc",
		".png",
		".qoi",
		".svg",
		".tga",
		".webp",
	)

	return extensions
}
