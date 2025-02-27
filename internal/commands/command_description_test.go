package commands

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestCommandDescription(t *testing.T) {
	t.Run("ApplyToCommand", func(t *testing.T) {
		cmd := &cobra.Command{}

		cd := &CommandDescription{
			Title: "Test Command",
			Short: "A test command",
			Long:  "A longer description of the test command",
			Examples: []string{
				"example1",
				"example2",
			},
			Sections: []Section{
				{
					Title: "Section 1",
					Text:  "Section 1 text",
				},
				{
					Title: "Section 2",
					Text:  "Section 2 text",
				},
			},
			Lists: []List{
				{
					Title:      "List 1",
					Indent:     2,
					IndentChar: "-",
					Items:      []string{"item1", "item2"},
				},
			},
		}

		cd.ApplyToCommand(cmd)

		assert.Contains(t, cmd.Long, cd.Title)
		assert.Equal(t, cd.Short, cmd.Short)

		for _, section := range cd.Sections {
			assert.Contains(t, cmd.Long, section.Title)
			assert.Contains(t, cmd.Long, section.Text)
		}

		// Verify examples are included
		for _, example := range cd.Examples {
			assert.Contains(t, cmd.Example, example)
		}
	})

	t.Run("EmptyDescription", func(t *testing.T) {
		cmd := &cobra.Command{}

		cd := &CommandDescription{}

		cd.ApplyToCommand(cmd)

		assert.NotPanics(t, func() {
			cd.ApplyToCommand(cmd)
		})
	})
}
