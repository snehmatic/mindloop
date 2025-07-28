package cli

import (
	"time"

	"github.com/snehmatic/mindloop/internal/utils"
	. "github.com/snehmatic/mindloop/internal/utils"
	"github.com/snehmatic/mindloop/models"
	"github.com/spf13/cobra"
)

// parent intent command
var intentCmd = &cobra.Command{
	Use:     "intent",
	Short:   "Manage your intents",
	Example: `mindloop intent start "Get this work done"`,
}

// start intent subcommand
var intentStartCmd = &cobra.Command{
	Use:     "start",
	Short:   "Start a new intent",
	Example: `mindloop intent start "Get this work done"`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		intent := &models.Intent{
			Name:   args[0],
			Status: "active",
		}
		// start the intent
		if err := gdb.Create(intent).Error; err != nil {
			PrintErrorln("Error starting intent:", err)
			ac.Logger.Error().Msgf("Error starting intent: %v", err)
			PrintInfoln("Please try again or check your database connection.")
			return
		}
		PrintSuccessf("Intent '%s' started successfully with id %d!\n", intent.Name, intent.ID)
		ac.Logger.Info().Msgf("Intent '%s' started successfully with id %d!", intent.Name, intent.ID)
	},
}

// list intent subcommand
var intentListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List all intents",
	Example: `mindloop intent list`,
	Run: func(cmd *cobra.Command, args []string) {
		var intents []models.Intent
		if err := gdb.Find(&intents).Error; err != nil {
			PrintErrorln("Error fetching intents:", err)
			ac.Logger.Error().Msgf("Error fetching intents: %v", err)
			PrintInfoln("Please check your database connection or try again later.")
			return
		}
		if len(intents) == 0 {
			PrintInfoln("No intents found... Try starting one with 'mindloop intent start <name>'")
			ac.Logger.Info().Msg("No intents found. Prompting user to start a new intent.")
			return
		}

		views := []models.IntentView{}
		for _, i := range intents {
			views = append(views, models.ToIntentView(i))
		}
		utils.PrintTable(views)
		ac.Logger.Info().Msgf("Listed %d intents successfully.", len(intents))
	},
}

// current intent subcommand
var intentCurrentCmd = &cobra.Command{
	Use:     "current",
	Short:   "Show current active intents",
	Example: `mindloop intent current`,
	Run: func(cmd *cobra.Command, args []string) {
		var intents []models.Intent
		if err := gdb.Where("Status = ?", "active").Find(&intents).Error; err != nil {
			PrintErrorln("Error fetching active intents:", err)
			ac.Logger.Error().Msgf("Error fetching active intents: %v", err)
			PrintInfoln("Please check your database connection or try again later.")
			return
		}
		if len(intents) == 0 {
			PrintInfoln("No active intents found. To get all intents, use 'mindloop intent list'")
			ac.Logger.Info().Msg("No active intents found. Prompting user to list all intents.")
			return
		}

		views := []models.IntentView{}
		for _, i := range intents {
			views = append(views, models.ToIntentView(i))
		}
		PrintTable(views)
		ac.Logger.Info().Msgf("Listed %d active intents successfully.", len(intents))
	},
}

// end intent subcommand
var intentEndCmd = &cobra.Command{
	Use:     "end",
	Short:   "End intent",
	Example: `mindloop intent end 10`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			PrintWarnln("Please provide the intent ID to end.")
			ac.Logger.Warn().Msg("No intent ID provided for ending intent.")
			return
		}
		var intent models.Intent
		if err := gdb.Where("id = ?", args[0]).First(&intent).Error; err != nil {
			PrintErrorln("Error fetching intent:", err)
			ac.Logger.Error().Msgf("Error fetching intent with ID %s: %v", args[0], err)
			return
		}
		if intent.ID == 0 {
			PrintWarnln("No intent found with the given ID.")
			ac.Logger.Warn().Msgf("No intent found with ID %s to finish.", args[0])
			return
		}

		now := time.Now()
		intent.Status = "done"
		intent.EndedAt = &now
		if err := gdb.Save(&intent).Error; err != nil {
			PrintErrorln("Error ending intent:", err)
			ac.Logger.Error().Msgf("Error ending intent with ID %d: %v", intent.ID, err)
			return
		}
		ac.Logger.Info().Msgf("Intent '%s' ended successfully!", intent.Name)
		intentView := models.ToIntentView(intent)
		PrintTable([]models.IntentView{intentView})
	},
}

func init() {
	rootCmd.AddCommand(intentCmd)
	intentCmd.AddCommand(intentStartCmd)
	intentCmd.AddCommand(intentListCmd)
	intentCmd.AddCommand(intentCurrentCmd)
	intentCmd.AddCommand(intentEndCmd)
}
