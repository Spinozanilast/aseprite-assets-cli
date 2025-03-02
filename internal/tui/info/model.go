package info

import (
	"os"

	"github.com/spinozanilast/aseprite-assets-cli/pkg/assets"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/consts"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/utils"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	AssetInfo *assets.AssetInfo
	Styles    *Styles
	Width     int
	Height    int
	Error     string
}

func NewInfoModel() Model {
	return Model{
		AssetInfo: &assets.AssetInfo{},
		Styles:    DefaultStyles(),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Height = msg.Height
		m.Width = msg.Width
	}

	return m, nil
}

func (m *Model) UpdateAssetInfo(assetPath string, assetType consts.AssetsType) {
	if assetPath == "" {
		m.AssetInfo = nil
		return
	}

	info, err := os.Stat(assetPath)
	if err != nil {
		m.Error = err.Error()
		return
	}

	ext := utils.GetFileExtension(assetPath)
	m.AssetInfo = &assets.AssetInfo{
		Name:      info.Name(),
		Path:      assetPath,
		Size:      info.Size(),
		ModTime:   info.ModTime(),
		Extension: ext,
		Type:      assetType,
	}

	m.Error = ""
}
