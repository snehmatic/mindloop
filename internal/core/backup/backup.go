package backup

import (
	"encoding/json"
	"os"
	"reflect"

	"github.com/snehmatic/mindloop/internal/core/points"
	"github.com/snehmatic/mindloop/internal/repository/focus"
	"github.com/snehmatic/mindloop/internal/repository/habit"
	"github.com/snehmatic/mindloop/internal/repository/habitlog"
	"github.com/snehmatic/mindloop/internal/repository/intent"
	"github.com/snehmatic/mindloop/internal/repository/journal"
	"github.com/snehmatic/mindloop/internal/repository/note"
	"github.com/snehmatic/mindloop/internal/repository/point"
	"github.com/snehmatic/mindloop/internal/repository/quest"
	"github.com/snehmatic/mindloop/internal/repository/routine"
	"github.com/snehmatic/mindloop/internal/repository/subtask"
	"github.com/snehmatic/mindloop/internal/repository/task"
	"github.com/snehmatic/mindloop/models"
	"gorm.io/gorm"
)

// Service handles backup and restore operations for Mindloop
type Service struct {
	focusRepo    focus.Repository
	habitRepo    habit.Repository
	habitLogRepo habitlog.Repository
	intentRepo   intent.Repository
	journalRepo  journal.Repository
	noteRepo     note.Repository
	pointRepo    point.Repository
	pointSvc     points.Service
	questRepo    quest.Repository
	routineRepo  routine.Repository
	subtaskRepo  subtask.Repository
	taskRepo     task.TaskRepository
	db           *gorm.DB
}

// NewService creates a new backup Service instance
func NewService(
	db *gorm.DB,
	pointSvc points.Service,
	focusRepo focus.Repository,
	habitRepo habit.Repository,
	habitLogRepo habitlog.Repository,
	intentRepo intent.Repository,
	journalRepo journal.Repository,
	noteRepo note.Repository,
	pointRepo point.Repository,
	questRepo quest.Repository,
	routineRepo routine.Repository,
	subtaskRepo subtask.Repository,
	taskRepo task.TaskRepository,
) *Service {
	return &Service{
		db:           db,
		pointSvc:     pointSvc,
		focusRepo:    focusRepo,
		habitRepo:    habitRepo,
		habitLogRepo: habitLogRepo,
		intentRepo:   intentRepo,
		journalRepo:  journalRepo,
		noteRepo:     noteRepo,
		pointRepo:    pointRepo,
		questRepo:    questRepo,
		routineRepo:  routineRepo,
		subtaskRepo:  subtaskRepo,
		taskRepo:     taskRepo,
	}
}

// Data represents the structure of the exported JSON backup file
type Data struct {
	Intents           []models.Intent           `json:"intents"`
	FocusSessions     []models.FocusSession     `json:"focus_sessions"`
	Routines          []models.Routine          `json:"routines,omitempty"`
	Habits            []models.Habit            `json:"habits"`
	HabitLogs         []models.HabitLog         `json:"habit_logs"`
	JournalEntries    []models.JournalEntry     `json:"journal_entries"`
	Notes             []models.Note             `json:"notes,omitempty"`
	SideQuests        []models.SideQuest        `json:"side_quests,omitempty"`
	Tasks             []models.Task             `json:"tasks,omitempty"`
	SubTasks          []models.SubTask          `json:"sub_tasks,omitempty"`
	PointTransactions []models.PointTransaction `json:"point_transactions,omitempty"`
}

// Export saves all application data to a JSON file
func (s *Service) Export(filePath string) error {
	var data Data

	intents, err := s.intentRepo.ListIntents()
	if err != nil {
		return err
	}
	data.Intents = intents

	focusSessions, err := s.focusRepo.ListSessions()
	if err != nil {
		return err
	}
	data.FocusSessions = focusSessions

	routines, err := s.routineRepo.FindRoutines()
	if err != nil {
		return err
	}
	data.Routines = routines

	habits, err := s.habitRepo.ListHabits("")
	if err != nil {
		return err
	}
	data.Habits = habits

	habitLogs, err := s.habitLogRepo.FindHabitLogs()
	if err != nil {
		return err
	}
	data.HabitLogs = habitLogs

	journalEntries, err := s.journalRepo.ListEntries()
	if err != nil {
		return err
	}
	data.JournalEntries = journalEntries

	notes, err := s.noteRepo.ListNotes()
	if err != nil {
		return err
	}
	data.Notes = notes

	sideQuests, err := s.questRepo.ListQuests()
	if err != nil {
		return err
	}
	data.SideQuests = sideQuests

	tasks, err := s.taskRepo.ListTasks()
	if err != nil {
		return err
	}
	data.Tasks = tasks

	subTasks, err := s.subtaskRepo.FindSubTasks()
	if err != nil {
		return err
	}
	data.SubTasks = subTasks

	// Use pointSvc to retrieve all point transactions
	pointTransactions, err := s.pointSvc.GetPointsInRange("", "")
	if err != nil {
		return err
	}
	data.PointTransactions = pointTransactions

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, jsonData, 0644)
}

// Import restores application data from a JSON file
func (s *Service) Import(filePath string) error {
	jsonData, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	var data Data
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return err
	}

	// Zero out IDs so SQLite/gorm generates fresh ones on re-insert
	data.stripIDs()

	// Use a transaction for import
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Clear existing data before restoring. Use Unscoped deletes so soft-deleted
		// rows do not keep primary keys occupied when the backup recreates records
		// with their original IDs.
		if err := tx.Session(&gorm.Session{AllowGlobalUpdate: true}).Unscoped().Delete(&models.PointTransaction{}).Error; err != nil {
			return err
		}
		if err := tx.Session(&gorm.Session{AllowGlobalUpdate: true}).Unscoped().Delete(&models.SubTask{}).Error; err != nil {
			return err
		}
		if err := tx.Session(&gorm.Session{AllowGlobalUpdate: true}).Unscoped().Delete(&models.Task{}).Error; err != nil {
			return err
		}
		if err := tx.Session(&gorm.Session{AllowGlobalUpdate: true}).Unscoped().Delete(&models.SideQuest{}).Error; err != nil {
			return err
		}
		if err := tx.Session(&gorm.Session{AllowGlobalUpdate: true}).Unscoped().Delete(&models.Intent{}).Error; err != nil {
			return err
		}
		if err := tx.Session(&gorm.Session{AllowGlobalUpdate: true}).Unscoped().Delete(&models.FocusSession{}).Error; err != nil {
			return err
		}
		if err := tx.Session(&gorm.Session{AllowGlobalUpdate: true}).Unscoped().Delete(&models.HabitLog{}).Error; err != nil {
			return err
		}
		if err := tx.Session(&gorm.Session{AllowGlobalUpdate: true}).Unscoped().Delete(&models.Habit{}).Error; err != nil {
			return err
		}
		if err := tx.Session(&gorm.Session{AllowGlobalUpdate: true}).Unscoped().Delete(&models.Routine{}).Error; err != nil {
			return err
		}
		if err := tx.Session(&gorm.Session{AllowGlobalUpdate: true}).Unscoped().Delete(&models.JournalEntry{}).Error; err != nil {
			return err
		}
		if err := tx.Session(&gorm.Session{AllowGlobalUpdate: true}).Unscoped().Delete(&models.Note{}).Error; err != nil {
			return err
		}

		if len(data.Intents) > 0 {
			if err := tx.Create(&data.Intents).Error; err != nil {
				return err
			}
		}
		if len(data.FocusSessions) > 0 {
			if err := tx.Create(&data.FocusSessions).Error; err != nil {
				return err
			}
		}
		if len(data.Routines) > 0 {
			if err := tx.Create(&data.Routines).Error; err != nil {
				return err
			}
		}
		if len(data.Habits) > 0 {
			if err := tx.Create(&data.Habits).Error; err != nil {
				return err
			}
		}
		if len(data.HabitLogs) > 0 {
			if err := tx.Create(&data.HabitLogs).Error; err != nil {
				return err
			}
		}
		if len(data.JournalEntries) > 0 {
			if err := tx.Create(&data.JournalEntries).Error; err != nil {
				return err
			}
		}
		if len(data.Notes) > 0 {
			if err := tx.Create(&data.Notes).Error; err != nil {
				return err
			}
		}
		if len(data.SideQuests) > 0 {
			if err := tx.Create(&data.SideQuests).Error; err != nil {
				return err
			}
		}
		if len(data.Tasks) > 0 {
			if err := tx.Create(&data.Tasks).Error; err != nil {
				return err
			}
		}
		if len(data.SubTasks) > 0 {
			if err := tx.Create(&data.SubTasks).Error; err != nil {
				return err
			}
		}
		if len(data.PointTransactions) > 0 {
			if err := tx.Create(&data.PointTransactions).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// stripIDs zeroes out the ID of every exported model so that gorm/SQLite
// assigns fresh autoincrement values on re-insert rather than colliding with
// still-existent keys when Unscoped deletes do not immediately release IDs.
func (d *Data) stripIDs() {
	strip := func(slicePtr interface{}) {
		v := reflect.ValueOf(slicePtr)
		if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Slice {
			return
		}
		sl := v.Elem()
		for i := 0; i < sl.Len(); i++ {
			if idField := sl.Index(i).FieldByName("ID"); idField.IsValid() {
				idField.SetUint(0)
			}
		}
	}
	strip(&d.PointTransactions)
	strip(&d.SubTasks)
	strip(&d.Tasks)
	strip(&d.SideQuests)
	strip(&d.Intents)
	strip(&d.FocusSessions)
	strip(&d.HabitLogs)
	strip(&d.Habits)
	strip(&d.Routines)
	strip(&d.JournalEntries)
	strip(&d.Notes)
}
