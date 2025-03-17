package files

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
)

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
