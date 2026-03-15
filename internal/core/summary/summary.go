package summary

import (
	"fmt"
	"time"

	"github.com/snehmatic/mindloop/internal/core/points"
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

	totalPoints, _ := points.GetTotalPoints(s.DB)
	pointStats := models.PointStats{
		TotalPoints: totalPoints,
	}

	return models.SummaryReport{
		DateRange: fmt.Sprintf("%s to %s", start.Format("02-Jan-2006"), end.Format("02-Jan-2006")),
		Focus:     focusStats,
		Habits:    habitStats,
		Intents:   intentStats,
		Points:    pointStats,
	}, nil
}

func (s *Service) GetPointSeries(start, end time.Time) ([]int, error) {
	days := int(end.Sub(start).Hours()/24) + 1
	if days < 1 {
		days = 1
	}

	stats := make([]int, days)

	var transactions []models.PointTransaction
	if err := s.DB.Where("CreatedAt >= ? AND CreatedAt <= ?", start, end).Find(&transactions).Error; err != nil {
		return nil, err
	}

	for _, tx := range transactions {
		txDate := tx.CreatedAt.Truncate(24 * time.Hour)
		startDate := start.Truncate(24 * time.Hour)
		diff := int(txDate.Sub(startDate).Hours() / 24)

		if diff >= 0 && diff < days {
			stats[diff] += tx.Points
		}
	}

	return stats, nil
}

func (s *Service) GetFocusStats(start, end time.Time) (models.FocusStats, error) {
	var sessions []models.FocusSession
	rangeQuery := "CreatedAt >= ? AND CreatedAt <= ?"

	if err := s.DB.Where(rangeQuery, start, end).Find(&sessions).Error; err != nil {
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
	var habits []models.Habit
	if err := s.DB.Order("CreatedAt DESC").Find(&habits).Error; err != nil {
		return nil, err
	}
	if len(habits) == 0 {
		return []models.HabitStats{}, nil
	}

	var habitLogs []models.HabitLog
	rangeQuery := "CreatedAt >= ? AND CreatedAt <= ?"
	if err := s.DB.Where(rangeQuery, start, end).Order("CreatedAt DESC").Find(&habitLogs).Error; err != nil {
		return nil, err
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
	rangeQuery := "CreatedAt >= ? AND CreatedAt <= ?"
	if err := s.DB.Where(rangeQuery, start, end).Order("CreatedAt DESC").Find(&intents).Error; err != nil {
		return nil, err
	}

	if len(intents) == 0 {
		return []models.IntentStats{}, nil
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

	var sessions []models.FocusSession
	if err := s.DB.Where("CreatedAt >= ? AND CreatedAt <= ?", start, end).Find(&sessions).Error; err != nil {
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

	var logs []models.HabitLog
	if err := s.DB.Where("CreatedAt >= ? AND CreatedAt <= ?", start, end).Find(&logs).Error; err != nil {
		return nil, err
	}

	for _, log := range logs {
		if log.ActualCount >= log.TargetCount {
			logDate := log.CreatedAt.Truncate(24 * time.Hour)
			startDate := start.Truncate(24 * time.Hour)
			diff := int(logDate.Sub(startDate).Hours() / 24)

			if diff >= 0 && diff < days {
				stats[diff]++
			}
		}
	}

	return stats, nil
}
