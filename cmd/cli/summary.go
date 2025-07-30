package cli

import (
	"fmt"
	"time"

	. "github.com/snehmatic/mindloop/internal/utils"
	"github.com/snehmatic/mindloop/models"
	"github.com/spf13/cobra"
)

var (
	year       *bool
	month      *bool
	week       *bool
	day        *bool
	rangeQuery = "CreatedAt >= ? AND CreatedAt <= ?"
)

var summaryCmd = &cobra.Command{
	Use:     "summary",
	Short:   "Get a summary of your productivity",
	Example: `mindloop summary --year`,
	Aliases: []string{"sum", "stats"},
	Run: func(cmd *cobra.Command, args []string) {
		PrettyPrintBanner()
		ac.Logger.Info().Msg("Requested productivity summary")

		start, end := GetTimeRangeFromFlags()

		PrintInfoln("Generating summary from", start.Format("02-Jan-2006"), "to", end.Format("02-Jan-2006"))

		summary, err := GenerateSummary(start, end)
		if err != nil {
			PrintErrorln("Error generating summary:", err)
			ac.Logger.Error().Msgf("Error generating summary: %v", err)
			return
		}
		PrintSummary(summary)

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

func GenerateSummary(start, end time.Time) (models.SummaryReport, error) {

	focusStats, err := GetFocusStats(start, end)
	if err != nil {
		PrintErrorln("Error fetching focus stats:", err)
		ac.Logger.Error().Msgf("Error fetching focus stats: %v", err)
		return models.SummaryReport{}, err
	}

	habitStats, err := GetHabitStats(start, end)
	if err != nil {
		PrintErrorln("Error fetching habit stats:", err)
		ac.Logger.Error().Msgf("Error fetching habit stats: %v", err)
		return models.SummaryReport{}, err
	}

	intentStats, err := GetIntentStats(start, end)
	if err != nil {
		PrintErrorln("Error fetching intent stats:", err)
		ac.Logger.Error().Msgf("Error fetching intent stats: %v", err)
		return models.SummaryReport{}, err
	}

	return models.SummaryReport{
		DateRange: fmt.Sprintf("%s to %s", start.Format("02-Jan-2006"), end.Format("02-Jan-2006")),
		Focus:     focusStats,
		Habits:    habitStats,
		Intents:   intentStats,
	}, nil
}

func GetFocusStats(start, end time.Time) (models.FocusStats, error) {
	var sessions []models.FocusSession

	if err := gdb.Where(rangeQuery, start, end).Find(&sessions).Error; err != nil {
		return models.FocusStats{}, err
	}
	if len(sessions) == 0 {
		return models.FocusStats{}, nil // No sessions found
	}
	totalDuration := 0.0
	longestSession := 0.0
	for _, session := range sessions {
		totalDuration += session.Duration
		if session.Duration > longestSession {
			longestSession = session.Duration
		}
	}
	return models.FocusStats{
		TotalSessions:  len(sessions),
		TotalDuration:  FormatMinutes(totalDuration),
		LongestSession: FormatMinutes(longestSession),
	}, nil
}

func GetHabitStats(start, end time.Time) ([]models.HabitStats, error) {

	var habits []models.Habit
	if err := gdb.Find(&habits).Error; err != nil {
		return nil, err
	}
	if len(habits) == 0 {
		return nil, nil // No habits found
	}

	var habitLogs []models.HabitLog
	res := gdb.Where(rangeQuery, start, end).Order("CreatedAt DESC").Find(&habitLogs)
	if res.Error != nil {
		ac.Logger.Error().
			Err(res.Error).
			Msg("Failed to retrieve habit logs")
		PrintErrorln("Failed to retrieve habit logs:", res.Error)
		return nil, res.Error
	}

	if len(habitLogs) == 0 {
		ac.Logger.Info().Msg("No habit logs found")
		return nil, nil
	}

	totalCompletedLogsForHabit := 0
	totalLogsForHabit := 0

	var stats []models.HabitStats
	for _, habit := range habits {
		for _, log := range habitLogs {
			if log.HabitID == habit.ID {
				totalLogsForHabit++
				if log.ActualCount >= log.TargetCount {
					totalCompletedLogsForHabit++
				}
			}
		}
		stats = append(stats, models.HabitStats{
			HabitName:      habit.Title,
			CompletionRate: float64(totalCompletedLogsForHabit) * 100 / float64(totalLogsForHabit),
			LogsTracked:    totalLogsForHabit,
			LogsCompleted:  totalCompletedLogsForHabit,
		})
	}
	return stats, nil
}

func GetIntentStats(start, end time.Time) ([]models.IntentStats, error) {
	var intents []models.Intent
	if err := gdb.Where(rangeQuery, start, end).Find(&intents).Error; err != nil {
		return nil, err
	}

	if len(intents) == 0 {
		return nil, nil // No intents found
	}

	var stats []models.IntentStats
	for _, intent := range intents {
		stats = append(stats, models.IntentStats{
			IntentName: intent.Name,
			Status:     intent.Status,
		})
	}
	return stats, nil
}
