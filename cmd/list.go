package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	config "github.com/spinozanilast/aseprite-assets-cli/config"
	assetsTui "github.com/spinozanilast/aseprite-assets-cli/tui/assets"
	utils "github.com/spinozanilast/aseprite-assets-cli/util"
	types "github.com/spinozanilast/aseprite-assets-cli/util/types"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	Short:   "List existing aseprite assets",
	RunE: func(cmd *cobra.Command, args []string) error {
		config, err := config.LoadConfig()
		if err != nil {
			return err
		}

		StartListAssetsTui(config.AsepritePath, findAssets(config.AssetsFolderPaths))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func StartListAssetsTui(appPath string, assetsFolders []types.AssetsDirs) {
	p := tea.NewProgram(assetsTui.InitialModel(appPath, assetsFolders))
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Oopsy, something went wrong and cli is not executing properly: '%s'", err)
		os.Exit(1)
	}
}

func findAssets(folderPaths []string) []types.AssetsDirs {
	var assetsDirs []types.AssetsDirs
	for _, folderPath := range folderPaths {
		assetsPaths, err := utils.FindFilesOfExtension(folderPath, ".aseprite")
		if err == nil {
			assetsDirs = append(assetsDirs, types.NewAssetsDirs(folderPath, assetsPaths))
		}
	}

	return assetsDirs
}
