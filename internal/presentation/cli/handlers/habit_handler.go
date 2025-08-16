package handlers

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/snehmatic/mindloop/internal/application/usecases"
	"github.com/snehmatic/mindloop/internal/domain/entities"
	"github.com/snehmatic/mindloop/internal/presentation/viewmodels"
	"github.com/snehmatic/mindloop/internal/shared/ui"
	"github.com/spf13/cobra"
)

// HabitHandler handles habit-related CLI commands
type HabitHandler struct {
	habitUseCase usecases.HabitUseCase
	ui           ui.Interface
}

// NewHabitHandler creates a new habit handler
func NewHabitHandler(habitUseCase usecases.HabitUseCase, ui ui.Interface) *HabitHandler {
	return &HabitHandler{
		habitUseCase: habitUseCase,
		ui:           ui,
	}
}

// CreateCommands creates all habit-related commands
func (h *HabitHandler) CreateCommands() *cobra.Command {
	// Parent habit command
	habitCmd := &cobra.Command{
		Use:     "habit",
		Short:   "Manage your habits",
		Example: `mindloop habit add "Exercise"`,
	}

	// Add subcommands
	habitCmd.AddCommand(h.createAddCommand())
	habitCmd.AddCommand(h.createListCommand())
	habitCmd.AddCommand(h.createDeleteCommand())
	habitCmd.AddCommand(h.createUpdateCommand())
	habitCmd.AddCommand(h.createLogCommands())

	return habitCmd
}

func (h *HabitHandler) createAddCommand() *cobra.Command {
	var daily, weekly, interactive bool

	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a new habit",
		Example: `mindloop habit add "Exercise" "Need to be fit!" 1 --daily
mindloop habit add -i`,
		RunE: func(cmd *cobra.Command, args []string) error {
			h.ui.ShowRocket("Great initiative! Adding a new habit...")

			var title, description string
			var targetCount int
			var interval entities.IntervalType

			if interactive {
				h.ui.ShowInfo("Interactive mode enabled for adding habit...")
				return h.handleInteractiveAdd()
			}

			// Non-interactive mode
			if len(args) < 3 {
				h.ui.ShowWarning("Please provide habit details. Ex. 'mindloop habit add <title> <description> <target_count>' --weekly or --daily(default)")
				return fmt.Errorf("missing required arguments")
			}

			title = args[0]
			description = args[1]
			var err error
			targetCount, err = strconv.Atoi(args[2])
			if err != nil {
				h.ui.ShowError("Invalid target count. Please provide a valid integer.")
				return fmt.Errorf("invalid target count: %w", err)
			}

			interval = h.getIntervalFromFlags(daily, weekly)

			habit, err := h.habitUseCase.CreateHabit(title, description, targetCount, interval)
			if err != nil {
				h.ui.ShowError(fmt.Sprintf("Failed to create habit: %v", err))
				return err
			}

			h.ui.ShowSuccess(fmt.Sprintf("Habit '%s' created successfully with ID %d!", habit.Title, habit.ID))
			return nil
		},
	}

	cmd.Flags().BoolVarP(&daily, "daily", "d", false, "Set habit as daily")
	cmd.Flags().BoolVarP(&weekly, "weekly", "w", false, "Set habit as weekly")
	cmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "Interactive mode for adding habit")

	return cmd
}

func (h *HabitHandler) createListCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "List all habits",
		Example: `mindloop habit list`,
		RunE: func(cmd *cobra.Command, args []string) error {
			habits, err := h.habitUseCase.GetAllHabits()
			if err != nil {
				h.ui.ShowError(fmt.Sprintf("Failed to get habits: %v", err))
				return err
			}

			if len(habits) == 0 {
				h.ui.ShowInfo("No habits found... Try adding one with 'mindloop habit add <title> <description> <target_count>'")
				return nil
			}

			h.ui.ShowInfo("Your habits:")
			views := viewmodels.ToHabitViews(habits)
			h.ui.ShowTable(views)
			return nil
		},
	}
}

func (h *HabitHandler) createDeleteCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "delete",
		Short:   "Delete a habit",
		Example: `mindloop habit delete <habit_id>`,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			habitID, err := strconv.ParseUint(args[0], 10, 32)
			if err != nil {
				h.ui.ShowError("Invalid habit ID. Please provide a valid number.")
				return fmt.Errorf("invalid habit ID: %w", err)
			}

			// Get habit details for confirmation
			habit, err := h.habitUseCase.GetHabit(uint(habitID))
			if err != nil {
				h.ui.ShowError(fmt.Sprintf("Habit not found: %v", err))
				return err
			}

			if !h.ui.ConfirmAction(fmt.Sprintf("Are you sure you want to delete habit '%s'?", habit.Title)) {
				h.ui.ShowWarning("Deletion cancelled.")
				return nil
			}

			if err := h.habitUseCase.DeleteHabit(uint(habitID)); err != nil {
				h.ui.ShowError(fmt.Sprintf("Failed to delete habit: %v", err))
				return err
			}

			h.ui.ShowSuccess(fmt.Sprintf("Habit '%s' deleted successfully!", habit.Title))
			return nil
		},
	}
}

func (h *HabitHandler) createUpdateCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "update",
		Short:   "Update a habit",
		Example: `mindloop habit update <habit_id>`,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			habitID, err := strconv.ParseUint(args[0], 10, 32)
			if err != nil {
				h.ui.ShowError("Invalid habit ID. Please provide a valid number.")
				return fmt.Errorf("invalid habit ID: %w", err)
			}

			habit, err := h.habitUseCase.GetHabit(uint(habitID))
			if err != nil {
				h.ui.ShowError(fmt.Sprintf("Habit not found: %v", err))
				return err
			}

			return h.handleInteractiveUpdate(habit)
		},
	}
}

func (h *HabitHandler) createLogCommands() *cobra.Command {
	logCmd := &cobra.Command{
		Use:   "log",
		Short: "Log habit completion",
	}

	// Log habit
	logCmd.AddCommand(&cobra.Command{
		Use:     "add",
		Short:   "Log habit completion",
		Example: `mindloop habit log add <habit_id> <actual_count>`,
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			habitID, err := strconv.ParseUint(args[0], 10, 32)
			if err != nil {
				h.ui.ShowError("Invalid habit ID. Please provide a valid number.")
				return fmt.Errorf("invalid habit ID: %w", err)
			}

			actualCount, err := strconv.Atoi(args[1])
			if err != nil {
				h.ui.ShowError("Invalid actual count. Please provide a valid integer.")
				return fmt.Errorf("invalid actual count: %w", err)
			}

			if err := h.habitUseCase.LogHabit(uint(habitID), actualCount); err != nil {
				h.ui.ShowError(fmt.Sprintf("Failed to log habit: %v", err))
				return err
			}

			h.ui.ShowSuccess("Habit logged successfully!")
			return nil
		},
	})

	// Show logs
	logCmd.AddCommand(&cobra.Command{
		Use:     "show",
		Short:   "Show habit logs",
		Example: `mindloop habit log show [habit_id]`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var logs []*entities.HabitLog
			var err error

			if len(args) > 0 {
				habitID, err := strconv.ParseUint(args[0], 10, 32)
				if err != nil {
					h.ui.ShowError("Invalid habit ID. Please provide a valid number.")
					return fmt.Errorf("invalid habit ID: %w", err)
				}
				logs, err = h.habitUseCase.GetHabitLogs(uint(habitID))
			} else {
				logs, err = h.habitUseCase.GetAllHabitLogs()
			}

			if err != nil {
				h.ui.ShowError(fmt.Sprintf("Failed to get habit logs: %v", err))
				return err
			}

			if len(logs) == 0 {
				h.ui.ShowInfo("No habit logs found.")
				return nil
			}

			h.ui.ShowInfo("Habit logs:")
			views := viewmodels.ToHabitLogViews(logs)
			h.ui.ShowTable(views)
			return nil
		},
	})

	return logCmd
}

func (h *HabitHandler) getIntervalFromFlags(daily, weekly bool) entities.IntervalType {
	if daily {
		return entities.Daily
	} else if weekly {
		return entities.Weekly
	}
	h.ui.ShowInfo("Defaulting to daily interval. Use -w or -d to set weekly or daily respectively.")
	return entities.Daily
}

func (h *HabitHandler) handleInteractiveAdd() error {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter habit title: ")
	scanner.Scan()
	title := strings.TrimSpace(scanner.Text())

	fmt.Print("Enter habit description: ")
	scanner.Scan()
	description := strings.TrimSpace(scanner.Text())

	fmt.Print("Enter target count: ")
	scanner.Scan()
	targetCountStr := strings.TrimSpace(scanner.Text())
	targetCount, err := strconv.Atoi(targetCountStr)
	if err != nil {
		return fmt.Errorf("invalid target count: %w", err)
	}

	fmt.Print("Enter interval (daily/weekly) [daily]: ")
	scanner.Scan()
	intervalStr := strings.TrimSpace(scanner.Text())
	if intervalStr == "" {
		intervalStr = "daily"
	}

	var interval entities.IntervalType
	switch intervalStr {
	case "daily":
		interval = entities.Daily
	case "weekly":
		interval = entities.Weekly
	default:
		return fmt.Errorf("invalid interval: %s", intervalStr)
	}

	habit, err := h.habitUseCase.CreateHabit(title, description, targetCount, interval)
	if err != nil {
		return err
	}

	h.ui.ShowSuccess(fmt.Sprintf("Habit '%s' created successfully with ID %d!", habit.Title, habit.ID))
	return nil
}

func (h *HabitHandler) handleInteractiveUpdate(habit *entities.Habit) error {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Printf("Current title: %s\nEnter new title (or press Enter to keep current): ", habit.Title)
	scanner.Scan()
	newTitle := strings.TrimSpace(scanner.Text())
	if newTitle != "" {
		habit.Title = newTitle
	}

	fmt.Printf("Current description: %s\nEnter new description (or press Enter to keep current): ", habit.Description)
	scanner.Scan()
	newDescription := strings.TrimSpace(scanner.Text())
	if newDescription != "" {
		habit.Description = newDescription
	}

	fmt.Printf("Current target count: %d\nEnter new target count (or press Enter to keep current): ", habit.TargetCount)
	scanner.Scan()
	newTargetCountStr := strings.TrimSpace(scanner.Text())
	if newTargetCountStr != "" {
		newTargetCount, err := strconv.Atoi(newTargetCountStr)
		if err != nil {
			return fmt.Errorf("invalid target count: %w", err)
		}
		habit.TargetCount = newTargetCount
	}

	if err := h.habitUseCase.UpdateHabit(habit); err != nil {
		return err
	}

	h.ui.ShowSuccess(fmt.Sprintf("Habit '%s' updated successfully!", habit.Title))
	return nil
}
