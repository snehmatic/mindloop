package entities

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Mood string

const (
	MoodHappy   Mood = "happy"
	MoodSad     Mood = "sad"
	MoodNeutral Mood = "neutral"
	MoodAngry   Mood = "angry"
	MoodExcited Mood = "excited"
)

var AllMoods = []Mood{MoodHappy, MoodSad, MoodNeutral, MoodAngry, MoodExcited}

// JournalEntry represents a journal entry entity in the domain
type JournalEntry struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Content   string         `gorm:"type:text;not null" json:"content"`
	Title     string         `gorm:"type:varchar(100);not null" json:"title"`
	Mood      Mood           `gorm:"type:varchar(50);not null" json:"mood"`
}

// NewJournalEntry creates a new journal entry
func NewJournalEntry(title, content string, mood Mood) *JournalEntry {
	entry := &JournalEntry{
		Title:   title,
		Content: content,
		Mood:    mood,
	}
	entry.SetDefaults()
	return entry
}

// SetDefaults sets default values for journal entry fields
func (je *JournalEntry) SetDefaults() {
	if je.Mood == "" {
		je.Mood = MoodNeutral
	}
}

// Validate validates the journal entry entity
func (je *JournalEntry) Validate() error {
	if je.Title == "" {
		return fmt.Errorf("journal entry title cannot be empty")
	}
	if je.Content == "" {
		return fmt.Errorf("journal entry content cannot be empty")
	}
	if !je.IsValidMood() {
		return fmt.Errorf("invalid mood: %s", je.Mood)
	}
	return nil
}

// IsValidMood checks if the mood is valid
func (je *JournalEntry) IsValidMood() bool {
	return IsValidMood(je.Mood)
}

// IsValidMood checks if the given mood is valid
func IsValidMood(mood Mood) bool {
	for _, validMood := range AllMoods {
		if validMood == mood {
			return true
		}
	}
	return false
}

// UpdateContent updates the content of the journal entry
func (je *JournalEntry) UpdateContent(content string) error {
	if content == "" {
		return fmt.Errorf("journal entry content cannot be empty")
	}
	je.Content = content
	return nil
}

// UpdateMood updates the mood of the journal entry
func (je *JournalEntry) UpdateMood(mood Mood) error {
	if !IsValidMood(mood) {
		return fmt.Errorf("invalid mood: %s", mood)
	}
	je.Mood = mood
	return nil
}
