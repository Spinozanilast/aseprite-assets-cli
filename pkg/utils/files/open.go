package files

import (
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
