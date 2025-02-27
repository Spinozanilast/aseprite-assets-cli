package commands

import (
	"github.com/spf13/cobra"
	cmddescription "github.com/spinozanilast/aseprite-assets-cli/pkg/cmd-description"
)

type CommandDescription struct {
	Title    string
	Short    string
	Long     string
	Examples []string
	Sections []Section
	Lists    []List
}

type Section struct {
	Title string
	Text  string
}

type List struct {
	Title      string
	Indent     int
	IndentChar string
	Items      []string
}

// ApplyToCommand applies the description to a cobra command
func (cd *CommandDescription) ApplyToCommand(cmd *cobra.Command) {
	builder := cmddescription.NewBuilder().
		WithStyles(cmddescription.DefaultStyles()).
		WithTitle(cd.Title).
		WithShort(cd.Short).
		WithLong(cd.Long)

	for _, section := range cd.Sections {
		builder.WithSection(section.Title, section.Text)
	}

	for _, list := range cd.Lists {
		builder.WithList(list.Title, list.Indent, list.IndentChar, list.Items...)
	}

	description := builder.Build()
	cmd.Long = description.String()
	cmd.Short = description.Short

	for _, example := range cd.Examples {
		cmd.Example += example + "\n"
	}
}
