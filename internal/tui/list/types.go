package list

import (
	"github.com/spinozanilast/aseprite-assets-cli/pkg/consts"
)

type Direction int

const (
	Left  Direction = -1
	Right Direction = 1
)

type AssetSource struct {
	directory  string
	assetNames []string
}

func (a *AssetSource) GetAssetsNames() []string {
	return a.assetNames
}

func (a *AssetSource) GetFolderPath() string {
	return a.directory
}

func NewAssetsSource(folderPath string, assetsNames []string) AssetSource {
	return AssetSource{
		directory:  folderPath,
		assetNames: assetsNames,
	}
}

type ListParams struct {
	Title         string
	AppPath       string
	ScriptsPath   string
	AssetsFolders []AssetSource
	AssetsType    consts.AssetsType
}
