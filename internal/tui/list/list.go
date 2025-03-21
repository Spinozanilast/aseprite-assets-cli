package list

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"

	tea "github.com/charmbracelet/bubbletea"
)

type item string

type itemDelegate struct{}

func (i item) FilterValue() string { return string(i) }

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	listStyles := DefaultListItemStyles()
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := listStyles.ItemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return listStyles.SelectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

func createListItems(texts []string) []list.Item {
	items := make([]list.Item, len(texts))
	for i, text := range texts {
		items[i] = item(text)
	}

	return items
}
