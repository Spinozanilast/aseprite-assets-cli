package list

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/aseprite"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/config"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/environment"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/tui"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/utils"

	list "github.com/spinozanilast/aseprite-assets-cli/internal/tui/list"
)

type ListType int

const (
	RecursiveList ListType = 1 << iota
	SpritesList
	PalettesList
	UnknownList = 0
)

type ListOptions struct {
	Config      func() (*config.Config, error)
	Recursive   bool
	SpriteList  bool
	PaletteList bool
}

type listHandler struct {
	config     *config.Config
	listType   ListType
	extensions []string
	folders    []string
}

func NewListCmd(env *environment.Environment) *cobra.Command {
	opts := &ListOptions{
		Config: env.Config,
	}

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"l"},
		Short:   "List existing aseprite assets",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList(opts)
		},
	}

	cmd.Flags().BoolVarP(&opts.Recursive, "recursive", "r", false, "recursive search for selected source")
	cmd.Flags().BoolVarP(&opts.SpriteList, "sprites", "s", false, "list sprites")
	cmd.Flags().BoolVarP(&opts.PaletteList, "palettes", "p", false, "list palettes")

	return cmd
}

func runList(opts *ListOptions) error {
	cfg, err := opts.Config()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	listType, err := opts.validateListType()
	if err != nil {
		return err
	}

	handler := &listHandler{
		config:   cfg,
		listType: listType,
	}

	if err := handler.specifyListParameters(); err != nil {
		return err
	}

	sources, err := handler.findSources()
	if err != nil {
		return err
	}

	if err = tui.StartListTui(WriteTitle(handler.listType), cfg.AsepritePath, sources); err != nil {

	}

	return nil
}

func (opts *ListOptions) validateListType() (ListType, error) {
	if opts.SpriteList && opts.PaletteList {
		return UnknownList, fmt.Errorf("cannot list both sprites and palettes at the same time")
	}

	if !opts.SpriteList && !opts.PaletteList {
		return UnknownList, fmt.Errorf("must specify either sprites or palettes")
	}

	var lt ListType
	if opts.SpriteList {
		lt |= SpritesList
	}
	if opts.PaletteList {
		lt |= PalettesList
	}
	if opts.Recursive {
		lt |= RecursiveList
	}

	return lt, nil
}

func (h *listHandler) specifyListParameters() error {
	switch {
	case h.listType&SpritesList != 0:
		h.folders = h.config.AssetsFolderPaths
		h.extensions = aseprite.SpritesExtensions()
	case h.listType&PalettesList != 0:
		h.folders = h.config.PalettesFolderPaths
		h.extensions = aseprite.PaletteExtensions()
	default:
		return fmt.Errorf("invalid list type configuration")
	}

	if len(h.folders) == 0 {
		return fmt.Errorf("no folder found for list of assets in config for searching assets type")
	}

	return nil
}

func (h *listHandler) findSources() ([]list.AssetsSource, error) {
	// if recursive flag was added
	if h.listType&RecursiveList != 0 {
		return h.findSourcesRecursive()
	}
	return h.findFlatSources()
}

// findSourcesRecursive finds assets recursively (watching inner folders for files)
func (h *listHandler) findSourcesRecursive() ([]list.AssetsSource, error) {
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

// findSourcesRecursive finds assets recursively (don't watching inner folders)
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
