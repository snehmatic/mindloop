package usecases

import (
	"fmt"
	"time"

	"github.com/snehmatic/mindloop/internal/domain/ports"
	"github.com/snehmatic/mindloop/internal/shared/utils"
)

// SummaryStats represents aggregated statistics
type SummaryStats struct {
	DateRange string
	Focus     FocusStats
	Habits    []HabitStats
	Intents   []IntentStats
}

type FocusStats struct {
	TotalSessions  int
	TotalDuration  string
	LongestSession string
}

type HabitStats struct {
	HabitName      string
	CompletionRate float64
	LogsTracked    int
	LogsCompleted  int
}

type IntentStats struct {
	IntentName string
	Status     string
}

// SummaryUseCase defines the interface for summary use cases
type SummaryUseCase interface {
	GenerateSummary(start, end time.Time) (*SummaryStats, error)
	GetDailySummary() (*SummaryStats, error)
	GetWeeklySummary() (*SummaryStats, error)
	GetMonthlySummary() (*SummaryStats, error)
	GetYearlySummary() (*SummaryStats, error)
}

type summaryUseCase struct {
	focusRepo    ports.FocusSessionRepository
	habitRepo    ports.HabitRepository
	habitLogRepo ports.HabitLogRepository
	intentRepo   ports.IntentRepository
}

// NewSummaryUseCase creates a new summary use case
func NewSummaryUseCase(
	focusRepo ports.FocusSessionRepository,
	habitRepo ports.HabitRepository,
	habitLogRepo ports.HabitLogRepository,
	intentRepo ports.IntentRepository,
) SummaryUseCase {
	return &summaryUseCase{
		focusRepo:    focusRepo,
		habitRepo:    habitRepo,
		habitLogRepo: habitLogRepo,
		intentRepo:   intentRepo,
	}
}

// GenerateSummary generates a summary for the given date range
func (s *summaryUseCase) GenerateSummary(start, end time.Time) (*SummaryStats, error) {
	focusStats, err := s.getFocusStats(start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get focus stats: %w", err)
	}

	habitStats, err := s.getHabitStats(start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get habit stats: %w", err)
	}

	intentStats, err := s.getIntentStats(start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get intent stats: %w", err)
	}

	return &SummaryStats{
		DateRange: fmt.Sprintf("%s to %s", start.Format("02-Jan-2006"), end.Format("02-Jan-2006")),
		Focus:     *focusStats,
		Habits:    habitStats,
		Intents:   intentStats,
	}, nil
}

// GetDailySummary generates a summary for the last 24 hours
func (s *summaryUseCase) GetDailySummary() (*SummaryStats, error) {
	end := time.Now()
	start := end.Add(-24 * time.Hour)
	return s.GenerateSummary(start, end)
}

// GetWeeklySummary generates a summary for the last week
func (s *summaryUseCase) GetWeeklySummary() (*SummaryStats, error) {
	end := time.Now()
	start := end.AddDate(0, 0, -7)
	start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
	return s.GenerateSummary(start, end)
}

// GetMonthlySummary generates a summary for the current month
func (s *summaryUseCase) GetMonthlySummary() (*SummaryStats, error) {
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	end := now.AddDate(0, 1, -now.Day())
	return s.GenerateSummary(start, end)
}

// GetYearlySummary generates a summary for the last year
func (s *summaryUseCase) GetYearlySummary() (*SummaryStats, error) {
	now := time.Now()
	start := time.Date(now.Year()-1, now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	return s.GenerateSummary(start, now)
}

func (s *summaryUseCase) getFocusStats(start, end time.Time) (*FocusStats, error) {
	sessions, err := s.focusRepo.GetByDateRange(start, end)
	if err != nil {
		return nil, err
	}

	if len(sessions) == 0 {
		return &FocusStats{}, nil
	}

	totalDuration := 0.0
	longestSession := 0.0
	for _, session := range sessions {
		duration := session.GetCurrentDuration()
		totalDuration += duration
		if duration > longestSession {
			longestSession = duration
		}
	}

	return &FocusStats{
		TotalSessions:  len(sessions),
		TotalDuration:  utils.FormatMinutes(totalDuration),
		LongestSession: utils.FormatMinutes(longestSession),
	}, nil
}

func (s *summaryUseCase) getHabitStats(start, end time.Time) ([]HabitStats, error) {
	habits, err := s.habitRepo.GetAll()
	if err != nil {
		return nil, err
	}

	if len(habits) == 0 {
		return nil, nil
	}

	habitLogs, err := s.habitLogRepo.GetByDateRange(start, end)
	if err != nil {
		return nil, err
	}

	if len(habitLogs) == 0 {
		return nil, nil
	}

	var stats []HabitStats
	for _, habit := range habits {
		totalLogs := 0
		completedLogs := 0

		for _, log := range habitLogs {
			if log.HabitID == habit.ID {
				totalLogs++
				if log.IsCompleted() {
					completedLogs++
				}
			}
		}

		if totalLogs > 0 {
			stats = append(stats, HabitStats{
				HabitName:      habit.Title,
				CompletionRate: float64(completedLogs) * 100 / float64(totalLogs),
				LogsTracked:    totalLogs,
				LogsCompleted:  completedLogs,
			})
		}
	}

	return stats, nil
}

func (s *summaryUseCase) getIntentStats(start, end time.Time) ([]IntentStats, error) {
	intents, err := s.intentRepo.GetByDateRange(start, end)
	if err != nil {
		return nil, err
	}

	if len(intents) == 0 {
		return nil, nil
	}

	var stats []IntentStats
	for _, intent := range intents {
		stats = append(stats, IntentStats{
			IntentName: intent.Name,
			Status:     string(intent.Status),
		})
	}

	return stats, nil
}
