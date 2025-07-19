package config

import (
	"fmt"
	"sync"

	"github.com/joho/godotenv"
	"github.com/snehmatic/mindloop/internal/utils"
)

// mindloop Application global configuration
type Config struct {
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

func InitConfig(name, port string) {
	once.Do(func() { // singleton

		// init DB Config
		err := godotenv.Load()
		if err != nil {
			fmt.Printf("error loading .env file: %v\n", err)
		}

		config = &Config{
			Name: name,
			Port: port,
			DBConfig: DBConfig{
				Host:     utils.GetEnvOrDie("DB_HOST"),
				Port:     utils.GetEnvOrDie("DB_PORT"),
				User:     utils.GetEnvOrDie("DB_USER"),
				Password: utils.GetEnvOrDie("DB_PASS"),
				Name:     utils.GetEnvOrDie("DB_NAME"),
			},
		}
		fmt.Println("Mindloop global config has been set!")
	})
}

func GetConfig() *Config {
	return config
}
