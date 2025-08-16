package ports

import (
	"time"

	"github.com/snehmatic/mindloop/internal/domain/entities"
)

// HabitRepository defines the interface for habit persistence operations
type HabitRepository interface {
	Create(habit *entities.Habit) error
	GetByID(id uint) (*entities.Habit, error)
	GetAll() ([]*entities.Habit, error)
	Update(habit *entities.Habit) error
	Delete(id uint) error
}

// HabitLogRepository defines the interface for habit log persistence operations
type HabitLogRepository interface {
	Create(habitLog *entities.HabitLog) error
	GetByID(id uint) (*entities.HabitLog, error)
	GetByHabitID(habitID uint) ([]*entities.HabitLog, error)
	GetByDateRange(start, end time.Time) ([]*entities.HabitLog, error)
	Update(habitLog *entities.HabitLog) error
	Delete(id uint) error
}

// IntentRepository defines the interface for intent persistence operations
type IntentRepository interface {
	Create(intent *entities.Intent) error
	GetByID(id uint) (*entities.Intent, error)
	GetAll() ([]*entities.Intent, error)
	GetActive() ([]*entities.Intent, error)
	GetByDateRange(start, end time.Time) ([]*entities.Intent, error)
	Update(intent *entities.Intent) error
	Delete(id uint) error
}

// FocusSessionRepository defines the interface for focus session persistence operations
type FocusSessionRepository interface {
	Create(session *entities.FocusSession) error
	GetByID(id uint) (*entities.FocusSession, error)
	GetAll() ([]*entities.FocusSession, error)
	GetActive() ([]*entities.FocusSession, error)
	GetByDateRange(start, end time.Time) ([]*entities.FocusSession, error)
	Update(session *entities.FocusSession) error
	Delete(id uint) error
}

// JournalRepository defines the interface for journal entry persistence operations
type JournalRepository interface {
	Create(entry *entities.JournalEntry) error
	GetByID(id uint) (*entities.JournalEntry, error)
	GetAll() ([]*entities.JournalEntry, error)
	GetByDateRange(start, end time.Time) ([]*entities.JournalEntry, error)
	Update(entry *entities.JournalEntry) error
	Delete(id uint) error
}

// UnitOfWork defines the interface for managing database transactions
type UnitOfWork interface {
	Begin() error
	Commit() error
	Rollback() error
	HabitRepository() HabitRepository
	HabitLogRepository() HabitLogRepository
	IntentRepository() IntentRepository
	FocusSessionRepository() FocusSessionRepository
	JournalRepository() JournalRepository
}
