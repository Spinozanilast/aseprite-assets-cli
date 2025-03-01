package assets

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	var mainContent strings.Builder

	// Title section
	title := m.styles.Title.Render(m.title)
	mainContent.WriteString(lipgloss.PlaceHorizontal(
		m.styles.App.GetWidth(),
		lipgloss.Center,
		title,
	) + "\n")

	// Folders section
	mainContent.WriteString(m.renderFolderNavigation() + "\n")

	// List section
	m.list.SetShowTitle(false)
	mainContent.WriteString(m.list.View() + "\n")

	//Error section
	if len(m.err) > 0 {
		mainContent.WriteString(m.styles.ErrorLabel.Render("Error: ") + m.styles.Error.Render(m.err) + "\n")
	}

	// Help section
	helpView := m.styles.Help.Render("\n" + m.help.View(m.keys))
	mainContent.WriteString(helpView)

	return m.styles.App.Render(mainContent.String())
}

func (m Model) renderFolderNavigation() string {
	sectionWidth := m.appWidth / 3
	remainder := m.appWidth % 3

	prevWidth := sectionWidth
	currentWidth := sectionWidth + remainder
	nextWidth := sectionWidth

	prevSection := m.styles.BeforeAfterFolders.
		Width(prevWidth).
		MaxWidth(prevWidth).
		Align(lipgloss.Left).
		Render("< " + m.prevFolderName)

	currentSection := m.styles.CurrentFolder.
		Width(currentWidth).
		MaxWidth(currentWidth).
		Align(lipgloss.Center).
		Render(m.activeFolderName)

	nextSection := m.styles.BeforeAfterFolders.
		Width(nextWidth).
		MaxWidth(nextWidth).
		Align(lipgloss.Right).
		Render(m.nextFolderName + " >")

	return lipgloss.JoinHorizontal(lipgloss.Top, prevSection, currentSection, nextSection)
}
