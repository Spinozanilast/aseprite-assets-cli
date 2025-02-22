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
	var inputs []string
	assetsFldIdx, paletteFolderFldIdx := 1, 1
	for i, field := range m.fields {
		label := field.description
		if field.fldType == AssetsFolderPathFld {
			label = fmt.Sprintf("Assets Directory #%d", assetsFldIdx)
			assetsFldIdx++
		} else if field.fldType == PalettesFolderPathFld {
			label = fmt.Sprintf("Palettes Directory #%d", paletteFolderFldIdx)
			paletteFolderFldIdx++
		}
		inputs = append(inputs, m.renderInputField(field, label, i))
	}

	inputContainer := lipgloss.NewStyle().
		Width(m.styles.App.GetWidth() - 4).
		Align(lipgloss.Left)

	mainContent.WriteString(inputContainer.Render(
		lipgloss.JoinVertical(lipgloss.Left, inputs...),
	))

	// Error section
	if len(m.err) > 0 {
		mainContent.WriteString(m.styles.ErrorLabel.Render("Error: ") + m.styles.Error.Render(m.err) + "\n")
	}

	// Help section
	helpView := m.styles.Help.Render("\n" + m.help.View(m.keys))
	mainContent.WriteString(helpView)

	return m.styles.App.Render(mainContent.String())
}

func (m model) renderInputField(field *inputField, label string, idx int) string {
	isActive := idx == m.activeFieldIndex
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
