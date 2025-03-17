package config

import (
	"errors"
	"fmt"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/aseprite"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/steam"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/utils/files"
	"go.uber.org/multierr"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/viper"
)

const (
	openAiApiUrl = "https://api.openai.com/v1"

	configName = ".aseprite-assets-cli"
	configType = "json"

	fromSteamKey    = "from_steam"
	appIdKey        = "app_id"
	asepritePathKey = "aseprite_path"
	scriptDirKey    = "scripts_dir"
	spriteDirsKey   = "assets_folder_paths"
	openAiConfigKey = "open_ai_api"
	palettesDirsKey = "palettes_folder_paths"
)

type OpenAiConfig struct {
	ApiKey string `mapstructure:"api_key" json:"api_key"`
	ApiUrl string `mapstructure:"api_url" json:"api_url"`
}

type Config struct {
	FromSteam            bool         `mapstructure:"from_steam" json:"from_steam"`
	AppId                string       `mapstructure:"app_id" json:"app_id"`
	AsepritePath         string       `mapstructure:"aseprite_path"`
	SpritesFoldersPaths  []string     `mapstructure:"assets_folder_paths"`
	PalettesFoldersPaths []string     `mapstructure:"palettes_folder_paths"`
	ScriptDirPath        string       `mapstructure:"scripts_dir"`
	OpenAiConfig         OpenAiConfig `mapstructure:"open_ai_api"`
}

// LoadConfig loads the configuration from the file system (use it again if you need config after updating)
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

func Info() string {
	var settingsBuilder strings.Builder

	for key, value := range viper.AllSettings() {
		settingsBuilder.WriteString(fmt.Sprintf("\t%v: %v\n", key, value))
	}

	return fmt.Sprintf("Config file used: %s\nAll settings: \n[\n%s]", viper.ConfigFileUsed(), settingsBuilder.String())
}

func SavePaths(appPath string, spritesDirs []string, palettesDirs []string) error {
	viper.Set(asepritePathKey, appPath)
	viper.Set(spriteDirsKey, spritesDirs)
	viper.Set(palettesDirsKey, palettesDirs)
	return saveConfig()
}

func SetScriptDirPath(path string) error {
	viper.Set(scriptDirKey, path)
	return saveConfig()
}

func SetDefaultScriptDirPath() error {
	pwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	scriptsDir := filepath.Join(pwd, "scripts")
	viper.Set(scriptDirKey, scriptsDir)
	return saveConfig()
}

func SetOpenAiConfig(apiKey string, apiUrl string) error {
	viper.Set(openAiConfigKey, OpenAiConfig{
		ApiKey: apiKey,
		ApiUrl: apiUrl,
	})
	return saveConfig()
}

func (c *Config) Validate() error {
	var errs []error

	if c.FromSteam && c.AppId == "" {
		errs = append(errs, errors.New("missing Steam AppID when using Steam configuration"))
	} else {
		if c.AsepritePath == "" {
			errs = append(errs, errors.New("missing required configuration: aseprite_path"))
		} else if !filepath.IsAbs(c.AsepritePath) || (runtime.GOOS == "windows" && !files.CheckFileExtension(c.AsepritePath, "exe")) {
			errs = append(errs, errors.New("aseprite path must be absolute executable path"))
		}
	}

	if len(c.SpritesFoldersPaths) == 0 {
		errs = append(errs, errors.New("at least one sprite folder path required"))
	} else {
		for _, path := range c.SpritesFoldersPaths {
			if !filepath.IsAbs(path) {
				errs = append(errs, fmt.Errorf("sprites path must be absolute: %s", path))
			}
		}
	}

	if len(c.PalettesFoldersPaths) > 0 {
		for _, path := range c.PalettesFoldersPaths {
			if !filepath.IsAbs(path) {
				errs = append(errs, fmt.Errorf("palettes path must be absolute: %s", path))
			}
		}
	}

	if c.ScriptDirPath == "" {
		errs = append(errs, errors.New("missing required configuration: scripts_dir"))
	} else if !filepath.IsAbs(c.ScriptDirPath) {
		errs = append(errs, errors.New("scripts dir path need to be absolute"))
	}

	return multierr.Combine(errs...)
}

func TryFindAsepritePath() (asePath string, appId string, fromSteam bool) {
	// Check environment variable first
	if envPath := os.Getenv("ASEPRITE"); envPath != "" {
		return envPath, "", false
	}

	// Fallback to Steam detection
	steamPath, _ := steam.FindSteamPath()
	aseSteamInfo, _ := steam.FindAppByName(steamPath, aseprite.Name, aseprite.Name)

	if aseSteamInfo == nil {
		return "", "", false
	}

	return aseSteamInfo.Executable, aseSteamInfo.AppID, true
}

func initConfig() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to determine home directory: %w", err)
	}

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetConfigName(configName)
	viper.SetConfigType(configType)
	viper.AddConfigPath(homeDir)

	pwd, _ := os.Getwd()
	asePath, appId, fromSteam := TryFindAsepritePath()

	viper.SetDefault(fromSteamKey, fromSteam)
	viper.SetDefault(appIdKey, appId)
	viper.SetDefault(asepritePathKey, asePath)
	viper.SetDefault(scriptDirKey, filepath.Join(pwd, "scripts"))
	viper.SetDefault(spriteDirsKey, "")
	viper.SetDefault(palettesDirsKey, "")
	viper.SetDefault(openAiConfigKey, OpenAiConfig{
		ApiUrl: openAiApiUrl,
		ApiKey: os.Getenv("OPENAI_API_KEY"),
	})

	if err := viper.ReadInConfig(); err != nil {
		var configNotFound viper.ConfigFileNotFoundError
		if errors.As(err, &configNotFound) {
			if err := viper.SafeWriteConfig(); err != nil {
				return fmt.Errorf("failed to write initial config: %w", err)
			}
		}
	}

	return nil
}

func saveConfig() error {
	if err := viper.WriteConfig(); err != nil {
		return fmt.Errorf("failed to persist configuration: %w", err)
	}

	return nil
}
