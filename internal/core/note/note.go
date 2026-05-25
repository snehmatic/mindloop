package note

import (
	"github.com/snehmatic/mindloop/internal/repository/note"
	"github.com/snehmatic/mindloop/models"
)

// Service handles business logic for markdown notes
type Service struct {
	repository note.Repository
}

// NewService creates a new note Service instance
func NewService(repo note.Repository) *Service {
	return &Service{repository: repo}
}

// CreateNote persists a new markdown note to the database
func (s *Service) CreateNote(title, content, labels string) (*models.Note, error) {
	if title == "" && content == "" {
		return nil, ErrNoteMustHaveTitleOrContent
	}
	return s.repository.CreateNote(title, content, labels)
}

// ListNotes retrieves all markdown notes from the database
func (s *Service) ListNotes() ([]models.Note, error) {
	return s.repository.ListNotes()
}

// GetNote retrieves a single markdown note by its ID
func (s *Service) GetNote(id int) (*models.Note, error) {
	return s.repository.GetNote(id)
}

// UpdateNote modifies an existing markdown note in the database
func (s *Service) UpdateNote(id int, title, content, labels string) (*models.Note, error) {
	return s.repository.UpdateNote(id, title, content, labels)
}

// DeleteNote removes a markdown note from the database by its ID
func (s *Service) DeleteNote(id int) error {
	return s.repository.DeleteNote(id)
}

// DeleteAll removes all markdown notes from the database
func (s *Service) DeleteAll() error {
	return s.repository.DeleteAll()
}
