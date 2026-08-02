package habit

import (
	"errors"
	"time"

	"github.com/snehmatic/mindloop/internal/core/points"
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
	var habits []models.Habit
	if err := s.DB.Where("id = ?", id).Limit(1).Find(&habits).Error; err != nil {
		return err
	}
	if len(habits) == 0 {
		return errors.New("habit not found")
	}
	return s.DB.Delete(&habits[0]).Error
}

func (s *Service) GetHabit(id string) (*models.Habit, error) {
	var habits []models.Habit
	if err := s.DB.Where("id = ?", id).Limit(1).Find(&habits).Error; err != nil {
		return nil, err
	}
	if len(habits) == 0 {
		return nil, errors.New("habit not found")
	}
	return &habits[0], nil
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

	result := query.Order("CreatedAt DESC").Find(&habits)
	return habits, result.Error
}

func (s *Service) ListEndedHabits() ([]models.Habit, error) {
	var habits []models.Habit
	now := time.Now()
	result := s.DB.Where("EndDate IS NOT NULL AND EndDate <= ?", now).Order("CreatedAt DESC").Find(&habits)
	return habits, result.Error
}

func (s *Service) LogHabit(habitID string, pointsToAward int) (*models.Habit, *models.HabitLog, bool, error) {
	habit, err := s.GetHabit(habitID)
	if err != nil {
		return nil, nil, false, err
	}

	var existingLogs []models.HabitLog
	now := time.Now()
	today := now.Truncate(24 * time.Hour)
	var startRange, endRange time.Time

	switch habit.Interval {
	case models.Daily:
		startRange = today
		endRange = today.AddDate(0, 0, 1)
	case models.Weekly:
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7 // Sunday
		}
		startRange = today.AddDate(0, 0, -(weekday - 1))
		endRange = startRange.AddDate(0, 0, 7)
	}

	res := s.DB.Where("HabitID = ? AND CreatedAt >= ? AND CreatedAt < ?", habit.ID, startRange, endRange).Limit(1).Find(&existingLogs)
	if res.Error != nil {
		return nil, nil, false, res.Error
	}

	milestoneReached := false

	if len(existingLogs) > 0 {
		existingLog := existingLogs[0]
		if existingLog.ActualCount >= habit.TargetCount {
			return habit, &existingLog, false, errors.New("habit already completed for interval")
		}

		existingLog.ActualCount++
		if err := s.DB.Save(&existingLog).Error; err != nil {
			return nil, nil, false, err
		}

		if existingLog.ActualCount == habit.TargetCount {
			milestoneReached, _ = points.AwardPoints(s.DB, models.CategoryHabit, habit.ID, pointsToAward)
		}

		return habit, &existingLog, milestoneReached, nil
	}

	// Create new log
	habitLog := &models.HabitLog{
		HabitID:     habit.ID,
		Title:       habit.Title,
		Interval:    habit.Interval,
		TargetCount: habit.TargetCount,
		ActualCount: 1,
		EndedAt:     endRange.AddDate(0, 0, -1), // Represents the last day of the interval
	}
	if err := s.DB.Create(habitLog).Error; err != nil {
		return nil, nil, false, err
	}

	if habitLog.ActualCount == habit.TargetCount {
		milestoneReached, _ = points.AwardPoints(s.DB, models.CategoryHabit, habit.ID, pointsToAward)
	}

	return habit, habitLog, milestoneReached, nil
}

func (s *Service) UnlogHabit(habitID string) (*models.Habit, error) {
	habit, err := s.GetHabit(habitID)
	if err != nil {
		return nil, err
	}

	var existingLogs []models.HabitLog
	now := time.Now()
	today := now.Truncate(24 * time.Hour)
	var startRange, endRange time.Time

	switch habit.Interval {
	case models.Daily:
		startRange = today
		endRange = today.AddDate(0, 0, 1)
	case models.Weekly:
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		startRange = today.AddDate(0, 0, -(weekday - 1))
		endRange = startRange.AddDate(0, 0, 7)
	}

	res := s.DB.Where("HabitID = ? AND CreatedAt >= ? AND CreatedAt < ?", habit.ID, startRange, endRange).Limit(1).Find(&existingLogs)
	if res.Error != nil {
		return nil, res.Error
	}

	if len(existingLogs) == 0 {
		return nil, errors.New("no existing log found for this interval")
	}

	existingLog := existingLogs[0]
	if existingLog.ActualCount <= 0 {
		return nil, errors.New("habit is already marked as undone")
	}

	existingLog.ActualCount--
	if err := s.DB.Save(&existingLog).Error; err != nil {
		return nil, err
	}

	return habit, nil
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

func (s *Service) RecalibrateAll() error {
	return s.DB.Session(&gorm.Session{AllowGlobalUpdate: true}).Model(&models.Habit{}).Update("RecalibratedAt", time.Now()).Error
}

func (s *Service) CalculateMomentumFromLogs(habit *models.Habit, logs []models.HabitLog) int {
	if habit.Interval != models.Daily {
		return 0
	}

	momentum := float64(0)
	startDate := habit.CreatedAt.Truncate(24 * time.Hour)
	if len(logs) > 0 {
		firstLogDate := logs[0].CreatedAt.Truncate(24 * time.Hour)
		if firstLogDate.Before(startDate) {
			startDate = firstLogDate
		}
	}
	today := time.Now().Truncate(24 * time.Hour)

	logMap := make(map[string]bool)
	for _, l := range logs {
		if l.ActualCount >= l.TargetCount {
			logDate := l.CreatedAt.Format("2006-01-02")
			logMap[logDate] = true
		}
	}

	for d := startDate; !d.After(today); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format("2006-01-02")
		isForgiven := false
		if habit.RecalibratedAt != nil {
			if d.Format("2006-01-02") < habit.RecalibratedAt.Format("2006-01-02") {
				isForgiven = true
			}
		}

		if logMap[dateStr] {
			momentum += 10
		} else if !isForgiven {
			momentum *= 0.9
		}
	}

	if momentum > 100 {
		momentum = 100
	}
	return int(momentum)
}

func (s *Service) CalculateMomentums(habits []models.Habit) (map[uint]int, error) {
	momentums := make(map[uint]int)
	if len(habits) == 0 {
		return momentums, nil
	}

	var habitIDs []uint
	for _, h := range habits {
		habitIDs = append(habitIDs, h.ID)
	}

	var allLogs []models.HabitLog
	if err := s.DB.Where("HabitID IN ?", habitIDs).Order("CreatedAt asc").Find(&allLogs).Error; err != nil {
		return nil, err
	}

	logsByHabit := make(map[uint][]models.HabitLog)
	for _, log := range allLogs {
		logsByHabit[log.HabitID] = append(logsByHabit[log.HabitID], log)
	}

	for _, h := range habits {
		momentums[h.ID] = s.CalculateMomentumFromLogs(&h, logsByHabit[h.ID])
	}

	return momentums, nil
}

func (s *Service) CalculateMomentum(habit *models.Habit) (int, error) {
	var logs []models.HabitLog
	if err := s.DB.Where("HabitID = ?", habit.ID).Order("CreatedAt asc").Find(&logs).Error; err != nil {
		return 0, err
	}
	return s.CalculateMomentumFromLogs(habit, logs), nil
}
