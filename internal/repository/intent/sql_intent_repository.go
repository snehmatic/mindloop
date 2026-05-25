package intent

import (
	"time"

	"github.com/snehmatic/mindloop/internal/core/points"
	"github.com/snehmatic/mindloop/models"
	"gorm.io/gorm"
)

// sqlRepository implements Repository using GORM
type sqlRepository struct {
	DB *gorm.DB
}

// NewSQLRepository creates a new SQL-based intent repository
func NewSQLRepository(db *gorm.DB) Repository {
	return &sqlRepository{DB: db}
}

// StartIntent creates a new intent
func (r *sqlRepository) StartIntent(name string) (*models.Intent, error) {
	if name == "" {
		return nil, ErrNameCannotBeEmpty
	}

	intent := &models.Intent{
		Name:   name,
		Status: models.IntentStatusActive,
	}

	if err := r.DB.Create(intent).Error; err != nil {
		return nil, err
	}
	return intent, nil
}

// ListIntents retrieves all intents from the database
func (r *sqlRepository) ListIntents() ([]models.Intent, error) {
	var intents []models.Intent
	result := r.DB.Order("CreatedAt DESC").Find(&intents)
	return intents, result.Error
}

// ListActiveIntents retrieves all active intents
func (r *sqlRepository) ListActiveIntents() ([]models.Intent, error) {
	var intents []models.Intent
	result := r.DB.Where("status = ?", models.IntentStatusActive).Order("CreatedAt DESC").Find(&intents)
	return intents, result.Error
}

// GetOngoingIntent retrieves the ongoing intent (active or paused)
func (r *sqlRepository) GetOngoingIntent() (*models.Intent, error) {
	var intents []models.Intent
	result := r.DB.Where("status IN ?", []string{string(models.IntentStatusActive), string(models.IntentStatusPaused)}).Limit(1).Find(&intents)
	if result.Error != nil {
		return nil, result.Error
	}
	if len(intents) == 0 {
		return nil, nil
	}
	return &intents[0], nil
}

// GetIntent retrieves a single intent by its ID
func (r *sqlRepository) GetIntent(id string) (*models.Intent, error) {
	var intent models.Intent
	if err := r.DB.Where("id = ?", id).First(&intent).Error; err != nil {
		return nil, err
	}
	return &intent, nil
}

// UpdateIntent modifies an existing intent in the database
func (r *sqlRepository) UpdateIntent(intent *models.Intent) error {
	return r.DB.Save(intent).Error
}

// EndIntent ends an intent and awards points if successful
func (r *sqlRepository) EndIntent(idStr string, pointsToAward int) (*models.Intent, bool, error) {
	var intent models.Intent
	if err := r.DB.Where("id = ?", idStr).First(&intent).Error; err != nil {
		return nil, false, err
	}

	now := time.Now()
	intent.Status = models.IntentStatusDone
	intent.EndedAt = &now

	if err := r.DB.Save(&intent).Error; err != nil {
		return nil, false, err
	}

	var milestoneReached bool
	var totalPoints int
	if err := r.DB.Model(&models.PointTransaction{}).Select("COALESCE(SUM(Points), 0)").Scan(&totalPoints).Error; err == nil {
		tx := models.PointTransaction{
			ActivityType: models.CategoryIntent,
			ActivityID:   intent.ID,
			Points:       pointsToAward,
		}
		if r.DB.Create(&tx).Error == nil {
			newTotal := totalPoints + pointsToAward
			if newTotal/points.MilestoneInterval > totalPoints/points.MilestoneInterval {
				milestoneReached = true
			}
		}
	}

	return &intent, milestoneReached, nil
}

// DeleteIntent removes an intent from the database
func (r *sqlRepository) DeleteIntent(id string) error {
	r.DB.Model(&models.Task{}).Where("IntentID = ?", id).Update("IntentID", nil)
	return r.DB.Delete(&models.Intent{}, "id = ?", id).Error
}

// DeleteAll removes all intents from the database
func (r *sqlRepository) DeleteAll() error {
	return r.DB.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.Intent{}).Error
}

// PauseIntent pauses an intent
func (r *sqlRepository) PauseIntent(id uint) (*models.Intent, error) {
	var intent models.Intent
	if err := r.DB.First(&intent, id).Error; err != nil {
		return nil, err
	}

	if intent.Status != models.IntentStatusActive {
		return nil, ErrIntentNotActive
	}

	intent.Status = models.IntentStatusPaused
	if err := r.DB.Save(&intent).Error; err != nil {
		return nil, err
	}
	return &intent, nil
}

// ResumeIntent resumes an intent
func (r *sqlRepository) ResumeIntent(id uint) (*models.Intent, error) {
	var intent models.Intent
	if err := r.DB.First(&intent, id).Error; err != nil {
		return nil, err
	}

	if intent.Status != models.IntentStatusPaused {
		return nil, ErrIntentNotPaused
	}

	intent.Status = models.IntentStatusActive
	if err := r.DB.Save(&intent).Error; err != nil {
		return nil, err
	}
	return &intent, nil
}

// GetIntentStats returns intent statistics for the given time range
func (r *sqlRepository) GetIntentStats(start, end time.Time) ([]models.IntentStats, error) {
	var intents []models.Intent
	rangeQuery := "CreatedAt >= ? AND CreatedAt <= ?"
	if err := r.DB.Where(rangeQuery, start, end).Order("CreatedAt DESC").Find(&intents).Error; err != nil {
		return nil, err
	}

	if len(intents) == 0 {
		return []models.IntentStats{}, nil
	}

	var stats []models.IntentStats
	for _, intent := range intents {
		stats = append(stats, models.IntentStats{
			IntentName: intent.Name,
			Status:     intent.Status,
		})
	}
	return stats, nil
}

// GetDB returns the GORM DB instance
func (r *sqlRepository) GetDB() *gorm.DB {
	return r.DB
}

// GetIntentsInRange retrieves all intents within the given time range
func (r *sqlRepository) GetIntentsInRange(start, end time.Time) ([]models.Intent, error) {
	var intents []models.Intent
	rangeQuery := "CreatedAt >= ? AND CreatedAt <= ?"
	if err := r.DB.Where(rangeQuery, start, end).Order("CreatedAt DESC").Find(&intents).Error; err != nil {
		return nil, err
	}
	return intents, nil
}
