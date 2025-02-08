package config

import (
	"github.com/charmbracelet/lipgloss"
)

type Styles struct {
	App                lipgloss.Style
	Title              lipgloss.Style
	CurrentFolder      lipgloss.Style
	BeforeAfterFolders lipgloss.Style
	Error              lipgloss.Style
	ErrorLabel         lipgloss.Style
	Help               lipgloss.Style
}

func DefaultStyles() *Styles {
	s := new(Styles)

	s.App = lipgloss.NewStyle().Padding(1, 2)
	s.Title = lipgloss.NewStyle().
		Bold(true).Foreground(lipgloss.Color("63")).
		MarginBottom(1)

	s.ErrorLabel = lipgloss.NewStyle().
		Bold(true).Foreground(lipgloss.Color("88"))
	s.Error = lipgloss.NewStyle().
		Background(lipgloss.Color("240")).Foreground(lipgloss.Color("160"))

	s.Help = lipgloss.NewStyle().Bold(true)

	s.CurrentFolder = lipgloss.NewStyle().
		Bold(true).Foreground(lipgloss.Color("61")).Background(lipgloss.Color("240"))
	s.BeforeAfterFolders = s.CurrentFolder.
		Bold(false).Foreground(lipgloss.Color("60")).UnsetBackground()

	return s
}

type ListItemStyles struct {
	ItemStyle         lipgloss.Style
	SelectedItemStyle lipgloss.Style
}

func DefaultListItemStyles() *ListItemStyles {
	baseStyle := lipgloss.NewStyle().PaddingLeft(1)
	return &ListItemStyles{
		ItemStyle:         baseStyle,
		SelectedItemStyle: baseStyle.Foreground(lipgloss.Color("170")).Bold(true),
	}
}

type ListLayoutStyles struct {
	ListHeight int
	ListWidth  int
}

func DefaultListLayoutStyles() *ListLayoutStyles {
	return &ListLayoutStyles{
		ListHeight: 14,
		ListWidth:  20,
	}
}
