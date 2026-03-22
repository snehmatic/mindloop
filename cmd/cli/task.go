package cli

import (
	"fmt"
	"strconv"

	"github.com/snehmatic/mindloop/db"
	"github.com/snehmatic/mindloop/internal/config"
	"github.com/snehmatic/mindloop/internal/core/task"
	"github.com/snehmatic/mindloop/internal/utils"
	"github.com/spf13/cobra"
)

var taskCmd = &cobra.Command{
	Use:   "task",
	Short: "Manage tasks and sub-tasks",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		config.ValidateUserConfig(cmd)
	},
}

var taskAddCmd = &cobra.Command{
	Use:   "add [title]",
	Short: "Add a new task",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		title := args[0]
		intentIDStr, _ := cmd.Flags().GetString("intent-id")
		focusIDStr, _ := cmd.Flags().GetString("focus-id")

		var intentID *uint
		if intentIDStr != "" {
			id, err := strconv.ParseUint(intentIDStr, 10, 32)
			if err == nil {
				uid := uint(id)
				intentID = &uid
			}
		}

		var focusID *uint
		if focusIDStr != "" {
			id, err := strconv.ParseUint(focusIDStr, 10, 32)
			if err == nil {
				uid := uint(id)
				focusID = &uid
			}
		}

		appConfig := config.GetConfig()
		database, err := db.ConnectToDb(*appConfig)
		if err != nil {
			utils.PrintErrorln("Failed to connect to database")
			return
		}

		svc := task.NewService(database)
		t, err := svc.CreateTask(title, intentID, focusID)
		if err != nil {
			utils.PrintErrorln(fmt.Sprintf("Failed to create task: %v", err))
			return
		}

		utils.PrintSuccessln(fmt.Sprintf("Created task '%s' (ID: %d)", t.Title, t.ID))
	},
}

var taskCompleteCmd = &cobra.Command{
	Use:   "complete [id]",
	Short: "Complete a task",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.ParseUint(args[0], 10, 32)
		if err != nil {
			utils.PrintErrorln("Invalid task ID")
			return
		}

		appConfig := config.GetConfig()
		database, err := db.ConnectToDb(*appConfig)
		if err != nil {
			utils.PrintErrorln("Failed to connect to database")
			return
		}

		svc := task.NewService(database)
		err = svc.CompleteTask(uint(id))
		if err != nil {
			utils.PrintErrorln(fmt.Sprintf("Failed to complete task: %v", err))
			return
		}

		utils.PrintSuccessln(fmt.Sprintf("Task %d marked as completed", id))
	},
}

var taskListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tasks",
	Run: func(cmd *cobra.Command, args []string) {
		appConfig := config.GetConfig()
		database, err := db.ConnectToDb(*appConfig)
		if err != nil {
			utils.PrintErrorln("Failed to connect to database")
			return
		}

		svc := task.NewService(database)
		tasks, err := svc.ListTasks()
		if err != nil {
			utils.PrintErrorln("Failed to list tasks")
			return
		}

		for _, t := range tasks {
			fmt.Printf("[%d] %s (Status: %s)\n", t.ID, t.Title, t.Status)
			for _, st := range t.SubTasks {
				fmt.Printf("  - [%d] %s (%s)\n", st.ID, st.Title, st.Status)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(taskCmd)
	taskCmd.AddCommand(taskAddCmd)
	taskCmd.AddCommand(taskCompleteCmd)
	taskCmd.AddCommand(taskListCmd)

	taskAddCmd.Flags().String("intent-id", "", "ID of the linked intent")
	taskAddCmd.Flags().String("focus-id", "", "ID of the linked focus session")
}
