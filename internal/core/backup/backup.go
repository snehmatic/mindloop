package backup

import (
	"encoding/json"
	"os"

	"github.com/snehmatic/mindloop/models"
	"gorm.io/gorm"
)

type Service struct {
	DB *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{DB: db}
}

type BackupData struct {
	Intents        []models.Intent       `json:"intents"`
	FocusSessions  []models.FocusSession `json:"focus_sessions"`
	Habits         []models.Habit        `json:"habits"`
	HabitLogs      []models.HabitLog     `json:"habit_logs"`
	JournalEntries []models.JournalEntry `json:"journal_entries"`
	Notes          []models.Note         `json:"notes,omitempty"`
}

func (s *Service) Export(filePath string) error {
	var data BackupData

	s.DB.Find(&data.Intents)
	s.DB.Find(&data.FocusSessions)
	s.DB.Find(&data.Habits)
	s.DB.Find(&data.HabitLogs)
	s.DB.Find(&data.JournalEntries)
	s.DB.Find(&data.Notes)

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, jsonData, 0644)
}

func (s *Service) Import(filePath string) error {
	jsonData, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	var data BackupData
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return err
	}

	// Use a transaction for import
	return s.DB.Transaction(func(tx *gorm.DB) error {
		// Clear existing data before restoring
		// This ensures we don't have primary key conflicts and matches "restore" behavior
		tx.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.Intent{})
		tx.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.FocusSession{})
		tx.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.Habit{})
		tx.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.HabitLog{})
		tx.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.JournalEntry{})
		tx.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.Note{})

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

		return nil
	})
}
