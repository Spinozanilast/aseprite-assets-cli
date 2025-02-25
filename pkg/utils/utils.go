package utils

import "github.com/spinozanilast/aseprite-assets-cli/pkg/consts"

func MaxLength(strings ...string) int {
	max := 0
	for _, s := range strings {
		if len(s) > max {
			max = len(s)
		}
	}
	return max
}

func ColorFormatFromString(format string) consts.ColorFormat {
	switch format {
	case "hex":
		return consts.HEX
	case "rgb":
		return consts.RGB
	default:
		return consts.HEX
	}
}
