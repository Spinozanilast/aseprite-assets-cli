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

func (info *AssetInfo) GeneratePreview(cli *aseprite.Cli, size int) (string, error) {
	params := preview.GenerateParams{
		Filename:     info.Path,
		ColorsPerRow: (size - 11) / 11,
		Size:         size / 4,
	}

	generator := preview.NewGenerator(cli)

	if output, err := generator.Generate(params); err != nil {
		return "", err
	} else {
		return output, nil
	}
}
