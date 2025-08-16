package entities

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type IntentStatus string

const (
	IntentActive IntentStatus = "active"
	IntentDone   IntentStatus = "done"
)

// Intent represents an intent entity in the domain
type Intent struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Name      string         `gorm:"not null" json:"name"`
	Status    IntentStatus   `gorm:"default:active" json:"status"`
	EndedAt   *time.Time     `json:"ended_at,omitempty"`
}

// NewIntent creates a new intent
func NewIntent(name string) *Intent {
	return &Intent{
		Name:   name,
		Status: IntentActive,
	}
}

// Validate validates the intent entity
func (i *Intent) Validate() error {
	if i.Name == "" {
		return fmt.Errorf("intent name cannot be empty")
	}
	return nil
}

// End marks the intent as done
func (i *Intent) End() {
	now := time.Now()
	i.Status = IntentDone
	i.EndedAt = &now
}

// IsActive checks if the intent is active
func (i *Intent) IsActive() bool {
	return i.Status == IntentActive
}

// IsDone checks if the intent is done
func (i *Intent) IsDone() bool {
	return i.Status == IntentDone
}
