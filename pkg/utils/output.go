package utils

import (
	"github.com/fatih/color"
)

func PrintFormatted(format string, a ...interface{}) {
	color.New(color.BgHiMagenta).Printf(format, a...)
}

func PrintlnSuccess(message string) {
	color.New(color.FgGreen).Println(message)
}

func PrintlnBold(message string) {
	color.New(color.Bold).Println(message)
}

func PrintError(message string) {
	color.New(color.FgRed).Println(message)
}
