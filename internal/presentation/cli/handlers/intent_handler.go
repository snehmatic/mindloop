package handlers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/snehmatic/mindloop/internal/application/usecases"
	"github.com/snehmatic/mindloop/internal/domain/entities"
	"github.com/snehmatic/mindloop/internal/presentation/viewmodels"
	"github.com/snehmatic/mindloop/internal/shared/ui"
	"github.com/spf13/cobra"
)

// IntentHandler handles intent-related CLI commands
type IntentHandler struct {
	intentUseCase usecases.IntentUseCase
	ui            ui.Interface
}

// NewIntentHandler creates a new intent handler
func NewIntentHandler(intentUseCase usecases.IntentUseCase, ui ui.Interface) *IntentHandler {
	return &IntentHandler{
		intentUseCase: intentUseCase,
		ui:            ui,
	}
}

// CreateCommands creates all intent-related commands
func (i *IntentHandler) CreateCommands() *cobra.Command {
	// Parent intent command
	intentCmd := &cobra.Command{
		Use:     "intent",
		Short:   "Manage your intents",
		Example: `mindloop intent start "Complete project documentation"`,
	}

	// Add subcommands
	intentCmd.AddCommand(i.createStartCommand())
	intentCmd.AddCommand(i.createEndCommand())
	intentCmd.AddCommand(i.createListCommand())
	intentCmd.AddCommand(i.createDeleteCommand())

	return intentCmd
}

func (i *IntentHandler) createStartCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start a new intent",
		Example: `mindloop intent start "Complete project documentation"
mindloop intent start "Learn Go programming"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			i.ui.ShowRocket("Starting a new intent...")

			if len(args) < 1 {
				i.ui.ShowWarning("Please provide a name for your intent.")
				return fmt.Errorf("missing intent name")
			}

			name := strings.Join(args, " ")

			intent, err := i.intentUseCase.StartIntent(name)
			if err != nil {
				i.ui.ShowError(fmt.Sprintf("Failed to start intent: %v", err))
				return err
			}

			i.ui.ShowSuccess(fmt.Sprintf("Intent '%s' started successfully with ID %d!", intent.Name, intent.ID))
			i.ui.ShowInfo("Use 'mindloop intent end <id>' to mark the intent as done when completed.")
			return nil
		},
	}

	return cmd
}

func (i *IntentHandler) createEndCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "end",
		Short:   "End an active intent",
		Example: `mindloop intent end 1`,
		RunE: func(cmd *cobra.Command, args []string) error {
			i.ui.ShowRocket("Ending intent...")

			if len(args) < 1 {
				i.ui.ShowWarning("Please provide the intent ID to end.")
				return fmt.Errorf("missing intent ID")
			}

			id, err := strconv.ParseUint(args[0], 10, 32)
			if err != nil {
				i.ui.ShowError("Invalid intent ID. Please provide a valid number.")
				return fmt.Errorf("invalid intent ID: %w", err)
			}

			intent, err := i.intentUseCase.EndIntent(uint(id))
			if err != nil {
				i.ui.ShowError(fmt.Sprintf("Failed to end intent: %v", err))
				return err
			}

			i.ui.ShowSuccess(fmt.Sprintf("Intent '%s' completed successfully!", intent.Name))
			return nil
		},
	}

	return cmd
}

func (i *IntentHandler) createListCommand() *cobra.Command {
	var active bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List intents",
		Example: `mindloop intent list
mindloop intent list --active`,
		RunE: func(cmd *cobra.Command, args []string) error {
			i.ui.ShowRocket("Fetching intents...")

			var intents []*entities.Intent
			var err error

			if active {
				intents, err = i.intentUseCase.GetActiveIntents()
			} else {
				intents, err = i.intentUseCase.GetAllIntents()
			}

			if err != nil {
				i.ui.ShowError(fmt.Sprintf("Failed to fetch intents: %v", err))
				return err
			}

			if len(intents) == 0 {
				if active {
					i.ui.ShowInfo("No active intents found.")
				} else {
					i.ui.ShowInfo("No intents found.")
				}
				return nil
			}

			i.ui.ShowInfo(fmt.Sprintf("Found %d intent(s):", len(intents)))

			for _, intent := range intents {
				vm := viewmodels.ToIntentView(intent)
				i.ui.ShowInfo(fmt.Sprintf("ID: %d | Name: %s | Status: %s | Created: %s",
					vm.ID, vm.Name, vm.Status, vm.CreatedAt))
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&active, "active", "a", false, "Show only active intents")
	return cmd
}

func (i *IntentHandler) createDeleteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete",
		Short:   "Delete an intent",
		Example: `mindloop intent delete 1`,
		RunE: func(cmd *cobra.Command, args []string) error {
			i.ui.ShowRocket("Deleting intent...")

			if len(args) < 1 {
				i.ui.ShowWarning("Please provide the intent ID to delete.")
				return fmt.Errorf("missing intent ID")
			}

			id, err := strconv.ParseUint(args[0], 10, 32)
			if err != nil {
				i.ui.ShowError("Invalid intent ID. Please provide a valid number.")
				return fmt.Errorf("invalid intent ID: %w", err)
			}

			err = i.intentUseCase.DeleteIntent(uint(id))
			if err != nil {
				i.ui.ShowError(fmt.Sprintf("Failed to delete intent: %v", err))
				return err
			}

			i.ui.ShowSuccess(fmt.Sprintf("Intent with ID %d deleted successfully!", id))
			return nil
		},
	}

	return cmd
}
