package cli

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/snehmatic/mindloop/internal/utils"
	"github.com/snehmatic/mindloop/models"
	"github.com/spf13/cobra"
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
		fmt.Println("Great initiative! Adding a new habit...")
		newHabit := &models.Habit{}
		newHabit.SetDefaults()

		if *interactive {
			fmt.Println("Interactive mode enabled for adding habit...")
			BuildHabitFromInteractiveMode(newHabit)
		} else {
			// non interactive mode
			if len(args) < 3 {
				fmt.Println("Please provide habit details. Ex. 'mindloop habit add <title> <description> <target_count>' --weekly or --daily(default)")
				logger.Error().
					Interface("habit", newHabit).
					Msg("Failed to add habit: missing arguments")
				return
			}
			newHabit.Title = args[0]
			newHabit.Description = args[1]
			targetCount, err := strconv.Atoi(args[2])
			if err != nil {
				logger.Error().
					Interface("habit", newHabit).
					Err(err).
					Msg("Failed to convert target count to integer")
				fmt.Println("Invalid target count. Please provide a valid integer.")
				return
			}
			newHabit.TargetCount = targetCount
			newHabit.Interval = GetIntervalFromFlag()
		}

		// validate habit
		err := newHabit.ValidateHabit()
		if err != nil {
			logger.Error().
				Interface("habit", newHabit).
				Err(err).
				Msgf("Habit validation failed: %v", err)
			fmt.Println("Habit validation failed: ", err)
			return
		}

		// persist new habit to db
		res := gdb.Create(newHabit)
		if res.Error != nil {
			logger.Error().
				Interface("habit", newHabit).
				Err(res.Error).
				Msg("Failed to add habit in db")
			fmt.Println("Failed to add habit:", res.Error)
			return
		}

		logger.Info().
			Interface("habit", newHabit).
			Msg("Habit added successfully")

		fmt.Printf("Habit '%s' added successfully with ID: %d\n", newHabit.Title, newHabit.ID)
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
			logger.Error().Msg("No habit ID provided for deletion")
			fmt.Println("Please provide the habit ID to delete.")
			return
		}
		habitID := args[0]
		var habit models.Habit
		res := gdb.Where("ID = ?", habitID).First(&habit)
		if res.Error != nil {
			logger.Error().
				Interface("habit", habit).
				Msg("Habit not found")
			fmt.Println("Habit not found:", res.Error)
			return
		}

		res = gdb.Delete(&habit)
		if res.Error != nil {
			logger.Error().
				Interface("habit", habit).
				Err(res.Error).
				Msg("Failed to delete habit")
			fmt.Println("Failed to delete habit:", res.Error)
			return
		}

		logger.Info().
			Interface("habit", habit).
			Msg("Habit deleted successfully")
		fmt.Printf("Habit '%s' deleted successfully.\n", habit.Title)
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
			logger.Error().Msg("No habit ID provided for update")
			fmt.Println("Please provide the habit ID to update.")
			return
		}
		habitId := args[0]
		var habit models.Habit
		res := gdb.Where("ID = ?", habitId).First(&habit)
		if res.Error != nil {
			logger.Error().
				Interface("habit", habit).
				Msg("Habit not found")
			fmt.Println("Habit not found:", res.Error)
			return
		}

		fmt.Printf("Updating habit '%s'...\n", habit.Title)
		var habitArr []models.HabitView
		habitArr = append(habitArr, models.ToHabitView(habit))
		utils.PrintTable(habitArr)
		fmt.Println("Entering interactive mode to update Habit (Press Enter to keep current field intact)")
		logger.Info().
			Interface("habit", habit).
			Msg("Entering interactive mode to update habit")
		BuildHabitFromInteractiveMode(&habit)

		err := habit.ValidateHabit()
		if err != nil {
			logger.Error().
				Interface("habit", habit).
				Err(err).
				Msg("Habit validation failed")
			fmt.Println("Habit validation failed: ", err)
			return
		}

		res = gdb.Save(&habit)
		if res.Error != nil {
			logger.Error().
				Interface("habit", habit).
				Err(res.Error).
				Msg("Failed to update habit")
			fmt.Println("Failed to update habit:", res.Error)
			return
		}

		logger.Info().
			Interface("habit", habit).
			Msg("Habit updated successfully")
		fmt.Printf("Habit '%s' updated successfully.\n", habit.Title)
	},
}

var habitListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List all habits",
	Example: `mindloop habit list`,
	Aliases: []string{"l"},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Keep calm, fetching all habits...")
		logger.Info().Msg("Fetching all habits...")

		var habits []models.Habit
		res := gdb.Find(&habits)
		if res.Error != nil {
			logger.Error().
				Err(res.Error).
				Msg("Failed to retrieve habits")
			fmt.Println("Failed to retrieve habits:", res.Error)
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
	Args:    cobra.ExactArgs(1),
	Short:   "Log a habit as done",
	Example: `mindloop habit log "Excercise"`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Logged habit as done...")
	},
}

// habit status
var habitLogStatsCmd = &cobra.Command{
	Use:     "stats",
	Aliases: []string{"status", "check"},
	Short:   "Check habit logs stats -w",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Here is your habit logs stats...")

		interval := GetIntervalFromFlag()

		var habitLogs []models.HabitLog
		res := gdb.Find(&habitLogs).Where("interval = ?", interval)
		if res.Error != nil {
			logger.Error().
				Err(res.Error).
				Msg("Failed to retrieve habit logs")
			fmt.Println("Failed to retrieve habit logs:", res.Error)
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
	habitCmd.AddCommand(habitListCmd)

	// flags
	all = habitCmd.PersistentFlags().BoolP("all", "A", false, "Select all habits") // not using now
	daily = habitAddCmd.Flags().BoolP("daily", "d", false, "Set habit as daily")
	weekly = habitAddCmd.Flags().BoolP("weekly", "w", false, "Set habit as weekly")
	interactive = habitAddCmd.Flags().BoolP("interactive", "i", false, "Interactive mode for adding habit")
}

// GetIntervalFromFlag returns the interval type based on the flags set
// Defaults to daily if no flags are set
func GetIntervalFromFlag() models.IntervalType {
	if *daily {
		return models.Daily
	} else if *weekly {
		return models.Weekly
	}
	fmt.Println("Defaulting to daily interval.")
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

	fmt.Print("Select interval (daily/weekly, default daily): ")
	var interval string
	fmt.Scanln(&interval)
	if interval != "" {
		if !models.IsValidIntervalType(interval) {
			logger.Error().
				Interface("habit", hb).
				Msg("Invalid interval type.")
			fmt.Println("Invalid interval type.")
			os.Exit(1)
		}
		hb.Interval = models.IntervalType(interval)
	}

	return hb
}
