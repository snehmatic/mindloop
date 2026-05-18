package journal

import (
	"github.com/snehmatic/mindloop/internal/core/points"
	"github.com/snehmatic/mindloop/models"
	"gorm.io/gorm"
)

// sqlRepository implements Repository using GORM
type sqlRepository struct {
	DB *gorm.DB
}

// NewSQLRepository creates a new SQL-based journal repository
func NewSQLRepository(db *gorm.DB) Repository {
	return &sqlRepository{DB: db}
}

// CreateEntry creates a new journal entry and awards points if successful
func (r *sqlRepository) CreateEntry(title, content, mood string, pointsToAward int) (bool, error) {
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

	err := r.DB.Create(&entry).Error
	milestoneReached := false
	if err == nil {
		var totalPoints int
		if err := r.DB.Model(&models.PointTransaction{}).Select("COALESCE(SUM(Points), 0)").Scan(&totalPoints).Error; err == nil {
			tx := models.PointTransaction{
				ActivityType: models.CategoryJournal,
				ActivityID:   entry.ID,
				Points:       pointsToAward,
			}
			if r.DB.Create(&tx).Error == nil {
				newTotal := totalPoints + pointsToAward
				if newTotal/points.MilestoneInterval > totalPoints/points.MilestoneInterval {
					milestoneReached = true
				}
			}
		}
	}
	return milestoneReached, err
}

// ListEntries retrieves all journal entries from the database
func (r *sqlRepository) ListEntries() ([]models.JournalEntry, error) {
	var entries []models.JournalEntry
	result := r.DB.Order("CreatedAt DESC").Find(&entries)
	return entries, result.Error
}

// GetEntry retrieves a single journal entry by its ID
func (r *sqlRepository) GetEntry(id string) (models.JournalEntry, error) {
	var entry models.JournalEntry
	result := r.DB.First(&entry, id)
	return entry, result.Error
}

// UpdateEntry modifies an existing journal entry in the database
func (r *sqlRepository) UpdateEntry(entry *models.JournalEntry) error {
	return r.DB.Save(entry).Error
}

// DeleteEntry removes a journal entry from the database by its ID
func (r *sqlRepository) DeleteEntry(id string) error {
	result := r.DB.Delete(&models.JournalEntry{}, id)
	return result.Error
}

// DeleteAll removes all journal entries from the database
func (r *sqlRepository) DeleteAll() error {
	return r.DB.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.JournalEntry{}).Error
}
