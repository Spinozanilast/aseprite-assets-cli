package utils

import "github.com/spinozanilast/aseprite-assets-cli/pkg/consts"

func MinLength(strings ...string) int {
	minl := 300
	for _, s := range strings {
		if len(s) < minl {
			minl = len(s)
		}
	}
	return minl
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
