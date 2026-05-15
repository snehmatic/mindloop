package cli

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/snehmatic/mindloop/internal/config"
	"github.com/snehmatic/mindloop/internal/core/ai"
	"github.com/snehmatic/mindloop/internal/gamification"
	"github.com/snehmatic/mindloop/internal/utils"
	"github.com/snehmatic/mindloop/models"
	"github.com/spf13/cobra"
)

// configure user command
var confCmd = &cobra.Command{
	Use:     "configure",
	Short:   "configure your mindloop profile",
	Example: `mindloop configure"`,
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)
		utils.PrintRocketln("Welcome to Mindloop configuration!")

		username := Prompt(reader, "Please enter your preferred username: ", "")
		var mode string
		for {
			mode = Prompt(reader, "Please enter your preferred mode [local/byodb]: ", "local")
			if models.IsValidMode(mode) {
				break
			}
			utils.PrintWarnln("Invalid mode. Please choose from: local, byodb.")
		}

		dbConfig := &config.DBConfig{}
		if mode == "byodb" {
			dbHost := Prompt(reader, "Please enter your database host name: ", "")
			dbPort := Prompt(reader, "Please enter your database port: ", "")
			dbUser := Prompt(reader, "Please enter your database user name: ", "")
			dbPass := Prompt(reader, "Please enter your database password: ", "")
			dbName := Prompt(reader, "Please enter your database name [mindloop]: ", "mindloop")

			dbConfig = &config.DBConfig{
				Host:     dbHost,
				Port:     dbPort,
				User:     dbUser,
				Password: dbPass,
				Name:     dbName,
			}
		}

		milestoneInterval := gamification.DefaultMilestoneInterval
		milestoneInput := Prompt(reader, fmt.Sprintf("Please enter your milestone interval [%d]: ", gamification.DefaultMilestoneInterval), "")
		if milestoneInput != "" {
			parsedInterval, err := strconv.Atoi(milestoneInput)
			if err == nil && parsedInterval > 0 {
				milestoneInterval = parsedInterval
			} else {
				utils.PrintWarnf("Invalid milestone interval. Using default %d.\n", gamification.DefaultMilestoneInterval)
			}
		}

		CreateUserConfigYAML(username, mode, dbConfig, milestoneInterval)

		utils.PrintSuccessf("Configuration complete! Your username is set to: %s, using mode: %s\n", username, mode)

		// AI Configuration Prompt
		configureAI := Prompt(reader, "\nWould you like to configure your AI provider now? [y/N]: ", "n")
		configureAI = strings.ToLower(configureAI)

		if configureAI == "y" || configureAI == "yes" {
			aiSvc := ai.NewService(gdb)
			if err := SetupAIConfig(reader, os.Stdout, aiSvc); err != nil {
				utils.PrintErrorln("Failed to configure AI: " + err.Error())
			}
		}
	},
}

// Prompt displays a message and reads a line of input from the reader
func Prompt(reader *bufio.Reader, message, defaultValue string) string {
	_, _ = fmt.Print(message)
	input, err := reader.ReadString('\n')
	if err != nil {
		return defaultValue
	}
	input = strings.TrimSpace(input)
	if input == "" {
		return defaultValue
	}
	return input
}

// SetupAIConfig extracts the AI setup logic for testability
func SetupAIConfig(reader *bufio.Reader, writer io.Writer, aiSvc *ai.Service) error {
	var aiProvider, aiModel, aiToken, aiBaseURL string

	for {
		_, _ = fmt.Fprint(writer, "Provider [gemini/openai/custom]: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		aiProvider = strings.TrimSpace(input)
		if aiProvider == "gemini" || aiProvider == "openai" || aiProvider == "custom" {
			break
		}
		_, _ = fmt.Fprintln(writer, "Invalid provider. Please choose from: gemini, openai, custom.")
	}

	if aiProvider == "custom" {
		for {
			_, _ = fmt.Fprint(writer, "Base URL (e.g. http://localhost:11434/v1): ")
			input, err := reader.ReadString('\n')
			if err != nil {
				return err
			}
			aiBaseURL = strings.TrimSpace(input)
			if strings.HasPrefix(aiBaseURL, "http://") || strings.HasPrefix(aiBaseURL, "https://") {
				break
			}
			_, _ = fmt.Fprintln(writer, "Invalid Base URL. Must start with http:// or https://")
		}
	}

	for {
		_, _ = fmt.Fprint(writer, "Model (e.g. gpt-4o-mini or llama3): ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		aiModel = strings.TrimSpace(input)
		if aiModel != "" {
			break
		}
		_, _ = fmt.Fprintln(writer, "Model name cannot be empty.")
	}

	_, _ = fmt.Fprint(writer, "API Token (Type 'none' or leave blank if using local without auth): ")
	input, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	aiToken = strings.TrimSpace(input)
	if aiToken == "none" {
		aiToken = ""
	}

	err = aiSvc.SaveSettings(aiProvider, aiModel, aiToken, aiBaseURL)
	if err != nil {
		return err
	}
	_, _ = fmt.Fprintln(writer, "AI Configuration saved successfully!")
	return nil
}

func init() {
	rootCmd.AddCommand(confCmd)
}

func CreateUserConfigYAML(username, mode string, dbConfig *config.DBConfig, milestoneInterval int) {
	uc := config.UserConfig{
		Name: username,
		Mode: mode,
		PointsConfig: config.PointsConfig{
			MilestoneInterval: milestoneInterval,
		},
	}

	if mode == "byodb" {
		if dbConfig == nil {
			utils.PrintWarnln("Database configuration is required for 'byodb' mode. Please try again.")
			return
		}
		uc.DbConfig = *dbConfig
	}

	uc.SetDefaults()
	uc.WriteToYAML()
	utils.PrintSuccessln("User config created successfully!")
	utils.PrintInfof("You can find your config at: %s\n", config.GetUserConfigPath())
}
