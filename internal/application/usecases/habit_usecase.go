package usecases

import (
	"fmt"
	"time"

	"github.com/snehmatic/mindloop/internal/domain/entities"
	"github.com/snehmatic/mindloop/internal/domain/ports"
)

// HabitUseCase defines the interface for habit use cases
type HabitUseCase interface {
	CreateHabit(title, description string, targetCount int, interval entities.IntervalType) (*entities.Habit, error)
	GetHabit(id uint) (*entities.Habit, error)
	GetAllHabits() ([]*entities.Habit, error)
	UpdateHabit(habit *entities.Habit) error
	DeleteHabit(id uint) error
	LogHabit(habitID uint, actualCount int) error
	UnlogHabit(logID uint) error
	GetHabitLogs(habitID uint) ([]*entities.HabitLog, error)
	GetAllHabitLogs() ([]*entities.HabitLog, error)
}

type habitUseCase struct {
	habitRepo    ports.HabitRepository
	habitLogRepo ports.HabitLogRepository
}

// NewHabitUseCase creates a new habit use case
func NewHabitUseCase(habitRepo ports.HabitRepository, habitLogRepo ports.HabitLogRepository) HabitUseCase {
	return &habitUseCase{
		habitRepo:    habitRepo,
		habitLogRepo: habitLogRepo,
	}
}

// CreateHabit creates a new habit
func (h *habitUseCase) CreateHabit(title, description string, targetCount int, interval entities.IntervalType) (*entities.Habit, error) {
	habit := entities.NewHabit(title, description, targetCount, interval)

	if err := habit.Validate(); err != nil {
		return nil, fmt.Errorf("habit validation failed: %w", err)
	}

	if err := h.habitRepo.Create(habit); err != nil {
		return nil, fmt.Errorf("failed to create habit: %w", err)
	}

	return habit, nil
}

// GetHabit retrieves a habit by ID
func (h *habitUseCase) GetHabit(id uint) (*entities.Habit, error) {
	habit, err := h.habitRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get habit: %w", err)
	}
	return habit, nil
}

// GetAllHabits retrieves all habits
func (h *habitUseCase) GetAllHabits() ([]*entities.Habit, error) {
	habits, err := h.habitRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get habits: %w", err)
	}
	return habits, nil
}

// UpdateHabit updates an existing habit
func (h *habitUseCase) UpdateHabit(habit *entities.Habit) error {
	if err := habit.Validate(); err != nil {
		return fmt.Errorf("habit validation failed: %w", err)
	}

	if err := h.habitRepo.Update(habit); err != nil {
		return fmt.Errorf("failed to update habit: %w", err)
	}

	return nil
}

// DeleteHabit deletes a habit
func (h *habitUseCase) DeleteHabit(id uint) error {
	// Check if habit exists
	_, err := h.habitRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("habit not found: %w", err)
	}

	if err := h.habitRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete habit: %w", err)
	}

	return nil
}

// LogHabit creates a habit log entry
func (h *habitUseCase) LogHabit(habitID uint, actualCount int) error {
	// Get the habit to validate it exists and get its properties
	habit, err := h.habitRepo.GetByID(habitID)
	if err != nil {
		return fmt.Errorf("habit not found: %w", err)
	}

	if actualCount < 0 {
		return fmt.Errorf("actual count cannot be negative")
	}

	habitLog := entities.NewHabitLog(habitID, habit.Title, habit.Interval, habit.TargetCount, actualCount)

	if err := h.habitLogRepo.Create(habitLog); err != nil {
		return fmt.Errorf("failed to log habit: %w", err)
	}

	return nil
}

// UnlogHabit removes a habit log entry
func (h *habitUseCase) UnlogHabit(logID uint) error {
	// Check if log exists
	_, err := h.habitLogRepo.GetByID(logID)
	if err != nil {
		return fmt.Errorf("habit log not found: %w", err)
	}

	if err := h.habitLogRepo.Delete(logID); err != nil {
		return fmt.Errorf("failed to unlog habit: %w", err)
	}

	return nil
}

// GetHabitLogs retrieves logs for a specific habit
func (h *habitUseCase) GetHabitLogs(habitID uint) ([]*entities.HabitLog, error) {
	logs, err := h.habitLogRepo.GetByHabitID(habitID)
	if err != nil {
		return nil, fmt.Errorf("failed to get habit logs: %w", err)
	}
	return logs, nil
}

// GetAllHabitLogs retrieves all habit logs
func (h *habitUseCase) GetAllHabitLogs() ([]*entities.HabitLog, error) {
	// Get logs from last 30 days by default
	end := time.Now()
	start := end.AddDate(0, 0, -30)

	logs, err := h.habitLogRepo.GetByDateRange(start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get habit logs: %w", err)
	}
	return logs, nil
}
