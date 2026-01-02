package cli

import (
	"fmt"
	"time"

	"github.com/snehmatic/mindloop/internal/core/summary"
	. "github.com/snehmatic/mindloop/internal/utils"
	"github.com/snehmatic/mindloop/models"
	"github.com/spf13/cobra"
)

var (
	year           *bool
	month          *bool
	week           *bool
	day            *bool
	summaryService *summary.Service
)

var summaryCmd = &cobra.Command{
	Use:     "summary",
	Short:   "Get a summary of your productivity",
	Example: `mindloop summary --year`,
	Aliases: []string{"sum", "stats"},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		summaryService = summary.NewService(gdb)
	},
	Run: func(cmd *cobra.Command, args []string) {
		PrettyPrintBanner()
		ac.Logger.Info().Msg("Requested productivity summary")

		start, end := GetTimeRangeFromFlags()

		PrintInfoln("Generating summary from", start.Format("02-Jan-2006"), "to", end.Format("02-Jan-2006"))

		report, err := summaryService.GenerateSummary(start, end)
		if err != nil {
			PrintErrorln("Error generating summary:", err)
			ac.Logger.Error().Msgf("Error generating summary: %v", err)
			return
		}
		PrintSummary(report)

		ac.Logger.Info().Msgf("Generated summary from %s to %s", start.Format("02-Jan-2006"), end.Format("02-Jan-2006"))
		PrintSuccessln("Summary generated successfully!")

		PrintInfoln("You can also use -d, -w, -m, or -y flags to specify the time range for these summaries. (-d is default)")
	},
}

func init() {
	rootCmd.AddCommand(summaryCmd)
	year = summaryCmd.Flags().BoolP("year", "y", false, "Show summary for the entire year")
	month = summaryCmd.Flags().BoolP("month", "m", false, "Show summary for the current month")
	week = summaryCmd.Flags().BoolP("week", "w", false, "Show summary for the current week")
	day = summaryCmd.Flags().BoolP("day", "d", false, "Show summary for today")
}

func GetTimeRangeFromFlags() (time.Time, time.Time) {
	now := time.Now()
	if *year {
		return time.Date(now.Year()-1, now.Month(), now.Day(), 0, 0, 0, 0, now.Location()), now
	} else if *month {
		return time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()), now.AddDate(0, 1, -now.Day())
	} else if *week {
		end := time.Now()
		start := end.AddDate(0, 0, -7)
		start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
		return start, end
	} else if *day {
		return now.Add(-24 * time.Hour), now
	}
	// default range "day"
	return now.Add(-24 * time.Hour), now
}

func PrintSummary(report models.SummaryReport) {
	fmt.Println("🧠 Mindloop Summary")
	fmt.Println("🗓️  Range:", report.DateRange)

	// Intent block
	fmt.Println("\n🎯 Intent Stats")
	for _, i := range report.Intents {
		fmt.Printf("- %s: %s\n", i.IntentName, i.Status)
	}

	// Focus block
	fmt.Println("\n⏱ Focus Stats")
	fmt.Printf("- Total Sessions: %d\n", report.Focus.TotalSessions)
	fmt.Printf("- Total Duration: %s\n", report.Focus.TotalDuration)
	fmt.Printf("- Longest Session: %s\n", report.Focus.LongestSession)

	// Habit block
	fmt.Println("\n📓 Habit Stats")
	for _, h := range report.Habits {
		fmt.Printf("- %s: %.0f%% (%d/%d)\n", h.HabitName, h.CompletionRate, h.LogsCompleted, h.LogsTracked)
	}
}
