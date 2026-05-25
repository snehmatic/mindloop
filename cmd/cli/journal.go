package cli

import (
	"fmt"
	"time"

	cfg 	"github.com/snehmatic/mindloop/internal/config"
	"github.com/snehmatic/mindloop/internal/core/ai"
	"github.com/snehmatic/mindloop/internal/core/journal"
	"github.com/snehmatic/mindloop/internal/core/summary"
	"github.com/snehmatic/mindloop/internal/repository/appsettings"
	"github.com/snehmatic/mindloop/internal/utils"
	"github.com/snehmatic/mindloop/models"
	"github.com/spf13/cobra"
)

var (
	mood              *string
	journalService    *journal.Service
	journalSummarySvc *summary.Service
)

func getJournalSvc() *journal.Service {
	if journalService == nil {
		journalService = journal.NewService(*jRepo)
	}
	return journalService
}

func getJournalSummarySvc() *summary.Service {
	if journalSummarySvc == nil {
		journalSummarySvc = summary.NewService(*fRepo, *hRepo, *iRepo, *pRepo, *tRepo, logger)
	}
	return journalSummarySvc
}

func getJournalTimeRange(period string) (time.Time, time.Time) {
	now := time.Now()
	switch period {
	case "yearly":
		return time.Date(now.Year()-1, now.Month(), now.Day(), 0, 0, 0, 0, now.Location()), now
	case "weekly":
		end := time.Now()
		start := end.AddDate(0, 0, -7)
		start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
		return start, end
	case "daily":
		return now.Add(-24 * time.Hour), now
	default:
		return now.Add(-24 * time.Hour), now
	}
}

var generateCmd = &cobra.Command{
	Use:     "generate",
	Short:   "Auto-generate a journal entry using AI",
	Example: `mindloop journal generate -w`,
	Run: func(cmd *cobra.Command, args []string) {
		weekly, _ := cmd.Flags().GetBool("weekly")
		yearly, _ := cmd.Flags().GetBool("yearly")

		period := "daily"
		if weekly {
			period = "weekly"
		} else if yearly {
			period = "yearly"
		}

		start, end := getJournalTimeRange(period)

		report, err := getJournalSummarySvc().GenerateSummary(start, end)
		if err != nil {
			utils.PrintErrorln("Failed to generate summary report:", err)
			return
		}

		utils.PrintLoadingln("✨ Generating AI journal entry...")
		settingsRepo := appsettings.NewSQLRepository(gdb)
		aiService := ai.NewService(settingsRepo)
		generatedText, err := aiService.GenerateJournal(report)
		if err != nil {
			utils.PrintErrorln("Failed to generate journal:", err)
			return
		}

		fmt.Println("\n" + generatedText + "\n")

		fmt.Print("Would you like to save this into the journal? (y/N): ")
		var response string
		_, _ = fmt.Scanln(&response)
		if response == "y" || response == "Y" {
			title := fmt.Sprintf("AI Summary: %s", report.DateRange)
			// Assuming journal points default to 5 if config isn't read fully
			uc := cfg.GetUserConfig()
			_ = uc.ReadFromYAML()
			pts := uc.PointsConfig.Journal
			if pts == 0 {
				pts = 5
			}
			_, err = getJournalSvc().CreateEntry(title, generatedText, "reflective", pts)
			if err != nil {
				utils.PrintErrorln("Failed to save journal:", err)
			} else {
				utils.PrintSuccessln("Saved successfully!")
			}
		}
	},
}

var journalCmd = &cobra.Command{
	Use:     "journal",
	Short:   "Journal your thoughts and progress",
	Long:    `Journal your thoughts, feelings, and progress to reflect on your journey.`,
	Example: `mindloop journal new "Here goes nothing..."`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		getJournalSvc()
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
			utils.PrintWarnln("Please provide a journal title.")
			return
		}
		utils.PrintRocketln("Let's capture your thoughts! Opening your editor...")
		content, err := utils.CaptureJournalWithEditor()
		if err != nil {
			utils.PrintErrorln("Error capturing journal:", err)
			return
		}
		if content == "" {
			utils.PrintWarnln("Empty journal. Nothing saved.")
			return
		}

		utils.PrintInfoln("Saving your journal entry...")
		// Mood handling is now done in the service if empty, but we pass the flag value
		uc := cfg.GetUserConfig()
		milestoneReached, err := getJournalSvc().CreateEntry(args[0], content, *mood, uc.PointsConfig.Journal)
		if err != nil {
			utils.PrintErrorln("Failed to save journal:", err)
			return
		}

		ac.Logger.Info(fmt.Sprintf("Journal entry '%s' saved with mood '%s'.", args[0], *mood))
		utils.PrintInfoln("Your journal entry has been saved successfully!")
		utils.PrintSuccessf("Journal entry saved. (+%d pts) 🎉\n", uc.PointsConfig.Journal)
		if milestoneReached {
			utils.PrintRocketln("🏆 MILESTONE REACHED! You're on fire! 🏆")
		}
	},
}

var journalListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List all journal entries",
	Example: `mindloop journal list`,
	Aliases: []string{"l"},
	Run: func(cmd *cobra.Command, args []string) {
		utils.PrintRocketln("Fetching your journal entries...")

		entries, err := getJournalSvc().ListEntries()
		if err != nil {
			utils.PrintErrorln("Failed to retrieve journal entries:", err)
			return
		}
		if len(entries) == 0 {
			utils.PrintInfoln("No journal entries found. Try creating one with 'mindloop journal new <title>'.")
			return
		}

		utils.PrintInfoln("Your journal entries:")

		var entryViews []models.JournalEntryView
		for _, entry := range entries {
			entryViews = append(entryViews, models.ToJournalEntryView(entry))
		}
		utils.PrintTable(entryViews)

		utils.PrintInfoln("To view a specific entry, use 'mindloop journal view <id>'.")
	},
}

var journalViewCmd = &cobra.Command{
	Use:     "view",
	Short:   "View a specific journal entry",
	Example: `mindloop journal view <id>`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		entry, err := getJournalSvc().GetEntry(id)
		if err != nil {
			utils.PrintErrorln("Journal entry not found:", err)
			ac.Logger.Error("Journal entry not found", err)
			return
		}

		PrintJournalEntry(entry)
		ac.Logger.Info(fmt.Sprintf("Viewed journal entry with ID %s.", id))
	},
}

var journalDeleteCmd = &cobra.Command{
	Use:     "delete",
	Short:   "Delete a specific journal entry",
	Example: `mindloop journal delete <id>`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]

		entry, err := getJournalSvc().GetEntry(id)
		if err != nil {
			utils.PrintErrorln("Journal entry not found:", err)
			ac.Logger.Error("Journal entry not found", err)
			return
		}

		utils.PrintInfof("Are you sure you want to delete journal entry with Title '%s'?\n", entry.Title)
		utils.PrintInfoln("This action cannot be undone. Type 'yes' to confirm.")
		var confirmation string
		_, _ = fmt.Scanln(&confirmation)
		if confirmation != "yes" {
			utils.PrintWarnln("Deletion cancelled.")
			ac.Logger.Warn(fmt.Sprintf("Deletion of journal entry with ID %s cancelled by user.", id))
			return
		}

		utils.PrintRocketf("Deleting journal entry '%s'\n", entry.Title)
		err = getJournalSvc().DeleteEntry(id)
		if err != nil {
			utils.PrintErrorln("Failed to delete journal entry:", err)
			ac.Logger.Error(fmt.Sprintf("Failed to delete journal entry with ID %s", id), err)
			return
		}

		utils.PrintSuccessln("Journal entry deleted successfully!")
		ac.Logger.Info(fmt.Sprintf("Deleted journal entry with ID %s.", id))
	},
}

func init() {
	journalCmd.AddCommand(journalNewCmd)
	journalCmd.AddCommand(journalListCmd)
	journalCmd.AddCommand(journalViewCmd)
	journalCmd.AddCommand(journalDeleteCmd)

	generateCmd.Flags().BoolP("daily", "d", false, "Generate daily summary")
	generateCmd.Flags().BoolP("weekly", "w", false, "Generate weekly summary")
	generateCmd.Flags().BoolP("yearly", "y", false, "Generate yearly summary")
	journalCmd.AddCommand(generateCmd)

	rootCmd.AddCommand(journalCmd)

	mood = journalNewCmd.Flags().StringP("mood", "m", "neutral", "Set journal mood")
}

func PrintJournalEntry(entry models.JournalEntry) {
	fmt.Println("-------------------------------")
	utils.PrintInfoln("Title:", entry.Title)
	utils.PrintInfoln("Mood:", entry.Mood)
	utils.PrintLoadingln("Date:", entry.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Println("-------------------------------")
	fmt.Println(entry.Content)
	fmt.Println("-------------------------------")
	utils.PrintInfoln("End of journal entry.")
}
