package aseprite

import (
	"fmt"
	"os/exec"
	"path/filepath"
)

type Cli struct {
	AsepritePath   string
	ScriptsDirPath string
}

func NewCLI(asepritePath string, scriptsDirPath string) *Cli {
	return &Cli{
		AsepritePath:   asepritePath,
		ScriptsDirPath: scriptsDirPath,
	}
}

func (a *Cli) CheckPrerequisites() error {
	cmd := exec.Command(a.AsepritePath, "--version")
	_, err := cmd.CombinedOutput()
	return err
}

func (a *Cli) Execute(scriptName string, args []string) (string, error) {
	scriptPath := filepath.Join(a.ScriptsDirPath, scriptName)
	args = append(args, "--script", scriptPath)

	cmd := exec.Command(a.AsepritePath, args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error executing aseprite command: %v\n%s", err, output)
	}
	return string(output), nil
}

func (a *Cli) ExecuteCommand(command Command) error {
	args := command.Args()

	_, err := a.Execute(command.ScriptName(), args)

	if err != nil {
		return err
	}

	return nil
}

func (a *Cli) ExecuteCommandOutput(command Command) (string, error) {
	args := command.Args()

	output, err := a.Execute(command.ScriptName(), args)

	if err != nil {
		return "", err
	}

	return output, nil
}
