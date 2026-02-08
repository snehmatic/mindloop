package cli

import (
	"fmt"

	"github.com/snehmatic/mindloop/internal/config"
	. "github.com/snehmatic/mindloop/internal/utils"
	"github.com/snehmatic/mindloop/models"
	"github.com/spf13/cobra"
)

// configure user command
var confCmd = &cobra.Command{
	Use:     "configure",
	Short:   "configure your mindloop profile",
	Example: `mindloop configure"`,
	Run: func(cmd *cobra.Command, args []string) {
		// Placeholder for configuration logic
		// This could involve setting user preferences, etc.
		_, _ = PrintRocketln("Welcome to Mindloop configuration!")
		fmt.Print("Please enter your preferred username: ")
		var username string
		_, _ = fmt.Scanln(&username)
		var mode string
		for {
			fmt.Print("Please enter your preferred mode [local/byodb]: ")
			_, _ = fmt.Scanln(&mode)
			if models.IsValidMode(mode) {
				break
			}
			_, _ = PrintWarnln("Invalid mode. Please choose from: local, byodb.")
		}

		dbConfig := &config.DBConfig{}
		if mode == "byodb" {
			fmt.Print("Please enter your database host name: ")
			var dbHost string
			_, _ = fmt.Scanln(&dbHost)
			fmt.Print("Please enter your database port: ")
			var dbPort string
			_, _ = fmt.Scanln(&dbPort)
			fmt.Print("Please enter your database user name: ")
			var dbUser string
			_, _ = fmt.Scanln(&dbUser)
			fmt.Print("Please enter your database password: ")
			var dbPass string
			_, _ = fmt.Scanln(&dbPass)
			fmt.Print("Please enter your database name [mindloop]: ")
			var dbName string
			_, _ = fmt.Scanln(&dbName)
			if dbName == "" {
				dbName = "mindloop" // default
			}
			dbConfig = &config.DBConfig{
				Host:     dbHost,
				Port:     dbPort,
				User:     dbUser,
				Password: dbPass,
				Name:     dbName,
			}
		}

		CreateUserConfigYAML(username, mode, dbConfig)

		_, _ = PrintSuccessf("Configuration complete! Your username is set to: %s, using mode: %s\n", username, mode)
	},
}

func init() {
	rootCmd.AddCommand(confCmd)
}

func CreateUserConfigYAML(username, mode string, dbConfig *config.DBConfig) {
	uc := config.UserConfig{
		Name: username,
		Mode: mode,
	}

	if mode == "byodb" {
		if dbConfig == nil {
			_, _ = PrintWarnln("Database configuration is required for 'byodb' mode. Please try again.")
			return
		}
		uc.DbConfig = *dbConfig
	}

	uc.WriteToYAML()
	_, _ = PrintSuccessln("User config created successfully!")
	_, _ = PrintInfoln("You can find your config as user_config.yaml")
}
