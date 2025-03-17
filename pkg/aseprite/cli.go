package aseprite

import (
	"fmt"
	"os/exec"
	"path/filepath"
)

type Cli struct {
	AsePath        string
	ScriptsDirPath string
	FromSteam      bool
}

type SteamInfo struct {
	AppId string
}

func NewCLI(asePath string, scriptsDirPath string, fromSteam bool) *Cli {
	return &Cli{
		AsePath:        asePath,
		ScriptsDirPath: scriptsDirPath,
		FromSteam:      fromSteam,
	}
}

func (ac *Cli) CheckPrerequisites() error {
	cmd := exec.Command(ac.AsePath, "--version")
	_, err := cmd.CombinedOutput()
	return err
}

func (ac *Cli) Execute(scriptName string, args []string) (string, error) {
	scriptPath := filepath.Join(ac.ScriptsDirPath, scriptName)
	args = append(args, "--script", scriptPath)

	cmd := exec.Command(ac.AsePath, args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error executing aseprite command: %v\n%s", err, output)
	}

	return string(output), nil
}

func (ac *Cli) ExecuteCommand(command Command) error {
	args := command.Args()

	_, err := ac.Execute(command.ScriptName(), args)

	if err != nil {
		return err
	}

	return nil
}

func (ac *Cli) ExecuteCommandOutput(command Command) (string, error) {
	args := command.Args()

	output, err := ac.Execute(command.ScriptName(), args)

	if err != nil {
		return "", err
	}

	return output, nil
}
