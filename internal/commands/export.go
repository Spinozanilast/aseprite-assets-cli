package commands

import (
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"
)

var exportCmd = &cobra.Command{
	Use:     "export [ARG]",
	Aliases: []string{"exp", "e"},
	Example: heredoc.Doc(`
	# Export aseprite asset to png format
	aseprite-assets export <asset-filename> --format png

	# Export aseprite asset to png format and save it to the specified path
	aseprite-assets export <asset-filename> --output-path ./output/asset.png`),
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	desc := &CommandDescription{
		Title: "Export aseprite asset command",
		Short: "Export aseprite asset to desired format",
		Long:  "Export aseprite asset to desired format with specified format and output path.",
		Lists: []List{
			{
				Title:      "Supported formats:",
				Indent:     2,
				IndentChar: "- ",
				Items: []string{
					"png",
					"jpg",
					"webp",
					"gif",
					"svg",
					"bmp",
				},
			},
		},
	}

	desc.ApplyToCommand(exportCmd)

	exportCmd.Flags().StringP("format", "f", "", "Export format")
	exportCmd.Flags().StringP("output", "o", "", "Output path")
	rootCmd.AddCommand(exportCmd)
}
