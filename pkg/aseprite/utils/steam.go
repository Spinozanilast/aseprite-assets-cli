package utils

import "github.com/spinozanilast/aseprite-assets-cli/pkg/steam"

// FindSteamAsepriteExecutable try find steam executable using registered in
// registry steam install path or default and accessing apps certificates
func FindSteamAsepriteExecutable() (string, error) {
	steamPath, err := steam.FindSteamPath()
	if err != nil {
		return "", err
	}

}
