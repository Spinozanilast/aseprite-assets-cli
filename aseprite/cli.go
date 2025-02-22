package aseprite

import (
	"fmt"
	"os/exec"
	"path/filepath"
)

const AsepriteFilesExtension = ".aseprite"

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

func (a *AsepriteCLI) ExecuteCommand(command AsepriteCommand) error {
	args := command.GetArgs()

	_, err := a.Execute(command.GetScriptName(), args)

	if err != nil {
		return err
	}

	return nil
}
