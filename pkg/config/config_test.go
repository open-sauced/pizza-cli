package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	t.Parallel()

	t.Run("Existing file", func(t *testing.T) {
		t.Parallel()
		tmpDir := t.TempDir()
		configFilePath := filepath.Join(tmpDir, ".sauced.yaml")

		fileContents := `# Configuration for attributing commits with emails to GitHub user profiles
# Used during codeowners generation.
# List the emails associated with the given username.
# The commits associated with these emails will be attributed to
# the username in this yaml map. Any number of emails may be listed.
attribution:
  brandonroberts:
    - robertsbt@gmail.com
  jpmcb:
    - john@opensauced.pizza`

		require.NoError(t, os.WriteFile(configFilePath, []byte(fileContents), 0644))

		config, err := LoadConfig(configFilePath, "")
		assert.NoError(t, err)
		assert.NotNil(t, config)

		// Assert that config contains all the Attributions in fileContents
		assert.Equal(t, 2, len(config.Attributions))

		// Check specific attributions
		assert.Equal(t, []string{"robertsbt@gmail.com"}, config.Attributions["brandonroberts"])
		assert.Equal(t, []string{"john@opensauced.pizza"}, config.Attributions["jpmcb"])
	})

	t.Run("Non-existent file", func(t *testing.T) {
		t.Parallel()
		tmpDir := t.TempDir()
		nonExistentPath := filepath.Join(tmpDir, ".sauced.yaml")

		config, err := LoadConfig(nonExistentPath, "")
		assert.Error(t, err)
		assert.Nil(t, config)
	})

	t.Run("Non-existent file with fallback", func(t *testing.T) {
		t.Parallel()
		fileContents := `# Configuration for attributing commits with emails to GitHub user profiles
# Used during codeowners generation.
# List the emails associated with the given username.
# The commits associated with these emails will be attributed to
# the username in this yaml map. Any number of emails may be listed.
attribution:
  brandonroberts:
    - robertsbt@gmail.com
  jpmcb:
    - john@opensauced.pizza
  nickytonline:
    - nick@nickyt.co
    - nick@opensauced.pizza
  zeucapua:
    - coding@zeu.dev`

		tmpDir := t.TempDir()
		fallbackPath := filepath.Join(tmpDir, ".sauced.yaml")
		require.NoError(t, os.WriteFile(fallbackPath, []byte(fileContents), 0644))

		// Print out the contents of the file we just wrote
		_, err := os.ReadFile(fallbackPath)
		require.NoError(t, err)

		nonExistentPath := filepath.Join(tmpDir, "non-existent.yaml")

		config, err := LoadConfig(nonExistentPath, fallbackPath)

		assert.NoError(t, err)
		assert.NotNil(t, config)

		// Assert that config contains all the Attributions in fileContents
		assert.Equal(t, 4, len(config.Attributions))

		// Check specific attributions
		assert.Equal(t, []string{"robertsbt@gmail.com"}, config.Attributions["brandonroberts"])
		assert.Equal(t, []string{"john@opensauced.pizza"}, config.Attributions["jpmcb"])
		assert.Equal(t, []string{"nick@nickyt.co", "nick@opensauced.pizza"}, config.Attributions["nickytonline"])
		assert.Equal(t, []string{"coding@zeu.dev"}, config.Attributions["zeucapua"])
	})

	//t.Run("Default path", func(t *testing.T) {
	//t.Parallel()
	//config, err := LoadConfig(DefaultConfigPath, "")
	//assert.Error(t, err)
	//assert.Nil(t, config)
	//})
}
