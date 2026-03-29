// Package models contains the data structures and UI views for Mindloop
package models

import (
	"fmt"
	"math"
	"time"

	"github.com/snehmatic/mindloop/internal/config"
	"gorm.io/gorm"
)

// model definitions reside here
// request/response structs, etc.

// IntervalType defines the frequency of a habit
type IntervalType string

var AllIntervalTypes = [...]string{"daily", "weekly"}

var (
	Daily  = IntervalType(AllIntervalTypes[0])
	Weekly = IntervalType(AllIntervalTypes[1])
)

// Habit represents a task that the user wants to perform regularly
type Habit struct {
	gorm.Model
	Title       string       `gorm:"type:varchar(100)" json:"title"`
	Description string       `gorm:"type:text" json:"description"`
	Interval    IntervalType `gorm:"type:varchar(100)" json:"interval"`
	TargetCount int          `gorm:"type:int" json:"target_count"`
	EndDate     *time.Time   `json:"end_date,omitempty"`
	RoutineID   *uint        `json:"routine_id,omitempty"`
}

// Routine groups multiple habits into a specific time of day
type Routine struct {
	gorm.Model
	Title     string  `gorm:"type:varchar(100)" json:"title"`
	TimeOfDay string  `gorm:"type:varchar(50)" json:"time_of_day"`
	Habits    []Habit `gorm:"foreignKey:RoutineID" json:"habits"`
}

// RoutineView is a simplified representation of a Routine for the UI
type RoutineView struct {
	ID        uint        `json:"id"`
	Title     string      `json:"title"`
	TimeOfDay string      `json:"time_of_day"`
	Habits    []HabitView `json:"habits"`
}

// Defaults for Habit
// TargetCount: 1
// Interval: Daily
// Description: "Default habit description"
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

func (h *Habit) ValidateHabit() error {
	if h.Title == "" {
		return fmt.Errorf("habit title cannot be empty")
	}
	if h.TargetCount <= 0 {
		return fmt.Errorf("target count must be greater than 0")
	}
	if !IsValidIntervalType(string(h.Interval)) {
		return fmt.Errorf("invalid interval type: %s", h.Interval)
	}
	return nil
}

// HabitLog records an instance of a habit being performed
type HabitLog struct {
	gorm.Model
	HabitID     uint         `gorm:"not null" json:"habit_id"`
	Title       string       `gorm:"not null" json:"title"`
	Interval    IntervalType `gorm:"type:varchar(100);not null" json:"interval"`
	TargetCount int          `gorm:"not null" json:"target_count"`
	ActualCount int          `gorm:"not null" json:"actual_count"`
	EndedAt     time.Time    `gorm:"not null" json:"ended_at"`
}

// HabitLogView is a simplified representation of a HabitLog for the UI
type HabitLogView struct {
	ID          uint         `json:"id"`
	HabitID     uint         `json:"habit_id"`
	Title       string       `json:"title"`
	TargetCount int          `json:"target_count"`
	ActualCount int          `json:"actual_count"`
	Interval    IntervalType `json:"interval"`
	StartedAt   string       `json:"started_at"`
	EndedAt     string       `json:"ended_at"`
}

func ToHabitLogViews(habitLogs []HabitLog) []HabitLogView {
	habitViews := make([]HabitLogView, len(habitLogs))
	for i, log := range habitLogs {
		habitViews[i] = HabitLogView{
			ID:          log.ID,
			HabitID:     log.HabitID,
			ActualCount: log.ActualCount,
			TargetCount: log.TargetCount,
			StartedAt:   log.CreatedAt.Format("2006-01-02"),
			EndedAt:     log.EndedAt.Format("2006-01-02"),
			Interval:    log.Interval,
			Title:       log.Title,
		}
	}
	return habitViews
}

// HabitView is a simplified representation of a Habit for the UI
type HabitView struct {
	ID          uint         `json:"id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Interval    IntervalType `json:"interval"`
	TargetCount int          `json:"target_count"`
	EndDate     string       `json:"end_date"`
	IsEnded     bool         `json:"is_ended"`
}

func ToHabitView(h Habit) HabitView {
	ended := ""
	isEnded := false
	if h.EndDate != nil {
		ended = h.EndDate.Format("2006-01-02")
		if h.EndDate.Before(time.Now()) {
			isEnded = true
		}
	}
	return HabitView{
		ID:          h.ID,
		Title:       h.Title,
		Description: h.Description,
		Interval:    h.Interval,
		TargetCount: h.TargetCount,
		EndDate:     ended,
		IsEnded:     isEnded,
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

// Intent represents a high-level goal for the user
type Intent struct {
	gorm.Model
	Name    string     `gorm:"not null" json:"name"`
	Status  string     `gorm:"default:active" json:"status"`
	EndedAt *time.Time `json:"ended_at,omitempty"`
}

// IntentView is a simplified representation of an Intent for the UI
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

// FocusSession records a period of deep work
type FocusSession struct {
	gorm.Model
	Title    string    `gorm:"not null" json:"title"`        // e.g., "Work on project"
	Status   string    `gorm:"default:active" json:"status"` // active, paused
	EndTime  time.Time `json:"end_time"`
	Duration float64   `json:"duration"`                 // in mins
	Rating   int       `gorm:"default:-1" json:"rating"` // 0 to 10, optional
}

// FocusSessionView is a simplified representation of a FocusSession for the UI
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

// JournalEntry represents a user's reflective writing
type JournalEntry struct {
	gorm.Model
	Content string `gorm:"type:text" json:"content"`
	Title   string `gorm:"type:varchar(100)" json:"title"`
	Mood    string `gorm:"type:varchar(50)" json:"mood"` // e.g., happy, sad, neutral
}

// JournalEntryView is a simplified representation of a JournalEntry for the UI
type JournalEntryView struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
	Mood  string `json:"mood"`
	Date  string `json:"date"` // formatted as "2006-01-02 15:04:05"
}

func ToJournalEntryView(entry JournalEntry) JournalEntryView {
	return JournalEntryView{
		ID:    entry.ID,
		Title: entry.Title,
		Mood:  entry.Mood,
		Date:  entry.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

// Note represents a quick markdown note
type Note struct {
	gorm.Model
	Title   string `gorm:"type:varchar(200)" json:"title"`
	Content string `gorm:"type:text" json:"content"`
	Labels  string `gorm:"type:varchar(200)" json:"labels"` // comma separated
}

// NoteView is a simplified representation of a Note for the UI
type NoteView struct {
	ID        uint   `json:"id"`
	Title     string `json:"title"`
	Labels    string `json:"labels"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func ToNoteView(n Note) NoteView {
	return NoteView{
		ID:        n.ID,
		Title:     n.Title,
		Labels:    n.Labels,
		CreatedAt: n.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: n.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func IsValidMode(mode string) bool {
	for _, item := range config.AllModes {
		if item == mode {
			return true
		}
	}
	return false
}

type FocusStats struct {
	TotalSessions  int
	TotalDuration  string
	RawDuration    float64
	LongestSession string
}

type HabitStats struct {
	HabitName      string
	CompletionRate float64
	LogsTracked    int
	LogsCompleted  int
}

type IntentStats struct {
	IntentName string
	Status     string
}

type SummaryReport struct {
	DateRange      string
	Focus          FocusStats
	Habits         []HabitStats
	Intents        []IntentStats
	Points         PointStats
	TasksCompleted int
}
// SideQuest represents an ad-hoc task during a focus session
type SideQuest struct {
	gorm.Model
	Title   string     `gorm:"not null" json:"title"`
	Status  string     `gorm:"default:active" json:"status"` // active, done
	Note    string     `gorm:"type:text" json:"note"`
	EndedAt *time.Time `json:"ended_at,omitempty"`
}

// Task represents a to-do item linked to an intent or focus session
type Task struct {
	gorm.Model
	Title          string    `gorm:"not null" json:"title"`
	Status         string    `gorm:"default:pending" json:"status"` // pending, completed
	IntentID       *uint     `json:"intent_id,omitempty"`
	FocusSessionID *uint     `json:"focus_session_id,omitempty"`
	Position       int       `gorm:"default:0" json:"position"`
	SubTasks       []SubTask `gorm:"foreignKey:TaskID" json:"sub_tasks"`
}

// SubTask is a smaller component of a Task
type SubTask struct {
	gorm.Model
	TaskID   uint   `gorm:"not null" json:"task_id"`
	Title    string `gorm:"not null" json:"title"`
	Status   string `gorm:"default:pending" json:"status"` // pending, completed
	Position int    `gorm:"default:0" json:"position"`
}

// SubTaskView is a simplified representation of a SubTask for the UI
type SubTaskView struct {
	ID     uint   `json:"id"`
	TaskID   uint   `json:"task_id"`
	Title    string `json:"title"`
	Status   string `json:"status"`
	Position int    `json:"position"`
}

func ToSubTaskView(st SubTask) SubTaskView {
	return SubTaskView{
		ID:     st.ID,
		TaskID:   st.TaskID,
		Title:    st.Title,
		Status:   st.Status,
		Position: st.Position,
	}
}

// TaskView is a simplified representation of a Task for the UI
type TaskView struct {
	ID             uint          `json:"id"`
	Title          string        `json:"title"`
	Status         string        `json:"status"`
	IntentID          *uint         `json:"intent_id,omitempty"`
	IntentName        string        `json:"intent_name,omitempty"`
	FocusSessionID    *uint         `json:"focus_session_id,omitempty"`
	FocusSessionTitle string        `json:"focus_session_title,omitempty"`
	Position          int           `json:"position"`
	SubTasks       []SubTaskView `json:"sub_tasks"`
	CreatedAt      string        `json:"created_at"`
}

func ToTaskView(t Task) TaskView {
	subTasks := make([]SubTaskView, len(t.SubTasks))
	for i, st := range t.SubTasks {
		subTasks[i] = ToSubTaskView(st)
	}

	return TaskView{
		ID:             t.ID,
		Title:          t.Title,
		Status:         t.Status,
		IntentID:       t.IntentID,
		FocusSessionID: t.FocusSessionID,
		Position:       t.Position,
		SubTasks:       subTasks,
		CreatedAt:      t.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

// SideQuestView is a simplified representation of a SideQuest for the UI
type SideQuestView struct {
	ID      uint   `json:"id"`
	Title   string `json:"title"`
	Status  string `json:"status"`
	Note    string `json:"note"`
	EndedAt string `json:"ended_at"`
}

func ToSideQuestView(sq SideQuest) SideQuestView {
	var ended string
	if sq.EndedAt != nil {
		ended = sq.EndedAt.Format("2006-01-02 15:04")
	} else {
		ended = "-"
	}
	return SideQuestView{
		ID:      sq.ID,
		Title:   sq.Title,
		Status:  sq.Status,
		Note:    sq.Note,
		EndedAt: ended,
	}
}

// PointCategory represents the type of activity that earned points
type PointCategory string

const (
	// CategoryFocus for focus sessions
	CategoryFocus PointCategory = "focus"
	// CategoryHabit for habit completions
	CategoryHabit PointCategory = "habit"
	// CategoryIntent for intent completions
	CategoryIntent PointCategory = "intent"
	// CategoryJournal for journal entries
	CategoryJournal PointCategory = "journal"
	// CategoryQuest for side quest completions
	CategoryQuest PointCategory = "quest"
	// CategoryTask for task completions
	CategoryTask PointCategory = "task"
	// CategorySubTask for subtask completions
	CategorySubTask PointCategory = "subtask"
)

// PointTransaction records points earned for a specific activity
type PointTransaction struct {
	gorm.Model
	ActivityType PointCategory `gorm:"type:varchar(50);not null" json:"activity_type"`
	ActivityID   uint          `gorm:"not null" json:"activity_id"`
	Points       int           `gorm:"not null" json:"points"`
}

// PointStats aggregates point information for reporting
type PointStats struct {
	TotalPoints int
	History     []PointTransaction
}
