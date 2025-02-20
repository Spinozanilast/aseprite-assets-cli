package cmd

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"

	openai "github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
	config "github.com/spinozanilast/aseprite-assets-cli/config"
)

type PaletteConfig struct {
	Description string
	NumColors   int
	Model       string
	OutputFile  string
}

type Color struct {
	R, G, B uint8
}

type Palette struct {
	Name   string
	Colors []Color
}

var cfg PaletteConfig

var paletteCmd = &cobra.Command{
	Use:     "palette [ARG]",
	Aliases: []string{"p"},
	Short:   "Create aseprite palette from request to LLM",
	RunE: func(cmd *cobra.Command, args []string) error {
		if cfg.Description == "" {
			fatalError("description cannot be empty")
		}

		config, err := config.LoadConfig()
		if err != nil {
			return err
		}

		apiKey := config.OpenAiApiKey
		if apiKey == "" {
			fatalError("OPENAI_API_KEY environment variable is not set")
		}

		colors, err := generateColorsWithAi(cfg.Description, cfg.NumColors, cfg.Model, apiKey)
		if err != nil {
			fatalError("failed to generate colors: %v", err)
		}

		palette := Palette{
			Name:   fmt.Sprintf("AI Palette: %s", cfg.Description),
			Colors: colors,
		}

		if cfg.OutputFile == "" {
			cfg.OutputFile = "ai-palette.gpl"
		}

		if err := generateGPL(palette, cfg.OutputFile); err != nil {
			fatalError("Error generating palette: %v", err)
		}

		fmt.Printf("Generated palette '%s' with %d colors\n", cfg.OutputFile, len(colors))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(paletteCmd)
	paletteCmd.Flags().StringVarP(&cfg.Description, "descriptin", "d", "", "Color palette description (e.g. 'love, robots, batman')")
	paletteCmd.Flags().IntVarP(&cfg.NumColors, "num-colors", "n", 5, "Number of colors to generate")
	paletteCmd.Flags().StringVarP(&cfg.Model, "model", "m", "gpt-3.5-turbo", "AI model to use")
	paletteCmd.Flags().StringVarP(&cfg.OutputFile, "output", "o", "ai-palette.gpl", "Output file path")
}

func generateColorsWithAi(description string, numColors int, model string, apiKey string) ([]Color, error) {
	config := openai.DefaultConfig(apiKey)
	client := openai.NewClientWithConfig(config)

	prompt := fmt.Sprintf(`Generate a color palette with %d colors for: "%s". 
		Return only hex color codes separated by commas. Example: #FF0000, #00FF00, #0000FF`,
		numColors, description)

	resp, err := client.CreateChatCompletion(
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
	if err != nil {
		return nil, fmt.Errorf("API error: %v", err)
	}

	return parseAIResponse(resp.Choices[0].Message.Content)
}

func parseAIResponse(response string) ([]Color, error) {
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

	if len(input) == 6 || len(input) == 3 {
		return parseHexColor("#" + input)
	}

	return Color{}, fmt.Errorf("unsupported color format: %s", input)
}

func parseHexColor(hex string) (Color, error) {
	hex = strings.TrimPrefix(hex, "#")

	var r, g, b uint8
	switch len(hex) {
	case 3:
		_, err := fmt.Sscanf(hex, "%1x%1x%1x", &r, &g, &b)
		r *= 17
		g *= 17
		b *= 17
		return Color{r, g, b}, err
	case 6:
		_, err := fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b)
		return Color{r, g, b}, err
	default:
		return Color{}, fmt.Errorf("invalid hex color length")
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
