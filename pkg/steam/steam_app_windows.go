package steam

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	AppsFolder         = "steamapps"
	LibraryFoldersFile = "libraryfolders.vdf"
	AppManifestPrefix  = "appmanifest_"
	AppManifestSuffix  = ".acf"
	CommonFolder       = "common"
)

type AppInfo struct {
	AppID      string
	Name       string
	InstallDir string
	Executable string
}

// FindAppByName searches app id by name it's name
func FindAppByName(steamDir, exeName, appName string) (*AppInfo, error) {
	libraryPaths, err := GetSteamLibraryPaths(steamDir)
	if err != nil {
		return &AppInfo{}, err
	}

	exeNameWithExt := exeName
	if !strings.HasSuffix(exeNameWithExt, ".exe") {
		exeNameWithExt += ".exe"
	}

	for _, libraryPath := range libraryPaths {
		appsPath := filepath.Join(libraryPath, AppsFolder)
		apps, err := ParseAppManifests(appsPath)
		if err != nil {
			continue
		}

		for _, app := range apps {
			if strings.EqualFold(app.Name, appName) {
				appDir := filepath.Join(libraryPath, AppsFolder, CommonFolder, app.InstallDir)
				executablePath := filepath.Join(appDir, exeNameWithExt)

				if _, err := os.Stat(executablePath); err != nil {
					return nil, err
				}

				app.Executable = executablePath
				return &app, nil
			}
		}
	}

	return &AppInfo{}, fmt.Errorf("app %s not found in any library", appName)
}

// GetSteamLibraryPaths finds all existing steam libraries files via libraryfolders.vdf
func GetSteamLibraryPaths(steamDir string) ([]string, error) {
	libraryFile := filepath.Join(steamDir, AppsFolder, LibraryFoldersFile)
	content, err := os.ReadFile(libraryFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", libraryFile, err)
	}

	paths := []string{steamDir}

	scanner := bufio.NewScanner(bytes.NewReader(content))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, `"path"`) {
			pathValue, _ := getVdfLineValue(line)
			if pathValue != "" {
				path := strings.Replace(pathValue, `\\'`, `\'`, -1)
				paths = append(paths, path)
			}
		}
	}

	return paths, nil
}

// ParseAppManifests parses all apps manifests
func ParseAppManifests(steamAppsPath string) ([]AppInfo, error) {
	var apps []AppInfo

	files, err := filepath.Glob(filepath.Join(steamAppsPath, AppManifestPrefix+"*"+AppManifestSuffix))
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			continue
		}

		app, err := ParseAppManifest(content)
		if err == nil {
			apps = append(apps, app)
		}
	}

	return apps, nil
}

// ParseAppManifest parses separate app manifest
func ParseAppManifest(content []byte) (AppInfo, error) {
	var app AppInfo

	scanner := bufio.NewScanner(bytes.NewReader(content))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if strings.HasPrefix(line, `"appid"`) {
			app.AppID, _ = getVdfLineValue(line)
		}

		if strings.HasPrefix(line, `"name"`) {
			app.Name, _ = getVdfLineValue(line)
		}

		if strings.HasPrefix(line, `"installdir"`) {
			app.InstallDir, _ = getVdfLineValue(line)
		}
	}

	if app.AppID == "" || app.Name == "" {
		return AppInfo{}, errors.New("invalid app manifest")
	}

	return app, nil
}

func getVdfLineValue(line string) (string, error) {
	parts := strings.SplitN(line, `"`, 5)
	if len(parts) >= 5 {
		return parts[3], nil
	}

	return "", errors.New("invalid app manifest")
}
