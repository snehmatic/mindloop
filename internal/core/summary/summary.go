package summary

import (
	"fmt"
	"time"

	"github.com/snehmatic/mindloop/internal/utils"
	"github.com/snehmatic/mindloop/models"
	"gorm.io/gorm"
)

type Service struct {
	DB *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{DB: db}
}

func (s *Service) GenerateSummary(start, end time.Time) (models.SummaryReport, error) {
	focusStats, err := s.GetFocusStats(start, end)
	if err != nil {
		return models.SummaryReport{}, err
	}

	habitStats, err := s.GetHabitStats(start, end)
	if err != nil {
		return models.SummaryReport{}, err
	}

	intentStats, err := s.GetIntentStats(start, end)
	if err != nil {
		return models.SummaryReport{}, err
	}

	return models.SummaryReport{
		DateRange: fmt.Sprintf("%s to %s", start.Format("02-Jan-2006"), end.Format("02-Jan-2006")),
		Focus:     focusStats,
		Habits:    habitStats,
		Intents:   intentStats,
	}, nil
}

func (s *Service) GetFocusStats(start, end time.Time) (models.FocusStats, error) {
	var sessions []models.FocusSession
	rangeQuery := "created_at >= ? AND created_at <= ?"

	if err := s.DB.Where(rangeQuery, start, end).Find(&sessions).Error; err != nil {
		return models.FocusStats{}, err
	}
	if len(sessions) == 0 {
		return models.FocusStats{}, nil
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
		LongestSession: utils.FormatMinutes(longestSession),
	}, nil
}

func (s *Service) GetHabitStats(start, end time.Time) ([]models.HabitStats, error) {
	var habits []models.Habit
	if err := s.DB.Find(&habits).Error; err != nil {
		return nil, err
	}
	if len(habits) == 0 {
		return nil, nil
	}

	var habitLogs []models.HabitLog
	rangeQuery := "created_at >= ? AND created_at <= ?"
	if err := s.DB.Where(rangeQuery, start, end).Order("created_at DESC").Find(&habitLogs).Error; err != nil {
		return nil, err
	}

	if len(habitLogs) == 0 {
		return nil, nil
	}

	totalCompletedLogsForHabit := 0
	totalLogsForHabit := 0

	var stats []models.HabitStats
	for _, habit := range habits {
		totalCompletedLogsForHabit = 0
		totalLogsForHabit = 0
		for _, log := range habitLogs {
			if log.HabitID == habit.ID {
				totalLogsForHabit++
				if log.ActualCount >= log.TargetCount {
					totalCompletedLogsForHabit++
				}
			}
		}
		if totalLogsForHabit > 0 {
			stats = append(stats, models.HabitStats{
				HabitName:      habit.Title,
				CompletionRate: float64(totalCompletedLogsForHabit) * 100 / float64(totalLogsForHabit),
				LogsTracked:    totalLogsForHabit,
				LogsCompleted:  totalCompletedLogsForHabit,
			})
		}
	}
	return stats, nil
}

func (s *Service) GetIntentStats(start, end time.Time) ([]models.IntentStats, error) {
	var intents []models.Intent
	rangeQuery := "created_at >= ? AND created_at <= ?"
	if err := s.DB.Where(rangeQuery, start, end).Find(&intents).Error; err != nil {
		return nil, err
	}

	if len(intents) == 0 {
		return nil, nil
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
