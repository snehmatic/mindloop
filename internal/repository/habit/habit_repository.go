package habit

import (
	"time"

	"github.com/snehmatic/mindloop/models"
	"gorm.io/gorm"
)

// Repository defines the interface for habit data access
type Repository interface {
	CreateHabit(habit *models.Habit) error
	DeleteHabit(id string) error
	GetHabit(id string) (*models.Habit, error)
	UpdateHabit(habit *models.Habit) error
	ListHabits(interval models.IntervalType) ([]models.Habit, error)
	ListEndedHabits() ([]models.Habit, error)
	LogHabit(habitID string, pointsToAward int) (*models.Habit, *models.HabitLog, bool, error)
	UnlogHabit(habitID string) (*models.Habit, error)
	ListHabitLogs(interval models.IntervalType) ([]models.HabitLog, error)
	ListLogsForHabit(habitID uint) ([]models.HabitLog, error)
	DeleteAll() error
	CalculateStreak(habitID uint, interval models.IntervalType) (int, error)
	GetHabitStats(start, end time.Time) ([]models.HabitStats, error)
	GetAllHabits() ([]models.Habit, error)
	GetHabitLogsInRange(start, end time.Time) ([]models.HabitLog, error)
	GetDB() *gorm.DB
}
