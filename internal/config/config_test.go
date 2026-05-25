//go:build test

//^ test tag is required here to resolve ResetForTest()

package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestSingletonLoadedOnce(t *testing.T) {
	ResetForTest()
	defer ResetForTest()

	// First initialization should succeed
	err := Init("test", "local", ":8080")
	require.NoError(t, err)

	// Second initialization should fail
	err2 := Init("test2", "local", ":8081")
	require.Error(t, err2)
	require.Contains(t, err2.Error(), "config already initialized")

	// Get should return the same instance every time
	c1 := GetConfig()
	c2 := GetConfig()
	require.Same(t, c1, c2)
	require.Equal(t, "test", c1.Name)
	require.Equal(t, ":8080", c1.Port)
	require.Equal(t, Local, c1.Mode)
}

func TestConcurrentAccess(t *testing.T) {
	ResetForTest()
	defer ResetForTest()

	// Initialize config
	err := Init("test", "local", ":8080")
	require.NoError(t, err)

	// Concurrent reads
	const numGoroutines = 100
	done := make(chan bool)
	var instances []*Config

	for i := 0; i < numGoroutines; i++ {
		go func() {
			instances = append(instances, GetConfig())
			done <- true
		}()
	}

	// Wait for all goroutines
	for i := 0; i < numGoroutines; i++ {
		<-done
	}

	// All instances should be the same
	first := instances[0]
	for _, inst := range instances {
		require.Same(t, first, inst)
	}
}

func TestCacheInvalidationOnChange(t *testing.T) {
	ResetForTest()
	defer ResetForTest()

	// Create a temporary directory and set it as the current directory for this test.
	tmpDir := t.TempDir()
	oldDir, err := os.Getwd()
	require.NoError(t, err)
	// Restore the original directory after the test.
	defer func() { _ = os.Chdir(oldDir) }()

	// Change to the temporary directory.
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	// Initialize with local mode (no DB config needed)
	err = Init("original", "local", ":8080")
	require.NoError(t, err)

	originalCfg := GetConfig()
	require.Equal(t, "original", originalCfg.Name)
	require.Equal(t, "", originalCfg.UserName) // no user config yet

	// Create initial user config in the current directory (which is tmpDir)
	initialUC := UserConfig{
		Name:            "initial_user",
		Mode:            "local",
		EditorWideWidth: false,
	}
	f, err := os.Create("user_config.yaml")
	require.NoError(t, err)
	defer f.Close()
	err = yaml.NewEncoder(f).Encode(&initialUC)
	require.NoError(t, err)

	// Reload should pick up the new user config
	err = Reload()
	require.NoError(t, err)

	reloadedCfg := GetConfig()
	require.Equal(t, reloadedCfg.Name, "original") // Name preserved from Init
	require.Equal(t, reloadedCfg.UserName, "initial_user")

	// Change user config
	changedUC := UserConfig{
		Name:            "changed_user",
		Mode:            "local",
		EditorWideWidth: true,
	}
	f, err = os.Create("user_config.yaml")
	require.NoError(t, err)
	defer f.Close()
	err = yaml.NewEncoder(f).Encode(&changedUC)
	require.NoError(t, err)

	// Track OnChange calls
	var changedCfg *Config
	OnChange(func(cfg *Config) {
		changedCfg = cfg
	})

	// Reload again
	err = Reload()
	require.NoError(t, err)

	// Check that the config was updated
	latestCfg := GetConfig()
	require.Equal(t, latestCfg.Name, "original")
	require.Equal(t, latestCfg.UserName, "changed_user")

	// Check that OnChange was called with the updated config
	require.NotNil(t, changedCfg)
	require.Equal(t, changedCfg.UserName, "changed_user")
}

func TestInitForTest(t *testing.T) {
	ResetForTest()
	defer ResetForTest()

	// Use the test helper to initialize
	err := Init("test", "local", ":9090")
	require.NoError(t, err)

	c := GetConfig()
	require.Equal(t, "test", c.Name)
	require.Equal(t, ":9090", c.Port)

	// After the test, ResetForTest is called (via defer) which should set instance to nil.
	// Let's verify that after ResetForTest, we can re-initialize with different values.
	ResetForTest()
	err2 := Init("test2", "local", ":9091")
	require.NoError(t, err2)

	c2 := GetConfig()
	require.Equal(t, "test2", c2.Name)
	require.Equal(t, ":9091", c2.Port)
}
