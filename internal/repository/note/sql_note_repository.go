package note

import (
	"github.com/snehmatic/mindloop/models"
	"gorm.io/gorm"
)

// sqlRepository implements Repository using GORM
type sqlRepository struct {
	DB *gorm.DB
}

// NewSQLRepository creates a new SQL-based note repository
func NewSQLRepository(db *gorm.DB) Repository {
	return &sqlRepository{DB: db}
}

// CreateNote persists a new markdown note to the database
func (r *sqlRepository) CreateNote(title, content, labels string) (*models.Note, error) {
	if title == "" && content == "" {
		return nil, ErrNoteMustHaveTitleOrContent
	}
	note := &models.Note{
		Title:   title,
		Content: content,
		Labels:  labels,
	}
	if err := r.DB.Create(note).Error; err != nil {
		return nil, err
	}
	return note, nil
}

// ListNotes retrieves all markdown notes from the database
func (r *sqlRepository) ListNotes() ([]models.Note, error) {
	var notes []models.Note
	if err := r.DB.Order("UpdatedAt desc").Find(&notes).Error; err != nil {
		return nil, err
	}
	return notes, nil
}

// GetNote retrieves a single markdown note by its ID
func (r *sqlRepository) GetNote(id int) (*models.Note, error) {
	var note models.Note
	if err := r.DB.First(&note, id).Error; err != nil {
		return nil, err
	}
	return &note, nil
}

// UpdateNote modifies an existing markdown note in the database
func (r *sqlRepository) UpdateNote(id int, title, content, labels string) (*models.Note, error) {
	note, err := r.GetNote(id)
	if err != nil {
		return nil, err
	}
	note.Title = title
	note.Content = content
	note.Labels = labels
	if err := r.DB.Save(note).Error; err != nil {
		return nil, err
	}
	return note, nil
}

// DeleteNote removes a markdown note from the database by its ID
func (r *sqlRepository) DeleteNote(id int) error {
	return r.DB.Delete(&models.Note{}, id).Error
}

// DeleteAll removes all markdown notes from the database
func (r *sqlRepository) DeleteAll() error {
	return r.DB.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.Note{}).Error
}
