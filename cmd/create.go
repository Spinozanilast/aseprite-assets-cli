package cmd

import (
	"errors"
	"fmt"
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

		if err := validateCreateOptions(opts); err != nil {
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

		showSummary(opts)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().StringVarP(&assetName, "name", "n", "", "name of the asset")
	createCmd.Flags().BoolVarP(&ui, "ui", "u", false, "open aseprite after asset creation")
	createCmd.Flags().IntVarP(&width, "width", "w", 32, "width of the asset")
	createCmd.Flags().IntVar(&height, "height", 32, "height of the asset") // Changed shorthand to 't'
	createCmd.Flags().StringVarP(&colorMode, "mode", "m", "rgb", "color mode of the asset (indexed, rgb, gray, tilemap)")
	createCmd.Flags().StringVarP(&outputPath, "path", "p", "", "output path for the asset")
}

func createAssetQuestions(dirs []string) []*survey.Question {
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
			},
		},
		{
			Name: "height",
			Prompt: &survey.Input{
				Message: "Height",
				Default: "32",
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
		return nil, err
	}

	return opts, nil
}

func validateCreateOptions(opts *AssetCreateOptions) error {
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

func showSummary(opts *AssetCreateOptions) {
	fmt.Printf("Asset configuration summary:\n")
	fmt.Printf("Name: %v\n", opts.AssetName)
	fmt.Printf("UI: %v\n", opts.Ui)
	fmt.Printf("Width: %v\n", opts.Width)
	fmt.Printf("Height: %v\n", opts.Height)
	fmt.Printf("Color mode: %v\n", opts.ColorMode)
	fmt.Printf("Output path: %v\n", opts.OutputPath)
	fmt.Println("✓ Asset created successfully")
}
