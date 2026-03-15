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

func GetUserConfigPath() string {
	localFile := "user_config.yaml"
	if _, err := os.Stat(localFile); err == nil {
		return localFile
	}
	return GetDataDir() + "/" + localFile
}

// Version is set during build time via ldflags
var Version = "dev"

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
	UserName string
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

		// Try to load user config
		uc := UserConfig{}
		if err := uc.ReadFromYAML(); err == nil {
			if uc.Name != "" {
				config.UserName = uc.Name
			}
			// Override mode if set in user config and not explicitly overridden by flag (which passed here)
			// For simplicity, we are not overriding mode here as it might conflict with flags
			// But we can check if DBConfig is needed
			if uc.Mode == "byodb" {
				config.DBConfig = uc.DbConfig
			}
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
	Name         string       `yaml:"name"`
	Mode         string       `yaml:"mode"`
	DbConfig     DBConfig     `yaml:"db_config"`
	FeatureFlags FeatureFlags `yaml:"feature_flags"`
}

type FeatureFlags struct {
        FocusCloud   bool `yaml:"focus_cloud"`
        HabitCloud   bool `yaml:"habit_cloud"`
        IntentCloud  bool `yaml:"intent_cloud"`
        JournalCloud bool `yaml:"journal_cloud"`
        NoteCloud    bool `yaml:"note_cloud"`
        Gamification bool `yaml:"gamification"`
}
func ValidateUserConfig(cmd *cobra.Command) {
	// check if user_config.yaml exists
	logger := log.Get()
	configPath := GetUserConfigPath()
	if utils.FileExists(configPath) {
		logger.Debug().Msgf("User config exists at %s", configPath)
	} else {
		if cmd.Use != "configure" {
			utils.PrintWarnln("Warn: user config does not exist, create a new one or run `mindloop configure`.")
			logger.Warn().Msg("User config does not exist, warned user")
			os.Exit(0)
		}
	}
}

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

func (uc *UserConfig) ReadFromYAML() error {
	data, err := os.ReadFile(GetUserConfigPath())
	if err != nil {
		return fmt.Errorf("failed to read user config file: %w", err)
	}
	err = yaml.Unmarshal(data, uc)
	if err != nil {
		return fmt.Errorf("failed to unmarshal user config: %w", err)
	}
	return nil
}
