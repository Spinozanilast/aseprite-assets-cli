package commands

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/aseprite"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/aseprite/commands"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/utils"

	config "github.com/spinozanilast/aseprite-assets-cli/pkg/config"
)

type AssetCreateOptions struct {
	AssetName  string `survey:"name"`
	Ui         bool
	Width      int
	Height     int
	ColorMode  string `survey:"mode"`
	OutputPath string `survey:"path"`
}

type assetHandler struct {
	config *config.Config
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

		h := &assetHandler{
			config: config,
		}

		opts, err := h.collectCreateOptions()
		if err != nil {
			fatalError("failed to collect create options: %w", err)
		}

		if err := h.createAsset(opts); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}

func (h *assetHandler) createAsset(opts *AssetCreateOptions) error {
	asepriteCli := aseprite.NewCLI(h.config.AsepritePath, h.config.ScriptDirPath)
	err := asepriteCli.CheckPrerequisites()
	if err != nil {
		return err
	}

	filename := filepath.Join(opts.OutputPath, strings.TrimSpace(opts.AssetName)+aseprite.Aseprite.String())

	if utils.СheckFileExists(filename, false) {
		return fmt.Errorf("file already exists: %s", filename)
	}

	err = asepriteCli.ExecuteCommand(&commands.AssetCreate{
		BatchMode:  !opts.Ui,
		Width:      opts.Width,
		Height:     opts.Height,
		ColorMode:  opts.ColorMode,
		OutputPath: filename,
	})

	if err != nil {
		return err
	}

	showSummary(opts)

	return nil
}

func createAssetQuestions(dirs []string) []*survey.Question {
	questions := []*survey.Question{
		{
			Name: "name",
			Prompt: &survey.Input{
				Message: "Asset name (without extension)",
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
				Options: aseprite.ColorModes(),
				Default: aseprite.ColorModeRGB.String(),
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

func (h *assetHandler) collectCreateOptions() (*AssetCreateOptions, error) {
	opts := &AssetCreateOptions{}
	err := survey.Ask(createAssetQuestions(h.config.AssetsFolderPaths), opts)

	if err != nil {
		return nil, err
	}

	return opts, nil
}

func showSummary(opts *AssetCreateOptions) {
	utils.PrintlnBold("\nAsset configuration summary:\n")
	fmt.Printf("Name: %v\n", opts.AssetName)
	fmt.Printf("UI: %v\n", opts.Ui)
	fmt.Printf("Width: %v\n", opts.Width)
	fmt.Printf("Height: %v\n", opts.Height)
	fmt.Printf("Color mode: %v\n", opts.ColorMode)
	fmt.Printf("Output path: %v\n", opts.OutputPath)
	utils.PrintlnSuccess("✓ Asset created successfully")
}
