package main

import (
	"fmt"
	"os"

	cmd "github.com/spinozanilast/aseprite-assets-cli/cmd"
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
		cmd.StartConfigInitializationTui(config)
		return
	}

	cmd.Execute()
}
