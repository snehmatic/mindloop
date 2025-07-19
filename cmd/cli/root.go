package cli

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"

	"github.com/snehmatic/mindloop/config"
	"github.com/snehmatic/mindloop/db"
	"github.com/spf13/cobra"
)

const (
	AppName = "MindLoop"
)

var rootCmd = &cobra.Command{
	Use:       "mindloop",
	Short:     "mindloop is a CLI tool for productivity tracking",
	Long:      `Mindloop helps track intent, focus sessions, and habits via CLI.`,
	Example:   `mindloop intent start "Get this work done"`,
	ValidArgs: []string{"intent", "focus", "habit", "log", "stats"},
	Args:      cobra.OnlyValidArgs,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	// define persistent flags here
}

func initConfig() {
	// Init global config
	config.InitConfig(AppName, "")
	appConfig := config.GetConfig()

	dbConnString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		appConfig.DBConfig.Host,
		appConfig.DBConfig.Port,
		appConfig.DBConfig.User,
		appConfig.DBConfig.Password,
		appConfig.DBConfig.Name,
	)

	_, err := db.Conn(dbConnString) // to be used later
	if err != nil {
		log.Fatal().Err(err).Msg("Error connecting to DB")
	}
}
