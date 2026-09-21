package journal

import (
	"github.com/snehmatic/mindloop/models"
	"gorm.io/gorm"
)

type Repository interface {
	CreateEntry(entry *models.JournalEntry) error
	ListEntries() ([]models.JournalEntry, error)
	GetEntry(id string) (*models.JournalEntry, error)
	UpdateEntry(entry *models.JournalEntry) error
	DeleteEntry(id string) error
	DeleteAll() error

	GetDB() *gorm.DB
}
