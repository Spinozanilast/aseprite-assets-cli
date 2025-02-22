package util

import (
	"fmt"
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

func FindFilesOfExtension(folderPath string, extension string) ([]string, error) {
	_, err := CheckAnyFileOfExtensionExists(folderPath, extension)
	if err != nil {
		return nil, err
	}

	extension = prefExtension(extension)

	pattern := filepath.Join(folderPath, "*"+extension)
	matches, err := filepath.Glob(pattern)

	if err != nil {
		return nil, err
	}

	return matches, nil
}

func CheckAnyFileOfExtensionExists(folderPath string, extension string) (bool, error) {
	extension = prefExtension(extension)

	pattern := filepath.Join(folderPath, "*"+extension)
	matches, err := filepath.Glob(pattern)

	if err != nil {
		return false, err
	}

	if matches != nil {
		return true, nil
	} else {
		return false, fmt.Errorf("folder (%s) doesnt has any file of %s extension", folderPath, extension)
	}
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
