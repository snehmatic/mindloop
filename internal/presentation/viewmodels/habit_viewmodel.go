package viewmodels

import (
	"github.com/snehmatic/mindloop/internal/domain/entities"
)

// HabitView represents a habit for display purposes
type HabitView struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Interval    string `json:"interval"`
	TargetCount int    `json:"target_count"`
	CreatedAt   string `json:"created_at"`
}

// ToHabitView converts a domain habit to view model
func ToHabitView(habit *entities.Habit) HabitView {
	return HabitView{
		ID:          habit.ID,
		Title:       habit.Title,
		Description: habit.Description,
		Interval:    string(habit.Interval),
		TargetCount: habit.TargetCount,
		CreatedAt:   habit.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

// ToHabitViews converts multiple domain habits to view models
func ToHabitViews(habits []*entities.Habit) []HabitView {
	views := make([]HabitView, len(habits))
	for i, habit := range habits {
		views[i] = ToHabitView(habit)
	}
	return views
}

// HabitLogView represents a habit log for display purposes
type HabitLogView struct {
	ID          uint   `json:"id"`
	HabitID     uint   `json:"habit_id"`
	Title       string `json:"title"`
	TargetCount int    `json:"target_count"`
	ActualCount int    `json:"actual_count"`
	Interval    string `json:"interval"`
	StartedAt   string `json:"started_at"`
	EndedAt     string `json:"ended_at"`
	Completed   string `json:"completed"`
}

// ToHabitLogView converts a domain habit log to view model
func ToHabitLogView(log *entities.HabitLog) HabitLogView {
	completed := "No"
	if log.IsCompleted() {
		completed = "Yes"
	}

	return HabitLogView{
		ID:          log.ID,
		HabitID:     log.HabitID,
		Title:       log.Title,
		TargetCount: log.TargetCount,
		ActualCount: log.ActualCount,
		Interval:    string(log.Interval),
		StartedAt:   log.CreatedAt.Format("2006-01-02"),
		EndedAt:     log.EndedAt.Format("2006-01-02"),
		Completed:   completed,
	}
}

// ToHabitLogViews converts multiple domain habit logs to view models
func ToHabitLogViews(logs []*entities.HabitLog) []HabitLogView {
	views := make([]HabitLogView, len(logs))
	for i, log := range logs {
		views[i] = ToHabitLogView(log)
	}
	return views
}
