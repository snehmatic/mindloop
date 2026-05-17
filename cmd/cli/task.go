package cli

import (
	"fmt"
	"strconv"

	"github.com/rs/zerolog"
	"github.com/snehmatic/mindloop/db"
	"github.com/snehmatic/mindloop/internal/config"
	"github.com/snehmatic/mindloop/internal/core/task"
	taskRepo "github.com/snehmatic/mindloop/internal/repository/task"
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

		taskRepository := taskRepo.NewSQLTaskRepository(database)
		svc := task.NewService(taskRepository, nil, new(zerolog.Nop()))
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
		uc := config.UserConfig{}
		_ = uc.ReadFromYAML()

		database, err := db.ConnectToDb(*appConfig)
		if err != nil {
			utils.PrintErrorln("Failed to connect to database")
			return
		}

		repo := taskRepo.NewSQLTaskRepository(database)
		svc := task.NewService(repo, &uc, new(zerolog.Nop()))
		_, err = svc.CompleteTask(uint(id), uc.PointsConfig.Task)
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

		taskRepository := taskRepo.NewSQLTaskRepository(database)
		svc := task.NewService(taskRepository, nil, new(zerolog.Nop()))
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

var subtaskCmd = &cobra.Command{
	Use:   "subtask",
	Short: "Manage subtasks",
}

var subtaskAddCmd = &cobra.Command{
	Use:   "add [task-id] [title]",
	Short: "Add a subtask to a task",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		taskID, err := strconv.ParseUint(args[0], 10, 32)
		if err != nil {
			utils.PrintErrorln("Invalid task ID")
			return
		}
		title := args[1]

		appConfig := config.GetConfig()
		database, err := db.ConnectToDb(*appConfig)
		if err != nil {
			utils.PrintErrorln("Failed to connect to database")
			return
		}

		taskRepository := taskRepo.NewSQLTaskRepository(database)
		svc := task.NewService(taskRepository, nil, new(zerolog.Nop()))
		st, err := svc.AddSubTask(uint(taskID), title)
		if err != nil {
			utils.PrintErrorln(fmt.Sprintf("Failed to add subtask: %v", err))
			return
		}

		utils.PrintSuccessln(fmt.Sprintf("Added subtask '%s' (ID: %d) to task %d", st.Title, st.ID, taskID))
	},
}

var subtaskCompleteCmd = &cobra.Command{
	Use:   "complete [id]",
	Short: "Complete a subtask",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.ParseUint(args[0], 10, 32)
		if err != nil {
			utils.PrintErrorln("Invalid subtask ID")
			return
		}

		appConfig := config.GetConfig()
		uc := config.UserConfig{}
		_ = uc.ReadFromYAML()

		database, err := db.ConnectToDb(*appConfig)
		if err != nil {
			utils.PrintErrorln("Failed to connect to database")
			return
		}

		taskRepository := taskRepo.NewSQLTaskRepository(database)
		svc := task.NewService(taskRepository, &uc, new(zerolog.Nop()))
		_, err = svc.CompleteSubTask(uint(id), uc.PointsConfig.SubTask)
		if err != nil {
			utils.PrintErrorln(fmt.Sprintf("Failed to complete subtask: %v", err))
			return
		}

		utils.PrintSuccessln(fmt.Sprintf("Subtask %d marked as completed", id))
	},
}

func init() {
	rootCmd.AddCommand(taskCmd)
	taskCmd.AddCommand(taskAddCmd)
	taskCmd.AddCommand(taskCompleteCmd)
	taskCmd.AddCommand(taskListCmd)
	taskCmd.AddCommand(subtaskCmd)
	subtaskCmd.AddCommand(subtaskAddCmd)
	subtaskCmd.AddCommand(subtaskCompleteCmd)

	taskAddCmd.Flags().String("intent-id", "", "ID of the linked intent")
	taskAddCmd.Flags().String("focus-id", "", "ID of the linked focus session")
}
