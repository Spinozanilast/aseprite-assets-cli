package assets

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	utils "github.com/spinozanilast/aseprite-assets-cli/util"
)

func (m model) View() string {
	var mainContent strings.Builder

	// Title section
	title := m.styles.Title.Render("Aseprite-assets View")
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

func (m model) renderFolderNavigation() string {
	maxLen := utils.MaxLength(m.prevFolderName, m.activeFolderName, m.nextFolderName)

	sectionWidth := maxLen + 2

	prevSection := m.styles.BeforeAfterFolders.Width(sectionWidth).AlignHorizontal(lipgloss.Left).Render("< " + m.prevFolderName)
	currentSection := m.styles.CurrentFolder.Width(sectionWidth).AlignHorizontal(lipgloss.Center).Render(m.activeFolderName)
	nextSection := m.styles.BeforeAfterFolders.Width(sectionWidth).AlignHorizontal(lipgloss.Right).Render(m.nextFolderName + " >")

	folderNavigation := lipgloss.JoinHorizontal(lipgloss.Center, prevSection, currentSection, nextSection)

	return lipgloss.PlaceHorizontal(m.appWidth, lipgloss.Left, folderNavigation)
}
