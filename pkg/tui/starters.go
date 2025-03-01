package tui

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	configTui "github.com/spinozanilast/aseprite-assets-cli/internal/tui/config"
	list "github.com/spinozanilast/aseprite-assets-cli/internal/tui/list"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/config"
)

func StartConfigTui(config *config.Config) error {
	p := tea.NewProgram(configTui.InitialModel(config))
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("config TUI error: %w", err)
	}

	return nil
}

func StartListTui(title string, executablePath string, sources []list.Source) error {
	p := tea.NewProgram(list.InitialModel(title, executablePath, sources))
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("list TUI error: '%s'", err)
	}
	return nil
}
