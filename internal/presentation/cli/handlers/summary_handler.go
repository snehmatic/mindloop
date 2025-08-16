package handlers

import (
	"fmt"
	"time"

	"github.com/snehmatic/mindloop/internal/application/usecases"
	"github.com/snehmatic/mindloop/internal/shared/ui"
	"github.com/spf13/cobra"
)

// SummaryHandler handles summary-related CLI commands
type SummaryHandler struct {
	summaryUseCase usecases.SummaryUseCase
	ui             ui.Interface
}

// NewSummaryHandler creates a new summary handler
func NewSummaryHandler(summaryUseCase usecases.SummaryUseCase, ui ui.Interface) *SummaryHandler {
	return &SummaryHandler{
		summaryUseCase: summaryUseCase,
		ui:             ui,
	}
}

// CreateCommands creates all summary-related commands
func (s *SummaryHandler) CreateCommands() *cobra.Command {
	// Parent summary command
	summaryCmd := &cobra.Command{
		Use:     "summary",
		Short:   "Generate productivity summaries",
		Example: `mindloop summary daily`,
	}

	// Add subcommands
	summaryCmd.AddCommand(s.createDailyCommand())
	summaryCmd.AddCommand(s.createWeeklyCommand())
	summaryCmd.AddCommand(s.createMonthlyCommand())
	summaryCmd.AddCommand(s.createYearlyCommand())
	summaryCmd.AddCommand(s.createCustomCommand())

	return summaryCmd
}

func (s *SummaryHandler) createDailyCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "daily",
		Short:   "Generate daily summary (last 24 hours)",
		Example: `mindloop summary daily`,
		RunE: func(cmd *cobra.Command, args []string) error {
			s.ui.ShowRocket("Generating daily summary...")

			stats, err := s.summaryUseCase.GetDailySummary()
			if err != nil {
				s.ui.ShowError(fmt.Sprintf("Failed to generate daily summary: %v", err))
				return err
			}

			s.displaySummary(stats, "Daily Summary")
			return nil
		},
	}

	return cmd
}

func (s *SummaryHandler) createWeeklyCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "weekly",
		Short:   "Generate weekly summary (last 7 days)",
		Example: `mindloop summary weekly`,
		RunE: func(cmd *cobra.Command, args []string) error {
			s.ui.ShowRocket("Generating weekly summary...")

			stats, err := s.summaryUseCase.GetWeeklySummary()
			if err != nil {
				s.ui.ShowError(fmt.Sprintf("Failed to generate weekly summary: %v", err))
				return err
			}

			s.displaySummary(stats, "Weekly Summary")
			return nil
		},
	}

	return cmd
}

func (s *SummaryHandler) createMonthlyCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "monthly",
		Short:   "Generate monthly summary (last 30 days)",
		Example: `mindloop summary monthly`,
		RunE: func(cmd *cobra.Command, args []string) error {
			s.ui.ShowRocket("Generating monthly summary...")

			stats, err := s.summaryUseCase.GetMonthlySummary()
			if err != nil {
				s.ui.ShowError(fmt.Sprintf("Failed to generate monthly summary: %v", err))
				return err
			}

			s.displaySummary(stats, "Monthly Summary")
			return nil
		},
	}

	return cmd
}

func (s *SummaryHandler) createYearlyCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "yearly",
		Short:   "Generate yearly summary (last 365 days)",
		Example: `mindloop summary yearly`,
		RunE: func(cmd *cobra.Command, args []string) error {
			s.ui.ShowRocket("Generating yearly summary...")

			stats, err := s.summaryUseCase.GetYearlySummary()
			if err != nil {
				s.ui.ShowError(fmt.Sprintf("Failed to generate yearly summary: %v", err))
				return err
			}

			s.displaySummary(stats, "Yearly Summary")
			return nil
		},
	}

	return cmd
}

func (s *SummaryHandler) createCustomCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "custom",
		Short:   "Generate custom date range summary",
		Example: `mindloop summary custom "2024-01-01" "2024-01-31"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			s.ui.ShowRocket("Generating custom summary...")

			if len(args) < 2 {
				s.ui.ShowWarning("Please provide start and end dates in YYYY-MM-DD format.")
				return fmt.Errorf("missing start and end dates")
			}

			startDate, err := time.Parse("2006-01-02", args[0])
			if err != nil {
				s.ui.ShowError("Invalid start date format. Please use YYYY-MM-DD.")
				return fmt.Errorf("invalid start date: %w", err)
			}

			endDate, err := time.Parse("2006-01-02", args[1])
			if err != nil {
				s.ui.ShowError("Invalid end date format. Please use YYYY-MM-DD.")
				return fmt.Errorf("invalid end date: %w", err)
			}

			if startDate.After(endDate) {
				s.ui.ShowError("Start date cannot be after end date.")
				return fmt.Errorf("invalid date range")
			}

			// Set end date to end of day
			endDate = endDate.Add(24*time.Hour - time.Second)

			stats, err := s.summaryUseCase.GenerateSummary(startDate, endDate)
			if err != nil {
				s.ui.ShowError(fmt.Sprintf("Failed to generate custom summary: %v", err))
				return err
			}

			s.displaySummary(stats, "Custom Summary")
			return nil
		},
	}

	return cmd
}

func (s *SummaryHandler) displaySummary(stats *usecases.SummaryStats, title string) {
	s.ui.ShowInfo(fmt.Sprintf("=== %s ===", title))
	s.ui.ShowInfo(fmt.Sprintf("Date Range: %s", stats.DateRange))
	s.ui.ShowInfo("")

	// Focus Stats
	s.ui.ShowInfo("🎯 Focus Sessions:")
	s.ui.ShowInfo(fmt.Sprintf("  • Total Sessions: %d", stats.Focus.TotalSessions))
	s.ui.ShowInfo(fmt.Sprintf("  • Total Duration: %s", stats.Focus.TotalDuration))
	s.ui.ShowInfo(fmt.Sprintf("  • Longest Session: %s", stats.Focus.LongestSession))
	s.ui.ShowInfo("")

	// Habit Stats
	if len(stats.Habits) > 0 {
		s.ui.ShowInfo("📈 Habits:")
		for _, habit := range stats.Habits {
			s.ui.ShowInfo(fmt.Sprintf("  • %s: %.1f%% completion (%d/%d logs)",
				habit.HabitName, habit.CompletionRate*100, habit.LogsCompleted, habit.LogsTracked))
		}
		s.ui.ShowInfo("")
	}

	// Intent Stats
	if len(stats.Intents) > 0 {
		s.ui.ShowInfo("🎯 Intents:")
		for _, intent := range stats.Intents {
			s.ui.ShowInfo(fmt.Sprintf("  • %s: %s", intent.IntentName, intent.Status))
		}
		s.ui.ShowInfo("")
	}

	s.ui.ShowSuccess("Summary generated successfully!")
}
