package info

import (
	"fmt"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/aseprite"
	"os"

	"github.com/spinozanilast/aseprite-assets-cli/pkg/assets"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/consts"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/utils"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	AssetInfo *assets.AssetInfo
	cli       *aseprite.AsepriteCLI
	Styles    *Styles
	Width     int
	Height    int
	Error     string
}

func NewInfoModel(cli *aseprite.AsepriteCLI) Model {
	return Model{
		AssetInfo: &assets.AssetInfo{},
		cli:       cli,
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

	assetInfo := &assets.AssetInfo{
		Name:      info.Name(),
		Path:      assetPath,
		Size:      info.Size(),
		ModTime:   info.ModTime(),
		Extension: ext,
		Type:      assetType,
	}

	if m.Width != 0 {
		preview, err := assetInfo.GeneratePreview(m.cli, (m.Width-2)/4)
		if err != nil {
			assetInfo.Preview = fmt.Sprintf("error generating preview for asset: %v", err)
		} else {
			assetInfo.Preview = preview
		}

	}
	m.AssetInfo = assetInfo

	m.Error = ""
}
