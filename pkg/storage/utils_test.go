package storage

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetConfigDir(t *testing.T) {
	t.Run("should get the config dir from XDG", func(t *testing.T) {
		tmpDir := os.TempDir()
		t.Setenv("XDG_CONFIG_HOME", tmpDir)

		dir, err := getConfigDir()

		assert.NoError(t, err)
		assert.Equal(t, dir, filepath.Join(tmpDir, "gostman"))
	})

	t.Run("should get the config dir using the home dir if XDG is not set", func(t *testing.T) {
		t.Setenv("XDG_CONFIG_HOME", "")
		homedir, _ := os.UserHomeDir()

		dir, err := getConfigDir()

		assert.NoError(t, err)
		assert.Equal(t, dir, homedir+"/.config/gostman")
	})

	t.Run("should return an error if UserHomeDir fails", func(t *testing.T) {
		t.Setenv("XDG_CONFIG_HOME", "")
		t.Setenv("HOME", "")

		dir, err := getConfigDir()

		assert.EqualError(t, err, "$HOME is not defined")
		assert.Empty(t, dir)
	})
}
