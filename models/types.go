package models

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
)

// model definitions reside here
// request/response structs, etc.

type MindloopMode string
type IntervalType string

var AllModes = [...]string{"local", "byodb", "api"}
var AllIntervalTypes = [...]string{"hour", "day", "week"}

var (
	Local MindloopMode = MindloopMode(AllModes[0])
	ByoDB MindloopMode = MindloopMode(AllModes[1])
	Api   MindloopMode = MindloopMode(AllModes[2])
)

var (
	Hour IntervalType = IntervalType(AllIntervalTypes[0])
	Day  IntervalType = IntervalType(AllIntervalTypes[1])
	Week IntervalType = IntervalType(AllIntervalTypes[2])
)

type UserConfig struct {
	Name string `yaml:"name"`
	Mode string `yaml:"mode"`
}

// UserConfigPath is the file path where the user configuration YAML will be written.
// ToDo: Make this configurable or use a constant
var UserConfigPath = "user_config.yaml"

type Habit struct {
	gorm.Model
	Name        string       `gorm:"type:varchar(100)" json:"name"`
	Interval    IntervalType `gorm:"type:varchar(100)" json:"internal"`
	ActualCount int          `gorm:"type:int" json:"actual_count"`
	TargetCount int          `gorm:"type:int" json:"target_count"`
}

type Intent struct {
	gorm.Model
	Message string         `gorm:"not null" json:"message"`
	Tags    string         `gorm:"type:text" json:"tags"`        // comma-separated or JSON later
	Status  string         `gorm:"default:active" json:"status"` // active, ended, archived
	EndedAt *time.Time     `json:"ended_at,omitempty"`
	Focuses []FocusSession `gorm:"foreignKey:IntentID" json:"focus_sessions,omitempty"`
}

type FocusSession struct {
	gorm.Model
	IntentID  int       `gorm:"not null;index" json:"intent_id"`
	StartTime time.Time `gorm:"autoCreateTime" json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Duration  int       `json:"duration"`                // in seconds
	Rating    int       `gorm:"default:0" json:"rating"` // 1 to 5, optional
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
	for _, item := range AllModes {
		if item == mode {
			return true
		}
	}
	return false
}

func (uc UserConfig) WriteToYAML() {
	marshalled, err := yaml.Marshal(uc)
	if err != nil {
		fmt.Println("Error marshalling user config to YAML")
		return
	}
	err = os.WriteFile(UserConfigPath, marshalled, 0644)
	if err != nil {
		fmt.Println("Error writing user config to file")
		return
	}
	fmt.Println("User config written to YAML successfully")
}
