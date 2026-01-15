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
	result := query.Find(&habits)
	return habits, result.Error
}

func (s *Service) LogHabit(habitID string) (*models.Habit, *models.HabitLog, error) {
	var habit models.Habit
	if err := s.DB.First(&habit, "id = ?", habitID).Error; err != nil {
		return nil, nil, err
	}

	var existingLog models.HabitLog
	today := time.Now().Truncate(24 * time.Hour)
	endedAt := today

	var res *gorm.DB
	switch habit.Interval {
	case models.Daily:
		res = s.DB.Where("HabitID = ? AND EndedAt = ?", habit.ID, today).First(&existingLog)
	case models.Weekly:
		startOfWeek := time.Now().AddDate(0, 0, -int(time.Now().Weekday()))
		endOfWeek := startOfWeek.AddDate(0, 0, 6).Truncate(24 * time.Hour)
		endedAt = endOfWeek
		res = s.DB.Where("HabitID = ? AND CreatedAt >= ? AND EndedAt <= ?", habit.ID, startOfWeek, endOfWeek).First(&existingLog)
	}

	if res.Error == nil {
		// Log found
		if existingLog.ActualCount >= habit.TargetCount {
			return &habit, &existingLog, errors.New("habit already completed for interval")
		}

		existingLog.ActualCount++
		existingLog.EndedAt = endedAt
		if err := s.DB.Save(&existingLog).Error; err != nil {
			return nil, nil, err
		}
		return &habit, &existingLog, nil
	} else if !errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, nil, res.Error
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
	today := time.Now().Format("2006-01-02")
	var res *gorm.DB

	switch habit.Interval {
	case models.Daily:
		res = s.DB.Where("HabitID = ? AND DATE(CreatedAt) = ?", habit.ID, today).First(&existingLog)
	case models.Weekly:
		startOfWeek := time.Now().AddDate(0, 0, -int(time.Now().Weekday()))
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
