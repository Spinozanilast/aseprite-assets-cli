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

const defaultPath string = ``

var (
	exeFileFilter = zenity.FileFilter{Name: "Executable files", Patterns: []string{"*.exe"}, CaseFold: true}
)

func OpenExecutableFilesDialog(title string) (string, error) {
	return zenity.SelectFile(zenity.Filename(defaultPath), zenity.FileFilters{exeFileFilter}, zenity.Title(title))
}

func OpenDirectoryDialog(title string) (string, error) {
	return zenity.SelectFile(
		zenity.Filename(defaultPath),
		zenity.Directory(),
		zenity.Title(title))
}

func OpenFile(path string) error {
	err := open.Start(path)
	if err != nil {
		return err
	}

	return nil
}

func OpenFileWith(path string, appname string) error {
	err := open.StartWith(path, appname)
	if err != nil {
		return err
	}

	return nil
}

func FindFilesOfExtensions(folderPath string, extensions ...string) ([]string, error) {
	var filenames []string

	for _, ext := range extensions {
		ext = prefExtension(ext)
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

func FindFilesOfExtensionsRecursive(folderPath string, extensions ...string) (map[string][]string, error) {
	results := make(map[string][]string)
	extMap := make(map[string]bool)

	for _, ext := range extensions {
		extMap[prefExtension(ext)] = true
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

func CheckAnyFileOfExtensionsExists(folderPath string, extensions ...string) (bool, error) {
	for _, ext := range extensions {
		ext = prefExtension(ext)
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

func EnsureFileExtension(filename, extension string) string {
	extension = prefExtension(extension)

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

func СheckFileExists(path string, watchDir bool) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir() == watchDir
}

func prefExtension(extension string) string {
	if extension[0] == '.' {
		return extension
	}

	return "." + extension
}

func appendDirFilesIfExists(dirsFiles map[string][]string, dir string, extension string) (map[string][]string, error) {
	files, err := FindFilesOfExtensions(dir, extension)

	if err != nil {
		return nil, err
	} else if len(files) > 0 {
		dirsFiles[dir] = files
	}

	return dirsFiles, nil
}
