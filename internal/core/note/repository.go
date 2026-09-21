package note

import (
	"github.com/snehmatic/mindloop/models"
	"gorm.io/gorm"
)

type Repository interface {
	CreateNote(note *models.Note) error
	ListNotes() ([]models.Note, error)
	GetNote(id int) (*models.Note, error)
	UpdateNote(note *models.Note) error
	DeleteNote(id int) error
	DeleteAll() error

	GetDB() *gorm.DB
}
