package config

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m model) View() string {
	if m.quitting {
		return "Configuration cancelled\n"
	}

	var mainContent strings.Builder

	// Title section
	title := m.styles.Title.Render("Aseprite-assets CLI Configuration Wizard")
	mainContent.WriteString(lipgloss.PlaceHorizontal(
		m.styles.App.GetWidth(),
		lipgloss.Center,
		title,
	) + "\n\n")

	// Input fields
	inputs := []string{
		m.renderInputField(&m.appPathFld, "Application Path", 0),
	}

	for i := range m.assetsDirsFlds {
		inputs = append(inputs, m.renderInputField(
			&m.assetsDirsFlds[i],
			fmt.Sprintf("Assets Directory #%d", i+1),
			i+1,
		))
	}

	inputContainer := lipgloss.NewStyle().
		Width(m.styles.App.GetWidth() - 4).
		Align(lipgloss.Left)

	mainContent.WriteString(inputContainer.Render(
		lipgloss.JoinVertical(lipgloss.Left, inputs...),
	))

	//Error section
	if len(m.err) > 0 {
		mainContent.WriteString(m.styles.ErrorLabel.Render("Error: ") + m.styles.Error.Render(m.err) + "\n")
	}

	// Help section
	helpView := m.styles.Help.Render("\n" + m.help.View(m.keys))
	mainContent.WriteString(helpView)

	return m.styles.App.Render(mainContent.String())
}

func (m model) renderInputField(field *inputField, label string, idx int) string {
	var isActive bool
	if idx == 0 && m.activeInputIdx == 0 {
		isActive = true
	} else {
		isActive = idx == m.activeInputIdx
	}
	var style lipgloss.Style

	if isActive {
		switch field.status {
		case statusValid:
			style = m.styles.ValidActiveInput
		case statusInvalid:
			style = m.styles.InvalidActiveInput
		default:
			style = m.styles.NeutralActiveInput
		}
	} else {
		switch field.status {
		case statusValid:
			style = m.styles.ValidInput
		case statusInvalid:
			style = m.styles.InvalidInput
		default:
			style = m.styles.NeutralInput
		}
	}

	return style.
		Width(m.styles.App.GetWidth()-6).
		Render(lipgloss.JoinVertical(
			lipgloss.Top,
			m.styles.InputLabel.Render(label+":"),
			field.View(),
		)) + "\n"
}
