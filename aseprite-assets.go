package main

import (
	"fmt"
	"os"

	"github.com/spinozanilast/aseprite-assets-cli/internal/commands"

	config "github.com/spinozanilast/aseprite-assets-cli/pkg/config"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Config error: %v\n", err)
		os.Exit(1)
	}

	if err := config.Validate(); err != nil {
		fmt.Printf("Invalid or empty config: %v\n", err)
		commands.StartConfigInitializationTui(config)
		return
	}

	commands.Execute()
}
