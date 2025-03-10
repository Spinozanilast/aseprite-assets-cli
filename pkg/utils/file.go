package utils

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/ncruces/zenity"
	"github.com/skratchdot/open-golang/open"
)

const defaultPath string = `` // Default starting path for file dialogs (empty string uses system default)

var (
	exeFileFilter = zenity.FileFilter{Name: "Executable files", Patterns: []string{"*.exe"}, CaseFold: true}
)

// OpenExecutableFilesDialog opens a file selection dialog filtered for executable files.
// Returns the selected file path or an error if the dialog is canceled or fails.
func OpenExecutableFilesDialog(title string) (string, error) {
	return zenity.SelectFile(zenity.Filename(defaultPath), zenity.FileFilters{exeFileFilter}, zenity.Title(title))
}

// OpenDirectoryDialog opens a directory selection dialog.
// Returns the selected directory path or an error if the dialog is canceled or fails.
func OpenDirectoryDialog(title string) (string, error) {
	return zenity.SelectFile(
		zenity.Filename(defaultPath),
		zenity.Directory(),
		zenity.Title(title))
}

// OpenFile opens the specified file using the system's default application.
// Returns an error if the operation fails.
func OpenFile(path string) error {
	err := open.Start(path)
	if err != nil {
		return err
	}

	return nil
}

// OpenFileWith opens the specified file using a specific application.
// Returns an error if the operation fails.
func OpenFileWith(path string, appExe string) error {
	err := open.StartWith(path, appExe)
	if err != nil {
		return err
	}

	return nil
}

// FindFilesOfExtensions searches a directory for files with specified extensions.
// Returns base names of matching files. Returns error if no matches found.
// Non-recursive search.
func FindFilesOfExtensions(folderPath string, extensions ...string) ([]string, error) {
	var filenames []string

	for _, ext := range extensions {
		ext = PrefExtension(ext)
		pattern := filepath.Join(folderPath, "*"+ext)

		fullPaths, err := filepath.Glob(pattern)
		if err != nil {
			return nil, err
		}

		for _, path := range fullPaths {
			filenames = append(filenames, filepath.Base(path))
		}
	}

	if len(filenames) == 0 {
		return nil, fmt.Errorf("no files found with extensions: %v in %s", extensions, folderPath)
	}

	return filenames, nil
}

// FindFilesOfExtensionsRecursive recursively searches a directory for files with specified extensions.
// Returns a map of directory paths to their contained files. Returns error if walk fails.
func FindFilesOfExtensionsRecursive(folderPath string, extensions ...string) (map[string][]string, error) {
	results := make(map[string][]string)
	extMap := make(map[string]bool)

	for _, ext := range extensions {
		extMap[PrefExtension(ext)] = true
	}

	err := filepath.WalkDir(folderPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		fileExt := strings.ToLower(filepath.Ext(path))
		if extMap[fileExt] {
			dir := filepath.Dir(path)
			results[dir] = append(results[dir], filepath.Base(path))
		}

		return nil
	})

	return results, err
}

// FindFilesOfExtensionsRecursiveFlatten recursively searches a directory for files with specified extensions.
// Returns a flat list of full file paths. Returns error if walk fails.
func FindFilesOfExtensionsRecursiveFlatten(folderPath string, extensions ...string) ([]string, error) {
	results := make([]string, 0)
	extMap := make(map[string]bool)

	for _, ext := range extensions {
		extMap[PrefExtension(ext)] = true
	}

	err := filepath.WalkDir(folderPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		fileExt := strings.ToLower(filepath.Ext(path))
		if extMap[fileExt] {
			results = append(results, path)
		}

		return nil
	})

	return results, err
}

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
	info, err := os.Stat(filepath.Dir(path))
	if err != nil {
		return err
	}

	if !info.IsDir() {
		os.Mkdir(path, os.ModePerm)
		_, err = os.Stat(path)
		if err != nil {
			return err
		}
	}

	return nil
}

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

// ChangeFilenameExtension replaces the extension of a filename with the specified extension.
// Adds extension if filename has none.
func ChangeFilenameExtension(filename, extension string) string {
	extension = PrefExtension(extension)
	dotIdx := strings.LastIndex(filename, ".")

	if dotIdx == -1 {
		return filename + extension
	}

	return filename[:dotIdx] + extension
}

// CheckFileExists checks if a file or directory exists at the specified path.
// watchDir=true checks for directory existence, false checks for file existence.
func CheckFileExists(path string, watchDir bool) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir() == watchDir
}

// GetFileExtension returns the extension of a filename (including the dot).
// Wrapper around filepath.Ext.
func GetFileExtension(filename string) string {
	return filepath.Ext(filename)
}

// PrefExtension ensures a filename end with extension.
// Adds leading dot if missing, preserves existing dots.
func PrefExtension(extension string) string {
	if extension[0] == '.' {
		return extension
	}

	return "." + extension
}

// appendDirFilesIfExists helper function adds files with specified extension from directory to map.
// Non-exported internal function.
func appendDirFilesIfExists(dirsFiles map[string][]string, dir string, extension string) (map[string][]string, error) {
	files, err := FindFilesOfExtensions(dir, extension)

	if err != nil {
		return nil, err
	} else if len(files) > 0 {
		dirsFiles[dir] = files
	}

	return dirsFiles, nil
}
