package note

import (
	"github.com/snehmatic/mindloop/models"
)

// Repository defines the interface for note data access
type Repository interface {
	CreateNote(title, content, labels string) (*models.Note, error)
	ListNotes() ([]models.Note, error)
	GetNote(id int) (*models.Note, error)
	UpdateNote(id int, title, content, labels string) (*models.Note, error)
	DeleteNote(id int) error
	DeleteAll() error
}
