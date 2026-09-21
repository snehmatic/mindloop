package journal

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

func (r *SQLRepository) CreateEntry(entry *models.JournalEntry) error {
	return r.db.Create(entry).Error
}

func (r *SQLRepository) ListEntries() ([]models.JournalEntry, error) {
	var entries []models.JournalEntry
	result := r.db.Order("CreatedAt DESC").Find(&entries)
	return entries, result.Error
}

func (r *SQLRepository) GetEntry(id string) (*models.JournalEntry, error) {
	var entry models.JournalEntry
	if err := r.db.First(&entry, id).Error; err != nil {
		return nil, err
	}
	return &entry, nil
}

func (r *SQLRepository) UpdateEntry(entry *models.JournalEntry) error {
	return r.db.Save(entry).Error
}

func (r *SQLRepository) DeleteEntry(id string) error {
	return r.db.Delete(&models.JournalEntry{}, id).Error
}

func (r *SQLRepository) DeleteAll() error {
	return r.db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.JournalEntry{}).Error
}
