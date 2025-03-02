package list

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	var listContent strings.Builder

	listContent.WriteString(m.renderTitle())
	listContent.WriteString(m.renderNavigation() + "\n")
	listContent.WriteString(m.renderList())
	listContent.WriteString(m.renderError())
	listContent.WriteString(m.renderHelp())

	return lipgloss.JoinHorizontal(lipgloss.Top,
		m.styles.App.Render(listContent.String()),
		m.infoPanel.View(),
	)
}

func (m Model) renderTitle() string {
	title := m.styles.Title.Render(m.title)
	return lipgloss.PlaceHorizontal(
		m.styles.App.GetWidth(),
		lipgloss.Center,
		title,
	) + "\n"
}

func (m Model) renderNavigation() string {
	sectionWidth := m.windowWidth/6 - 1
	remainder := m.windowWidth % 6

	return lipgloss.JoinHorizontal(lipgloss.Top,
		m.renderFolderSection(sectionWidth, m.nav.prev, lipgloss.Left),
		m.renderFolderSection(sectionWidth+remainder, "<-- "+m.nav.active+" -->", lipgloss.Center),
		m.renderFolderSection(sectionWidth, m.nav.next, lipgloss.Right),
	)
}

func (m Model) renderFolderSection(width int, content string, align lipgloss.Position) string {
	return m.styles.BeforeAfterFolders.
		Width(width).
		MaxWidth(width).
		Align(align).
		Render(content)
}

func (m Model) renderList() string {
	m.list.SetShowTitle(false)
	return m.list.View() + "\n"
}

func (m Model) renderError() string {
	if len(m.err) > 0 {
		return m.styles.ErrorLabel.Render("Error: ") +
			m.styles.Error.Render(m.err) + "\n"
	}
	return ""
}

func (m Model) renderHelp() string {
	return m.styles.Help.Render("\n" + m.help.View(m.keys))
}
