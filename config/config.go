package config

import (
	"fmt"
	"sync"

	"github.com/joho/godotenv"
	"github.com/snehmatic/mindloop/internal/utils"
	"github.com/snehmatic/mindloop/models"
)

// mindloop Application global configuration
type Config struct {
	Mode     models.MindloopMode
	Port     string
	Name     string
	DBConfig DBConfig
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

		// check if user_config.yaml exists
		if utils.FileExists(models.UserConfigPath) {
			fmt.Println("User config already exists at", models.UserConfigPath)
			return
		}
		config = &Config{
			Name:     name,
			Port:     port,
			Mode:     models.MindloopMode(mode),
			DBConfig: DBConfig{},
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

		fmt.Println("Mindloop global config has been set!")
	})
}

func GetConfig() *Config {
	return config
}
