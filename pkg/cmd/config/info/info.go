package info

import (
	"github.com/spf13/cobra"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/config"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/environment"
)

func NewConfigOpenCmd(env *environment.Environment) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "info",
		Short: "Show current configuration",
		Run: func(cmd *cobra.Command, args []string) {
			config.ConfigInfo()
		},
	}

	return cmd
}
