package project

import (
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestCreateConfigFileGeneratesIndependentSecretsPerProject(t *testing.T) {
	keys := []string{"security.secret_key", "security.session_secret", "security.csrf_secret_key"}
	valuesByProject := make([]map[string]string, 0, 2)

	for i := 0; i < 2; i++ {
		projectPath := filepath.Join(t.TempDir(), "project")
		require.NoError(t, os.MkdirAll(filepath.Join(projectPath, "config"), 0755))
		require.NoError(t, createConfigFile(projectPath, "sqlite"))

		cfg := viper.New()
		cfg.SetConfigFile(filepath.Join(projectPath, "config", "config.yaml"))
		require.NoError(t, cfg.ReadInConfig())
		values := make(map[string]string, len(keys))
		for _, key := range keys {
			value := cfg.GetString(key)
			require.Len(t, value, 64)
			_, err := hex.DecodeString(value)
			require.NoError(t, err)
			require.False(t, strings.HasPrefix(strings.ToLower(strings.TrimSpace(value)), "change-me"))
			require.NotEqual(t, "secret", strings.ToLower(strings.TrimSpace(value)))
			require.NotEqual(t, "default", strings.ToLower(strings.TrimSpace(value)))
			values[key] = value
		}
		require.NotEqual(t, values[keys[0]], values[keys[1]])
		require.NotEqual(t, values[keys[0]], values[keys[2]])
		require.NotEqual(t, values[keys[1]], values[keys[2]])
		valuesByProject = append(valuesByProject, values)
	}

	for _, key := range keys {
		require.NotEqual(t, valuesByProject[0][key], valuesByProject[1][key])
	}
}
