package cli

import (
	"fmt"

	"github.com/snehmatic/mindloop/internal/core/journal"
	. "github.com/snehmatic/mindloop/internal/utils"
	"github.com/snehmatic/mindloop/models"
	"github.com/spf13/cobra"
)

var (
	mood           *string
	journalService *journal.Service
)

var journalCmd = &cobra.Command{
	Use:     "journal",
	Short:   "Journal your thoughts and progress",
	Long:    `Journal your thoughts, feelings, and progress to reflect on your journey.`,
	Example: `mindloop journal new "Here goes nothing..."`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		journalService = journal.NewService(gdb)
	},
}

var journalNewCmd = &cobra.Command{
	Use:     "new",
	Short:   "Create a new journal entry using your default $EDITOR",
	Example: `mindloop journal new <title>`,
	Aliases: []string{"n", "create", "add"},
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			PrintWarnln("Please provide a journal title.")
			return
		}
		PrintRocketln("Let's capture your thoughts! Opening your editor...")
		content, err := CaptureJournalWithEditor()
		if err != nil {
			PrintErrorln("Error capturing journal:", err)
			return
		}
		if content == "" {
			PrintWarnln("Empty journal. Nothing saved.")
			return
		}

		PrintInfoln("Saving your journal entry...")
		// Mood handling is now done in the service if empty, but we pass the flag value
		err = journalService.CreateEntry(args[0], content, *mood)
		if err != nil {
			PrintErrorln("Failed to save journal:", err)
			return
		}

		ac.Logger.Info().Msgf("Journal entry '%s' saved with mood '%s'.", args[0], *mood)
		PrintInfoln("Your journal entry has been saved successfully!")
		PrintSuccessln("Journal entry saved.")
	},
}

var journalListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List all journal entries",
	Example: `mindloop journal list`,
	Aliases: []string{"l"},
	Run: func(cmd *cobra.Command, args []string) {
		PrintRocketln("Fetching your journal entries...")

		entries, err := journalService.ListEntries()
		if err != nil {
			PrintErrorln("Failed to retrieve journal entries:", err)
			return
		}
		if len(entries) == 0 {
			PrintInfoln("No journal entries found. Try creating one with 'mindloop journal new <title>'.")
			return
		}

		PrintInfoln("Your journal entries:")

		var entryViews []models.JournalEntryView
		for _, entry := range entries {
			entryViews = append(entryViews, models.ToJournalEntryView(entry))
		}
		PrintTable(entryViews)

		PrintInfoln("To view a specific entry, use 'mindloop journal view <id>'.")
	},
}

var journalViewCmd = &cobra.Command{
	Use:     "view",
	Short:   "View a specific journal entry",
	Example: `mindloop journal view <id>`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		entry, err := journalService.GetEntry(id)
		if err != nil {
			PrintErrorln("Journal entry not found:", err)
			ac.Logger.Error().Msgf("Journal entry not found: %v", err)
			return
		}

		PrintJournalEntry(entry)

		ac.Logger.Info().Msgf("Viewed journal entry with ID %s.", id)
	},
}

var journalDeleteCmd = &cobra.Command{
	Use:     "delete",
	Short:   "Delete a specific journal entry",
	Example: `mindloop journal delete <id>`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]

		entry, err := journalService.GetEntry(id)
		if err != nil {
			PrintErrorln("Journal entry not found:", err)
			ac.Logger.Error().Msgf("Journal entry not found: %v", err)
			return
		}

		PrintInfof("Are you sure you want to delete journal entry with Title '%s'?\n", entry.Title)
		PrintInfoln("This action cannot be undone. Type 'yes' to confirm.")
		var confirmation string
		fmt.Scanln(&confirmation)
		if confirmation != "yes" {
			PrintWarnln("Deletion cancelled.")
			ac.Logger.Warn().Msgf("Deletion of journal entry with ID %s cancelled by user.", id)
			return
		}

		PrintRocketf("Deleting journal entry '%s'\n", entry.Title)
		err = journalService.DeleteEntry(id)
		if err != nil {
			PrintErrorln("Failed to delete journal entry:", err)
			ac.Logger.Error().Msgf("Failed to delete journal entry with ID %s: %v", id, err)
			return
		}

		PrintSuccessln("Journal entry deleted successfully!")
		ac.Logger.Info().Msgf("Deleted journal entry with ID %s.", id)
	},
}

func init() {
	journalCmd.AddCommand(journalNewCmd)
	journalCmd.AddCommand(journalListCmd)
	journalCmd.AddCommand(journalViewCmd)
	journalCmd.AddCommand(journalDeleteCmd)
	rootCmd.AddCommand(journalCmd)

	mood = journalNewCmd.Flags().StringP("mood", "m", "neutral", "Set journal mood")
}

func PrintJournalEntry(entry models.JournalEntry) {
	fmt.Println("-------------------------------")
	PrintInfoln("Title:", entry.Title)
	PrintInfoln("Mood:", entry.Mood)
	PrintLoadingln("Date:", entry.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Println("-------------------------------")
	fmt.Println(entry.Content)
	fmt.Println("-------------------------------")
	PrintInfoln("End of journal entry.")
}
