package config_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type DefaultWrapper struct {
	testFunc func(t *testing.T, cfg *config.Config)
}

func (w DefaultWrapper) testWrapped(t *testing.T) {
	cfg := SetupTestConfig(t)

	w.testFunc(t, cfg)

	assertConfigExists(t)
}

// SetupTestConfig configures a temporary USERPROFILE directory
func SetupTestConfig(t *testing.T) (cfg *config.Config) {
	viper.Reset()
	tempDir := t.TempDir()

	t.Setenv("USERPROFILE", tempDir)

	cfg, err := config.LoadConfig()
	require.NoError(t, err)

	return cfg
}

func TestSavePaths(t *testing.T) {
	w := DefaultWrapper{
		testFunc: func(t *testing.T, _ *config.Config) {

			expectedAsepritePath := "/fake/aseprite.exe"
			expectedSpritesPaths := []string{"/assets1", "/assets2"}
			expectedPalettesDirs := []string{"/palettes1", "/palettes2"}

			err := config.SavePaths(expectedAsepritePath, expectedSpritesPaths, expectedPalettesDirs)
			require.NoError(t, err)

			cfg, err := config.LoadConfig()
			require.NoError(t, err)

			assert.Equal(t, expectedAsepritePath, cfg.AsepritePath)
			assert.Equal(t, expectedSpritesPaths, cfg.SpritesFoldersPaths)
			assert.Equal(t, expectedPalettesDirs, cfg.PalettesFoldersPaths)

		},
	}
	w.testWrapped(t)
}

func TestSetScriptsDirPath(t *testing.T) {
	w := DefaultWrapper{
		testFunc: func(t *testing.T, _ *config.Config) {
			pwd, err := os.Getwd()
			require.NoError(t, err)
			expected := filepath.Join(filepath.Join(pwd, "outer"), "scripts")

			err = config.SetScriptDirPath(expected)
			require.NoError(t, err)

			cfg, err := config.LoadConfig()
			require.NoError(t, err)

			assert.Equal(t, expected, cfg.ScriptDirPath)
		},
	}
	w.testWrapped(t)
}

func TestSetDefaultAfterChangingScriptsDirPath(t *testing.T) {
	w := DefaultWrapper{
		testFunc: func(t *testing.T, _ *config.Config) {
			pwd, err := os.Getwd()
			require.NoError(t, err)
			preDefault := filepath.Join(filepath.Join(pwd, "outer"), "scripts")

			err = config.SetScriptDirPath(preDefault)
			require.NoError(t, err)

			err = config.SetDefaultScriptDirPath()
			require.NoError(t, err)

			cfg, err := config.LoadConfig()
			require.NoError(t, err)

			expected := filepath.Join(pwd, "scripts")
			assert.NotEqual(t, preDefault, cfg.ScriptDirPath)
			assert.Equal(t, expected, cfg.ScriptDirPath)
		},
	}
	w.testWrapped(t)
}

func TestSetOpenAiConfig(t *testing.T) {
	w := DefaultWrapper{
		testFunc: func(t *testing.T, _ *config.Config) {
			apiKey := "new-key"
			apiUrl := "new-url"

			err := config.SetOpenAiConfig(apiKey, apiUrl)
			require.NoError(t, err)

			cfg, err := config.LoadConfig()
			require.NoError(t, err)

			assert.Equal(t, apiKey, cfg.OpenAiConfig.ApiKey)
			assert.Equal(t, apiUrl, cfg.OpenAiConfig.ApiUrl)
		},
	}
	w.testWrapped(t)
}

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  *config.Config
		wantErr string
	}{
		{
			name:    "missing aseprite path",
			config:  &config.Config{AsepritePath: "", SpritesFoldersPaths: []string{"/assets"}, ScriptDirPath: "/scripts"},
			wantErr: "aseprite_path",
		},
		{
			name:    "missing assets paths",
			config:  &config.Config{AsepritePath: "D:\\aseprite.exe", SpritesFoldersPaths: nil, ScriptDirPath: "/scripts"},
			wantErr: "at least one sprite folder path required",
		},
		{
			name:    "missing script dir",
			config:  &config.Config{AsepritePath: "D:\\aseprite.exe", SpritesFoldersPaths: []string{"/assets"}, ScriptDirPath: ""},
			wantErr: "scripts_dir",
		},
		{
			name:    "invalid not absolute sprites paths",
			config:  &config.Config{AsepritePath: "D:\\aseprite.exe", SpritesFoldersPaths: []string{"/assets"}, ScriptDirPath: "D:\\scripts"},
			wantErr: "sprites path must be absolute",
		},
		{
			name:    "valid all paths are absolute",
			config:  &config.Config{AsepritePath: "D:\\aseprite.exe", SpritesFoldersPaths: []string{"D:\\assets"}, ScriptDirPath: "D:\\scripts", PalettesFoldersPaths: []string{"D:\\palettes"}},
			wantErr: "",
		},
		{
			name:    "invalid aseprite path",
			config:  &config.Config{AsepritePath: "D:\\aseprite", SpritesFoldersPaths: []string{"D:\\assets"}, ScriptDirPath: "D:\\scripts", PalettesFoldersPaths: []string{"D:\\palettes"}},
			wantErr: "aseprite path must be absolute executable path",
		},
		{
			name:    "invalid palette paths if not empty and are not absolute",
			config:  &config.Config{AsepritePath: "D:\\aseprite.exe", SpritesFoldersPaths: []string{"D:\\assets"}, ScriptDirPath: "D:\\scripts", PalettesFoldersPaths: []string{"/palettes"}},
			wantErr: "palettes path must be absolute",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tt.wantErr)
			}

		})
	}
}

func TestInfo(t *testing.T) {
	w := DefaultWrapper{
		testFunc: func(t *testing.T, _ *config.Config) {
			expectedAsepritePath := "/aseprite.exe"

			err := config.SavePaths(expectedAsepritePath, []string{"/assets"}, []string{"/palettes"})
			require.NoError(t, err)

			info := config.Info()

			homeDir := os.Getenv("USERPROFILE")
			expectedConfigPath := filepath.Join(homeDir, ".aseprite-assets-cli.json")
			assert.Contains(t, info, fmt.Sprintf("Config file used: %s", expectedConfigPath))
			assert.Contains(t, info, fmt.Sprintf("aseprite_path: %s", expectedAsepritePath))
		},
	}
	w.testWrapped(t)
}

func TestDefaultAsepritePathIfEnvExists(t *testing.T) {
	tempDir := t.TempDir()
	t.Setenv("USERPROFILE", tempDir)

	t.Setenv("ASEPRITE", "/env/aseprite")

	viper.Reset()
	cfg, err := config.LoadConfig()
	require.NoError(t, err)

	assert.Equal(t, "/env/aseprite", cfg.AsepritePath)
}

func TestOpenAiConfigDefaultFromEnv(t *testing.T) {
	tempDir := t.TempDir()
	t.Setenv("USERPROFILE", tempDir)

	t.Setenv("OPENAI_API_KEY", "test-key")

	viper.Reset()
	cfg, err := config.LoadConfig()
	require.NoError(t, err)

	assert.Equal(t, "test-key", cfg.OpenAiConfig.ApiKey)
}

func TestNotErrorWhenConfigFileHasInvalidContent(t *testing.T) {
	tempDir := t.TempDir()
	t.Setenv("USERPROFILE", tempDir)

	configPath := filepath.Join(tempDir, ".aseprite-assets-cli.json")

	err := os.WriteFile(configPath, []byte("{invalid}"), 0644)
	require.NoError(t, err)

	cfg, err := config.LoadConfig()

	assert.NotNil(t, cfg)
	assert.NoError(t, err)
}

func assertConfigExists(t *testing.T) {
	homeDir := os.Getenv("USERPROFILE")
	configPath := filepath.Join(homeDir, ".aseprite-assets-cli.json")
	assert.FileExists(t, configPath)
}
