package palette

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/environment"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/aseprite"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/aseprite/commands"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/utils"

	"github.com/sashabaranov/go-openai"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/config"
)

type paletteHandler struct {
	config          *config.Config
	openAiClient    *openai.Client
	asepriteCli     *aseprite.AsepriteCLI
	availableModels []string
}

type PaletteSaveVariant int

const (
	SaveAsPreset PaletteSaveVariant = iota
	SaveFile
	Both
)

type PaletteOptions struct {
	Description  string `survey:"description"`
	NumColors    int    `survey:"number-of-colors"`
	Model        string `survey:"model"`
	Advanced     bool   `survey:"advanced"`
	Transparency bool   `survey:"transparency"`
}

type Color struct {
	R, G, B, A uint8
}

type Palette struct {
	Name   string
	Colors []Color
}

type PaletteOutputOptions struct {
	Ui                 bool               `survey:"ui"`
	Directory          string             `survey:"directory"`
	PaletteName        string             `survey:"name"`
	FileType           string             `survey:"file-type"`
	PaletteSaveVariant PaletteSaveVariant `survey:"save-variant"`
	PresetName         string             `survey:"preset-name"`
}

func NewPaletteCmd(env *environment.Environment) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "palette [ARG]",
		Aliases: []string{"p"},
		Short:   "Create aseprite palette from request to LLM",
		RunE:    runPaletteCommand,
	}

	return cmd
}

func runPaletteCommand(cmd *cobra.Command, args []string) error {
	config, err := config.LoadConfig()
	if err != nil {
		return err
	}

	handler := &paletteHandler{
		config:       config,
		openAiClient: initOpenAIClient(config.OpenAiConfig),
		asepriteCli:  aseprite.NewCLI(config.AsepritePath, config.ScriptDirPath),
		availableModels: []string{
			openai.GPT3Dot5Turbo,
			openai.GPT4oMini,
			openai.GPT4o,
			openai.O1Mini,
			openai.GPT4Turbo,
			openai.GPT4,
		},
	}

	return handler.generatePalette()
}

type generationParams struct {
	description  string
	numColors    int
	model        string
	transparency bool
}

func (opts *PaletteOptions) toGenerationParams() generationParams {
	return generationParams{
		description:  opts.Description,
		numColors:    opts.NumColors,
		model:        opts.Model,
		transparency: opts.Transparency,
	}
}

func (h *paletteHandler) generatePalette() error {
	paletteOpts, err := h.collectPaletteOptions()

	if err != nil {
		return err
	}

	utils.PrintlnBold("\n⚡ Generating colors, please wait...")

	colors, err := h.generateColors(paletteOpts.toGenerationParams())
	if err != nil {
		fmt.Errorf("❌ Failed to generate colors: %v", err)
	}

	presentResults(colors, 5)

	outputOpts := &PaletteOutputOptions{}

	paletteConfirmed := collectConfirmPaletteOptions(outputOpts)

	if !paletteConfirmed {
		fmt.Println("⚠️ Palette generation is not confirmed, exiting...")
		return nil
	}

	err = h.collectSaveOptions(outputOpts, paletteOpts.Transparency)

	h.savePalette(outputOpts, paletteOpts, colors)

	if err != nil {
		return err
	}

	return nil
}

func initOpenAIClient(cfg config.OpenAiConfig) *openai.Client {
	apiKey := cfg.ApiKey
	if apiKey == "" {
		fmt.Errorf("OPENAI_API_KEY environment variable is not set\nWrite `asseprite-cli config edit -k <key> -u <url> to set it`")
	}

	apiUrl := cfg.ApiUrl
	if apiUrl == "" {
		fmt.Errorf("Open api url environment variable is not set\nWrite `asseprite-cli config edit -k <key> -u <url> to set it`")
	}

	clientConfig := openai.DefaultConfig(apiKey)
	clientConfig.BaseURL = apiUrl
	return openai.NewClientWithConfig(clientConfig)
}

func (h *paletteHandler) collectPaletteOptions() (*PaletteOptions, error) {
	opts := &PaletteOptions{}
	questions := []*survey.Question{
		{
			Name: "description",
			Prompt: &survey.Input{
				Message: "Color palette description (e.g. 'love, robots, batman'):",
			},
		},
		{
			Name: "number-of-colors",
			Prompt: &survey.Input{
				Message: "Number of colors to generate (if 0 - generate all colors):",
				Default: "10",
			},
			Validate: survey.Required,
		},
		{
			Name: "model",
			Prompt: &survey.Select{
				Message: "AI model to use:",
				Options: h.availableModels,
				Default: h.availableModels[0],
			},
		},
		{
			Name: "advanced",
			Prompt: &survey.Confirm{
				Message: "Enable advanced mode?",
				Default: false,
			},
		},
	}

	err := survey.Ask(questions, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to collect palette options: %w", err)
	}

	if opts.Advanced {
		advancedQuestions := []*survey.Question{
			{
				Name: "transparency",
				Prompt: &survey.Confirm{
					Message: "Include transparency?",
					Default: false,
				},
			},
		}
		err = survey.Ask(advancedQuestions, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to collect advanced options: %w", err)
		}
	}

	return opts, nil
}

func (h *paletteHandler) collectSaveOptions(opts *PaletteOutputOptions, transparencyEnabled bool) error {
	selectedSaveVariant := opts.PaletteSaveVariant

	var questions []*survey.Question
	if selectedSaveVariant != SaveAsPreset {

		outputDirQuestion := &survey.Question{
			Name: "directory",
			Prompt: &survey.Input{
				Message: "Directory to save palettes to:",
				Default: "palettes",
				Suggest: func(_ string) []string {
					return h.config.PalettesFolderPaths
				},
			},
			Validate: func(val interface{}) error {
				dir := val.(string)
				if _, err := os.Stat(dir); os.IsNotExist(err) {
					return errors.New("directory does not exist")
				}
				return nil
			},
		}

		outputNameQuestion := &survey.Question{
			Name: "name",
			Prompt: &survey.Input{
				Message: "Palette name:",
				Default: "palette",
			},
			Validate: func(val interface{}) error {
				name := val.(string)
				dir := opts.Directory
				if dir == "" {
					dir = "palettes"
				}
				path := fmt.Sprintf("%s/%s", dir, name)
				for _, ext := range aseprite.PaletteExtensions() {
					if _, err := os.Stat(path + ext); err == nil {
						return errors.New("file already exists")
					}
				}
				return nil
			},
		}

		questions = append(questions, outputDirQuestion, outputNameQuestion)
	}

	if transparencyEnabled {
		opts.FileType = aseprite.PNG.String()
	} else {
		fileTypeQuestion := &survey.Question{
			Name: "file-type",
			Prompt: &survey.Select{
				Message: "Select file type:",
				Options: aseprite.PaletteExtensions(),
				Default: aseprite.GPL.String(),
			},
		}
		questions = append(questions, fileTypeQuestion)
	}

	err := survey.Ask(questions, opts)
	if err != nil {
		return fmt.Errorf("failed to collect palette options: %w", err)
	}

	return nil
}

func collectConfirmPaletteOptions(opts *PaletteOutputOptions) (confirm bool) {
	confirmGenerationPrompt := &survey.Confirm{
		Message: "Are you want to save this palette?",
		Default: true,
	}
	survey.AskOne(confirmGenerationPrompt, &confirm)

	if confirm {
		saveVariantPrompt := &survey.Select{
			Message: "Select save variant:",
			Options: SaveVariants(),
			Default: SaveFile.String(),
		}

		var selectedVariantS string
		survey.AskOne(saveVariantPrompt, &selectedVariantS)
		selectedVariant := SaveVariantFromString(selectedVariantS)

		if selectedVariant != SaveFile {
			prompt := &survey.Input{
				Message: "Palette preset name:",
				Default: opts.PaletteName,
			}
			survey.AskOne(prompt, &opts.PresetName)
		}

		opts.PaletteSaveVariant = selectedVariant
	}

	return confirm
}

func (h *paletteHandler) savePalette(outputOpts *PaletteOutputOptions, paletteOpts *PaletteOptions, colors []Color) {
	outputPath := filepath.Join(outputOpts.Directory, outputOpts.PaletteName)
	outputPath = utils.EnsureFileExtension(outputPath, outputOpts.FileType)

	palette := Palette{
		Name:   fmt.Sprintf("AI Palette: %s", paletteOpts.Description),
		Colors: colors,
	}

	if outputOpts.FileType == aseprite.PNG.String() {
		if err := generatePNG(palette, outputPath); err != nil {
			fmt.Errorf("Error generating palette: %v", err)
		}
	} else {
		if err := generateGPL(palette, outputPath); err != nil {
			fmt.Errorf("Error generating palette: %v", err)
		}
	}

	utils.PrintFormatted("Generated palette was saved to %s\n", outputPath)

	if outputOpts.PaletteSaveVariant != SaveFile {
		asepriteCli := aseprite.NewCLI(h.config.AsepritePath, h.config.ScriptDirPath)

		asepriteCli.ExecuteCommand(&commands.SavePalette{
			BatchMode:       true,
			PresetName:      outputOpts.PresetName,
			PaletteFilename: outputPath,
		})
	}
}

func (h *paletteHandler) generateColors(params generationParams) ([]Color, error) {
	var basePrompt string

	if params.numColors == 0 {
		basePrompt = "Generate a color palette for: "
	} else {
		basePrompt = fmt.Sprintf("Generate a color palette with exactly %d colors for: ", params.numColors)
	}

	prompt := fmt.Sprintf(`%s: "%s". 
		Use color theory principles to ensure a harmonious palette. 
		If the description asks for shades, arrange them from light to dark.
		Return only hex color codes separated by commas (Upper case only). Example: #FF0000, #00FF00, #0000FF`,
		basePrompt, params.description)

	if params.transparency {
		prompt += "\nInclude transparency in the colors."
	}

	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				for _, r := range `-\|/` {
					fmt.Printf("\r%c Generating colors...", r)
					time.Sleep(100 * time.Millisecond)
				}
			}
		}
	}()

	resp, err := h.openAiClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)
	done <- true
	fmt.Println()

	if err != nil {
		return nil, fmt.Errorf("API error: %v", err)
	}

	response := resp.Choices[0].Message.Content

	logOpenAIResponse(response)
	return parseResponse(response)
}

func logOpenAIResponse(response string) {
	logFile, err := os.OpenFile("openai_responses.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error opening log file: %v\n", err)
		return
	}

	defer logFile.Close()

	logEntry := fmt.Sprintf("Timestamp: %s\nResponse: %s\n\n", time.Now().Format(time.RFC3339), response)
	if _, err := logFile.WriteString(logEntry); err != nil {
		fmt.Printf("Error writing to log file: %v\n", err)
	}
}

func parseResponse(response string) ([]Color, error) {
	var colors []Color
	response = strings.ReplaceAll(response, " ", "")
	colorStrings := strings.Split(response, ",")

	for _, cs := range colorStrings {
		color, err := parseColor(cs)
		if err != nil {
			return nil, fmt.Errorf("invalid AI response format: %v", err)
		}
		colors = append(colors, color)
	}

	return colors, nil
}

func parseColor(input string) (Color, error) {
	input = strings.ToLower(strings.TrimSpace(input))

	if strings.HasPrefix(input, "#") {
		return parseHexColor(input)
	}

	return Color{}, fmt.Errorf("unsupported color format: %s", input)
}

func parseHexColor(hex string) (Color, error) {
	hex = strings.TrimPrefix(hex, "#")

	var r, g, b, a uint8
	switch len(hex) {
	case 4:
		_, err := fmt.Sscanf(hex, "%1x%1x%1x%1x", &r, &g, &b, &a)
		r *= 17
		g *= 17
		b *= 17
		a *= 17
		return Color{r, g, b, a}, err
	case 3:
		_, err := fmt.Sscanf(hex, "%1x%1x%1x", &r, &g, &b)
		r *= 17
		g *= 17
		b *= 17
		return Color{r, g, b, 255}, err
	case 6:
		_, err := fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b)
		return Color{r, g, b, 255}, err
	case 8:
		_, err := fmt.Sscanf(hex, "%02x%02x%02x%02x", &r, &g, &b, &a)
		return Color{r, g, b, a}, err
	default:
		return Color{}, errors.New("invalid hex color length")
	}
}

func generateGPL(palette Palette, path string) error {
	var buf bytes.Buffer

	buf.WriteString("GIMP Palette\n")
	buf.WriteString(fmt.Sprintf("Name: %s\n", palette.Name))
	buf.WriteString("Columns: 0\n#\n")

	for i, color := range palette.Colors {
		buf.WriteString(fmt.Sprintf("%-3d %-3d %-3d Color %d\n",
			color.R, color.G, color.B, i+1))
	}

	return os.WriteFile(path, buf.Bytes(), 0644)
}

func generatePNG(palette Palette, path string) error {
	width := len(palette.Colors)
	height := 1
	img := image.NewNRGBA(image.Rect(0, 0, width, height))

	for i, c := range palette.Colors {
		img.Set(i, 0, color.NRGBA{R: c.R, G: c.G, B: c.B, A: c.A})
	}

	file, err := os.Create(path)

	if err != nil {
		return fmt.Errorf("failed to create PNG file: %w", err)
	}
	defer file.Close()

	if err := png.Encode(file, img); err != nil {
		return fmt.Errorf("failed to encode PNG: %w", err)
	}

	return nil
}

func presentResults(colors []Color, colorsPerRow int) error {
	if len(colors) == 0 {
		return errors.New("no colors generated")
	}

	fmt.Printf("✅ Generated palette with %d colors\n", len(colors))

	for i, color := range colors {
		hex := fmt.Sprintf("#%02x%02x%02x", color.R, color.G, color.B)
		style := lipgloss.NewStyle().Background(lipgloss.Color(hex)).Foreground(lipgloss.Color("#000000")).Padding(0, 1)
		fmt.Print(style.Render(hex))

		if (i+1)%colorsPerRow == 0 {
			fmt.Println()
		}
	}

	fmt.Println()
	return nil
}

func (s PaletteSaveVariant) String() string {
	switch s {
	case SaveAsPreset:
		return "Save as preset"
	case SaveFile:
		return "Save as file"
	case Both:
		return "Save as preset and save palette file"
	default:
		return "Unknown"
	}
}

func SaveVariantFromString(variant string) PaletteSaveVariant {
	switch variant {
	case SaveAsPreset.String():
		return SaveAsPreset
	case SaveFile.String():
		return SaveFile
	case Both.String():
		return Both
	default:
		return SaveAsPreset
	}
}

func SaveVariants() []string {
	return []string{SaveAsPreset.String(), SaveFile.String(), Both.String()}
}
