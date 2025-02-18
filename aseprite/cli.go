package aseprite

import (
	"bytes"
	"fmt"
	"os/exec"
)

type AsepriteCLI struct {
	AsepritePath string
}

func NewAsepriteCLI(asepritePath string) *AsepriteCLI {
	return &AsepriteCLI{
		AsepritePath: asepritePath,
	}
}

func (a *AsepriteCLI) CheckPrerequisites() error {
	cmd := exec.Command(a.AsepritePath, "--version")
	_, err := cmd.CombinedOutput()
	return err
}

func (a *AsepriteCLI) Execute(args []string) (string, error) {
	cmd := exec.Command(a.AsepritePath, args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	if err != nil {
		return "", fmt.Errorf("error executing aseprite command: %v\n%s", err, stderr.String())
	}

	return stdout.String(), nil
}

func (a *AsepriteCLI) CreateAsset(command AsepriteAssetCreateCommand) error {
	args := command.GetArgs()

	_, err := a.Execute(args)
	if err != nil {
		return err
	}

	return nil
}
