package types

type AssetsDirs struct {
	folderPath  string
	assetsNames []string
}

func (a *AssetsDirs) GetAssetsNames() []string {
	return a.assetsNames
}

func (a *AssetsDirs) GetFolderPath() string {
	return a.folderPath
}

func NewAssetsDirs(folderPath string, assetsNames []string) AssetsDirs {
	return AssetsDirs{
		folderPath:  folderPath,
		assetsNames: assetsNames,
	}
}
