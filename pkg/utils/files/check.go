package files

import (
	"fmt"
	"os"
	"path/filepath"
)

// CheckAnyFileOfExtensionsExists checks if any files with specified extensions exist in a directory.
// Returns true if at least one match found. Returns error if no matches found.
// Non-recursive search.
func CheckAnyFileOfExtensionsExists(folderPath string, extensions ...string) (bool, error) {
	for _, ext := range extensions {
		ext = PrefExtension(ext)
		pattern := filepath.Join(folderPath, "*"+ext)
		matches, err := filepath.Glob(pattern)
		if err != nil {
			return false, err
		}
		if len(matches) > 0 {
			return true, nil
		}
	}
	return false, fmt.Errorf("no files found with extensions %v in %s", extensions, folderPath)
}

// EnsureDirExists verifies and creates the directory containing the specified file path if necessary.
// Returns error if directory creation fails or path is invalid.
func EnsureDirExists(path string) error {
	dir := filepath.Dir(path)
	_, err := os.Stat(dir)

	if err != nil {
		err := os.Mkdir(dir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}

// CheckFileExtension checks if a file path has one of the specified extensions.
// Comparison is case-insensitive. Returns true if extension matches.
func CheckFileExtension(path string, extensions ...string) bool {
	for _, ext := range extensions {
		ext = PrefExtension(ext)
		if filepath.Ext(path) == ext {
			return true
		}
	}

	return false
}

// CheckFileExists checks if a file or directory exists at the specified path.
// watchDir=true checks for directory existence, false checks for file existence.
func CheckFileExists(path string, watchDir bool) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir() == watchDir
}

// PrefExtension ensures a filename end with extension.
// Adds leading dot if missing, preserves existing dots.
func PrefExtension(extension string) string {
	if extension[0] == '.' {
		return extension
	}

	return "." + extension
}
