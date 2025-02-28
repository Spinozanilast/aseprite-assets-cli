package cli

import (
	"fmt"
	"os"

	"github.com/spinozanilast/aseprite-assets-cli/pkg/cmd/root"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/config"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/environment"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/tui"
)

func Main() {
	env := environment.NewEnvironment(config.LoadConfig)

	if cfg, err := env.Config(); err != nil {
		ExitDueToError(err)
	} else {
		if err := cfg.Validate(); err != nil {
			fmt.Printf("Invalid or empty config: %v\n", err)
			if err = tui.StartConfigTui(cfg); err != nil {
				ExitDueToError(err)
			}
			ExitOk()
		}
	}

	rootCmd := root.NewRootCmd(&env)
	err := rootCmd.Execute()

	if err != nil {
		ExitDueToError(err)
	}
}

func ExitDueToError(err error) {
	fmt.Printf("CLI execution pipeline error: %v\n", err)
	os.Exit(1)
}

func ExitOk() {
	os.Exit(0)
}
