package export

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
	"github.com/spinozanilast/aseprite-assets-cli/pkg/aseprite/commands/helpers"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/config"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/environment"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/utils"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/utils/files"
)

type SpriteScaleMode string

const (
	ScaledMode SpriteScaleMode = "scaled"
	SizedMode  SpriteScaleMode = "sized"
)

type exportHandler struct {
	config       *config.Config
	asepriteCli  *aseprite.Cli
	options      *exportOptions
	spriteLayers []string
}

type exportOptions struct {
	SpriteFilename string `survey:"sprite-filename"`
	OutputFilename string `survey:"output-filename"`
	FramesIncluded string `survey:"frames-included"`
	SelectedLayer  string
	Format         string
	Sizes          string
	Scales         string
}

func NewExportCmd(env *environment.Environment) *cobra.Command {
	options := &exportOptions{}

	cmd := &cobra.Command{
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
			cfg, err := env.Config()
			if err != nil {
				return err
			}

			h := &exportHandler{
				config:      cfg,
				options:     options,
				asepriteCli: aseprite.NewCLI(cfg.AsepritePath, cfg.ScriptDirPath, cfg.FromSteam),
			}

			if h.options.needsSurvey() {
				utils.PrintlnBold("Do not have enough data to export sprite\n")
				if err := h.collect(); err != nil {
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

	cmd.Flags().StringVarP(&options.SpriteFilename, "sprite-filename", "s", "", "aseprite asset filename")
	cmd.Flags().StringVarP(&options.OutputFilename, "output-filename", "o", "", "output filename")
	cmd.Flags().StringVarP(&options.Format, "format", "f", "", "output format")
	cmd.Flags().StringVarP(&options.SelectedLayer, "layer", "l", "", "separate layer name to export")
	cmd.Flags().StringVar(&options.Sizes, "sizes", "", "comma separated list of sizes (e.g., \"64x64,128x128\")")
	cmd.Flags().StringVar(&options.Scales, "scales", "", "comma separated list of scales (e.g., \"1,2,3\")")
	cmd.Flags().StringVar(&options.FramesIncluded, "frames", "0", "frames included template - zero based (e.g. '0:2', '0', '*'")

	_ = cmd.RegisterFlagCompletionFunc("layer", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		writtenFilename := options.SpriteFilename
		if writtenFilename == "" {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		if !files.CheckFileExists(writtenFilename, false) || !files.CheckFileExtension(writtenFilename, aseprite.SpritesExtensions()...) {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		cfg, err := env.Config()
		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		aseCli := aseprite.NewCLI(cfg.AsepritePath, cfg.ScriptDirPath, cfg.FromSteam)
		cmdLayers := &helpers.SpriteLayersNames{SpriteFilename: writtenFilename}
		output, err := aseCli.ExecuteCommandOutput(cmdLayers)
		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		layers := strings.Split(strings.TrimSpace(output), "\n")
		var suggestions []string
		for _, layer := range layers {
			if layer != "" && strings.HasPrefix(layer, toComplete) {
				suggestions = append(suggestions, layer)
			}
		}

		return suggestions, cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}

func (h *exportHandler) collectSpriteLayers() error {
	spriteLayersCmd := &helpers.SpriteLayersNames{
		SpriteFilename: h.options.SpriteFilename,
	}

	output, err := h.asepriteCli.ExecuteCommandOutput(spriteLayersCmd)
	if err != nil {
		return errors.New("Failed to export sprite layers: " + err.Error())
	}

	h.spriteLayers = strings.Split(output, "\n")

	if len(h.spriteLayers) == 0 {
		return errors.New("failed to export sprite layers")
	}

	return nil
}

func (h *exportHandler) export() error {
	opts := h.options

	if opts.Scales != "" && opts.Sizes != "" {
		return fmt.Errorf("cannot specify both scales and sizes, choose one")
	}

	if !opts.IsSpriteFilenameValid() {
		return fmt.Errorf("invalid sprite filename: %q", opts.SpriteFilename)
	}

	outputPath := opts.OutputFilename
	if outputPath == "" {
		if !opts.IsFormatValid() {
			return errors.New("format required when output filename is not specified")
		}

		outputPath = files.ChangeFilenameExtension(opts.SpriteFilename, opts.Format)
	}

	if err := files.EnsureDirExists(outputPath); err != nil {
		return fmt.Errorf("failed to sprite output directory: %w", err)
	}

	fmt.Printf("Exporting sprite: %s to output: %s\n", opts.SpriteFilename, outputPath)

	exportCmd := &commands.ExportSprite{
		SpriteFilename: opts.SpriteFilename,
		OutputFilename: outputPath,
		Format:         opts.Format,
		FramesIncluded: opts.FramesIncluded,
	}

	if len(opts.SelectedLayer) != 0 {
		if strings.HasSuffix(opts.SelectedLayer, "\r") {
			opts.SelectedLayer = opts.SelectedLayer[:len(opts.SelectedLayer)-1]
		}

		exportCmd.SelectedLayerName = opts.SelectedLayer
	}

	switch {
	case opts.Scales != "":
		if err := ValidateScalesInput(opts.Scales); err != nil {
			return fmt.Errorf("invalid scales input: %w", err)
		}
		exportCmd.Scales = opts.Scales

	case opts.Sizes != "":
		if err := ValidateSizesInput(opts.Sizes); err != nil {
			return fmt.Errorf("invalid sizes input: %w", err)
		}
		exportCmd.Sizes = opts.Sizes

	}

	output, err := h.asepriteCli.ExecuteCommandOutput(exportCmd)

	if err != nil {
		return fmt.Errorf("failed to export sprite: %w", err)
	}

	if output != "" {
		fmt.Printf("Export result:\n%s", output)
	}
	return nil
}

func (h *exportHandler) collect() error {
	o := h.options
	if err := o.collectSourceInfo(h.config); err != nil {
		return err
	}

	if err := o.collectOutputInfo(); err != nil {
		return err
	}

	if err := o.collectFramesInfo(); err != nil {
		return err
	}

	err := h.collectSpriteLayers()
	if err != nil {
		return err
	}

	if err := o.collectLayerInfo(h.spriteLayers); err != nil {
		return err
	}

	if shouldProceed, err := o.confirmExportOptions(); err != nil || !shouldProceed {
		return err
	}

	return o.collectResizeParameters()
}

func (o *exportOptions) collectSourceInfo(cfg *config.Config) error {
	qs := []*survey.Question{
		{
			Name: "SpriteFilename",
			Prompt: &survey.Input{
				Message: "Sprite filename:",
				Suggest: o.spriteSuggestions(cfg),
			},
			Validate: o.validateSpriteFile,
		},
	}
	return survey.Ask(qs, o)
}

func (o *exportOptions) collectOutputInfo() error {
	if o.IsOutputFilenameValid() {
		o.Format = files.GetFileExtension(o.OutputFilename)
		return nil
	}

	if o.IsFormatValid() {
		o.OutputFilename = files.ChangeFilenameExtension(o.SpriteFilename, o.Format)
		return nil
	}

	return o.askOutputDetails()
}

func (o *exportOptions) collectFramesInfo() error {
	if err := survey.AskOne(
		&survey.Input{
			Message: "Frames to export (e.g. '1:3', '*', '5'):",
			Default: "1",
			Suggest: func(toComplete string) []string {
				return []string{"*", "1", "1:5", "2:10"}
			},
		},
		&o.FramesIncluded,
		survey.WithValidator(func(ans interface{}) error {
			input := ans.(string)
			return ValidateFramesInput(input)
		}),
	); err != nil {
		return err
	}

	return nil
}

func (o *exportOptions) collectLayerInfo(spriteLayers []string) error {
	var askLayerSelection bool
	question := &survey.Confirm{
		Message: "Do you want to choose one layer to export?",
		Default: false,
	}

	if err := survey.AskOne(question, &askLayerSelection); err != nil {
		return err
	}

	if askLayerSelection {
		if err := survey.AskOne(
			&survey.Select{
				Message: "Choose a layer to export:",
				Default: spriteLayers[0],
				Options: spriteLayers,
			},
			&o.SelectedLayer,
		); err != nil {
			return err
		}
	}

	return nil
}

func (o *exportOptions) spriteSuggestions(cfg *config.Config) func(string) []string {
	return func(toComplete string) []string {
		var suggestions []string
		for _, folder := range cfg.SpritesFoldersPaths {
			fs, _ := files.FindFilesOfExtensionsRecursiveFlatten(folder, aseprite.SpritesExtensions()...)

			for _, file := range fs {
				if strings.Contains(file, toComplete) {
					suggestions = append(suggestions, file)
				}
			}
		}
		return suggestions
	}
}

func (o *exportOptions) validateSpriteFile(val interface{}) error {
	filename, ok := val.(string)
	if !ok {
		return errors.New("invalid filename type")
	}

	if !files.CheckFileExists(filename, false) {
		return errors.New("file does not exist")
	}

	if !files.CheckFileExtension(filename, aseprite.SpritesExtensions()...) {
		return fmt.Errorf("invalid file extension, allowed: %v", aseprite.SpritesExtensions())
	}

	return nil
}

func (o *exportOptions) askOutputDetails() error {
	var askFormat bool

	question := &survey.Confirm{
		Message: "Do you want to just choose format? (if no, go to write full output path)",
		Default: true,
	}

	if err := survey.AskOne(question, &askFormat); err != nil {
		return err
	}

	if askFormat {
		if err := survey.AskOne(
			&survey.Select{
				Message: "Output format:",
				Options: aseprite.AvailableExportExtensions(),
			},
			&o.Format,
		); err != nil {
			return err
		}
		o.OutputFilename = ""
	} else {
		if err := survey.Ask([]*survey.Question{{
			Name: "OutputFilename",
			Prompt: &survey.Input{
				Message: "Output path:",
				Suggest: o.outputSuggestions(),
				Default: files.ChangeFilenameExtension(o.SpriteFilename, aseprite.PNG.String()),
			},
			Validate: o.validateOutputFile,
		}}, o); err != nil {
			return err
		}
		o.Format = ""
	}

	return nil
}

func (o *exportOptions) outputSuggestions() func(string) []string {
	return func(toComplete string) []string {
		var suggestions []string

		if toComplete == "" {
			return suggestions
		}

		base := files.ChangeFilenameExtension(o.SpriteFilename, "")

		for _, ext := range aseprite.AvailableExportExtensions() {
			suggestion := fmt.Sprintf("%s%s", base, ext)
			if strings.HasPrefix(suggestion, toComplete) {
				suggestions = append(suggestions, suggestion)
			}
		}
		return suggestions
	}
}

func (o *exportOptions) validateOutputFile(val interface{}) error {
	filename, ok := val.(string)
	if !ok {
		return errors.New("invalid filename type")
	}

	if files.CheckFileExists(filename, false) {
		return errors.New("file already exists")
	}

	if !files.CheckFileExtension(filename, aseprite.AvailableExportExtensions()...) {
		return fmt.Errorf("invalid output format, allowed: %v", aseprite.AvailableExportExtensions())
	}

	return nil
}

func (o *exportOptions) confirmExportOptions() (bool, error) {
	var proceed bool
	err := survey.AskOne(&survey.Confirm{
		Message: "Configure scaling/sizing options?",
		Default: true,
	}, &proceed)
	return proceed, err
}

func (o *exportOptions) collectResizeParameters() error {
	var mode SpriteScaleMode
	if err := survey.AskOne(
		&survey.Select{
			Message: "Resize mode:",
			Options: []string{string(ScaledMode), string(SizedMode)},
			Default: string(ScaledMode),
		},
		&mode,
	); err != nil {
		return err
	}

	switch mode {
	case ScaledMode:
		return o.collectScales()
	case SizedMode:
		return o.collectSizes()
	default:
		return fmt.Errorf("unknown resize mode: %s", mode)
	}
}

func (o *exportOptions) collectScales() error {
	return survey.AskOne(
		&survey.Input{
			Message: "Enter scales (comma-separated):",
			Suggest: o.scaleSuggestions(),
		},
		&o.Scales,
		survey.WithValidator(ValidateScalesInputValidator),
	)
}

func (o *exportOptions) scaleSuggestions() func(string) []string {
	return func(toComplete string) []string {
		presets := []string{"1", "2", "3", "0.5,1,2", "1,1.5,2", "2,3,4"}
		return filterSuggestions(toComplete, presets)
	}
}

func (o *exportOptions) collectSizes() error {
	return survey.AskOne(
		&survey.Input{
			Message: "Enter sizes (comma-separated WxH):",
			Suggest: o.sizeSuggestions(),
		},
		&o.Sizes,
		survey.WithValidator(ValidateSizesInputValidator),
	)
}

func (o *exportOptions) sizeSuggestions() func(string) []string {
	return func(toComplete string) []string {
		presets := []string{
			"64x64", "128x128", "256x256",
			"32x32", "48x48", "96x96",
			"64x32", "128x64", "256x128",
		}
		return filterSuggestions(toComplete, presets)
	}
}

func (o *exportOptions) IsOutputInfoValid() bool {
	return o.IsOutputFilenameValid() || o.IsFormatValid()
}

func (o *exportOptions) IsSpriteFilenameValid() bool {
	return o.SpriteFilename != "" &&
		files.CheckFileExists(o.SpriteFilename, false) &&
		files.CheckFileExtension(o.SpriteFilename, aseprite.SpritesExtensions()...)
}

func (o *exportOptions) IsOutputFilenameValid() bool {
	return o.OutputFilename != "" &&
		files.CheckFileExtension(o.OutputFilename, aseprite.AvailableExportExtensions()...)
}

func (o *exportOptions) IsFormatValid() bool {
	return o.Format != "" &&
		slices.Contains(aseprite.AvailableExportExtensions(), files.PrefExtension(o.Format))
}

func ValidateScalesInputValidator(ans interface{}) error {
	input, ok := ans.(string)
	if !ok {
		return errors.New("invalid scales input type")
	}

	return ValidateScalesInput(input)
}

func ValidateSizesInputValidator(ans interface{}) error {
	input, ok := ans.(string)
	if !ok {
		return errors.New("invalid sizes input type")
	}

	return ValidateSizesInput(input)
}

func ValidateScalesInput(input string) error {
	if input == "" {
		return errors.New("scales cannot be empty")
	}

	elements := strings.Split(input, ",")

	if err := ValidateNumberList(elements); err == nil {
		return nil
	}

	return errors.New("invalid format: scales must be a comma-separated list of numbers (e.g., \"1,2,3\")")
}

func ValidateSizesInput(input string) error {
	if input == "" {
		return errors.New("sizes cannot be empty")
	}

	elements := strings.Split(input, ",")

	if err := ValidatePairList(elements); err == nil {
		return nil
	}

	return errors.New("invalid format: sizes must be a comma-separated list of pairs (e.g., \"64x64,128x128\")")
}

func ValidateNumberList(elements []string) error {
	for _, elem := range elements {
		elem = strings.TrimSpace(elem)
		if _, err := strconv.Atoi(elem); err != nil {
			return fmt.Errorf("invalid number: %q", elem)
		}
	}
	return nil
}

func ValidatePairList(elements []string) error {
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

func ValidateFramesInput(input string) error {
	if input == "" {
		return errors.New("frames cannot be empty")
	}

	if input == "*" {
		return nil
	}

	if _, err := strconv.Atoi(input); err == nil {
		return nil
	}

	parts := strings.Split(input, ":")
	if len(parts) != 2 {
		return fmt.Errorf("invalid frames format: %s", input)
	}

	start, err := strconv.Atoi(parts[0])
	if err != nil {
		return fmt.Errorf("invalid start frame: %s", parts[0])
	}

	end, err := strconv.Atoi(parts[1])
	if err != nil {
		return fmt.Errorf("invalid end frame: %s", parts[1])
	}

	if start > end {
		return fmt.Errorf("start frame %d > end frame %d", start, end)
	}

	return nil
}

func filterSuggestions(input string, options []string) []string {
	var matches []string
	for _, opt := range options {
		if strings.HasPrefix(opt, input) {
			matches = append(matches, opt)
		}
	}
	return matches
}

func (o *exportOptions) needsSurvey() bool {
	return !o.IsSpriteFilenameValid() ||
		!o.IsOutputInfoValid() ||
		ValidateFramesInput(o.FramesIncluded) != nil ||
		(o.Scales != "" && ValidateScalesInput(o.Scales) != nil) ||
		(o.Sizes != "" && ValidateSizesInput(o.Sizes) != nil)
}
