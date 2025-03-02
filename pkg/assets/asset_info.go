package assets

import (
	"github.com/spinozanilast/aseprite-assets-cli/pkg/consts"
	"time"
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

func (info *AssetInfo) GeneratePreview() string {
	panic("unimplemented")
}
