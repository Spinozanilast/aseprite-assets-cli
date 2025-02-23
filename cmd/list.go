package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spinozanilast/aseprite-assets-cli/utils"

	tea "github.com/charmbracelet/bubbletea"
	config "github.com/spinozanilast/aseprite-assets-cli/config"
	list "github.com/spinozanilast/aseprite-assets-cli/tui/list"
)

type ListType int

const (
	RecursiveList ListType = 1 << iota
	SpritesList
	PalettesList
	UnknownList = 0
)

type listHandler struct {
	listTitle  string
	config     *config.Config
	listType   ListType
	extensions []string
	folders    []string
}

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	Short:   "List existing aseprite assets",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.LoadConfig()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		handler, err := setupListHandler(cmd, cfg)
		if err != nil {
			return err
		}

		sources, err := handler.findAssetSources()
		if err != nil {
			return err
		}

		handler.startTui(sources)
		return nil
	},
}

func init() {
	listCmd.Flags().BoolP("recursive", "r", false, "recursive search for assets")
	listCmd.Flags().BoolP("sprites", "s", false, "assets extension")
	listCmd.Flags().BoolP("palettes", "p", false, "assets extension")
	rootCmd.AddCommand(listCmd)
}

func setupListHandler(cmd *cobra.Command, cfg *config.Config) (*listHandler, error) {
	recursive, _ := cmd.Flags().GetBool("recursive")
	sprites, _ := cmd.Flags().GetBool("sprites")
	palettes, _ := cmd.Flags().GetBool("palettes")

	listType, err := validateListType(sprites, palettes, recursive)
	if err != nil {
		return nil, err
	}

	handler := &listHandler{
		config:   cfg,
		listType: listType,
	}

	if err := handler.setAssetParameters(); err != nil {
		return nil, err
	}

	return handler, nil
}

func validateListType(sprites, palettes, recursive bool) (ListType, error) {
	if sprites && palettes {
		return UnknownList, fmt.Errorf("cannot list both sprites and palettes at the same time")
	}

	var lt ListType
	if sprites {
		lt |= SpritesList
	}
	if palettes {
		lt |= PalettesList
	}
	if recursive {
		lt |= RecursiveList
	}

	if lt&(SpritesList|PalettesList) == 0 {
		return UnknownList, fmt.Errorf("must specify either sprites or palettes")
	}

	return lt, nil
}

func (h *listHandler) startTui(sources []list.AssetsSource) {
	title := WriteTitle(h.listType)
	p := tea.NewProgram(list.InitialModel(title, h.config.AsepritePath, sources))
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Oopsy, something went wrong and cli is not executing properly: '%s'", err)
		os.Exit(1)
	}
}

func (h *listHandler) setAssetParameters() error {
	switch {
	case h.listType&SpritesList != 0:
		h.folders = h.config.AssetsFolderPaths
		h.extensions = []string{".aseprite"}
	case h.listType&PalettesList != 0:
		h.folders = h.config.PalettesFolderPaths
		h.extensions = []string{".gpl", ".png"}
	default:
		return fmt.Errorf("invalid list type configuration")
	}
	return nil
}

func (h *listHandler) findAssetSources() ([]list.AssetsSource, error) {
	if h.listType&RecursiveList != 0 {
		return h.findRecursiveSources()
	}
	return h.findFlatSources()
}

func (h *listHandler) findRecursiveSources() ([]list.AssetsSource, error) {
	var sources []list.AssetsSource

	for _, dir := range h.folders {
		if !utils.СheckFileExists(dir, true) {
			return nil, fmt.Errorf("directory not found: %s", dir)
		}

		files, err := utils.FindFilesOfExtensionsRecursive(dir, h.extensions...)
		if err != nil {
			return nil, fmt.Errorf("search failed in %s: %w", dir, err)
		}

		for folder, paths := range files {
			sources = append(sources, list.NewAssetsSource(folder, paths))
		}
	}

	return sources, nil
}

func (h *listHandler) findFlatSources() ([]list.AssetsSource, error) {
	var sources []list.AssetsSource

	for _, dir := range h.folders {
		if !utils.СheckFileExists(dir, true) {
			return nil, fmt.Errorf("directory not found: %s", dir)
		}

		files, err := utils.FindFilesOfExtensions(dir, h.extensions...)
		if err != nil {
			return nil, fmt.Errorf("search failed in %s: %w", dir, err)
		}

		if len(files) > 0 {
			sources = append(sources, list.NewAssetsSource(dir, files))
		}
	}

	return sources, nil
}

func WriteTitle(t ListType) (title string) {
	if t&SpritesList != 0 {
		title = "🖌️ List of sprites"
	} else if t&PalettesList != 0 {
		title = "🎨 List of palettes"
	}

	if t&RecursiveList != 0 {
		title += " | 🔍 Recursive"
	}

	return title
}
