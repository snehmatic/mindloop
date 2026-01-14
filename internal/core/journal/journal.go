package journal

import (
	"errors"

	"github.com/snehmatic/mindloop/models"
	"gorm.io/gorm"
)

type Service struct {
	DB *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{DB: db}
}

func (s *Service) CreateEntry(title, content, mood string) error {
	if title == "" {
		return errors.New("title cannot be empty")
	}
	if content == "" {
		return errors.New("content cannot be empty")
	}
	if mood == "" {
		mood = "neutral"
	}

	entry := models.JournalEntry{
		Title:   title,
		Content: content,
		Mood:    mood,
	}

	return s.DB.Create(&entry).Error
}

func (s *Service) ListEntries() ([]models.JournalEntry, error) {
	var entries []models.JournalEntry
	result := s.DB.Order("CreatedAt DESC").Find(&entries)
	return entries, result.Error
}

func (s *Service) GetEntry(id string) (models.JournalEntry, error) {
	var entry models.JournalEntry
	result := s.DB.First(&entry, "id = ?", id)
	return entry, result.Error
}

func (s *Service) DeleteEntry(id string) error {
	result := s.DB.Delete(&models.JournalEntry{}, "id = ?", id)
	return result.Error
}

func (s *Service) DeleteAll() error {
	return s.DB.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.JournalEntry{}).Error
}
