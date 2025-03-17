package tui

import (
	"fmt"

	"github.com/spinozanilast/aseprite-assets-cli/internal/tui/list"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/config"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/steam"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/utils/files"

	tea "github.com/charmbracelet/bubbletea"
	configTui "github.com/spinozanilast/aseprite-assets-cli/internal/tui/config"
)

// StartConfigTui starts terminal user interface for configuring cli
func StartConfigTui(cfg *config.Config) error {
	p := tea.NewProgram(configTui.InitialModel(cfg))
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("cfg TUI error: %w", err)
	}

	return nil
}

// StartListTui starts terminal user interface for listing existing assets
func StartListTui(params list.ListParams, cfg *config.Config, startWithSteam bool) error {
	openAppFunc := SelectOpenHandleFunc(cfg, startWithSteam)
	p := tea.NewProgram(list.NewModel(params, openAppFunc, cfg.FromSteam), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("list TUI error: '%s'", err)
	}
	return nil
}

// SelectOpenHandleFunc returns func for opening assets.
// Result depends and could be steam based or direct app cli based
func SelectOpenHandleFunc(config *config.Config, startWithSteam bool) func(string) error {
	if config.FromSteam && startWithSteam {
		return func(filename string) error {
			if _, err := steam.StartSteamAppById(config.AsepritePath, filename); err != nil {
				return err
			}

			return nil
		}
	} else {
		return func(filename string) error {
			if err := files.OpenFileWith(filename, config.AsepritePath); err != nil {
				return err
			}
			return nil
		}
	}
}
