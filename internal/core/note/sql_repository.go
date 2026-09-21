package note

import (
	"github.com/snehmatic/mindloop/models"
	"gorm.io/gorm"
)

type SQLRepository struct {
	db *gorm.DB
}

func NewSQLRepository(db *gorm.DB) Repository {
	return &SQLRepository{db: db}
}

func (r *SQLRepository) GetDB() *gorm.DB {
	return r.db
}

func (r *SQLRepository) CreateNote(note *models.Note) error {
	return r.db.Create(note).Error
}

func (r *SQLRepository) ListNotes() ([]models.Note, error) {
	var notes []models.Note
	result := r.db.Order("UpdatedAt DESC").Find(&notes)
	return notes, result.Error
}

func (r *SQLRepository) GetNote(id int) (*models.Note, error) {
	var note models.Note
	if err := r.db.First(&note, id).Error; err != nil {
		return nil, err
	}
	return &note, nil
}

func (r *SQLRepository) UpdateNote(note *models.Note) error {
	return r.db.Save(note).Error
}

func (r *SQLRepository) DeleteNote(id int) error {
	return r.db.Delete(&models.Note{}, id).Error
}

func (r *SQLRepository) DeleteAll() error {
	return r.db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.Note{}).Error
}
