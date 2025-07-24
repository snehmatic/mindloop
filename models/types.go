package models

import (
	"math"
	"time"

	"github.com/snehmatic/mindloop/internal/config"
	"gorm.io/gorm"
)

// model definitions reside here
// request/response structs, etc.

type IntervalType string

var AllIntervalTypes = [...]string{"hour", "day", "week"}

var (
	Hour IntervalType = IntervalType(AllIntervalTypes[0])
	Day  IntervalType = IntervalType(AllIntervalTypes[1])
	Week IntervalType = IntervalType(AllIntervalTypes[2])
)

type Habit struct {
	gorm.Model
	Name        string       `gorm:"type:varchar(100)" json:"name"`
	Interval    IntervalType `gorm:"type:varchar(100)" json:"internal"`
	ActualCount int          `gorm:"type:int" json:"actual_count"`
	TargetCount int          `gorm:"type:int" json:"target_count"`
}

type Intent struct {
	gorm.Model
	Name    string     `gorm:"not null" json:"name"`         // comma-separated or JSON later
	Status  string     `gorm:"default:active" json:"status"` // active, done, archived
	EndedAt *time.Time `json:"ended_at,omitempty"`
}

type IntentView struct {
	ID      uint
	Name    string
	Status  string
	EndedAt string
}

func ToIntentView(i Intent) IntentView {
	var ended string
	if i.EndedAt != nil {
		ended = i.EndedAt.Format("2006-01-02 15:04")
	} else {
		ended = "-"
	}
	return IntentView{
		ID:      i.ID,
		Name:    i.Name,
		Status:  i.Status,
		EndedAt: ended,
	}
}

type FocusSession struct {
	gorm.Model
	Title    string    `gorm:"not null" json:"title"`        // e.g., "Work on project"
	Status   string    `gorm:"default:active" json:"status"` // active, paused
	EndTime  time.Time `json:"end_time"`
	Duration float64   `json:"duration"`                 // in seconds
	Rating   int       `gorm:"default:-1" json:"rating"` // 1 to 5, optional
}

type FocusSessionView struct {
	ID        uint    `json:"id"`
	Title     string  `json:"title"`
	Status    string  `json:"status"`
	EndTime   string  `json:"end_time"`   // formatted as "2006-01-02 15:04:05"
	Duration  float64 `json:"duration"`   // in seconds
	Rating    int     `json:"rating"`     // 1 to 10, -1 if not rated
	CreatedAt string  `json:"created_at"` // formatted as "2006-01-02 15:04:05"
}

func ToFocusSessionView(fs FocusSession) FocusSessionView {
	fsv := FocusSessionView{
		ID:        fs.ID,
		Title:     fs.Title,
		Status:    fs.Status,
		EndTime:   fs.EndTime.Format("2006-01-02 15:04:05"),
		Duration:  fs.Duration,
		Rating:    fs.Rating,
		CreatedAt: fs.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	if fs.EndTime.IsZero() {
		fsv.EndTime = "Focus on!"
	}
	if fs.Rating == 0 {
		fsv.Rating = -1 // indicate no rating given
	}
	now := time.Now()
	fsv.Duration = now.Sub(fs.CreatedAt).Minutes()
	fsv.Duration = math.Floor(fsv.Duration) // todo: fix decimals
	return fsv
}

type HabitLog struct {
	gorm.Model
	HabitID   int       `gorm:"not null;index:idx_habit_day,unique" json:"habit_id"`
	Count     int       `gorm:"not null" json:"count"` // number of times the habit was done
	Completed bool      `gorm:"default:false" json:"completed"`
	HabitName string    `gorm:"not null;index:idx_habit_day,unique" json:"habit_name"`
	Date      time.Time `gorm:"not null;index:idx_habit_day,unique" json:"date"` // YYYY-MM-DD only
}

type JournalEntry struct {
	gorm.Model
	Date  time.Time `gorm:"uniqueIndex" json:"date"` // one entry per day
	Entry string    `gorm:"type:text" json:"entry"`
	Mood  string    `gorm:"type:varchar(50)" json:"mood"` // e.g., happy, sad, neutral
}

func IsValidMode(mode string) bool {
	for _, item := range config.AllModes {
		if item == mode {
			return true
		}
	}
	return false
}
