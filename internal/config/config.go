// Package config manages application-wide and user-specific configurations
package config

import (
	"fmt"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/snehmatic/mindloop/internal/log"
	"github.com/snehmatic/mindloop/internal/utils"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// GetUserConfigPath returns the path to the user configuration file
func GetUserConfigPath() string {
	localFile := "user_config.yaml"
	if _, err := os.Stat(localFile); err == nil {
		return localFile
	}
	return GetDataDir() + "/" + localFile
}

// Version is set during build time via ldflags
var Version = "dev"

// GetDataDir returns the directory where mindloop data is stored
func GetDataDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "."
	}
	dir := homeDir + "/.mindloop"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		_ = os.MkdirAll(dir, 0755)
	}
	return dir
}

// MindloopMode defines the operating mode of the application
type MindloopMode string

// AllModes contains the supported operating modes
var AllModes = [...]string{"local", "byodb", "api"}

var (
	Local = MindloopMode(AllModes[0])
	ByoDB = MindloopMode(AllModes[1])
	API   = MindloopMode(AllModes[2])
)

// Config mindloop Application global configuration
type Config struct {
	Mode     MindloopMode
	Port     string
	Name     string
	UserName string
	DBConfig DBConfig
	Logger   zerolog.Logger
}

// DBConfig holds database connection parameters
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

var once sync.Once
var cfg *Config

// InitConfig initializes the global application configuration
func InitConfig(name, mode, port string) {
	once.Do(func() { // singleton
		cfg = &Config{
			Name:     name,
			Port:     port,
			Mode:     MindloopMode(mode),
			DBConfig: DBConfig{},
			Logger:   log.Get(),
		}

		// Try to load user config
		uc := UserConfig{}
		if err := uc.ReadFromYAML(); err == nil {
			if uc.Name != "" {
				cfg.UserName = uc.Name
			}
			// Override mode if set in user config and not explicitly overridden by flag (which passed here)
			// For simplicity, we are not overriding mode here as it might conflict with flags
			// But we can check if DBConfig is needed
			if uc.Mode == "byodb" {
				cfg.DBConfig = uc.DbConfig
			}
		}

		if mode == "api" {
			// init DB Config
			err := godotenv.Load()
			if err != nil {
				fmt.Printf("error loading .env file: %v\n", err)
			}
			cfg.DBConfig = DBConfig{
				Host:     utils.GetEnvOrDie("DB_HOST"),
				Port:     utils.GetEnvOrDie("DB_PORT"),
				User:     utils.GetEnvOrDie("DB_USER"),
				Password: utils.GetEnvOrDie("DB_PASS"),
				Name:     utils.GetEnvOrDie("DB_NAME"),
			}
		}

		cfg.Logger.Info().Msg("Mindloop global config has been set!")
	})
}

// GetConfig returns the global configuration object
func GetConfig() *Config {
	return cfg
}

// UserConfig represents the persistent user preferences
type UserConfig struct {
	Name            string       `yaml:"name"`
	Mode            string       `yaml:"mode"`
	EditorWideWidth bool         `yaml:"editor_wide_width"`
	DbConfig        DBConfig     `yaml:"db_config"`
	FeatureFlags    FeatureFlags `yaml:"feature_flags"`
	PointsConfig    PointsConfig `yaml:"points_config"`
}

// FeatureFlags toggles specific application functionalities
type FeatureFlags struct {
	FocusCloud   bool `yaml:"focus_cloud"`
	HabitCloud   bool `yaml:"habit_cloud"`
	IntentCloud  bool `yaml:"intent_cloud"`
	JournalCloud bool `yaml:"journal_cloud"`
	NoteCloud    bool `yaml:"note_cloud"`
	// Gamification toggles the points and celebration system
	Gamification bool `yaml:"gamification"`
}

// PointsConfig stores the user-defined point values for each activity type
type PointsConfig struct {
	Focus   int `yaml:"focus"`
	Habit   int `yaml:"habit"`
	Intent  int `yaml:"intent"`
	Journal int `yaml:"journal"`
	Quest   int `yaml:"quest"`
}

// SetDefaults ensures that the UserConfig has sensible default values
func (uc *UserConfig) SetDefaults() {
	// Feature flags defaults
	// Note: in YAML, a missing bool is false.
	// If we want it true by default, we'd need a more complex check or just assume if file is missing.
	// For now, let's just force it to true if the config is newly initialized or missing the key.
	// Actually, if we just want it 'default enabled' for new users:
	if !uc.FeatureFlags.Gamification && uc.PointsConfig.Focus == 0 {
		uc.FeatureFlags.Gamification = true
	}

	if uc.PointsConfig.Focus == 0 {
		uc.PointsConfig.Focus = 10
	}
	if uc.PointsConfig.Habit == 0 {
		uc.PointsConfig.Habit = 5
	}
	if uc.PointsConfig.Intent == 0 {
		uc.PointsConfig.Intent = 10
	}
	if uc.PointsConfig.Journal == 0 {
		uc.PointsConfig.Journal = 5
	}
	if uc.PointsConfig.Quest == 0 {
		uc.PointsConfig.Quest = 5
	}
}

// ValidateUserConfig checks if the user configuration is valid and exists
func ValidateUserConfig(cmd *cobra.Command) {
	// check if user_config.yaml exists
	logger := log.Get()
	configPath := GetUserConfigPath()
	if utils.FileExists(configPath) {
		logger.Debug().Msgf("User config exists at %s", configPath)
	} else if cmd.Use != "configure" {
		utils.PrintWarnln("Warn: user config does not exist, create a new one or run `mindloop configure`.")
		logger.Warn().Msg("User config does not exist, warned user")
		os.Exit(0)
	}
}

// WriteToYAML persists the current UserConfig to a YAML file
func (uc UserConfig) WriteToYAML() {
	marshalled, err := yaml.Marshal(uc)
	if err != nil {
		utils.PrintErrorln("Error marshalling user config to YAML")
		return
	}
	err = os.WriteFile(GetUserConfigPath(), marshalled, 0644)
	if err != nil {
		utils.PrintErrorln("Error writing user config to file")
		return
	}
	utils.PrintSuccessln("User config written to YAML successfully")
}

// ReadFromYAML loads the UserConfig from a YAML file
func (uc *UserConfig) ReadFromYAML() error {
	data, err := os.ReadFile(GetUserConfigPath())
	if err != nil {
		uc.SetDefaults() // Set defaults even if file doesn't exist
		return fmt.Errorf("failed to read user config file: %w", err)
	}
	err = yaml.Unmarshal(data, uc)
	if err != nil {
		return fmt.Errorf("failed to unmarshal user config: %w", err)
	}
	uc.SetDefaults() // Ensure defaults for missing fields
	return nil
}
