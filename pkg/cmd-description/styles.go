package cmddescription

import (
	"github.com/charmbracelet/lipgloss"
)

type DescriptionStyles struct {
	Title               lipgloss.Style
	SectionTitle        lipgloss.Style
	SectionText         lipgloss.Style
	ListTitle           lipgloss.Style
	ListItem            lipgloss.Style
	ListItemDescription lipgloss.Style
	ListSubItem         lipgloss.Style
}

func DefaultStyles() *DescriptionStyles {
	return &DescriptionStyles{
		Title: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#00ff00")).
			MarginTop(1).
			MarginBottom(1),

		SectionTitle: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#ffff00")).
			MarginBottom(1).
			MarginTop(1),

		SectionText: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ffffff")),
		ListTitle: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#ffff00")).
			MarginBottom(1).
			MarginTop(1),

		ListItem: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#aaaaaa")),

		ListItemDescription: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#555555")),
		ListSubItem: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#aaaaaa")).
			BorderStyle(lipgloss.Border{
				Left: "│",
			}).
			BorderForeground(lipgloss.Color("#444444")),
	}
}

func (s *DescriptionStyles) WithColors(titleColor, textColor, sectionColor, listColor, subtextColor string) *DescriptionStyles {
	newStyles := *s
	newStyles.Title = newStyles.Title.Foreground(lipgloss.Color(titleColor))
	newStyles.SectionTitle = newStyles.SectionTitle.Foreground(lipgloss.Color(sectionColor))
	newStyles.SectionText = newStyles.SectionText.Foreground(lipgloss.Color(textColor))
	newStyles.ListTitle = newStyles.ListTitle.Foreground(lipgloss.Color(sectionColor))
	newStyles.ListItem = newStyles.ListItem.Foreground(lipgloss.Color(listColor))
	newStyles.ListSubItem = newStyles.ListSubItem.Foreground(lipgloss.Color(subtextColor))
	newStyles.ListItemDescription = newStyles.ListItemDescription.Foreground(lipgloss.Color(subtextColor))
	return &newStyles
}

func (s *DescriptionStyles) WithMargins(vertical, horizontal int) *DescriptionStyles {
	newStyles := *s
	newStyles.Title = newStyles.Title.Margin(vertical, horizontal)
	newStyles.SectionTitle = newStyles.SectionTitle.Margin(vertical, horizontal)
	newStyles.SectionText = newStyles.SectionText.Margin(vertical, horizontal)
	newStyles.ListTitle = newStyles.ListTitle.Margin(vertical, horizontal)
	newStyles.ListItem = newStyles.ListItem.Margin(vertical, horizontal)
	newStyles.ListSubItem = newStyles.ListSubItem.Margin(vertical, horizontal)
	return &newStyles
}
