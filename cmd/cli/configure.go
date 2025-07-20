package cli

import (
	"fmt"

	"github.com/snehmatic/mindloop/internal/config"
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
		cmd.Println("Welcome to Mindloop configuration!")
		fmt.Print("Please enter your preferred username: ")
		var username string
		fmt.Scanln(&username)
		var mode string
		for {
			fmt.Print("Please enter your preferred mode [local/byodb]: ")
			fmt.Scanln(&mode)
			if models.IsValidMode(mode) {
				break
			}
			cmd.Println("Invalid mode. Please choose from: local, byodb.")
		}

		dbConfig := &config.DBConfig{}
		if mode == "byodb" {
			fmt.Print("Please enter your database host name: ")
			var dbHost string
			fmt.Scanln(&dbHost)
			fmt.Print("Please enter your database port: ")
			var dbPort string
			fmt.Scanln(&dbPort)
			fmt.Print("Please enter your database user name: ")
			var dbUser string
			fmt.Scanln(&dbUser)
			fmt.Print("Please enter your database password: ")
			var dbPass string
			fmt.Scanln(&dbPass)
			fmt.Print("Please enter your database name [mindloop]: ")
			var dbName string
			fmt.Scanln(&dbName)
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

		cmd.Printf("Configuration complete! Your username is set to: %s, using mode: %s\n", username, mode)
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
			fmt.Println("Database configuration is required for 'byodb' mode. Please try again.")
			return
		}
		uc.DbConfig = *dbConfig
	}

	uc.WriteToYAML()
	fmt.Println("User config created successfully!")
	fmt.Println("You can find your config as user_config.yaml")
}
