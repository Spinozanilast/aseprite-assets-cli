package files

import (
	"fmt"
	"os"
)

// RemoveFile removes a file from the filesystem if exists.
func RemoveFile(filename string) error {
	if CheckFileExists(filename, false) {
		return os.Remove(filename)
	}

	return fmt.Errorf("file %s does not exist", filename)
}
