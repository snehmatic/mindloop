package cli

import (
	"fmt"
	"strconv"

	"github.com/snehmatic/mindloop/db"
	"github.com/snehmatic/mindloop/internal/config"
	"github.com/snehmatic/mindloop/internal/core/routine"
	"github.com/snehmatic/mindloop/internal/utils"
	"github.com/spf13/cobra"
)

var routineCmd = &cobra.Command{
	Use:   "routine",
	Short: "Manage daily and weekly routines",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		config.ValidateUserConfig(cmd)
	},
}

var routineCreateCmd = &cobra.Command{
	Use:   "create [title]",
	Short: "Create a new routine",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		title := args[0]
		timeOfDay, _ := cmd.Flags().GetString("time")

		appConfig := config.GetConfig()
		database, err := db.ConnectToDb(*appConfig)
		if err != nil {
			utils.PrintErrorln("Failed to connect to database")
			return
		}

		svc := routine.NewService(database)
		r, err := svc.CreateRoutine(title, timeOfDay)
		if err != nil {
			utils.PrintErrorln(fmt.Sprintf("Failed to create routine: %v", err))
			return
		}

		utils.PrintSuccessln(fmt.Sprintf("Created routine '%s' (ID: %d)", r.Title, r.ID))
	},
}

var routineAddHabitCmd = &cobra.Command{
	Use:   "add-habit [routine_id] [habit_id]",
	Short: "Add a habit to a routine",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		routineID, err := strconv.ParseUint(args[0], 10, 32)
		if err != nil {
			utils.PrintErrorln("Invalid routine ID")
			return
		}
		habitID, err := strconv.ParseUint(args[1], 10, 32)
		if err != nil {
			utils.PrintErrorln("Invalid habit ID")
			return
		}

		appConfig := config.GetConfig()
		database, err := db.ConnectToDb(*appConfig)
		if err != nil {
			utils.PrintErrorln("Failed to connect to database")
			return
		}

		svc := routine.NewService(database)
		err = svc.AddHabitToRoutine(uint(routineID), uint(habitID))
		if err != nil {
			utils.PrintErrorln(fmt.Sprintf("Failed to add habit to routine: %v", err))
			return
		}

		utils.PrintSuccessln("Added habit to routine successfully")
	},
}

var routineListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all routines",
	Run: func(cmd *cobra.Command, args []string) {
		appConfig := config.GetConfig()
		database, err := db.ConnectToDb(*appConfig)
		if err != nil {
			utils.PrintErrorln("Failed to connect to database")
			return
		}

		svc := routine.NewService(database)
		routines, err := svc.ListRoutines()
		if err != nil {
			utils.PrintErrorln("Failed to list routines")
			return
		}

		for _, r := range routines {
			fmt.Printf("[%d] %s (%s) - %d habits\n", r.ID, r.Title, r.TimeOfDay, len(r.Habits))
		}
	},
}

func init() {
	rootCmd.AddCommand(routineCmd)
	routineCmd.AddCommand(routineCreateCmd)
	routineCmd.AddCommand(routineAddHabitCmd)
	routineCmd.AddCommand(routineListCmd)

	routineCreateCmd.Flags().String("time", "", "Time of day for the routine (e.g. Morning, Evening)")
}
