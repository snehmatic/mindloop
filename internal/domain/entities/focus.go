package entities

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type FocusStatus string

const (
	FocusActive FocusStatus = "active"
	FocusPaused FocusStatus = "paused"
	FocusEnded  FocusStatus = "ended"
)

// FocusSession represents a focus session entity in the domain
type FocusSession struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Title     string         `gorm:"not null" json:"title"`
	Status    FocusStatus    `gorm:"default:active" json:"status"`
	EndTime   time.Time      `json:"end_time"`
	Duration  float64        `json:"duration"`                 // in minutes
	Rating    int            `gorm:"default:-1" json:"rating"` // 0 to 10, -1 if not rated
}

// NewFocusSession creates a new focus session
func NewFocusSession(title string) *FocusSession {
	return &FocusSession{
		Title:  title,
		Status: FocusActive,
		Rating: -1, // not rated
	}
}

// Validate validates the focus session entity
func (fs *FocusSession) Validate() error {
	if fs.Title == "" {
		return fmt.Errorf("focus session title cannot be empty")
	}
	return nil
}

// End ends the focus session
func (fs *FocusSession) End() {
	if fs.Status != FocusActive {
		return
	}
	fs.Status = FocusEnded
	fs.EndTime = time.Now()
	fs.Duration = fs.EndTime.Sub(fs.CreatedAt).Minutes()
}

// Pause pauses the focus session
func (fs *FocusSession) Pause() {
	if fs.Status == FocusActive {
		fs.Status = FocusPaused
	}
}

// Resume resumes the focus session
func (fs *FocusSession) Resume() {
	if fs.Status == FocusPaused {
		fs.Status = FocusActive
	}
}

// Rate sets the rating for the focus session
func (fs *FocusSession) Rate(rating int) error {
	if rating < 0 || rating > 10 {
		return fmt.Errorf("rating must be between 0 and 10")
	}
	if fs.Status != FocusEnded {
		return fmt.Errorf("can only rate ended focus sessions")
	}
	fs.Rating = rating
	return nil
}

// IsActive checks if the focus session is active
func (fs *FocusSession) IsActive() bool {
	return fs.Status == FocusActive
}

// IsEnded checks if the focus session is ended
func (fs *FocusSession) IsEnded() bool {
	return fs.Status == FocusEnded
}

// GetCurrentDuration returns the current duration of the session in minutes
func (fs *FocusSession) GetCurrentDuration() float64 {
	if fs.IsEnded() {
		return fs.Duration
	}
	return time.Since(fs.CreatedAt).Minutes()
}
