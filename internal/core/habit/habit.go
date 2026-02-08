package habit

import (
	"errors"
	"time"

	"github.com/snehmatic/mindloop/models"
	"gorm.io/gorm"
)

type Service struct {
	DB *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{DB: db}
}

func (s *Service) CreateHabit(habit *models.Habit) error {
	if habit == nil {
		return errors.New("habit cannot be nil")
	}
	if err := habit.ValidateHabit(); err != nil {
		return err
	}
	return s.DB.Create(habit).Error
}

func (s *Service) DeleteHabit(id string) error {
	var habit models.Habit
	if err := s.DB.First(&habit, "id = ?", id).Error; err != nil {
		return err
	}
	return s.DB.Delete(&habit).Error
}

func (s *Service) GetHabit(id string) (*models.Habit, error) {
	var habit models.Habit
	if err := s.DB.First(&habit, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &habit, nil
}

func (s *Service) UpdateHabit(habit *models.Habit) error {
	if habit == nil {
		return errors.New("habit cannot be nil")
	}
	if err := habit.ValidateHabit(); err != nil {
		return err
	}
	return s.DB.Save(habit).Error
}

func (s *Service) ListHabits(interval models.IntervalType) ([]models.Habit, error) {
	var habits []models.Habit
	query := s.DB
	if interval != "" {
		query = query.Where("interval = ?", interval)
	}
	// Only show habits that haven't ended yet
	now := time.Now()
	query = query.Where("EndDate IS NULL OR EndDate > ?", now)

	result := query.Find(&habits)
	return habits, result.Error
}

func (s *Service) ListEndedHabits() ([]models.Habit, error) {
	var habits []models.Habit
	now := time.Now()
	result := s.DB.Where("EndDate IS NOT NULL AND EndDate <= ?", now).Find(&habits)
	return habits, result.Error
}

func (s *Service) LogHabit(habitID string) (*models.Habit, *models.HabitLog, error) {
	var habit models.Habit
	if err := s.DB.First(&habit, "id = ?", habitID).Error; err != nil {
		return nil, nil, err
	}

	var existingLog models.HabitLog
	now := time.Now()
	today := now.Truncate(24 * time.Hour)
	endedAt := today

	var res *gorm.DB
	switch habit.Interval {
	case models.Daily:
		res = s.DB.Where("HabitID = ? AND EndedAt = ?", habit.ID, today).Limit(1).Find(&existingLog)
	case models.Weekly:
		// ISO Week: Monday to Sunday
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7 // Sunday
		}
		startOfWeek := today.AddDate(0, 0, -(weekday - 1))
		endOfWeek := startOfWeek.AddDate(0, 0, 6)
		endedAt = endOfWeek
		res = s.DB.Where("HabitID = ? AND CreatedAt >= ? AND CreatedAt <= ?", habit.ID, startOfWeek, endOfWeek.Add(23*time.Hour+59*time.Minute+59*time.Second)).Limit(1).Find(&existingLog)
	}

	if res.Error != nil {
		return nil, nil, res.Error
	}

	if res.RowsAffected > 0 {
		// Log found
		if existingLog.ActualCount >= habit.TargetCount {
			return &habit, &existingLog, errors.New("habit already completed for interval")
		}

		existingLog.ActualCount++
		// If it's weekly, keep the same endedAt (end of week)
		if habit.Interval == models.Daily {
			existingLog.EndedAt = today
		}
		if err := s.DB.Save(&existingLog).Error; err != nil {
			return nil, nil, err
		}
		return &habit, &existingLog, nil
	}

	// Create new log
	habitLog := &models.HabitLog{
		HabitID:     habit.ID,
		Title:       habit.Title,
		Interval:    habit.Interval,
		TargetCount: habit.TargetCount,
		ActualCount: 1,
		EndedAt:     endedAt,
	}
	if err := s.DB.Create(habitLog).Error; err != nil {
		return nil, nil, err
	}

	return &habit, habitLog, nil
}

func (s *Service) UnlogHabit(habitID string) (*models.Habit, error) {
	var habit models.Habit
	if err := s.DB.First(&habit, "id = ?", habitID).Error; err != nil {
		return nil, err
	}

	var existingLog models.HabitLog
	now := time.Now()
	today := now.Truncate(24 * time.Hour)
	var res *gorm.DB

	switch habit.Interval {
	case models.Daily:
		res = s.DB.Where("HabitID = ? AND EndedAt = ?", habit.ID, today).First(&existingLog)
	case models.Weekly:
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		startOfWeek := today.AddDate(0, 0, -(weekday - 1))
		res = s.DB.Where("HabitID = ? AND CreatedAt >= ?", habit.ID, startOfWeek).First(&existingLog)
	}

	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("no existing log found for this habit")
		}
		return nil, res.Error
	}

	if existingLog.ActualCount <= 0 {
		return nil, errors.New("habit is already marked as undone")
	}

	existingLog.ActualCount = 0 // resetting progress for the day/week.

	if err := s.DB.Save(&existingLog).Error; err != nil {
		return nil, err
	}

	return &habit, nil
}

func (s *Service) ListHabitLogs(interval models.IntervalType) ([]models.HabitLog, error) {
	var habitLogs []models.HabitLog
	query := s.DB
	if interval != "" {
		query = query.Where("interval = ?", interval)
	}
	result := query.Order("CreatedAt DESC").Find(&habitLogs)
	return habitLogs, result.Error
}

func (s *Service) ListLogsForHabit(habitID uint) ([]models.HabitLog, error) {
	var habitLogs []models.HabitLog
	result := s.DB.Where("HabitID = ?", habitID).Order("CreatedAt ASC").Find(&habitLogs)
	return habitLogs, result.Error
}

func (s *Service) DeleteAll() error {
	// Transaction to delete both logs and habits
	return s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.HabitLog{}).Error; err != nil {
			return err
		}
		if err := tx.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.Habit{}).Error; err != nil {
			return err
		}
		return nil
	})
}
func (s *Service) CalculateStreak(habitID uint, interval models.IntervalType) (int, error) {
	if interval != models.Daily {
		return 0, nil // fast track for non-daily for now
	}

	var logs []models.HabitLog
	// Fetch all logs for this habit ordered by date descending
	if err := s.DB.Where("HabitID = ?", habitID).Order("CreatedAt desc").Find(&logs).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil
		}
		return 0, err
	}

	if len(logs) == 0 {
		return 0, nil
	}

	streak := 0
	today := time.Now().Truncate(24 * time.Hour)
	// yesterday := today.AddDate(0, 0, -1)

	// We need to check if the sequence is unbroken.
	// The most recent log could be today or yesterday to keep the streak alive.
	// If the most recent log is older than yesterday, streak is 0.

	lastLogDate := logs[0].CreatedAt.Truncate(24 * time.Hour)
	daysDiff := today.Sub(lastLogDate).Hours() / 24

	if daysDiff > 1 {
		return 0, nil
	}

	// Iterate and count
	// We expect dates to be consecutive
	expectedDate := lastLogDate
	for _, log := range logs {
		logDate := log.CreatedAt.Truncate(24 * time.Hour)

		// If multiple logs on same day (shouldn't happen with current logic but safeguards), skip
		if logDate.Equal(expectedDate) {
			if log.ActualCount >= log.TargetCount {
				streak++
				expectedDate = expectedDate.AddDate(0, 0, -1)
			}
		} else if logDate.After(expectedDate) {
			// duplicate or error, ignore
			continue
		} else {
			// Gap found
			break
		}
	}

	return streak, nil

}
