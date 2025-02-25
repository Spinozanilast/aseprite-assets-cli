package commands

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spinozanilast/aseprite-assets-cli/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigCommand(t *testing.T) {
	t.Run("SetValidScriptsDir", func(t *testing.T) {
		tmpHome := t.TempDir()
		t.Setenv("HOME", tmpHome)

		_, err := config.LoadConfig()
		require.NoError(t, err)

		scriptsDir := filepath.Join(tmpHome, "scripts")
		require.NoError(t, os.Mkdir(scriptsDir, 0755))

		cmd := configCmd
		cmd.SetArgs([]string{"edit", "--scripts-dir", scriptsDir})
		require.NoError(t, cmd.Execute())

		cfg, err := config.LoadConfig()
		require.NoError(t, err)
		assert.Equal(t, scriptsDir, cfg.ScriptDirPath)
	})

	t.Run("SetDefaultScriptsDir", func(t *testing.T) {
		tmpHome := t.TempDir()
		t.Setenv("HOME", tmpHome)

		cmd := configCmd
		cmd.SetArgs([]string{"edit", "--scripts-dir", "default"})
		require.NoError(t, cmd.Execute())

		cfg, err := config.LoadConfig()
		require.NoError(t, err)
		pwd, err := os.Getwd()
		require.NoError(t, err)
		expected := filepath.Join(pwd, "scripts")
		assert.Equal(t, expected, cfg.ScriptDirPath)
	})
}
