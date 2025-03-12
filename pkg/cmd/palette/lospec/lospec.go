package lospec

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/environment"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/utils"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"unicode"
)

type Options struct {
	FolderDestination string
	Format            string
}

func NewPaletteLospecCmd(env *environment.Environment) *cobra.Command {
	opts := &Options{}

	cmd := &cobra.Command{
		Use:     "lospec [ARG] [FLAGS]",
		Aliases: []string{"l", "lp"},
		Short:   "Palette command to import palettes from Lopsec by names",
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := env.Config()
			if err != nil {
				return err
			}

			if opts.FolderDestination == "" {
				utils.PrintlnBold("⚠ no palettes folder specified, choosing first from config")

				if len(cfg.PalettesFolderPaths) == 0 {
					return errors.New("❌ config palettes folders are not exists")
				}
				opts.FolderDestination = cfg.PalettesFolderPaths[0]
			}

			err = ImportLospecPalettes(opts, args)
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&opts.FolderDestination, "destination", "d", "", "Folder to store palette files")
	cmd.Flags().StringVarP(&opts.Format, "format", "f", "gpl", "Format of downloadable palette (gpl, hex, png, pal, ase, txt)")

	return cmd
}

func ImportLospecPalettes(opts *Options, palettesNames []string) error {
	if !slices.Contains(supportedFormats(), opts.Format) {
		return fmt.Errorf("unsupported format: %s", opts.Format)
	}

	for _, paletteName := range palettesNames {
		processedName, err := parsePaletteName(paletteName)
		if err != nil {
			return fmt.Errorf("invalid palette name '%s': %w", paletteName, err)
		}

		stopSpinner := make(chan bool)
		go utils.CreateSpinner("⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏", stopSpinner, fmt.Sprintf("Downloading %s...", paletteName))

		filename := fmt.Sprintf("%s.%s", processedName, opts.Format)
		url := fmt.Sprintf("https://lospec.com/palette-list/%s", filename)
		filePath := filepath.Join(opts.FolderDestination, processedName)
		err = downloadPalette(url, filePath)

		stopSpinner <- true
		<-stopSpinner

		if err != nil {
			utils.PrintError(fmt.Sprintf("Failed download and save %s: %v\n", paletteName, err))
			continue
		}

		utils.PrintlnSuccess(fmt.Sprintf("Downloaded %s ➡ %s\n", paletteName, filePath))
	}

	return nil
}

func parsePaletteName(name string) (string, error) {
	nameSeparator := '-'
	processed := strings.ToLower(name)
	processed = strings.Join(strings.Fields(processed), string(nameSeparator))
	processed = strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsNumber(r) || r == nameSeparator {
			return r
		}
		return -1
	}, processed)

	if len(processed) == 0 {
		return "", errors.New("invalid characters in palette name")
	}

	return processed, nil
}

func downloadPalette(url, path string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	return err
}

func supportedFormats() []string {
	return []string{"gpl", "hex", "png", "pal", "ase", "txt"}
}
