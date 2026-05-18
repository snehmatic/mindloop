package journal

import (
	"github.com/snehmatic/mindloop/models"
)

// Repository defines the interface for journal data access
type Repository interface {
	CreateEntry(title, content, mood string, pointsToAward int) (bool, error)
	ListEntries() ([]models.JournalEntry, error)
	GetEntry(id string) (models.JournalEntry, error)
	UpdateEntry(entry *models.JournalEntry) error
	DeleteEntry(id string) error
	DeleteAll() error
}
