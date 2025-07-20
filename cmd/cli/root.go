package cli

import (
	"fmt"
	"os"

	"gorm.io/gorm"

	"github.com/snehmatic/mindloop/config"
	"github.com/snehmatic/mindloop/db"
	"github.com/snehmatic/mindloop/internal/utils"
	"github.com/spf13/cobra"
)

const (
	AppName = "Mindloop"
)

var gdb *gorm.DB

var rootCmd = &cobra.Command{
	Use:       "mindloop",
	Short:     "mindloop is a CLI tool for productivity tracking",
	Long:      `Mindloop helps track intent, focus sessions, and habits via CLI.`,
	Example:   `mindloop intent start "Get this work done"`,
	ValidArgs: []string{"intent", "focus", "habit", "log", "stats"},
	Args:      cobra.OnlyValidArgs,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		utils.ValidateUserConfig(cmd)
		appConfig := config.GetConfig()

		if db.LocalDBFileExists() {
			fmt.Println("Found local DB file, using it for local mode.")
		} else {
			fmt.Println("No local DB file found, a new one will be created.")
		}

		// Initialize local db
		_, err := db.ConnectToDb(*appConfig) // to be used later
		if err != nil {
			fmt.Printf("Error connecting to DB: %v\n", err)
			os.Exit(1)
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
	// Init global config
	config.InitConfig(AppName, "local", "")
}
