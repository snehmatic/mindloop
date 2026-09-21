package journal

import (
	"github.com/snehmatic/mindloop/internal/core/points"
	"github.com/snehmatic/mindloop/models"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateEntry(title, content, mood string, pointsToAward int) (bool, error) {
	if title == "" {
		return false, ErrTitleCannotBeEmpty
	}
	if content == "" {
		return false, ErrContentCannotBeEmpty
	}
	if mood == "" {
		mood = "neutral"
	}

	entry := models.JournalEntry{
		Title:   title,
		Content: content,
		Mood:    mood,
	}

	err := s.repo.CreateEntry(&entry)
	milestoneReached := false
	if err == nil {
		milestoneReached, _ = points.AwardPoints(s.repo.GetDB(), models.CategoryJournal, entry.ID, pointsToAward)
	}
	return milestoneReached, err
}

func (s *Service) ListEntries() ([]models.JournalEntry, error) {
	return s.repo.ListEntries()
}

func (s *Service) GetEntry(id string) (models.JournalEntry, error) {
	e, err := s.repo.GetEntry(id)
	if e != nil {
		return *e, err
	}
	return models.JournalEntry{}, err
}

func (s *Service) UpdateEntry(entry *models.JournalEntry) error {
	return s.repo.UpdateEntry(entry)
}

func (s *Service) DeleteEntry(id string) error {
	return s.repo.DeleteEntry(id)
}

func (s *Service) DeleteAll() error {
	return s.repo.DeleteAll()
}
