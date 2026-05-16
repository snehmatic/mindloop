package backup

import (
	"encoding/json"
	"os"

	"github.com/snehmatic/mindloop/models"
	"gorm.io/gorm"
)

// Service handles backup and restore operations for Mindloop
type Service struct {
	DB *gorm.DB
}

// NewService creates a new backup Service instance
func NewService(db *gorm.DB) *Service {
	return &Service{DB: db}
}

// Data represents the structure of the exported JSON backup file
type Data struct {
	Intents           []models.Intent           `json:"intents"`
	FocusSessions     []models.FocusSession     `json:"focus_sessions"`
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

	s.DB.Find(&data.Intents)
	s.DB.Find(&data.FocusSessions)
	s.DB.Find(&data.Habits)
	s.DB.Find(&data.HabitLogs)
	s.DB.Find(&data.JournalEntries)
	s.DB.Find(&data.Notes)
	s.DB.Find(&data.SideQuests)
	s.DB.Find(&data.Tasks)
	s.DB.Find(&data.SubTasks)
	s.DB.Find(&data.PointTransactions)

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

	// Use a transaction for import
	return s.DB.Transaction(func(tx *gorm.DB) error {
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