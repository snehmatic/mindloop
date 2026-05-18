package habit

import (
	"time"

	"github.com/snehmatic/mindloop/internal/repository/habit"
	"github.com/snehmatic/mindloop/models"
)

type Service struct {
	repo habit.Repository
}

func NewService(repo habit.Repository) *Service {
	return &Service{repo: repo}
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
	return s.repo.ListHabits(interval)
}

func (s *Service) ListEndedHabits() ([]models.Habit, error) {
	return s.repo.ListEndedHabits()
}

func (s *Service) LogHabit(habitID string, pointsToAward int) (*models.Habit, *models.HabitLog, bool, error) {
	return s.repo.LogHabit(habitID, pointsToAward)
}

func (s *Service) UnlogHabit(habitID string) (*models.Habit, error) {
	return s.repo.UnlogHabit(habitID)
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

func (s *Service) CalculateStreak(habitID uint, interval models.IntervalType) (int, error) {
	return s.repo.CalculateStreak(habitID, interval)
}

func (s *Service) GetHabitStats(start, end time.Time) ([]models.HabitStats, error) {
	return s.repo.GetHabitStats(start, end)
}

func (s *Service) GetAllHabits() ([]models.Habit, error) {
	return s.repo.GetAllHabits()
}

func (s *Service) GetHabitLogsInRange(start, end time.Time) ([]models.HabitLog, error) {
	return s.repo.GetHabitLogsInRange(start, end)
}
