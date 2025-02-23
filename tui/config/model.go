package config

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/spinozanilast/aseprite-assets-cli/utils"

	tea "github.com/charmbracelet/bubbletea"
	config "github.com/spinozanilast/aseprite-assets-cli/config"
)

type FldType int
type Direction int

type inputField struct {
	textinput.Model
	status      string
	description string
	fldType     FldType
}

type AppState int

const RequiredFieldsNum = 3

const (
	AppPathFld FldType = iota
	OpenAiKeyFld
	OpenAiUrlFld
	PalettesFolderPathFld
	AssetsFolderPathFld
)

const (
	statusValid   = "valid"
	statusInvalid = "invalid"
	statusNeutral = "neutral"
)

const (
	Up   Direction = -1
	Down Direction = 1
)

const (
	StateConfiguring AppState = iota
	StateCompleted
)

type model struct {
	state            AppState
	fields           []*inputField
	activeFieldIndex int
	quitting         bool
	styles           *Styles
	keys             keyMap
	help             help.Model
	err              string
}

type keyMap struct {
	Up            key.Binding
	Down          key.Binding
	Enter         key.Binding
	Clear         key.Binding
	Help          key.Binding
	Save          key.Binding
	AddAssetDir   key.Binding
	AddPaletteDir key.Binding
	Quit          key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.AddAssetDir, k.AddPaletteDir, k.Enter, k.Clear, k.Help, k.Save, k.Quit},
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
	AddAssetDir: key.NewBinding(
		key.WithKeys("ctrl+a"),
		key.WithHelp("Ctrl+A", "add assets directory"),
	),
	AddPaletteDir: key.NewBinding(
		key.WithKeys("ctrl+p"),
		key.WithHelp("Ctrl+P", "add palette directory"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("Enter", "confirm or browse file/directory (if empty)"),
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

	if config.AsepritePath != "" {
		model.fields[AppPathFld].SetValue(config.AsepritePath)
	}

	for _, path := range config.AssetsFolderPaths {
		inputField := newAssetInputField()
		inputField.SetValue(path)
		model.fields = append(model.fields, inputField)
	}

	model.fields[OpenAiKeyFld].SetValue(config.OpenAiConfig.ApiKey)
	model.fields[OpenAiUrlFld].SetValue(config.OpenAiConfig.ApiUrl)

	model.validateAllFields()

	return model
}

func blankInitialModel() model {
	appPathField := newInputField("Enter Aseprite executable path or press Enter to browse...", "Asseprite Executable Location", AppPathFld)
	appPathField.Placeholder = "Enter aseprite executable path or tap enter to open file dialog"
	appPathField.Focus()

	h := help.New()
	h.ShowAll = true

	return model{
		state:            StateConfiguring,
		fields:           []*inputField{appPathField, newInputField("Enter OpenAI API key", "OpenAI API Key", OpenAiKeyFld), newInputField("Enter OpenAI API URL", "OpenAI API URL", OpenAiUrlFld)},
		activeFieldIndex: 0,
		styles:           DefaultStyles(),
		keys:             keys,
		help:             h,
	}
}

func newAssetInputField() *inputField {
	return newInputField("Enter Aseprite assets directory or press Enter to browse", "Assets Directory Path", AssetsFolderPathFld)
}

func newPalettesInputField() *inputField {
	return newInputField("Enter Palettes directory or press Enter to browse", "Palettes Directory Path", PalettesFolderPathFld)
}

func newInputField(placeholder string, description string, fldType FldType) *inputField {
	textInput := textinput.New()
	textInput.Placeholder = placeholder
	textInput.ShowSuggestions = true

	return &inputField{
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
		case key.Matches(msg, m.keys.AddAssetDir):
			return m.handleAddDirectory(AssetsFolderPathFld)
		case key.Matches(msg, m.keys.AddPaletteDir):
			return m.handleAddDirectory(PalettesFolderPathFld)
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

			duplicatesExist, err := m.checkFieldsDuplicatesExist()

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
	if m.activeFieldIndex >= 0 && m.activeFieldIndex < len(m.fields) {
		return m.fields[m.activeFieldIndex]
	}
	return nil
}

func (m *model) validateCurrentInput() {
	field := m.currentField()
	m.validateField(field)
}

func (m *model) validateField(fld *inputField) {
	if fld == nil {
		return
	}

	value := fld.Value()
	fldType := fld.fldType

	if strings.TrimSpace(value) == "" || value == "\\" {
		fld.status = statusInvalid
		return
	}

	if fldType == AppPathFld && utils.СheckFileExists(value, false) {
		fld.status = statusValid
	} else if fldType == OpenAiUrlFld {
		_, err := url.ParseRequestURI(value)
		if err == nil {
			fld.status = statusValid
		}
	} else if fldType == OpenAiKeyFld {
		if strings.HasPrefix(value, "sk-") {
			fld.status = statusValid
		}
	} else if fldType.IsInTypes(AssetsFolderPathFld, PalettesFolderPathFld) && utils.СheckFileExists(value, true) {
		fld.status = statusValid
	} else {
		fld.status = statusInvalid
	}
}

func (m *model) validateAllFields() {
	for _, field := range m.fields {
		m.validateField(field)
	}
}

func (m *model) handleEnterKey() (tea.Model, tea.Cmd) {
	current := m.currentField()

	if current == nil {
		return m, nil
	}

	if current.isEmpty() && !current.fldType.IsInTypes(OpenAiKeyFld, OpenAiUrlFld) {
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

	fldType := current.fldType
	var path string
	var err error

	if fldType == AppPathFld {
		path, err = utils.OpenExecutableFilesDialog("Select Aseprite executable file")
	} else if fldType == PalettesFolderPathFld {
		path, err = utils.OpenDirectoryDialog("Select palettes store folder directory")
	} else {
		path, err = utils.OpenDirectoryDialog("Select Aseprite assets directory")
	}

	if err == nil && path != "" {
		current.SetValue(path)
		m.validateCurrentInput()
	}

	return m, nil
}

func (m model) handleAddDirectory(fldType FldType) (tea.Model, tea.Cmd) {
	duplicatesExist, err := m.checkFieldsDuplicatesExist()

	if m.allFieldsValid() && !duplicatesExist {
		if fldType == AssetsFolderPathFld {
			m.fields = append(m.fields, newAssetInputField())
		} else if fldType == PalettesFolderPathFld {
			m.fields = append(m.fields, newPalettesInputField())
		}
		m.clearError()
	}

	if duplicatesExist {
		m.err = err.Error()
	}

	return m, nil
}

func (m model) handleClearInput() (tea.Model, tea.Cmd) {
	currentField := m.currentField()

	isFolderFld := currentField.fldType.IsInTypes(AssetsFolderPathFld, PalettesFolderPathFld)
	if currentField.isEmpty() && isFolderFld && len(m.fields) > RequiredFieldsNum {
		m.removeInput(currentField)
		return m.moveFocus(Up), nil
	}

	currentField.Reset()
	m.validateField(currentField)

	return m, nil
}

func (m model) updateCurrentInput(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	currentField := m.currentField()

	if currentField == nil {
		return m, nil
	}

	currentField.Model, cmd = currentField.Update(msg)

	m.validateField(currentField)

	return m, cmd
}

func (m *model) moveFocus(direction Direction) *model {
	newFocusIndex := m.activeFieldIndex + int(direction)

	if newFocusIndex < 0 {
		newFocusIndex = len(m.fields) - 1
	} else if newFocusIndex >= len(m.fields) {
		newFocusIndex = 0
	}

	if current := m.currentField(); current != nil {
		current.Blur()
	}

	m.activeFieldIndex = newFocusIndex

	if newField := m.currentField(); newField != nil {
		newField.Focus()
	}
	return m
}

func (m *model) ToConfig() error {
	appPath := m.fields[AppPathFld].Value()

	assetsDirs := make([]string, 0)
	palettesDirs := make([]string, 0)
	for _, field := range m.fields {
		fldType := field.fldType
		if fldType == AssetsFolderPathFld {
			assetsDirs = append(assetsDirs, field.Value())
		} else if fldType == PalettesFolderPathFld {
			palettesDirs = append(palettesDirs, field.Value())
		}
	}

	openAiKeyFld := m.fields[OpenAiKeyFld]
	openAiUrlFld := m.fields[OpenAiUrlFld]

	if openAiKeyFld.status == statusValid || openAiUrlFld.status == statusValid {
		openAiKey := m.fields[OpenAiKeyFld].Value()
		openAiUrl := m.fields[OpenAiUrlFld].Value()
		config.SetOpenAiConfig(openAiKey, openAiUrl)
	}

	return config.SavePaths(appPath, assetsDirs, palettesDirs)
}

func (m *model) allFieldsValid() bool {
	for _, field := range m.fields {
		// Open Ai fields are optional
		isOpenAiFld := field.fldType.IsInTypes(OpenAiKeyFld, OpenAiUrlFld)
		if field.status != statusValid && !isOpenAiFld {
			return false
		}
	}

	m.clearError()
	return true
}

func (m *model) removeInput(fldToRemove *inputField) {
	removeIdx := -1
	for i, field := range m.fields {
		if field == fldToRemove {
			removeIdx = i
			break
		}
	}

	if removeIdx == -1 {
		return
	}

	if len(m.fields) > 1 {
		m.fields = append(m.fields[:removeIdx], m.fields[removeIdx+1:]...)
	}
}

func (m *model) checkFieldsDuplicatesExist() (bool, error) {
	uniqueDirs := make(map[string]struct{})

	for _, field := range m.fields {
		if field.status == statusValid && field.fldType.IsInTypes(AssetsFolderPathFld, PalettesFolderPathFld) {
			dir := field.Value()
			uniqueDirs[dir] = struct{}{}
		}
	}

	areDuplicatesExists := len(uniqueDirs) != (len(m.fields) - RequiredFieldsNum)

	if areDuplicatesExists {
		return true, fmt.Errorf("duplicate assets (or empty) directories found")
	}

	return false, nil
}

func (m *model) clearError() {
	m.err = ""
}

func (fldType FldType) IsInTypes(fieldTypes ...FldType) bool {
	if len(fieldTypes) == 0 {
		return false
	}

	if len(fieldTypes) == 1 {
		return fldType == fieldTypes[0]
	}

	for _, t := range fieldTypes {
		if fldType == t {
			return true
		}
	}
	return false
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
