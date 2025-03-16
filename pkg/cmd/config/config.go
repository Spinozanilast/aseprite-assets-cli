package config

import (
	"github.com/spf13/cobra"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/cmd/config/edit"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/cmd/config/info"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/cmd/config/open"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/environment"
)

func NewConfigCmd(env *environment.Environment) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "config [command]",
		Aliases: []string{"cfg"},
		Short:   "Manage configuration settings",
		Long: `Configure various settings for the aseprite-assets CLI tool.
Subcommands allow you to:
- View current configuration (info)
- Modify settings interactively or via flags (edit)
- Open config file in your editor (open)
`,
	}

	cmd.AddCommand(
		edit.NewConfigEditCmd(env),
		info.NewConfigInfoCmd(env),
		open.NewConfigOpenCmd(env),
	)

	return cmd
}
