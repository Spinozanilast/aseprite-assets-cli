package files

import (
	"path/filepath"
	"strings"
)

// EnsureFileExtension ensures a filename has the specified extension.
// Adds extension if missing, replaces existing extension if different.
func EnsureFileExtension(filename, extension string) string {
	extension = PrefExtension(extension)

	if !strings.HasPrefix(extension, ".") {
		extension = "." + extension
	}

	if filepath.Ext(filename) != "" {
		if filepath.Ext(filename) != extension {
			return strings.TrimSuffix(filename, filepath.Ext(filename)) + extension
		}
		return filename
	}

	return filename + extension
}

// ChangeFilenameExtension replaces the extension of a filename with the specified extension.
// Adds extension if filename has none.
func ChangeFilenameExtension(filename, extension string) string {
	if extension != "" {
		extension = PrefExtension(extension)
	}
	dotIdx := strings.LastIndex(filename, ".")

	if dotIdx == -1 {
		return filename + extension
	}

	return filename[:dotIdx] + extension
}

// GetFileExtension returns the extension of a filename (including the dot).
// Wrapper around filepath.Ext.
func GetFileExtension(filename string) string {
	return filepath.Ext(filename)
}
