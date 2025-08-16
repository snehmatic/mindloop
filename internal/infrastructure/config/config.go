package config

import (
	"fmt"
	"os"
	"slices"
	"sync"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/snehmatic/mindloop/internal/infrastructure/logging"
	"gopkg.in/yaml.v3"
)

type MindloopMode string

const (
	Local MindloopMode = "local"
	ByoDB MindloopMode = "byodb"
	Api   MindloopMode = "api"
)

var AllModes = []string{string(Local), string(ByoDB), string(Api)}

// Config represents the application configuration
type Config struct {
	Mode     MindloopMode `yaml:"mode"`
	Port     string       `yaml:"port"`
	Name     string       `yaml:"name"`
	DBConfig DBConfig     `yaml:"db_config"`
	Logger   zerolog.Logger
}

// DBConfig represents database configuration
type DBConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

// UserConfig represents user-specific configuration
type UserConfig struct {
	Name     string   `yaml:"name"`
	Mode     string   `yaml:"mode"`
	DbConfig DBConfig `yaml:"db_config"`
}

var (
	once     sync.Once
	instance *Config
)

// UserConfigPath is the file path where the user configuration YAML will be written
var UserConfigPath = "user_config.yaml"

// NewConfig creates a new configuration instance
func NewConfig(name, mode, port string) *Config {
	config := &Config{
		Name:     name,
		Port:     port,
		Mode:     MindloopMode(mode),
		DBConfig: DBConfig{},
		Logger:   logging.GetLogger(),
	}

	if mode == string(Api) || mode == string(ByoDB) {
		// Load environment variables for database configuration
		if err := godotenv.Load(); err != nil {
			config.Logger.Warn().Err(err).Msg("Could not load .env file")
		}

		config.DBConfig = DBConfig{
			Host:     getEnvOrDefault("DB_HOST", "localhost"),
			Port:     getEnvOrDefault("DB_PORT", "5432"),
			User:     getEnvOrDefault("DB_USER", "postgres"),
			Password: getEnvOrDefault("DB_PASS", ""),
			Name:     getEnvOrDefault("DB_NAME", "mindloop"),
		}
	}

	return config
}

// InitConfig initializes the global configuration (singleton pattern)
func InitConfig(name, mode, port string) {
	once.Do(func() {
		instance = NewConfig(name, mode, port)
		instance.Logger.Info().Msg("Mindloop global config has been set!")
	})
}

// GetConfig returns the global configuration instance
func GetConfig() *Config {
	if instance == nil {
		// Initialize with defaults if not already initialized
		InitConfig("mindloop", string(Local), "8080")
	}
	return instance
}

// IsValidMode checks if the given mode is valid
func IsValidMode(mode string) bool {
	return slices.Contains(AllModes, mode)
}

// LoadUserConfig loads user configuration from YAML file
func LoadUserConfig() (*UserConfig, error) {
	if !FileExists(UserConfigPath) {
		return nil, fmt.Errorf("user config file does not exist at %s", UserConfigPath)
	}

	data, err := os.ReadFile(UserConfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read user config file: %w", err)
	}

	var userConfig UserConfig
	if err := yaml.Unmarshal(data, &userConfig); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user config: %w", err)
	}

	return &userConfig, nil
}

// SaveUserConfig saves user configuration to YAML file
func SaveUserConfig(userConfig *UserConfig) error {
	data, err := yaml.Marshal(userConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal user config: %w", err)
	}

	if err := os.WriteFile(UserConfigPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write user config file: %w", err)
	}

	return nil
}

// ValidateUserConfig checks if user configuration exists and is valid
func ValidateUserConfig() error {
	if !FileExists(UserConfigPath) {
		return fmt.Errorf("user config does not exist, please run 'mindloop configure'")
	}

	userConfig, err := LoadUserConfig()
	if err != nil {
		return fmt.Errorf("invalid user config: %w", err)
	}

	if !IsValidMode(userConfig.Mode) {
		return fmt.Errorf("invalid mode in user config: %s", userConfig.Mode)
	}

	return nil
}

// getEnvOrDefault gets an environment variable or returns a default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// FileExists checks if a file exists
func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}
