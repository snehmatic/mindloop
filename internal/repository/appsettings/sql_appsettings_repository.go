package appsettings

import (
	"github.com/snehmatic/mindloop/models"
	"gorm.io/gorm"
)

// sqlRepository implements Repository using GORM.
type sqlRepository struct {
	DB *gorm.DB
}

// NewSQLRepository creates a new SQL-based repository.
func NewSQLRepository(db *gorm.DB) Repository {
	return &sqlRepository{DB: db}
}

// GetSetting retrieves a single setting by key.
func (r *sqlRepository) GetSetting(key string) (*models.AppSetting, error) {
	var setting models.AppSetting
	result := r.DB.Where("key = ?", key).Limit(1).Find(&setting)
	if result.RowsAffected == 0 {
		return nil, result.Error
	}
	return &setting, result.Error
}

// SaveSetting creates or updates a setting record.
func (r *sqlRepository) SaveSetting(setting *models.AppSetting) error {
	return r.DB.Save(setting).Error
}
