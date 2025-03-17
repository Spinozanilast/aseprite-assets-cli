package list

import "github.com/charmbracelet/bubbles/key"

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

type keyMap struct {
	Left      key.Binding
	Right     key.Binding
	Enter     key.Binding
	Tab       key.Binding
	Help      key.Binding
	Quit      key.Binding
	showEnter bool // Controls visibility of Enter in help
}

// FullHelp to conditionally include Enter
func (k keyMap) FullHelp() [][]key.Binding {
	row := []key.Binding{k.Left, k.Right, k.Tab}
	if k.showEnter {
		row = append(row, k.Enter)
	}
	row = append(row, k.Help, k.Quit)
	return [][]key.Binding{row}
}

var keys = keyMap{
	Right: key.NewBinding(
		key.WithKeys("right"),
		key.WithHelp("→", "move to right directory"),
	),
	Left: key.NewBinding(
		key.WithKeys("left"),
		key.WithHelp("←", "move to left directory"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("Enter ", "open selected asset with default app"),
	),
	Help: key.NewBinding(
		key.WithKeys("ctrl+h"),
		key.WithHelp("Ctrl+H", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("esc", "ctrl+c"),
		key.WithHelp("ESC/Ctrl+C", "quit"),
	),
	showEnter: true, // Default to showing Enter
}
