package assets

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/utils"

	tea "github.com/charmbracelet/bubbletea"
)

type Direction int

const (
	Left  Direction = -1
	Right Direction = 1
)

type item string

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {

	listStyles := DefaultListItemStyles()
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := listStyles.ItemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return listStyles.SelectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

type Source struct {
	folderPath  string
	assetsNames []string
}

func (a *Source) GetAssetsNames() []string {
	return a.assetsNames
}

func (a *Source) GetFolderPath() string {
	return a.folderPath
}

func NewAssetsSource(folderPath string, assetsNames []string) Source {
	return Source{
		folderPath:  folderPath,
		assetsNames: assetsNames,
	}
}

type Model struct {
	list             list.Model
	appPath          string
	assetsFolders    []Source
	assetsActive     []int
	activeFolderIdx  int
	activeFolderName string
	prevFolderName   string
	nextFolderName   string
	title            string
	styles           *Styles
	keys             keyMap
	help             help.Model
	err              string
	quitting         bool
	appWidth         int
}

type keyMap struct {
	Left  key.Binding
	Right key.Binding
	Enter key.Binding
	Tab   key.Binding
	Help  key.Binding
	Quit  key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Left, k.Right, k.Tab, k.Enter, k.Help, k.Quit},
	}
}

var keys = keyMap{
	Right: key.NewBinding(
		key.WithKeys("right"),
		key.WithHelp("→", "move to right folder"),
	),
	Left: key.NewBinding(
		key.WithKeys("left"),
		key.WithHelp("←", "move to left folder"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("Enter ", "confirm or browse file/directory"),
	),
	Help: key.NewBinding(
		key.WithKeys("ctrl+h"),
		key.WithHelp("Ctrl+H", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("esc", "ctrl+c"),
		key.WithHelp("ESC/Ctrl+C", "quit"),
	),
}

func InitialModel(title string, appPath string, assetsFolders []Source) Model {
	h := help.New()
	h.ShowAll = true
	folderLength := len(assetsFolders)

	listLayoutStyles := DefaultListLayoutStyles()
	assetsActive := make([]int, folderLength)
	assetsActive[0] = 0
	items := createListItems(assetsFolders[0].GetAssetsNames())

	activeFolderName := assetsFolders[0].GetFolderPath()

	var prevFolderName string
	var nextFolderName string

	switch {
	case folderLength == 2:
		prevFolderName = assetsFolders[1].GetFolderPath()
	case folderLength > 2:
		prevFolderName = assetsFolders[folderLength-1].GetFolderPath()
		nextFolderName = assetsFolders[1].GetFolderPath()
	}

	list := list.New(items, itemDelegate{}, listLayoutStyles.ListWidth, listLayoutStyles.ListHeight)

	return Model{
		appPath:          appPath,
		list:             list,
		assetsFolders:    assetsFolders,
		assetsActive:     assetsActive,
		activeFolderIdx:  0,
		activeFolderName: activeFolderName,
		prevFolderName:   prevFolderName,
		nextFolderName:   nextFolderName,
		styles:           DefaultStyles(),
		keys:             keys,
		help:             h,
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
		m.list.SetWidth(msg.Width)
		m.appWidth = msg.Width
		return m, nil
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
			return m.moveBetweenFoldersFocus(Left), nil
		case key.Matches(msg, m.keys.Right):
			return m.moveBetweenFoldersFocus(Right), nil
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
			return m, nil
		}
	}

	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m *Model) handleEnterKey() (tea.Model, tea.Cmd) {
	selectedItem := m.currentAssets()[m.list.Index()]
	filename := filepath.Join(m.activeFolderName, selectedItem)
	utils.OpenFileWith(filename, m.appPath)
	return m, nil
}

func (m *Model) moveBetweenFoldersFocus(direction Direction) *Model {
	total := len(m.assetsFolders)
	m.assetsActive[m.activeFolderIdx] = m.list.Index()

	m.activeFolderIdx = (m.activeFolderIdx + int(direction) + total) % total
	m.updateFolderMetadata()
	m.updateListContent()

	return m
}

func (m *Model) updateListContent() {
	storedPos := m.assetsActive[m.activeFolderIdx]

	items := createListItems(m.currentAssets())
	m.list.SetItems(items)
	m.list.Select(storedPos)

	m.list.SetSize(m.list.Width(), m.list.Height())
}

func (m *Model) currentAssets() []string {
	return m.assetsFolders[m.activeFolderIdx].GetAssetsNames()
}

func (m *Model) updateFolderMetadata() {
	total := len(m.assetsFolders)

	if total == 1 {
		m.prevFolderName = ""
		m.nextFolderName = ""
		return
	}

	prevIndex := (m.activeFolderIdx - 1 + total) % total
	nextIndex := (m.activeFolderIdx + 1) % total

	m.prevFolderName = m.assetsFolders[prevIndex].GetFolderPath()
	m.activeFolderName = m.assetsFolders[m.activeFolderIdx].GetFolderPath()
	m.nextFolderName = m.assetsFolders[nextIndex].GetFolderPath()
}

func (i item) FilterValue() string { return string(i) }

func createListItems(texts []string) []list.Item {
	items := make([]list.Item, len(texts))
	for i, text := range texts {
		items[i] = item(text)
	}

	return items
}
