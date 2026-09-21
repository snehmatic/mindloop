package note

import (
	"fmt"
	"unicode/utf8"

	"github.com/snehmatic/mindloop/models"
)

const maxLabelsLength = 200

// UpdateInput describes fields to change. A nil field is left unchanged.
type UpdateInput struct {
	Title   *string
	Content *string
	Labels  *string
}

func validateLabels(labels string) error {
	if utf8.RuneCountInString(labels) > maxLabelsLength {
		return fmt.Errorf("labels cannot exceed %d characters", maxLabelsLength)
	}
	return nil
}

// Service handles business logic for markdown notes
type Service struct {
	repo Repository
}

// NewService creates a new note Service instance
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// CreateNote persists a new markdown note to the database
func (s *Service) CreateNote(title, content, labels string) (*models.Note, error) {
	if title == "" && content == "" {
		return nil, ErrNoteMustHaveATitleOrContent
	}
	if err := validateLabels(labels); err != nil {
		return nil, err
	}
	note := &models.Note{
		Title:   title,
		Content: content,
		Labels:  labels,
	}
	if err := s.repo.CreateNote(note); err != nil {
		return nil, err
	}
	return note, nil
}

// ListNotes retrieves all markdown notes from the database
func (s *Service) ListNotes() ([]models.Note, error) {
	return s.repo.ListNotes()
}

// GetNote retrieves a single markdown note by its ID
func (s *Service) GetNote(id int) (*models.Note, error) {
	return s.repo.GetNote(id)
}

// UpdateNote modifies an existing markdown note in the database
func (s *Service) UpdateNote(id int, title, content, labels string) (*models.Note, error) {
	return s.UpdateNoteFields(id, UpdateInput{
		Title:   &title,
		Content: &content,
		Labels:  &labels,
	})
}

// UpdateNoteFields updates only the fields supplied by the caller.
func (s *Service) UpdateNoteFields(id int, input UpdateInput) (*models.Note, error) {
	note, err := s.GetNote(id)
	if err != nil {
		return nil, err
	}
	if input.Title != nil {
		note.Title = *input.Title
	}
	if input.Content != nil {
		note.Content = *input.Content
	}
	if input.Labels != nil {
		if err := validateLabels(*input.Labels); err != nil {
			return nil, err
		}
		note.Labels = *input.Labels
	}
	if err := s.repo.UpdateNote(note); err != nil {
		return nil, err
	}
	return note, nil
}

// DeleteNote removes a markdown note from the database by its ID
func (s *Service) DeleteNote(id int) error {
	return s.repo.DeleteNote(id)
}

// DeleteAll removes all markdown notes from the database
func (s *Service) DeleteAll() error {
	return s.repo.DeleteAll()
}
