package cmd

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"github.com/spinozanilast/aseprite-assets-cli/aseprite"
	config "github.com/spinozanilast/aseprite-assets-cli/config"
	"github.com/spinozanilast/aseprite-assets-cli/util"
)

type AssetCreateOptions struct {
	AssetName  string `survey:"name"`
	Ui         bool
	Width      int
	Height     int
	ColorMode  string `survey:"mode"`
	OutputPath string `survey:"path"`
}

var (
	assetName  string
	ui         bool
	width      int
	height     int
	colorMode  string
	outputPath string
)

var createCmd = &cobra.Command{
	Use:     "create [ARG]",
	Aliases: []string{"cr"},
	Short:   "Create aseprite asset",
	Long: `Create a new aseprite asset with the specified options.

Available arguments:
  name       - The name of the asset.
  ui         - Whether to open aseprite after asset creation.
  width      - The width of the asset.
  height     - The height of the asset.
  mode       - The color mode of the asset (indexed, rgb, gray, tilemap).
  path       - The output path for the asset.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		config, err := config.LoadConfig()
		if err != nil {
			return err
		}

		opts := &AssetCreateOptions{
			AssetName:  assetName,
			Ui:         ui,
			Width:      width,
			Height:     height,
			ColorMode:  colorMode,
			OutputPath: outputPath,
		}

		if opts.AssetName == "" || opts.OutputPath == "" {
			opts, err = collectCreateOptions(config.AssetsFolderPaths)
			if err != nil {
				return err
			}
		}

		if err := opts.validateOptions(); err != nil {
			return err
		}

		asepriteCli := aseprite.NewAsepriteCLI(config.AsepritePath, config.ScriptDirPath)
		err = asepriteCli.CheckPrerequisites()
		if err != nil {
			return err
		}

		err = asepriteCli.ExecuteCommand(&aseprite.AsepriteAssetCreateCommand{
			Ui:         opts.Ui,
			Width:      opts.Width,
			Height:     opts.Height,
			ColorMode:  opts.ColorMode,
			OutputPath: opts.OutputPath + "\\" + util.EnsureFileExtension(strings.TrimSpace(opts.AssetName), aseprite.AsepriteFilesExtension),
		})

		if err != nil {
			return err
		}

		opts.showSummary()

		return nil
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().StringVarP(&assetName, "name", "n", "", "Name of the asset")
	createCmd.Flags().BoolVarP(&ui, "ui", "u", false, "Open aseprite after asset creation")
	createCmd.Flags().IntVarP(&width, "width", "w", 32, "Width of the asset")
	createCmd.Flags().IntVar(&height, "height", 32, "Height of the asset") // Changed shorthand to 't'
	createCmd.Flags().StringVarP(&colorMode, "mode", "m", "rgb", "Color mode of the asset (indexed, rgb, gray, tilemap)")
	createCmd.Flags().StringVarP(&outputPath, "path", "p", "", "Output path for the asset")
}

func createAssetQuestions(dirs []string) []*survey.Question {
	sizesSuggestions := func(toComplete string) []string {
		num, _ := strconv.Atoi(toComplete)
		if toComplete == "" || (num <= 0 && num >= 10) {
			return []string{"16", "32", "64", "128", "256", "512", "1024", "Custom"}
		}

		floated := float64(num)
		var suggestions []string
		for i := 1; i <= 7; i++ {
			suggestion := fmt.Sprintf("%d", int(math.Pow(floated, float64(i))))
			suggestions = append(suggestions, suggestion)
		}
		return suggestions
	}

	questions := []*survey.Question{
		{
			Name: "name",
			Prompt: &survey.Input{
				Message: "Asset name",
				Default: "asset",
			},
		},
		{
			Name: "ui",
			Prompt: &survey.Confirm{
				Message: "Open aseprite after asset creation?",
				Default: false,
			},
		},
		{
			Name: "width",
			Prompt: &survey.Input{
				Message: "Width",
				Default: "32",
				Suggest: sizesSuggestions,
			},
		},
		{
			Name: "height",
			Prompt: &survey.Input{
				Message: "Height",
				Default: "32",
				Suggest: sizesSuggestions,
			},
		},
		{
			Name: "mode",
			Prompt: &survey.Select{
				Message: "Color mode",
				Options: []string{"indexed", "rgb", "gray", "tilemap"},
				Default: "rgb",
			},
		},
		{
			Name: "path",
			Prompt: &survey.Select{
				Message: "Output path",
				Options: dirs,
				Default: dirs[0],
			},
		},
	}
	return questions
}

func collectCreateOptions(saveDirs []string) (*AssetCreateOptions, error) {
	opts := &AssetCreateOptions{}
	err := survey.Ask(createAssetQuestions(saveDirs), opts)

	if err != nil {
		return nil, fmt.Errorf("failed to collect create options: %w", err)
	}

	return opts, nil
}

func (opts *AssetCreateOptions) validateOptions() error {
	if opts.AssetName == "" {
		return errors.New("asset name cannot be empty")
	}
	if opts.Width <= 0 {
		return errors.New("width must be greater than 0")
	}
	if opts.Height <= 0 {
		return errors.New("height must be greater than 0")
	}
	if opts.OutputPath == "" {
		return errors.New("output path cannot be empty")
	}
	return nil
}

func (opts *AssetCreateOptions) showSummary() {
	fmt.Printf("\nAsset configuration summary:\n")
	fmt.Printf("Name: %v\n", opts.AssetName)
	fmt.Printf("UI: %v\n", opts.Ui)
	fmt.Printf("Width: %v\n", opts.Width)
	fmt.Printf("Height: %v\n", opts.Height)
	fmt.Printf("Color mode: %v\n", opts.ColorMode)
	fmt.Printf("Output path: %v\n", opts.OutputPath)
	fmt.Println("\n✓ Asset created successfully")
}
