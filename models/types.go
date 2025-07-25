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

var AllIntervalTypes = [...]string{"daily", "weekly"}

var (
	Daily  IntervalType = IntervalType(AllIntervalTypes[0])
	Weekly IntervalType = IntervalType(AllIntervalTypes[1])
)

type Habit struct {
	gorm.Model
	Title       string       `gorm:"type:varchar(100)" json:"title"`
	Description string       `gorm:"type:text" json:"description"`
	Interval    IntervalType `gorm:"type:varchar(100)" json:"interval"`
	TargetCount int          `gorm:"type:int" json:"target_count"`
}

type HabitLog struct {
	gorm.Model
	HabitID         int       `gorm:"not null;index:idx_habit_day" json:"habit_id"`
	TargetCount     int       `gorm:"not null" json:"target_count"`      // number of times the habit was done
	ActualCount     int       `gorm:"not null" json:"actual_count"`      // number of times the habit was actually done
	CompletedOnDate time.Time `gorm:"not null" json:"completed_on_date"` // YYYY-MM-DD only

}

type HabitLogView struct {
	ID          uint         `json:"id"`
	Title       string       `json:"title"`
	TargetCount int          `json:"target_count"`
	ActualCount int          `json:"actual_count"`
	Interval    IntervalType `json:"interval"`
	CompletedOn string       `json:"completed_on"`
}

func ToHabitLogViews(habitLogs []HabitLog) []HabitLogView {
	habitViews := make([]HabitLogView, len(habitLogs))
	for i, log := range habitLogs {
		habitViews[i] = HabitLogView{
			ID:          log.ID,
			ActualCount: log.ActualCount,
			TargetCount: log.TargetCount,
			CompletedOn: log.CompletedOnDate.Format("2006-01-02"),
			Interval:    Daily, // default to daily, can be changed later
			Title:       "",    // title can be fetched from Habit model if needed
		}
	}
	return habitViews
}

type HabitView struct {
	ID          uint         `json:"id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Interval    IntervalType `json:"interval"`
	ActualCount int          `json:"actual_count"`
	TargetCount int          `json:"target_count"`
}

func ToHabitView(h Habit) HabitView {
	return HabitView{
		ID:          h.ID,
		Title:       h.Title,
		Description: h.Description,
		Interval:    h.Interval,
		TargetCount: h.TargetCount,
	}
}

func IsValidIntervalType(interval string) bool {
	for _, item := range AllIntervalTypes {
		if item == interval {
			return true
		}
	}
	return false
}

type Intent struct {
	gorm.Model
	Name    string     `gorm:"not null" json:"name"`
	Status  string     `gorm:"default:active" json:"status"`
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
	Duration float64   `json:"duration"`                 // in mins
	Rating   int       `gorm:"default:-1" json:"rating"` // 0 to 10, optional
}

type FocusSessionView struct {
	ID        uint    `json:"id"`
	Title     string  `json:"title"`
	Status    string  `json:"status"`
	EndTime   string  `json:"end_time"`   // formatted as "2006-01-02 15:04:05"
	Duration  float64 `json:"duration"`   // in mins
	Rating    int     `json:"rating"`     // 0 to 10, -1 if not rated
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
