package entities

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type IntervalType string

const (
	Daily  IntervalType = "daily"
	Weekly IntervalType = "weekly"
)

var AllIntervalTypes = []IntervalType{Daily, Weekly}

// Habit represents a habit entity in the domain
type Habit struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Title       string         `gorm:"type:varchar(100);not null" json:"title"`
	Description string         `gorm:"type:text" json:"description"`
	Interval    IntervalType   `gorm:"type:varchar(100);not null" json:"interval"`
	TargetCount int            `gorm:"type:int;not null" json:"target_count"`
}

// NewHabit creates a new habit with defaults
func NewHabit(title, description string, targetCount int, interval IntervalType) *Habit {
	habit := &Habit{
		Title:       title,
		Description: description,
		TargetCount: targetCount,
		Interval:    interval,
	}
	habit.SetDefaults()
	return habit
}

// SetDefaults sets default values for habit fields
func (h *Habit) SetDefaults() {
	if h.TargetCount <= 0 {
		h.TargetCount = 1
	}
	if h.Interval == "" {
		h.Interval = Daily
	}
	if h.Description == "" {
		h.Description = "Default habit description"
	}
}

// Validate validates the habit entity
func (h *Habit) Validate() error {
	if h.Title == "" {
		return fmt.Errorf("habit title cannot be empty")
	}
	if h.TargetCount <= 0 {
		return fmt.Errorf("target count must be greater than 0")
	}
	if !h.IsValidInterval() {
		return fmt.Errorf("invalid interval type: %s", h.Interval)
	}
	return nil
}

// IsValidInterval checks if the interval is valid
func (h *Habit) IsValidInterval() bool {
	return IsValidIntervalType(h.Interval)
}

// IsValidIntervalType checks if the given interval type is valid
func IsValidIntervalType(interval IntervalType) bool {
	for _, validInterval := range AllIntervalTypes {
		if validInterval == interval {
			return true
		}
	}
	return false
}

// HabitLog represents a habit log entry
type HabitLog struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	HabitID     uint           `gorm:"not null" json:"habit_id"`
	Title       string         `gorm:"not null" json:"title"`
	Interval    IntervalType   `gorm:"type:varchar(100);not null" json:"interval"`
	TargetCount int            `gorm:"not null" json:"target_count"`
	ActualCount int            `gorm:"not null" json:"actual_count"`
	EndedAt     time.Time      `gorm:"not null" json:"ended_at"`
}

// NewHabitLog creates a new habit log
func NewHabitLog(habitID uint, title string, interval IntervalType, targetCount, actualCount int) *HabitLog {
	return &HabitLog{
		HabitID:     habitID,
		Title:       title,
		Interval:    interval,
		TargetCount: targetCount,
		ActualCount: actualCount,
		EndedAt:     time.Now(),
	}
}

// IsCompleted checks if the habit log is completed
func (hl *HabitLog) IsCompleted() bool {
	return hl.ActualCount >= hl.TargetCount
}
