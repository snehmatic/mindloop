package cli

import (
	"fmt"

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

		CreateUserConfigYAML(username, mode)

		cmd.Printf("Configuration complete! Your username is set to: %s, using mode: %s\n", username, mode)
	},
}

func init() {
	rootCmd.AddCommand(confCmd)
}

func CreateUserConfigYAML(username, mode string) {
	models.UserConfig{
		Name: username,
		Mode: mode,
	}.WriteToYAML()
	fmt.Println("User config created successfully!")
	fmt.Println("You can find your config as user_config.yaml")
}
