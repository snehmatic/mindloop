package cli

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"time"

	. "github.com/snehmatic/mindloop/internal/utils"
	"github.com/snehmatic/mindloop/models"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

var (
	all         *bool
	daily       *bool
	weekly      *bool
	interactive *bool
)

// parent habit command
var habitCmd = &cobra.Command{
	Use:     "habit",
	Short:   "Manage your habits",
	Example: `mindloop habit add "Excercise"`,
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

		// validate habit
		err := newHabit.ValidateHabit()
		if err != nil {
			ac.Logger.Error().
				Interface("habit", newHabit).
				Err(err).
				Msgf("Habit validation failed: %v", err)
			PrintErrorln("Habit validation failed: ", err)
			return
		}

		// persist new habit to db
		res := gdb.Create(newHabit)
		if res.Error != nil {
			ac.Logger.Error().
				Interface("habit", newHabit).
				Err(res.Error).
				Msg("Failed to add habit in db")
			PrintErrorln("Failed to add habit:", res.Error)
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
		var habit models.Habit
		res := gdb.Where("ID = ?", habitID).First(&habit)
		if res.Error != nil {
			ac.Logger.Error().
				Interface("habit", habit).
				Msg("Habit not found")
			PrintErrorln("Habit not found:", res.Error)
			return
		}

		res = gdb.Delete(&habit)
		if res.Error != nil {
			ac.Logger.Error().
				Interface("habit", habit).
				Err(res.Error).
				Msg("Failed to delete habit")
			PrintErrorln("Failed to delete habit:", res.Error)
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
		var habit models.Habit
		res := gdb.Where("ID = ?", habitId).First(&habit)
		if res.Error != nil {
			ac.Logger.Error().
				Interface("habit", habit).
				Msg("Habit not found")
			PrintErrorln("Habit not found:", res.Error)
			return
		}

		PrintInfof("Updating habit '%s'...\n", habit.Title)
		PrintTable([]models.HabitView{models.ToHabitView(habit)})
		PrintInfoln("Entering interactive mode to update Habit (Press Enter to keep current field intact)")
		ac.Logger.Info().
			Interface("habit", habit).
			Msg("Entering interactive mode to update habit")
		BuildHabitFromInteractiveMode(&habit)

		err := habit.ValidateHabit()
		if err != nil {
			ac.Logger.Error().
				Interface("habit", habit).
				Err(err).
				Msg("Habit validation failed")
			PrintErrorln("Habit validation failed: ", err)
			return
		}

		res = gdb.Save(&habit)
		if res.Error != nil {
			ac.Logger.Error().
				Interface("habit", habit).
				Err(res.Error).
				Msg("Failed to update habit")
			PrintErrorln("Failed to update habit:", res.Error)
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

		intervalFilter := ""
		if !*daily && !*weekly { // nothing selected via flags
			PrintInfoln("No interval filter applied. Showing all habit logs.")
			ac.Logger.Info().Msg("No interval filter applied. Showing all habit logs.")
			intervalFilter = "" // no filter
		} else {
			intervalFilter = fmt.Sprintf("Interval = '%s'", GetIntervalFromFlag())
		}

		var habits []models.Habit
		res := gdb.Where(intervalFilter).Find(&habits)
		if res.Error != nil {
			ac.Logger.Error().
				Err(res.Error).
				Msg("Failed to retrieve habits")
			PrintErrorln("Failed to retrieve habits:", res.Error)
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
		var habit models.Habit
		res := gdb.Where("ID = ?", habitID).First(&habit)
		if res.Error != nil {
			ac.Logger.Error().
				Interface("habit", habit).
				Msg("Habit not found")
			PrintErrorln("Habit not found:", res.Error)
			return
		}

		// check if habit is already logged today
		var existingLog models.HabitLog
		today := time.Now().Truncate(24 * time.Hour)
		endedAt := today

		switch habit.Interval {
		case models.Daily:
			// if the log exists for today
			res = gdb.Where("HabitID = ? AND EndedAt = ?", habit.ID, today, today).First(&existingLog)
		case models.Weekly:
			// if the log exists for the current week
			startOfWeek := time.Now().AddDate(0, 0, -int(time.Now().Weekday()))
			endOfWeek := startOfWeek.AddDate(0, 0, 6).Truncate(24 * time.Hour)
			endedAt = endOfWeek
			res = gdb.Where("HabitID = ? AND CreatedAt >= ? AND EndedAt <= ?", habit.ID, startOfWeek, endOfWeek).First(&existingLog)
		}
		if res.Error == nil {
			// habit log record found, update it if required
			if existingLog.ActualCount >= habit.TargetCount {
				PrintRocketf("Habit already completed for %s interval. No need to log again.\n", habit.Interval)
				ac.Logger.Info().Msgf("Habit already completed for %s interval. No need to log again.", habit.Interval)
				return
			}

			// habit not complete, increment the count
			existingLog.ActualCount++
			existingLog.EndedAt = endedAt
			res = gdb.Save(&existingLog)
			if res.Error != nil {
				ac.Logger.Error().
					Interface("habitLog", existingLog).
					Err(res.Error).
					Msg("Failed to update existing habit log")
				PrintErrorln("Failed to update existing habit log:", res.Error)
				return
			}

			ac.Logger.Info().
				Interface("habit", habit).
				Msgf("Habit %s logged %d/%d times in %s interval", habit.Title, existingLog.ActualCount, habit.TargetCount, habit.Interval)
			PrintLoadingf("Habit %s logged %d/%d times in %s interval.\n", habit.Title, existingLog.ActualCount, habit.TargetCount, habit.Interval)
			PrintInfof("Use 'mindloop habit unlog <id>' to mark it as undone, and reset to 0/%d.\n", habit.TargetCount)

			return
		} else if res.Error != gorm.ErrRecordNotFound {
			ac.Logger.Error().
				Interface("habit", habit).
				Err(res.Error).
				Msg("Failed to check existing habit log")
			PrintErrorln("Failed to check existing habit log:", res.Error)
			return
		}

		// habit log record not found in db
		// create a new one
		habitLog := &models.HabitLog{
			HabitID:     habit.ID,
			Title:       habit.Title,
			Interval:    habit.Interval,
			TargetCount: habit.TargetCount,
			ActualCount: 1,
			EndedAt:     endedAt,
		}
		res = gdb.Create(habitLog)
		if res.Error != nil {
			ac.Logger.Error().
				Interface("habitLog", habitLog).
				Err(res.Error).
				Msg("Failed to log habit")
			PrintErrorln("Failed to log habit:", res.Error)
			return
		}

		ac.Logger.Info().
			Interface("habit", habit).
			Msg("Habit logged successfully")
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
		var habit models.Habit
		res := gdb.Where("ID = ?", habitID).First(&habit)
		if res.Error != nil {
			ac.Logger.Error().
				Interface("habit", habit).
				Msg("Habit not found")
			PrintErrorln("Habit not found:", res.Error)
			return
		}

		var existingLog models.HabitLog
		today := time.Now().Format("2006-01-02")
		switch habit.Interval {
		case models.Daily:
			res = gdb.Where("HabitID = ? AND DATE(CreatedAt) = ?", habit.ID, today).First(&existingLog)
		case models.Weekly:
			startOfWeek := time.Now().AddDate(0, 0, -int(time.Now().Weekday()))
			res = gdb.Where("HabitID = ? AND CreatedAt >= ?", habit.ID, startOfWeek).First(&existingLog)
		}
		if res.Error != nil {
			if res.Error == gorm.ErrRecordNotFound {
				PrintWarnln("No existing log found for this habit.")
				ac.Logger.Warn().Msg("No existing log found for this habit.")
				return
			}
			ac.Logger.Error().
				Err(res.Error).
				Msg("Failed to check existing habit log")
			PrintErrorln("Failed to check existing habit log:", res.Error)
			return
		}

		if existingLog.ActualCount <= 0 {
			PrintWarnln("Habit is already marked as undone.")
			ac.Logger.Warn().Msg("Habit is already marked as undone.")
			return
		}

		existingLog.ActualCount = 0
		res = gdb.Save(&existingLog)
		if res.Error != nil {
			ac.Logger.Error().
				Err(res.Error).
				Msg("Failed to unlog habit")
			PrintErrorln("Failed to unlog habit:", res.Error)
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

		intervalFilter := ""
		if !*daily && !*weekly { // nothing selected via flags
			PrintInfoln("No interval filter applied. Showing all habit logs.")
			ac.Logger.Info().Msg("No interval filter applied. Showing all habit logs.")
			intervalFilter = "" // no filter
		} else {
			intervalFilter = fmt.Sprintf("Interval = '%s'", GetIntervalFromFlag())
		}

		var habitLogs []models.HabitLog
		res := gdb.Where(intervalFilter).Order("CreatedAt DESC").Find(&habitLogs)
		if res.Error != nil {
			ac.Logger.Error().
				Err(res.Error).
				Msg("Failed to retrieve habit logs")
			PrintErrorln("Failed to retrieve habit logs:", res.Error)
			return
		}

		if len(habitLogs) == 0 {
			PrintInfoln("Ruh-roh! No habit logs found. Start logging habits with 'mindloop habit log <id>'")
			ac.Logger.Info().Msg("No habit logs found. Prompting user to log habits.")
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
