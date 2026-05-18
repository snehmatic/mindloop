package journal

import (
	"github.com/snehmatic/mindloop/internal/repository/journal"
	"github.com/snehmatic/mindloop/models"
)

// Service handles business logic for journal entries
type Service struct {
	repository journal.Repository
}

// NewService creates a new journal Service instance
func NewService(repo journal.Repository) *Service {
	return &Service{repository: repo}
}

func (s *Service) CreateEntry(title, content, mood string, pointsToAward int) (bool, error) {
	return s.repository.CreateEntry(title, content, mood, pointsToAward)
}

func (s *Service) ListEntries() ([]models.JournalEntry, error) {
	return s.repository.ListEntries()
}

func (s *Service) GetEntry(id string) (models.JournalEntry, error) {
	return s.repository.GetEntry(id)
}

func (s *Service) UpdateEntry(entry *models.JournalEntry) error {
	return s.repository.UpdateEntry(entry)
}

func (s *Service) DeleteEntry(id string) error {
	return s.repository.DeleteEntry(id)
}

func (s *Service) DeleteAll() error {
	return s.repository.DeleteAll()
}
