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

// UserConfigPath is the file path where the user configuration YAML will be written.
// ToDo: Make this configurable or use a constant
var UserConfigPath = "user_config.yaml"

type MindloopMode string

var AllModes = [...]string{"local", "byodb", "api"}

var (
	Local MindloopMode = MindloopMode(AllModes[0])
	ByoDB MindloopMode = MindloopMode(AllModes[1])
	Api   MindloopMode = MindloopMode(AllModes[2])
)

// mindloop Application global configuration
type Config struct {
	Mode     MindloopMode
	Port     string
	Name     string
	DBConfig DBConfig
	Logger   zerolog.Logger
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

var once sync.Once
var config *Config

func InitConfig(name, mode, port string) {
	once.Do(func() { // singleton

		config = &Config{
			Name:     name,
			Port:     port,
			Mode:     MindloopMode(mode),
			DBConfig: DBConfig{},
			Logger:   log.Get(),
		}

		if mode == "api" {
			// init DB Config
			err := godotenv.Load()
			if err != nil {
				fmt.Printf("error loading .env file: %v\n", err)
			}
			config.DBConfig = DBConfig{
				Host:     utils.GetEnvOrDie("DB_HOST"),
				Port:     utils.GetEnvOrDie("DB_PORT"),
				User:     utils.GetEnvOrDie("DB_USER"),
				Password: utils.GetEnvOrDie("DB_PASS"),
				Name:     utils.GetEnvOrDie("DB_NAME"),
			}
		}

		config.Logger.Info().Msg("Mindloop global config has been set!")
	})
}

func GetConfig() *Config {
	return config
}

type UserConfig struct {
	Name     string   `yaml:"name"`
	Mode     string   `yaml:"mode"`
	DbConfig DBConfig `yaml:"db_config"`
}

func ValidateUserConfig(cmd *cobra.Command) {
	// check if user_config.yaml exists
	logger := log.Get()
	if utils.FileExists(UserConfigPath) {
		logger.Debug().Msgf("User config exists at %s", UserConfigPath)
	} else {
		if cmd.Use != "configure" {
			fmt.Println("Warn: user config does not exist, create a new one or run `mindloop configure`.")
			logger.Warn().Msg("User config does not exist, warned user")
			os.Exit(0)
		}
	}
}

func (uc UserConfig) WriteToYAML() {
	marshalled, err := yaml.Marshal(uc)
	if err != nil {
		fmt.Println("Error marshalling user config to YAML")
		return
	}
	err = os.WriteFile(UserConfigPath, marshalled, 0644)
	if err != nil {
		fmt.Println("Error writing user config to file")
		return
	}
	fmt.Println("User config written to YAML successfully")
}

func (uc *UserConfig) ReadFromYAML() error {
	data, err := os.ReadFile(UserConfigPath)
	if err != nil {
		return fmt.Errorf("failed to read user config file: %w", err)
	}
	err = yaml.Unmarshal(data, uc)
	if err != nil {
		return fmt.Errorf("failed to unmarshal user config: %w", err)
	}
	return nil
}
