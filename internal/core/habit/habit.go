package habit

import (
	"time"

	"github.com/snehmatic/mindloop/internal/core/points"
	"github.com/snehmatic/mindloop/models"
	"gorm.io/gorm"
)

type Service struct {
	repo Repository
}

func NewService(db *gorm.DB) *Service {
	return &Service{repo: NewSQLRepository(db)}
}

func (s *Service) CreateHabit(habit *models.Habit) error {
	if habit == nil {
		return ErrHabitCannotBeNil
	}
	if err := habit.ValidateHabit(); err != nil {
		return err
	}
	return s.repo.CreateHabit(habit)
}

func (s *Service) DeleteHabit(id string) error {
	return s.repo.DeleteHabit(id)
}

func (s *Service) GetHabit(id string) (*models.Habit, error) {
	return s.repo.GetHabit(id)
}

func (s *Service) UpdateHabit(habit *models.Habit) error {
	if habit == nil {
		return ErrHabitCannotBeNil
	}
	if err := habit.ValidateHabit(); err != nil {
		return err
	}
	return s.repo.UpdateHabit(habit)
}

func (s *Service) ListHabits(interval models.IntervalType) ([]models.Habit, error) {
	return s.repo.ListHabits(interval, time.Now())
}

func (s *Service) ListEndedHabits() ([]models.Habit, error) {
	return s.repo.ListEndedHabits(time.Now())
}

func (s *Service) LogHabit(habitID string, pointsToAward int) (*models.Habit, *models.HabitLog, bool, error) {
	habit, err := s.GetHabit(habitID)
	if err != nil {
		return nil, nil, false, err
	}

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

	existingLog, err := s.repo.GetHabitLogForInterval(habit.ID, startRange, endRange)
	if err != nil {
		return nil, nil, false, err
	}

	milestoneReached := false

	if existingLog != nil {
		if existingLog.ActualCount >= habit.TargetCount {
			return habit, existingLog, false, ErrHabitAlreadyCompletedForInterval
		}

		existingLog.ActualCount++
		if err := s.repo.UpdateHabitLog(existingLog); err != nil {
			return nil, nil, false, err
		}

		if existingLog.ActualCount == habit.TargetCount {
			milestoneReached, _ = points.AwardPoints(s.repo.GetDB(), models.CategoryHabit, habit.ID, pointsToAward)
		}

		return habit, existingLog, milestoneReached, nil
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
	if err := s.repo.CreateHabitLog(habitLog); err != nil {
		return nil, nil, false, err
	}

	if habitLog.ActualCount == habit.TargetCount {
		milestoneReached, _ = points.AwardPoints(s.repo.GetDB(), models.CategoryHabit, habit.ID, pointsToAward)
	}

	return habit, habitLog, milestoneReached, nil
}

func (s *Service) UnlogHabit(habitID string) (*models.Habit, error) {
	habit, err := s.GetHabit(habitID)
	if err != nil {
		return nil, err
	}

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

	existingLog, err := s.repo.GetHabitLogForInterval(habit.ID, startRange, endRange)
	if err != nil {
		return nil, err
	}

	if existingLog == nil {
		return nil, ErrNoExistingLogFoundForThisInterval
	}

	if existingLog.ActualCount <= 0 {
		return nil, ErrHabitIsAlreadyMarkedAsUndone
	}

	existingLog.ActualCount--
	if err := s.repo.UpdateHabitLog(existingLog); err != nil {
		return nil, err
	}

	return habit, nil
}

func (s *Service) ListHabitLogs(interval models.IntervalType) ([]models.HabitLog, error) {
	return s.repo.ListHabitLogs(interval)
}

func (s *Service) ListLogsForHabit(habitID uint) ([]models.HabitLog, error) {
	return s.repo.ListLogsForHabit(habitID)
}

func (s *Service) DeleteAll() error {
	return s.repo.DeleteAll()
}

func (s *Service) RecalibrateAll() error {
	return s.repo.RecalibrateAll()
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

	allLogs, err := s.repo.ListLogsForHabits(habitIDs)
	if err != nil {
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
	logs, err := s.repo.ListLogsForHabit(habit.ID)
	if err != nil {
		return 0, err
	}
	return s.CalculateMomentumFromLogs(habit, logs), nil
}
