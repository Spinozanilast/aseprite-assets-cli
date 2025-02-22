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
	OpenAiConfigName = "open_ai_api"
	PalettesDirsName = "palettes_folder_paths"
)

type OpenAiConfig struct {
	ApiKey string `mapstructure:"api_key" json:"api_key"`
	ApiUrl string `mapstructure:"api_url" json:"api_url"`
}

type Config struct {
	AsepritePath        string       `mapstructure:"aseprite_path"`
	AssetsFolderPaths   []string     `mapstructure:"assets_folder_paths"`
	ScriptDirPath       string       `mapstructure:"scripts_dir"`
	OpenAiConfig        OpenAiConfig `mapstructure:"open_ai_api"`
	PalettesFolderPaths []string     `mapstructure:"palettes_folder_paths"`
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

func SavePaths(appPath string, assetsDirs []string, palettesDirs []string) error {
	viper.Set(AsepritePathName, appPath)
	viper.Set(AssetsDirsName, assetsDirs)
	viper.Set(PalettesDirsName, palettesDirs)
	return saveConfig()
}

func SetScriptDirPath(path string) error {
	viper.Set(ScriptDirPath, path)
	return saveConfig()
}

func SetDefaultScriptDirPath() {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	scriptsDir := pwd + "\\scripts"
	viper.Set(ScriptDirPath, scriptsDir)
	fmt.Printf("Default scripts directory path set to: %s\n", scriptsDir)
	saveConfig()
}

func SetOpenAiConfig(apiKey string, apiUrl string) {
	viper.Set(OpenAiConfigName, OpenAiConfig{ApiKey: apiKey, ApiUrl: apiUrl})
	saveConfig()
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
	viper.SetDefault(AssetsDirsName, []string{})
	viper.SetDefault(ScriptDirPath, pwd+"\\scripts")
	viper.SetDefault(OpenAiConfigName, OpenAiConfig{ApiKey: os.Getenv("OPENAI_API_KEY"), ApiUrl: "https://api.openai.com/v1"})
	viper.SetDefault(PalettesDirsName, []string{})

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			if err := viper.SafeWriteConfig(); err != nil {
				return fmt.Errorf("fatal error writing config file: %w", err)
			}
		} else {
			return fmt.Errorf("fatal error reading file: %w", err)
		}
	}

	return nil
}

func saveConfig() error {
	if err := viper.WriteConfig(); err != nil {
		return fmt.Errorf("fatal error writing config file: %w", err)
	}
	return nil
}
