package config

import (
	"fmt"
	"sync"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/snehmatic/mindloop/internal/log"
	"github.com/snehmatic/mindloop/internal/utils"
	"github.com/snehmatic/mindloop/models"
)

// mindloop Application global configuration
type Config struct {
	Mode     models.MindloopMode
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
			Mode:     models.MindloopMode(mode),
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
