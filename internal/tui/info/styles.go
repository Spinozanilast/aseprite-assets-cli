package info

import "github.com/charmbracelet/lipgloss"

type Styles struct {
	Base      lipgloss.Style
	Title     lipgloss.Style
	Label     lipgloss.Style
	Value     lipgloss.Style
	Error     lipgloss.Style
	NoContent lipgloss.Style
}

func DefaultStyles() *Styles {
	s := new(Styles)

	s.Base = lipgloss.NewStyle().Padding(1, 2).Border(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("63"))
	s.Title = lipgloss.NewStyle().Foreground(lipgloss.Color("63")).MarginBottom(1).Bold(true)
	s.Label = lipgloss.NewStyle().Foreground(lipgloss.Color("243")).Width(12)
	s.Value = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	s.Error = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)
	s.NoContent = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Italic(true)

	return s
}
