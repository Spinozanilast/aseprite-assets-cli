package config

import (
	"github.com/charmbracelet/lipgloss"
)

type Styles struct {
	App                lipgloss.Style
	Title              lipgloss.Style
	InputLabel         lipgloss.Style
	Error              lipgloss.Style
	ErrorLabel         lipgloss.Style
	Help               lipgloss.Style
	ActiveInput        lipgloss.Style
	ValidInput         lipgloss.Style
	ValidActiveInput   lipgloss.Style
	InvalidInput       lipgloss.Style
	InvalidActiveInput lipgloss.Style
	NeutralInput       lipgloss.Style
	NeutralActiveInput lipgloss.Style
}

func DefaultStyles() *Styles {
	s := new(Styles)
	baseBorder := lipgloss.NormalBorder()
	minInputWidth := 60

	s.App = lipgloss.NewStyle().Padding(1, 2).Width(minInputWidth + 20)
	s.Title = lipgloss.NewStyle().
		Bold(true).Foreground(lipgloss.Color("63")).
		MarginBottom(1)
	s.ErrorLabel = lipgloss.NewStyle().
		Bold(true).Foreground(lipgloss.Color("88"))
	s.Error = lipgloss.NewStyle().Background(lipgloss.Color("240")).Foreground(lipgloss.Color("160"))
	s.InputLabel = lipgloss.NewStyle().MarginBottom(1)
	s.Help = lipgloss.NewStyle().Bold(true)

	s.ValidInput = lipgloss.NewStyle().BorderForeground(lipgloss.Color("42")).
		Border(baseBorder).Bold(false).Width(minInputWidth)
	s.InvalidInput = s.ValidInput.BorderForeground(lipgloss.Color("160"))
	s.NeutralInput = s.ValidInput.BorderForeground(lipgloss.Color("240"))

	s.ValidActiveInput = s.ValidInput.Border(lipgloss.ThickBorder())
	s.InvalidActiveInput = s.InvalidInput.Border(lipgloss.ThickBorder())
	s.NeutralActiveInput = s.NeutralInput.Border(lipgloss.ThickBorder())

	return s
}
