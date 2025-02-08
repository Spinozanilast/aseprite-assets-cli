package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

const (
	configName       = ".aseprite-assets-cli"
	configType       = "json"
	AsepritePathName = "aseprite_path"
	AssetsDirsName   = "assets_folder_paths"
)

type Config struct {
	AsepritePath      string   `mapstructure:"aseprite_path"`
	AssetsFolderPaths []string `mapstructure:"assets_folder_paths"`
}

func LoadConfig() (*Config, error) {
	if err := initConfig(); err != nil {
		return nil, fmt.Errorf("config initialization failed: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("config unmarshal failed: %w", err)
	}

	return &config, nil
}

func ConfigInfo() {
	fmt.Printf("Config file used: %s\n", viper.ConfigFileUsed())
	fmt.Printf("All settings: %+v\n", viper.AllSettings())
}

func SaveConfig(appPath string, assetsDirs []string) error {
	viper.Set(AsepritePathName, appPath)
	viper.Set(AssetsDirsName, assetsDirs)

	err := viper.WriteConfig()

	if err != nil {
		return fmt.Errorf("fatal error config file: %w", err)
	}

	return nil
}

func initConfig() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	viper.SetConfigName(configName)
	viper.SetConfigType(configType)
	viper.AddConfigPath(homeDir)

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetDefault(AsepritePathName, "")
	viper.SetDefault(AssetsDirsName, []string{})

	err = viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			err = viper.SafeWriteConfig()
			if err != nil {
				return fmt.Errorf("fatal error writing config file: %w", err)
			}
		} else {
			return fmt.Errorf("fatal error reading file: %w", err)

		}
	}

	return nil
}

func (c *Config) Validate() error {
	if c.AsepritePath == "" {
		return fmt.Errorf("aseprite path is required")
	}

	if len(c.AssetsFolderPaths) == 0 {
		return fmt.Errorf("at least one assets folder path is required")
	}
	return nil
}
