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

// JournalHandler handles journal-related CLI commands
type JournalHandler struct {
	journalUseCase usecases.JournalUseCase
	ui             ui.Interface
}

// NewJournalHandler creates a new journal handler
func NewJournalHandler(journalUseCase usecases.JournalUseCase, ui ui.Interface) *JournalHandler {
	return &JournalHandler{
		journalUseCase: journalUseCase,
		ui:             ui,
	}
}

// CreateCommands creates all journal-related commands
func (j *JournalHandler) CreateCommands() *cobra.Command {
	// Parent journal command
	journalCmd := &cobra.Command{
		Use:     "journal",
		Short:   "Manage your journal entries",
		Example: `mindloop journal create "Today's thoughts" "Had a great day!" --mood happy`,
	}

	// Add subcommands
	journalCmd.AddCommand(j.createCreateCommand())
	journalCmd.AddCommand(j.createListCommand())
	journalCmd.AddCommand(j.createUpdateCommand())
	journalCmd.AddCommand(j.createDeleteCommand())

	return journalCmd
}

func (j *JournalHandler) createCreateCommand() *cobra.Command {
	var mood string
	var interactive bool

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new journal entry",
		Example: `mindloop journal create "Today's thoughts" "Had a great day!" --mood happy
mindloop journal create -i`,
		RunE: func(cmd *cobra.Command, args []string) error {
			j.ui.ShowRocket("Creating a new journal entry...")

			var title, content string
			var moodEnum entities.Mood

			if interactive {
				j.ui.ShowInfo("Interactive mode enabled for creating journal entry...")
				return j.handleInteractiveCreate()
			}

			// Non-interactive mode
			if len(args) < 2 {
				j.ui.ShowWarning("Please provide title and content. Ex. 'mindloop journal create <title> <content>' --mood <mood>")
				return fmt.Errorf("missing required arguments")
			}

			title = args[0]
			content = strings.Join(args[1:], " ")

			// Parse mood
			if mood == "" {
				mood = "neutral"
			}
			moodEnum = entities.Mood(mood)
			if !entities.IsValidMood(moodEnum) {
				j.ui.ShowError(fmt.Sprintf("Invalid mood: %s. Valid moods are: %v", mood, entities.AllMoods))
				return fmt.Errorf("invalid mood")
			}

			entry, err := j.journalUseCase.CreateEntry(title, content, moodEnum)
			if err != nil {
				j.ui.ShowError(fmt.Sprintf("Failed to create journal entry: %v", err))
				return err
			}

			j.ui.ShowSuccess(fmt.Sprintf("Journal entry '%s' created successfully with ID %d!", entry.Title, entry.ID))
			return nil
		},
	}

	cmd.Flags().StringVarP(&mood, "mood", "m", "neutral", fmt.Sprintf("Mood for the entry (%v)", entities.AllMoods))
	cmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "Interactive mode for creating journal entry")
	return cmd
}

func (j *JournalHandler) handleInteractiveCreate() error {
	// Get title
	j.ui.ShowInfo("Enter the title for your journal entry:")
	var title string
	fmt.Scanln(&title)
	if title == "" {
		return fmt.Errorf("title cannot be empty")
	}

	// Get content using editor
	j.ui.ShowInfo("Opening editor for journal content...")
	content, err := j.ui.CaptureJournalWithEditor()
	if err != nil {
		j.ui.ShowError(fmt.Sprintf("Failed to capture journal content: %v", err))
		return err
	}
	if content == "" {
		return fmt.Errorf("content cannot be empty")
	}

	// Get mood
	j.ui.ShowInfo(fmt.Sprintf("Select mood (%v):", entities.AllMoods))
	var moodInput string
	fmt.Scanln(&moodInput)
	if moodInput == "" {
		moodInput = "neutral"
	}

	mood := entities.Mood(moodInput)
	if !entities.IsValidMood(mood) {
		j.ui.ShowError(fmt.Sprintf("Invalid mood: %s. Using neutral.", moodInput))
		mood = entities.MoodNeutral
	}

	entry, err := j.journalUseCase.CreateEntry(title, content, mood)
	if err != nil {
		j.ui.ShowError(fmt.Sprintf("Failed to create journal entry: %v", err))
		return err
	}

	j.ui.ShowSuccess(fmt.Sprintf("Journal entry '%s' created successfully with ID %d!", entry.Title, entry.ID))
	return nil
}

func (j *JournalHandler) createListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List journal entries",
		Example: `mindloop journal list`,
		RunE: func(cmd *cobra.Command, args []string) error {
			j.ui.ShowRocket("Fetching journal entries...")

			entries, err := j.journalUseCase.GetAllEntries()
			if err != nil {
				j.ui.ShowError(fmt.Sprintf("Failed to fetch journal entries: %v", err))
				return err
			}

			if len(entries) == 0 {
				j.ui.ShowInfo("No journal entries found.")
				return nil
			}

			j.ui.ShowInfo(fmt.Sprintf("Found %d journal entry(ies):", len(entries)))

			for _, entry := range entries {
				vm := viewmodels.ToJournalEntryView(entry)
				j.ui.ShowInfo(fmt.Sprintf("ID: %d | Title: %s | Mood: %s | Preview: %s | Created: %s",
					vm.ID, vm.Title, vm.Mood, vm.Preview, vm.CreatedAt))
			}

			return nil
		},
	}

	return cmd
}

func (j *JournalHandler) createUpdateCommand() *cobra.Command {
	var mood string

	cmd := &cobra.Command{
		Use:     "update",
		Short:   "Update a journal entry",
		Example: `mindloop journal update 1 "Updated title" "Updated content" --mood happy`,
		RunE: func(cmd *cobra.Command, args []string) error {
			j.ui.ShowRocket("Updating journal entry...")

			if len(args) < 3 {
				j.ui.ShowWarning("Please provide ID, title, and content. Ex. 'mindloop journal update <id> <title> <content>' --mood <mood>")
				return fmt.Errorf("missing required arguments")
			}

			id, err := strconv.ParseUint(args[0], 10, 32)
			if err != nil {
				j.ui.ShowError("Invalid journal entry ID. Please provide a valid number.")
				return fmt.Errorf("invalid journal entry ID: %w", err)
			}

			title := args[1]
			content := strings.Join(args[2:], " ")

			// Get existing entry
			entry, err := j.journalUseCase.GetEntry(uint(id))
			if err != nil {
				j.ui.ShowError(fmt.Sprintf("Failed to get journal entry: %v", err))
				return err
			}

			// Update fields
			entry.Title = title
			if err := entry.UpdateContent(content); err != nil {
				j.ui.ShowError(fmt.Sprintf("Failed to update content: %v", err))
				return err
			}

			// Update mood if provided
			if mood != "" {
				moodEnum := entities.Mood(mood)
				if !entities.IsValidMood(moodEnum) {
					j.ui.ShowError(fmt.Sprintf("Invalid mood: %s. Valid moods are: %v", mood, entities.AllMoods))
					return fmt.Errorf("invalid mood")
				}
				if err := entry.UpdateMood(moodEnum); err != nil {
					j.ui.ShowError(fmt.Sprintf("Failed to update mood: %v", err))
					return err
				}
			}

			// Save changes
			if err := j.journalUseCase.UpdateEntry(entry); err != nil {
				j.ui.ShowError(fmt.Sprintf("Failed to update journal entry: %v", err))
				return err
			}

			j.ui.ShowSuccess(fmt.Sprintf("Journal entry '%s' updated successfully!", entry.Title))
			return nil
		},
	}

	cmd.Flags().StringVarP(&mood, "mood", "m", "", fmt.Sprintf("Update mood (%v)", entities.AllMoods))
	return cmd
}

func (j *JournalHandler) createDeleteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete",
		Short:   "Delete a journal entry",
		Example: `mindloop journal delete 1`,
		RunE: func(cmd *cobra.Command, args []string) error {
			j.ui.ShowRocket("Deleting journal entry...")

			if len(args) < 1 {
				j.ui.ShowWarning("Please provide the journal entry ID to delete.")
				return fmt.Errorf("missing journal entry ID")
			}

			id, err := strconv.ParseUint(args[0], 10, 32)
			if err != nil {
				j.ui.ShowError("Invalid journal entry ID. Please provide a valid number.")
				return fmt.Errorf("invalid journal entry ID: %w", err)
			}

			err = j.journalUseCase.DeleteEntry(uint(id))
			if err != nil {
				j.ui.ShowError(fmt.Sprintf("Failed to delete journal entry: %v", err))
				return err
			}

			j.ui.ShowSuccess(fmt.Sprintf("Journal entry with ID %d deleted successfully!", id))
			return nil
		},
	}

	return cmd
}
