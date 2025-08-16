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

// FocusHandler handles focus session-related CLI commands
type FocusHandler struct {
	focusUseCase usecases.FocusUseCase
	ui           ui.Interface
}

// NewFocusHandler creates a new focus handler
func NewFocusHandler(focusUseCase usecases.FocusUseCase, ui ui.Interface) *FocusHandler {
	return &FocusHandler{
		focusUseCase: focusUseCase,
		ui:           ui,
	}
}

// CreateCommands creates all focus-related commands
func (f *FocusHandler) CreateCommands() *cobra.Command {
	// Parent focus command
	focusCmd := &cobra.Command{
		Use:     "focus",
		Short:   "Manage your focus sessions",
		Example: `mindloop focus start "Complete project documentation"`,
	}

	// Add subcommands
	focusCmd.AddCommand(f.createStartCommand())
	focusCmd.AddCommand(f.createEndCommand())
	focusCmd.AddCommand(f.createPauseCommand())
	focusCmd.AddCommand(f.createResumeCommand())
	focusCmd.AddCommand(f.createRateCommand())
	focusCmd.AddCommand(f.createListCommand())
	focusCmd.AddCommand(f.createDeleteCommand())

	return focusCmd
}

func (f *FocusHandler) createStartCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start a new focus session",
		Example: `mindloop focus start "Complete project documentation"
mindloop focus start "Read technical documentation"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			f.ui.ShowRocket("Starting a new focus session...")

			if len(args) < 1 {
				f.ui.ShowWarning("Please provide a title for your focus session.")
				return fmt.Errorf("missing focus session title")
			}

			title := strings.Join(args, " ")

			session, err := f.focusUseCase.StartFocusSession(title)
			if err != nil {
				f.ui.ShowError(fmt.Sprintf("Failed to start focus session: %v", err))
				return err
			}

			f.ui.ShowSuccess(fmt.Sprintf("Focus session '%s' started successfully with ID %d!", session.Title, session.ID))
			f.ui.ShowInfo("Use 'mindloop focus end <id>' to end the session when you're done.")
			return nil
		},
	}

	return cmd
}

func (f *FocusHandler) createEndCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "end",
		Short:   "End an active focus session",
		Example: `mindloop focus end 1`,
		RunE: func(cmd *cobra.Command, args []string) error {
			f.ui.ShowRocket("Ending focus session...")

			if len(args) < 1 {
				f.ui.ShowWarning("Please provide the focus session ID to end.")
				return fmt.Errorf("missing focus session ID")
			}

			id, err := strconv.ParseUint(args[0], 10, 32)
			if err != nil {
				f.ui.ShowError("Invalid focus session ID. Please provide a valid number.")
				return fmt.Errorf("invalid focus session ID: %w", err)
			}

			session, err := f.focusUseCase.EndFocusSession(uint(id))
			if err != nil {
				f.ui.ShowError(fmt.Sprintf("Failed to end focus session: %v", err))
				return err
			}

			duration := session.GetCurrentDuration()
			f.ui.ShowSuccess(fmt.Sprintf("Focus session '%s' ended successfully!", session.Title))
			f.ui.ShowInfo(fmt.Sprintf("Duration: %.1f minutes", duration))
			f.ui.ShowInfo("Use 'mindloop focus rate <id> <rating>' to rate your session (0-10).")
			return nil
		},
	}

	return cmd
}

func (f *FocusHandler) createPauseCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pause",
		Short:   "Pause an active focus session",
		Example: `mindloop focus pause 1`,
		RunE: func(cmd *cobra.Command, args []string) error {
			f.ui.ShowRocket("Pausing focus session...")

			if len(args) < 1 {
				f.ui.ShowWarning("Please provide the focus session ID to pause.")
				return fmt.Errorf("missing focus session ID")
			}

			id, err := strconv.ParseUint(args[0], 10, 32)
			if err != nil {
				f.ui.ShowError("Invalid focus session ID. Please provide a valid number.")
				return fmt.Errorf("invalid focus session ID: %w", err)
			}

			session, err := f.focusUseCase.PauseFocusSession(uint(id))
			if err != nil {
				f.ui.ShowError(fmt.Sprintf("Failed to pause focus session: %v", err))
				return err
			}

			f.ui.ShowSuccess(fmt.Sprintf("Focus session '%s' paused successfully!", session.Title))
			f.ui.ShowInfo("Use 'mindloop focus resume <id>' to resume the session.")
			return nil
		},
	}

	return cmd
}

func (f *FocusHandler) createResumeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "resume",
		Short:   "Resume a paused focus session",
		Example: `mindloop focus resume 1`,
		RunE: func(cmd *cobra.Command, args []string) error {
			f.ui.ShowRocket("Resuming focus session...")

			if len(args) < 1 {
				f.ui.ShowWarning("Please provide the focus session ID to resume.")
				return fmt.Errorf("missing focus session ID")
			}

			id, err := strconv.ParseUint(args[0], 10, 32)
			if err != nil {
				f.ui.ShowError("Invalid focus session ID. Please provide a valid number.")
				return fmt.Errorf("invalid focus session ID: %w", err)
			}

			session, err := f.focusUseCase.ResumeFocusSession(uint(id))
			if err != nil {
				f.ui.ShowError(fmt.Sprintf("Failed to resume focus session: %v", err))
				return err
			}

			f.ui.ShowSuccess(fmt.Sprintf("Focus session '%s' resumed successfully!", session.Title))
			return nil
		},
	}

	return cmd
}

func (f *FocusHandler) createRateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "rate",
		Short:   "Rate a completed focus session (0-10)",
		Example: `mindloop focus rate 1 8`,
		RunE: func(cmd *cobra.Command, args []string) error {
			f.ui.ShowRocket("Rating focus session...")

			if len(args) < 2 {
				f.ui.ShowWarning("Please provide the focus session ID and rating (0-10).")
				return fmt.Errorf("missing focus session ID or rating")
			}

			id, err := strconv.ParseUint(args[0], 10, 32)
			if err != nil {
				f.ui.ShowError("Invalid focus session ID. Please provide a valid number.")
				return fmt.Errorf("invalid focus session ID: %w", err)
			}

			rating, err := strconv.Atoi(args[1])
			if err != nil {
				f.ui.ShowError("Invalid rating. Please provide a number between 0 and 10.")
				return fmt.Errorf("invalid rating: %w", err)
			}

			if rating < 0 || rating > 10 {
				f.ui.ShowError("Rating must be between 0 and 10.")
				return fmt.Errorf("rating out of range")
			}

			session, err := f.focusUseCase.RateFocusSession(uint(id), rating)
			if err != nil {
				f.ui.ShowError(fmt.Sprintf("Failed to rate focus session: %v", err))
				return err
			}

			f.ui.ShowSuccess(fmt.Sprintf("Focus session '%s' rated %d/10 successfully!", session.Title, rating))
			return nil
		},
	}

	return cmd
}

func (f *FocusHandler) createListCommand() *cobra.Command {
	var active bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List focus sessions",
		Example: `mindloop focus list
mindloop focus list --active`,
		RunE: func(cmd *cobra.Command, args []string) error {
			f.ui.ShowRocket("Fetching focus sessions...")

			var sessions []*entities.FocusSession
			var err error

			if active {
				sessions, err = f.focusUseCase.GetActiveFocusSessions()
			} else {
				sessions, err = f.focusUseCase.GetAllFocusSessions()
			}

			if err != nil {
				f.ui.ShowError(fmt.Sprintf("Failed to fetch focus sessions: %v", err))
				return err
			}

			if len(sessions) == 0 {
				if active {
					f.ui.ShowInfo("No active focus sessions found.")
				} else {
					f.ui.ShowInfo("No focus sessions found.")
				}
				return nil
			}

			f.ui.ShowInfo(fmt.Sprintf("Found %d focus session(s):", len(sessions)))

			for _, session := range sessions {
				vm := viewmodels.ToFocusSessionView(session)
				f.ui.ShowInfo(fmt.Sprintf("ID: %d | Title: %s | Status: %s | Duration: %s | Rating: %s",
					vm.ID, vm.Title, vm.Status, vm.Duration, vm.Rating))
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&active, "active", "a", false, "Show only active focus sessions")
	return cmd
}

func (f *FocusHandler) createDeleteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete",
		Short:   "Delete a focus session",
		Example: `mindloop focus delete 1`,
		RunE: func(cmd *cobra.Command, args []string) error {
			f.ui.ShowRocket("Deleting focus session...")

			if len(args) < 1 {
				f.ui.ShowWarning("Please provide the focus session ID to delete.")
				return fmt.Errorf("missing focus session ID")
			}

			id, err := strconv.ParseUint(args[0], 10, 32)
			if err != nil {
				f.ui.ShowError("Invalid focus session ID. Please provide a valid number.")
				return fmt.Errorf("invalid focus session ID: %w", err)
			}

			err = f.focusUseCase.DeleteFocusSession(uint(id))
			if err != nil {
				f.ui.ShowError(fmt.Sprintf("Failed to delete focus session: %v", err))
				return err
			}

			f.ui.ShowSuccess(fmt.Sprintf("Focus session with ID %d deleted successfully!", id))
			return nil
		},
	}

	return cmd
}
