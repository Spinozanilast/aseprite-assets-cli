package aseprite

import (
	"fmt"
	"os/exec"
	"path/filepath"
)

type AsepriteCLI struct {
	AsepritePath   string
	ScriptsDirPath string
}

func NewCLI(asepritePath string, scriptsDirPath string) *AsepriteCLI {
	return &AsepriteCLI{
		AsepritePath:   asepritePath,
		ScriptsDirPath: scriptsDirPath,
	}
}

func (a *AsepriteCLI) CheckPrerequisites() error {
	cmd := exec.Command(a.AsepritePath, "--version")
	_, err := cmd.CombinedOutput()
	return err
}

func (a *AsepriteCLI) Execute(scriptName string, args []string) (string, error) {
	scriptPath := filepath.Join(a.ScriptsDirPath, scriptName)
	args = append(args, "--script", scriptPath)

	cmd := exec.Command(a.AsepritePath, args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error executing aseprite command: %v\n%s", err, output)
	}
	return string(output), nil
}

func (a *AsepriteCLI) ExecuteCommand(command Command) error {
	args := command.Args()

	_, err := a.Execute(command.ScriptName(), args)

	if err != nil {
		return err
	}

	return nil
}

func (a *AsepriteCLI) ExecuteCommandOutput(command Command) (string, error) {
	args := command.Args()

	output, err := a.Execute(command.ScriptName(), args)

	if err != nil {
		return "", err
	}

	return output, nil
}
