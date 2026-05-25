package cli

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	cfg "github.com/snehmatic/mindloop/internal/config"
	"github.com/snehmatic/mindloop/internal/core/habit"
	"github.com/snehmatic/mindloop/internal/log"
	h "github.com/snehmatic/mindloop/internal/repository/habit"
	"github.com/snehmatic/mindloop/internal/utils"
	"github.com/snehmatic/mindloop/models"
	"github.com/spf13/cobra"
)

var (
	daily        *bool
	weekly       *bool
	interactive  *bool
	habitService *habit.Service
)

// parent habit command
var habitCmd = &cobra.Command{
	Use:     "habit",
	Short:   "Manage your habits",
	Example: `mindloop habit add "Exercise"`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		habitService = habit.NewService(h.NewSQLRepository(gdb, logger))
	},
}

// add habit subcommand
var habitAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new habit",
	Example: `mindloop habit add "Exercise" "Need to be fit!" 1 --daily
	mindloop habit add -i`,
	Run: func(cmd *cobra.Command, args []string) {
		utils.PrintRocketln("Great initiative! Adding a new habit...")
		newHabit := &models.Habit{}
		newHabit.SetDefaults()

		if *interactive {
			utils.PrintInfoln("Interactive mode enabled for adding habit...")
			BuildHabitFromInteractiveMode(newHabit)
		} else {
			// non interactive mode
			if len(args) < 3 {
				utils.PrintWarnln("Please provide habit details. Ex. 'mindloop habit add <title> <description> <target_count>' --weekly or --daily(default)")
				ac.Logger.Error("Failed to add habit: missing arguments", nil, log.Field{Key: "habit", Value: newHabit})
				return
			}
			newHabit.Title = args[0]
			newHabit.Description = args[1]
			targetCount, err := strconv.Atoi(args[2])
			if err != nil {
				ac.Logger.Error("Failed to convert target count to integer", err, log.Field{Key: "habit", Value: newHabit})
				utils.PrintErrorln("Invalid target count. Please provide a valid integer.")
				return
			}
			newHabit.TargetCount = targetCount
			newHabit.Interval = GetIntervalFromFlag()
		}

		// Service call replaces direct validation and creation
		err := habitService.CreateHabit(newHabit)
		if err != nil {
			ac.Logger.Error("Failed to add habit", err, log.Field{Key: "habit", Value: newHabit})
			utils.PrintErrorln("Failed to add habit:", err)
			return
		}

		ac.Logger.Info("Habit added successfully", log.Field{Key: "habit", Value: newHabit})

		utils.PrintSuccessf("Habit '%s' added successfully with ID: %d\n", newHabit.Title, newHabit.ID)
	},
}

// delete habit subcommand
var habitDeleteCmd = &cobra.Command{
	Use:     "delete",
	Short:   "Delete a habit",
	Aliases: []string{"rm", "remove", "del"},
	Args:    cobra.ExactArgs(1),
	Example: `mindloop habit delete "Exercise"`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			ac.Logger.Error("No habit ID provided for deletion", nil)
			utils.PrintWarnln("Please provide the habit ID to delete.")
			return
		}
		habitID := args[0]

		// Optional: Fetch to print title before deleting, or just delete.
		// The service delete might error if not found.
		// Let's fetch first to keep UI consistent with previous version (showing title)

		habit, err := habitService.GetHabit(habitID)
		if err != nil {
			ac.Logger.Error("Habit not found", nil)
			utils.PrintErrorln("Habit not found:", err)
			return
		}

		err = habitService.DeleteHabit(habitID)
		if err != nil {
			ac.Logger.Error("Failed to delete habit", err)
			utils.PrintErrorln("Failed to delete habit:", err)
			return
		}

		ac.Logger.Info("Habit deleted successfully", log.Field{Key: "habit", Value: habit})
		utils.PrintSuccessf("Habit '%s' deleted successfully.\n", habit.Title)
	},
}

// update habit subcommand
var habitUpdateCmd = &cobra.Command{
	Use:     "update",
	Short:   "Update a habit",
	Aliases: []string{"edit", "modify"},
	Example: `mindloop habit update "Exercise" --time "1:00 PM"`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			ac.Logger.Error("No habit ID provided for update", nil)
			utils.PrintWarnln("Please provide the habit ID to update.")
			return
		}
		habitId := args[0]

		habit, err := habitService.GetHabit(habitId)
		if err != nil {
			ac.Logger.Error("Habit not found", nil)
			utils.PrintErrorln("Habit not found:", err)
			return
		}

		utils.PrintInfof("Updating habit '%s'...\n", habit.Title)
		utils.PrintTable([]models.HabitView{models.ToHabitView(*habit)})
		utils.PrintInfoln("Entering interactive mode to update Habit (Press Enter to keep current field intact)")
		ac.Logger.Info("Entering interactive mode to update habit", log.Field{Key: "habit", Value: habit})

		// Modifies habit in place
		BuildHabitFromInteractiveMode(habit)

		err = habitService.UpdateHabit(habit)
		if err != nil {
			ac.Logger.Error("Failed to update habit", err)
			utils.PrintErrorln("Failed to update habit:", err)
			return
		}

		ac.Logger.Info("Habit updated successfully", log.Field{Key: "habit", Value: habit})
		utils.PrintSuccessf("Habit '%s' updated successfully.\n", habit.Title)
	},
}

var habitListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List all habits",
	Example: `mindloop habit list`,
	Aliases: []string{"l"},
	Run: func(cmd *cobra.Command, args []string) {
		utils.PrintInfoln("Keep calm, fetching habits...")
		ac.Logger.Info("Fetching habits...")

		intervalFilter := models.IntervalType("")
		if !*daily && !*weekly { // nothing selected via flags
			utils.PrintInfoln("No interval filter applied. Showing all habit logs.")
			ac.Logger.Info("No interval filter applied. Showing all habit logs.")
		} else {
			intervalFilter = GetIntervalFromFlag()
		}

		habits, err := habitService.ListHabits(intervalFilter)
		if err != nil {
			ac.Logger.Error("Failed to retrieve habits", err)
			utils.PrintErrorln("Failed to retrieve habits:", err)
			return
		}

		var habitViews []models.HabitView
		for _, habit := range habits {
			habitViews = append(habitViews, models.ToHabitView(habit))
		}
		utils.PrintTable(habitViews)
	},
}

// log habit as done subcommand
var habitLogCmd = &cobra.Command{
	Use:     "log",
	Aliases: []string{"done", "complete", "mkd"},
	Short:   "Log a habit as done",
	Example: `mindloop habit log "Exercise"`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			ac.Logger.Error("No habit ID provided for logging", nil)
			utils.PrintWarnln("Please provide the habit ID to log.")
			return
		}
		habitID := args[0]

		uc := cfg.GetUserConfig()

		habit, logEntry, milestoneReached, err := habitService.LogHabit(habitID, uc.PointsConfig.Habit)
		if err != nil {
			if err.Error() == "habit already completed for interval" {
				utils.PrintRocketf("Habit already completed. No need to log again.\n")
				return
			}
			ac.Logger.Error("Failed to log habit", err)
			utils.PrintErrorln("Failed to log habit:", err)
			return
		}

		msg := fmt.Sprintf("Habit %s logged %d/%d times in %s interval", habit.Title, logEntry.ActualCount, habit.TargetCount, habit.Interval)
		ac.Logger.Info(msg, log.Field{Key: "habit", Value: habit})
		utils.PrintLoadingf("Habit %s logged %d/%d times in %s interval.\n", habit.Title, logEntry.ActualCount, habit.TargetCount, habit.Interval)
		utils.PrintInfof("Use 'mindloop habit unlog <id>' to mark it as undone, and reset to 0/%d.\n", habit.TargetCount)

		if logEntry.ActualCount == habit.TargetCount {
			utils.PrintSuccessf("Habit '%s' marked done! (+%d pts) 🎉\n", habit.Title, uc.PointsConfig.Habit)
			if milestoneReached {
				utils.PrintRocketln("🏆 MILESTONE REACHED! You're on fire! 🏆")
			}
		} else {
			utils.PrintSuccessf("Habit '%s' logged successfully.\n", habit.Title)
		}
	},
}

// log habit as done subcommand
var habitUnLogCmd = &cobra.Command{
	Use:     "unlog",
	Aliases: []string{"undone", "incomplete", "mkud"},
	Args:    cobra.ExactArgs(1),
	Short:   "Log a habit as undone",
	Example: `mindloop habit unlog "Exercise"`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			ac.Logger.Error("No habit ID provided for unlogging", nil)
			utils.PrintWarnln("Please provide the habit ID to unlog.")
			return
		}
		habitID := args[0]

		habit, err := habitService.UnlogHabit(habitID)
		if err != nil {
			ac.Logger.Error("Failed to unlog habit", err)
			utils.PrintErrorln("Failed to unlog habit:", err)
			return
		}

		ac.Logger.Info("Habit unlogged successfully", log.Field{Key: "habit", Value: habit})
		utils.PrintSuccessf("Habit '%s' unlogged successfully. Reset to 0/%d.\n", habit.Title, habit.TargetCount)
		utils.PrintInfoln("Use 'mindloop habit log <id>' to mark it as done again.")
	},
}

// habit show subcommand
var habitLogShowCmd = &cobra.Command{
	Use:     "show",
	Aliases: []string{"status", "check", "stats"},
	Short:   "Check habit logs show -w",
	Run: func(cmd *cobra.Command, args []string) {
		utils.PrintRocketln("'show me the logs'? Here you go Chief...")

		intervalFilter := models.IntervalType("")
		if !*daily && !*weekly { // nothing selected via flags
			utils.PrintInfoln("No interval filter applied. Showing all habit logs.")
			ac.Logger.Info("No interval filter applied. Showing all habit logs.")
			intervalFilter = "" // no filter
		} else {
			intervalFilter = GetIntervalFromFlag()
		}

		habitLogs, err := habitService.ListHabitLogs(intervalFilter)
		if err != nil {
			ac.Logger.Error("Failed to retrieve habit logs", err)
			utils.PrintErrorln("Failed to retrieve habit logs:", err)
			return
		}

		if len(habitLogs) == 0 {
			utils.PrintInfoln("Ruh-roh! No habit logs found. Start logging habits with 'mindloop habit log <id>'")
			return
		}

		habitLogViews := models.ToHabitLogViews(habitLogs)
		utils.PrintTable(habitLogViews)
	},
}

func init() {
	// cmds
	rootCmd.AddCommand(habitCmd)
	habitCmd.AddCommand(habitAddCmd)
	habitCmd.AddCommand(habitDeleteCmd)
	habitCmd.AddCommand(habitUpdateCmd)
	habitCmd.AddCommand(habitLogCmd)
	habitCmd.AddCommand(habitUnLogCmd)
	habitCmd.AddCommand(habitListCmd)
	habitLogCmd.AddCommand(habitLogShowCmd)

	// flags
	// all = habitCmd.PersistentFlags().BoolP("all", "A", false, "Select all habits") // not using now
	daily = habitCmd.PersistentFlags().BoolP("daily", "d", false, "Set habit as daily")
	weekly = habitCmd.PersistentFlags().BoolP("weekly", "w", false, "Set habit as weekly")
	interactive = habitCmd.PersistentFlags().BoolP("interactive", "i", false, "Interactive mode for adding habit")
}

// GetIntervalFromFlag returns the interval type based on the flags set
// Defaults to daily if no flags are set
func GetIntervalFromFlag() models.IntervalType {
	if *daily {
		return models.Daily
	} else if *weekly {
		return models.Weekly
	}
	utils.PrintInfoln("Defaulting to daily interval. Use -w or -d to set weekly or daily respectively.")
	return models.Daily
}

// BuildHabitFromInteractiveMode builds a Habit from user input in interactive mode
// If a nil pointer is passed, it initializes a new Habit
// Returns the updated Habit pointer
func BuildHabitFromInteractiveMode(hb *models.Habit) *models.Habit {
	if hb == nil {
		hb = &models.Habit{}
		hb.SetDefaults()
	}

	fmt.Print("Enter habit name: ")
	inputReader := bufio.NewReader(os.Stdin)
	input, _ := inputReader.ReadString('\n')
	title := input[:len(input)-1]
	if title != "" {
		hb.Title = title
	}

	fmt.Print("Enter habit description: ")
	inputReader = bufio.NewReader(os.Stdin)
	input, _ = inputReader.ReadString('\n')
	desc := input[:len(input)-1]
	if desc != "" {
		hb.Description = desc
	}

	fmt.Print("Enter target count (default 1): ")
	var targetCount int
	_, _ = fmt.Scanln(&targetCount)
	if targetCount > 0 {
		hb.TargetCount = targetCount
	}

	for {
		fmt.Print("Select interval (daily/weekly, default daily): ")
		var interval string
		_, _ = fmt.Scanln(&interval)
		if interval != "" {
			if !models.IsValidIntervalType(interval) {
				ac.Logger.Error("Invalid interval type.", nil, log.Field{Key: "habit", Value: hb})
				utils.PrintWarnln("Invalid interval type. Retry with 'daily' or 'weekly'.")

				continue
			}
			hb.Interval = models.IntervalType(interval)

			break
		}

		break
	}

	return hb
}
