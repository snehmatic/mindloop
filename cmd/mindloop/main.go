package main

import (
	"fmt"
	"os"

	"github.com/snehmatic/mindloop/internal/application"
	"github.com/snehmatic/mindloop/internal/infrastructure/config"
	"github.com/snehmatic/mindloop/internal/presentation/cli/handlers"
	"github.com/spf13/cobra"
)

func main() {
	if err := execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func execute() error {
	// Initialize configuration
	appConfig := config.GetConfig()

	// Validate user configuration
	if err := config.ValidateUserConfig(); err != nil {
		// Allow configure command to run without user config
		if len(os.Args) > 1 && os.Args[1] == "configure" {
			return executeConfigureCommand()
		}
		return fmt.Errorf("configuration validation failed: %w\nPlease run 'mindloop configure' to set up your profile", err)
	}

	// Load user configuration and update app config
	userConfig, err := config.LoadUserConfig()
	if err != nil {
		return fmt.Errorf("failed to load user config: %w", err)
	}

	// Update app config with user settings
	config.InitConfig(userConfig.Name, userConfig.Mode, "8080")
	appConfig = config.GetConfig()

	// Initialize dependency injection container
	container, err := application.NewContainer(appConfig)
	if err != nil {
		return fmt.Errorf("failed to initialize application: %w", err)
	}
	defer container.Close()

	// Create root command
	rootCmd := createRootCommand(container)

	return rootCmd.Execute()
}

func createRootCommand(container *application.Container) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:       "mindloop",
		Short:     "mindloop is a CLI tool for productivity tracking",
		Long:      `Mindloop helps track intent, focus sessions, and habits via CLI.`,
		Example:   `mindloop intent start "Get this work done"`,
		ValidArgs: []string{"intent", "focus", "habit", "journal", "summary", "configure"},
		Args:      cobra.OnlyValidArgs,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// Check database health
			if err := container.DB.Health(); err != nil {
				container.UI.ShowWarning("Database connection issue detected")
				container.Config.Logger.Warn().Err(err).Msg("Database health check failed")
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			container.UI.ShowBanner()
			container.UI.ShowRocket("Welcome to Mindloop! Use 'mindloop help' to see available commands.")
			container.UI.ShowInfo("For starters, try 'mindloop configure' to set up your profile.")
			container.Config.Logger.Info().Msg("User accessed root command, prompting for help.")
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			fmt.Println("🧠 Thank you for using Mindloop! This is still a work in progress.")
		},
	}

	// Add command handlers
	habitHandler := handlers.NewHabitHandler(container.HabitUseCase, container.UI)
	rootCmd.AddCommand(habitHandler.CreateCommands())
	intentHandler := handlers.NewIntentHandler(container.IntentUseCase, container.UI)
	rootCmd.AddCommand(intentHandler.CreateCommands())
	focusHandler := handlers.NewFocusHandler(container.FocusUseCase, container.UI)
	rootCmd.AddCommand(focusHandler.CreateCommands())
	journalHandler := handlers.NewJournalHandler(container.JournalUseCase, container.UI)
	rootCmd.AddCommand(journalHandler.CreateCommands())
	summaryHandler := handlers.NewSummaryHandler(container.SummaryUseCase, container.UI)
	rootCmd.AddCommand(summaryHandler.CreateCommands())

	// Add configure command
	rootCmd.AddCommand(createConfigureCommand(container))

	return rootCmd
}

func createConfigureCommand(container *application.Container) *cobra.Command {
	return &cobra.Command{
		Use:     "configure",
		Short:   "Configure your mindloop profile",
		Example: `mindloop configure`,
		RunE: func(cmd *cobra.Command, args []string) error {
			container.UI.ShowRocket("Welcome to Mindloop configuration!")

			fmt.Print("Please enter your preferred username: ")
			var username string
			fmt.Scanln(&username)

			var mode string
			for {
				fmt.Print("Please enter your preferred mode [local/byodb]: ")
				fmt.Scanln(&mode)
				if config.IsValidMode(mode) {
					break
				}
				container.UI.ShowWarning("Invalid mode. Please choose from: local, byodb.")
			}

			userConfig := &config.UserConfig{
				Name: username,
				Mode: mode,
			}

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
					dbName = "mindloop"
				}

				userConfig.DbConfig = config.DBConfig{
					Host:     dbHost,
					Port:     dbPort,
					User:     dbUser,
					Password: dbPass,
					Name:     dbName,
				}
			}

			if err := config.SaveUserConfig(userConfig); err != nil {
				container.UI.ShowError(fmt.Sprintf("Failed to save configuration: %v", err))
				return err
			}

			container.UI.ShowSuccess(fmt.Sprintf("Configuration complete! Your username is set to: %s, using mode: %s", username, mode))
			container.UI.ShowInfo("You can find your config as user_config.yaml")
			return nil
		},
	}
}

func executeConfigureCommand() error {
	// Simple configure command execution without full container initialization
	appConfig := config.GetConfig()
	container := &application.Container{
		Config: appConfig,
	}

	configCmd := createConfigureCommand(container)
	return configCmd.Execute()
}
