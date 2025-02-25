package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

const (
	configName      = ".aseprite-assets-cli"
	configType      = "json"
	AsepritePathKey = "aseprite_path"
	ScriptDirKey    = "scripts_dir"
	AssetsDirsKey   = "assets_folder_paths"
	OpenAiConfigKey = "open_ai_api"
	PalettesDirsKey = "palettes_folder_paths"
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

func ConfigInfo() string {
	return fmt.Sprintf("Config file used: %s\nAll settings: %+v", viper.ConfigFileUsed(), viper.AllSettings())
}

func SavePaths(appPath string, assetsDirs []string, palettesDirs []string) error {
	viper.Set(AsepritePathKey, appPath)
	viper.Set(AssetsDirsKey, assetsDirs)
	viper.Set(PalettesDirsKey, palettesDirs)
	return saveConfig()
}

func SetScriptDirPath(path string) error {
	viper.Set(ScriptDirKey, path)
	return saveConfig()
}

func SetDefaultScriptDirPath() error {
	pwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	scriptsDir := filepath.Join(pwd, "scripts")
	viper.Set(ScriptDirKey, scriptsDir)
	return saveConfig()
}

func SetOpenAiConfig(apiKey string, apiUrl string) error {
	viper.Set(OpenAiConfigKey, OpenAiConfig{
		ApiKey: apiKey,
		ApiUrl: apiUrl,
	})
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

	// Bind environment variables for OpenAiConfig
	viper.BindEnv(OpenAiConfigKey+".api_key", "OPENAI_API_KEY")
	viper.BindEnv(OpenAiConfigKey+".api_url", "OPENAI_API_URL")

	pwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	viper.SetDefault(AsepritePathKey, "")
	viper.SetDefault(AssetsDirsKey, []string{})
	viper.SetDefault(ScriptDirKey, filepath.Join(pwd, "scripts"))
	viper.SetDefault(OpenAiConfigKey, OpenAiConfig{
		ApiUrl: "https://api.openai.com/v1",
	})
	viper.SetDefault(PalettesDirsKey, []string{})

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			if err := viper.SafeWriteConfig(); err != nil {
				return fmt.Errorf("failed to write initial config: %w", err)
			}
		} else {
			return fmt.Errorf("failed to read config: %w", err)
		}
	}

	return nil
}

func (c *Config) Validate() error {
	if c.AsepritePath == "" {
		return fmt.Errorf("missing required configuration: aseprite_path")
	}

	if len(c.AssetsFolderPaths) == 0 {
		return fmt.Errorf("at least one path required in assets_folder_paths")
	}

	if c.ScriptDirPath == "" {
		return fmt.Errorf("missing required configuration: scripts_dir")
	}

	return nil
}

func saveConfig() error {
	if err := viper.WriteConfig(); err != nil {
		return fmt.Errorf("failed to persist configuration: %w", err)
	}
	return nil
}
