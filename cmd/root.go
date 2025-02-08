package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "aseprite-assets-cli",
	Short: "A CLI tool to manage aseprite assets",
	Long: `A CLI tool to manage aseprite assets. 
			This tool allows you to manage aseprite assets by listing, adding, removing and renaming them.`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Oopsy, something went wrong and cli is not executing properly: '%s'", err)
		os.Exit(1)
	}
}
