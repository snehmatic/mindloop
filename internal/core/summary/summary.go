package summary

import (
	"fmt"
	"time"

	"github.com/snehmatic/mindloop/internal/log"
	"github.com/snehmatic/mindloop/internal/repository/focus"
	"github.com/snehmatic/mindloop/internal/repository/habit"
	"github.com/snehmatic/mindloop/internal/repository/intent"
	"github.com/snehmatic/mindloop/internal/repository/point"
	"github.com/snehmatic/mindloop/internal/repository/task"
	"github.com/snehmatic/mindloop/internal/utils"
	"github.com/snehmatic/mindloop/models"
)

// Service handles the logic for generating summary reports
type Service struct {
	focusRepo  focus.Repository
	habitRepo  habit.Repository
	intentRepo intent.Repository
	pointRepo  point.Repository
	taskRepo   task.TaskRepository
	logger     log.Logger
}

// NewService creates a new summary Service instance
func NewService(
	focusRepo focus.Repository,
	habitRepo habit.Repository,
	intentRepo intent.Repository,
	pointRepo point.Repository,
	taskRepo task.TaskRepository,
	logger log.Logger,
) *Service {
	return &Service{
		focusRepo:  focusRepo,
		habitRepo:  habitRepo,
		intentRepo: intentRepo,
		pointRepo:  pointRepo,
		taskRepo:   taskRepo,
		logger:     logger,
	}
}

func (s *Service) GenerateSummary(start, end time.Time) (models.SummaryReport, error) {
	focusStats, err := s.focusRepo.GetFocusStats(start, end)
	if err != nil {
		return models.SummaryReport{}, err
	}

	habitStats, err := s.habitRepo.GetHabitStats(start, end)
	if err != nil {
		return models.SummaryReport{}, err
	}

	intentStats, err := s.intentRepo.GetIntentStats(start, end)
	if err != nil {
		return models.SummaryReport{}, err
	}

	totalPoints, err := s.pointRepo.GetTotalPoints()
	if err != nil {
		return models.SummaryReport{}, err
	}
	pointStats := models.PointStats{
		TotalPoints: totalPoints,
	}

	tasksCompleted, err := s.taskRepo.CountCompletedTasks(start, end)
	if err != nil {
		return models.SummaryReport{}, err
	}

	return models.SummaryReport{
		DateRange:      fmt.Sprintf("%s to %s", start.Format("02-Jan-2006"), end.Format("02-Jan-2006")),
		Focus:          focusStats,
		Habits:         habitStats,
		Intents:        intentStats,
		Points:         pointStats,
		TasksCompleted: tasksCompleted,
	}, nil
}

func (s *Service) GetPointSeries(start, end time.Time) ([]int, error) {
	return s.pointRepo.GetPointSeries(start, end)
}

func (s *Service) GetFocusStats(start, end time.Time) (models.FocusStats, error) {
	sessions, err := s.focusRepo.GetSessionsInRange(start, end)
	if err != nil {
		return models.FocusStats{}, err
	}
	if len(sessions) == 0 {
		return models.FocusStats{
			TotalSessions:  0,
			TotalDuration:  "0 mins",
			LongestSession: "0 mins",
		}, nil
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
		TotalDuration:  utils.FormatMinutes(totalDuration),
		RawDuration:    totalDuration,
		LongestSession: utils.FormatMinutes(longestSession),
	}, nil
}

func (s *Service) GetHabitStats(start, end time.Time) ([]models.HabitStats, error) {
	habits, err := s.habitRepo.GetAllHabits()
	if err != nil {
		return nil, err
	}
	if len(habits) == 0 {
		return []models.HabitStats{}, nil
	}

	habitLogs, err := s.habitRepo.GetHabitLogsInRange(start, end)
	if err != nil {
		return nil, err
	}

	var stats []models.HabitStats
	for _, habitEntry := range habits {
		totalCompletedLogsForHabit := 0
		totalLogsForHabit := 0
		for _, habitLog := range habitLogs {
			if habitLog.HabitID == habitEntry.ID {
				totalLogsForHabit++
				if habitLog.ActualCount >= habitLog.TargetCount {
					totalCompletedLogsForHabit++
				}
			}
		}
		if totalLogsForHabit > 0 {
			stats = append(stats, models.HabitStats{
				HabitName:      habitEntry.Title,
				CompletionRate: float64(totalCompletedLogsForHabit) * 100 / float64(totalLogsForHabit),
				LogsTracked:    totalLogsForHabit,
				LogsCompleted:  totalCompletedLogsForHabit,
			})
		}
	}
	return stats, nil
}

func (s *Service) GetIntentStats(start, end time.Time) ([]models.IntentStats, error) {
	intents, err := s.intentRepo.GetIntentsInRange(start, end)
	if err != nil {
		return nil, err
	}

	if len(intents) == 0 {
		return []models.IntentStats{}, nil
	}

	var stats []models.IntentStats
	for _, intentEntry := range intents {
		stats = append(stats, models.IntentStats{
			IntentName: intentEntry.Name,
			Status:     intentEntry.Status,
		})
	}
	return stats, nil
}

// GetFocusSeries returns daily focus duration for the given range
func (s *Service) GetFocusSeries(start, end time.Time) ([]float64, []string, error) {
	// Calculate number of days
	days := int(end.Sub(start).Hours()/24) + 1
	if days < 1 {
		days = 1
	}

	stats := make([]float64, days)
	labels := make([]string, days)

	for i := 0; i < days; i++ {
		date := start.AddDate(0, 0, i)
		labels[i] = date.Format("Jan 02")
	}

	sessions, err := s.focusRepo.GetSessionsInRange(start, end)
	if err != nil {
		return nil, nil, err
	}

	for _, session := range sessions {
		// Calculate index
		// We truncate everything to midnight to compare dates
		sessionDate := session.CreatedAt.Truncate(24 * time.Hour)
		startDate := start.Truncate(24 * time.Hour)

		diff := int(sessionDate.Sub(startDate).Hours() / 24)
		if diff >= 0 && diff < days {
			stats[diff] += session.Duration
		}
	}

	return stats, labels, nil
}

// GetHabitSeries returns daily habit completion count for the given range
func (s *Service) GetHabitSeries(start, end time.Time) ([]int, error) {
	days := int(end.Sub(start).Hours()/24) + 1
	if days < 1 {
		days = 1
	}

	stats := make([]int, days)

	habitLogs, err := s.habitRepo.GetHabitLogsInRange(start, end)
	if err != nil {
		return nil, err
	}

	for _, habitLog := range habitLogs {
		if habitLog.ActualCount >= habitLog.TargetCount {
			logDate := habitLog.CreatedAt.Truncate(24 * time.Hour)
			startDate := start.Truncate(24 * time.Hour)
			diff := int(logDate.Sub(startDate).Hours() / 24)

			if diff >= 0 && diff < days {
				stats[diff]++
			}
		}
	}

	return stats, nil
}
