package habit

import (
	"time"

	"github.com/snehmatic/mindloop/models"
	"gorm.io/gorm"
)

// Repository defines the data access methods for Habits and HabitLogs
type Repository interface {
	CreateHabit(habit *models.Habit) error
	GetHabit(id string) (*models.Habit, error)
	UpdateHabit(habit *models.Habit) error
	DeleteHabit(id string) error
	ListHabits(interval models.IntervalType, now time.Time) ([]models.Habit, error)
	ListEndedHabits(now time.Time) ([]models.Habit, error)

	GetHabitLogForInterval(habitID uint, startRange, endRange time.Time) (*models.HabitLog, error)
	CreateHabitLog(log *models.HabitLog) error
	UpdateHabitLog(log *models.HabitLog) error
	ListHabitLogs(interval models.IntervalType) ([]models.HabitLog, error)
	ListLogsForHabit(habitID uint) ([]models.HabitLog, error)
	ListLogsForHabits(habitIDs []uint) ([]models.HabitLog, error)

	DeleteAll() error
	RecalibrateAll() error

	GetDB() *gorm.DB
}
