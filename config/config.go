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
	ScriptDirPath    = "scripts_dir"
	AssetsDirsName   = "assets_folder_paths"
	OpenApiKeyName   = "openai_api_key"
)

type Config struct {
	AsepritePath      string   `mapstructure:"aseprite_path"`
	AssetsFolderPaths []string `mapstructure:"assets_folder_paths"`
	ScriptDirPath     string   `mapstructure:"scripts_dir"`
	OpenAiApiKey      string   `mapstructure:"openai_api_key"`
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

func SavePaths(appPath string, assetsDirs []string) error {
	viper.Set(AsepritePathName, appPath)
	viper.Set(AssetsDirsName, assetsDirs)

	err := viper.WriteConfig()

	if err != nil {
		return fmt.Errorf("fatal error config file: %w", err)
	}

	return nil
}

func SetScriptDirPath(path string) error {
	viper.Set(ScriptDirPath, path)

	return saveConfig()
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

	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	viper.SetDefault(AsepritePathName, "")
	viper.SetDefault(OpenApiKeyName, "")
	viper.SetDefault(AssetsDirsName, []string{})
	viper.SetDefault(ScriptDirPath, pwd+"\\scripts")

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

func SetDefaultScriptDirPath() {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Print(err.Error())
	}

	scriptsDir := pwd + "\\scripts"
	viper.Set(ScriptDirPath, scriptsDir)
	fmt.Printf("Default scripts directory path set to: %s\n", scriptsDir)

	saveConfig()
}

func SetOpenAiApiKey(key string) {
	viper.Set(OpenApiKeyName, key)
	saveConfig()
}

func saveConfig() error {
	err := viper.WriteConfig()
	if err != nil {
		return fmt.Errorf("fatal error config file: %v", err)
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

	if c.ScriptDirPath == "" {
		return fmt.Errorf("scripts directory path is required")
	}
	return nil
}
