package commands

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/aseprite"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/aseprite/commands"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/config"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/utils"
)

type SpriteScaleMode string

const (
	ScaledMode SpriteScaleMode = "scaled"
	SizedMode  SpriteScaleMode = "sized"
)

type exportHandler struct {
	config      *config.Config
	options     *exportOptions
	asepriteCli *aseprite.AsepriteCLI
}

type exportOptions struct {
	spriteFilename string
	outputFilename string
	format         string
	sizes          string
	scales         string
}

var exportCmd = &cobra.Command{
	Use:     "export [ARG]",
	Aliases: []string{"exp", "e"},
	Example: heredoc.Doc(`
	# Export aseprite asset to png format
	aseprite-assets export <asset-filename> --format png

	# Export aseprite asset to png format and save it to the specified path
	aseprite-assets export <asset-filename> --output-path ./output/asset.png
	
	# Export aseprite asset to png format in scales 1,2,3
	aseprite-assets export <asset-filename> --format png --scales 1,2,3
	
	# Export aseprite asset to png format in sizes 64x64,128x128
	aseprite-assets export <asset-filename> --format png --sizes 64x64,128x128`),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.LoadConfig()
		if err != nil {
			return err
		}

		h := &exportHandler{
			config:      cfg,
			options:     &exportOptions{},
			asepriteCli: aseprite.NewCLI(cfg.AsepritePath, cfg.ScriptDirPath),
		}

		h.options.spriteFilename, _ = cmd.Flags().GetString("sprite-filename")
		h.options.outputFilename, _ = cmd.Flags().GetString("output-filename")
		h.options.format, _ = cmd.Flags().GetString("format")
		h.options.sizes, _ = cmd.Flags().GetString("sizes")
		h.options.scales, _ = cmd.Flags().GetString("scales")

		if h.options.needsSurvey() {
			if err := h.options.collect(); err != nil {
				return err
			}
		}

		err = h.export()

		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	desc := &CommandDescription{
		Title: "Export aseprite asset command",
		Short: "Export aseprite asset to desired format",
		Long:  "Export aseprite asset to desired format with (if you want) specified scale(s) or size(s).",
		Lists: []List{
			{
				Title:      "Supported formats:",
				Indent:     2,
				IndentChar: "- ",
				Items:      aseprite.AvailableExportExtensions(),
			},
		},
	}

	desc.ApplyToCommand(exportCmd)

	exportCmd.Flags().StringP("sprite-filename", "s", "", "aseprite asset filename")
	exportCmd.Flags().StringP("output-filename", "o", "", "output filename")
	exportCmd.Flags().StringP("format", "f", "", "output format")
	exportCmd.Flags().String("sizes", "", "comma separated list of sizes (e.g., \"64x64,128x128\")")
	exportCmd.Flags().String("scales", "", "comma separated list of scales (e.g., \"1,2,3\")")

	rootCmd.AddCommand(exportCmd)
}

func (h *exportHandler) export() error {
	opts := h.options

	if opts.scales != "" && opts.sizes != "" {
		return fmt.Errorf("cannot specify both scales and sizes, choose one")
	}

	if !opts.isSpriteFilenameValid() {
		return fmt.Errorf("invalid sprite filename: %q", opts.spriteFilename)
	}

	outputPath := opts.outputFilename
	if outputPath == "" {
		if !opts.isFormatValid() {
			return errors.New("format required when output filename is not specified")
		}

		outputPath = utils.ChangeFilenameExtension(opts.spriteFilename, opts.format)
	}

	if err := utils.EnsureDirExists(outputPath); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	exportCmd := &commands.ExportSprite{
		BatchMode:      true,
		SpriteFilename: opts.spriteFilename,
		OutputFilename: outputPath,
		Format:         opts.format,
	}

	switch {
	case opts.scales != "":
		if err := ValidateScalesInput(opts.scales); err != nil {
			return fmt.Errorf("invalida scales input: %w", err)
		}
		exportCmd.Scales = opts.scales

	case opts.sizes != "":
		if err := ValidateSizesInput(opts.sizes); err != nil {
			return fmt.Errorf("invalid sizes input: %w", err)
		}
		exportCmd.Sizes = opts.sizes

	}

	if err := h.asepriteCli.ExecuteCommand(exportCmd); err != nil {
		return fmt.Errorf("failed to export sprite: %w", err)
	}

	fmt.Printf("Successfully exported %s to %s\n", opts.spriteFilename, outputPath)
	return nil
}

func (o *exportOptions) collect() error {
	filenamesQuestions := []*survey.Question{
		{
			Name: "sprite-filename",
			Prompt: &survey.Input{
				Message: "Sprite filename:",
				Default: o.spriteFilename,
			},
			Validate: func(val interface{}) error {
				filename := val.(string)
				if !utils.СheckFileExists(filename, false) {
					return errors.New("file does not exist")
				}
				return nil
			},
		},
		{
			Name: "output-filename",
			Prompt: &survey.Input{
				Message: "Output filename:",
				Default: o.outputFilename,
			},
			Validate: func(val interface{}) error {
				filename := val.(string)
				if utils.СheckFileExists(filename, false) {
					return errors.New("file already exists")
				}
				return nil
			},
		},
	}

	formatQuestion := &survey.Select{
		Message: "Select format:",
		Options: aseprite.AvailableExportExtensions(),
		Default: o.format,
	}

	err := survey.Ask(filenamesQuestions, o)
	if err != nil {
		return err
	}

	err = survey.AskOne(formatQuestion, o.format)
	if err != nil {
		return err
	}

	saveOrContinuteQuestion := &survey.Confirm{
		Message: "Do you want to save this sprite? By negating this option, you will be able to continue with exporting (mode selection).",
		Default: true,
	}

	var saveOptions bool
	err = survey.AskOne(saveOrContinuteQuestion, saveOptions)
	if err != nil {
		return err
	} else if saveOptions {
		return nil
	}

	sizingModeQuestion := &survey.Select{
		Message: "Select sizing mode:",
		Options: []string{string(ScaledMode), string(SizedMode)},
		Default: string(ScaledMode),
	}

	var sizingMode SpriteScaleMode
	err = survey.AskOne(sizingModeQuestion, sizingMode)
	if err != nil {
		return err
	}

	if sizingMode == ScaledMode {
		scalesQuestion := &survey.Input{
			Message: "Scales separated by comma (e.g., \"1,2,3\"):",
			Default: o.scales,
		}
		err = survey.AskOne(scalesQuestion, o.scales)
		if err != nil {
			return err
		}
	} else {
		sizesQuestion := &survey.Input{
			Message: "Sizes separated by comma (e.g., \"64x64,128x128\"):",
			Default: o.sizes,
		}
		err = survey.AskOne(sizesQuestion, o.sizes)
		if err != nil {
			return err
		}
	}

	return nil
}

func (o *exportOptions) needsSurvey() bool {
	if !o.isSpriteFilenameValid() || !o.isOutputInfoValid() {
		return true
	}
	if (o.scales != "" && ValidateScalesInput(o.scales) != nil) ||
		(o.sizes != "" && ValidateSizesInput(o.sizes) != nil) {
		return true
	}

	return false
}

func (o *exportOptions) isOutputInfoValid() bool {
	return o.isOutputFilenameValid() || o.isFormatValid()
}

func (o *exportOptions) isSpriteFilenameValid() bool {
	return o.spriteFilename != "" &&
		utils.СheckFileExists(o.spriteFilename, false) &&
		utils.СheckFileExtension(o.spriteFilename, aseprite.SpritesExtensions()...)
}

func (o *exportOptions) isOutputFilenameValid() bool {
	return o.outputFilename != "" &&
		utils.СheckFileExtension(o.outputFilename, aseprite.AvailableExportExtensions()...)
}

func (o *exportOptions) isFormatValid() bool {
	return o.format != "" &&
		slices.Contains(aseprite.AvailableExportExtensions(), utils.PrefExtension(o.format))
}

func ValidateScalesInput(input string) error {
	if input == "" {
		return errors.New("sizes cannot be empty")
	}

	elements := strings.Split(input, ",")

	if err := validateNumberList(elements); err == nil {
		return nil
	}

	return errors.New("invalid format, must be numbers (1,2,3)")
}

func ValidateSizesInput(input string) error {
	if input == "" {
		return errors.New("scales cannot be empty")
	}

	elements := strings.Split(input, ",")

	if err := validatePairList(elements); err == nil {
		return nil
	}

	return errors.New("invalid format, must be pairs (64x64,128x128)")
}

func validateNumberList(elements []string) error {
	for _, elem := range elements {
		elem = strings.TrimSpace(elem)
		if _, err := strconv.Atoi(elem); err != nil {
			return fmt.Errorf("invalid number: %q", elem)
		}
	}
	return nil
}

func validatePairList(elements []string) error {
	for _, elem := range elements {
		elem = strings.TrimSpace(elem)
		parts := strings.Split(elem, "x")
		if len(parts) != 2 {
			return fmt.Errorf("invalid pair format: %q", elem)
		}

		// Validate width
		if _, err := strconv.Atoi(strings.TrimSpace(parts[0])); err != nil {
			return fmt.Errorf("invalid number in pair: %q", parts[0])
		}

		// Validate height
		if _, err := strconv.Atoi(strings.TrimSpace(parts[1])); err != nil {
			return fmt.Errorf("invalid number in pair: %q", parts[1])
		}
	}
	return nil
}
