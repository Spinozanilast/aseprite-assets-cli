package list

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spinozanilast/aseprite-assets-cli/internal/tui/info"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/aseprite"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/consts"
	"path/filepath"
)

type folderNavigation struct {
	activeIdx int
	active    string
	prev      string
	next      string
}

type Model struct {
	list      list.Model
	infoPanel info.Model
	help      help.Model

	assetsType    consts.AssetsType
	assetsFolders []AssetSource
	assetsActive  []int

	nav folderNavigation

	title  string
	styles *Styles

	keys        keyMap
	err         string
	quitting    bool
	windowWidth int

	openAsset           func(string) error
	openActionAvailable bool
}

func NewModel(p ListParams, openAssetFunc func(string) error, fromSteam bool) Model {
	h := help.New()
	h.ShowAll = true

	folderLength := len(p.AssetsFolders)
	assetsActive := make([]int, folderLength)
	assetsActive[0] = 0
	items := createListItems(p.AssetsFolders[0].GetAssetsNames())

	nav := initFolderNavigation(p.AssetsFolders)

	listLayoutStyles := DefaultListLayoutStyles()
	listModel := list.New(items, itemDelegate{}, listLayoutStyles.ListWidth, listLayoutStyles.ListHeight)

	cli := aseprite.NewCLI(p.AppPath, p.ScriptsPath, fromSteam)
	infoModel := info.NewInfoModel(cli)

	modelKeys := keys
	modelKeys.showEnter = p.OpenActionAvailable

	return Model{
		list:                listModel,
		infoPanel:           infoModel,
		nav:                 nav,
		assetsActive:        assetsActive,
		styles:              DefaultStyles(),
		assetsFolders:       p.AssetsFolders,
		assetsType:          p.AssetsType,
		title:               p.Title,
		openActionAvailable: p.OpenActionAvailable,
		keys:                modelKeys,
		help:                h,
		openAsset:           openAssetFunc,
	}
}

func initFolderNavigation(sources []AssetSource) folderNavigation {
	folderLength := len(sources)
	activeFolderName := sources[0].GetFolderPath()

	var prevFolderName string
	var nextFolderName string

	switch {
	case folderLength == 2:
		prevFolderName = sources[1].GetFolderPath()
	case folderLength > 2:
		prevFolderName = sources[folderLength-1].GetFolderPath()
		nextFolderName = sources[1].GetFolderPath()
	}

	return folderNavigation{
		activeIdx: 0,
		active:    activeFolderName,
		prev:      prevFolderName,
		next:      nextFolderName,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.quitting {
		return m, tea.Quit
	}

	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m.handleResize(msg)
	case tea.KeyMsg:
		if m.list.FilterState() == list.Filtering {
			break
		}
		switch {
		case key.Matches(msg, m.keys.Quit):
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, m.keys.Enter):
			return m.handleEnterKey()
		case key.Matches(msg, m.keys.Left):
			return m.changeFolderFocus(Left), nil
		case key.Matches(msg, m.keys.Right):
			return m.changeFolderFocus(Right), nil
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
			return m, nil
		}
	}

	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel
	cmds = append(cmds, cmd)

	m.updateItemInfo()

	return m, tea.Batch(cmds...)
}

func (m *Model) handleEnterKey() (tea.Model, tea.Cmd) {
	if !m.openActionAvailable {
		return m, nil
	}

	err := m.openAsset(m.currentItemFilename())

	if err != nil {
		m.err = err.Error()
	}

	return m, nil
}

func (m *Model) handleResize(msg tea.WindowSizeMsg) (Model, tea.Cmd) {
	m.windowWidth = msg.Width

	m.list.SetWidth(msg.Width / 2)
	m.infoPanel.Width = msg.Width / 2

	m.updateItemInfo()
	return *m, nil
}

func (m *Model) updateListContent() {
	storedPos := m.assetsActive[m.nav.activeIdx]

	items := createListItems(m.currentAssets())
	m.list.SetItems(items)
	m.list.Select(storedPos)

	m.list.SetSize(m.list.Width(), m.list.Height())
}

func (m *Model) changeFolderFocus(direction Direction) *Model {
	activeIdx := m.nav.activeIdx
	total := len(m.assetsFolders)

	m.assetsActive[activeIdx] = m.list.Index()

	m.nav.activeIdx = (activeIdx + int(direction) + total) % total
	m.updateFoldersNavigation()
	m.updateListContent()
	m.updateItemInfo()
	return m
}

func (m *Model) updateFoldersNavigation() {
	total := len(m.assetsFolders)

	if total == 1 {
		m.nav.prev = ""
		m.nav.next = ""
		return
	}

	activeIdx := m.nav.activeIdx

	prevIndex := (activeIdx - 1 + total) % total
	nextIndex := (activeIdx + 1) % total

	m.nav.prev = m.assetsFolders[prevIndex].GetFolderPath()
	m.nav.active = m.assetsFolders[activeIdx].GetFolderPath()
	m.nav.next = m.assetsFolders[nextIndex].GetFolderPath()
}

func (m *Model) updateItemInfo() {
	m.infoPanel.UpdateAssetInfo(m.currentItemFilename(), m.assetsType)
}

func (m *Model) currentItemFilename() string {
	selectedItem := m.currentAssets()[m.list.Index()]
	return filepath.Join(m.nav.active, selectedItem)
}

func (m *Model) currentAssets() []string {
	return m.assetsFolders[m.nav.activeIdx].GetAssetsNames()
}
