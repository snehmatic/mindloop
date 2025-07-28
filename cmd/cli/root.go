package cli

import (
	"fmt"
	"os"

	"gorm.io/gorm"

	"github.com/snehmatic/mindloop/db"
	"github.com/snehmatic/mindloop/internal/config"
	"github.com/spf13/cobra"
)

var gdb *gorm.DB
var ac *config.Config

var rootCmd = &cobra.Command{
	Use:       "mindloop",
	Short:     "mindloop is a CLI tool for productivity tracking",
	Long:      `Mindloop helps track intent, focus sessions, and habits via CLI.`,
	Example:   `mindloop intent start "Get this work done"`,
	ValidArgs: []string{"intent", "focus", "habit", "log", "stats"},
	Args:      cobra.OnlyValidArgs,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		config.ValidateUserConfig(cmd)

		if db.LocalDBFileExists() {
			ac.Logger.Info().Msg("Found local DB file, using it for local mode.")
		} else {
			ac.Logger.Warn().Msg("No local DB file found, a new one will be created.")
		}
	},
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
	ac = config.GetConfig()
	// Initialize local db
	db, err := db.ConnectToDb(*ac)
	if err != nil {
		fmt.Printf("Error connecting to DB: %v\n", err)
		ac.Logger.Error().Msgf("Error connecting to DB: %v", err)
		fmt.Println("Please check your database connection or configuration.")
		ac.Logger.Warn().Msg("Exiting due to DB connection error.")
		os.Exit(1)
	}
	gdb = db
}
