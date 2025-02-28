package environment

import "github.com/spinozanilast/aseprite-assets-cli/pkg/config"

type Environment struct {
	Config func() (*config.Config, error)
}

func NewEnvironment(config func() (*config.Config, error)) Environment {
	return Environment{config}
}
