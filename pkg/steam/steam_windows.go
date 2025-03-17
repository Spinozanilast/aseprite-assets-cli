package steam

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"golang.org/x/sys/windows/registry"
)

const (
	DefaultExe       = "steam.exe"
	DefaultSteamPath = `C:\Program Files (x86)\Steam`
	DefaultSubKey    = `SOFTWARE\WOW6432Node\Valve\Steam`
	InstallPathKey   = "InstallPath"
)

var PossibleRegParentKeys = []registry.Key{
	registry.LOCAL_MACHINE,
	registry.CURRENT_USER,
}

// StartSteamAppByIdOnDir starts steam app using id
// that can be found using FindAppByName
func StartSteamAppByIdOnDir(steamDir, id string, args ...string) (string, error) {
	steamExe := filepath.Join(steamDir, DefaultExe)
	if _, err := os.Stat(steamExe); os.IsNotExist(err) {
		return "", fmt.Errorf("steamExe %s does not exist", steamExe)
	}

	execArgs := []string{"-applaunch", id}
	execArgs = append(execArgs, args...)

	cmd := exec.Command(steamExe, execArgs...)
	if err := cmd.Start(); err != nil {
		return "", fmt.Errorf("failde to start steam app: %w", err)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error starting app with command: %v\n%s", err, output)
	}

	return string(output), nil
}

// StartSteamAppById wraps StartSteamAppByIdOnDir with pre search of steam directory
func StartSteamAppById(id string, args ...string) (string, error) {
	steamDir, err := FindSteamPath()

	if err != nil {
		return "", err
	}

	output, err := StartSteamAppByIdOnDir(steamDir, id, args...)
	if err != nil {
		return "", err
	}

	return output, nil
}

// FindSteamPath locates the Steam installation directory by checking:
// 1. Windows registry entries (both HKLM and HKCU)
// 2. Default installation path (if exists)
// Returns found path or error if neither location is valid
func FindSteamPath() (string, error) {
	//Try registry lookup first
	path, registryErr := FindSteamPathFromRegistry()
	if registryErr == nil {
		return path, nil
	}

	//Check if steam.exe exists in default steam path
	defaultExePath := filepath.Join(DefaultSteamPath, DefaultExe)
	if _, err := os.Stat(defaultExePath); err == nil {
		return DefaultSteamPath, nil
	}

	return "", fmt.Errorf(
		"steam installation not found: registry error (%v), default path (%s) error: %w",
		registryErr, DefaultSteamPath, os.ErrNotExist)
}

// FindSteamPathFromRegistry searches through possible registry locations
func FindSteamPathFromRegistry() (string, error) {
	for _, parentKey := range PossibleRegParentKeys {
		path, err := FindSteamPathByRegistryKey(parentKey, DefaultSubKey)

		if err == nil {
			return path, nil
		}
	}

	return "", errors.New("registry entry not found in any monitored hives")
}

// FindSteamPathByRegistryKey reads Steam path from a specific registry key
func FindSteamPathByRegistryKey(parentKey registry.Key, subKey string) (string, error) {
	k, err := registry.OpenKey(parentKey, subKey, registry.QUERY_VALUE)

	if err != nil {
		return "", fmt.Errorf("registry key access failed: %w", err)
	}

	defer k.Close()

	path, _, err := k.GetStringValue(InstallPathKey)
	if err != nil {
		return "", fmt.Errorf("installpath value read failed: %w", err)
	}

	return path, nil
}
