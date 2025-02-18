package cmd

import (
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"github.com/spinozanilast/aseprite-assets-cli/aseprite"
	config "github.com/spinozanilast/aseprite-assets-cli/config"
)

type AssetCreateOptions struct {
	AssetName  string `survey:"name"`
	Ui         bool
	Width      int
	Height     int
	ColorMode  string `survey:"mode"`
	OutputPath string `survey:"path"`
}

var createCmd = &cobra.Command{
	Use:     "create [ARG]",
	Aliases: []string{"cr"},
	Short:   "Create aseprite asset",
	RunE: func(cmd *cobra.Command, args []string) error {
		config, err := config.LoadConfig()
		if err != nil {
			return err
		}

		asepritePath := config.AsepritePath
		opts, err := collectCreateOptions(config.AssetsFolderPaths)
		if err != nil {
			return err
		}

		asepriteCli := aseprite.NewAsepriteCLI(asepritePath)
		err = asepriteCli.CheckPrerequisites()
		if err != nil {
			return err
		}

		err = asepriteCli.CreateAsset(aseprite.AsepriteAssetCreateCommand{
			Ui:         opts.Ui,
			Width:      opts.Width,
			Height:     opts.Height,
			ColorMode:  opts.ColorMode,
			OutputPath: opts.OutputPath + "\\" + strings.TrimSpace(opts.AssetName) + ".aseprite",
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
