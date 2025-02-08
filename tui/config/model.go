package config

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"

	tea "github.com/charmbracelet/bubbletea"
	config "github.com/spinozanilast/aseprite-assets-cli/config"
	utils "github.com/spinozanilast/aseprite-assets-cli/util"
)

const (
	statusValid   = "valid"
	statusInvalid = "invalid"
	statusNeutral = "neutral"
)

type FldType int

const (
	AppPathFld          FldType = 0
	AssetsFolderPathFld FldType = 1
)

type Direction int

const (
	Up   Direction = -1
	Down Direction = 1
)

type inputField struct {
	textinput.Model
	status      string
	description string
	fldType     FldType
}

type AppState int

const (
	StateConfiguring AppState = iota
	StateCompleted
)

type ErrorField struct {
	Field *inputField
	Err   error
}

type model struct {
	state          AppState
	appPathFld     inputField
	assetsDirsFlds []inputField
	activeInputIdx int
	quitting       bool
	styles         *Styles
	keys           keyMap
	help           help.Model
	err            string
}

type keyMap struct {
	Up    key.Binding
	Down  key.Binding
	Tab   key.Binding
	Enter key.Binding
	Clear key.Binding
	Help  key.Binding
	Save  key.Binding
	Quit  key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Tab, k.Enter, k.Clear, k.Help, k.Save, k.Quit},
	}
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up"),
		key.WithHelp("↑", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down"),
		key.WithHelp("↓", "move down"),
	),
	Tab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("TAB", "add directory (when on last and it is valid)"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("Enter ", "confirm or browse file/directory"),
	),
	Clear: key.NewBinding(
		key.WithKeys("ctrl+d"),
		key.WithHelp("Ctrl+D", "clear (remove current assets folder input if empty)"),
	),
	Help: key.NewBinding(
		key.WithKeys("ctrl+h"),
		key.WithHelp("Ctrl+H", "toggle help"),
	),
	Save: key.NewBinding(
		key.WithKeys("ctrl+s"),
		key.WithHelp("Ctrl+S", "save configuration"),
	),
	Quit: key.NewBinding(
		key.WithKeys("esc", "ctrl+c"),
		key.WithHelp("ESC/Ctrl+C", "quit"),
	),
}

func InitialModel(config *config.Config) model {
	model := blankInitialModel()
	model.assetsDirsFlds = nil

	if (config.AsepritePath != "") || (len(config.AssetsFolderPaths) > 0) {
		model.appPathFld.SetValue(config.AsepritePath)
		for _, path := range config.AssetsFolderPaths {
			inputField := newAssetInputField()
			inputField.SetValue(path)
			model.validateFld(&inputField)
			model.assetsDirsFlds = append(model.assetsDirsFlds, inputField)
		}
	}

	return model
}

func blankInitialModel() model {
	appPathField := newInputField("Enter Aseprite executable path or press Enter to browse...", "Asseprite Executable Location", AppPathFld)
	appPathField.Placeholder = "Enter aseprite executable path or tap enter to open file dialog"
	appPathField.Focus()

	h := help.New()
	h.ShowAll = true

	return model{
		state:          StateConfiguring,
		appPathFld:     appPathField,
		assetsDirsFlds: []inputField{newAssetInputField()},
		activeInputIdx: 0,
		styles:         DefaultStyles(),
		keys:           keys,
		help:           h,
	}
}

func newAssetInputField() inputField {
	return newInputField("Enter Aseprite assets directory or press Enter to browse", "Assets Directory Path", AssetsFolderPathFld)
}

func newInputField(placeholder string, description string, fldType FldType) inputField {
	textInput := textinput.New()
	textInput.Placeholder = placeholder

	return inputField{
		Model:       textInput,
		status:      statusNeutral,
		description: description,
		fldType:     fldType,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.quitting {
		return m, tea.Quit
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.help.Width = msg.Width
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, m.keys.Enter):
			return m.handleEnterKey()
		case key.Matches(msg, m.keys.Tab):
			return m.handleTabKey()
		case key.Matches(msg, m.keys.Clear):
			return m.handleClearInput()
		case key.Matches(msg, m.keys.Up):
			return m.moveFocus(Up), nil
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
			return m, nil
		case key.Matches(msg, m.keys.Save):
			if !m.allFieldsValid() {
				m.err = "Fill all fields with valid paths before saving"
				return m, nil
			}

			duplicatesExist, err := m.checkFldsDuplicatesExist()

			if duplicatesExist {
				m.err = err.Error()
				return m, nil
			}

			err = m.ToConfig()
			if err != nil {
				m.err = err.Error()
				return m, nil
			} else if m.err != "" {
				return m, nil
			}
			return m, tea.Quit
		case key.Matches(msg, m.keys.Down):
			return m.moveFocus(Down), nil
		}
	}

	return m.updateCurrentInput(msg)
}

func (m *model) currentField() *inputField {
	if m.activeInputIdx == 0 {
		return &m.appPathFld
	} else if m.activeInputIdx > 0 && m.activeInputIdx <= len(m.assetsDirsFlds) {
		return &m.assetsDirsFlds[m.activeInputIdx-1]
	}

	return nil
}

func (m *model) validateCurrentInput() {
	field := m.currentField()
	m.validateFld(field)
}

func (m *model) validateFld(fld *inputField) {
	if fld == nil {
		return
	}

	value := fld.Value()
	if value == "" || value == "/" {
		fld.status = statusInvalid
		return
	}

	if checkIfFileExists(value) {
		if fld.fldType == AppPathFld {
			fld.status = statusValid
			return
		}

		valid, err := utils.CheckAnyFileOfExtensionExists(value, ".aseprite")
		if valid {
			fld.status = statusValid
		} else {
			m.err = err.Error() + fmt.Errorf("\nChoose directory with assets inside").Error()
			fld.status = statusInvalid
		}
	} else {
		fld.status = statusInvalid
	}
}

func (m *model) handleEnterKey() (tea.Model, tea.Cmd) {
	current := m.currentField()
	if current == nil {
		return m, nil
	}

	if current.isEmpty() {
		return m.handleBrowse()
	}

	if current.status == statusValid {
		return m.moveFocus(Down), nil
	}

	return m, nil
}

func (m *model) handleBrowse() (tea.Model, tea.Cmd) {
	current := m.currentField()

	if current == nil {
		return m, nil
	}

	var path string
	var err error

	if m.activeInputIdx == 0 {
		path, err = utils.OpenExecutableFilesDialog("Select Aseprite executable file")
	} else {
		path, err = utils.OpenDirectoryDialog("Select Aseprite assets directory")
	}

	if err == nil && path != "" {
		current.SetValue(path)
		m.validateCurrentInput()
	}

	return m, nil
}

func (m model) handleTabKey() (tea.Model, tea.Cmd) {
	duplicatesExist, err := m.checkFldsDuplicatesExist()

	if m.allFieldsValid() && !duplicatesExist {
		m.assetsDirsFlds = append(m.assetsDirsFlds, newAssetInputField())
		m.clearError()
	}

	if duplicatesExist {
		m.err = err.Error()
	}

	return m, nil
}

func (m model) handleClearInput() (tea.Model, tea.Cmd) {
	currentField := m.currentField()

	if currentField.isEmpty() && currentField.fldType == AssetsFolderPathFld && len(m.assetsDirsFlds) > 1 {
		m.removeInput(currentField)
		return m.moveFocus(Up), nil
	}

	currentField.Reset()
	m.validateFld(currentField)

	return m, nil
}

func (m model) updateCurrentInput(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	currentField := m.currentField()

	if currentField == nil {
		return m, nil
	}

	currentField.Model, cmd = currentField.Update(msg)

	m.validateFld(currentField)

	return m, cmd
}

func (m *model) moveFocus(direction Direction) *model {
	maxInputIdx := len(m.assetsDirsFlds)
	newFocusIdx := m.activeInputIdx + int(direction)

	if newFocusIdx < 0 {
		newFocusIdx = maxInputIdx
	} else if newFocusIdx > maxInputIdx {
		newFocusIdx = 0
	}

	if current := m.currentField(); current != nil {
		current.Blur()
	}

	m.activeInputIdx = newFocusIdx

	if newField := m.currentField(); newField != nil {
		newField.Focus()
	}
	return m
}

func (m *model) ToConfig() error {
	appPath := m.appPathFld.Value()

	assetsDirs := make([]string, len(m.assetsDirsFlds))
	for i, dirFld := range m.assetsDirsFlds {
		assetsDirs[i] = dirFld.Value()
	}

	return config.SaveConfig(appPath, assetsDirs)
}

func (m *model) allFieldsValid() bool {
	if m.appPathFld.status != statusValid {
		return false
	}

	for _, field := range m.assetsDirsFlds {
		if field.status != statusValid {
			return false
		}
	}

	m.clearError()
	return true
}

func (m *model) removeInput(fldToRemove *inputField) {
	removeIdx := -1
	for i := range m.assetsDirsFlds {
		if &m.assetsDirsFlds[i] == fldToRemove {
			removeIdx = i
			break
		}
	}

	if removeIdx == -1 {
		return
	}

	if len(m.assetsDirsFlds) > 1 {
		m.assetsDirsFlds = append(
			m.assetsDirsFlds[:removeIdx],
			m.assetsDirsFlds[removeIdx+1:]...,
		)
	}
}

func (m *model) checkFldsDuplicatesExist() (bool, error) {
	uniqueAssetsDirs := make(map[string]struct{})

	for _, dirFld := range m.assetsDirsFlds {
		if dirFld.status == statusValid {
			dir := dirFld.Value()
			uniqueAssetsDirs[dir] = struct{}{}
		}
	}

	return len(uniqueAssetsDirs) != len(m.assetsDirsFlds), fmt.Errorf("duplicate assets (or empty) directories found")
}

func (m *model) clearError() {
	m.err = ""
}

func checkIfFileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func (fld *inputField) isEmpty() bool {
	return fld.Value() == ""
}

func (f *inputField) Focus() {
	f.Model.Focus()
}

func (f *inputField) Blur() {
	f.Model.Blur()
}
