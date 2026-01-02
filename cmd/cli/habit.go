package cli

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/snehmatic/mindloop/internal/core/habit"
	. "github.com/snehmatic/mindloop/internal/utils"
	"github.com/snehmatic/mindloop/models"
	"github.com/spf13/cobra"
)

var (
	all          *bool
	daily        *bool
	weekly       *bool
	interactive  *bool
	habitService *habit.Service
)

// parent habit command
var habitCmd = &cobra.Command{
	Use:     "habit",
	Short:   "Manage your habits",
	Example: `mindloop habit add "Excercise"`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		habitService = habit.NewService(gdb)
	},
}

// add habit subcommand
var habitAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new habit",
	Example: `mindloop habit add "Excercise" "Need to be fit!" 1 --daily
	mindloop habit add -i`,
	Run: func(cmd *cobra.Command, args []string) {
		PrintRocketln("Great initiative! Adding a new habit...")
		newHabit := &models.Habit{}
		newHabit.SetDefaults()

		if *interactive {
			PrintInfoln("Interactive mode enabled for adding habit...")
			BuildHabitFromInteractiveMode(newHabit)
		} else {
			// non interactive mode
			if len(args) < 3 {
				PrintWarnln("Please provide habit details. Ex. 'mindloop habit add <title> <description> <target_count>' --weekly or --daily(default)")
				ac.Logger.Error().
					Interface("habit", newHabit).
					Msg("Failed to add habit: missing arguments")
				return
			}
			newHabit.Title = args[0]
			newHabit.Description = args[1]
			targetCount, err := strconv.Atoi(args[2])
			if err != nil {
				ac.Logger.Error().
					Interface("habit", newHabit).
					Err(err).
					Msg("Failed to convert target count to integer")
				PrintErrorln("Invalid target count. Please provide a valid integer.")
				return
			}
			newHabit.TargetCount = targetCount
			newHabit.Interval = GetIntervalFromFlag()
		}

		// Service call replaces direct validation and creation
		err := habitService.CreateHabit(newHabit)
		if err != nil {
			ac.Logger.Error().
				Interface("habit", newHabit).
				Err(err).
				Msg("Failed to add habit")
			PrintErrorln("Failed to add habit:", err)
			return
		}

		ac.Logger.Info().
			Interface("habit", newHabit).
			Msg("Habit added successfully")

		PrintSuccessf("Habit '%s' added successfully with ID: %d\n", newHabit.Title, newHabit.ID)
	},
}

// delete habit subcommand
var habitDeleteCmd = &cobra.Command{
	Use:     "delete",
	Short:   "Delete a habit",
	Aliases: []string{"rm", "remove", "del"},
	Args:    cobra.ExactArgs(1),
	Example: `mindloop habit delete "Excercise"`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			ac.Logger.Error().Msg("No habit ID provided for deletion")
			PrintWarnln("Please provide the habit ID to delete.")
			return
		}
		habitID := args[0]

		// Optional: Fetch to print title before deleting, or just delete.
		// The service delete might error if not found.
		// Let's fetch first to keep UI consistent with previous version (showing title)

		habit, err := habitService.GetHabit(habitID)
		if err != nil {
			ac.Logger.Error().Msg("Habit not found")
			PrintErrorln("Habit not found:", err)
			return
		}

		err = habitService.DeleteHabit(habitID)
		if err != nil {
			ac.Logger.Error().Err(err).Msg("Failed to delete habit")
			PrintErrorln("Failed to delete habit:", err)
			return
		}

		ac.Logger.Info().
			Interface("habit", habit).
			Msg("Habit deleted successfully")
		PrintSuccessf("Habit '%s' deleted successfully.\n", habit.Title)
	},
}

// update habit subcommand
var habitUpdateCmd = &cobra.Command{
	Use:     "update",
	Short:   "Update a habit",
	Aliases: []string{"edit", "modify"},
	Example: `mindloop habit update "Excercise" --time "1:00 PM"`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			ac.Logger.Error().Msg("No habit ID provided for update")
			PrintWarnln("Please provide the habit ID to update.")
			return
		}
		habitId := args[0]

		habit, err := habitService.GetHabit(habitId)
		if err != nil {
			ac.Logger.Error().Msg("Habit not found")
			PrintErrorln("Habit not found:", err)
			return
		}

		PrintInfof("Updating habit '%s'...\n", habit.Title)
		PrintTable([]models.HabitView{models.ToHabitView(*habit)})
		PrintInfoln("Entering interactive mode to update Habit (Press Enter to keep current field intact)")
		ac.Logger.Info().
			Interface("habit", habit).
			Msg("Entering interactive mode to update habit")

		// Modifies habit in place
		BuildHabitFromInteractiveMode(habit)

		err = habitService.UpdateHabit(habit)
		if err != nil {
			ac.Logger.Error().Err(err).Msg("Failed to update habit")
			PrintErrorln("Failed to update habit:", err)
			return
		}

		ac.Logger.Info().
			Interface("habit", habit).
			Msg("Habit updated successfully")
		PrintSuccessf("Habit '%s' updated successfully.\n", habit.Title)
	},
}

var habitListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List all habits",
	Example: `mindloop habit list`,
	Aliases: []string{"l"},
	Run: func(cmd *cobra.Command, args []string) {
		PrintInfoln("Keep calm, fetching habits...")
		ac.Logger.Info().Msg("Fetching habits...")

		intervalFilter := models.IntervalType("")
		if !*daily && !*weekly { // nothing selected via flags
			PrintInfoln("No interval filter applied. Showing all habit logs.")
			ac.Logger.Info().Msg("No interval filter applied. Showing all habit logs.")
		} else {
			intervalFilter = GetIntervalFromFlag()
		}

		habits, err := habitService.ListHabits(intervalFilter)
		if err != nil {
			ac.Logger.Error().Err(err).Msg("Failed to retrieve habits")
			PrintErrorln("Failed to retrieve habits:", err)
			return
		}

		var habitViews []models.HabitView
		for _, habit := range habits {
			habitViews = append(habitViews, models.ToHabitView(habit))
		}
		PrintTable(habitViews)
	},
}

// log habit as done subcommand
var habitLogCmd = &cobra.Command{
	Use:     "log",
	Aliases: []string{"done", "complete", "mkd"},
	Short:   "Log a habit as done",
	Example: `mindloop habit log "Excercise"`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			ac.Logger.Error().Msg("No habit ID provided for logging")
			PrintWarnln("Please provide the habit ID to log.")
			return
		}
		habitID := args[0]

		habit, log, err := habitService.LogHabit(habitID)
		if err != nil {
			if err.Error() == "habit already completed for interval" {
				PrintRocketf("Habit already completed. No need to log again.\n")
				return
			}
			ac.Logger.Error().Err(err).Msg("Failed to log habit")
			PrintErrorln("Failed to log habit:", err)
			return
		}

		ac.Logger.Info().
			Interface("habit", habit).
			Msgf("Habit %s logged %d/%d times in %s interval", habit.Title, log.ActualCount, habit.TargetCount, habit.Interval)
		PrintLoadingf("Habit %s logged %d/%d times in %s interval.\n", habit.Title, log.ActualCount, habit.TargetCount, habit.Interval)
		PrintInfof("Use 'mindloop habit unlog <id>' to mark it as undone, and reset to 0/%d.\n", habit.TargetCount)
		PrintSuccessf("Habit '%s' logged successfully.\n", habit.Title)
	},
}

// log habit as done subcommand
var habitUnLogCmd = &cobra.Command{
	Use:     "unlog",
	Aliases: []string{"undone", "incomplete", "mkud"},
	Args:    cobra.ExactArgs(1),
	Short:   "Log a habit as undone",
	Example: `mindloop habit unlog "Excercise"`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			ac.Logger.Error().Msg("No habit ID provided for unlogging")
			PrintWarnln("Please provide the habit ID to unlog.")
			return
		}
		habitID := args[0]

		habit, err := habitService.UnlogHabit(habitID)
		if err != nil {
			ac.Logger.Error().Err(err).Msg("Failed to unlog habit")
			PrintErrorln("Failed to unlog habit:", err)
			return
		}

		ac.Logger.Info().
			Interface("habit", habit).
			Msg("Habit unlogged successfully")
		PrintSuccessf("Habit '%s' unlogged successfully. Reset to 0/%d.\n", habit.Title, habit.TargetCount)
		PrintInfoln("Use 'mindloop habit log <id>' to mark it as done again.")
	},
}

// habit show subcommand
var habitLogShowCmd = &cobra.Command{
	Use:     "show",
	Aliases: []string{"status", "check", "stats"},
	Short:   "Check habit logs show -w",
	Run: func(cmd *cobra.Command, args []string) {
		PrintRocketln("'show me the logs'? Here you go Chief...")

		intervalFilter := models.IntervalType("")
		if !*daily && !*weekly { // nothing selected via flags
			PrintInfoln("No interval filter applied. Showing all habit logs.")
			ac.Logger.Info().Msg("No interval filter applied. Showing all habit logs.")
			intervalFilter = "" // no filter
		} else {
			intervalFilter = GetIntervalFromFlag()
		}

		habitLogs, err := habitService.ListHabitLogs(intervalFilter)
		if err != nil {
			ac.Logger.Error().Err(err).Msg("Failed to retrieve habit logs")
			PrintErrorln("Failed to retrieve habit logs:", err)
			return
		}

		if len(habitLogs) == 0 {
			PrintInfoln("Ruh-roh! No habit logs found. Start logging habits with 'mindloop habit log <id>'")
			return
		}

		habitLogViews := models.ToHabitLogViews(habitLogs)
		PrintTable(habitLogViews)
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
	all = habitCmd.PersistentFlags().BoolP("all", "A", false, "Select all habits") // not using now
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
	PrintInfoln("Defaulting to daily interval. Use -w or -d to set weekly or daily respectively.")
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
	fmt.Scanln(&targetCount)
	if targetCount > 0 {
		hb.TargetCount = targetCount
	}

	for {
		fmt.Print("Select interval (daily/weekly, default daily): ")
		var interval string
		fmt.Scanln(&interval)
		if interval != "" {
			if !models.IsValidIntervalType(interval) {
				ac.Logger.Error().
					Interface("habit", hb).
					Msg("Invalid interval type.")
				PrintWarnln("Invalid interval type. Retry with 'daily' or 'weekly'.")
				continue
			}
			hb.Interval = models.IntervalType(interval)
			break
		}
		break
	}

	return hb
}
