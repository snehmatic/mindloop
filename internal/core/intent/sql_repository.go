package intent

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

func (r *SQLRepository) CreateIntent(intent *models.Intent) error {
	return r.db.Create(intent).Error
}

func (r *SQLRepository) ListIntents() ([]models.Intent, error) {
	var intents []models.Intent
	result := r.db.Order("CreatedAt DESC").Find(&intents)
	return intents, result.Error
}

func (r *SQLRepository) ListActiveIntents() ([]models.Intent, error) {
	var intents []models.Intent
	result := r.db.Where("status = ?", models.StatusActive).Order("CreatedAt DESC").Find(&intents)
	return intents, result.Error
}

func (r *SQLRepository) GetOngoingIntent() (*models.Intent, error) {
	var intents []models.Intent
	result := r.db.Where("status IN ?", []string{models.StatusActive, models.StatusPaused}).Limit(1).Find(&intents)
	if result.Error != nil {
		return nil, result.Error
	}
	if len(intents) == 0 {
		return nil, nil
	}
	return &intents[0], nil
}

func (r *SQLRepository) GetIntent(id string) (*models.Intent, error) {
	var intent models.Intent
	if err := r.db.Where("id = ?", id).First(&intent).Error; err != nil {
		return nil, err
	}
	return &intent, nil
}

func (r *SQLRepository) GetIntentByID(id uint) (*models.Intent, error) {
	var intent models.Intent
	if err := r.db.First(&intent, id).Error; err != nil {
		return nil, err
	}
	return &intent, nil
}

func (r *SQLRepository) UpdateIntent(intent *models.Intent) error {
	return r.db.Save(intent).Error
}

func (r *SQLRepository) DeleteIntent(id string) error {
	r.db.Model(&models.Task{}).Where("IntentID = ?", id).Update("IntentID", nil)
	return r.db.Delete(&models.Intent{}, "id = ?", id).Error
}

func (r *SQLRepository) DeleteAll() error {
	return r.db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.Intent{}).Error
}
