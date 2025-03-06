package assets

import (
	"github.com/spinozanilast/aseprite-assets-cli/pkg/aseprite"
	"time"

	"github.com/spinozanilast/aseprite-assets-cli/pkg/aseprite/preview"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/consts"
)

type AssetInfo struct {
	Name      string
	Path      string
	Size      int64
	ModTime   time.Time
	Extension string
	Preview   string
	Type      consts.AssetsType
}

func (info *AssetInfo) GeneratePreview(cli *aseprite.AsepriteCLI, size int) (string, error) {
	params := preview.GenerateParams{
		Filename: info.Path,
		Size:     size,
	}

	generator := preview.NewGenerator(cli)

	if output, err := generator.Generate(params); err != nil {
		return "", err
	} else {
		return output, nil
	}
}
