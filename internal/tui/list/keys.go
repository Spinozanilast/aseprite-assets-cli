package list

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Left  key.Binding
	Right key.Binding
	Enter key.Binding
	Tab   key.Binding
	Help  key.Binding
	Quit  key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Left, k.Right, k.Tab, k.Enter, k.Help, k.Quit},
	}
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
		key.WithHelp("Enter ", "confirm or browse file/directory"),
	),
	Help: key.NewBinding(
		key.WithKeys("ctrl+h"),
		key.WithHelp("Ctrl+H", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("esc", "ctrl+c"),
		key.WithHelp("ESC/Ctrl+C", "quit"),
	),
}
