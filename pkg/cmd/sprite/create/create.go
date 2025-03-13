package create

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/aseprite"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/aseprite/commands"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/config"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/environment"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/utils"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/utils/files"
)

type SpriteCreateOptions struct {
	AssetName         string `survey:"name"`
	OpenAfterCreation bool   `survey:"ui"`
	Width             int
	Height            int
	ColorMode         string `survey:"mode"`
	OutputPath        string `survey:"path"`
}

type spriteCreationHandler struct {
	config         *config.Config
	outputFilename string
}

func NewSpriteCreateCmd(env *environment.Environment) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Aliases: []string{"c", "cr"},
		Short:   "Create aseprite sprite",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := env.Config()
			if err != nil {
				return err
			}

			h := &spriteCreationHandler{
				config: cfg,
			}

			opts, err := h.collectCreateOptions()
			if err != nil {

				return fmt.Errorf("failed to collect sprite options: %w", err)
			}

			if err := h.createAsset(opts); err != nil {
				return err
			}

			if files.CheckFileExists(h.outputFilename, false) {
				showSummary(opts)
			} else {
				showFailure()
			}

			return nil
		}}
	return cmd
}

func (h *spriteCreationHandler) createAsset(opts *SpriteCreateOptions) error {
	aseCli := aseprite.NewCLI(h.config.AsepritePath, h.config.ScriptDirPath)
	err := aseCli.CheckPrerequisites()
	if err != nil {
		return err
	}

	filename := filepath.Join(opts.OutputPath, strings.TrimSpace(opts.AssetName)+aseprite.Aseprite.String())

	if files.CheckFileExists(filename, false) {
		return fmt.Errorf("file already exists: %s", filename)
	}

	aseCommand := &commands.CreateSprite{
		OpenAfterCreation: opts.OpenAfterCreation,
		Width:             opts.Width,
		Height:            opts.Height,
		ColorMode:         opts.ColorMode,
		OutputPath:        filename,
	}

	err = aseCli.ExecuteCommand(aseCommand)

	// Callback that opens or not file after creation
	openingCallback := aseCommand.ScriptCallback(h.config.AsepritePath)
	openingCallback()

	h.outputFilename = filename

	if err != nil {
		return err
	}

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

func (h *spriteCreationHandler) collectCreateOptions() (*SpriteCreateOptions, error) {
	opts := &SpriteCreateOptions{}
	err := survey.Ask(createAssetQuestions(h.config.SpriteFolderPaths), opts)

	if err != nil {
		return nil, err
	}

	return opts, nil
}

func showSummary(opts *SpriteCreateOptions) {
	utils.PrintlnBold("\nAsset configuration summary:\n")
	fmt.Printf("Name: %v\n", opts.AssetName)
	fmt.Printf("UI: %v\n", opts.OpenAfterCreation)
	fmt.Printf("Width: %v\n", opts.Width)
	fmt.Printf("Height: %v\n", opts.Height)
	fmt.Printf("Color mode: %v\n", opts.ColorMode)
	fmt.Printf("Output path: %v\n", opts.OutputPath)
	utils.PrintlnSuccess("✓ Asset created successfully")
}

func showFailure() {
	utils.PrintError("❌ Asset creation is failed")
}
