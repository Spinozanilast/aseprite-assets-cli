package commands

// import (
// 	"fmt"

// 	"github.com/spf13/cobra"
// 	config "github.com/spinozanilast/aseprite-assets-cli/config"
// )

// var removeCmd = &cobra.Command{
// 	Use:     "remove",
// 	Aliases: []string{"r", "rm"},
// 	Short:   "Remove existing aseprite assets",
// 	RunE: func(cmd *cobra.Command, args []string) error {
// 		config, err := config.LoadConfig()
// 		if err != nil {
// 			return err
// 		}

// 		// TODO: remove assets with autocomplete and confirmation

// 		if len(args) == 0 {
// 			return fmt.Errorf("no assets specified")
// 		}
// 	},
// }

// func init() {
// 	rootCmd.AddCommand(removeCmd)
// }
