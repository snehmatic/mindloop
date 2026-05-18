package cli

import (
	"fmt"

	cfg "github.com/snehmatic/mindloop/internal/config"
	"github.com/snehmatic/mindloop/internal/core/intent"
	"github.com/snehmatic/mindloop/internal/utils"
	"github.com/snehmatic/mindloop/models"
	"github.com/spf13/cobra"
)

var (
	intentService *intent.Service
)

// parent intent command
var intentCmd = &cobra.Command{
	Use:     "intent",
	Short:   "Manage your intents",
	Example: `mindloop intent start "Get this work done"`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		intentService = intent.NewService(*iRepo)
	},
}

// start intent subcommand
var intentStartCmd = &cobra.Command{
	Use:     "start",
	Short:   "Start a new intent",
	Example: `mindloop intent start "Get this work done"`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// start the intent
		intent, err := intentService.StartIntent(args[0])
		if err != nil {
			utils.PrintErrorln("Error starting intent:", err)
			ac.Logger.Error("Error starting intent", err)
			utils.PrintInfoln("Please try again or check your database connection.")
			return
		}
		utils.PrintSuccessf("Intent '%s' started successfully with id %d!\n", intent.Name, intent.ID)
		ac.Logger.Info(fmt.Sprintf("Intent '%s' started successfully with id %d!", intent.Name, intent.ID))
	},
}

// list intent subcommand
var intentListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List all intents",
	Example: `mindloop intent list`,
	Run: func(cmd *cobra.Command, args []string) {
		intents, err := intentService.ListIntents()
		if err != nil {
			utils.PrintErrorln("Error fetching intents:", err)
			ac.Logger.Error("Error fetching intents", err)
			utils.PrintInfoln("Please check your database connection or try again later.")
			return
		}
		if len(intents) == 0 {
			utils.PrintInfoln("No intents found... Try starting one with 'mindloop intent start <name>'")
			ac.Logger.Info("No intents found. Prompting user to start a new intent.")
			return
		}

		views := []models.IntentView{}
		for _, i := range intents {
			views = append(views, models.ToIntentView(i))
		}
		utils.PrintTable(views)
		ac.Logger.Info(fmt.Sprintf("Listed %d intents successfully.", len(intents)))
	},
}

// current intent subcommand
var intentCurrentCmd = &cobra.Command{
	Use:     "current",
	Short:   "Show current active intents",
	Example: `mindloop intent current`,
	Run: func(cmd *cobra.Command, args []string) {
		intents, err := intentService.ListActiveIntents()
		if err != nil {
			utils.PrintErrorln("Error fetching active intents:", err)
			ac.Logger.Error("Error fetching active intents", err)
			utils.PrintInfoln("Please check your database connection or try again later.")
			return
		}
		if len(intents) == 0 {
			utils.PrintInfoln("No active intents found. To get all intents, use 'mindloop intent list'")
			ac.Logger.Info("No active intents found. Prompting user to list all intents.")
			return
		}

		views := []models.IntentView{}
		for _, i := range intents {
			views = append(views, models.ToIntentView(i))
		}
		utils.PrintTable(views)
		ac.Logger.Info(fmt.Sprintf("Listed %d active intents successfully.", len(intents)))
	},
}

// end intent subcommand
var intentEndCmd = &cobra.Command{
	Use:     "end",
	Short:   "End intent",
	Example: `mindloop intent end 10`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			utils.PrintWarnln("Please provide the intent ID to end.")
			ac.Logger.Warn("No intent ID provided for ending intent.")
			return
		}

		uc := cfg.GetUserConfig()

		intent, milestoneReached, err := intentService.EndIntent(args[0], uc.PointsConfig.Intent)
		if err != nil {
			utils.PrintErrorln("Error ending intent:", err)
			ac.Logger.Error(fmt.Sprintf("Error ending intent with ID %s", args[0]), err)
			return
		}

		utils.PrintSuccessf("Intent '%s' ended successfully! (+%d pts) 🎉\n", intent.Name, uc.PointsConfig.Intent)
		if milestoneReached {
			utils.PrintRocketln("🏆 MILESTONE REACHED! You're on fire! 🏆")
		}
		ac.Logger.Info(fmt.Sprintf("Intent '%s' ended successfully!", intent.Name))
		intentView := models.ToIntentView(*intent)
		utils.PrintTable([]models.IntentView{intentView})
	},
}

func init() {
	rootCmd.AddCommand(intentCmd)
	intentCmd.AddCommand(intentStartCmd)
	intentCmd.AddCommand(intentListCmd)
	intentCmd.AddCommand(intentCurrentCmd)
	intentCmd.AddCommand(intentEndCmd)
}
